/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.PrintStream;
import java.util.logging.Logger;

import com.tangosol.net.CacheFactory;

public interface OperatorLogger {

    /**
     * The system property to use to set the health logging should use Java logger.
     */
    String PROP_LOGGER = "coherence.k8s.operator.health.logger";

    String LOGGER_STD_ERR = "err";

    String LOGGER_STD_OUT = "out";

    String LOGGER_JAVA = "jdk";

    String LOGGER_COHERENCE = "coh";

    void debug(String msg, Object... args);

    void info(String msg, Object... args);

    void warn(String msg, Object... args);

    void error(String msg, Object... args);

    static OperatorLogger getLogger() {
        switch (System.getProperty(PROP_LOGGER, LOGGER_COHERENCE)) {
        case LOGGER_JAVA:
            return new JavaLogger();
        case LOGGER_STD_ERR:
            return new PrintLogger(System.err);
        case LOGGER_STD_OUT:
            return new PrintLogger(System.out);
        case LOGGER_COHERENCE:
        default:
            return new CoherenceLogger();
        }
    }

    class CoherenceLogger implements OperatorLogger {
        @Override
        public void debug(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            //noinspection deprecation
            CacheFactory.log(String.format(msg, args), CacheFactory.LOG_DEBUG);
        }

        @Override
        public void info(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            //noinspection deprecation
            CacheFactory.log(String.format(msg, args), CacheFactory.LOG_INFO);
        }

        @Override
        public void warn(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            //noinspection deprecation
            CacheFactory.log(String.format(msg, args), CacheFactory.LOG_WARN);
        }

        @Override
        public void error(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            //noinspection deprecation
            CacheFactory.log(String.format(msg, args), CacheFactory.LOG_WARN);
        }
    }

    class JavaLogger implements OperatorLogger {
        private static final Logger LOGGER = Logger.getLogger(OperatorLogger.class.getName());

        @Override
        public void debug(String msg, Object... args) {
            LOGGER.fine(String.format(msg, args));
        }

        @Override
        public void info(String msg, Object... args) {
            LOGGER.info(String.format(msg, args));
        }

        @Override
        public void warn(String msg, Object... args) {
            LOGGER.warning(String.format(msg, args));
        }

        @Override
        public void error(String msg, Object... args) {
            LOGGER.severe(String.format(msg, args));
        }
    }

    class PrintLogger implements OperatorLogger {
        private static final Logger LOGGER = Logger.getLogger(OperatorLogger.class.getName());

        private final PrintStream out;

        public PrintLogger(PrintStream out) {
            this.out = out;
        }

        @Override
        public void debug(String msg, Object... args) {
            out.printf("[DEBUG] " + msg, args);
        }

        @Override
        public void info(String msg, Object... args) {
            out.printf("[INFO] " + msg, args);
        }

        @Override
        public void warn(String msg, Object... args) {
            out.printf("[WARNING] " + msg, args);
        }

        @Override
        public void error(String msg, Object... args) {
            out.printf("[ERROR] " + msg, args);
        }
    }
}
