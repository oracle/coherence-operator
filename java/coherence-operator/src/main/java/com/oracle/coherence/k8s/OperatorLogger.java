/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.PrintStream;
import java.util.logging.Level;
import java.util.logging.Logger;

import com.tangosol.net.CacheFactory;

/**
 * A logger for the Coherence Operator that sends messages to a
 * specific logger implementation.
 */
public interface OperatorLogger {

    /**
     * The system property to use to set the health logging should use Java logger.
     */
    String PROP_LOGGER = "coherence.k8s.operator.health.logger";

    /**
     * The {@link #PROP_LOGGER} value to log to std-err.
     */
    String LOGGER_STD_ERR = "err";

    /**
     * The {@link #PROP_LOGGER} value to log to std-out.
     */
    String LOGGER_STD_OUT = "out";

    /**
     * The {@link #PROP_LOGGER} value to log to a Java util logger.
     */
    String LOGGER_JAVA = "jdk";

    /**
     * The {@link #PROP_LOGGER} value to log to the Coherence logger.
     */
    String LOGGER_COHERENCE = "coh";

    /**
     * Log a debug message.
     *
     * @param msg   the log message
     * @param args  any arguments to apply to the log message using {@link String#format(String, Object...)}
     */
    void debug(String msg, Object... args);

    /**
     * Log an info message.
     *
     * @param msg   the log message
     * @param args  any arguments to apply to the log message using {@link String#format(String, Object...)}
     */
    void info(String msg, Object... args);

    /**
     * Log a warning message.
     *
     * @param msg   the log message
     * @param args  any arguments to apply to the log message using {@link String#format(String, Object...)}
     */
    void warn(String msg, Object... args);

    /**
     * Log an error message.
     *
     * @param msg   the log message
     * @param args  any arguments to apply to the log message using {@link String#format(String, Object...)}
     */
    void error(String msg, Object... args);

    /**
     * Log an error message.
     *
     * @param t     the excpetion to log
     * @param msg   the log message
     * @param args  any arguments to apply to the log message using {@link String#format(String, Object...)}
     */
    void error(Throwable t, String msg, Object... args);

    /**
     * Returns the {@link OperatorLogger} to use.
     *
     * @return the {@link OperatorLogger} to use
     */
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
            try {
                return new CoherenceLogger();
            }
            catch (Throwable ignored) {
                return new CacheFactoryLogger();
            }
        }
    }

    /**
     * An {@link OperatorLogger} that logs to the Coherence logger.
     */
    class CoherenceLogger implements OperatorLogger {

        public CoherenceLogger() {
            // will throw if Logger is not available (i.e. earlier Coherence version)
            com.oracle.coherence.common.base.Logger.isEnabled(com.oracle.coherence.common.base.Logger.INFO);
        }

        @Override
        public void debug(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            com.oracle.coherence.common.base.Logger.fine(String.format(msg, args));
        }

        @Override
        public void info(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            com.oracle.coherence.common.base.Logger.info(String.format(msg, args));
        }

        @Override
        public void warn(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            com.oracle.coherence.common.base.Logger.warn(String.format(msg, args));
        }

        @Override
        public void error(String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            com.oracle.coherence.common.base.Logger.err(String.format(msg, args));
        }

        @Override
        public void error(Throwable t, String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            com.oracle.coherence.common.base.Logger.err(String.format(msg, args), t);
        }
    }

    /**
     * An {@link OperatorLogger} that logs to the Coherence logger.
     */
    class CacheFactoryLogger implements OperatorLogger {
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
            CacheFactory.log(String.format(msg, args), CacheFactory.LOG_ERR);
        }

        @Override
        public void error(Throwable t, String msg, Object... args) {
            // The log method is deprecated but we need to work with earlier Coherence versions so we use it.
            //noinspection deprecation
            CacheFactory.log(String.format(msg, args), CacheFactory.LOG_ERR);
            CacheFactory.err(t);
        }
    }

    /**
     * An {@link OperatorLogger} that logs to the Java util logger.
     */
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

        @Override
        public void error(Throwable t, String msg, Object... args) {
            LOGGER.log(Level.SEVERE, String.format(msg, args), t);
        }
    }

    /**
     * An {@link OperatorLogger} that logs to a {@link PrintStream}.
     */
    class PrintLogger implements OperatorLogger {

        private final PrintStream out;

        /**
         * Create a {@link PrintLogger}.
         *
         * @param out  the {@link PrintStream} to log to
         */
        PrintLogger(PrintStream out) {
            this.out = out;
        }

        @Override
        public void debug(String msg, Object... args) {
            out.printf("[DEBUG] " + msg, args);
            out.println();
            out.flush();
        }

        @Override
        public void info(String msg, Object... args) {
            out.printf("[INFO] " + msg, args);
            out.println();
            out.flush();
        }

        @Override
        public void warn(String msg, Object... args) {
            out.printf("[WARNING] " + msg, args);
            out.println();
            out.flush();
        }

        @Override
        public void error(String msg, Object... args) {
            out.printf("[ERROR] " + msg, args);
            out.println();
            out.flush();
        }

        @Override
        public void error(Throwable t, String msg, Object... args) {
            out.printf("[ERROR] " + msg, args);
            out.println();
            t.printStackTrace(out);
            out.flush();
        }
    }
}
