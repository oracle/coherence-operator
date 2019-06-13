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
import com.oracle.bedrock.runtime.console.SystemApplicationConsole;
import com.oracle.bedrock.runtime.k8s.helm.HelmInstall;
import com.oracle.bedrock.runtime.k8s.helm.HelmUpgrade;
import com.oracle.bedrock.testsupport.deferred.Eventually;

import com.tangosol.net.ConfigurableCacheFactory;
import com.tangosol.net.NamedCache;
import org.junit.Before;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import javax.naming.Name;
import java.io.File;
import java.util.Collection;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test the rolling-upgrade-sample.
 *
 * Any changes to the arguments of the helm install commands in the README.md, should be
 * also made to the corresponding yaml files in test/resources.
 *
 * @author tam  2019.05.22
 */
@RunWith(Parameterized.class)
public class RollingUpgradeSampleIT
       extends BaseSampleTest
    {
    // ----- constructor ----------------------------------------------------

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public RollingUpgradeSampleIT(String sOperatorChartURL, String sCoherenceChartURL)
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

            String sTag = System.getProperty("docker.push.tag.prefix") + System.getProperty("project.artifactId") + ":1.0.0";

            // process yaml file to replace user artifacts image
            String sProcessedCoherenceYaml = getProcessedYamlFile("coherence-v1.0.0.yaml", sTag, null);
            assertThat(sProcessedCoherenceYaml, is(notNullValue()));

            String sClusterRelease = installCoherence(s_k8sCluster, toURL(m_sCoherenceChartURL), sCohNamespace, sProcessedCoherenceYaml);
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
    @Ignore("Needs more investigation")
    public void testRollingUpgradeSample() throws Exception
        {
        if (testShouldRun())
            {
            String sCoherenceRelease = m_asReleases[0];
            String sCohNamespace     = getTargetNamespaces()[0];

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
                assertThat(nc.get("key-1"), is("value-1"));
                ccf.dispose();
                }

            // Next, issue a helm update to update the user artifacts image to 2.0.0
            String sTag = System.getProperty("docker.push.tag.prefix") + System.getProperty("project.artifactId") + ":2.0.0";

            File fileChartDir = extractChart("coherence", toURL(m_sCoherenceChartURL));
            
            int nExitCode = s_helm.upgrade(sCoherenceRelease, fileChartDir.getPath() + "/coherence")
                        .namespace(sCohNamespace)
                        .withFlags("--reuse-values")
                        .set("imagePullSecrets=ocr-k8s-operator-development-secret")
                        .set("userArtifacts.image=" + sTag)
                        .set("userArtifacts.imagePullPolicy=Never")
                        .executeAndWait(SystemApplicationConsole.builder());

            assertThat("Helm upgrade failed", nExitCode, is(0));

            Eventually.assertThat(invoking(this).isValueUppercase(sCoherenceRelease),
                  is(true),
                  Timeout.after(600, TimeUnit.SECONDS),
                  InitialDelay.of(180, TimeUnit.SECONDS),
                  RetryFrequency.every(10, TimeUnit.SECONDS));
            }
        }

    /**
     * Returns true when the value of key is uppercase.
     *
     * @param sCoherenceRelease release to check
     *                              
     * @return {@code true} when the value of key is uppercase
     */
    public boolean isValueUppercase(String sCoherenceRelease)
        {
        try (Application application = portForwardExtend(sCoherenceRelease, 20000))
            {
            PortMapping portMapping = application.get(PortMapping.class);
            int         nActualPort = portMapping.getPort().getActualPort();
            
            ConfigurableCacheFactory ccf  = getCacheFactory("client-cache-config.xml", nActualPort);
            Eventually.assertThat(invoking(ccf).ensureCache("test", null), is(notNullValue()));
            
            NamedCache nc = ccf.ensureCache("test", null);

            // this should always succeed
            assertThat(nc.get("key-1"), is("value-1"));
            
            nc.put("key-2", "value-2");
            boolean fResult = "VALUE-2".equals(nc.get("key-2"));
            System.err.println("value for key-2 is " + nc.get("key-2"));

            ccf.dispose();

            return fResult;
            }
        catch (Exception e)
            {
            // ignore exceptions as they may be because a pod is gone.
            System.err.println("Unable to run portForward");
            e.printStackTrace();
            return false;
            }
        }
    }
