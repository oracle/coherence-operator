/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.net.management.annotation.Description;
import com.tangosol.net.management.annotation.MetricsValue;

/**
 * A custom MBean interface to track heap usage that Coherence will publish as a metric.
 *
 * @author Jonathan Knight  2020.09.01
 */
public interface HeapUsageMBean {

    /**
     * Obtain the amount of heap in use after the last GC.
     *
     * @return the amount of heap (in bytes) in use after the last GC
     */
    @MetricsValue
    @Description("The total heap in use after the last GC")
    long getUsed();

    /**
     * Obtain the amount of heap in use after the last GC as a percentage of the maximum heap.
     *
     * @return the amount of heap in use after the last GC as a percentage of the maximum heap
     */
    @MetricsValue
    @Description("The total heap in use as after the last GC a percentage of the maximum heap")
    double getPercentageUsed();
}
