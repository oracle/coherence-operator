/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import java.net.URLClassLoader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;

class HealthSecuritySupportTest {
    @Test
    void recognizesPinnedHardenedCoherenceHandler() throws Exception {
        String configured = System.getProperty("health.mutator.coherence.jar", "");
        Assumptions.assumeTrue(!configured.isEmpty(), "pinned hardened Coherence jar is required by the dedicated CI lane");
        Path jar = Paths.get(configured);
        Assumptions.assumeTrue(Files.isRegularFile(jar), "configured hardened Coherence jar does not exist");
        ClassLoader platformLoader = ClassLoader.getSystemClassLoader().getParent();
        try (URLClassLoader loader = new URLClassLoader(new java.net.URL[] {jar.toUri().toURL()}, platformLoader)) {
            Class<?> handler = Class.forName("com.tangosol.internal.health.HealthHttpHandler", false, loader);
            assertEquals(jar.toUri(), handler.getProtectionDomain().getCodeSource().getLocation().toURI());
            assertDoesNotThrow(() -> HealthSecuritySupport.verify(loader));
        }
    }
}
