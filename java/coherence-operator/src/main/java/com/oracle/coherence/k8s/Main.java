/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.File;
import java.io.PrintWriter;
import java.lang.reflect.Method;
import java.util.concurrent.CompletableFuture;

import com.tangosol.coherence.config.Config;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.Coherence;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.Member;

/**
 * A main class that is used to run some initialisation code before
 * running another main class.
 */
public class Main {

    private static final String DEFAULT_MAIN = "$DEFAULT$";

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
        else if (DEFAULT_MAIN.equals(args[0])) {
            args[0] = getMainClass();
        }

        init();
        CompletableFuture.runAsync(Main::initCohCtl);

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
     * @throws Exception if there is an error starting the REST server.
     */
    public static synchronized void init() throws Exception {
        if (initialised) {
            return;
        }
        initialised = true;
        CoherenceOperator.printBanner(System.out);
        OperatorRestServer server = new OperatorRestServer(System.getProperties());
        server.start();
    }

    private static String getMainClass() {
        try {
            return Coherence.class.getCanonicalName();
        }
        catch (Throwable e) {
            return DefaultCacheServer.class.getCanonicalName();
        }
    }

    private static void initCohCtl() {
        try {
            Cluster cluster = CacheFactory.getCluster();
            Member member = cluster.getLocalMember();
            String clusterName = member.getClusterName();
            String port = Config.getProperty("coherence.management.http.port", "30000");
            String provider = Config.getProperty("coherence.management.http.provider");
            String defaultProtocol = provider == null || provider.isEmpty() ? "http" : "https";
            String protocol = Config.getProperty("coherence.operator.cli.protocol", defaultProtocol);
            String home = System.getProperty("user.home");
            String connectionType = "http";

            File cohctlHome = new File(home + File.separator + ".cohctl");
            File configFile = new File(cohctlHome, "cohctl.yaml");

            if (!configFile.exists()) {
                System.out.println("CoherenceOperator: creating default cohctl config at " + configFile.getAbsolutePath());
                if (!cohctlHome.exists()) {
                    cohctlHome.mkdirs();
                }
                try (PrintWriter out = new PrintWriter(configFile)) {
                    out.println("clusters:");
                    out.println("    - name: default");
                    out.println("      discoverytype: manual");
                    out.println("      connectiontype: " + connectionType);
                    out.println("      connectionurl: " + protocol + "://127.0.0.1:" + port + "/management/coherence/cluster");
                    out.println("      nameservicediscovery: \"\"");
                    out.println("      clusterversion: \"" + CacheFactory.VERSION + "\"");
                    out.println("      clustername: \"" + clusterName + "\"");
                    out.println("      clustertype: Standalone");
                    out.println("      manuallycreated: false");
                    out.println("      baseclasspath: \"\"");
                    out.println("      additionalclasspath: \"\"");
                    out.println("      arguments: \"\"");
                    out.println("      managementport: 0");
                    out.println("      persistencemode: \"\"");
                    out.println("      loggingdestination: \"\"");
                    out.println("      managementavailable: false");
                    out.println("color: \"on\"");
                    out.println("currentcontext: default");
                    out.println("debug: false");
                    out.println("defaultbytesformat: m");
                    out.println("ignoreinvalidcerts: false");
                    out.println("requesttimeout: 30");
                }
            }
        }
        catch (Exception e) {
            System.err.println("Coherence Operator: Failed to create default cohctl config");
            e.printStackTrace();
        }
    }
}
