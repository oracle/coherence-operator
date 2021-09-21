/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.coherence.persistence.PersistenceManager;
import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.PartitionedService$PersistenceControl;
import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.PartitionedService$PersistenceControl$SnapshotController;
import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.partitionedService.PartitionedCache;
import com.tangosol.coherence.component.util.safeService.safeCacheService.SafeDistributedCacheService;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.Service;
import com.tangosol.net.WrapperService;

/**
 * A simple utility that wraps away the use of TDE classes that show in source as compile errors
 * even though they will actually compile.
 */
public final class PersistenceHelper {

    /**
     * Private constructor for final utility class.
     */
    private PersistenceHelper() {
    }

    /**
     * Returns {@code true} if the specified service is a storage enabled cache service configured
     * with persistence.
     *
     * @param service  the  service to check
     *
     * @return {@code true} if the specified service is configured with persistence
     */
    public static boolean isActivePersistenceEnabled(Service service) {
        PartitionedService$PersistenceControl persistenceControl = getPersistenceControl(service);
        PersistenceManager<?> activeManager = persistenceControl == null ? null : persistenceControl.getActiveManager();
        return activeManager != null;
    }

    /**
     * Returns {@code true} if the specified service is a storage enabled cache service configured
     * with persistence and the persistence controller is not idle.
     *
     * @param service  the  service to check
     *
     * @return {@code true} if the specified service is configured with persistence and the persistence
     *         controller is not idle
     */
    public static boolean isActive(Service service) {
        PartitionedService$PersistenceControl persistenceControl = getPersistenceControl(service);
        if (persistenceControl != null) {
            // IntelliJ underlines this code red as it thinks it will not compile, but these are TDE
            // classes and will compile fine.
            PartitionedService$PersistenceControl$SnapshotController snapshotController
                    = persistenceControl.getSnapshotController();
            return snapshotController != null && !snapshotController.isIdle();
        }
        return false;
    }

    private static PartitionedService$PersistenceControl getPersistenceControl(Service service) {
        if (service instanceof DistributedCacheService && ((DistributedCacheService) service).isLocalStorageEnabled()) {
            while (true) {
                if (service instanceof SafeDistributedCacheService) {
                    service = ((SafeDistributedCacheService) service).getService();
                }
                else if (service instanceof WrapperService) {
                    service = ((WrapperService) service).getService();
                }
                else {
                    break;
                }
            }

            if (service instanceof PartitionedCache) {
                PartitionedCache partitionedCache = (PartitionedCache) service;
                if (partitionedCache.isOwnershipEnabled()) {
                    // IntelliJ underlines this code red as it thinks it will not compile, but these are TDE
                    // classes and will compile fine.
                    return partitionedCache.getPersistenceControl();
                }
            }
        }
        return null;
    }
}
