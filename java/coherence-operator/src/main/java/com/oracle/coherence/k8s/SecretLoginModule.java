/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.Arrays;
import java.util.Map;

import javax.security.auth.Subject;
import javax.security.auth.callback.Callback;
import javax.security.auth.callback.CallbackHandler;
import javax.security.auth.callback.NameCallback;
import javax.security.auth.callback.PasswordCallback;
import javax.security.auth.login.FailedLoginException;
import javax.security.auth.login.LoginException;
import javax.security.auth.spi.LoginModule;

/**
 * Validates Basic credentials against customer-managed, mounted Secret files.
 * Files are read for each login; supported rotation uses immutable Secret names.
 */
public class SecretLoginModule implements LoginModule {
    private CallbackHandler handler;
    private Map<String, ?> options;
    private boolean authenticated;

    @Override
    public void initialize(Subject subject, CallbackHandler callbacks,
                           Map<String, ?> sharedState, Map<String, ?> configuration) {
        handler = callbacks;
        options = configuration;
    }

    @Override
    public boolean login() throws LoginException {
        NameCallback name = new NameCallback("username");
        PasswordCallback password = new PasswordCallback("password", false);
        byte[] expectedUser = null;
        byte[] expectedPassword = null;
        byte[] suppliedPassword = null;
        char[] chars = null;
        try {
            handler.handle(new Callback[] {name, password});
            expectedUser = read("usernameFile");
            expectedPassword = read("passwordFile");
            chars = password.getPassword();
            if (name.getName() == null || chars == null) {
                throw new FailedLoginException("Invalid credentials");
            }
            suppliedPassword = new String(chars).getBytes(StandardCharsets.UTF_8);
            boolean userMatches = MessageDigest.isEqual(expectedUser, name.getName().getBytes(StandardCharsets.UTF_8));
            boolean passwordMatches = MessageDigest.isEqual(expectedPassword, suppliedPassword);
            authenticated = userMatches & passwordMatches;
            if (!authenticated) {
                throw new FailedLoginException("Invalid credentials");
            }
            return true;
        }
        catch (Exception failure) {
            authenticated = false;
            throw new FailedLoginException("Invalid credentials");
        }
        finally {
            password.clearPassword();
            if (chars != null) {
                Arrays.fill(chars, '\0');
            }
            for (byte[] value : new byte[][] {expectedUser, expectedPassword, suppliedPassword}) {
                if (value != null) {
                    Arrays.fill(value, (byte) 0);
                }
            }
        }
    }

    private byte[] read(String option) throws Exception {
        Object path = options.get(option);
        if (!(path instanceof String)) {
            throw new LoginException("Credential file is required");
        }
        if (Files.size(Paths.get((String) path)) > 8192) {
            throw new LoginException("Invalid credential file");
        }
        byte[] value = Files.readAllBytes(Paths.get((String) path));
        if (value.length == 0) {
            throw new LoginException("Invalid credential file");
        }
        return value;
    }

    @Override
    public boolean commit() {
        return authenticated;
    }

    @Override
    public boolean abort() {
        authenticated = false;
        return true;
    }

    @Override
    public boolean logout() {
        authenticated = false;
        return true;
    }
}
