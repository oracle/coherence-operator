/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package helm;

import com.oracle.bedrock.runtime.k8s.K8sCluster;
import org.junit.After;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Test;

/**
 * Verify that the Coherence chart can be installed using default values
 * @author jk  2019.05.28
 */
public class DefaultsIT
        extends BaseHelmChartTest
    {
    @BeforeClass
    public static void setup()
        {
        assertPreconditions(s_k8sCluster);
        ensureNamespace(s_k8sCluster);
        ensureSecret(s_k8sCluster);
        }

    @AfterClass
    public static void cleanup()
        {
        cleanupPullSecrets(s_k8sCluster);
        cleanupNamespace(s_k8sCluster);

        for (String sNamespace : getTargetNamespaces())
            {
            cleanupNamespace(s_k8sCluster, sNamespace);
            }
        }

    @After
    public void cleanUpCoherence()
        {
        if (m_sRelease != null)
            {
            deleteCoherence(s_k8sCluster, m_sNamespace, m_sRelease, false);
            }
        }

    @Test
    public void shouldInstall() throws Exception
        {
        m_sRelease = installChart(s_k8sCluster, COHERENCE_HELM_CHART_NAME, COHERENCE_HELM_CHART_URL, m_sNamespace);

        assertCoherence(s_k8sCluster, m_sNamespace, m_sRelease);
        }

    // ----- data members ---------------------------------------------------

    private static K8sCluster s_k8sCluster = getDefaultCluster();

    private String m_sNamespace = getK8sNamespace();

    private String m_sRelease;
    }
