/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.LaunchLogging;
import com.oracle.bedrock.options.Timeout;

import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.ApplicationConsole;
import com.oracle.bedrock.runtime.ApplicationConsoleBuilder;

import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.console.EventsApplicationConsole;
import com.oracle.bedrock.runtime.console.FileWriterApplicationConsole;
import com.oracle.bedrock.runtime.console.SystemApplicationConsole;

import com.oracle.bedrock.runtime.k8s.K8sCluster;

import com.oracle.bedrock.runtime.k8s.helm.Helm;
import com.oracle.bedrock.runtime.k8s.helm.HelmCommand;
import com.oracle.bedrock.runtime.k8s.helm.HelmInstall;

import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.runtime.options.Console;

import com.oracle.bedrock.runtime.options.DisplayName;
import com.oracle.bedrock.testsupport.deferred.Eventually;

import com.oracle.bedrock.testsupport.junit.TestLogs;

import com.oracle.coherence.k8s.CoherenceVersion;

import com.tangosol.util.AssertionException;
import com.tangosol.util.Resources;
import com.tangosol.util.WrapperException;
import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;
import org.hamcrest.Matcher;
import org.junit.After;
import org.junit.Assume;
import org.junit.Before;
import org.junit.ClassRule;

import org.junit.Rule;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;
import util.Kubernetes;

import javax.management.MBeanServerConnection;
import javax.management.ObjectName;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;

import java.io.PrintStream;
import java.net.URL;

import java.nio.file.Files;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import java.util.Queue;
import java.util.concurrent.TimeUnit;

import java.util.function.Predicate;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;

import static helm.HelmUtils.HELM_TIMEOUT;

import static helm.HelmUtils.getK8sObject;
import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;

/**
 * A base class for executing Helm chart tests.
 *
 * @author jk
 * @author sc
 */
public abstract class BaseHelmChartTest
    {
    // ----- test lifecycle -------------------------------------------------

    @Before
    public void logTestStart()
        {
        System.err.println(">>>>>> Starting test " + getClass().getSimpleName() + " " + m_testName.getMethodName()
                                   + " >>>>>>>>>>>>>>>>");
        }

    @After
    public void logTestEnd()
        {
        System.err.println("<<<<<< Finished test " + getClass().getSimpleName() + " " + m_testName.getMethodName()
                                   + " <<<<<<<<<<<<<<<<");
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Install the Helm chart.
     *
     * @param cluster         the k8s cluster to use
     * @param sHelmChartName  the Helm chart name
     * @param sNamespace      the k8s namespace being used
     * @param sURLValue       the values file to use
     * @param asSetValues     the array of Helm set values
     *
     * @return  the name of the Helm release
     *
     * @throws Exception if there is an error
     */
    public static String installChart(K8sCluster cluster,
                                      String     sHelmChartName,
                                      URL        urlChartPackage,
                                      String     sNamespace,
                                      String     sURLValue,
                                      String...  asSetValues) throws Exception
        {
        String[] aURLValues = (sURLValue == null) ? null : new String[] { sURLValue };
        return installChart(cluster, sHelmChartName, urlChartPackage, sNamespace, aURLValues, asSetValues);
        }

    /**
     * Install the Helm chart.
     *
     * @param cluster         the k8s cluster to use
     * @param sHelmChartName  the Helm chart name
     * @param sNamespace      the k8s namespace being used
     * @param aURLValues      an array of the values file to use
     * @param asSetValues     an optional array of Helm set values
     *
     * @return  the name of the Helm release
     *
     * @throws Exception if there is an error
     */
    public static String installChart(K8sCluster cluster,
                                      String     sHelmChartName,
                                      URL        urlChartPackage,
                                      String     sNamespace,
                                      String[]   aURLValues,
                                      String...  asSetValues) throws Exception
        {
        URL[] aURL = (aURLValues == null) ? null :
                                            Arrays.stream(aURLValues)
                                                  .map(BaseHelmChartTest::toURL)
                                                  .toArray(URL[]::new);

        return installChart(cluster, sHelmChartName, urlChartPackage, sNamespace, aURL, asSetValues);
        }

    /**
     * Install the Helm chart.
     *
     * @param cluster         the k8s cluster to use
     * @param sHelmChartName  the Helm chart name
     * @param urlHelmPackage  the URL of the .tar.gz file containing the chart
     * @param sNamespace      the k8s namespace being used
     * @param aURLValues      an array of the values file to use
     * @param asSetValues     an optional array of Helm set values
     *
     * @return  the name of the Helm release
     *
     * @throws Exception if there is an error
     */
    public static String installChart(K8sCluster cluster,
                                      String     sHelmChartName,
                                      URL        urlHelmPackage,
                                      String     sNamespace,
                                      URL[]      aURLValues,
                                      String...  asSetValues) throws Exception
        {
        assertPreconditions(cluster);

        assertThat("Chart name is null", sHelmChartName, is(notNullValue()));
        assertThat("Chart tar.gz package URL is null",  urlHelmPackage, is(notNullValue()));

        // fetch and unzip helm chart
        File fileChartDir = extractChart(sHelmChartName, urlHelmPackage);

        // helm install dry run and get the release name used for the dry-run
        String sRelease = installDryRun(fileChartDir, sHelmChartName, sNamespace, aURLValues, asSetValues);

        try
            {
            // helm install real using the release name from the dry-run
            int nExitCode = install(fileChartDir, sHelmChartName, sNamespace, sRelease, aURLValues, asSetValues);

            if (nExitCode != 0)
                {
                try
                    {
                    System.err.println("Clean up Helm install '" + sRelease + "' with exit code " + nExitCode);
                    cleanupHelmReleases(sRelease);
                    cleanupPersistentVolumeClaims(cluster, sRelease, sNamespace);
                    }
                catch (Throwable t)
                    {
                    System.err.println("Error in clean up Helm release '" + sRelease + "': " + t);
                    }
                throw new Exception("Helm install '" + sRelease + "' failed with exit code: " + nExitCode);
                }
            }
        catch (Throwable t)
            {
            // cleanup helm release artifacts when an exception is thrown during helm install.
            System.err.println("Handled exception " + t.getClass().getName() + " during helm install " + sHelmChartName + " namespace=" + sNamespace + " release=" + sRelease);
            try
                {
                System.err.println("Clean up Helm install '" + sRelease + "' after handled exception: " + t);
                t.printStackTrace();
                cleanupHelmReleases(sRelease);
                cleanupPersistentVolumeClaims(cluster, sRelease, sNamespace);
                }
            catch (Throwable t1)
                {
                System.err.println("Error in clean up Helm release '" + sRelease + "': " + t1);
                }
            throw new Exception("Helm install '" + sRelease + "' failed with exception: " + t);
            }

        // Wait for the StatefulSet to be ready
        System.err.println("Installed Helm chart - release=" + sRelease);

        return sRelease;
        }

    /**
     * Extract a Helm chart .tar.gz file and verify that the chart is valid.
     *
     * @param sHelmChartName  the name of the chart
     * @param urlHelmPackage  the URL of the .tar.gz file
     *
     * @return  the directory that the .tar.gz was extracted to
     *
     * @throws IOException     if an error occurs extracting the chart
     * @throws AssertionError  if the chart is not valid
     */
    public static File extractChart(String sHelmChartName, URL urlHelmPackage) throws IOException
        {
        // fetch and unzip helm chart
        File fileTempDir = s_temp.newFolder();

        assertThat(fileTempDir.exists(), is(true));

        System.err.println("Extracting Helm chart " + urlHelmPackage + " into " + fileTempDir);

        HelmUtils.extractTarGZ(fileTempDir, urlHelmPackage);

        // Run Helm lint to verify that the chart is valid.
        assertHelmLint(fileTempDir, sHelmChartName);

        return fileTempDir;
        }

    /**
     * Run Helm lint on the chart.
     *
     * @param fileDir  the directory containing the chart
     * @param sChart   the name of the chart
     */
    protected static void assertHelmLint(File fileDir, String sChart)
        {
        File fileChartDir = new File(fileDir, sChart);

        assertThat(fileChartDir.exists(), is(true));
        assertThat(fileChartDir.isDirectory(), is(true));

        // Run Helm lint to verify that the chart is valid.
        int nExitCode = s_helm.lint(fileDir, sChart)
                          .executeAndWait();

        assertThat("Helm lint failed", nExitCode, is(0));
        }

    /**
     * Perform a dry-run install.
     *
     * @param fileDir         the directory containing the chart
     * @param sHelmChartName  the Helm chart name
     * @param sNamespace      the k8s namespace being used
     * @param aURLValues      the values files to use
     * @param asSetValues     an optional array of Helm set values
     *
     * @return  the name of the release
     */
    protected static String installDryRun(File fileDir, String sHelmChartName, String sNamespace, URL[] aURLValues, String... asSetValues) throws Exception
        {
        CapturingApplicationConsole consoleInstall = new CapturingApplicationConsole();

        HelmInstall install = s_helm.install(fileDir, sHelmChartName)
                                    .dryRun()
                                    .debug()
                                    .timeout(HELM_TIMEOUT)
                                    .set("coherenceK8sOperatorTesting=true")
                                    .set(asSetValues);

        int nExitCode = install(install, sNamespace, Console.of(consoleInstall), aURLValues);

        if (nExitCode != 0)
            {
            System.err.println("Helm dry-run install failed with non-zero exit code. Helm dry-run install will be retried.");

            for (int i = 0 ; i < Integer.parseInt(HELM_INSTALL_MAX_RETRY); i++)
                {
                nExitCode = install(install, sNamespace, Console.of(consoleInstall), aURLValues);

                if (nExitCode == 0)
                    {
                    System.err.println(String.format("Helm dry-run install successful with retry attempt: %d", i + 1));
                    break;
                    }
                }
            }

        HelmUtils.logConsoleOutput("helm-install", consoleInstall);
        String reason = "Install dry-run failed for helm chart " + sHelmChartName + " namespace " + sNamespace;
        assertThat(reason, nExitCode, is(0));

        List<String> listLines = new ArrayList<>(consoleInstall.getCapturedOutputLines());

        String sNameLine = listLines.stream()
                                    .filter(s -> s.startsWith("NAME:"))
                                    .findFirst()
                                    .orElse(null);

        assertThat("installDryRun: failed to get a release name", sNameLine, is(notNullValue()));

        String sRelease = sNameLine.substring(5).trim();

        assertThat("installDryRun: failed to get a release name", sRelease.isEmpty(), is(false));

        File   fileRoot = s_testLogs.getOutputFolder();
        String sName    = "helm-install-" + sRelease;
        File   fileLog  = new File(fileRoot, sName);

        try (PrintStream out = new PrintStream(fileLog))
            {
            HelmUtils.logConsoleOutput(sName, consoleInstall, out);
            out.flush();
            }

        return sRelease;
        }

    /**
     * Perform a Helm install.
     *
     * @param fileDir         the directory containing the chart to install
     * @param sHelmChartName  the Helm chart name
     * @param sNamespace      the k8s namespace being used
     * @param sRelease        the name to use for the release
     * @param aURLValues      an array of the values files to use
     * @param asSetValues     an optional array of Helm set values
     *
     * @return the exit code of Helm install
     *
     * @throws Exception if an error occurs during install
     */
    protected static int install(File fileDir, String sHelmChartName, String sNamespace, String sRelease, URL[] aURLValues, String... asSetValues) throws Exception
        {
        HelmInstall install = s_helm.install(fileDir, sHelmChartName)
                                    .set("coherenceK8sOperatorTesting=true")
                                    .set(asSetValues)
                                    .timeout(HELM_TIMEOUT)
                                    .name(sRelease);

        int nExitCode = install(install, sNamespace, SystemApplicationConsole.builder(), aURLValues);

        if (nExitCode != 0)
            {
            System.err.println(String.format("Helm install of \"%s\" failed with non-zero exit code."
                               + "Helm install of \"%s\" will be retried.", sRelease, sRelease));

            for (int i = 0 ; i < Integer.parseInt(HELM_INSTALL_MAX_RETRY); i++)
                {
                try
                    {
                    System.err.println("Clean up resources before retrying install of release: '" + sRelease);
                    cleanupHelmReleases(sRelease);
                    cleanupPersistentVolumeClaims(getDefaultCluster(), sRelease, sNamespace);
                    System.err.println(String.format("Finish cleaning existing release \"%s\" before retry attempt: %d",
                            sRelease, i + 1));
                    }
                catch (Throwable t)
                    {
                    System.err.println("Error in clean up Helm release '" + sRelease + "': " + t);
                    }

                nExitCode = install(install, sNamespace, SystemApplicationConsole.builder(), aURLValues);

                if (nExitCode == 0)
                    {
                    System.err.println(String.format("Helm install of \"%s\" successful with retry attempt: %d",
                            sRelease, i + 1));
                    break;
                    }
                }
            }

        return nExitCode;
        }

    /**
     * Execute a {@link HelmInstall} command.
     *
     * @param install     the {@link HelmInstall} command to execute
     * @param sNamespace  the k8s namespace being used
     * @param console     the {@link ApplicationConsoleBuilder} to use
     * @param aURLValues  an optional set of values files to use
     *
     * @throws Exception if an error occurs during install
     */
    private static int install(HelmInstall install, String sNamespace, ApplicationConsoleBuilder console, URL... aURLValues) throws Exception
        {
        if (aURLValues != null && aURLValues.length > 0)
            {
            File file;

            for (URL url : aURLValues)
                {
                if ("file".equalsIgnoreCase(url.getProtocol()))
                    {
                    file = new File(url.getFile());
                    }
                else
                    {
                    file = s_temp.newFile("values.yaml");
                    Files.copy(url.openStream(), file.toPath());
                    }

                // dump the values to the log
                System.err.println("Running Helm install using values from " + url);
                System.err.println("---------------------------------------");
                Files.lines(file.toPath())
                        .forEach(System.err::println);
                System.err.println("---------------------------------------");
                System.err.flush();

                install = install.values(file);
                }
            }

        if (sNamespace != null && sNamespace.trim().length() > 0)
            {
            install = install.namespace(sNamespace.trim());
            }

        if (K8S_IMAGE_PULL_SECRET != null && K8S_IMAGE_PULL_SECRET.trim().length() > 0)
            {
            install = install.set("imagePullSecrets={" + K8S_IMAGE_PULL_SECRET + "}");
            }

        return install.executeAndWait(console);
        }

    public static void captureInstalledPodLogs(K8sCluster cluster, String sNamespace, String sRelease)
        {
        try
            {
            CapturingApplicationConsole consolePods       = new CapturingApplicationConsole();
            CapturingApplicationConsole consoleContainers = new CapturingApplicationConsole();
            int                         nExitCode;

            nExitCode = cluster.kubectlAndWait(Arguments.of("get", "pods", "--all-namespaces=true", "--selector", "release=" + sRelease,  "-o", "name"),
                                               Console.of(consolePods),
                                               LaunchLogging.disabled());

            if (nExitCode != 0)
                {
                System.err.println("Cannot obtain Pod names so cannot capture logs");
                return;
                }

            Queue<String> queuePods = consolePods.getCapturedOutputLines();

            for (String sPod : queuePods)
                {
                if (sPod.equals("(terminated)"))
                    {
                    continue;
                    }

                nExitCode = cluster.kubectlAndWait(Arguments.of("get", sPod, "--all-namespaces=true", "-o", "jsonpath={.spec.containers[*].name}"),
                                                   Console.of(consoleContainers),
                                                   LaunchLogging.disabled());

                String[] asContainers;

                if (nExitCode == 0)
                    {
                    asContainers = consoleContainers.getCapturedOutputLines().poll().split("\\s");
                    }
                else
                    {
                    System.err.println("Cannot obtain Pod container names so cannot capture logs");
                    asContainers = new String[]{null};
                    }

                for (String sContainer : asContainers)
                    {
                    Arguments arguments = Arguments.of("logs", sPod, "-f");

                    if (sNamespace != null)
                        {
                        arguments = arguments.with("--namespace", sNamespace);
                        }

                    if (sContainer != null)
                        {
                        arguments = arguments.with("-c", sContainer);
                        }

                    String sName = sContainer == null ? sPod.substring(4) : sPod.substring(4) + "-" + sContainer;

                    cluster.kubectl(arguments,
                                    s_testLogs.builder(),
                                    DisplayName.of(sName),
                                    LaunchLogging.disabled());
                    }
                }
            }
        catch (Throwable thrown)
            {
            System.err.println("Cannot capure POD logs: " + thrown.getMessage());
            }
        }


    /**
     * Obtain all of the k8s resources created for the release.
     *
     * @param cluster   the k8s cluster
     * @param sRelease  the name of the release
     *
     * @return a {@link Map} of K8s resources
     */
    @SuppressWarnings("unchecked")
    protected Map<String, ?> getK8sResources(K8sCluster cluster, String sRelease) throws Exception
        {
        CapturingApplicationConsole consoleK8sGetAll = new CapturingApplicationConsole();

        cluster.kubectlAndWait(Arguments.of("get", "all", "--selector", "release=" + sRelease,  "-o", "json"),
                               Console.of(consoleK8sGetAll));

        HelmUtils.logConsoleOutput("kubectl-get", consoleK8sGetAll);

        String sJson  = consoleK8sGetAll.getCapturedOutputLines().stream().collect(Collectors.joining());

        return HelmUtils.JSON_MAPPER.readValue(sJson, Map.class);
        }

    /**
     * Assert the test pre-conditions.
     *
     * @param cluster  the k8s cluster
     */
    protected static void assertPreconditions(K8sCluster cluster)
        {
        Assume.assumeThat("The k8s cluster to use has not been set", cluster, is(notNullValue()));
        Assume.assumeThat("Helm executable " + s_helm.getHelmLocation()
                                  + " does not exists or is not executable", s_helm.helmExists(), is(true));

        }

    /**
     * Assert whether the Helm deployment is ready or not (i.e all Pods are ready).
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the k8s namespace being used
     * @param sSelector   the selector
     * @param fReady      the readiness of the deployment
     *
     * @exception AssertionError  if the deployment is not ready
     */
    protected static void assertDeploymentReady(K8sCluster cluster, String sNamespace, String sSelector, boolean fReady)
        {
        Eventually.assertThat("assertDeploymentReady namespace=" + sNamespace + " selector=" + sSelector,
            invoking(STUB).isDeploymentReady(cluster, sNamespace, sSelector), is(fReady),
                Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS));
        }

    /**
     * Assert that the Coherence stateful set's service is running.
     *
     * @param k8sCluster  the k8s cluster to use
     * @param sNamespace  the k8s namespace the service is in
     * @param sRelease    the name of the Helm release
     */
    protected void assertCoherenceService(K8sCluster k8sCluster, String sNamespace, String sRelease)
        {
        String       sCoherenceServiceSelector = getCoherenceServiceSelector(sRelease);
        List<String> listSvcs                  = getK8sObject(k8sCluster,
                                                              "svc",
                                                              sNamespace,
                                                              sCoherenceServiceSelector);

        assertThat("assertCoherenceService assertThat Coherence stateful set service is running", listSvcs.size(), is(1));

        String sReleaseName = sRelease + "-" + COHERENCE_HELM_CHART_NAME;
        assertThat("assertCoherenceService assert release name", listSvcs.get(0), is(sReleaseName));
        }

    /**
     * Return a Coherence Pod selector.
     *
     * @param sRelease  release name
     *
     * @return  a Coherence Pod selector
     */
    protected static String getCoherenceSelector(String sRelease)
        {
        return getSelector("coherence", sRelease);
        }

    /**
     * Return a Coherence Pod selector.
     *
     * @param sRelease  release name
     *
     * @return  a Coherence Pod selector
     */
    protected static String getCoherencePodSelector(String sRelease)
        {
        return getSelector("coherencePod", sRelease);
        }

    /**
     * Return a Coherence JMX Pod selector.
     *
     * @param sRelease  release name
     *
     * @return  a Coherence JMX Pod selector
     */
    protected static String getCoherenceJmxPodSelector(String sRelease)
        {
        return getSelector("coherenceJMXPod", sRelease);
        }

    /**
     * Return a Coherence JMX Deployment selector.
     *
     * @param sRelease  release name
     *
     * @return  a Coherence JMX Deployment selector
     */
    protected static String getCoherenceJmxDeploymentSelector(String sRelease)
        {
        return getSelector("coherence-jmx", sRelease);
        }

    /**
     * Return a Coherence Service selector.
     *
     * @param sRelease  release name
     *
     * @return  a Coherence Service selector
     */
    protected static String getCoherenceServiceSelector(String sRelease)
        {
        return getSelector("coherence-service", sRelease);
        }

    /**
     * Return a Coherence Operator selector.
     *
     * @param sRelease  release name
     *
     * @return  a Coherence Operator selector
     */
    protected static String getCoherenceOperatorSelector(String sRelease)
        {
        return getSelector("coherence-operator", sRelease);
        }

    /**
     * Return a Kibana selector.
     *
     * @param sRelease  release name
     *
     * @return  a Kibana selector
     */
    protected static String getKibanaSelector(String sRelease)
        {
        return getSelector("kibana", sRelease);
        }


    /**
     * Return a ElasticSearch selector.
     *
     * @param sRelease  release name
     *
     * @return  a ElasticSearch selector
     */
    protected static String getElasticSearchSelector(String sRelease)
        {
        return getSelector("elasticsearch", sRelease);
        }

    /**
     * Return a selector with the given app name and release name.
     *
     * @param sComp     component name
     * @param sRelease  release name
     * @return  the selector with the given app name and release name
     */
    protected static String getSelector(String sComp, String sRelease)
        {
        return "component=" + sComp + ",release=" + sRelease;
        }

    /**
     * Determine whether the Helm deployment is ready (i.e all Pods are ready).
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the k8s namespace being used
     * @param sSelector   the selector
     *
     * @return  {@code true} if the deployment is ready
     */
    @SuppressWarnings("unchecked")
    // must be public - used in Eventually.assertThat call.
    public boolean isDeploymentReady(K8sCluster cluster, String sNamespace, String sSelector)
        {
        Map<String, ?> map = getJson(cluster, sNamespace, sSelector, "deployment", false);

        if (map == null)
            {
            return false;
            }

        Map<String, ?> mapSpec    = (Map<String, ?>) map.get("spec");
        Number         cReplicas  = mapSpec   == null ? null : (Number) mapSpec.get("replicas");
        Map<String, ?> mapStatus  = (Map<String, ?>) map.get("status");
        Number         cAvailable = mapStatus == null ? null : (Number) mapStatus.get("availableReplicas");

        return cAvailable != null && cReplicas != null && cAvailable.intValue() == cReplicas.intValue();
        }

    /**
     * Assert whether the Helm stateful set is ready or not (i.e all Pods are ready).
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the k8s namespace being used
     * @param sSelector   the selector
     * @param fReady      readiness of the statefulset
     *
     * @throws AssertionError  if the stateful set is not ready
     */
    static void assertStatefulSetReady(K8sCluster cluster, String sNamespace, String sSelector, boolean fReady)
        {
        Eventually.assertThat(invoking(STUB).isStatefulSetReady(cluster, sNamespace, sSelector), is(fReady),
                Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS));
        }

    /**
     * Determine whether the Helm stateful set is ready (i.e all Pods are ready).
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the k8s namespace being used
     * @param sSelector   the selector
     *
     * @return  {@code true} if the deployment is ready
     */
    @SuppressWarnings("unchecked")
    // must be public - used in Eventually.assertThat call.
    public boolean isStatefulSetReady(K8sCluster cluster, String sNamespace, String sSelector)
        {
        Map<String, ?> map = getJson(cluster, sNamespace, sSelector, "statefulset", false);

        if (map == null)
            {
            return false;
            }

        Map<String, ?> mapSpec    = (Map<String, ?>) map.get("spec");
        Number         cReplicas  = mapSpec   == null ? null : (Number) mapSpec.get("replicas");
        Map<String, ?> mapStatus  = (Map<String, ?>) map.get("status");
        Number         cReady     = mapStatus == null ? null : (Number) mapStatus.get("readyReplicas");

        return cReady != null && cReplicas != null && cReady.intValue() == cReplicas.intValue();
        }

    /**
     * Determine whether the coherence start up log is ready.
     *
     * param cluster      the k8s cluster
     * @param sNamespace  the namespace of coherence
     * @param sPod        the pod name of coherence
     *
     * @return  {@code true} if the Logstash is ready
     */
    // must be public - used in Eventually.assertThat call.
    public boolean hasDefaultCacheServerStarted(K8sCluster cluster, String sNamespace, String sPod)
        {
        try
            {
            Queue<String> sLogs = getPodLog(cluster, sNamespace, sPod);
            return sLogs.stream().anyMatch(l -> l.contains("Started DefaultCacheServer"));
            }
        catch (Exception ex)
            {
            return false;
            }
        }

    /**
     * Determine whether the Helm deployment is ready (i.e all Pods are ready).
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the k8s namespace being used
     * @param sSelector   the selector
     * @param fLogException  only log exception if true. can fail if deployment or statefulset is not deployed yet.
     *
     * @return  {@code true} if the deployment is ready
     */
    @SuppressWarnings("unchecked")
    protected Map<String, ?> getJson(K8sCluster cluster, String sNamespace, String sSelector, String sType, boolean fLogException)
        {
        CapturingApplicationConsole console   = new CapturingApplicationConsole();
        Arguments                   arguments = Arguments.of("get", sType, "--selector", sSelector, "-o", "json");

        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        int nExitCode = cluster.kubectlAndWait(arguments, LaunchLogging.disabled(), Console.of(console));

        if (nExitCode == 0)
            {
            try
                {
                String               sJson     = console.getCapturedOutputLines().stream().collect(Collectors.joining());
                Map<String, ?>       map       = HelmUtils.JSON_MAPPER.readValue(sJson, Map.class);
                List<Map<String, ?>> listItems = (List<Map<String, ?>>) map.get("items");

                if (listItems != null && listItems.size() > 0)
                    {
                    return listItems.get(0);
                    }
                }
            catch (IOException e)
                {
                if (fLogException)
                    {
                    System.err.println(getClass().getSimpleName() + "is" + sType + "Ready() failed: " + e.getMessage());
                    }
                }
            }

        return null;
        }

    /**
     * Obtain the yaml to use to install a client.
     *
     * @param name          the client name
     * @param sRelease      the Helm release name
     * @param sClusterName  the cluster name
     *
     * @return the yaml to use
     *
     * @throws IOException if there is an error creating the template
     */
    protected String getClientYaml(String name, String sRelease, String sClusterName) throws IOException
        {
        final Map<String, String> templateParams = new HashMap<>();

        templateParams.put("%%TEST_REGISTRY_PREFIX%%", System.getProperty("test.image.prefix"));
        templateParams.put("%%NAME%%", name);
        templateParams.put("%%NAMESPACE%%", getK8sNamespace());
        templateParams.put("%%WKA%%", sRelease + "-coherence-headless");
        templateParams.put("%%LISTEN_PORT%%", "20000");
        templateParams.put("%%CLUSTER%%", sClusterName);
        templateParams.put("%%IMAGE_PULL_POLICY%%", (OP_IMAGE_PULL_POLICY == null) ? "IfNotPresent" : OP_IMAGE_PULL_POLICY);
        templateParams.put("%%IMAGE_PULL_SECRETS%%", K8S_IMAGE_PULL_SECRET);

        String clientTemplateFile = Resources.findFileOrResource("coh-client-template.yaml", null).getPath();
        String clientYamlContents = IOUtils.toString(new FileInputStream(clientTemplateFile), "UTF-8");

        for (Map.Entry<String, String> entry : templateParams.entrySet())
            {
            String sValue = entry.getValue();

            if (sValue == null)
                {
                sValue = "";
                }

            clientYamlContents = clientYamlContents.replaceAll(entry.getKey(), sValue);
            }

        File clientYaml = new File(clientTemplateFile.substring(0, clientTemplateFile.lastIndexOf("/")), name);

        FileUtils.writeStringToFile(clientYaml, clientYamlContents, "UTF-8");

        return clientYaml.getPath();
        }

    /**
     * Remove all k8s resources for the Helm releases.
     *
     * @param asRelease  the Helm releases to delete
     */
    public static void cleanupHelmReleases(String... asRelease)
        {
        int nAttempt = 1;
        int nMax     = 10;

        if (asRelease != null)
            {
            for (String sRelease : asRelease)
                {
                int nExitCode = s_helm.delete(sRelease)
                                      .purge()
                                      .executeAndWait(Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS));

                if (nExitCode == 0)
                    {
                    break;
                    }

                if (nAttempt >= nMax)
                    {
                    System.err.println("Maximum number of attempts reached trying to run helm delete for " + sRelease);
                    }

                System.err.println("Clean-up: (attempt " + nAttempt + ") non-zero return from Helm delete for release "
                                           + sRelease + " [" + nExitCode + "]");
                nAttempt++;
                }
            }
        }

    /**
     * Delete prometheus-operator CRDs.
     *
     * @param cluster  the k8s cluster
     */
    public static void cleanupPrometheusOperatorCRDs(K8sCluster cluster)
        {
        // Prometheus operator is a subchart of coherence-operator.
        /*
         * See Uninstalling the Chart on link: https://github.com/helm/charts/tree/master/stable/prometheus-operator
         *
         * CRDs created by this chart are not removed by default and should be manually cleaned up:
         *
         * kubectl delete crd prometheus.monitoring.coreos.com
         * kubectl delete crd prometheusrules.monitoring.coreos.com
         * kubectl delete crd servicemonitors.monitoring.coreos.com
         * kubectl delete crd alertmanagers.monitoring.coreos.com
         */
        for (String sCRD : PROMETHEUS_OPERATOR_CRDS)
            {
            deleteCRD(cluster, sCRD);
            }
        }

    /**
     * Delete CRD.
     *
     * @param cluster  the k8s cluster
     * @param sName    CRD name
     */
    public static void deleteCRD(K8sCluster cluster, String sName)
        {
        int nExitCode = cluster.kubectlAndWait(Arguments.of("delete", "crd", "--ignore-not-found=true", sName));

        if (nExitCode != 0)
            {
            // not an error but print out so can catch an unneeded deleteCRD in functional-test.
            System.err.println("kubectl returned non-zero exit deleting crd named " + sName);
            }
        }

    public static void cleanupPrometheusOperator(K8sCluster cluster)
        {
        cleanupPrometheusOperatorCRDs(cluster);
        }

    /**
     * Capture the Pod logs for the specified release in the current namespace.
     *
     * @param clsTest     the test class
     * @param cluster     the K8s cluster
     * @param sSelector   the selector
     */
    static void capturePodLogs(Class<?> clsTest, K8sCluster cluster, String sSelector)
        {
        capturePodLogs(clsTest, cluster, getK8sNamespace(), sSelector, new String[0]);
        }

    /**
     * Capture the Pod logs for the specified release in the current namespace.
     *
     * @param clsTest     the test class
     * @param cluster     the K8s cluster
     * @param sSelector   the selector
     * @param sContainer  the container
     */
    static void capturePodLogs(Class<?> clsTest, K8sCluster cluster, String sSelector, String sContainer)
        {
        capturePodLogs(clsTest, cluster, getK8sNamespace(), sSelector, sContainer);
        }

    /**
     * Capture the Pod logs for the specified release.
     *
     * @param clsTest     the test class
     * @param cluster     the K8s cluster
     * @param sNamespace  the k8s namespace
     * @param sSelector   the selector
     * @param asContainer the container to capture logs for
     */
    static void capturePodLogs(Class<?> clsTest, K8sCluster cluster, String sNamespace, String sSelector, String... asContainer)
        {
        try
            {
            List<String> listPods = getPods(cluster, sNamespace, sSelector);

            for (String sPod : listPods)
                {
                capturePodLog(clsTest, cluster, sNamespace, sPod, asContainer);
                }
            }
        catch (Throwable t)
            {
            System.err.println("Error capturing Pod logs");
            t.printStackTrace();
            }
        }

    /**
     * Capture the Pod logs for the specified release and Pod.
     *
     * @param clsTest     the test class
     * @param cluster     the K8s cluster
     * @param sNamespace  the k8s namespace
     * @param sPod        the pod
     * @param asContainer the container to capture logs for
     */
    static void capturePodLog(Class<?> clsTest, K8sCluster cluster, String sNamespace, String sPod, String... asContainer)
        {
        if (asContainer == null || asContainer.length == 0)
            {
            asContainer = new String[]{null};
            }

        for (String sContainer : asContainer)
            {
            try
                {
                File   fileOut   = s_testLogs.getOutputFolder();
                File   fileRoot  = new File(new File(fileOut, "k8slogs"), clsTest.getSimpleName());
                fileRoot.mkdirs();

                String                       sLogName  = sPod + (sContainer == null ? "" : "-" + sContainer);
                File                         fileLog   = new File(fileRoot, sLogName + ".log");
                FileWriterApplicationConsole console   = new FileWriterApplicationConsole(fileLog.getCanonicalPath());

                Arguments arguments = Arguments.empty();

                if (sNamespace != null)
                    {
                    arguments = arguments.with("--namespace", sNamespace);
                    }

                arguments = arguments.with("logs", sPod);

                if (sContainer != null)
                    {
                    arguments = arguments.with(sContainer);
                    }

                int nExitCode = cluster.kubectlAndWait(arguments, Console.of(console));

                if (nExitCode != 0)
                    {
                    System.err.println("kubectl returned non-zero exit code capturing logs for Pod " + sPod);
                    }
                }
            catch (Exception e)
                {
                System.err.println("Error capturing Pod logs for Pod " + sPod + " container " + sContainer);
                e.printStackTrace();
                }
            }
        }

    protected static Queue<String> getPodLog(K8sCluster cluster, String sNamespace, String sPod)
        {
        CapturingApplicationConsole console   = new CapturingApplicationConsole();

        Arguments arguments = Arguments.empty();

        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        arguments = arguments.with("logs", sPod);

        int nExitCode = cluster.kubectlAndWait(arguments, Console.of(console), LaunchLogging.disabled());

        if (nExitCode != 0)
            {
                throw new IllegalStateException("kubectl returned non-zero exit code capturing logs for Pod " + sPod);
            }

        return console.getCapturedOutputLines();
        }

    protected static Queue<String> getPodLog(K8sCluster cluster, String sNamespace, String sPod, String sContainer)
        {
        CapturingApplicationConsole console = new CapturingApplicationConsole();

        getPodLog(cluster, sNamespace, sPod, sContainer, console);

        return console.getCapturedOutputLines();
        }

    private static void getPodLog(K8sCluster         cluster,
                                  String             sNamespace,
                                  String             sPod,
                                  String             sContainer,
                                  ApplicationConsole console)
        {
        Arguments arguments = Arguments.empty();

        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        arguments = arguments.with("logs", sPod);

        if (sContainer != null)
            {
            arguments = arguments.with("-c", sContainer);
            }

        cluster.kubectlAndWait(arguments, Console.of(console), LaunchLogging.disabled());
        }

    protected static void dumpPodLog(K8sCluster cluster, String sNamespace, String sPod)
        {
        dumpPodLog(cluster, sNamespace, sPod, null);
        }

    protected static void dumpPodLog(K8sCluster cluster, String sNamespace, String sPod, String sContainer)
        {
        Arguments arguments = Arguments.empty();

        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        arguments = arguments.with("logs", sPod);

        if (sContainer != null)
            {
            arguments = arguments.with("-c", sContainer);
            }

        System.err.println("-----------------------------------------------------------");
        System.err.println("Logs for Pod " + sPod + " container " + sContainer);
        System.err.println("-----------------------------------------------------------");
        cluster.kubectlAndWait(arguments, SystemApplicationConsole.builder(), LaunchLogging.disabled());
        System.err.println("-----------------------------------------------------------");
        }

    /**
     * Obtain the {@link URL} of the specified resource.
     *
     * @param sResource  the name of the resource to locate
     *
     * @return  the {@link URL} of the resource
     */
    static URL toURL(String sResource)
        {
        return Resources.findFileOrResource(sResource, null);
        }

    /**
     * Obtain the default {@link K8sCluster} based
     * on the System properties passed to the test.
     *
     * @return  the default {@link K8sCluster}.
     */
    protected static K8sCluster getDefaultCluster()
        {
        // try the property first to get the configuration file
        String sConfig     = getK8sConfig();
        File   fileConfig  = sConfig == null ? null : new File(sConfig);
        File   fileKubectl = KUBECTL == null ? null : new File(KUBECTL);

        return new Kubernetes()
                .logRetries(s_k8sTestName)
                .withKubectlAt(fileKubectl)
                .withKubectlConfig(fileConfig)
                .withKubectlContext(getPropertyOrNull("k8s.context"))
                .withKubectlInsecure(true);
        }

    /**
     * Ensure that the namespace exists, creating it if the {@link #CREATE_NAMESPACE}
     * flag is {@code true}.
     *
     * @param cluster  the k8s cluster
     *
     * @throws  AssertionError if the namespace does not exist and {@link #CREATE_NAMESPACE} is false
     */
    protected static void ensureNamespace(K8sCluster cluster)
        {
        String sNamespace = getK8sNamespace();

        if (sNamespace != null && sNamespace.length() > 0)
            {
            ensureNamespace(cluster, sNamespace);
            }
        }

    /**
     * Ensure that the namespace exists, creating it if the {@link #CREATE_NAMESPACE}
     * flag is {@code true}.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the name of the namespace
     *
     * @throws  AssertionError if the namespace does not exist and {@link #CREATE_NAMESPACE} is false
     */
    static void ensureNamespace(K8sCluster cluster, String sNamespace)
        {
        System.err.println("Ensuring existence of namespace " + sNamespace + " - create flag is " + CREATE_NAMESPACE);

        boolean fExists = STUB.hasNamespace(cluster, sNamespace);

        if (!fExists)
            {
            if (CREATE_NAMESPACE)
                {
                // namespace does not exist, try creating it
                System.err.println("Creating namespace " + sNamespace);

                int nExitCode = cluster.kubectlAndWait(Arguments.of("create", "namespace", sNamespace));

                assertThat("Could not create namespace " + sNamespace, nExitCode, is(0));
                }
            else
                {
                throw new AssertionError("Namespace " + sNamespace
                                         + " does not exist and create namespace flag is false");
                }
            }
        else
            {
            System.err.println("Namespace " + sNamespace + " already exists.");
            }
        }

    /**
     * If the k8s namespace was created as part of the test then delete it.
     *
     * @param cluster the k8s cluster
     */
    protected static void cleanupNamespace(K8sCluster cluster)
        {
        if (CREATE_NAMESPACE)
            {
            String sNamespace = getK8sNamespace();

            if (sNamespace != null && sNamespace.length() > 0)
                {
                cleanupNamespace(cluster, sNamespace);
                }
            }
        }

    /**
     * If the k8s namespace was created as part of the test then delete it.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the name of the namespace
     */
    protected static void cleanupNamespace(K8sCluster cluster, String sNamespace)
        {
        if (CREATE_NAMESPACE)
            {
            Arguments arguments = Arguments.of("delete", "namespace", sNamespace);

            cluster.kubectlAndWait(arguments);

            Timeout timeout = Eventually.within(5, TimeUnit.MINUTES);
            Eventually.assertThat(invoking(STUB).hasNamespace(cluster, sNamespace), is(false), timeout, RetryFrequency.every(12, TimeUnit.SECONDS));
            }
        }

    /**
     * Determine whether the specified namespace exists.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the name of the namespace
     *
     * @return  {@code true} if the namespace exists
     */
    public boolean hasNamespace(K8sCluster cluster, String sNamespace)
        {
        CapturingApplicationConsole console   = new CapturingApplicationConsole();
        Arguments                   arguments = Arguments.of("get", "namespace");
        int                         nExitCode = cluster.kubectlAndWait(arguments,
                                                                       Console.of(console),
                                                                       LaunchLogging.disabled());

        if (nExitCode == 0)
            {
            return console.getCapturedOutputLines().stream()
                          .anyMatch(line -> line.startsWith(sNamespace + " "));
            }

        return false;
        }

    /**
     * Return true iff all of the {@link #PROMETHEUS_OPERATOR_CRDS Prometheus Operator Custom Resource Definition} (CRD) exist.
     * If partial number of Prometheus Operator CRD installed, delete all of the remaining ones and return false.
     *
     * @param cluster  the k8s cluster
     *
     * @return  {@code true} iff all of the Prometheus Operator CRDs are installed
     */
    public boolean hasPrometheusOperatorCRD(K8sCluster cluster)
        {
        CapturingApplicationConsole console   = new CapturingApplicationConsole();
        Arguments                   arguments = Arguments.of("get", "crd");
        int                         nExitCode = cluster.kubectlAndWait(arguments,
                                                                       Console.of(console),
                                                                       LaunchLogging.enabled());
        if (nExitCode == 0)
            {
            int nCRDMatches = 0;
            for (String CRD : PROMETHEUS_OPERATOR_CRDS)
                {
                for (String line: console.getCapturedOutputLines())
                    {
                    if (line.contains(CRD))
                        {
                        nCRDMatches++;
                        break;
                        }
                    }
                }

            if (nCRDMatches == PROMETHEUS_OPERATOR_CRDS.length)
                {
                // all PROMETHEUS_OPERATOR CRDs are installed.
                return true;
                }
            else if (nCRDMatches == 0)
                {
                return false;
                }
            else
                {
                // some CRDs installed, some missing.
                // delete remaining ones and return false.
                cleanupPrometheusOperatorCRDs(cluster);
                return false;
                }
            }

        return false;
        }

    static void cleanupPersistentVolumeClaims(K8sCluster cluster, String sRelease, String sNamespace)
        {
        CapturingApplicationConsole console   = new CapturingApplicationConsole();
        String                      sCluster  = sRelease + "-coherence";
        Arguments                   arguments = Arguments.of("get", "pvc", "-o", "name", "--selector",
                                                             "coherenceDeployment=" + sCluster, "--namespace", sNamespace);
        int                         nExitCode = cluster.kubectlAndWait(arguments,
                                                                       Console.of(console),
                                                                       LaunchLogging.enabled());

        if (nExitCode == 0)
            {
            console.getCapturedOutputLines().stream()
                    .filter(s -> !"(terminated)".equals(s))
                    .forEach(s -> deleteResource(cluster, s, sNamespace));
            }
        }

    /**
     * Delete the specified resource.
     *
     * @param cluster    the k8s cluster
     * @param sName      the name of the resource in the format resource-type/name
     * @param sNamespace the namespace
     */
    protected static void deleteResource(K8sCluster cluster, String sName, String sNamespace)
        {
        cluster.kubectlAndWait(Arguments.of("delete", sName, "--namespace", sNamespace));
        }

    /**
     * Obtain the name of the K8s configuration file to use.
     *
     * @return  the name of the K8s configuration to use
     */
    static String getK8sConfig()
        {
        // try the property first to get the configuration file
        String sConfig = getPropertyOrNull("k8s.config");

        if (sConfig == null || sConfig.trim().length() == 0)
            {
            // no property so try the env variable to get the configuration file
            sConfig = System.getenv("KUBECONFIG");
            }

        if (sConfig == null || sConfig.trim().length() == 0)
            {
            // no property so try the env variable to get the configuration file
            String sFileName = System.getProperty("user.home") + File.separator + ".kube/config";
            File   file      = new File(sFileName);

            if (file.exists() && file.isFile())
                {
                sConfig = sFileName;
                }
            }

        return sConfig == null || sConfig.trim().length() == 0 ? null : sConfig;
        }

    /**
     * Obtain the k8s namespace to use for tests.
     *
     * @return  the k8s namespace to use for tests
     */
    protected static String getK8sNamespace()
        {
        String sNamespace = getPropertyOrNull("k8s.namespace");

        if (sNamespace == null)
            {
            sNamespace = System.getenv("K8S_NAMESPACE");
            }

        if (sNamespace == null || sNamespace.trim().length() == 0)
            {
            sNamespace = System.getenv("CI_BUILD_ID");
            }

        return sNamespace == null || sNamespace.trim().length() == 0 ? "default" : sNamespace.trim();
        }

    /**
     * Obtain the namespaces for coherence.
     *
     * @return an array of target namespaces
     */
    protected static String[] getTargetNamespaces()
        {
        String sNamespaces = getPropertyOrNull("k8s.target.namespaces");

        if (sNamespaces == null)
            {
            sNamespaces = System.getenv("K8S_TARGET_NAMESPACES");
            }

        if (sNamespaces == null || sNamespaces.trim().length() == 0)
            {
            String sCINamespace = System.getenv("CI_BUILD_ID");
            if (sCINamespace != null)
                {
                sNamespaces = sCINamespace + "," + sCINamespace + "2";
                }
            }

        return sNamespaces == null || sNamespaces.trim().length() == 0 ? new String[] { null } : sNamespaces.trim().split("\\s*,\\s*");
        }

    /**
     * Obtain the default Helm set values.
     *
     * @return an array of helm set values
     */
    static String[] getDefaultHelmSetValues()
        {
        return (OP_IMAGE_PULL_POLICY == null) ? new String[0] :
            new String[] { "coherenceOperator.imagePullPolicy=" + OP_IMAGE_PULL_POLICY,
                           "coherence.imagePullPolicy=" + OP_IMAGE_PULL_POLICY,
                           "coherenceUtils.imagePullPolicy=" + OP_IMAGE_PULL_POLICY,
                           "userArtifacts.imagePullPolicy=" + OP_IMAGE_PULL_POLICY };
        }

    /**
     * Obtain the default Helm set values combined with additional set values.
     *
     * @param asSetValues  additional set values
     *
     * @return string array containing default and additional specified Helm set values
     */
    static String[] withDefaultHelmSetValues(String... asSetValues)
        {
        final Predicate<String> isNotBlank = s -> s != null && !s.trim().isEmpty();

        return Stream.concat(Arrays.stream(getDefaultHelmSetValues()), Arrays.stream(asSetValues))
            .filter(isNotBlank).toArray(String[]::new);
        }

    /**
     * Ensure that the k8s image pull secrets exist, creating it if
     * the {@link #CREATE_SECRET} flag is {@code true}.
     *
     * @param cluster      the k8s cluster
     *
     * @throws  AssertionError if the secret does not exist and
     *                         {@link #CREATE_SECRET} is {@code false}
     */
    protected static void ensureSecret(K8sCluster cluster)
        {
        ensureSecret(cluster, getK8sNamespace());
        }

    /**
     * Ensure that the k8s image pull secrets exist, creating it if
     * the {@link #CREATE_SECRET} flag is {@code true}.
     *
     * @param cluster      the k8s cluster
     * @param sNamespace   the k8s namespace
     *
     * @throws  AssertionError if the secret does not exist and
     *                         {@link #CREATE_SECRET} is {@code false}
     */
    static void ensureSecret(K8sCluster cluster, String sNamespace)
        {
        String sSecretName = getPropertyOrNull(PROP_K8S_PULL_SECRET);

        ensureSecret(cluster, sNamespace, sSecretName);
        }

    /**
     * Ensure that the k8s image pull secrets exist, creating it if
     * the {@link #CREATE_SECRET} flag is {@code true}.
     *
     * @param cluster      the k8s cluster
     * @param sNamespace   the k8s namespace
     * @param sSecretName  the name of the secret
     *
     * @throws  AssertionError if the secret does not exist and
     *                         {@link #CREATE_SECRET} is {@code false}
     */
    static void ensureSecret(K8sCluster cluster, String sNamespace, String sSecretName)
        {
        if (sSecretName == null || sSecretName.length() == 0)
            {
            return;
            }

        Arguments arguments = Arguments.empty();

        if (sNamespace != null && sNamespace.length() > 0)
            {
            arguments = arguments.with("--namespace=" + sNamespace);
            }

        arguments = arguments.with("get", "secret", sSecretName);

        int nExitCode = cluster.kubectlAndWait(arguments);

        if (nExitCode != 0 && CREATE_SECRET)
            {
            System.err.println("Creating secret " + sSecretName + " in namespace " + sNamespace);

            String    sRepo     = getPropertyOrDefault(PROP_DOCKER_REPO, "");
            String    sToken    = getPropertyOrNull(PROP_SECRET_CREDENTIALS);
            String    sUsername = getPropertyOrDefault(PROP_SECRET_USER, "moreply");
            String    sEmail    = getPropertyOrDefault(PROP_SECRET_EMAIL, "noreply@oracle.com");

            arguments = Arguments.of("create",
                                     "secret",
                                     "docker-registry",
                                     sSecretName,
                                     "--docker-server=" + sRepo,
                                     "--docker-password=" + sToken);

            if (sUsername != null && sUsername.length() > 0)
                {
                arguments = arguments.with("--docker-username=" + sUsername);
                }

            if (sEmail != null && sEmail.length() > 0)
                {
                arguments = arguments.with("--docker-email=" + sEmail);
                }

            if (sNamespace != null && sNamespace.length() > 0)
                {
                arguments = arguments.with("--namespace=" + sNamespace);
                }

            assertThat(cluster.kubectlAndWait(arguments), is(0));
            }
        }

    /**
     * Remove the image pull secrets if the {@link #CREATE_SECRET} flag is {@code true}.
     *
     * @param cluster     the k8s cluster
     */
    protected static void cleanupPullSecrets(K8sCluster cluster)
        {
        cleanupPullSecrets(cluster, getK8sNamespace());
        }

    /**
     * Remove the image pull secrets if the {@link #CREATE_SECRET} flag is {@code true}.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the k8s namespace
     */
    static void cleanupPullSecrets(K8sCluster cluster, String sNamespace)
        {
        if (CREATE_SECRET)
            {
            String sSecretName = getPropertyOrNull(PROP_K8S_PULL_SECRET);

            if (sSecretName != null && sSecretName.length() > 0)
                {
                Arguments arguments = Arguments.empty();

                if (sNamespace != null && sNamespace.length() > 0)
                    {
                    arguments = arguments.with("--namespace=" + sNamespace);
                    }

                arguments = arguments.with("delete", "secret", sSecretName);

                cluster.kubectlAndWait(arguments);
                }
            }
        }

    /**
     * Install coherence with given helm values and target namespace.
     *
     * @param cluster       the k8s cluster
     * @param sNamespace    the namespace to install into
     * @param sHelmValues   Helm values file
     *
     * @return the name of the Helm release
     *
     * @throws Exception if installation fails
     */
    protected static String installCoherence(K8sCluster cluster,
                                      String     sNamespace,
                                      String     sHelmValues) throws Exception
        {
        String[] asReleases = installCoherence(cluster, new String[]{sNamespace}, sHelmValues);
        return asReleases[0];
        }

    /**
     * Install coherence with given helm values and target namespace.
     *
     * @param cluster      the k8s cluster
     * @param sNamespace   the namespace to install into
     * @param sHelmValues  the Helm values file
     * @param asSetValues  the array of Helm set values
     *
     * @return the name of the Helm release
     *
     * @throws Exception if installation fails
     */
    protected static String installCoherence(K8sCluster cluster,
                                      String     sNamespace,
                                      String     sHelmValues,
                                      String...  asSetValues) throws Exception
        {
        String[] asReleases = installCoherence(cluster, new String[]{sNamespace}, sHelmValues, asSetValues);
        return asReleases[0];
        }

    /**
     * Install coherence with given helm values and target namespaces.
     *
     * @param cluster          the k8s cluster
     * @param asCohNamespaces  an array of namespace for installing coherence
     * @param sHelmValues      Helm values file
     *
     * @return an array of release names
     */
    protected static String[] installCoherence(K8sCluster cluster,
                                        String[]   asCohNamespaces,
                                        String     sHelmValues) throws Exception
        {
        return installCoherence(cluster, asCohNamespaces, sHelmValues, new String[0]);
        }

    /**
     * Install coherence with given helm values and target namespaces.
     *
     * @param cluster          the k8s cluster
     * @param asCohNamespaces  an array of namespace for installing coherence
     * @param sHelmValues      Helm values file
     * @param asSetValues      the array of Helm set values
     *
     * @return an array of release names
     */
    protected static String[] installCoherence(K8sCluster cluster,
                                        String[]   asCohNamespaces,
                                        String     sHelmValues,
                                        String...  asSetValues) throws Exception
        {
        String[] asReleases = new String[asCohNamespaces.length];

        for (int i = 0; i < asReleases.length; i++)
            {
            String sNamespace = asCohNamespaces[i];

            try
                {
                String[] asActualSetValues;

                if (asSetValues == null || asSetValues.length == 0)
                    {
                    asActualSetValues = getDefaultHelmSetValues();
                    }
                else
                    {
                    asActualSetValues = Stream.concat(Arrays.stream(getDefaultHelmSetValues()), Arrays.stream(asSetValues))
                                              .toArray(String[]::new);
                    }

                asReleases[i] = installChart(cluster, COHERENCE_HELM_CHART_NAME, COHERENCE_HELM_CHART_URL, sNamespace, sHelmValues, asActualSetValues);
                }
            catch(Throwable throwable)
                {
                System.err.println("Install Coherence in namespace " + sNamespace + " failed.");
                // clean up previous installed Coherence
                for (String sReleaseName : asReleases)
                    {
                    if (sReleaseName != null)
                        {
                        try
                            {
                            System.err.println("Clean up Helm Coherence install '" + sReleaseName);
                            cleanupHelmReleases(sReleaseName);
                            cleanupPersistentVolumeClaims(cluster, sReleaseName, sNamespace);
                            }
                        catch (Throwable t)
                            {
                            System.err.println("Error in clean up Helm release '" + sReleaseName + "': " + t);
                            }
                        }
                    }

                if (throwable instanceof Exception)
                    {
                    throw (Exception) throwable;
                    }
                else
                    {
                    throw new Exception(throwable);
                    }
                }
            }

        return asReleases;
        }

    /**
     * Assert that the Coherence JMX Deployment is up.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  namespace for installing coherence
     * @param sRelease    the Helm release name
     */
    protected void assertCoherenceJMX(K8sCluster cluster, String sNamespace, String sRelease)
        throws Exception
        {
        String selector = getCoherenceJmxDeploymentSelector(sRelease);
        assertDeploymentReady(cluster, sNamespace, selector, true);

        ensureJMXServerPartOfCluster(cluster, sNamespace, sRelease);
        }

    protected void ensureJMXServerPartOfCluster(K8sCluster k8sCluster, String sNamespace, String sRelease)
        throws Exception
        {
        try
            {
            Eventually.assertThat("Unable to confirm JMX server is part of cluster",
                invoking(this).getClusterSizeViaJMX(k8sCluster, sNamespace, sRelease), greaterThanOrEqualTo(2));
            }
        catch (Throwable t)
            {
            String sCoherenceJmxSelector = getCoherenceJmxPodSelector(sRelease);
            List<String> listJmxPods = getPods(k8sCluster, sNamespace, sCoherenceJmxSelector);
            if (listJmxPods.size() > 0)
                {
                dumpPodLog(k8sCluster, sNamespace, listJmxPods.get(0));
                }
            throw t;
            }
        }

    /**
     * Assert that coherence is up with given helm values, target namespaces and persistence.
     *
     * @param cluster        the k8s cluster
     * @param sCohNamespace  namespace for installing coherence
     * @param sRelease       the Helm release name
     */
    protected static void assertCoherence(K8sCluster cluster, String sCohNamespace, String sRelease)
        {
        assertCoherence(cluster, new String[] { sCohNamespace }, new String[] { sRelease });
        }

    /**
     * Assert that coherence is up with given helm values, target namespaces and persistence.
     *
     * @param cluster          the k8s cluster
     * @param asCohNamespaces  an array of namespace for installing coherence
     * @param asReleases       an array of Helm release names
     */
    protected static void assertCoherence(K8sCluster cluster, String[] asCohNamespaces, String[] asReleases)
        {
        assertThat("Some of Helm installation are failed", asReleases.length, is(asCohNamespaces.length));
        for (int i = 0; i < asCohNamespaces.length; i++)
            {
            String sNamespace = asCohNamespaces[i];
            String sRelease   = asReleases[i];
            String sSelector  = getCoherenceSelector(sRelease);

            assertStatefulSetReady(cluster, sNamespace, sSelector, true);
            }
        }

    /**
     * Delete given Coherence for target namespace.
     *
     * @param cluster        the k8s cluster
     * @param sCohNamespace  the coherence target namespace
     * @param sRelease       the release name of coherence
     * @param fPersistence   indicates whether coherence cache data is persisted
     */
    protected void deleteCoherence(K8sCluster cluster, String sCohNamespace, String sRelease, boolean fPersistence)
        {
        deleteCoherence(getClass(), cluster, new String[] { sCohNamespace }, new String[] { sRelease }, fPersistence);
        }

    /**
     * Delete given Coherence for target namespace.
     *
     * @param cluster        the k8s cluster
     * @param sCohNamespace  the coherence target namespace
     * @param sRelease       the release name of coherence
     * @param fPersistence   indicates whether coherence cache data is persisted
     */
    protected static void deleteCoherence(Class clsTest, K8sCluster cluster, String sCohNamespace, String sRelease, boolean fPersistence)
        {
        deleteCoherence(clsTest, cluster, new String[] { sCohNamespace }, new String[] { sRelease }, fPersistence);
        }

    /**
     * Delete given Coherence for target namespace.
     *
     * @param cluster          the k8s cluster
     * @param asCohNamespaces  an array of coherence target namespaces
     * @param asReleases        an array of the release names of coherence
     * @param fPersistence     indicates whether coherence cache data is persisted
     */
    protected void deleteCoherence(K8sCluster cluster, String[] asCohNamespaces, String[] asReleases, boolean fPersistence)
        {
        deleteCoherence(getClass(), cluster, asCohNamespaces, asReleases, fPersistence);
        }

    /**
     * Delete given Coherence for target namespace.
     *
     * @param cluster          the k8s cluster
     * @param asCohNamespaces  an array of coherence target namespaces
     * @param asReleases        an array of the release names of coherence
     * @param fPersistence     indicates whether coherence cache data is persisted
     */
    protected static void deleteCoherence(Class clsTest, K8sCluster cluster, String[] asCohNamespaces, String[] asReleases, boolean fPersistence)
        {
        for (int i = 0; i < asReleases.length; i++)
            {
            String sRelease   = asReleases[i];
            String sNamespace = asCohNamespaces[i];
            String sSelector  = getCoherencePodSelector(sRelease);
            capturePodLogs(clsTest, cluster, getK8sNamespace(), sSelector, "coherence", "user-artifacts", "coherence-k8s-utils");
            cleanupHelmReleases(sRelease);
            cleanupPersistentVolumeClaims(cluster, sRelease, sNamespace);

            try
                {
                assertStatefulSetReady(cluster, sNamespace, getCoherenceSelector(sRelease), false);
                }
            catch(Throwable t)
                {
                System.err.println("Fail to wait for deleting Coherence release '" + sRelease + "'.");
                }
            }
        }

    /**
     * Retrieve http result as a queue of String
     *
     * @param cluster      the k8s cluster
     * @param sPod         the pod name
     * @param sHttpMethod  the http method name
     * @param sHost        the http host
     * @param nPort        the http port
     * @param sPath        the http path
     *
     * @return http result as a queue of String
     */
    protected Queue<String> processHttpRequest(K8sCluster cluster, String sPod, String sHttpMethod, String sHost, int nPort, String sPath)
        {
        // Workaround intermittent kubectl exec returning non-zero due to a communication failure by retrying
        final int          MAX_RETRY     = 5;
        AssertionException lastException = null;

        for (int i=0; i < MAX_RETRY; i++)
            {
            try
                {
                return processHttpRequest(cluster, sPod, sHttpMethod, sHost, nPort, sPath, /*fConsole*/ true);
                }
            catch (AssertionException e)
                {
                lastException = e;
                try
                    {
                    // backoff before retry
                    Thread.sleep(500);
                    }
                catch (InterruptedException e1)
                    {
                    }
                }
            }
        throw lastException;
        }

    /**
     * Retrieve http result as a queue of String by executing curl in specified sPod
     *
     * @param cluster      the k8s cluster
     * @param sPod         the pod name to execute curl request in
     * @param sHttpMethod  the http method name
     * @param sHost        the http host
     * @param nPort        the http port
     * @param sPath        the http path
     * @param fConsole     if false, do not print result to console.  (for large results needing post processing)
     *
     * @return http result as a queue of String
     */
    protected Queue<String> processHttpRequest(K8sCluster cluster, String sPod, String sHttpMethod, String sHost, int nPort, String sPath, boolean fConsole)
        {
        CapturingApplicationConsole console = new CapturingApplicationConsole();

        Arguments arguments = Arguments.of("exec", "-i", sPod);

        String sNamespace = getK8sNamespace();
        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        arguments = arguments.with("--", "curl", "-X", sHttpMethod, sHost + ":" + nPort + sPath);

        int nExitCode = cluster.kubectlAndWait(arguments, LaunchLogging.disabled(), Console.of(console));

        if (fConsole)
            {
            HelmUtils.logConsoleOutput("exec-i-http", console);
            }

        assertThat("kubectl returned non-zero exit code", nExitCode, is(0));

        return console.getCapturedOutputLines();
        }

    /**
     * Retrieve http result as a queue of String by executing curl in a temporary pod created to run curl in k8 environment.
     *
     * @param cluster      the k8s cluster
     * @param sHttpMethod  the http method name
     * @param sHost        the http host
     * @param nPort        the http port
     * @param sPath        the http path
     * @param fConsole     if false, do not print result to console.  (for large results needing post processing)
     *
     * @return http result as a queue of String
     */
    protected Queue<String> processHttpRequest(K8sCluster cluster, String sHttpMethod, String sHost, int nPort, String sPath, boolean fConsole)
        {
        String                      sNamespace = getK8sNamespace();
        CapturingApplicationConsole console    = new CapturingApplicationConsole();

        Arguments arguments = Arguments.of("run", "--generator=run-pod/v1", "--namespace=" + sNamespace, "--image=library/oraclelinux:7.6", "--restart=Never", "--rm", "-it", "curl");

        arguments = arguments.with("--", "curl", "-X", sHttpMethod, sHost + ":" + nPort + sPath);

        int nExitCode = cluster.kubectlAndWait(arguments, LaunchLogging.disabled(), Console.of(console));

        if (nExitCode != 0)
            {
            // recover from failure pod/curl already exists.
            deleteResource(cluster, "pod/curl", sNamespace);
            }

        if (fConsole)
            {
            HelmUtils.logConsoleOutput("curl-http-request", console);
            }

        return console.getCapturedOutputLines();
        }

    /**
     * Obtain the value of the System property or {@code null}
     * if the property is not set or is set to empty String.
     *
     * @param sPropertyName  the name of the property
     *
     * @return  the value of the property or {@code null} if the
     *          property is not set or is an empty String.
     */
    private static String getPropertyOrNull(String sPropertyName)
        {
        String sValue = System.getProperty(sPropertyName);

        return sValue == null || sValue.trim().length() == 0 ? null : sValue;
        }

    /**
     * Obtain the value of the System property or the default
     * if the property is not set or is set to empty String.
     *
     * @param sPropertyName  the name of the property
     * @param sDefault       the default value
     *
     * @return  the value of the property or the default if the
     *          property is not set or is an empty String.
     */
    private static String getPropertyOrDefault(String sPropertyName, String sDefault)
        {
        String sValue = System.getProperty(sPropertyName);

        return sValue == null || sValue.trim().length() == 0 ? sDefault : sValue;
        }

    /**
     * Run a kubectl port-forward process to forward a port to a Coherence Pod.
     *
     * @param k8sCluster  the {@link K8sCluster} to use
     * @param sNamespace  the namespace that the Pod is in
     * @param sRelease    the name of the Helm release that installed the Pod
     * @param nPort       the Pod's port to forward
     * @return  an {@link Application} running the port-forward process
     * @throws Exception  if the port-forward application fails to start
     */
    protected static Application portForwardCoherencePod(K8sCluster k8sCluster, String sNamespace, String sRelease, int nPort) throws Exception
        {
        String sSelector = getCoherencePodSelector(sRelease);
        return portForward(k8sCluster, sNamespace, sSelector, nPort);
        }

    /**
     * Run a kubectl port-forward process to forward a port to a Pod.
     *
     * @param k8sCluster  the {@link K8sCluster} to use
     * @param sNamespace  the namespace that the Pod is in
     * @param sSelector   the name of the k8s selector to use to locate a Pod
     * @param nPort       the Pod's port to forward
     * @return  an {@link Application} running the port-forward process
     * @throws Exception  if the port-forward application fails to start
     */
    protected static Application internalPortForward(K8sCluster k8sCluster, String sNamespace, String sSelector, int nPort) throws Exception
        {
        List<String>                               listPod        = getPods(k8sCluster, sNamespace, sSelector);
        String                                     sPod           = listPod.get(0);
        EventsApplicationConsole.CountDownListener listener       = new EventsApplicationConsole.CountDownListener(1);
        Predicate<String>                          predicate      = s -> s.contains("Forwarding from 127.0.0.1");
        TestLogs.ConsoleBuilder                    consoleBuilder = s_testLogs.builder();

        // Add a listener to the log that will be triggered when it sees the proxy start message
        consoleBuilder.clearListeners();
        consoleBuilder.addStdOutListener(predicate, listener);
        consoleBuilder.addStdErrListener(predicate, listener);

        Application application = HelmUtils.portForward(k8sCluster, sPod, sNamespace, nPort, consoleBuilder);

        assertThat("Failed to see port forward start message",
                   listener.await(1, TimeUnit.MINUTES), is(true));

        return application;
        }

    /**
     * Run a kubectl port-forward process to forward a port to a Pod.
     *
     * @param k8sCluster  the {@link K8sCluster} to use
     * @param sNamespace  the namespace that the Pod is in
     * @param sSelector   the name of the k8s selector to use to locate a Pod
     * @param nPort       the Pod's port to forward
     *
     * @return  an {@link Application} running the port-forward process
     *
     * @throws Exception  if the port-forward application fails to start
     */
    protected static Application portForward(K8sCluster k8sCluster, String sNamespace, String sSelector, int nPort) throws Exception
        {
        // Workaround intermittent kubectl port-forward failure by retrying
        Throwable lastThrowable = null;
        final int MAX_RETRY     = 8;
        for (int i = 0; i < MAX_RETRY; i++)
            {
            try
                {
                return internalPortForward(k8sCluster, sNamespace, sSelector, nPort);
                }
            catch (Throwable t)
                {
                lastThrowable = t;
                try
                    {
                    // backoff before retry
                    Thread.sleep(500);
                    }
                catch (InterruptedException e1)
                    {
                    }
                }
            }
        throw new WrapperException(lastThrowable);
        }

    /**
     * Perform a JMX query using JMXMP transport.
     * <p>
     * This method will run a kubectl port-forward process to expose the JMX port
     * then perform the JMX query and close the port-forward process.
     *
     * @param k8sCluster    the {@link K8sCluster} to use
     * @param sNamespace    the namespace that the JMX Pod is in
     * @param sRelease      the name of the Helm release that installed the JMX Pod
     * @param objectName    the ObjectName of the MBean to query
     * @param attributeName the name of the attribute to get
     * @param <T>           the type of the attribute value
     *
     * @return the value of the attribute from the MBean
     * @throws Exception if the query fails
     */
    public <T> T jmxQuery(K8sCluster k8sCluster, String sNamespace, String sRelease, String objectName, String attributeName) throws Exception
        {
        String sSelector = getCoherenceJmxPodSelector(sRelease);
        try (Application application = portForward(k8sCluster, sNamespace, sSelector, 9099))
            {
            PortMapping portMapping = application.get(PortMapping.class);
            int         nPort       = portMapping.getPort().getActualPort();

            return jmxQuery("127.0.0.1", nPort, objectName, attributeName);
            }
        }

    /**
     * Perform a JMX query using JMXMP transport.
     *
     * @param hostName      the address that the jmxmp JMX server is bound to
     * @param port          the port that the jmxmp JMX server is listening on
     * @param objectName    the ObjectName of the MBean to query
     * @param attributeName the name of the attribute to get
     * @param <T>           the type of the attribute value
     *
     * @return the value of the attribute from the MBean
     * @throws Exception if the query fails
     */
    @SuppressWarnings("unchecked")
    public <T> T jmxQuery(String hostName, int port, String objectName, String attributeName) throws Exception
        {
        JMXServiceURL jmxURL = new JMXServiceURL("jmxmp", hostName, port);

        try (JMXConnector jmxc = JMXConnectorFactory.connect(jmxURL, null))
            {
            MBeanServerConnection serverConnection = jmxc.getMBeanServerConnection();

            return (T) serverConnection.getAttribute(new ObjectName(objectName), attributeName);
            }
        }

    /**
     * Invoke a method on a JMX MBean.
     * <p>
     * This method will run a kubectl port-forward process to expose the JMX port
     * then perform the JMX invoke and close the port-forward process.
     *
     * @param k8sCluster    the {@link K8sCluster} to use
     * @param sNamespace    the namespace that the JMX Pod is in
     * @param sRelease      the name of the Helm release that installed the JMX Pod
     * @param objectName  the ObjectName of the MBean to invoke operation
     * @param methodName  the operation to invoke
     * @param params      an array containing the parameters to be set when
     *                    the operation is invoked
     * @param signature   an array containing the signature of the operation,
     *                    an array of class names in the format returned by
     *                    {@link Class#getName()}. The class objects will be
     *                    loaded using the same class loader as the one used
     *                    for loading the MBean on which the operation was invoked
     *
     * @param <T>         the type of the attribute value
     *
     * @return the result of the MBean methodName operation
     * @throws Exception if the query fails
     */
    public <T> T jmxInvoke(K8sCluster k8sCluster,
                          String sNamespace,
                          String sRelease,
                          String objectName,
                          String methodName,
                          Object[] params,
                          String[] signature) throws Exception
        {
        String sSelector = getCoherenceJmxPodSelector(sRelease);
        try (Application application = portForward(k8sCluster, sNamespace, sSelector, 9099))
            {
            PortMapping portMapping = application.get(PortMapping.class);
            int         nPort       = portMapping.getPort().getActualPort();

            return jmxInvoke("127.0.0.1", nPort, objectName, methodName, params, signature);
            }
        }

    /**
     * Invoke a method on a JMX MBean.
     *
     * @param hostName    the address that the jmxmp JMX server is bound to
     * @param port        the port that the jmxmp JMX server is listening on
     * @param objectName  the ObjectName of the MBean to query
     * @param methodName  the name of the attribute to get
     * @param params      an array containing the parameters to be set when
     *                    the operation is invoked
     * @param signature   an array containing the signature of the operation,
     *                    an array of class names in the format returned by
     *                    {@link Class#getName()}. The class objects will be
     *                    loaded using the same class loader as the one used
     *                    for loading the MBean on which the operation was invoked
     *
     * @param <T>         the type of the invocation result
     *
     * @return the value of the attribute from the MBean
     *
     * @throws Exception if the query fails
     */
    @SuppressWarnings("unchecked")
    public <T> T jmxInvoke(String hostName, int port, String objectName, String methodName,
                           Object[] params, String[] signature) throws Exception
        {
        JMXServiceURL jmxURL = new JMXServiceURL("jmxmp", hostName, port);

        try (JMXConnector jmxc = JMXConnectorFactory.connect(jmxURL, null))
            {
            MBeanServerConnection serverConnection = jmxc.getMBeanServerConnection();

            return (T) serverConnection.invoke(new ObjectName(objectName), methodName, params, signature);
            }
        }

    public Integer getClusterSizeViaJMX(K8sCluster k8sCluster, String sNamespace, String sRelease) throws Exception
        {
        String sClusterMBean = "Coherence:type=Cluster";
        String sAttribute = "ClusterSize";
        return jmxQuery(k8sCluster, sNamespace, sRelease, sClusterMBean, sAttribute);
        }

    /**
     * Return true if current Coherence version is greater than or equal to minimal version.
     *
     * @param sMinimalVersion  minimal version of current Coherence
     *
     * @return true iff current Coherence version is equal to or greater than minimal version.
     */
    protected boolean versionCheck(String sMinimalVersion)
        {
        return versionCheck(COHERENCE_VERSION, sMinimalVersion);
        }

    /**
     * Return true if specified Coherence version is greater than or equal to minimal version.
     *
     * @param sVersion         target Coherence version
     * @param sMinimalVersion  minimal version of target Coherence
     *
     * @return true iff target Coherence version is equal to or greater than minimal version.
     */
    protected boolean versionCheck(String sVersion, String sMinimalVersion)
        {
        return CoherenceVersion.versionCheck(sVersion, sMinimalVersion);
        }

    // ----- inner class  ---------------------------------------------------

    /**
     * The class provides a stub class for com.oracle.bedrock.deferred.DeferredHelper.invoking.
     */
    private static class StubHelmChartTest extends BaseHelmChartTest
        {
        }

    // ----- data members ---------------------------------------------------

    /**
     * The property for helm install max retry attempt. Default to 3.
     */
    public static final String HELM_INSTALL_MAX_RETRY = System.getProperty("helm.install.maxRetry", "3");

    /**
     * The default name for the Coherence Operator Helm chart to test.
     */
    public static final String DEFAULT_OPERATOR_HELM_CHART = "coherence-operator";

    /**
     * The default name for the Coherence Helm chart to test.
     */
    public static final String DEFAULT_COHERENCE_HELM_CHART = "coherence";

    /**
     * The System property to use to set the name of the Helm package to test.
     */
    public static final String PROP_OPERATOR_HELM_PACKAGE = "operator.helm.chart.package";

    /**
     * The System property to use to set the name of the Helm package to test.
     */
    public static final String PROP_COHERENCE_HELM_PACKAGE = "coherence.helm.chart.package";

    /**
     * The System property to use to set the name of the Operator Helm chart to test.
     */
    public static final String PROP_OPERATOR_HELM_CHART = "operator.helm.chart.name";

    /**
     * The System property to use to set the name of the Coherence Helm chart to test.
     */
    public static final String PROP_COHERENCE_HELM_CHART = "coherence.helm.chart.name";

    /**
     * The System property to use to obtain the name of the optional k8s docker-registry secret.
     */
    public static final String PROP_K8S_PULL_SECRET = "k8s.image.pull.secret";

    /**
     * The name of the System property to use to determine whether to create and destroy the k8s test namespace.
     */
    public static final String PROP_CREATE_NAMESPACE = "k8s.create.namespace";

    /**
     * The name of the System property to use to determine whether to create and destroy the k8s pull secret.
     */
    public static final String PROP_CREATE_SECRET = "k8s.create.secret";

    /**
     * The name of the System property to use to determine the location of kubectl.
     */
    public static final String PROP_KUBECTL = "k8s.kubectl";

    /**
     * The name of the System property to use to determine the k8s secret user.
     */
    public static final String PROP_SECRET_USER = "k8s.secret.user";

    /**
     * The name of the System property to use to determine the k8s secret credentials.
     */
    public static final String PROP_SECRET_CREDENTIALS = "k8s.secret.credentials";

    /**
     * The name of the System property to use to determine the k8s secret email.
     */
    public static final String PROP_SECRET_EMAIL = "k8s.secret.email";

    /**
     * The name of the System property to use to determine the docker repo name.
     */
    public static final String PROP_DOCKER_REPO = "docker.repo";

    /**
     * The name of the System property to use to determine the operator image pull policy.
     */
    public static final String PROP_OP_IMAGE_PULL_POLICY = "op.image.pull.policy";

    /**
     * The name of the Operator Helm chart package being tested.
     */
    private static final String OPERATOR_HELM_CHART_PACKAGE = System.getProperty(PROP_OPERATOR_HELM_PACKAGE);

    /**
     * The URL of the Operator Helm chart package.
     */
    protected static final URL OPERATOR_HELM_CHART_URL = toURL(OPERATOR_HELM_CHART_PACKAGE);

    /**
     * The name of the Helm chart to test.
     */
    protected static final String OPERATOR_HELM_CHART_NAME = getPropertyOrDefault(PROP_OPERATOR_HELM_CHART, DEFAULT_OPERATOR_HELM_CHART);

    /**
     * The name of the Coherence Helm chart package being tested.
     */
    private static final String COHERENCE_HELM_CHART_PACKAGE = System.getProperty(PROP_COHERENCE_HELM_PACKAGE);

    /**
     * The URL of the Coherence Helm chart package.
     */
    protected static final URL COHERENCE_HELM_CHART_URL = toURL(COHERENCE_HELM_CHART_PACKAGE);

    /**
     * The name of the Helm chart to test.
     */
    protected static final String COHERENCE_HELM_CHART_NAME = getPropertyOrDefault(PROP_COHERENCE_HELM_CHART, DEFAULT_COHERENCE_HELM_CHART);

    /**
     * The operator image pull policy.
     */
    protected static final String OP_IMAGE_PULL_POLICY = getPropertyOrNull(PROP_OP_IMAGE_PULL_POLICY);

    /**
     * The name of the optional k8s docker-registry secret.
     */
    protected static final String K8S_IMAGE_PULL_SECRET = getPropertyOrNull(PROP_K8S_PULL_SECRET);

    /**
     * Flag indicating whether to create and destroy the test namespace.
     */
    public static final boolean CREATE_NAMESPACE = Boolean.getBoolean(PROP_CREATE_NAMESPACE);

    /**
     * Flag indicating whether to create and destroy the test namespace.
     */
    public static final boolean CREATE_SECRET = Boolean.getBoolean(PROP_CREATE_SECRET);

    /**
     * The location of the kubectl.
     */
    public static final String KUBECTL = getPropertyOrNull(PROP_KUBECTL);

    /**
     * The name of the Coherence container in the Coherence Pod.
     */
    public static final String COHERENCE_CONTAINER_NAME = "coherence";

    /**
     * The component name of Coherence Kubernetes Operator.
     */
    public static final String COHERENCE_K8S_OPERATOR = "coherence-operator";

    /**
     * The parameters value to use in a JMX MBean invocation for a method that takes no parameters.
     */
    public static final Object[] NO_JMX_PARAMS = new Object[0];

    /**
     * The signature value to use in a JMX MBean invocation for a method that takes no parameters.
     */
    public static final String[] EMPTY_JMX_SIGNATURE = new String[0];

    /**
     * The version (tag) for the latest Coherence image version being tested.
     */
    public static final String COHERENCE_VERSION = System.getProperty("coherence.docker.version");

    /**
     * A JUnit class rule to create temporary files and folders.
     */
    @ClassRule
    public static TemporaryFolder s_temp = new TemporaryFolder();

    /**
     * A JUnit class rule to obtain the current test name.
     */
    @Rule
    public final TestName m_testName = new TestName();

    /**
     * Customer Resource Definitions (CRDs) created when prometheus-operator is installed.
     */
    public static final String[] PROMETHEUS_OPERATOR_CRDS =
        {
        "prometheuses.monitoring.coreos.com",
        "prometheusrules.monitoring.coreos.com",
        "servicemonitors.monitoring.coreos.com",
        "alertmanagers.monitoring.coreos.com"
        };

    /**
     * The mapper to parse json.
     */
    public static final ObjectMapper MAPPER = new ObjectMapper();

    /**
     * A Bedrock JUnit rule to configure {@link ApplicationConsole}s for tests.
     */
    @ClassRule
    @Rule
    public static final TestLogs s_testLogs = new TestLogs();

    @ClassRule
    @Rule
    public static final Kubernetes.TestName s_k8sTestName = new Kubernetes.TestName();

    /**
     * The Helm command template.
     */
    protected static HelmCommand.Template s_helm = Helm.template()
                       .home(getPropertyOrNull("helm.home"))
                       .host(getPropertyOrNull("helm.tiller.host"))
                       .kubeConfig(getK8sConfig())
                       .kubeContext(getPropertyOrNull("k8s.context"))
                       .tillerNamespace(getPropertyOrNull("helm.tiller.namespace"));

    /**
     * A stub {@link BaseHelmChartTest} to use to call methods via {@link Eventually#assertThat(String, Object, Matcher)}
     */
    private static final BaseHelmChartTest STUB = new StubHelmChartTest();
    }
