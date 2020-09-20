/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.boot.context.event.ApplicationPreparedEvent;
import org.springframework.context.ApplicationListener;

/**
 * A Spring Boot application listener that is used to initialise
 * the application when running in Spring Boot as we assume that
 * {@link Main#main(String[])} did not run.
 * If {@link Main#main(String[])} did run this listener's
 * {@link #onApplicationEvent(ApplicationPreparedEvent)}
 * method will be a no-op.
 *
 * @author Jonathan Knight  2020.09.10
 */
public class SpringBootListener
        implements ApplicationListener<ApplicationPreparedEvent> {

    private static final Log LOGGER = LogFactory.getLog(SpringBootListener.class);

    @Override
    public void onApplicationEvent(ApplicationPreparedEvent event) {
        if (Boolean.parseBoolean(System.getProperty("coherence.operator.springboot.listener", "true").toLowerCase())) {
            try {
                LOGGER.info("Initialising Coherence Operator REST endpoint");
                Main.init();
            }
            catch (Throwable t) {
                LOGGER.error("Failed to initialise the Coherence Operator", t);
            }
        }
    }
}

