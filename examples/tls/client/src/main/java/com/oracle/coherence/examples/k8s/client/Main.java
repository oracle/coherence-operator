/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.k8s.client;

import com.tangosol.net.Coherence;
import com.tangosol.net.NamedMap;
import com.tangosol.net.Session;

public class Main {

    public static void main(String[] args) throws Exception {
        try {
            Coherence coherence = Coherence.client();
            coherence.start().join();

            Session session = coherence.getSession();
            NamedMap<String, String> cache = session.getMap("test");

            for (int i=0; i<10000; i++) {
                String old = cache.get("key-" + i);
                cache.put("key-" + i, "value-" + i);
                System.out.println("Put " + i);
                try {
                    Thread.sleep(1000);
                }
                catch (InterruptedException e) {
                    break;
                }
            }
        }
        catch (Throwable thrown) {
            thrown.printStackTrace();
        }

        // wait forever so that the container never stops, even after an exception.
        synchronized (Main.class) {
            Main.class.wait();
        }
    }
}
