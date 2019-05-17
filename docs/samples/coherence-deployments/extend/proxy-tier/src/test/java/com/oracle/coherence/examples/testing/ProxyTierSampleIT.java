/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.net.ConfigurableCacheFactory;
import com.tangosol.net.NamedCache;
import com.tangosol.util.Resources;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

import java.util.Collection;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the proxy-tier-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the coherence.yaml file.
 *
 * @author tam  2019.05.14
 */
@RunWith(Parameterized.class)
public class ProxyTierSampleIT
        extends BaseProxySampleTest
    {

    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public ProxyTierSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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

            // process yaml file for proxy tier
            String sProcessedProxyYaml = getProcessedYamlFile("coherence-proxy-tier.yaml", sTag, sClusterRelease);
            assertThat(sProcessedProxyYaml, is(notNullValue()));

            String sProxyRelease = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), sCohNamespace, sProcessedProxyYaml);
            assertCoherence(s_k8sCluster,sCohNamespace, sProxyRelease);

            m_asReleases = new String[] { sClusterRelease, sProxyRelease };
            }
        }

    // ----- tests ----------------------------------------------------------

    /**
     * Test the proxy tier sample.
     *
     * @throws Exception
     */
    @Test
    public void testProxyTierSample() throws Exception
        {
        if (testShouldRun())
            {
            // connect to proxy tier - m_asReleases[1]
            testProxyConnection(m_asReleases[1]);
            }
        }
    }
