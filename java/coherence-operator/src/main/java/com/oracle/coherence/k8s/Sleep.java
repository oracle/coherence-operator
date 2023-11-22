/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import java.util.Properties;
import java.util.Set;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;

/**
 * A main class that sleeps for a specific duration.
 */
public class Sleep {
    /**
     * The default number of milliseconds to sleep.
     */
    public static final long DEFAULT_SLEEP = 60000;

    /**
     * Private constructor for utility class.
     */
    private Sleep() {
    }

    /**
     * Sleep main entry point.
     *
     * @param args  the arguments
     *
     * @throws Exception if an error occurs
     */
    public static void main(String[] args) throws Exception {

        Properties props = new Properties();
        try (ReadinessServer rest = new ReadinessServer()) {
            rest.start();

            long cMillis;
            if (args.length == 0) {
                cMillis = DEFAULT_SLEEP;
            }
            else {
                cMillis = new Duration(args[0]).asJavaDuration().toMillis();
            }

            cMillis = Math.max(10000, cMillis);

            long cSeconds = cMillis / 1000;
            CacheFactory.log("Sleeping for " + cSeconds + " seconds");
            Thread.sleep(cMillis);
        }
    }

    /**
     * A simple readiness probe server.
     */
    private static class ReadinessServer
            extends OperatorRestServer {

        /**
         * Create a {@link ReadinessServer}.
         */
        private ReadinessServer() {
            super(() -> null, () -> {}, new Properties());
        }

        @Override
        protected boolean hasClusterMembers() {
            return true;
        }

        @Override
        boolean isStatusHA() {
            return true;
        }

        @Override
        boolean isStatusHA(String exclusions) {
            return true;
        }

        @Override
        boolean areCacheServicesHA(Cluster cluster, Set<String> allowEndangered) {
            return true;
        }

        @Override
        boolean isPersistenceIdle() {
            return true;
        }
    }
}
