/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.tls;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.URL;

import com.tangosol.net.PasswordProvider;
import com.tangosol.util.Base;
import com.tangosol.util.Resources;

/**
 * A file based Coherence {@link com.tangosol.net.PasswordProvider}.
 * <p>
 * If the file name passed to the constructor is {@code null} then
 * an empty password value is returned from the {@link #get()}
 * method.
 */
public class FileBasedPasswordProvider
        implements PasswordProvider {

    private static final char[] EMPTY_PASSWORD = new char[0];

    /**
     * The name of the file containing the password.
     */
    private final String fileName;

    /**
     * Create a {@link com.oracle.coherence.examples.tls.FileBasedPasswordProvider}.
     *
     * @param file the name of the file containing the password
     */
    public FileBasedPasswordProvider(String file) {
        fileName = file;
    }

    @Override
    public char[] get() {
        return readPassword(fileName);
    }

    /**
     * Read a password from a file.
     *
     * @param fileName  the password file name
     *
     * @return the password
     */
    public static char[] readPassword(String fileName) {
        return readPassword(fileName, EMPTY_PASSWORD);
    }

    /**
     * Read a password from a file.
     *
     * @param fileName         the password file name
     * @param defaultPassword  the default password
     *
     * @return the password or the default password if the file was not found
     */
    public static char[] readPassword(String fileName, char[] defaultPassword) {
        if (fileName == null || fileName.trim().length() == 0) {
            return defaultPassword;
        }

        URL url = Resources.findFileOrResource(fileName, FileBasedPasswordProvider.class.getClassLoader());
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
