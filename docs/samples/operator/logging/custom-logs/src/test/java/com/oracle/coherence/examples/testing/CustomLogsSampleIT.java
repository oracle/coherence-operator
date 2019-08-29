/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.deferred.options.MaximumRetryDelay;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.net.ConfigurableCacheFactory;
import com.tangosol.net.NamedCache;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.io.IOException;
import java.util.Collection;
import java.util.List;
import java.util.Map;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static com.oracle.coherence.examples.testing.HelmUtils.HELM_TIMEOUT;
import static com.oracle.coherence.examples.testing.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the customer-logger-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.22
 */
@RunWith(Parameterized.class)
public class CustomLogsSampleIT
    extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public CustomLogsSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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
            installChartsSingleTier(m_sOperatorChartURL, m_sCoherenceChartURL);

            String sNamespace = getTargetNamespaces()[0];
            
            assertDeploymentReady(s_k8sCluster, sNamespace, getKibanaSelector(s_sOperatorRelease), true);

            // get the Elastic search pod
            String sElasticSearchSelector = getElasticSearchSelector(s_sOperatorRelease);

            assertDeploymentReady(s_k8sCluster, sNamespace, sElasticSearchSelector, true);

            List<String> listPods = getPods(s_k8sCluster, sNamespace, sElasticSearchSelector);
            assertThat(listPods.size(), is(1));

            s_sElasticsearchPod = listPods.get(0);

            // get the Kibana Pod
            String       sKibanaSelector = getKibanaSelector(s_sOperatorRelease);
            List<String> listPodsKibana  = getPods(s_k8sCluster, sNamespace, sKibanaSelector);
            assertThat(listPodsKibana.size(), is(1));

            s_sKibanaPod = listPodsKibana.get(0);

            Eventually.assertThat(invoking(this).isFluentdReady(),
                              is(true),
                              Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS),
                              InitialDelay.of(10, TimeUnit.SECONDS),
                              MaximumRetryDelay.of(RETRY_FREQUENCEY_SECONDS, TimeUnit.SECONDS),
                              RetryFrequency.every(RETRY_FREQUENCEY_SECONDS, TimeUnit.SECONDS));
            }
        }

    // ----- tests ----------------------------------------------------------

    /**
     * Test the default proxy sample.
     *
     * @throws Exception
     */
    @Test
    public void testCustomLogsSample() throws Exception
        {
        if (testShouldRun())
            {
            String sCoherenceRelease = m_asReleases[0];
            // test proxy connection to coherence pod
            try (Application application = portForwardExtend(sCoherenceRelease, 20000))
                {
                PortMapping              portMapping = application.get(PortMapping.class);
                int                      nActualPort = portMapping.getPort().getActualPort();
                ConfigurableCacheFactory ccf         = getCacheFactory("client-cache-config.xml", nActualPort);

                Eventually.assertThat(invoking(ccf).ensureCache("test", null), is(notNullValue()));

                NamedCache nc = ccf.ensureCache("test", null);
                nc.put("key-1", "value-1");

                // server side interceptor should make uppercase and log this to custom "cloud" logger
                assertThat(nc.get("key-1"), is("VALUE-1"));
                ccf.dispose();
                }

            // wait until the log makes its way to elastic search
            Eventually.assertThat(invoking(this).isCustomLogInElasticSearch(),
                      is(true),
                      Timeout.after(600, TimeUnit.SECONDS),
                      InitialDelay.of(10, TimeUnit.SECONDS),
                      MaximumRetryDelay.of(RETRY_FREQUENCEY_SECONDS, TimeUnit.SECONDS),
                      RetryFrequency.every(RETRY_FREQUENCEY_SECONDS, TimeUnit.SECONDS));
            }
        }

    // ----- helpers --------------------------------------------------------

    /**
     * Determine whether the Helm fluentd is ready (i.e all Pods are ready).
     *
     * @return  {@code true} if the fluentd is ready
     */
    // must be public - used in Eventually.assertThat call.
    public boolean isFluentdReady()
        {
        Queue<String> queueLines = processElasticsearchQuery("/_cat/indices");
        return queueLines.stream().anyMatch(s -> s.contains("coherence"));
        }

    /**
     * Determine whethere the log message has got to fluentd.
     *
     * @return {@code true} if the log message is in fluentd
     */
     // must be public - used in Eventually.assertThat call.
    public boolean isCustomLogInElasticSearch()
        {
        String sIndex = getESIndex();
        if (sIndex != null)
            {
            // we have the index
            Queue<String> queueLogs = new ConcurrentLinkedQueue<>();

            queueLogs.addAll(processElasticsearchQuery("/" + sIndex + "/_search?q=log:Before"));
            Map<String, ?> map = null;
            try
                {
                map = HelmUtils.JSON_MAPPER.readValue(queueLogs.stream().collect(Collectors.joining()), Map.class);
                }
            catch (IOException e)
                {
                return false;
                }

            Map<String, List<Map<String, ?>>> mapHits = (Map<String, List<Map<String, ?>>>) map.get("hits");

            // "hist" element should always wbe there
            if (mapHits != null)
                {
                List<Map<String, ?>> list = mapHits.get("hits");

                // log is present if the "hits" element is present and has at least 1 entry
                return list != null && list.size() > 0;
                }
            }
        System.err.println("Index not yet present");
        return false;
        }

    /**
     * Return the index name for the cloud- prefix.
     *
     * @return the index name for the cloud- prefix
     */
    protected String getESIndex()
        {
        Queue<String>  queueIndices = processElasticsearchQuery("/_cat/indices");
        String         sIndexName   = queueIndices.stream().filter(s -> s.contains("cloud-"))
                .map(s -> s.split(" ")[2]).findFirst().orElse(null);
        return sIndexName;
        }

    /**
     * Process a query from elastic search.
     *
     * @param sPath path to search for.
     *
     * @return the results
     */
    protected Queue<String> processElasticsearchQuery(String sPath)
        {
        return processHttpRequest(s_k8sCluster, s_sElasticsearchPod, "GET", "localhost", 9200, sPath);
        }

    // ----- constants ------------------------------------------------------

    /**
     * Kibana pod.
     */
    private static String s_sKibanaPod;

    /**
     * Elastic search pod.
     */
    private static String s_sElasticsearchPod;

    /**
     * The retry frequency in seconds.
     */
    private static final int RETRY_FREQUENCEY_SECONDS = 10;

    }

