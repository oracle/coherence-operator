/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.tangosol.net.ConfigurableCacheFactory;
import com.tangosol.net.ExtensibleConfigurableCacheFactory;
import org.junit.After;
import org.junit.Assume;
import org.junit.BeforeClass;

import java.net.URL;

/**
 * Base test class for Proxy samples.
 *
 * @author  tam 2019.05.14
 */
public class BaseProxySampleTest
        extends BaseHelmChartTest
    {

    // ----- test lifecycle -------------------------------------------------
    
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

        String sNamespace   = getK8sNamespace();

        String sRelease = installChart(s_k8sCluster,
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
     * Port forward to a Coherence pod.
     *
     * @param sCoherenceRelease release name
     *                          
     * @return
     * @throws Exception
     */
    protected Application portForwardExtend(String sCoherenceRelease) throws Exception
        {
        return portForwardCoherencePod(s_k8sCluster, getK8sNamespace(), sCoherenceRelease, 20000);
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
