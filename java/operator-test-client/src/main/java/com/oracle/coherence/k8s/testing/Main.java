/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

/**
 * A test client runner that runs specific types of client.
 */
public class Main {

    /**
     * Private constructor for utility class.
     */
    private Main() {
    }

    /**
     * Run the gRPC client.
     *
     * @param args  the program arguments.
     */
    public static void main(String[] args) {
        String clientType = System.getenv("CLIENT_TYPE");

        if (clientType == null) {
            clientType = "extend";
        }

        switch (clientType.toLowerCase()) {
        case "grpc":
            GrpcClient.main(args);
            break;
        case "extend":
        default:
            ExtendClient.main(args);
        }
    }
}
