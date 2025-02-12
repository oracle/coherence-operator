/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing.spring;

import java.util.Map;
import java.util.Properties;
import java.util.TreeMap;

import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.NamedCache;
import com.tangosol.net.Session;
import com.tangosol.net.partition.SimplePartitionKey;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

/**
 * A test Spring controller.
 *
 * @author Jonathan Knight  2020.09.15
 */
@RestController
public class ApplicationController {

    /**
     * The name of the canary cache.
     */
    public static final String CACHE_NAME_CANARY = "canary";

    @Autowired
    private Session session;

    @Autowired
    @Qualifier("commandLineArguments")
    private String[] commandLineArguments;

    /**
     * Obtain the program arguments.
     *
     * @return the program arguments
     */
    @RequestMapping("/args")
    public String[] args() {
        return commandLineArguments;
    }

    /**
     * Perform the canary test.
     *
     * @return the test results.
     */
    @RequestMapping("/canaryCheck")
    public String canaryCheck() {
        NamedCache<SimplePartitionKey, String> cache = session.getCache(CACHE_NAME_CANARY);
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();
        int nSize = cache.size();

        if (nSize != nPart) {
            throw new CanaryFailure(nSize, nPart);
        }
        return "OK " + nSize + " entries";
    }

    /**
     * Clear the canary cache.
     *
     * @return  the request response
     */
    @RequestMapping("/canaryClear")
    public String canaryClear() {
        NamedCache<SimplePartitionKey, String> cache = session.getCache(CACHE_NAME_CANARY);
        cache.truncate();
        return "OK";
    }

    /**
     * Initialise the canary cache.
     *
     * @return the request response
     */
    @RequestMapping("/canaryStart")
    public String canaryStart() {
        NamedCache<SimplePartitionKey, String> cache = session.getCache(CACHE_NAME_CANARY);
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();

        for (int i = 0; i < nPart; i++) {
            SimplePartitionKey key = SimplePartitionKey.getPartitionKey(i);
            cache.put(key, "data");
        }

        return "OK";
    }

    /**
     * Obtain the environment variables.
     *
     * @return the environment variables
     */
    @RequestMapping("/env")
    public Map<String, String> env() {
        return new TreeMap<>(System.getenv());
    }

    /**
     * Obtain the Spring Boot message.
     *
     * @return the Spring Boot message
     */
    @RequestMapping("/")
    public String index() {
        return "Greetings from Spring Boot!";
    }

    /**
     * Obtain the system properties.
     *
     * @return the system properties
     */
    @RequestMapping("/props")
    public Map<String, String> props() {
        Map<String, String> map = new TreeMap<>();
        Properties props = System.getProperties();
        for (String name : props.stringPropertyNames()) {
            map.put(name, props.getProperty(name));
        }
        return map;
    }

    /**
     * A canary test failure exception.
     */
    @ResponseStatus(value = HttpStatus.BAD_REQUEST)
    public static class CanaryFailure
            extends RuntimeException {
        /**
         * Create a canary test failure exception.
         *
         * @param actual   the actual cache size
         * @param expected the expected cache size
         */
        public CanaryFailure(int actual, int expected) {
            super("Canary check failed. Expected " + expected + " entries but there are only " + actual);
        }
    }
}
