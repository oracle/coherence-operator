/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.console.SystemApplicationConsole;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.k8s.helm.Helm;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.testsupport.deferred.Eventually;

import org.junit.*;

import util.AssumingCoherenceVersion;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Map;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static helm.HelmUtils.HELM_TIMEOUT;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.core.Is.is;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.assertTrue;

/**
 * Test that Coherence K8s operator deploys Prometheus Operator as a subchart.
 *
 * @author jf
 * @author sc
 */
@Ignore
public class PrometheusOperatorHelmSubChartIT
    extends BaseHelmChartTest
    {
    // ----- test lifecycle --------------------------------------------------

    @BeforeClass
    public static void setup()
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

    @After
    public void cleanUpInstalledReleases()
        throws InterruptedException
        {
        if (m_asReleases != null)
            {
            deleteCoherence(s_k8sCluster, getTargetNamespaces(), m_asReleases, PERSISTENCE);
            m_asReleases = null;
            }

        if (m_sRelease != null)
            {
            cleanupHelmReleases(m_sRelease);
            m_sRelease = null;
            }
        }

    // ----- test methods ---------------------------------------------------

    /**
     * Test with default prometheus scraping using coherence-service-monitor.
     */
    @Test
    public void testDefaultPrometheusOperatorSubchart()
        throws Exception
        {
        testPrometheusOperatorSubchart("values/helm-values-prometheus-operator-subchart.yaml");
        }

    /**
     * Test with user supplied prometheus_io annotations with coherence-service-monitor disabled.
     * @throws Exception
     */
    @Test
    public void testPrometheusOperatorSubchartWithUserSuppliedAdditionalScrapeConfig()
        throws Exception
        {
        testPrometheusOperatorSubchart("values/helm-values-prometheus-operator-subchart-use-annotations.yaml");
        }

    // ----- helpers --------------------------------------------------------

    /**
     * Test with specified values file.
     *
     * @param sOperatorValuesFile  prometheus enabled values file.
     */
    private void testPrometheusOperatorSubchart(String sOperatorValuesFile)
        throws Exception
        {
        String sNamespace               = getK8sNamespace();
        String sSetValueSkipInstallCrd  = hasPrometheusOperatorCRD(s_k8sCluster) ?
            "prometheusoperator.prometheusOperator.createCustomResource=false" : null;

        System.err.println("Deploying " + OPERATOR_HELM_CHART_NAME + " with " + sOperatorValuesFile + " in namespace " + sNamespace);

        try
            {
            m_sRelease = installChart(s_k8sCluster,
                OPERATOR_HELM_CHART_NAME,
                OPERATOR_HELM_CHART_URL,
                sNamespace,
                sOperatorValuesFile,
                withDefaultHelmSetValues(sSetValueSkipInstallCrd));
            }
        catch (Throwable t)
            {
            System.err.println("failed installing " + OPERATOR_HELM_CHART_NAME);
            t.printStackTrace();
            throw t;
            }

        assertDeploymentReady(s_k8sCluster, sNamespace, getCoherenceOperatorSelector(m_sRelease), true);

        // override values extracted from this test's value file, resources/values/helm-values-coh-metrics.yaml
        final String CLUSTER_NAME          = "myClusterWithMetricsEnabled";
        final Long   EXPECTED_CLUSTER_SIZE = 2L;
        String[]     asCohNamespaces       = getTargetNamespaces();
        String       sOpNamespace          = getK8sNamespace();

        try
            {
            Eventually.assertThat("assert prometheus server is running in namespace " + sOpNamespace,
                invoking(this).getPrometheusPodsCount(sOpNamespace),
                greaterThanOrEqualTo(1),
                RetryFrequency.every(10, TimeUnit.SECONDS),
                Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS));
            }
        finally
            {
            dumpInfo("after assertions");
            }

        m_asReleases = installCoherence(s_k8sCluster, asCohNamespaces, "values/helm-values-coh-metrics.yaml");
        assertCoherence(s_k8sCluster, asCohNamespaces, m_asReleases);

        for (int i = 0; i < m_asReleases.length; i++)
            {
            System.err.println("validating cluster size of " + CLUSTER_NAME + " in namespace " + asCohNamespaces[i] + " for coherence release " + m_asReleases[i]);
            Eventually.assertThat(invoking(this).getPrometheusMetricAsLong("coherence_cluster_size", CLUSTER_NAME, asCohNamespaces[i], sOpNamespace),
                greaterThanOrEqualTo(EXPECTED_CLUSTER_SIZE),
                RetryFrequency.every(10, TimeUnit.SECONDS),
                Timeout.after(4, TimeUnit.MINUTES));
            }
        }

    private void dumpInfo(String sHeading)
        {
        System.err.println("--------------------------------------------------");
        System.err.println("All Pods " + sHeading);
        System.err.println("--------------------------------------------------");
        s_k8sCluster.kubectlAndWait(Arguments.of("get", "pods", "--all-namespaces=true", "-o", "custom-columns='NAME:.metadata.name,NAMESPACE:.metadata.namespace,LABELS:.metadata.labels'"),
                                    SystemApplicationConsole.builder());
        System.err.println("--------------------------------------------------");
        System.err.println("Helm status " + sHeading);
        System.err.println("--------------------------------------------------");
        try
            {
            Helm.status(m_sRelease).executeAndWait(SystemApplicationConsole.builder());
            }
        catch (Throwable t)
            {
            System.err.println("Cannot execute helm status " + t);
            }
        System.err.println("--------------------------------------------------");
        }

    /**
     * Return list of prometheus pods in namespace.
     * @param sNamespace namespace
     *
     * @return  the number of Prometheus Pods
     */
    public Integer getPrometheusPodsCount(String sNamespace)
        {
        return HelmUtils.getPods(s_k8sCluster, sNamespace, "app=prometheus").size();
        }

    /**
     * Query prometheus server rest api for MetricName value.
     *
     * @param sMetricName                   query for this prometheus metric name
     * @param sClusterName                  cluster name to match metric tag cluster
     * @param sNamespace                    namespace to match metric tag kubernetes_namespace
     * @param sPrometheusOperatorNamespace  namespace to look for prometheus service to scrape.
     *
     * @return prometheus metric name's value or -1L.
     */
    public Long getPrometheusMetricAsLong(String sMetricName,
                                          String sClusterName,
                                          String sNamespace,
                                          String sPrometheusOperatorNamespace)
        {
        Long lMetricValue = -1L;

        try
            {
            String        sHost          = "prometheus-operated." + sPrometheusOperatorNamespace + ".svc.cluster.local";
            String        sRestQueryPath = "/api/v1/query?query=" + sMetricName;

            System.err.println("making http request to " + sHost + ":" + 9090 + sRestQueryPath);

            Queue<String> result         = processHttpRequest(s_k8sCluster,"GET", sHost, 9090, sRestQueryPath, true);
            System.err.println("completed http request to " + sHost + ":" + 9090 + sRestQueryPath);

            String        sJson          = result != null ? result.element() : null;

            if (sNamespace == null)
                {
                sNamespace = "default";
                }

            assertTrue("prometheus rest query for metric " + sMetricName + " did not succeed",
                sJson != null && sJson.contains("status") && sJson.contains("success"));

            Map<String, ?> mapResponse = HelmUtils.JSON_MAPPER.readValue(sJson, Map.class);
            Map            mapData     = (Map) mapResponse.get("data");
            ArrayList<Map> arrResult   = (ArrayList<Map>) mapData.get("result");

            assertThat("prometheus rest query did not return any results", arrResult.size(), greaterThanOrEqualTo(1));

            Set<String> setTargetNamespaces = new HashSet<String>(Arrays.asList(getTargetNamespaces()));
            for (Map mapResult : arrResult)
                {
                Map<String, String> mapMetricTags = (Map) mapResult.get("metric");

                assertThat("assertion: only scrape coherence servers in target namespaces",
                    setTargetNamespaces.contains(mapMetricTags.get("namespace")), is(true));
                }

            for (Map mapResult : arrResult)
                {
                Map mapMetricTags = (Map) mapResult.get("metric");

                if (sNamespace.equals(mapMetricTags.get("namespace")) &&
                    sClusterName.equals(mapMetricTags.get("cluster")))
                    {
                    ArrayList arrMetricValue = (ArrayList) mapResult.get("value");

                    lMetricValue = (arrMetricValue != null && arrMetricValue.size() == 2) ?
                            Long.valueOf((String) arrMetricValue.get(1)) : 0L;

                    System.err.println("metric name: " + mapMetricTags.get("__name__") +
                        " for coherence cluster " + mapMetricTags.get("cluster") +
                        " in namespace " + mapMetricTags.get("namespace") + " has a value of " + lMetricValue);
                    return lMetricValue;
                    }
                }
            }
        catch (Throwable t)
            {
            System.out.println("getPrometheusMetricAsLong: unexpected exception " + t.getClass().getCanonicalName() + ": "+ t.getMessage());
            }

        // not found
        System.out.println("getPrometheusMetricAsLong: not found for metric " + sMetricName + " for cluster " + sClusterName + " in namespace " + sNamespace);
        return lMetricValue;
        }

    // ----- data members ---------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The boolean indicates whether coherence cache data is persisted.
     */
    private static boolean PERSISTENCE = false;

    /**
     * The name of the deployed Helm release.
     */
    private String m_sRelease;

    /**
     * The name of the deployed Coherence Helm releases.
     */
    private String[] m_asReleases;

    /**
     * Prometheus metrics are only available from Coherence 12.2.1.4.0 and greater.
     */
    @ClassRule
    public static AssumingCoherenceVersion assumingCoherenceVersion = new AssumingCoherenceVersion(COHERENCE_VERSION, "12.2.1.4.0");
    }
