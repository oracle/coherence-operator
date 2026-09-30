/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.Map;

import javax.security.auth.Subject;
import javax.security.auth.callback.NameCallback;
import javax.security.auth.callback.PasswordCallback;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class SecretLoginModuleTest {
    @Test
    void validatesMountedCredentials() throws Exception {
        Path directory = Files.createTempDirectory("health-login-");
        Path username = directory.resolve("username");
        Path password = directory.resolve("password");
        Files.write(username, "health-operator".getBytes(StandardCharsets.UTF_8));
        Files.write(password, "secret".getBytes(StandardCharsets.UTF_8));

        SecretLoginModule module = module(username, password, "health-operator", "secret");
        assertTrue(module.login());
        assertTrue(module.commit());
        assertTrue(module.logout());

        SecretLoginModule rejected = module(username, password, "health-operator", "wrong");
        assertThrows(javax.security.auth.login.FailedLoginException.class, rejected::login);
    }

    private SecretLoginModule module(Path usernameFile, Path passwordFile, String username, String password) {
        Map<String, Object> options = new HashMap<>();
        options.put("usernameFile", usernameFile.toString());
        options.put("passwordFile", passwordFile.toString());
        SecretLoginModule module = new SecretLoginModule();
        module.initialize(new Subject(), callbacks -> {
            for (javax.security.auth.callback.Callback callback : callbacks) {
                if (callback instanceof NameCallback) {
                    ((NameCallback) callback).setName(username);
                }
                else if (callback instanceof PasswordCallback) {
                    ((PasswordCallback) callback).setPassword(password.toCharArray());
                }
            }
        }, new HashMap<>(), options);
        return module;
    }
}
