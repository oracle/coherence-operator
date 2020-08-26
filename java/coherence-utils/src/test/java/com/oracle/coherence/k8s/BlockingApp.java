/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

/**
 * An application that does nothing but wait.
 * <p>
 * This can be started as an application by Bedrock
 * so that we have a JVM in a plain state that
 * {@link Runnable}s and {@link java.util.concurrent.Callable}s
 * can be invoke against.
 */
public class BlockingApp {

    private static final Object BLOCKER = new Object();

    public static void main(String[] args) throws Exception {
        synchronized (BLOCKER) {
            BLOCKER.wait();
        }
    }
}
