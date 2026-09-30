/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class HealthSSLProviderConfigTest {
    @Test
    void mapsManagedHealthStoreProperties() throws Exception {
        try (InputStream input = getClass().getClassLoader().getResourceAsStream("k8s-coherence-override.xml")) {
            assertNotNull(input);
            String xml = new String(readAll(input), StandardCharsets.UTF_8);
            int start = xml.indexOf("<socket-provider id=\"HealthSSLProvider\">");
            assertTrue(start >= 0, "HealthSSLProvider is missing");
            String provider = xml.substring(start, xml.indexOf("</socket-provider>", start));
            for (String expected : new String[] {
                    "<socket-provider id=\"HealthSSLProvider\">",
                    "coherence.health.security.keystore",
                    "coherence.health.security.keystore.password",
                    "coherence.health.security.keystore.type",
                    "coherence.health.security.keystore.algorithm",
                    "coherence.health.security.keystore.provider",
                    "coherence.health.security.key.password"}) {
                assertTrue(provider.contains(expected), "missing " + expected);
            }
            assertFalse(provider.contains("<trust-manager>"));
            assertFalse(provider.contains("coherence.health.security.truststore"));
            assertFalse(provider.contains("<client-auth>"));
        }
    }

    private byte[] readAll(InputStream input) throws java.io.IOException {
        java.io.ByteArrayOutputStream output = new java.io.ByteArrayOutputStream();
        byte[] buffer = new byte[4096];
        int count;
        while ((count = input.read(buffer)) >= 0) {
            output.write(buffer, 0, count);
        }
        return output.toByteArray();
    }
}
