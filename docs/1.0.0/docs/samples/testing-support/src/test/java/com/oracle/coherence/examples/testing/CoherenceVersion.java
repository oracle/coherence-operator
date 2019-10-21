/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.tangosol.net.CacheFactory;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * A simple class that either prints out the current Coherence version
 * or if a version is passed in prints whether the current Coherence
 * version is greater than or equal to that version.
 * 
 * @author jk
 */
public class CoherenceVersion
    {
    /**
     * Print the Coherence version to standard out.
     *
     * @param asArg  the program command line arguments
     */
    public static void main(String[] asArg)
        {
        int nExitCode = 0;

        if (asArg != null && asArg.length > 0)
            {
            nExitCode = versionCheck(CacheFactory.VERSION, asArg) ? 0 : 1;
            }
        else
            {
            System.out.println(CacheFactory.VERSION);
            }
        
        System.exit(nExitCode);
        }

    static public boolean versionCheck(String sCoherence, String... asArg)
        {
        if (sCoherence.contains(":"))
            {
            sCoherence = sCoherence.substring(sCoherence.indexOf(":") + 1);
            }

        boolean fResult     = true;
        int[]   anCoherence = splitVersion(sCoherence);
        int[]   anVersion   = splitVersion(asArg[0]);
        int     cPart       = Math.min(anCoherence.length, anVersion.length);

        if (cPart > 0)
            {
            for (int i = 0; i < cPart && fResult; i++)
                {
                fResult = anCoherence[i] >= anVersion[i];
                }
            }

        return fResult;
        }

    private static int[] splitVersion(String sVersion)
        {
        Matcher matcher = pattern.matcher(sVersion);
        int[]   anPart;

        if (matcher.matches())
            {
            int   cGroup = matcher.groupCount();
            anPart = new int[cGroup];

            for (int i = 1; i <= cGroup; i++)
                {
                try
                    {
                    anPart[i-1] = Integer.parseInt(matcher.group(i));
                    }
                catch (NumberFormatException e)
                    {
                    anPart[i-1] = 0;
                    }
                }
            }
        else
            {
            anPart = new int[0];
            }

        return anPart;
        }

    private static final Pattern pattern = Pattern.compile("(\\d*)\\D*(\\d*)\\D*(\\d*)\\D*(\\d*)\\D*(\\d*)\\D*");
    }
