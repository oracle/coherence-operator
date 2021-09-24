/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import java.lang.management.ManagementFactory;
import java.lang.management.MemoryManagerMXBean;
import java.lang.management.MemoryPoolMXBean;
import java.lang.management.MemoryType;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

import javax.management.NotCompliantMBeanException;
import javax.management.openmbean.CompositeData;
import javax.management.openmbean.CompositeDataSupport;
import javax.management.openmbean.TabularDataSupport;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.management.AnnotatedStandardMBean;
import com.tangosol.net.management.MBeanServerProxy;

/**
 * A custom MBean to track heap usage that Coherence will publish as a metric.
 *
 * @author Jonathan Knight  2020.09.01
 */
public class HeapUsage
        implements HeapUsageMBean {

    /**
     * The Garbage Collector MBean name pattern used by Coherence.
     */
    private static final String GC_PATTERN = "type=Platform,Domain=java.lang,name=%s,subType=GarbageCollector,nodeId=%d";

    /**
     * The names of the garbage collectors that should be queried for last GC info.
     */
    private final Set<String> gcNames;

    /**
     * The names of the heap memory pools used by the JVM.
     */
    private final Set<String> pools;

    /**
     * Create a {@link HeapUsage} MBean.
     *
     * @param gcNames  the names of the garbage collectors that should be queried for last GC info
     * @param pools    the names of the heap memory pools used by the JVM
     */
    public HeapUsage(Set<String> gcNames, Set<String> pools) {
        this.gcNames = gcNames;
        this.pools = pools;
    }

    /**
     * Create a {@link HeapUsage} MBean instance.
     * <p>
     * This factory method will be called by Coherence to register the MBean.
     *
     * @return the heap usage MBean
     * @throws NotCompliantMBeanException if there is an error creating the MBean
     */
    public static AnnotatedStandardMBean create() throws NotCompliantMBeanException {
        // Find the heap memory pools
        Set<String> pools = ManagementFactory.getMemoryPoolMXBeans().stream()
                .filter(m -> m.getType() == MemoryType.HEAP)
                .map(MemoryPoolMXBean::getName)
                .collect(Collectors.toSet());

        Set<String> gcNames = ManagementFactory.getGarbageCollectorMXBeans().stream()
                .map(MemoryManagerMXBean::getName)
                .collect(Collectors.toSet());

        // Create the MBean - we need to wrap our MBean in a Coherence AnnotatedStandardMBean
        // so that the metrics annotations are correctly processed.
        return new AnnotatedStandardMBean(new HeapUsage(gcNames, pools), HeapUsageMBean.class);
    }

    /**
     * Find the {@link List} of {@link CompositeDataSupport} instances for
     * the latest gc.
     *
     * @return the {@link List} of {@link CompositeDataSupport} instances for
     *         the latest gc
     */
    private List<CompositeDataSupport> getLastGcInfo() {
        List<CompositeDataSupport> lastGcInfos = new ArrayList<>();

        Cluster cluster = CacheFactory.getCluster();
        int id = cluster.getLocalMember().getId();
        TabularDataSupport usage = null;

        // find the latest GC from the collector MBeans
        MBeanServerProxy proxy = cluster.getManagement().getMBeanServerProxy();
        long nLast = 0;
        for (String gcName : gcNames) {
            String sMBeanName = String.format(GC_PATTERN, gcName, id);
            if (!proxy.isMBeanRegistered(sMBeanName)) {
                continue;
            }

            // Get the LastGcInfo attribute from the GC MBean (this may be null)
            CompositeData data = (CompositeData) proxy.getAttribute(sMBeanName, "LastGcInfo");
            if (data != null) {
                // see when it last did a GC
                Long nEnd = (Long) data.get("endTime");
                if (nEnd != null && nEnd > nLast) {
                    // this last GC is the latest that we've found so far
                    nLast = nEnd;
                    usage = (TabularDataSupport) data.get("memoryUsageAfterGc");
                }
            }
        }

        // If we found a latest gc, iterate over this to get the use after gc values for each memory pool
        if (usage != null) {
            for (String pool : pools) {
                CompositeData cd = usage.get(new Object[] {pool});
                if (cd != null) {
                    CompositeDataSupport cds = (CompositeDataSupport) cd.get("value");
                    if (cds != null) {
                        lastGcInfos.add(cds);
                    }
                }
            }
        }

        return lastGcInfos;
    }

    @Override
    public double getPercentageUsed() {
        List<CompositeDataSupport> list = getLastGcInfo();
        // sum up the heap use values for the memory pools.
        long used = list.stream()
                .map(cds -> (Number) cds.get("used"))
                .mapToLong(Number::longValue)
                .sum();
        // sum up the max heap values for the memory pools (some of these may be -1 in G1 so we ignore those).
        long max = list.stream()
                .map(cds -> (Number) cds.get("max"))
                .mapToLong(Number::longValue)
                .filter(n -> n >= 0)
                .sum();

        // although max shouldn't ever be <= 0 we don;t want a divide by zero error
        if (max <= 0) {
            return 0;
        }

        // calculate the heap use percentage rounded to two decimal places.
        return ((double) Math.round((double) used * 10000.0d / (double) max)) / 100.0d;
    }

    @Override
    public long getUsed() {
        // sum up the heap use values for the memory pools.
        return getLastGcInfo().stream()
                .map(cds -> (Number) cds.get("used"))
                .mapToLong(Number::longValue)
                .sum();
    }
}
