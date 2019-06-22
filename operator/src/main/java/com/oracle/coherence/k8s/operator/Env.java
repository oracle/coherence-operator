/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import java.util.Optional;

/**
 * A helper class for environment.
 *
 * @author as
 */
public class Env
    {
    public static String get(String name, String defaultValue)
        {
        return Optional.ofNullable(System.getenv(name)).orElse(defaultValue);
        }
    }
