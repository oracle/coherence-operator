/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.DefaultCacheServer;

import java.lang.reflect.Method;

/**
 * A main class that is used to run some initialisation code before
 * running another main class.
 *
 * @author jk
 */
public class Main
    {
    // ----- Main methods ------------------------------------------------

    /**
     * Program entry point.
     *
     * @param asArgs the program command line arguments
     */
    public static void main(String[] asArgs) throws Exception
        {
        if (asArgs.length == 0)
            {
            asArgs = new String[]{DefaultCacheServer.class.getCanonicalName()};
            }

        registerHealthMBean();

        String   sMainClass = asArgs[0];
        String[] asArgsReal = new String[asArgs.length - 1];
        System.arraycopy(asArgs, 1, asArgsReal, 0, asArgsReal.length);

        Class<?> clsMain = Class.forName(sMainClass);
        Method   method  = clsMain.getMethod("main", asArgsReal.getClass());
        method.invoke(null, (Object) asArgsReal);
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Register the health MBean.
     *
     * @throws Exception if the MBean cannot be registered.
     */
    static void registerHealthMBean()
        {
        HealthMBeanFactory.registerHealthMBean();
        }
    }
