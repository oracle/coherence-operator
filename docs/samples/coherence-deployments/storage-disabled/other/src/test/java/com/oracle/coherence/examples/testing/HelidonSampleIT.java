/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.console.SystemApplicationConsole;

import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.util.Resources;
import org.hamcrest.MatcherAssert;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.io.File;
import java.util.Collection;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static com.oracle.coherence.examples.testing.HelmUtils.getPods;

import static org.hamcrest.CoreMatchers.is;;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the helidon-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.14
 */
@RunWith(Parameterized.class)
public class HelidonSampleIT
        extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public HelidonSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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

            // install the Helidon chart, which will connect to the Coherence cluster
            String sCoherenceRelease = m_asReleases[0];
            String sNamespace        = getTargetNamespaces()[0];
            String sChartDirName     = "target/webserver-" + System.getProperty("project.version") + "-helm/";

            File fileChartDir = new File(Resources.findFileOrResource(sChartDirName, null).getFile());
            System.err.println("ChartDir=" + fileChartDir.getPath());
            
            int nExitCode = s_helm.install(fileChartDir, "webserver")
                        .namespace(sNamespace)
                        .name("helidon-web-app")
                        .set("wka=" + sCoherenceRelease + "-coherence-headless")
                        .set("clusterName=helidon-cluster")
                        .executeAndWait(SystemApplicationConsole.builder());
            assertThat("Helm install failed for helidon-web-app", nExitCode, is(0));

             // validate that the metrics turn up in prometheus
            // retrieve the client tier release pod
            String       sCoherenceSelector = getCoherencePodSelector(sCoherenceRelease);
            List<String> listPods           = getPods(s_k8sCluster, sNamespace, sCoherenceSelector);

            MatcherAssert.assertThat(listPods.size(), is(2));
            String sCoherencePod = listPods.get(0);

            // wait for the following message to appear, which indicates the helidon web app has joined the cluster
            // Role=OracleCoherenceExamplesMain
            Eventually.assertThat(invoking(this).hasLogMessageAppeared(s_k8sCluster, sNamespace, sCoherencePod,
                    "Role=OracleCoherenceExamplesMain"),
                    is(true), Timeout.after(300, TimeUnit.SECONDS), InitialDelay.of(20, TimeUnit.SECONDS));
            }
        }
    
    @After
    public void cleanupHelidonRelease()
        {
        try
            {
            if (testShouldRun())
                {
                cleanupHelmReleases("helidon-web-app");
                }
            }
        catch (Exception e)
            {
            // ignore
            }
        }

    // ----- tests ----------------------------------------------------------

    /**
     * Test the helidon-sample.
     *
     * @throws Exception
     */
    @Test
    public void testHelidonSample()
        {
        if (testShouldRun())
            {
            try {
                System.err.println("Success");
            }
            finally
                {
                cleanupHelmReleases("helidon-web-app");
                }
            }
        }

    // ----- constants ------------------------------------------------------

    /**
     * Selector for helidon webapp.
     */
    private static final String SELECTOR = "app=webserver,component=WebServerPod";

    /**
     * Query to run against app.
     */
    private static final String QUERY = "/query -d '{\"query\":\"create cache foo\"}";
    }

