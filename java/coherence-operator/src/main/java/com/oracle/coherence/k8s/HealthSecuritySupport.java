/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

/** Checks Coherence health security capability without calling an administrative endpoint. */
public final class HealthSecuritySupport {
    private HealthSecuritySupport() {
    }

    /**
     * Verify the selected implementation exposes the supported security contract.
     *
     * @param args ignored
     * @throws Exception if capability cannot be established
     */
    public static void main(String[] args) throws Exception {
        verify(Thread.currentThread().getContextClassLoader());
    }

    /**
     * Verify the hardened Coherence handler visible to the supplied application class loader.
     *
     * @param loader the application class loader
     * @throws Exception if the required contract is absent
     */
    static void verify(ClassLoader loader) throws Exception {
        Class<?> handler = Class.forName("com.tangosol.internal.health.HealthHttpHandler", false, loader);
        if (!"coherence.health.http.mutators.enabled".equals(handler.getField("PROP_MUTATORS_ENABLED").get(null))) {
            throw new IllegalStateException("Unsupported health server");
        }
        if (!"coherence.health.http.auth".equals(handler.getField("PROP_AUTH").get(null))) {
            throw new IllegalStateException("Unsupported health authentication");
        }
        handler.getDeclaredMethod("ensureMutatorAccess",
                Class.forName("com.tangosol.internal.http.HttpRequest", false, loader));
    }
}
