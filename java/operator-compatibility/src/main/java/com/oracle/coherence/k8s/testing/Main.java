/*
 * Copyright (c) 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import com.tangosol.net.DefaultCacheServer;

/**
 * Main class.
 *
 * @author Jonathan Knight
 */
public class Main {
    /**
     * Private constructor.
     */
    private Main() {
    }

    public static void main(String[] args) {
        DefaultCacheServer.main(args);
    }
}
