/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing.spring;

import com.tangosol.util.MapListener;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.context.annotation.Bean;

/**
 * The test Spring storage application.
 *
 * @author Jonathan Knight  2020.09.10
 */
@SpringBootApplication
@EnableCaching
public class StorageApplication {
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
     * Return an instance of a {@link com.tangosol.util.MapListener} bean.
     *
     * @param <K>  the cache key type
     * @param <V>  the cache value type
     *
     * @return an instance of a {@link com.tangosol.util.MapListener} bean
     */
    @Bean(name = "mapListenerBean")
    public <K, V> MapListener<K, V> createListener() {
        return new MapListenerBean<>();
    }
}
