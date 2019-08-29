/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.oracle.bedrock.options.LaunchLogging;
import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.ApplicationConsoleBuilder;
import com.oracle.bedrock.runtime.LocalPlatform;
import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.console.SystemApplicationConsole;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.network.AvailablePortIterator;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.runtime.options.Console;
import com.oracle.bedrock.runtime.options.DisplayName;
import com.oracle.bedrock.util.Capture;
import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;
import org.apache.commons.compress.compressors.gzip.GzipCompressorInputStream;

import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.PrintStream;
import java.io.PrintWriter;
import java.net.URL;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * @author jk
 */
public class HelmUtils
    {
    /**
     * Private constructor for utility method class.
     */
    private HelmUtils()
        {
        }

    /**
     * Extract the specified tar.gz file.
     *
     * @param fileTarget  the directory to extract to
     * @param url         the tar.gz to extract
     *
     * @throws IOException  if there is an error
     */
    public static void extractTarGZ(File fileTarget, URL url) throws IOException
        {
        fileTarget.mkdirs();

        try (InputStream in = url.openStream())
            {
            int bufferSize                   = 1024 * 1024;
            GzipCompressorInputStream gzipIn = new GzipCompressorInputStream(in);

            try (TarArchiveInputStream tarIn = new TarArchiveInputStream(gzipIn))
                {
                TarArchiveEntry entry;

                while ((entry = (TarArchiveEntry) tarIn.getNextEntry()) != null)
                    {
                    if (entry.isDirectory())
                        {
                        File    file     = new File(fileTarget, entry.getName());
                        boolean fCreated = file.mkdir();

                        if (!fCreated)
                            {
                            throw new IOException("Unable to create directory during extraction of archive contents. " +
                                                  file.getAbsolutePath());
                            }
                        }
                    else
                        {
                        byte data[] = new byte[bufferSize];
                        int  count;

                        File             file = new File(fileTarget, entry.getName());
                        FileOutputStream out  = new FileOutputStream(file, false);

                        try (BufferedOutputStream dest = new BufferedOutputStream(out, bufferSize))
                            {
                            while ((count = tarIn.read(data, 0, bufferSize)) != -1)
                                {
                                dest.write(data, 0, count);
                                }
                            }
                        }
                    }
                }
            }
        }

    /**
     * Dump the specified console output to the logs.
     *
     * @param sPrefix  the prefix to add to each line
     * @param console  the console to dump
     */
    public static void logConsoleOutput(String sPrefix, CapturingApplicationConsole console)
        {
        logConsoleOutput(sPrefix, console, System.out, System.err);
        }

    /**
     * Dump the specified console output to the logs.
     *
     * @param sPrefix  the prefix to add to each line
     * @param console  the console to dump
     * @param out      the {@link PrintStream} to write to console output to
     */
    public static void logConsoleOutput(String sPrefix, CapturingApplicationConsole console, PrintStream out)
        {
        logConsoleOutput(sPrefix, console, out, out);
        }

    /**
     * Dump the specified console output to the logs.
     *
     * @param sPrefix  the prefix to add to each line
     * @param console  the console to dump
     * @param out      the {@link PrintStream} to write to console std-out to
     * @param err      the {@link PrintStream} to write to console std-err to
     */
    public static void logConsoleOutput(String                      sPrefix,
                                        CapturingApplicationConsole console,
                                        PrintStream                 out,
                                        PrintStream                 err)
        {
        AtomicInteger cOut = new AtomicInteger();
        AtomicInteger cErr = new AtomicInteger();

        console.getCapturedOutputLines()
                .forEach(s -> out.printf("[%s:out] %4d: %s\n", sPrefix, cOut.getAndIncrement(), s));
        out.flush();
        console.getCapturedErrorLines()
                .forEach(s -> err.printf("[%s:err] %4d: %s\n", sPrefix, cErr.getAndIncrement(), s));
        err.flush();
        }


    /**
     * Remove all k8s resources with the cloudCollectionsTesting=true label.
     *
     * @param cluster    the k8s cluster
     */
    public static void cleanupTestResources(K8sCluster cluster, String sNamespace)
        {
        Arguments args = Arguments.of("delete", "all", "--selector", "cloudCollectionsTesting=true");

        if (sNamespace != null)
            {
            args = args.with("--namespace", sNamespace);
            }

        int nExitCode = cluster.kubectlAndWait(args);

        if (nExitCode != 0)
            {
            System.err.println("Clean-up: non-zero return from kubectl [" + nExitCode + "]");
            }
        }
    
    /**
     * Start a kubectl port-forward process.
     *
     * @param cluster         the k8s cluster
     * @param sPod            the pod name
     * @param sNamespace      the k8s namespace
     * @param nPort           the port to forward
     * @param consoleBuilder  the {@link ApplicationConsoleBuilder}
     *
     * @return  the kubectl {@link Application}
     */
    public static Application portForward(K8sCluster                cluster,
                                          String                    sPod,
                                          String                    sNamespace,
                                          int                       nPort,
                                          ApplicationConsoleBuilder consoleBuilder)
        {
        AvailablePortIterator ports       = LocalPlatform.get().getAvailablePorts();
        Capture<Integer>      port        = new Capture<>(ports);
        PortMapping           portMapping = new PortMapping(String.valueOf(nPort), port, nPort);
        String                sName       = "kubectl-port-forward[" + nPort + "@" + sPod + "]";
        Arguments             arguments   = Arguments.of("port-forward", sPod, portMapping);

        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        if (consoleBuilder == null)
            {
            consoleBuilder = SystemApplicationConsole.builder();
            }

        Application application = cluster.kubectl(DisplayName.of(sName), consoleBuilder, arguments);

        application.add(portMapping);

        return application;
        }


    /**
     * Obtain the list of Pod names for the given Helm release.
     *
     * @param cluster     the k8s cluster running the Pods
     * @param sNamespace  the k8s namespace to use
     * @param sSelector   the selector
     *
     * @return  the {@link List} of Pod names
     */
    public static List<String> getPods(K8sCluster cluster, String sNamespace, String sSelector)
        {
        return getK8sObject(cluster, "pods", sNamespace, sSelector);
        }

    /**
     * Obtain the list k8s object name for the given Helm release.
     *
     * @param cluster     the k8s cluster running
     * @param objectName  the k8s object name
     * @param sNamespace  the k8s namespace to use
     * @param sSelector   the selector
     *
     * @return  the {@link List} of Pod names
     */
    public static List<String> getK8sObject(K8sCluster cluster, String objectName, String sNamespace, String sSelector)
        {
        return getK8sObject(cluster, objectName, sNamespace, sSelector, "{.items[*].metadata.name}");
        }

    /**
     * Obtain the list k8s object name for the given Helm release.
     *
     * @param cluster     the k8s cluster running
     * @param objectName  the k8s object name
     * @param sNamespace  the k8s namespace to use
     * @param sSelector   the selector
     * @param sJsonPath   the json path
     *
     * @return  the {@link List} of Pod names
     */
    public static List<String> getK8sObject(K8sCluster cluster, String objectName, String sNamespace, String sSelector, String sJsonPath)
        {
        CapturingApplicationConsole console = new CapturingApplicationConsole();

        Arguments                   args    = Arguments.of("get", objectName);

        if (sNamespace != null && sNamespace.trim().length() > 0)
            {
            args = args.with("--namespace", sNamespace);
            }

        args = args.with("-o", "jsonpath=\"" + sJsonPath +"\"", "-l", sSelector);

        int nExitCode = cluster.kubectlAndWait(args, LaunchLogging.disabled(), Console.of(console));

        logConsoleOutput("get-" + objectName + "-jsonpath-" + sJsonPath, console);

        assertThat("kubectl returned non-zero exit code", nExitCode, is(0));

        String sList = console.getCapturedOutputLines().poll();

        // strip any leading quote
        if (sList.charAt(0) == '"')
            {
            sList = sList.substring(1);
            }

        // strip any trailing quote
        if (sList.endsWith("\""))
            {
            sList = sList.substring(0, sList.length() - 1);
            }

        List<String> list = new ArrayList<>((sList.trim().length() == 0) ? Collections.emptyList() : Arrays.asList(sList.split(" ")));

        Collections.sort(list);

        return list;
        }

    // ----- data members ---------------------------------------------------

    /**
     * The timeout to use for Helm commands, default 300 seconds.
     */
    public static final int HELM_TIMEOUT = Integer.parseInt(System.getProperty("helm.timeout", "300"));

    /**
     *
     */
    public static final int K8S_MIN_FORWARD_PORT = Integer.getInteger("k8s.min.forward.port", 40000);

    /**
     * An {@link ObjectMapper} to convert json to Objects and vice-versa.
     */
    public static final ObjectMapper JSON_MAPPER = new ObjectMapper();
    }
