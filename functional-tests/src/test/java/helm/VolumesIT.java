/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.runtime.options.Console;
import com.tangosol.util.Resources;
import org.junit.*;
import org.junit.rules.TemporaryFolder;

import java.io.File;
import java.io.PrintWriter;
import java.net.URL;
import java.sql.SQLOutput;
import java.util.List;
import java.util.Queue;

import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test installing the Coherence chart with various logging configuration settings.
 *
 * @author jk  2019.02.12
 */
@Ignore
public class VolumesIT
        extends BaseHelmChartTest
    {
    // ----- test lifecycle -------------------------------------------------

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
            deleteCoherence(s_k8sCluster, getK8sNamespace(), m_sRelease, false);
            }
        }

    // ----- test methods ---------------------------------------------------

    /**
     * Add a custom volume mount that maps to a custom config map and
     * verify that the Coherence Pod can see that data from the config map.
     *
     * @throws Exception  if the test fails
     */
    @Test
    public void shouldMapExtraVolume() throws Exception
        {
        String   sNamespace = getK8sNamespace();
        String   sConfigMap = "test-mount";
        File     file       = s_temporaryFolder.newFile();
        String   sData      = "test data";
        String   sValues    = "values/helm-values-extra-volumes.yaml";

        try (PrintWriter writer = new PrintWriter(file))
            {
            writer.println(sData);
            }

        s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "delete", "configmap", sConfigMap));
        s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "create", "configmap",
                                                 sConfigMap, "--from-file", file.getCanonicalPath()));

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValues);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        String       sSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPod   = getPods(s_k8sCluster, sNamespace, sSelector);

        assertThat(listPod.isEmpty(), is(false));

        String                      sPodName  = listPod.get(0);
        String                      sFileName = "/extra-data/" + file.getName();
        CapturingApplicationConsole console   = new CapturingApplicationConsole();
        int                         nExitCode = s_k8sCluster.kubectlAndWait(Console.of(console),
                                                                            Arguments.of("-n", sNamespace,
                                                                                         "exec", sPodName,
                                                                                         "cat", sFileName));

        assertThat(nExitCode, is(0));
        assertThat(console.getCapturedOutputLines().poll(), is(sData));
        }

    /**
     * Add a custom PVC mount.
     *
     * @throws Exception  if the test fails
     */
    @Test
    public void shouldMapExtraPersistentVolumeClaim() throws Exception
        {
        String   sNamespace = getK8sNamespace();
        String   sValues    = "values/helm-values-extra-pvc.yaml";

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValues);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        String       sSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPod   = getPods(s_k8sCluster, sNamespace, sSelector);

        assertThat(listPod.isEmpty(), is(false));

        String                      sPodName  = listPod.get(0);
        CapturingApplicationConsole console   = new CapturingApplicationConsole();
        int                         nExitCode = s_k8sCluster.kubectlAndWait(Console.of(console),
                                                                            Arguments.of("-n", sNamespace,
                                                                                         "exec", sPodName, "-it",
                                                                                         "stat", "/extra-data"));

        assertThat(nExitCode, is(0));
        }

    // ----- data members ---------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The name of the deployed Coherence Helm release.
     */
    private String m_sRelease;

    /**
     * JUnit rule to create temporary files and folders.
     */
    @ClassRule
    public static final TemporaryFolder s_temporaryFolder = new TemporaryFolder();
    }
