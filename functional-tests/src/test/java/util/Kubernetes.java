/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package util;

import com.oracle.bedrock.Option;
import com.oracle.bedrock.OptionsByType;

import com.oracle.bedrock.deferred.options.RetryFrequency;

import com.oracle.bedrock.options.Timeout;

import com.oracle.bedrock.runtime.ApplicationConsole;
import com.oracle.bedrock.runtime.ApplicationConsoleBuilder;

import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.console.EventsApplicationConsole;
import com.oracle.bedrock.runtime.console.SystemApplicationConsole;

import com.oracle.bedrock.runtime.k8s.K8sCluster;

import com.oracle.bedrock.runtime.options.Console;
import com.oracle.bedrock.testsupport.MavenProjectFileUtils;

import com.oracle.bedrock.util.Duration;

import org.junit.rules.TestWatcher;

import org.junit.runner.Description;

import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;

import java.nio.file.Files;
import java.nio.file.StandardOpenOption;

import java.util.Iterator;

import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeUnit;


/**
 * An extension of a Bedrock {@link K8sCluster} that allows retrying
 * of failed kubectl commands when calling the {@link #kubectlAndWait(Option...)}
 * or {@link #kubectlAndWait(Timeout, Option...)} methods.
 *
 * @author jk  2019.05.16
 */
public class Kubernetes
        extends K8sCluster<Kubernetes>
    {
    public Kubernetes()
        {
        File fileDir = MavenProjectFileUtils.ensureTestOutputBaseFolder(Kubernetes.class);
        fileDir.mkdirs();
        m_fileLog = new File(fileDir, "kubectl-retries.log");
        }

    /**
     * Set whether operations retry attempts should be logged.
     *
     * @param testName  the {@link TestName} JUnit rule to use to capture the name of the current test
     *
     * @return  this {@link Kubernetes} instance
     */
    public Kubernetes logRetries(TestName testName)
        {
        m_testName = testName;
        return this;
        }

    /**
     * Execute a kubectl command against the k8s cluster
     * and wait for the command to complete.
     * <p>
     * This implementation allows retrying of failed commands (kubectl
     * returns a non-zero exit code) by specifying a {@link MaxRetries}
     * option.
     * <p>
     * The default is to retry a maximum of five times if no {@link MaxRetries}
     * option is present.
     * <p>
     * If no {@link RetryFrequency} option is present the default is to use
     * a fibonacci back-off (see {@link RetryFrequency#fibonacci()}.
     *
     * @param timeout  the time to wait for the command to complete
     * @param options  the options to use to run the kubectl command
     *
     * @return  the exit code from the kubectl command
     */
    @Override
    public int kubectlAndWait(Timeout timeout, Option... options)
        {
        OptionsByType               opts              = OptionsByType.of(options);
        RetryFrequency              defaultFrequency  = RetryFrequency.fibonacci();
        Iterator<Duration>          retryFrequency    = opts.getOrDefault(RetryFrequency.class, defaultFrequency)
                                                            .get().iterator();
        MaxRetries                  defaultRetries    = MaxRetries.of(5);
        int                         nMaxRetries       = opts.getOrDefault(MaxRetries.class, defaultRetries)
                                                            .getMaxRetryCount();
        CapturingApplicationConsole console           = new CapturingApplicationConsole();
        ApplicationConsoleBuilder   appConsoleBuilder = opts.getOrDefault(ApplicationConsoleBuilder.class,
                                                                          SystemApplicationConsole.builder());



        opts.remove(ApplicationConsoleBuilder.class);
        opts.add(MultiplexingApplicationConsole.builder(appConsoleBuilder, Console.of(console)));

        int nExitCode = super.kubectlAndWait(timeout, opts.asArray());

        while (nExitCode != 0 && nMaxRetries > 0)
            {
            nMaxRetries--;
            long nMillis = retryFrequency.hasNext() ? retryFrequency.next().to(TimeUnit.MILLISECONDS) : 250L;

            if (m_testName != null)
                {
                try
                    {
                    StringBuilder sMessage = new StringBuilder();

                    sMessage.append("------------------------------------------------------------------------\n")
                            .append("Test: ").append(m_testName.getName()).append("\n")
                            .append("Kubectl returned a non-zero exit code (").append(nExitCode).append(")\n");

                    console.getCapturedOutputLines()
                           .forEach(sLine -> sMessage.append(sLine).append("\n"));
                    console.getCapturedErrorLines()
                           .forEach(sLine -> sMessage.append(sLine).append("\n"));

                    Files.write(m_fileLog.toPath(),
                                sMessage.toString().getBytes(),
                                StandardOpenOption.CREATE,
                                StandardOpenOption.APPEND);
                    }
                catch (IOException e)
                    {
                    System.err.println("Could not write to retry log: " + e.getMessage());
                    }
                }

            System.err.println("Kubectl returned a non-zero exit code (" + nExitCode + ") "
                               + "- retrying command in " + nMillis + " millis. "
                               + "Remaining retry attempts " + nMaxRetries);

            try
                {
                Thread.sleep(nMillis);
                }
            catch (InterruptedException e)
                {
                System.err.println("Kubectl retry thread interrupted, stopping retries");
                break;
                }

            nExitCode = super.kubectlAndWait(timeout, options);
            }

        return nExitCode;
        }


    // ----- inner class: TestName ------------------------------------------

    /**
     * A JUnit rule for capturing the description of the current test.
     */
    public static class TestName
            extends TestWatcher
        {
        /**
         * Obtain the name of the current test.
         *
         * @return   the name of the current test
         */
        public String getName()
            {
            return m_sName;
            }

        @Override
        protected void starting(Description description)
            {
            m_sName = description.getDisplayName();
            }

        /**
         *  The name of the current test
         */
        private String m_sName;
        }

    // ----- data members ---------------------------------------------------

    /**
     * The file to log retry attempts to.
     */
    private File m_fileLog;

    /**
     * The {@link TestName} JUnit rule used to capture the name of the current test.
     */
    private TestName m_testName;
    }
