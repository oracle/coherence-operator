/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.runtime.options.Console;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import helm.BaseHelmChartTest;

import org.junit.After;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.ClassRule;
import org.junit.Test;

import util.AssumingCoherenceVersion;

import java.util.Arrays;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static com.oracle.bedrock.testsupport.deferred.Eventually.within;
import static helm.HelmUtils.HELM_TIMEOUT;
import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.Assert.assertNotNull;

/**
 * Test that validates coherence container prometheus metrics endpoint is publishing
 * coherence metrics. Test validates coherence_cluster_size and coherence_cache_size.
 * <p>
 * This test is an adaption of CustomJarInClasspathIT since it was
 * only existing functional-test that created and added to a coherence cache.
 * Since this test does not install coherence-operator,
 * there is no default prometheus server in target namespace.
 * This test is ideal for running with an independently installed prometheus-operator
 * that can remain running across multiple test runs.
 *
 * @author jf
 * @author cp
 */
public class MetricsUsingCustomJarInClasspathIT
    extends BaseHelmChartTest
    {
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

    @After
    public void cleanUpCoherence()
        {
        if (m_sRelease != null)
            {
            deleteCoherence(s_k8sCluster, getK8sNamespace(), m_sRelease, PERSISTENCE);
            }
        }

    @Test
    public void testMetrics() throws Exception
        {
        testCoherenceWithUserSuppliedJarInClasspath("values/helm-values-coh-metrics-user-artifacts.yaml");
        }

    /**
     * Run the test using supplied list of values files.
     *
     * @param values  value file names
     *
     * @throws Exception
     */
    public void testCoherenceWithUserSuppliedJarInClasspath(String... values) throws Exception
        {
        final String sNamespace = getK8sNamespace();

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, values[0]);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);
        assertCoherenceService(s_k8sCluster, sNamespace, m_sRelease);

        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPods           = getPods(s_k8sCluster, sNamespace, sCoherenceSelector);

        assertThat(listPods.size(), is(2));

        // wait for ALL pods to be up.
        for (String sPod : listPods)
            {
            System.err.println("Waiting for Coherence Pod " + sPod + "...");
            Eventually.assertThat(invoking(this).isCoherenceMetricsReady(s_k8sCluster, sNamespace, sPod),
                                  is(true), Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS),
                                  InitialDelay.of(10, TimeUnit.SECONDS));
            }

        System.err.println("Coherence Pods started");

        try
            {
            installClient(s_k8sCluster, CLIENT1, sNamespace, m_sRelease, CLUSTER1);
            installClient(s_k8sCluster, CLIENT2, sNamespace, m_sRelease, CLUSTER1);

            System.err.println("Waiting for Client-1 initial state...");
            Eventually.assertThat(invoking(this).isRequiredClientStateReached(s_k8sCluster, sNamespace, CLIENT1),
                                  is(true),
                                  within(TIMEOUT, TimeUnit.SECONDS),
                                  RetryFrequency.fibonacci());

            // validate metrics ports
            for (String sPod : listPods)
                {
                Queue<String> metricsScrape = getMetricsScrape(sPod);
                assertNotNull("missing expected metrics scrape from pod " + sPod, metricsScrape);
                System.err.println("Metrics scrape from pod " + sPod + metricsScrape.poll() + "\n" + metricsScrape.poll() + "...");
                }

            Eventually.assertThat(invoking(this).isRequiredClientStateReached(s_k8sCluster, sNamespace, CLIENT1),
                                  is(true),
                                  within(120, TimeUnit.SECONDS),
                                  RetryFrequency.fibonacci());

            Eventually.assertThat(invoking(this).isRequiredClientStateReached(s_k8sCluster, sNamespace, CLIENT2),
                                  is(true),
                                  within(120, TimeUnit.SECONDS),
                                  RetryFrequency.fibonacci());

            for (int i = 0; i < 6; i++)
                {
                Thread.sleep(5000L);
                for (String sPod : listPods)
                    {
                    Queue<String> metricsScrape = getMetricsScrape(sPod);
                    assertNotNull("Missing expected metrics scrape from pod " + sPod, metricsScrape);
                    System.err.println("Metrics scrape from pod " + sPod);
                    System.err.println(getMetric(metricsScrape, "coherence_cluster_size"));
                    System.err.println(getMetric(metricsScrape, "coherence_cache_size"));
                    System.err.println(getMetric(metricsScrape, "coherence_cache_hits"));
                    System.err.println(getMetric(metricsScrape, "coherence_cache_misses"));
                    System.err.println(getMetric(metricsScrape, "coherence_cache_misses_millis"));
                    System.err.println(getMetric(metricsScrape, "coherence_cache_total_puts_millis"));
                    System.err.println(getMetric(metricsScrape, "coherence_cache_total_gets_millis"));
                    System.err.println();
                    }
                }
            }
        finally
            {
            deleteClients();
            }
        }

    /**
     * Metric scrape of coherence container.
     *
     * @param sPod  pod name running coherence container
     *
     * @return a {@link Queue} of metrics, one per line.
     */
    Queue<String> getMetricsScrape(String sPod)
        {
        return processHttpRequest(s_k8sCluster, sPod, "GET", "localhost", 9095, "/metrics");
        }

    /**
     * Return line containing provided metric name from provided metrics.
     * Very rudimentary method that does not differentiate between any metric tags.
     *
     * @param metrics      metric response
     * @param sMetricName  metric name to search for
     *
     * @return first metric line in metrics that matches sMetricName
     */
    String getMetric(Queue<String> metrics, String sMetricName)
        {
        for (String sMetricLine : metrics)
            {
            if (sMetricLine.contains(sMetricName))
                {
                if (sMetricLine.contains("cache"))
                    {
                    // filter out tier=front for now.  no near cache.
                    if (sMetricLine.contains("tier=\"back\""))
                        {
                        return sMetricLine;
                        }
                    }
                else
                    {
                    return sMetricLine;
                    }
                }
            }
        return "metricName " + sMetricName + " not found";
        }

    /**
     * Determine whether the coherence service metrics service is is ready. param cluster the k8s
     * cluster
     *
     * @param sNamespace  the namespace of coherence
     * @param sPod        the pod name of coherence
     *
     * @return {@code true} if coherence container is ready
     */
    public boolean isCoherenceMetricsReady(K8sCluster cluster, String sNamespace, String sPod)
        {
        try
            {
            Queue<String> sLogs = getPodLog(cluster, sNamespace, sPod);
            return sLogs.stream().anyMatch(l -> l.contains("Service MetricsHttpProxy joined"));
            }
        catch (Exception ex)
            {
            return false;
            }
        }

    private Queue<String> deleteClients()
        {
        CapturingApplicationConsole console = new CapturingApplicationConsole();

        Arguments arguments = Arguments.empty();
        arguments = arguments.with("delete", "pod", CLIENT1, "-n", getK8sNamespace());
        int nExitCode = s_k8sCluster.kubectlAndWait(arguments, Console.of(console));
        if (nExitCode != 0)
            {
            throw new IllegalStateException("kubectl delete coherence client pod returned non-zero exit code.");
            }

        arguments = Arguments.empty();
        arguments = arguments.with("delete", "pod", CLIENT2, "-n", getK8sNamespace());
        nExitCode = s_k8sCluster.kubectlAndWait(arguments, Console.of(console));
        if (nExitCode != 0)
            {
            throw new IllegalStateException("kubectl delete coherence client pod returned non-zero exit code.");
            }
        return console.getCapturedOutputLines();
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
     * The name of the deployed Coherence Helm release.
     */
    private String                m_sRelease;

    private static final String   CLIENT1  = "coh-client-1";
    private static final String   CLIENT2  = "coh-client-2";
    private static final String[] CLIENTS  = new String[] { CLIENT1, CLIENT2 };
    private static final String   CLUSTER1 = "MyCoherenceCluster";

    /**
     * Prometheus metrics are only available from Coherence 12.2.1.4.0 and greater.
     */
    @ClassRule
    public static AssumingCoherenceVersion assumingCoherenceVersion
            = new AssumingCoherenceVersion(COHERENCE_IMAGE, "12.2.1.4.0");
    }
