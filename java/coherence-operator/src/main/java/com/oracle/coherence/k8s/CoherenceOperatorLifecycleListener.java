/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.util.concurrent.CompletableFuture;

import com.tangosol.application.Context;
import com.tangosol.application.LifecycleListener;
import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.partitionedService.PartitionedCache;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.PartitionedService;
import com.tangosol.net.events.Event;
import com.tangosol.net.events.EventDispatcher;
import com.tangosol.net.events.EventDispatcherAwareInterceptor;
import com.tangosol.net.events.partition.PartitionedServiceDispatcher;
import com.tangosol.util.Service;
import com.tangosol.util.ServiceEvent;
import com.tangosol.util.ServiceListener;

/**
 * A Coherence LifecycleListener that initializes internal Operator functionality.
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
     * The system property to enable or disable the Operator ready check resuming services.
     */
    public static final String PROP_CAN_RESUME = "coherence.k8s.operator.can.resume.services";

    /**
     * A flag that when {@code true}, allows the Operator to resume suspended services on start-up.
     */
    public static final boolean CAN_RESUME = Boolean.parseBoolean(System.getProperty(PROP_CAN_RESUME, "true"));

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
    }

    @Override
    public void preStop(Context context) {
    }

    @Override
    public void postStop(Context context) {
    }
}
