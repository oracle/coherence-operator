/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

import com.tangosol.net.CacheFactory;

/**
 * A simple class that either prints out the current Coherence version
 * or if a version is passed in prints whether the current Coherence
 * version is greater than or equal to that version.
 */
public class CoherenceVersion {

    private static final Pattern PATTERN = Pattern.compile("(\\d*)\\D*(\\d*)\\D*(\\d*)\\D*(\\d*)\\D*(\\d*)\\D*");

    /**
     * Private constructor for utility class.
     */
    private CoherenceVersion() {
    }

    /**
     * Print the Coherence version to standard out.
     *
     * @param args the program command line arguments
     */
    public static void main(String[] args) {
        int exitCode = 0;

        if (args != null && args.length > 0) {
            exitCode = versionCheck(CacheFactory.VERSION, args) ? 0 : 1;
        }
        else {
            System.out.println(CacheFactory.VERSION);
        }

        System.exit(exitCode);
    }

    /**
     * Check the Coherence version.
     *
     * @param coherenceVersion the actual Coherence version
     * @param args             the version to validate against
     * @return {@code true} if the actual Coherence version is at least the check version
     */
    public static boolean versionCheck(String coherenceVersion, String... args) {
        if (coherenceVersion.contains(":")) {
            coherenceVersion = coherenceVersion.substring(coherenceVersion.indexOf(":") + 1);
        }

        int[] coherenceParts = splitVersion(coherenceVersion);
        int[] versionParts = splitVersion(args[0]);
        int partCount = Math.min(coherenceParts.length, versionParts.length);

        if (partCount > 0) {
            for (int i = 0; i < partCount; i++) {
                if (coherenceParts[i] == versionParts[i]) {
                    continue;
                }
                // else versions differ
                return coherenceParts[i] > versionParts[i];
            }
        }

        // versions are equal
        return true;
    }

    private static int[] splitVersion(String version) {
        Matcher matcher = PATTERN.matcher(version);
        int[] count;

        if (matcher.matches()) {
            int groupCount = matcher.groupCount();
            count = new int[groupCount];

            for (int i = 1; i <= groupCount; i++) {
                try {
                    count[i - 1] = Integer.parseInt(matcher.group(i));
                }
                catch (NumberFormatException e) {
                    count[i - 1] = 0;
                }
            }
        }
        else {
            count = new int[0];
        }

        return count;
    }
}
