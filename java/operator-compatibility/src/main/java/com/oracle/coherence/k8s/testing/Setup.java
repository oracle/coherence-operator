/*
 * Copyright (c) 2021, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.nio.file.attribute.PosixFilePermission;
import java.nio.file.attribute.PosixFilePermissions;
import java.util.HashSet;
import java.util.Set;

/**
 * Set-up class run by the compatibility image build.
 */
public class Setup {

    /**
     * Private constructor.
     */
    private Setup() {
    }

    /**
     * Entry point.
     *
     * @param args arguments
     *
     * @throws java.lang.Exception if the process fails to run
     */
    public static void main(String[] args) throws Exception {
        String home = System.getenv("COHERENCE_HOME");
        if (home == null || home.isEmpty()) {
            home = "/u01/coherence";
        }
        System.out.println("COHERENCE_HOME=" + home);

        copy(Paths.get(home, "lib", "coherence.jar"), Paths.get("/app/libs"));
    }

    private static void copy(Path source, Path targetDir) throws Exception {
        System.out.println("Copying " + source + " to " + targetDir);
        if (source.toFile().exists()) {
            Path target = targetDir.resolve(source.getFileName());
            Files.copy(source, target, StandardCopyOption.REPLACE_EXISTING);
            Files.setPosixFilePermissions(target, PosixFilePermissions.fromString("rw-r--r--"));
        } else {
            System.out.println("Nothing to copy, source does not exist: " + source);
        }
    }
}
