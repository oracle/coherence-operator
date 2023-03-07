/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.net.InetAddress;
import java.net.UnknownHostException;

/**
 * A utility to perform a DNS lookup in Java.
 */
public class InetLookup {
    /**
     * Private constructor for utility class.
     */
    private InetLookup() {
    }

    /**
     * The entry point for running a DNS lookup.
     *
     * @param args  the array of DNS names to look up
     */
    public static void main(String[] args) {
        for (String sName: args) {
            System.out.println("Looking up: " + sName);
            try {
                InetAddress[] aAddress = InetAddress.getAllByName(sName);
                for (InetAddress address : aAddress) {
                    System.out.println(address);
                }
                System.out.println();
            }
            catch (UnknownHostException e) {
                System.out.println(e.getMessage());
            }
        }
    }
}
