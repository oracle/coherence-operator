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
import org.junit.Assume;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Arrays;
import java.util.Collection;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the default-proxy-example.
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
        f_sOperatorChartURL  = sOperatorChartURL;
        f_sCoherenceChartURL = sCoherenceChartURL;
        System.err.println("Operator Chart:   " + f_sOperatorChartURL + "\nCoherence Chart: " + f_sCoherenceChartURL);
        }

    @Parameterized.Parameters
    public static Collection testParameters()
        {
        return Arrays.asList(new Object[][] {
            {
            "https://oracle.github.io/coherence-operator/charts/coherence-operator-0.9.4.tgz",
            "https://oracle.github.io/coherence-operator/charts/coherence-0.9.4.tgz"
            },
            {
            OPERATOR_HELM_CHART_PACKAGE, COHERENCE_HELM_CHART_PACKAGE
            }
            });
        }

    /**
     * Indicates if a particular instance of a test should be run. The case
     * where it does not run is when the OPERATOR_HELM_CHART_URL and
     * COHERENCE_HELM_CHART_URL are both null.
     *
     * @return true if a test should be run
     */
    private boolean testShouldRun()
        {
        return f_sCoherenceChartURL != null && f_sOperatorChartURL != null;
        }

    // ----- test lifecycle -------------------------------------------------

    /**
     * Install the charts required for the test.
     * 
     * @throws Exception
     */
    @Before
    public void installCharts() throws Exception
        {
        Assume.assumeTrue(testShouldRun());

        // install Coherence Operator chart
        s_sOperatorRelease = installOperator("coherence-operator.yaml",toURL(f_sOperatorChartURL));

        // install Coherence chart
        String[] asCohNamespaces = getTargetNamespaces();

        m_asReleases = installCoherence(s_k8sCluster, toURL(f_sCoherenceChartURL), asCohNamespaces,"coherence.yaml");

        assertCoherence(s_k8sCluster, asCohNamespaces, m_asReleases);
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
        Assume.assumeTrue(testShouldRun());

        try (Application application = portForwardExtend(m_asReleases[0]))
            {
            PortMapping              portMapping = application.get(PortMapping.class);
            int                      nPort       = portMapping.getPort().getActualPort();
            ConfigurableCacheFactory ccf         = getCacheFactory("client-cache-config.xml", nPort);

            Eventually.assertThat(invoking(ccf).ensureCache("test", null), is(notNullValue()));

            NamedCache nc = ccf.ensureCache("test", null);
            nc.put("key-1", "value-1");
            assertThat(nc.get("key-1"), is("value-1"));
            ccf.dispose();
            }
        }

    // ----- data members ---------------------------------------------------

    /**
     * Operator chart to run test with.
     */
    private final String f_sOperatorChartURL;

    /**
     * Coherence chart to run test with.
     */
    private final String f_sCoherenceChartURL;
    }
