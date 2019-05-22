/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.net.ConfigurableCacheFactory;
import com.tangosol.net.ExtensibleConfigurableCacheFactory;
import com.tangosol.net.NamedCache;
import org.junit.After;
import org.junit.BeforeClass;

import java.net.URL;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Base test class for samples.
 *
 * @author  tam 2019.05.14
 */
public class BaseSampleTest
        extends BaseHelmChartTest
    {

    // ----- constructors ---------------------------------------------------

    /**
     * No-args constructor.
     */
    public BaseSampleTest()
        {
        }

    /**
     * Constructor for Parameterized test
     *
     * @param sOperatorChartURL   Operator chart URL
     * @param sCoherenceChartURL  Coherence chart URL
     */
    public BaseSampleTest(String sOperatorChartURL, String sCoherenceChartURL)
        {
        m_sOperatorChartURL  = sOperatorChartURL;
        m_sCoherenceChartURL = sCoherenceChartURL;
        System.err.println("Operator Chart:   " + m_sOperatorChartURL + "\nCoherence Chart: " + m_sCoherenceChartURL);
        }

    // ----- test lifecycle -------------------------------------------------

    protected static void cleanupExtra(K8sCluster cluster)
        {
        }

    /**
     * Ensure Kubernetes cluster and namespaces are setup.
     *
     */
    @BeforeClass
    public static void setupNamespaceAndOperator()
        {
        assertPreconditions(s_k8sCluster);
        ensureNamespace(s_k8sCluster);
        ensureSecret(s_k8sCluster);

        String sOpNamespace = getK8sNamespace();
        for (String sNamespace : getTargetNamespaces())
            {
            if (sNamespace != null && !sNamespace.equals(sOpNamespace))
                {
                ensureNamespace(s_k8sCluster, sNamespace);
                }
            }
        }

    /**
     * Cleanup any charts created by test.
     */
    @After
    public void cleanupCharts()
        {
        if (m_asReleases != null && m_asReleases.length > 0)
            {
            for (int i = 0; i < m_asReleases.length; i++)
                {
                deleteCoherence(s_k8sCluster, getK8sNamespace(), m_asReleases[i], false);
                }
            }

        if (s_sOperatorRelease != null)
            {
            try
                {
                capturePodLogs(this.getClass(), s_k8sCluster, getCoherenceOperatorSelector(s_sOperatorRelease),
                    "coherence-operator", "fluentd");
                cleanupHelmReleases(s_sOperatorRelease);
                }
            catch (Throwable t)
                {
                System.err.println("Error in cleaning up Helm Releases: " + s_sOperatorRelease);
                }
            }

        cleanupPullSecrets(s_k8sCluster);
        cleanupNamespace(s_k8sCluster);

        cleanupExtra(s_k8sCluster);

        for (String sNamespace : getTargetNamespaces())
            {
            cleanupNamespace(s_k8sCluster, sNamespace);
            }
        }

    /**
     * Install the Coherence Operator and return the release name.
     *
     * @param sValuesFile  values file for additional chart values
     * @param urlOperator  URL of Coherence Operator to install
     *
     * @return the release name
     *
     * @throws Exception
     */
    protected String installOperator(String sValuesFile, URL urlOperator)
            throws Exception
        {
        System.err.println("Deploying " + OPERATOR_HELM_CHART_NAME + " with " + sValuesFile);

        String sNamespace = getK8sNamespace();
        String sRelease   = installChart(s_k8sCluster,
                                          OPERATOR_HELM_CHART_NAME,
                                          urlOperator,
                                          sNamespace,
                                          sValuesFile,
                                          getDefaultHelmSetValues());

        assertDeploymentReady(s_k8sCluster, sNamespace, getCoherenceOperatorSelector(sRelease), true);

        return sRelease;
        }

    // ----- helpers --------------------------------------------------------

    /**
     * Install operator and Coherence charts for a single tier install.
     */
    protected void installChartsSingleTier(String sOperatorChartURL, String sCoherenceChartURL) throws Exception
        {
        // install Coherence Operator chart
        s_sOperatorRelease = installOperator("coherence-operator.yaml",toURL(sOperatorChartURL));

        // install Coherence chart
        String sCohNamespace = getTargetNamespaces()[0];

        String sTag = System.getProperty("docker.push.tag.prefix") + System.getProperty("project.artifactId") + ":" +
                      System.getProperty("project.version");

        // process yaml file to replace user artifacts image
        String sProcessedCoherenceYaml = getProcessedYamlFile("coherence.yaml", sTag, null);
        assertThat(sProcessedCoherenceYaml, is(notNullValue()));

        String sClusterRelease = installCoherence(s_k8sCluster, toURL(sCoherenceChartURL), sCohNamespace, sProcessedCoherenceYaml);
        assertCoherence(s_k8sCluster, sCohNamespace, sClusterRelease);

        m_asReleases = new String[] { sClusterRelease };
        }

    /**
     * Test a given Proxy connection to a release and default port and cache name.
     *
     * @param sRelease  release to connect to
     *
     * @throws Exception
     */
    protected void testProxyConnection(String sRelease) throws Exception
        {
        testProxyConnection(sRelease, 20000, "test");
        }

    /**
     * Test a given Proxy connection to a release, default cache name and a port.
     *
     * @param sRelease  release to connect to
     *
     * @throws Exception
     */
    protected void testProxyConnection(String sRelease, int nPort) throws Exception
        {
        testProxyConnection(sRelease, nPort, "test");
        }
    
    /**
     * Test a given Proxy connection to a release and specific port.
     *
     * @param sRelease  release to connect to
     * @param nPort     port to connect to
     * @param sCache    cache to create
     *
     * @throws Exception
     */
    protected void testProxyConnection(String sRelease, int nPort, String sCache) throws Exception
        {
        try (Application application = portForwardExtend(sRelease, nPort))
            {
            PortMapping              portMapping = application.get(PortMapping.class);
            int                      nActualPort = portMapping.getPort().getActualPort();
            ConfigurableCacheFactory ccf         = getCacheFactory("client-cache-config.xml", nActualPort);

            Eventually.assertThat(invoking(ccf).ensureCache(sCache, null), is(notNullValue()));

            NamedCache nc = ccf.ensureCache(sCache, null);
            nc.put("key-1", "value-1");
            assertThat(nc.get("key-1"), is("value-1"));
            ccf.dispose();
            }
        }
    /**
     * Indicates if a particular instance of a test should be run. The case
     * where it does not run is when the OPERATOR_HELM_CHART_URL and
     * COHERENCE_HELM_CHART_URL are both null.
     *
     * @return true if a test should be run
     */
    protected boolean testShouldRun()
        {
        return m_sCoherenceChartURL != null && m_sCoherenceChartURL.length() > 0 &&
               m_sOperatorChartURL  != null && m_sOperatorChartURL.length() > 0;
        }

    /**
     * Port forward to a Coherence pod on the default extend port.
     *
     * @param sCoherenceRelease release name
     *                          
     * @return the {@link Application} to use
     * @throws Exception
     */
    protected Application portForwardExtend(String sCoherenceRelease) throws Exception
        {
        return portForwardCoherencePod(s_k8sCluster, getK8sNamespace(), sCoherenceRelease, 20000);
        }

    /**
     * Port forward to a Coherence pod on a specific port.
     *
     * @param sCoherenceRelease release name
     * @param nPort             port
     *
     * @return the {@link Application} to use
     * @throws Exception
     */
    protected Application portForwardExtend(String sCoherenceRelease, int nPort) throws Exception
        {
        return portForwardCoherencePod(s_k8sCluster, getK8sNamespace(), sCoherenceRelease, nPort);
        }

    /**
     * Returns a {@link ConfigurableCacheFactory} for a given cache config and port
     * 
     * @param sCacheConfig  cache config to load
     * @param nPort         port to connect to (127.0.0.1 is assumed for host)
     * @return
     */
    protected ConfigurableCacheFactory getCacheFactory(String sCacheConfig, int nPort)
        {
        System.setProperty("coherence.serializer", "java");
        System.setProperty("proxy.address", "127.0.0.1");
        System.setProperty("proxy.port", String.valueOf(nPort));

        ExtensibleConfigurableCacheFactory.Dependencies deps
                = ExtensibleConfigurableCacheFactory.DependenciesHelper.newInstance(sCacheConfig);

        return new ExtensibleConfigurableCacheFactory(deps);
        }


    // ----- data members ---------------------------------------------------

    /**
     * Operator chart to run test with.
     */
    protected String m_sOperatorChartURL;

    /**
     * Coherence chart to run test with.
     */
    protected String m_sCoherenceChartURL;

    // ----- constants ------------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    protected static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The name of the deployed Operator Helm release.
     */
    protected static String s_sOperatorRelease;

    /**
     * The name of the deployed Coherence Helm releases.
     */
    protected String[] m_asReleases;
    }
