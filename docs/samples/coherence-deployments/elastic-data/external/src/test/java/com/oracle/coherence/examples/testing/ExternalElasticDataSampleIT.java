/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.util.Resources;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Collection;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static com.oracle.coherence.examples.testing.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the elastic-data-sample-external.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.21
 */
@RunWith(Parameterized.class)
public class ExternalElasticDataSampleIT
    extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public ExternalElasticDataSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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

            String sCohNamespace = getTargetNamespaces()[0];

            String sTag = System.getProperty("docker.push.tag.prefix") + System.getProperty("project.artifactId") + ":" +
                          System.getProperty("project.version");

            // process yaml file to replace user artifacts image
            String sProcessedCoherenceYaml = getProcessedYamlFile("coherence.yaml", sTag, null);
            assertThat(sProcessedCoherenceYaml, is(notNullValue()));

            // append the contents of the src/main/yaml/volumes.yaml to the sProcessedCoherenceYaml file
            String sVolumesYaml = Resources.findFileOrResource("src/main/yaml/volumes.yaml", null).getPath();
            assertThat(sVolumesYaml, is(notNullValue()));

            String sYamlContent1 = new String(Files.readAllBytes(Paths.get(sProcessedCoherenceYaml)));
            String sYamlContent2 = new String(Files.readAllBytes(Paths.get(sVolumesYaml)));

            File fileTemp = File.createTempFile("final-processed-file",".yaml");
            StringBuilder sb = new StringBuilder(sYamlContent1).append('\n').append(sYamlContent2);

            String sNewYamlFile = fileTemp.getPath();
            Files.write(Paths.get(sNewYamlFile), sb.toString().getBytes());

            // install Coherence chart
            String sClusterRelease = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), sCohNamespace, sNewYamlFile);
            assertCoherence(s_k8sCluster, sCohNamespace, sClusterRelease);

            m_asReleases = new String[] { sClusterRelease };
            }
        }

    // ----- tests ----------------------------------------------------------

    /**
     * Test the proxy tier sample.
     *
     * @throws Exception
     */
    @Test
    public void testExternalElasticDataSample() throws Exception
        {
        if (testShouldRun())
            {
            // connect to proxy tier - m_asReleases[0]
            testProxyConnection(m_asReleases[0], 20000, "flash-01");


            // retrieve the client tier release pod
            String       sNamespace         = getTargetNamespaces()[0];
            String       sCoherenceSelector = getCoherencePodSelector(m_asReleases[0]);
            List<String> listPods           = getPods(s_k8sCluster, getTargetNamespaces()[0], sCoherenceSelector);

            assertThat(listPods.size(), is(1));
            String sCoherencePod = listPods.get(0);

            // wait for the following message to appear, which indicates the server is using /elastic-data
            // TODO: Fix
//            Eventually.assertThat(invoking(this).hasLogMessageAppeared(s_k8sCluster, sNamespace, sCoherencePod, "/elastic-data"),
//                    is(true), Timeout.after(120, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));
            }
        }
    }

