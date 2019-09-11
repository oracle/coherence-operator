/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

/**
 * @author jk
 */
public interface Probe
        extends AutoCloseable
    {
    /**
     * Determine whether this probe can be used.
     *
     * @return  {@code true} if this probe can be used.
     */
    boolean isAvailable();

    /**
     * Perform the readiness test.
     *
     * @return  {@code true} if the readiness test passes.
     */
    boolean isReady();

    /**
     * Perform the liveness test.
     *
     * @return  {@code true} if the liveness test passes.
     */
    boolean isLive();

    /**
     * Determine whether the services are all HA.
     *
     * @return  {@code true} if the status is HA
     */
    boolean isStatusHA();

    @Override
    void close();
    }
