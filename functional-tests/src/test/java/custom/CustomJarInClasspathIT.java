/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static helm.HelmUtils.HELM_TIMEOUT;
import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

import java.io.File;
import java.net.URL;
import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.TimeUnit;

import com.oracle.bedrock.runtime.console.SystemApplicationConsole;
import com.tangosol.util.Resources;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Ignore;
import org.junit.Test;

import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.k8s.helm.HelmUpgrade;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.testsupport.deferred.Eventually;

import helm.BaseHelmChartTest;
import helm.HelmUtils;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import util.CustomParameterizedRunner;

import static org.hamcrest.CoreMatchers.is;

/**
 * @author cp
 */
@RunWith(Parameterized.class)
@Parameterized.UseParametersRunnerFactory(CustomParameterizedRunner.Factory.class)
@Ignore
public class CustomJarInClasspathIT extends BaseHelmChartTest
    {
    // ----- test lifecycle --------------------------------------------------

    /**
     * Create the test parameters (the versions of the Coherence image to test).
     *
     * @return  the test parameters
     */
    @Parameterized.Parameters(name = "{0}")
    public static Collection<Object[]> parameters()
        {
        return Arrays.asList(new Object[][] {
                {COHERENCE_VERSION},
                {"12.2.1.1"},
                {"12.2.1.2"},
                {"12.2.1.3"}
            });
        }

    @BeforeClass
    public static void setup()
        {
        assertPreconditions(s_k8sCluster);
        ensureNamespace(s_k8sCluster);
        ensureSecret(s_k8sCluster);
        }

    @AfterClass
    public static void cleanup()
        {
        cleanupPullSecrets(s_k8sCluster);
        cleanupNamespace(s_k8sCluster);

        for (String sNamespace : getTargetNamespaces())
            {
            cleanupNamespace(s_k8sCluster, sNamespace);
            }
        }

    @CustomParameterizedRunner.AfterParmeterizedRun
    public void cleanUpCoherence()
        {
        if (m_sRelease != null)
            {
            deleteCoherence(s_k8sCluster, getK8sNamespace(), m_sRelease, PERSISTENCE);
            }
        }

    // ----- constructors ---------------------------------------------------

    /**
     * Create a test class instance.
     *
     * @param sCoherenceTag  the tag to use when pulling the Coherence image for this test
     */
    public CustomJarInClasspathIT(String sCoherenceTag)
        {
        m_sCoherenceTag = sCoherenceTag;
        }

    // ----- test methods ---------------------------------------------------

    /**
     * Run the test using supplied list of values files.
     *
     * @throws Exception  if the test fails
     */
    @Test
    public void testCoherenceWithUserSuppliedJarInClasspath() throws Exception
        {
        String sNamespace      = getK8sNamespace();
        String sValuesOriginal = "values/helm-values-coh-user-artifacts.yaml";
        String sValuesUpgrade  = "values/helm-values-coh-user-artifacts-upgrade.yaml";
        String sCoherenceImage = COHERENCE_IMAGE_PREFIX + m_sCoherenceTag;

        m_sRelease = installCoherence(s_k8sCluster,
                                      sNamespace,
                                      sValuesOriginal,
                                      "coherence.image=" + sCoherenceImage);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        assertCoherenceService(s_k8sCluster, sNamespace, m_sRelease);

        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPods           = getPods(s_k8sCluster, sNamespace, sCoherenceSelector);

        assertThat(listPods.size(), is(3));

        System.err.println("Waiting for Coherence Pods");

        for (String sPod : listPods)
            {
            System.err.println("Waiting for Coherence Pod " + sPod + "...");
            Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, sNamespace, sPod),
                                  is(true), Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS));
            }

        System.err.println("Coherence Pods started");

        try
            {
            installClient(s_k8sCluster, CLIENT1, sNamespace, m_sRelease);

            installClient(s_k8sCluster, CLIENT2, sNamespace, m_sRelease);

            System.err.println("Waiting for Client-1 initial state ...");
            Eventually.assertThat(invoking(this).isRequiredClientStateReached(s_k8sCluster, sNamespace, CLIENT1),
                                  is(true),
                                  Eventually.within(TIMEOUT, TimeUnit.SECONDS));

            System.err.println("****************************************************************");
            System.err.println("********[STARTING UPGRADE OF " + m_sRelease + "]***********");
            System.err.println("****************************************************************");

            URL urlValues = Resources.findFileOrResource(sValuesUpgrade, null);

            upgradeUserArtifactsInCoherenceClasspath(m_sRelease, urlValues.getPath());

            System.err.println("Waiting for required Client-1 state after upgrade ...");
            Eventually.assertThat(invoking(this).isRequiredClientStateReachedAfterUpgrade(s_k8sCluster, sNamespace, CLIENT1),
                                  is(true),
                                  Eventually.within(TIMEOUT, TimeUnit.SECONDS));
            }
        finally
            {
            dumpPodLog(s_k8sCluster, sNamespace, CLIENT1, null);
            deleteClients();
            }
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Install coherence with upgraded user supplied artifacts
     *
     * @param sRelease     the release that is being upgraded to newer version of artifacts
     * @param sHelmValues  Helm values file
     */
    private void upgradeUserArtifactsInCoherenceClasspath(String sRelease, String sHelmValues) throws Exception
        {
        File fileTempDir = s_temp.newFolder();

        System.err.println("Extracting Helm chart " + COHERENCE_HELM_CHART_URL + " into " + fileTempDir);

        HelmUtils.extractTarGZ(fileTempDir, COHERENCE_HELM_CHART_URL);

        assertHelmLint(fileTempDir, COHERENCE_HELM_CHART_NAME);

        HelmUpgrade cohUpgrade = s_helm.upgrade(sRelease, fileTempDir + "/" + COHERENCE_HELM_CHART_NAME)
                                       .values(sHelmValues);

        if (K8S_IMAGE_PULL_SECRET != null && K8S_IMAGE_PULL_SECRET.trim().length() > 0)
            {
            cohUpgrade = cohUpgrade.set("imagePullSecrets={" + K8S_IMAGE_PULL_SECRET + "}");
            }

        int nExitCode = cohUpgrade.executeAndWait();

        if (nExitCode == 0)
            {
            captureInstalledPodLogs(getDefaultCluster(), getK8sNamespace(), sRelease);
            }

        assertThat("Helm upgrade returned non-zero exit code.", nExitCode, is(0));
        }

    /**
     * Check for required client state after the upgrade.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the namespace name
     * @param sClientPod  the pod name
     *
     * @return {@code true} if required client state is reached.
     */
    public boolean isRequiredClientStateReachedAfterUpgrade(K8sCluster cluster, String sNamespace, String sClientPod)
        {
        try
            {
            Queue<String> sLogs = getPodLog(cluster, sNamespace, sClientPod, null);
            return sLogs.stream().anyMatch(l -> l.contains("Cache Value Before Cloud EntryProcessor: AWS"))
                    && sLogs.stream().anyMatch(l -> l.contains("Cache Value After Cloud EntryProcessor: OCI"));
            }
        catch (Exception ex)
            {
            return false;
            }
        }

    /**
     * Check for required client state.
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the namespace name
     * @param sClientPod  the pod name
     *
     * @return {@code true} if required client state is reached.
     */
    public boolean isRequiredClientStateReached(K8sCluster cluster, String sNamespace, String sClientPod)
        {
        try
            {
            Queue<String> sLogs = getPodLog(cluster, sNamespace, sClientPod, null);

            return sLogs.stream().anyMatch(l -> l.contains("Cache Value Before Cloud EntryProcessor: AWS"))
                        && sLogs.stream().anyMatch(l -> l.contains("Cache Value After Cloud EntryProcessor: GCP"));
            }
        catch (Exception ex)
            {
            return false;
            }
        }

    private void installClient(K8sCluster cluster, String name, String sNamespace, String sRelease) throws Exception
        {
        Arguments arguments = Arguments.of("apply");

        if (sNamespace != null)
            {
            arguments = arguments.with("--namespace", sNamespace);
            }

        arguments = arguments.with("-f", getClientYaml(name, sRelease, CLUSTER1));

        int nExitCode = cluster.kubectlAndWait(arguments, SystemApplicationConsole.builder());

        assertThat("kubectl create coherence client pod returned non-zero exit code", nExitCode, is(0));
        }

    /**
     * Determine whether the coherence start up log is ready. param cluster the k8s
     * cluster
     *
     * @param cluster     the k8s cluster
     * @param sNamespace  the namespace of coherence
     * @param sPod        the pod name of coherence
     *
     * @return {@code true} if coherence container is ready
     */
    public boolean hasDefaultCacheServerStarted(K8sCluster cluster, String sNamespace, String sPod)
        {
        try
            {
            Queue<String> sLogs = getPodLog(cluster, sNamespace, sPod, COHERENCE_CONTAINER_NAME);
            return sLogs.stream().anyMatch(l -> l.contains("Started DefaultCacheServer"));
            }
        catch (Exception ex)
            {
            return false;
            }
        }

    private void deleteClients()
        {
        deleteClient(CLIENT1);
        deleteClient(CLIENT2);
        }

    private void deleteClient(String sClient)
        {
        String sNamespace = getK8sNamespace();
        Arguments arguments = Arguments.of("delete", "pod", sClient);

        if (sNamespace != null)
            {
            arguments = arguments.with("-n", sNamespace);
            }

        s_k8sCluster.kubectlAndWait(arguments, SystemApplicationConsole.builder());
        }

    // ----- data members ---------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster     s_k8sCluster = getDefaultCluster();

    /**
     * Time out value for checking the required condition.
     */
    private static final int      TIMEOUT      = 300;

    /**
     * The boolean indicates whether coherence cache data is persisted.
     */
    private static boolean        PERSISTENCE  = false;

    /**
     * The Docker registry name to use to pull Coherence images.
     */
    public static final String DOCKER_REGISTRY = System.getProperty("docker.repo");

    /**
     * The version (tag) for the latest Coherence image version.
     */
    public static final String COHERENCE_VERSION = System.getProperty("coherence.docker.version");

    /**
     * The base Coherence image name without a tag.
     */
    public static final String COHERENCE_IMAGE_PREFIX = DOCKER_REGISTRY + "oracle/coherence:";

    /**
     * The tag to use when pulling the Coherence image for this test.
     */
    private String m_sCoherenceTag;

    /**
     * The name of the deployed Coherence Helm release.
     */
    private String                m_sRelease;

    private static final String   CLIENT1  = "coh-client-1";

    private static final String   CLIENT2  = "coh-client-2";

    private static final String   CLUSTER1 = "MyCoherenceCluster";
    }
