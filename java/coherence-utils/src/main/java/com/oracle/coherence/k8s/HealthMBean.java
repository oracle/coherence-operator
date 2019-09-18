/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

/**
 * A health check MBean.
 *
 * @author jk
 */
public interface HealthMBean
    {
    default String getFoo() {
    return "foo";
    }
    /**
     * Returns {@code true} if the JVM is ready.
     *
     * @return {@code true} if the JVM is ready
     */
    boolean ready();

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    boolean statusHA();

    /**
     * Returns {@code true} if the JVM is healthy.
     *
     * @return {@code true} if the JVM is healthy
     */
    boolean healthy();

    // ----- constants ------------------------------------------------------

    /**
     * The object name of the health MBean.
     */
    String HEALTH_OBJECT_NAME = "CoherenceOperator:type=Health";
    }
