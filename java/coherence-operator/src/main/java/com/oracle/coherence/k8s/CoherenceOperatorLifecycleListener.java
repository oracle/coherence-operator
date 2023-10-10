/*
 * Copyright (c) 2021, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.File;
import java.io.PrintWriter;
import java.util.Base64;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import com.tangosol.application.Context;
import com.tangosol.application.LifecycleListener;
import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.partitionedService.PartitionedCache;
import com.tangosol.coherence.config.Config;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.Member;
import com.tangosol.net.PartitionedService;
import com.tangosol.net.events.Event;
import com.tangosol.net.events.EventDispatcher;
import com.tangosol.net.events.EventDispatcherAwareInterceptor;
import com.tangosol.net.events.partition.PartitionedServiceDispatcher;
import com.tangosol.util.Service;
import com.tangosol.util.ServiceEvent;
import com.tangosol.util.ServiceListener;

/**
 * A Coherence {@link com.tangosol.net.DefaultCacheServer} {@link LifecycleListener} that
 * initializes internal Operator functionality.
 *
 * @author Jonathan Knight  2021.09.21
 */
@SuppressWarnings("rawtypes")
public class CoherenceOperatorLifecycleListener
        implements LifecycleListener, EventDispatcherAwareInterceptor, ServiceListener {

    /**
     * The operator logger to use.
     */
    private static final OperatorLogger LOGGER = OperatorLogger.getLogger();

    /**
     * The system property to enable or disable the Operator resuming services.
     */
    public static final String PROP_CAN_RESUME = "coherence.k8s.operator.can.resume.services";

    /**
     * The system property to enable or disable the Operator resuming individual services.
     */
    public static final String PROP_RESUME_SERVICES = "coherence.k8s.operator.resume.services";

    /**
     * A flag that when {@code true}, allows the Operator to resume suspended services on start-up.
     */
    public static final boolean CAN_RESUME;

    /**
     * The default flag indicating whether to resume a service.
     */
    public static final boolean DEFAULT_RESUME_SERVICE;

    /**
     * An optional map of service names and whether they can be resumed on start-up.
     */
    public static final Map<String, Boolean> SERVICE_RESUME_MAP;

    /*
     * Initialise the can resume flag and resume services map.
     */
    static {
        boolean resumeServices = Boolean.parseBoolean(System.getProperty(PROP_CAN_RESUME, "true"));
        DEFAULT_RESUME_SERVICE = resumeServices;
        Map<String, Boolean> map = getResumeMap();
        if (map != null && !map.isEmpty()) {
            CAN_RESUME = true;
            SERVICE_RESUME_MAP = map;
        }
        else {
            CAN_RESUME = resumeServices;
            SERVICE_RESUME_MAP = null;
        }
    }

    @Override
    public void introduceEventDispatcher(String s, EventDispatcher eventDispatcher) {
        if (CAN_RESUME && eventDispatcher instanceof PartitionedServiceDispatcher) {
            PartitionedService service = ((PartitionedServiceDispatcher) eventDispatcher).getService();
            if (service instanceof DistributedCacheService && ((DistributedCacheService) service).isLocalStorageEnabled()) {
                if (service.isRunning()) {
                    ensureResumed(service);
                }
                else {
                    service.addServiceListener(this);
                }
            }
        }
    }

    @Override
    public void onEvent(Event event) {
    }

    @Override
    public void serviceStarting(ServiceEvent serviceEvent) {
    }

    @Override
    public void serviceStarted(ServiceEvent serviceEvent) {
        if (CAN_RESUME) {
            ensureResumed(serviceEvent.getService());
        }
    }

    private void ensureResumed(Service service) {
        if (service instanceof PartitionedCache && ((PartitionedCache) service).isSuspended()) {
            String serviceName = ((PartitionedCache) service).getServiceName();
            if (SERVICE_RESUME_MAP == null
                || SERVICE_RESUME_MAP.getOrDefault(serviceName, DEFAULT_RESUME_SERVICE) == Boolean.TRUE) {
                // We need to resume the service on another thread so that we do not block start-up,
                // in this case we'll just use the fork-join pool.
                CompletableFuture.runAsync(() -> {
                    LOGGER.info("CoherenceOperator: is automatically resuming suspended service %s", serviceName);
                    ((PartitionedCache) service).getCluster().resumeService(serviceName);
                }).handle((ignored, err) -> {
                    if (err != null) {
                        LOGGER.error(err, "CoherenceOperator: failed to resume service %s", serviceName);
                    }
                    return null;
                });
            }
            else {
                LOGGER.info("CoherenceOperator: not resuming service %s as it is in the exclusion list", serviceName);
            }
        }
    }

    @Override
    public void serviceStopping(ServiceEvent serviceEvent) {
    }

    @Override
    public void serviceStopped(ServiceEvent serviceEvent) {
    }

    @Override
    public void preStart(Context context) {
        try {
            LOGGER.info("Ensuring initialisation of Coherence Operator");
            Main.init();
            context.getConfigurableCacheFactory().getInterceptorRegistry().registerEventInterceptor(this);
        }
        catch (Throwable t) {
            LOGGER.error("Failed to initialise the Coherence Operator", t);
        }
    }

    @Override
    public void postStart(Context context) {
        initCohCtl();
    }

    @Override
    public void preStop(Context context) {
    }

    @Override
    public void postStop(Context context) {
    }

    /**
     * Load a map of services and whether they can be resumed from the
     * {@link #PROP_RESUME_SERVICES} system property.
     *
     * @return a map of service names and whether they can be resumed
     */
    static Map<String, Boolean> getResumeMap() {
        return getResumeMap(System.getProperty(PROP_RESUME_SERVICES));
    }

    /**
     * Load a map of services and whether they can be resumed.
     *
     * @param services the string of service names
     *
     * @return a map of service names and whether they can be resumed
     */
    static Map<String, Boolean> getResumeMap(String services) {
        try {
            if (services == null) {
                return null;
            }

            services = services.trim();
            if (services.isEmpty()) {
                return null;
            }

            if (services.startsWith("base64:")) {
                services = new String(Base64.getDecoder().decode(services.substring(7)));
            }

            Map<String, Boolean> map = new HashMap<>();
            StringBuilder serviceName = new StringBuilder();
            StringBuilder enabled = new StringBuilder();
            StringBuilder builder = serviceName;
            boolean inQuotes = false;

            for (int i = 0; i < services.length(); i++) {
                char ch = services.charAt(i);
                switch (ch) {
                    case '"':
                        inQuotes = !inQuotes;
                        break;
                    case '\\':
                        int next = i + 1;
                        if (next < services.length() && services.charAt(next) == '"') {
                            builder.append('"');
                            i++;
                        }
                        else {
                            builder.append('\\');
                        }
                        break;
                    case '=':
                        if (inQuotes) {
                            builder.append('=');
                        }
                        else {
                            builder = enabled;
                        }
                        break;
                    case ',':
                        if (inQuotes) {
                            builder.append(',');
                        }
                        else {
                            String s = serviceName.toString();
                            boolean resume = Boolean.parseBoolean(enabled.toString());
                            if (!s.trim().isEmpty()) {
                                map.put(s, resume);
                            }
                            serviceName = new StringBuilder();
                            enabled = new StringBuilder();
                            builder = serviceName;
                        }
                        break;
                    default:
                        builder.append(ch);
                        break;
                }
            }

            if (serviceName.length() > 0) {
                String s = serviceName.toString();
                boolean resume = Boolean.parseBoolean(enabled.toString());
                if (!s.trim().isEmpty()) {
                    map.put(s, resume);
                }
            }
            return Collections.unmodifiableMap(map);
        }
        catch (Throwable t) {
            LOGGER.error(t, "CoherenceOperator: Error decoding service resume list %s", services);
            return null;
        }
    }

    void initCohCtl() {
        try {
            Cluster cluster        = CacheFactory.getCluster();
            Member  member         = cluster.getLocalMember();
            String  clusterName    = member.getClusterName();
            String  port           = Config.getProperty("coherence.management.http.port", "30000");
            String  provider       = Config.getProperty("coherence.management.http.provider");
            String defaultProtocol = provider == null || provider.isEmpty() ? "http" : "https";
            String protocol        = Config.getProperty("coherence.operator.cli.protocol", defaultProtocol);
            String home            = System.getProperty("user.home");
            String connectionType  = "http";

            File cohctlHome = new File(home + File.separator + ".cohctl");
            File configFile = new File(cohctlHome, "cohctl.yaml");

            if (!configFile.exists()) {
                LOGGER.info("CoherenceOperator: creating default cohctl config at " + configFile.getAbsolutePath());
                if (!cohctlHome.exists()) {
                    cohctlHome.mkdirs();
                }
                try (PrintWriter out = new PrintWriter(configFile)) {
                    out.println("clusters:");
                    out.println("    - name: default");
                    out.println("      discoverytype: manual");
                    out.println("      connectiontype: " + connectionType);
                    out.println("      connectionurl: " + protocol + "://127.0.0.1:" + port + "/management/coherence/cluster");
                    out.println("      nameservicediscovery: \"\"");
                    out.println("      clusterversion: \"" + CacheFactory.VERSION + "\"");
                    out.println("      clustername: \"" + clusterName + "\"");
                    out.println("      clustertype: Standalone");
                    out.println("      manuallycreated: false");
                    out.println("      baseclasspath: \"\"");
                    out.println("      additionalclasspath: \"\"");
                    out.println("      arguments: \"\"");
                    out.println("      managementport: 0");
                    out.println("      persistencemode: \"\"");
                    out.println("      loggingdestination: \"\"");
                    out.println("      managementavailable: false");
                    out.println("color: \"on\"");
                    out.println("currentcontext: default");
                    out.println("debug: false");
                    out.println("defaultbytesformat: m");
                    out.println("ignoreinvalidcerts: false");
                    out.println("requesttimeout: 30");
                }
            }
        }
        catch (Exception e) {
            LOGGER.error(e, "Coherence Operator: Failed to create default cohctl config");
        }
    }
}
