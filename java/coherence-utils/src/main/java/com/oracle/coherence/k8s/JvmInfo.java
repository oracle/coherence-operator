/*
 * Copyright (c) 2019, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

/**
 * Determine the JVM information.
 */
public class JvmInfo {

    private JvmInfo() {
    }

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     */
    public static void main(String[] args) {
        if (args.length > 0) {
            System.out.println(System.getProperty(args[0]));
        }
    }
}
