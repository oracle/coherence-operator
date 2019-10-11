/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.lang.reflect.Method;

import com.tangosol.net.DefaultCacheServer;

/**
 * A main class that is used to run some initialisation code before
 * running another main class.
 */
public class Main {

    /**
     * Private constructor for utility class.
     */
    private Main() {
    }

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     * @throws java.lang.Exception if an error occurs
     */
    public static void main(String[] args) throws Exception {
        if (args.length == 0) {
            args = new String[] {DefaultCacheServer.class.getCanonicalName()};
        }

        HealthServer server = new HealthServer();
        server.start();

        String sMainClass = args[0];
        String[] asArgsReal = new String[args.length - 1];
        System.arraycopy(args, 1, asArgsReal, 0, asArgsReal.length);

        Class<?> clsMain = Class.forName(sMainClass);
        Method method = clsMain.getMethod("main", asArgsReal.getClass());
        method.invoke(null, (Object) asArgsReal);
    }
}
