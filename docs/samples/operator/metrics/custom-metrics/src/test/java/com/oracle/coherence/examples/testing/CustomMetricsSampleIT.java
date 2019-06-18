/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.options.Arguments;

import com.oracle.bedrock.testsupport.deferred.Eventually;
import org.hamcrest.MatcherAssert;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static com.oracle.coherence.examples.testing.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.Matchers.notNullValue;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.assertTrue;


/**
 * Test the custom-metrics-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.22
 */
@RunWith(Parameterized.class)
public class CustomMetricsSampleIT
        extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public CustomMetricsSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
        {
        super(sOperatorChartURL, sCoherenceChartURL);
        }

    // ----- test lifecycle -------------------------------------------------

    @Parameterized.Parameters
    public static Collection testParameters()
        {
        return buildTestParameters();
        }

    /**
     * Install the charts required for the test.
     *
     * @throws Exception
     */
    @Before
    public void installCharts() throws Exception
        {
        if (testShouldRun())
            {
            // install Coherence Operator chart
            s_sOperatorRelease = installOperator("coherence-operator.yaml",toURL(m_sOperatorChartURL));

            String sCohNamespace = getTargetNamespaces()[0];

            assertDeploymentReady(s_k8sCluster, sCohNamespace, getCoherenceOperatorSelector(s_sOperatorRelease), true);

            // process the file servicemonitoring.yaml and replace sample-coherence-n with our namespace from above
            String sProcessedYaml = getProcessedYamlInstallFile(
                    "src/main/yaml/servicemonitoring.yaml", sCohNamespace, s_sOperatorRelease);
            assertThat(sProcessedYaml, is(notNullValue()));

            // install the service monitor
            Arguments args = Arguments.of("create", "-f", sProcessedYaml);
            int nExitCode = s_k8sCluster.kubectlAndWait(args);
            assertThat("Unable to install servicemonitor", nExitCode, is(0));
            
            String sTag = System.getProperty("docker.push.tag.prefix") + System.getProperty("project.artifactId") + ":" +
                          System.getProperty("project.version");

            // process yaml file to replace user artifacts image
            String sProcessedCoherenceYaml = getProcessedYamlFile("coherence.yaml", sTag, null);
            assertThat(sProcessedCoherenceYaml, is(notNullValue()));

            // install Coherence chart
            String sClusterRelease = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), sCohNamespace, sProcessedCoherenceYaml);
            assertCoherence(s_k8sCluster, sCohNamespace, sClusterRelease);
            
            m_asReleases = new String[] { sClusterRelease };
            }
        }


    // ----- tests ----------------------------------------------------------

    /**
     * Test the custom metrics sample.
     *
     * @throws Exception
     */
    @Test
    public void testCustomMetricsSample() throws Exception
        {
        if (testShouldRun())
            {
            try {
                // connect to cluster - m_asReleases[0]
                testProxyConnection(m_asReleases[0]);

                // validate that the metrics turn up in prometheus
                // retrieve the client tier release pod
                String       sNamespace         = getTargetNamespaces()[0];
                String       sCoherenceSelector = getCoherencePodSelector(m_asReleases[0]);
                List<String> listPods           = getPods(s_k8sCluster, getTargetNamespaces()[0], sCoherenceSelector);

                MatcherAssert.assertThat(listPods.size(), is(1));
                String sCoherencePod = listPods.get(0);

                // wait for the following message to appear, which indicates the metrics are enabled
                // Inserted key=40, value=08:33:43
                Eventually.assertThat(invoking(this).hasLogMessageAppeared(s_k8sCluster, sNamespace, sCoherencePod, "Creating a custom metrics endpoint"),
                        is(true), Timeout.after(300, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

                // wait until we get metrics ending up in prometheus
                Eventually.assertThat(invoking(this).getPrometheusMetricAsLong("PromInterceptor_custom_scrape_total_PromInterceptor_updates", CLUSTER_NAME, sNamespace, sNamespace),
                        greaterThanOrEqualTo(5L),
                        RetryFrequency.every(10, TimeUnit.SECONDS),
                        Timeout.after(4, TimeUnit.MINUTES));
            }
            finally
                {
                // remove the service monitor
                Arguments args = Arguments.of("delete",  "servicemonitor",  "custom-coherence-monitor", "--namespace", getTargetNamespaces()[0]);
                int nExitCode = s_k8sCluster.kubectlAndWait(args);
                assertThat("Unable to remove servicemonitor", nExitCode, is(0));
                }
            }
        }

    // ----- helpers --------------------------------------------------------

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

            Queue<String> result = processHttpRequest(s_k8sCluster,"GET", sHost, 9090, sRestQueryPath, true);
            String sJson          = result != null ? result.element() : null;

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

        return lMetricValue;
        }


    // ----- constants ------------------------------------------------------

    private static final String CLUSTER_NAME = "metrics-cluster";
    }

