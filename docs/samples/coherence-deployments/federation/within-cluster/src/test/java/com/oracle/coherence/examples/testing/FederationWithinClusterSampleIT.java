/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Collection;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the federation-within-cluster-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.14
 */
@RunWith(Parameterized.class)
public class FederationWithinClusterSampleIT
         extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public FederationWithinClusterSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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
            // testProxyConnection(m_asReleases[1]);
            }
        }
    }
