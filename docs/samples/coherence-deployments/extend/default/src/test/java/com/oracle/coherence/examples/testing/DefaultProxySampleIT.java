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

import org.junit.Before;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Collection;
/**
 * Test the default-proxy-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the coherence.yaml file.
 *
 * @author tam  2019.05.14
 */
@RunWith(Parameterized.class)
public class DefaultProxySampleIT
        extends BaseProxySampleTest
    {

    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public DefaultProxySampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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
            String[] asCohNamespaces = getTargetNamespaces();

            m_asReleases = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), asCohNamespaces,"coherence.yaml");

            assertCoherence(s_k8sCluster, asCohNamespaces, m_asReleases);
            }
        }

    // ----- tests ----------------------------------------------------------

    /**
     * Test the default proxy sample.
     * 
     * @throws Exception
     */
    @Test
    public void testDefaultProxySample() throws Exception
        {
        if (testShouldRun())
            {
            // test proxy connection to coherence pod
            testProxyConnection(m_asReleases[0]);
            }
        }
    }
