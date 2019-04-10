/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.PasswordProvider;
import com.tangosol.util.Base;
import com.tangosol.util.Resources;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;

/**
 * A file based Coherence {@link PasswordProvider}.
 * <p>
 * If the file name passed to the constructor is {@code null} then
 * a {@code null} password value is returned from the {@link #get()}
 * method.
 *
 * @author jk
 */
public class FileBasedPasswordProvider
        implements PasswordProvider
    {
    // ----- constructors ---------------------------------------------------

    /**
     * Create a {@link FileBasedPasswordProvider}.
     *
     * @param sFile  the name of the file containing the password
     */
    public FileBasedPasswordProvider(String sFile)
        {
        m_sFile = sFile;
        }

    // ----- PasswordProvider methods ---------------------------------------

    @Override
    public char[] get()
        {

        if (m_sFile == null || m_sFile.trim().length() == 0)
            {
            return null;
            }

        try (InputStream in = Resources.findFileOrResource(m_sFile, getClass().getClassLoader()).openStream())
            {
            BufferedReader reader = new BufferedReader(new InputStreamReader(in));
            return reader.readLine().toCharArray();
            }
        catch (IOException e)
            {
            throw Base.ensureRuntimeException(e);
            }
        }

    // ----- data members ------------------------------------------------

    /**
     * The name of the file containing the password.
     */
    private String m_sFile;
    }
