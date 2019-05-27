/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.runtime.Application;

import com.oracle.coherence.examples.SampleClient;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Collection;

/**
 * Test the sidecar-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.14
 */
@RunWith(Parameterized.class)
public class SidecarSampleIT
    extends BaseSampleTest
    {
     // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public SidecarSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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

    @Test
    public void testSidecarSample()
        {
        if (testShouldRun())
            {
            try (Application application = portForwardExtend(m_asReleases[0], 20000))
                {
                PortMapping portMapping = application.get(PortMapping.class);
                int         nActualPort = portMapping.getPort().getActualPort();

                System.err.println("Started: " + application.getName());

                System.setProperty("proxy.address", "127.0.0.1");
                System.setProperty("proxy.port", Integer.toString(nActualPort));
                System.setProperty("coherence.pof.config", "conf/storage-pof-config.xml");
                System.setProperty("coherence.cacheconfig", "client-cache-config.xml");
                System.setProperty("coherence.tcmpenabled", "false");
                System.setProperty("coherence.distributed.localstorage", "false");
                
                SampleClient.main();
                }
            catch (Exception e)
                {
                e.printStackTrace();
                }
            }
        }
    }
