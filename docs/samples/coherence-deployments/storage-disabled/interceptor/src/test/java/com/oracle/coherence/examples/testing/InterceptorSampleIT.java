/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.testsupport.deferred.Eventually;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Collection;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static com.oracle.coherence.examples.testing.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the interceptor-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.14
 */
@RunWith(Parameterized.class)
public class InterceptorSampleIT
      extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public InterceptorSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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

            // install Coherence chart
            String sCohNamespace = getTargetNamespaces()[0];

            String sTag = System.getProperty("docker.push.tag.prefix") + System.getProperty("project.artifactId") + ":" +
                          System.getProperty("project.version");

            // process yaml file to replace user artifacts image
            String sProcessedCoherenceYaml = getProcessedYamlFile("coherence.yaml", sTag, null);
            assertThat(sProcessedCoherenceYaml, is(notNullValue()));

            String sClusterRelease = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), sCohNamespace, sProcessedCoherenceYaml);
            assertCoherence(s_k8sCluster, sCohNamespace, sClusterRelease);

            // process yaml file for storage-disabled client tier
            String sProcessedProxyYaml = getProcessedYamlFile("coherence-client-tier.yaml", sTag, sClusterRelease);
            assertThat(sProcessedProxyYaml, is(notNullValue()));

            String sClientTierRelease = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), sCohNamespace, sProcessedProxyYaml);
            assertCoherence(s_k8sCluster,sCohNamespace, sClientTierRelease);

            m_asReleases = new String[] { sClusterRelease, sClientTierRelease };
            }
        }

    // ----- tests ----------------------------------------------------------

    /**
     * Test the proxy tier sample.
     *
     * @throws Exception
     */
    @Test
    public void testInterceptorSample() throws Exception
        {
        if (testShouldRun())
            {
            // retrieve the client tier release pod
            String sNamespace         = getTargetNamespaces()[0];
            String sCoherenceSelector = getCoherencePodSelector(m_asReleases[1]);
            List<String> listPods     = getPods(s_k8sCluster, getTargetNamespaces()[0], sCoherenceSelector);

            assertThat(listPods.size(), is(1));
            String sCoherencePod = listPods.get(0);

            // wait for the following message to appear, which indicates the interceptor is running
            // Inserted key=40, value=08:33:43
            Eventually.assertThat(invoking(this).hasLogMessageAppeared(s_k8sCluster, sNamespace, sCoherencePod, "Inserted key"),
                    is(true), Timeout.after(120, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

            }
        }
    }
