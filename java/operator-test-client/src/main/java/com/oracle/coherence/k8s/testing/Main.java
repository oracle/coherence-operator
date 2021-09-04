/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

public class Main {
    public static void main(String[] args) {
        String clientType = System.getenv("CLIENT_TYPE");

        if (clientType == null) {
            clientType = "extend";
        }

        switch (clientType.toLowerCase()) {
        case "grpc":
            GrpcClient.main(args);
        case "extend":
        default:
            ExtendClient.main(args);
        }
    }
}
