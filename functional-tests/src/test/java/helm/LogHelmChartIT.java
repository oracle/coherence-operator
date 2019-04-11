/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import org.junit.*;

import java.util.List;
import java.util.Queue;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static helm.HelmUtils.HELM_TIMEOUT;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test log aspects of the Helm chart values file.
 * <p>
 * This test depends on certain artifacts produced by the Maven build
 * so although this test can be run from an IDE it requires that at least
 * a Maven build with at least the package phase being run first.
 *
 * @author sc
 */
@Ignore
public class LogHelmChartIT
        extends BaseHelmChartTest {

    // ----- test lifecycle --------------------------------------------------

    @BeforeClass
    public static void setup() throws Exception
        {
        assertPreconditions(s_k8sCluster);
        ensureNamespace(s_k8sCluster);
        ensureSecret(s_k8sCluster);

        String sOpNamespace = getK8sNamespace();
        for (String sNamespace : getTargetNamespaces())
            {
            if (sNamespace != null && !sNamespace.equals(sOpNamespace))
                {
                ensureNamespace(s_k8sCluster, sNamespace);
                }
            }

        String sOperatorValuesFile = "values/helm-values-log.yaml";
        System.err.println("Deploying " + OPERATOR_HELM_CHART_NAME + " with " + sOperatorValuesFile);

        String sNamespace = getK8sNamespace();

        s_sOperatorRelease = installChart(s_k8sCluster,
                                          OPERATOR_HELM_CHART_NAME,
                                          OPERATOR_HELM_CHART_URL,
                                          sNamespace,
                                          sOperatorValuesFile,
                                          getDefaultHelmSetValues());

        assertDeploymentReady(s_k8sCluster, sNamespace, getCoherenceOperatorSelector(s_sOperatorRelease), true);
        }

    @AfterClass
    public static void cleanup()
        {
        if (s_sOperatorRelease != null)
            {
            try
                {
                capturePodLogs(LogHelmChartIT.class, s_k8sCluster, getCoherenceOperatorSelector(s_sOperatorRelease), null);
                cleanupHelmReleases(s_sOperatorRelease);
                }
            catch (Throwable t)
                {
                System.err.println("Error in cleaning up Helm Releases: " + s_sOperatorRelease);
                }
            }

        cleanupPullSecrets(s_k8sCluster);
        cleanupNamespace(s_k8sCluster);

        for (String sNamespace : getTargetNamespaces())
            {
            cleanupNamespace(s_k8sCluster, sNamespace);
            }
        }

    @After
    public void cleanUpCoherence()
        {
        if (m_asReleases != null)
            {
            deleteCoherence(s_k8sCluster, getTargetNamespaces(), m_asReleases, PERSISTENCE);
            }
        }

    // ----- test methods ---------------------------------------------------
    /**
     * Test default coherence values.yaml.
     * Test the log messages for create and delete coherence.
     *
     * @throws Exception
     */
    @Test
    public void testDefaultCoherenceYaml() throws Exception
        {
        String[] asCohNamespaces = getTargetNamespaces();
                 m_asReleases    = installCoherence(s_k8sCluster, asCohNamespaces, null);

        assertCoherence(s_k8sCluster, asCohNamespaces, m_asReleases);

        // verify Coherence started
        for (int i = 0; i < m_asReleases.length; i++)
            {
            String sCoherenceSelector = getCoherencePodSelector(m_asReleases[i]);
            List<String> listPods = HelmUtils.getPods(s_k8sCluster, asCohNamespaces[i], sCoherenceSelector);

            Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, asCohNamespaces[i], listPods.get(0)), is(true),
                    Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));
            }
        }


    /**
     * Test the log messages for jvm options and start.
     *
     * @throws Exception
     */
    @Test
    public void testCoherenceJavaOpts() throws Exception
        {
        String[] asCohNamespaces = getTargetNamespaces();
                 m_asReleases    = installCoherence(s_k8sCluster, asCohNamespaces,"values/helm-values-coh-jvm.yaml");

        assertCoherence(s_k8sCluster, asCohNamespaces, m_asReleases);

        // verify the jvm option, role and cluster are set for each Coherence
        for (int i = 0; i < m_asReleases.length; i++)
            {
            String sCoherenceSelector = getCoherencePodSelector(m_asReleases[i]);
            List<String> listPods = HelmUtils.getPods(s_k8sCluster, asCohNamespaces[i], sCoherenceSelector);

            Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, asCohNamespaces[i], listPods.get(0)), is(true),
                                  Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

            dumpPodLog(s_k8sCluster, asCohNamespaces[i], listPods.get(0));
            Queue<String> sLogs = getPodLog(s_k8sCluster, asCohNamespaces[i], listPods.get(0));

            assertThat(sLogs.stream().anyMatch(l -> l.contains("-Dtest1=dummy")), is(true));
            assertThat(sLogs.stream().anyMatch(l -> l.contains("-Dcoherence.role=myrole")), is(true));
            assertThat(sLogs.stream().anyMatch(l -> l.contains("-Dcoherence.cluster=mycluster")), is(true));
            assertThat(sLogs.stream().anyMatch(l -> l.contains("Role=myrole")), is(true));
            assertThat(sLogs.stream().anyMatch(l -> l.contains("Started cluster Name=mycluster")), is(true));
            }
        }

    // ----- data members ---------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The name of the deployed Operator Helm release.
     */
    private static String s_sOperatorRelease;

    /**
     * The boolean indicates whether Coherence cache data is persisted.
     */
    private static boolean PERSISTENCE = false;

    /**
     * The name of the deployed Coherence Helm releases.
     */
    private String[] m_asReleases;
    }
