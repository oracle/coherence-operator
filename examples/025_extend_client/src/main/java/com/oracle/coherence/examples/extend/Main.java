/*
 * Copyright (c) 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.examples.extend;

import com.tangosol.net.Coherence;
import com.tangosol.net.NamedMap;
import com.tangosol.net.Session;

/**
 * A simple Coherence*Extend client test application.
 *
 * @author Jonathan Knight
 */
public class Main {
    public static void main(String[] args) {
        try {
            // Start Coherence
            Coherence coherence = Coherence.client().start().join();
            // Obtain the default Coherence Session
            Session session = coherence.getSession();
            // Obtain the test NamedMap
            NamedMap<String, String> map = session.getMap("test");

            // Put a random value into the cache
            String value = String.valueOf(Math.random());
            String key = "key-1";
            String previous = map.put(key, value);
            // Display the key, value and previous cached value
            System.out.printf("Put key=%s value=%s previous=%s", key, value, previous);

            // exit with a successful return code
            System.exit(0);
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
    }
}
