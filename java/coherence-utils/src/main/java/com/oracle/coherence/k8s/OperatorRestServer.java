/*
 * Copyright (c) 2019, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.function.Supplier;
import java.util.stream.Collectors;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.Service;
import com.tangosol.net.management.MBeanHelper;
import com.tangosol.net.management.MBeanServerProxy;
import com.tangosol.net.management.Registry;
import com.tangosol.util.Filters;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

/**
 * Simple http endpoint for heath checking.
 */
public class OperatorRestServer {
    // ----- constants ------------------------------------------------------

    /**
     * The system property to use to set the health port.
     */
    public static final String PROP_HEALTH_PORT = "coherence.k8s.operator.health.port";

    /**
     * The system property to use to determine whether to wait for DCS to start.
     */
    public static final String PROP_WAIT_FOR_DCS = "coherence.k8s.operator.health.wait.dcs";

    /**
     * The path to the ready endpoint.
     */
    public static final String PATH_READY = "/ready";

    /**
     * The path to the health endpoint.
     */
    public static final String PATH_HEALTH = "/healthz";

    /**
     * The path to the HA endpoint.
     */
    public static final String PATH_HA = "/ha";

    /**
     * The path to the status endpoint.
     */
    public static final String PATH_STATUS = "/status";

    /**
     * The path to the suspend endpoint.
     */
    public static final String PATH_SUSPEND = "/suspend";

    /**
     * The path to the resume endpoint.
     */
    public static final String PATH_RESUME = "/resume";

    /**
     * The MBean name of the PartitionAssignment MBean.
     */
    public static final String MBEAN_PARTITION_ASSIGNMENT = Registry.PARTITION_ASSIGNMENT_TYPE
            + ",service=*,responsibility=DistributionCoordinator";

    /**
     * Service MBean Attributes required to compute HAStatus.
     *
     * @see #isServiceStatusHA(java.util.Map)
     */
    public static final String[] SERVICE_STATUS_HA_ATTRIBUTES =
            {
                    "HAStatus",
                    "HAStatusCode",
                    "BackupCount",
                    "ServiceNodeCount"
            };

    /**
     * The value of the Status HA attribute to signify endangered.
     */
    public static final String STATUS_ENDANGERED = "ENDANGERED";

    /**
     * The name of the HA status MBean attribute.
     */
    public static final String ATTRIB_HASTATUS = "hastatus";

    /**
     * The name of the HA status code MBean attribute.
     */
    public static final String ATTRIB_HASTATUS_CODE = "hastatuscode";

    /**
     * The name of the backup count MBean attribute.
     */
    public static final String ATTRIB_BACKUPS = "backupcount";

    /**
     * The name of the service node count MBean attribute.
     */
    public static final String ATTRIB_NODE_COUNT = "servicenodecount";

    /**
     * The error message in an exception due to there being no management member in the cluster.
     */
    public static final String NO_MANAGED_NODES = "None of the nodes are managed";

    /**
     * An empty response body.
     */
    private static final byte[] EMPTY_BODY = new byte[0];

    /**
     * System property to specify service names to be skipped in the StatusHA test.
     */
    public static final String PROP_ALLOW_ENDANGERED = "coherence.k8s.operator.statusha.allowendangered";

    // ----- data members ---------------------------------------------------

    /**
     * The http server.
     */
    private HttpServer httpServer;

    /**
     * The {@link Cluster} supplier.
     */
    private final Supplier<Cluster> clusterSupplier;

    private final Runnable waitForServiceStart;

    // ----- constructors ---------------------------------------------------

    OperatorRestServer() {
        this(CacheFactory::getCluster, OperatorRestServer::waitForDCS);
    }

    OperatorRestServer(Supplier<Cluster> supplier, Runnable waitForServiceStart) {
        this.clusterSupplier = supplier;
        this.waitForServiceStart = waitForServiceStart;
    }

    // ----- HealthServer methods ------------------------------------------------

    /**
     * Start a http server.
     *
     * @throws IOException if an error occurs
     */
    public synchronized void start() throws IOException {
        if (httpServer == null) {
            int port = Integer.getInteger(PROP_HEALTH_PORT, 6676);
            HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);

            server.createContext(PATH_READY, this::ready);
            server.createContext(PATH_HEALTH, this::health);
            server.createContext(PATH_HA, this::statusHA);
            server.createContext(PATH_STATUS, this::status);
            server.createContext(PATH_SUSPEND, this::suspend);
            server.createContext(PATH_RESUME, this::resume);

            server.setExecutor(null); // creates a default executor
            server.start();

            CacheFactory.log("Coherence Operator REST server is listening on http://localhost:"
                                     + server.getAddress().getPort(), CacheFactory.LOG_INFO);

            httpServer = server;
        }
    }

    // ----- helper methods -------------------------------------------------

    /**
     * Send a http response.
     *
     * @param t      the {@link HttpExchange} to send the response to
     * @param status the response status
     */
    private static void send(HttpExchange t, int status) {
        try {
            t.sendResponseHeaders(status, 0);
            OutputStream os = t.getResponseBody();
            os.write(EMPTY_BODY);
            os.close();
        }
        catch (IOException e) {
            e.printStackTrace();
        }
    }

    /**
     * Send a http response.
     *
     * @param t      the {@link HttpExchange} to send the response to
     * @param status the response status
     * @param body   the response body
     */
    private static void send(HttpExchange t, int status, String body) {
        try {
            byte[] bytes = body == null ? EMPTY_BODY : body.getBytes(StandardCharsets.UTF_8);
            t.sendResponseHeaders(status, bytes.length);
            OutputStream os = t.getResponseBody();
            os.write(bytes);
            os.close();
        }
        catch (IOException e) {
            e.printStackTrace();
        }
    }

    /**
     * Process a ready request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void ready(HttpExchange exchange) {
        try {
            int response = hasClusterMembers() && isStatusHA() ? 200 : 400;
            send(exchange, response);
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "Ready check");
        }
    }

    /**
     * Process a health request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void health(HttpExchange exchange) {
        try {
            int response = hasClusterMembers() ? 200 : 400;
            send(exchange, response);
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "Health check");
        }
    }

    /**
     * Process a status HA request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void statusHA(HttpExchange exchange) {
        try {
            CacheFactory.log("CoherenceOperator: StatusHA check request", CacheFactory.LOG_INFO);
            int response = isStatusHA() ? 200 : 400;
            CacheFactory.log("CoherenceOperator: StatusHA check response " + response, CacheFactory.LOG_INFO);
            send(exchange, response);
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "StatusHA check");
        }
    }

    /**
     * Process a status request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void status(HttpExchange exchange) {
        try {
            String status = getStatusName();
            send(exchange, 200, status.toLowerCase());
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "Status check");
        }
    }

    /**
     * Process a suspend request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    @SuppressWarnings("unchecked")
    void suspend(HttpExchange exchange) {
        try {
            String   path  = exchange.getRequestURI().getPath();
            String   name  = "";
            String[] parts = path.split("/");
            if (parts.length > 2) {
                name = parts[2].trim();
            }

            Cluster cluster = clusterSupplier.get();
            Registry registry = cluster.getManagement();
            MBeanServerProxy proxy = registry.getMBeanServerProxy();
            Map<Integer, String> identityMap = new HashMap<>();

            Set<String> mbeanNames = proxy.queryNames(":type=KubernetesOperator,nodeId=*", Filters.always());
            for (String mbeanName : mbeanNames) {
                Map<String, Object> attributes = proxy.getAttributes(mbeanName, Filters.always());
                identityMap.put((Integer) attributes.get(CoherenceOperatorMBean.ATTRIBUTE_NODE),
                                (String) attributes.get(CoherenceOperatorMBean.ATTRIBUTE_IDENTITY));
            }

            if (!name.isEmpty()) {
                Service service = cluster.getService(name);
                if (service == null) {
                    send(exchange, 404);
                    return;
                }
                CacheFactory.log("CoherenceOperator: Suspending service " + name, CacheFactory.LOG_WARN);
                cluster.suspendService(name);
            }
            else {
                CacheFactory.log("CoherenceOperator: Suspending all services", CacheFactory.LOG_WARN);
                Enumeration<String> names = cluster.getServiceNames();
                while (names.hasMoreElements()) {
                    name = names.nextElement();
                    Service svc = cluster.getService(name);
                    if (svc instanceof DistributedCacheService && ((DistributedCacheService) svc).isLocalStorageEnabled()) {
                        DistributedCacheService dcs = (DistributedCacheService) svc;
                        long count = dcs.getOwnershipEnabledMembers().stream()
                                        .map(m -> identityMap.get(m.getId()))
                                        .distinct()
                                        .count();
                        if (count == 1) {
                            CacheFactory.log("CoherenceOperator: Suspending service " + name, CacheFactory.LOG_INFO);
                            cluster.suspendService(name);
                        }
                        else {
                            CacheFactory.log("CoherenceOperator: Not suspending service "
                                    + name + " - is storage enabled in other deployments",
                                    CacheFactory.LOG_INFO);
                        }
                    }
                }
            }
            send(exchange, 200);
        }
        catch (Exception e) {
            CacheFactory.err(e);
            send(exchange, 500);
        }
    }

    /**
     * Process a resume request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    @SuppressWarnings("unchecked")
    void resume(HttpExchange exchange) {
        try {
            String   path  = exchange.getRequestURI().getPath();
            String   name  = "";
            String[] parts = path.split("/");
            if (parts.length > 2) {
                name = parts[2].trim();
            }

            Cluster cluster = clusterSupplier.get();
            if (!name.isEmpty()) {
                Service service = cluster.getService(name);
                if (service == null) {
                    send(exchange, 404);
                    return;
                }
                CacheFactory.log("CoherenceOperator: Resuming service " + name, CacheFactory.LOG_WARN);
                cluster.resumeService(name);
            }
            else {
                CacheFactory.log("CoherenceOperator: Resuming all services", CacheFactory.LOG_WARN);
                Enumeration<String> names = cluster.getServiceNames();
                while (names.hasMoreElements()) {
                    cluster.resumeService(names.nextElement());
                }
            }
            send(exchange, 200);
        }
        catch (Exception e) {
            CacheFactory.err(e);
            send(exchange, 500);
        }
    }

    private void handleError(HttpExchange t, Throwable thrown, String action) {
        String msg = thrown.getMessage();
        CacheFactory.log("CoherenceOperator: " + action + " failed due to '" + thrown.getMessage() + "'", CacheFactory.LOG_ERR);
        if (msg != null && msg.contains(NO_MANAGED_NODES)) {
            send(t, 400);
        }
        else {
            CacheFactory.err(thrown);
            send(t, 500);
        }
    }

    /**
     * Determine whether there are any members in the cluster.
     *
     * @return {@code true} if the Coherence cluster has members
     */
    private boolean hasClusterMembers() {
        Cluster cluster = clusterSupplier.get();
        return cluster != null && cluster.isRunning() && !cluster.getMemberSet().isEmpty();
    }

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    boolean isStatusHA() {
        String exclusions = System.getProperty(PROP_ALLOW_ENDANGERED);
        return isStatusHA(exclusions);
    }

    boolean isStatusHA(String exclusions) {
        try {
            waitForServiceStart.run();
        }
        catch (IllegalStateException e) {
            // there is probably no DCS
            CacheFactory.log("CoherenceOperator: StatusHA check failed, " + e.getMessage(), CacheFactory.LOG_ERR);
            return false;
        }

        Set<String> allowEndangered = null;
        if (exclusions != null) {
            allowEndangered = Arrays.stream(exclusions.split(","))
                .map(this::quoteMBeanName)
                .map(s -> ",service=" + s + ",")
                .collect(Collectors.toSet());
        }
System.err.println("****** allowEndangered \"" + allowEndangered + "\" exclusions=\"" + exclusions + "\"");

        Cluster cluster = clusterSupplier.get();
        if (cluster != null && cluster.isRunning()) {
            for (String mBean : getPartitionAssignmentMBeans()) {
                if (allowEndangered != null && allowEndangered.stream().anyMatch(mBean::contains)) {
                    // this service is allowed to be endangered so skip it.
                    continue;
                }
                Map<String, Object> attributes = getMBeanServiceStatusHAAttributes(mBean);
System.err.println("****** Attributes for " + mBean + " " + attributes);
                if (!isServiceStatusHA(attributes)) {
                    CacheFactory.log("CoherenceOperator: StatusHA check failed for MBean " + mBean, CacheFactory.LOG_DEBUG);
                    return false;
                }
            }
            return true;
        }
        else {
            CacheFactory.log("CoherenceOperator: StatusHA check failed - cluster is null", CacheFactory.LOG_ERR);
            return false;
        }
    }

    private String quoteMBeanName(String sMBean) {
        if (MBeanHelper.isQuoteRequired(sMBean)) {
            return MBeanHelper.quote(sMBean);
        }
        return sMBean;
    }

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    private String getStatusName() {
        Cluster cluster = clusterSupplier.get();
        int lowestStatus = Integer.MAX_VALUE;
        String status = "n/a";

        if (cluster != null && cluster.isRunning()) {
            for (String mBean : getPartitionAssignmentMBeans()) {
                Map<String, Object> attributes = getMBeanServiceStatusHAAttributes(mBean);
                // convert the attribute case as MBeanProxy or REST return them with different cases
                Map<String, Object> map = attributes.entrySet()
                        .stream()
                        .collect(Collectors.toMap(e -> e.getKey().toLowerCase(), Map.Entry::getValue));

                Integer code = (Integer) map.get(ATTRIB_HASTATUS_CODE);
                if (code != null && code < lowestStatus) {
                    status = (String) map.get(ATTRIB_HASTATUS);
                }
            }
        }

        return status;
    }

    private static void waitForDCS() {
        String s = System.getProperty(PROP_WAIT_FOR_DCS, "false");
        if (Boolean.parseBoolean(s)) {
            DefaultCacheServer dcs = DefaultCacheServer.getInstance();
            // Wait for service start to ensure that we will get back any partition cache MBeans
            dcs.waitForServiceStart();
        }
    }

    /**
     * Determine whether the Status HA state of the specified service is endangered.
     * <p>
     * If the service only has a single member then it will always be endangered but
     * this method will return {@code false}.
     *
     * @param attributes    the MBean attributes to use to determine whether the service is HA
     * @return {@code true} if the service is endangered
     */
    private boolean isServiceStatusHA(Map<String, Object> attributes) {
        boolean statusHA = true;

        // convert the attribute case as MBeanProxy or REST return them with different cases
        Map<String, Object> map = attributes.entrySet()
                .stream()
                .collect(Collectors.toMap(e -> e.getKey().toLowerCase(), Map.Entry::getValue));

        Number nodeCount = (Number) map.get(ATTRIB_NODE_COUNT);
        Number backupCount = (Number) map.get(ATTRIB_BACKUPS);
        Object status = map.get(ATTRIB_HASTATUS);

        if (nodeCount != null && nodeCount.intValue() > 1 && backupCount != null && backupCount.intValue() > 0) {
            statusHA = !Objects.equals(STATUS_ENDANGERED, status);
        }
        return statusHA;
    }

    /**
     * Obtain the {@link MBeanServerProxy} to use to query Coherence MBeans.
     *
     * @return the {@link MBeanServerProxy} to use to query Coherence MBeans
     */
    private MBeanServerProxy getMBeanServerProxy() {
        Cluster cluster = clusterSupplier.get();
        Registry registry = cluster.getManagement();

        return registry.getMBeanServerProxy();
    }

    private Set<String> getPartitionAssignmentMBeans() {
        return getMBeanServerProxy().queryNames(MBEAN_PARTITION_ASSIGNMENT, null);
    }

    /**
     * Return attribute and values needed to compute service HA Status.
     *
     * @param mBeanName the MBean name
     * @return attribute/value pairs of the specified MBean needed to compute service HAStatus.
     */
    private Map<String, Object> getMBeanServiceStatusHAAttributes(String mBeanName) {
        return getMBeanAttributes(mBeanName, SERVICE_STATUS_HA_ATTRIBUTES);
    }

    private Map<String, Object> getMBeanAttributes(String sMBean, String[] asAttributes) {
        Map<String, Object> mapAttrValue = new HashMap<>();

        for (String attribute : asAttributes) {
            mapAttrValue.put(attribute, getMBeanServerProxy().getAttribute(sMBean, attribute));
        }

        return mapAttrValue;
    }
}
