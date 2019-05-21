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

/**
 * Test the elastic-data-sample-default.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.21
 */
@RunWith(Parameterized.class)
public class DefaultElasticDataSampleIT
        extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public DefaultElasticDataSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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
            }
        }

    // ----- tests ----------------------------------------------------------

        /**
     * Test the proxy tier sample.
     *
     * @throws Exception
     */
    @Test
    public void testMultipleProxiesSample() throws Exception
        {
        if (testShouldRun())
            {
            // connect to proxy tier - m_asReleases[1] and test both proxy ports
            testProxyConnection(m_asReleases[0], 20000, "flash-01");
            }
        }
    }

