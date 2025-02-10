/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing.spring;

import java.util.Arrays;
import java.util.logging.Logger;

import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.Session;

import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.core.annotation.Order;

/**
 * The test Spring storage application.
 *
 * @author Jonathan Knight  2020.09.10
 */
@SpringBootApplication
public class StorageApplication {
    private static final Logger LOGGER = Logger.getLogger(StorageApplication.class.getName());

    private static String[] arguments;

    /**
     * Run the application.
     *
     * @param args the application arguments
     */
    public static void main(String[] args) {
        arguments = args;
        SpringApplication.run(StorageApplication.class, args);
    }

    /**
     * Obtain the application arguments.
     *
     * @return the application arguments
     */
    @Bean(name = "commandLineArguments")
    public String[] commandLineArguments() {
        return arguments;
    }

    /**
     * Obtain the Coherence {@link DefaultCacheServer} starter.
     *
     * @return the Coherence {@link DefaultCacheServer} starter
     */
    @Bean
    @Order(1)
    public ApplicationRunner runCoherence() {
        return (args) -> {
            LOGGER.info("Starting DefaultCacheServer with args " + Arrays.toString(args.getSourceArgs()));
            DefaultCacheServer.main(args.getSourceArgs());
        };
    }

    /**
     * Obtain a Coherence {@link Session}.
     *
     * @return a Coherence {@link Session}
     */
    @Bean
    public Session createCoherenceSession() {
        return Session.create();
    }
}
