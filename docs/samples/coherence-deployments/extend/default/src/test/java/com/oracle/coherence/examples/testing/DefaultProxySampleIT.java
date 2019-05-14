/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.runtime.k8s.K8sCluster;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Ignore;
import org.junit.Test;

@Ignore
public class DefaultProxySampleIT
        extends BaseHelmChartTest
    {

    // ----- test lifecycle -------------------------------------------------

    @BeforeClass
    public static void setup() throws Exception
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

    @AfterClass
    public static void cleanup()
        {
        if (s_sOperatorRelease != null)
            {
            try
                {
                capturePodLogs(DefaultProxySampleIT.class, s_k8sCluster, getCoherenceOperatorSelector(s_sOperatorRelease),
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

    // ----- tests ----------------------------------------------------------


    @Test
    public void testDefaultProxySample()
        {
        System.out.println("Ready");
        }

    // ----- constants ------------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The name of the deployed Operator Helm release.
     */
    private static String s_sOperatorRelease;

    /**
     * The boolean indicates whether Coherence cache data is persisted.
     */
    private static boolean PERSISTENCE = false;

    /**
     * The name of the deployed Coherence Helm releases.
     */
    private String[] m_asReleases;
    }
