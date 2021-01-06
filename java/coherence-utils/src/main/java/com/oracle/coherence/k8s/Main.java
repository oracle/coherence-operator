/*
 * Copyright (c) 2019, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.lang.reflect.Method;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Coherence;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.run.xml.XmlElement;

/**
 * A main class that is used to run some initialisation code before
 * running another main class.
 */
public class Main {

    private static boolean initialised = false;

    /**
     * Private constructor for utility class.
     */
    private Main() {
    }

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     * @throws java.lang.Exception if an error occurs
     */
    public static void main(String[] args) throws Exception {
        if (args.length == 0) {
            args = new String[] {getMainClass()};
        }

        init();

        String sMainClass = args[0];
        String[] asArgsReal = new String[args.length - 1];
        System.arraycopy(args, 1, asArgsReal, 0, asArgsReal.length);

        Class<?> clsMain = Class.forName(sMainClass);
        Method method = clsMain.getMethod("main", asArgsReal.getClass());
        method.invoke(null, (Object) asArgsReal);
    }

    /**
     * Initialise the application.
     *
     * @throws IOException if there is an error starting the REST server.
     */
    public static synchronized void init() throws IOException {
        if (initialised) {
            return;
        }
        initialised = true;
        // ensure that we add the operator MBean to the management configuration
        XmlElement xml    = CacheFactory.getManagementConfig();
        XmlElement mbeans = xml.getSafeElement("mbeans");
        XmlElement mbean  = mbeans.addElement("mbean");
        mbean.addAttribute("id").setString("coherence.operator");
        mbean.addElement("mbean-class").setString(CoherenceOperator.class.getName());
        mbean.addElement("mbean-name").setString(CoherenceOperator.OBJECT_NAME);
        mbean.addElement("enabled").setBoolean(true);
        CacheFactory.setManagementConfig(xml);

        OperatorRestServer server = new OperatorRestServer();
        server.start();
    }

    private static String getMainClass() {
        try {
            return Coherence.class.getCanonicalName();
        }
        catch (Exception e) {
            return DefaultCacheServer.class.getCanonicalName();
        }
    }
}
