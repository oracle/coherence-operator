/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DefaultCacheServer;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.lang.reflect.Method;
import java.util.function.Supplier;

/**
 * A main class that is used to run some initialisation code before
 * running another main class.
 *
 * @author jk
 */
public class Main
    {
    // ----- Main methods ------------------------------------------------

    /**
     * Program entry point.
     *
     * @param asArgs the program command line arguments
     */
    public static void main(String[] asArgs) throws Exception
        {
        if (asArgs.length == 0)
            {
            asArgs = new String[]{DefaultCacheServer.class.getCanonicalName()};
            }

        registerHealthMBean();

        String   sMainClass = asArgs[0];
        String[] asArgsReal = new String[asArgs.length - 1];
        System.arraycopy(asArgs, 1, asArgsReal, 0, asArgsReal.length);

        Class<?> clsMain = Class.forName(sMainClass);
        Method   method  = clsMain.getMethod("main", asArgsReal.getClass());
        method.invoke(null, (Object) asArgsReal);
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Register the health MBean.
     *
     * @throws Exception if the MBean cannot be registered.
     */
    static void registerHealthMBean() throws Exception
        {
        MBeanServer server = ManagementFactory.getPlatformMBeanServer();
        Health      health = new Health();
        server.registerMBean(health,  new ObjectName(HealthObjectName));
        }

    // ----- inner interface HealthMBean ------------------------------------

    /**
     * A health check MBean.
     */
    public interface HealthMBean
        {
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
        }

    // ----- inner interface Health -----------------------------------------

    /**
     * The health MBean implementation.
     */
    public static class Health implements HealthMBean
        {
        // ----- constructors -----------------------------------------------

        Health()
            {
            this(CacheFactory::ensureCluster);
            }

        Health(Supplier<Cluster> supplier)
            {
            f_supplier = supplier;
            f_probe    = new ClusterMemberProbe(supplier);
            }

        // ----- HealthMBean methods ----------------------------------------

        @Override
        public boolean ready()
            {
            return hasClusterMembers() && statusHA();
            }

        @Override
        public boolean statusHA()
            {
            return f_probe.isStatusHA();
            }

        @Override
        public boolean healthy()
            {
            return hasClusterMembers();
            }

        // ----- helper methods ---------------------------------------------

        /**
         * Determine whether there are any members in the cluster.
         *
         * @return  {@code true} if the Coherence cluster has members
         */
        private boolean hasClusterMembers()
            {
            return !f_supplier.get().getMemberSet().isEmpty();
            }

        // ----- data members -----------------------------------------------

        /**
         * The {@link Cluster} supplier.
         */
        private final Supplier<Cluster> f_supplier;

        /**
         * The probe to use to obtain StatusHA.
         */
        private final ClusterMemberProbe f_probe;
        }

    // ----- constants ------------------------------------------------------

    /**
     * The object name of the health MBean.
     */
    public static final String HealthObjectName = "CoherenceOperator:type=Health";
    }
