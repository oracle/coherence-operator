/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.NamedCache;

/**
 * A simple Coherence Extend client.
 */
public class ExtendClient {

    /**
     * Private constructor for utility class.
     */
    private ExtendClient() {
    }

    /**
     * Run the Extend client.
     *
     * @param args  the program arguments.
     */
    public static void main(String[] args) {
        try {
            System.out.println("Getting cache 'test' from Extend client session");
            NamedCache<String, String> cache = CacheFactory.getCache("test");
            System.out.println("Putting key and value into cache 'test'");
            cache.put("key-1", "value-1");

            System.out.println("Test completed successfully");
            System.exit(0);
        }
        catch (Throwable e) {
            e.printStackTrace();
            System.out.println("Test failed");
            System.exit(1);
        }
    }
}
