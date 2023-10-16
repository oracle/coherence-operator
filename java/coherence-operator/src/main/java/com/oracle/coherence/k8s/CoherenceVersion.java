/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.coherence.component.net.memberSet.actualMemberSet.ServiceMemberSet;
import com.tangosol.net.CacheFactory;

/**
 * A simple class that either prints out the current Coherence version
 * or if a version is passed in prints whether the current Coherence
 * version is greater than or equal to that version.
 */
public class CoherenceVersion {
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

        String version = System.getenv().getOrDefault("COH_VERSION_CHECK", CacheFactory.VERSION);

        if (args != null && args.length > 0) {
            exitCode = versionCheck(version, args) ? 0 : 1;
        }
        else {
            System.out.println(version);
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
        System.out.print("CoherenceOperator: version check actual=\"" + coherenceVersion + "\" required=\"" + args[0] + '"');
        if (coherenceVersion.contains(" ")) {
            coherenceVersion = coherenceVersion.substring(0, coherenceVersion.indexOf(" "));
        }
        if (coherenceVersion.contains(":")) {
            coherenceVersion = coherenceVersion.substring(coherenceVersion.indexOf(":") + 1);
        }

        int[] nCoherenceParts = ServiceMemberSet.toVersionArray(coherenceVersion);
        int nActual;
        if (nCoherenceParts[0] > 20) {
            nActual = ServiceMemberSet.encodeVersion(nCoherenceParts[0], nCoherenceParts[1], nCoherenceParts[2]);
        }
        else {
            nActual = ServiceMemberSet.parseVersion(coherenceVersion);
        }

        int[] nParts = ServiceMemberSet.toVersionArray(args[0]);
        int nRequired;
        if (nParts[0] > 20) {
            nRequired = ServiceMemberSet.encodeVersion(nParts[0], nParts[1], nParts[2]);
        }
        else {
            nRequired = ServiceMemberSet.parseVersion(args[0]);
        }
        boolean fResult = nActual >= nRequired;

        // versions are equal
        System.out.println(" result=" + fResult);
        return fResult;
    }
}
