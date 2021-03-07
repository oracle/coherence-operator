/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.k8s.client;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.URL;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.PasswordProvider;
import com.tangosol.util.Base;
import com.tangosol.util.Resources;

/**
 * A file based Coherence {@link com.tangosol.net.PasswordProvider}.
 * <p>
 * If the file name passed to the constructor is {@code null} then
 * a empty password value is returned from the {@link #get()}
 * method.
 */
public class FileBasedPasswordProvider
        implements PasswordProvider {

    /**
     * The name of the file containing the password.
     */
    private final String fileName;

    /**
     * Create a {@link FileBasedPasswordProvider}.
     *
     * @param file the name of the file containing the password
     */
    public FileBasedPasswordProvider(String file) {
        fileName = file;
    }

    @Override
    public char[] get() {

        if (fileName == null || fileName.trim().length() == 0) {
            return new char[0];
        }

        URL url = Resources.findFileOrResource(fileName, getClass().getClassLoader());
        if (url == null) {
            throw new IllegalStateException("Could not find password file " + fileName);
        }

        try (InputStream in = url.openStream()) {
            BufferedReader reader = new BufferedReader(new InputStreamReader(in));
            String line = reader.readLine();
            return line == null ? new char[0] : line.toCharArray();
        }
        catch (IOException e) {
            throw Base.ensureRuntimeException(e);
        }
    }
}
