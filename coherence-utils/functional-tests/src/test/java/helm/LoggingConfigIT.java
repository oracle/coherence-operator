/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.util.Resources;
import org.junit.*;

import java.io.File;
import java.net.URL;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static helm.HelmUtils.HELM_TIMEOUT;
import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test installing the Coherence chart with various logging configuration settings.
 *
 * @author jk  2019.02.12
 */
public class LoggingConfigIT
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

    @Test
    public void shouldSetCoherenceLogLevel() throws Exception
        {
        String   sNamespace  = getK8sNamespace();
        String   sValues     = "values/helm-values-coh-user-artifacts.yaml";
        String[] asSetValues = {"clusterSize=1", "store.logging.level=9"};

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValues, asSetValues);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        String       sSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPod   = getPods(s_k8sCluster, sNamespace, sSelector);

        assertThat(listPod.isEmpty(), is(false));

        Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, sNamespace, listPod.get(0)),
            is(true), Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

        Queue<String> queueLog = getPodLog(s_k8sCluster, sNamespace, listPod.get(0));
        boolean       fEnvVar  = queueLog.stream().anyMatch(s -> s.equals("COH_LOG_LEVEL=9"));
        boolean       fSysProp = queueLog.stream().anyMatch(s -> s.contains(" -Dcoherence.log.level=9 "));

        assertThat("Env variable COH_LOGGING_CONFIG is not correct", fEnvVar, is(true));
        assertThat("System property java.util.logging.config.file is not correct", fSysProp, is(true));
        }

    @Test
    public void shouldUseDefaultLoggingConfiguration() throws Exception
        {
        String   sNamespace  = getK8sNamespace();
        String   sValues     = "values/helm-values-coh.yaml";
        String[] asSetValues = {"clusterSize=1"};
        String   sLogConfig  = "/scripts/logging.properties";

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValues, asSetValues);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        String       sSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPod   = getPods(s_k8sCluster, sNamespace, sSelector);

        assertThat(listPod.isEmpty(), is(false));

        Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, sNamespace, listPod.get(0)),
            is(true), Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

        Queue<String> queueLog = getPodLog(s_k8sCluster, sNamespace, listPod.get(0));
        boolean       fEnvVar  = queueLog.stream().anyMatch(s -> s.equals("COH_LOGGING_CONFIG=" + sLogConfig));
        boolean       fSysProp = queueLog.stream().anyMatch(s -> s.contains(" -Djava.util.logging.config.file=" + sLogConfig));

        assertThat("Env variable COH_LOGGING_CONFIG is not correct", fEnvVar, is(true));
        assertThat("System property java.util.logging.config.file is not correct", fSysProp, is(true));
        }


    @Test
    public void shouldUseCustomLoggingConfiguration() throws Exception
        {
        String   sNamespace  = getK8sNamespace();
        String   sValues     = "values/helm-values-coh-user-artifacts.yaml";
        String   sCfgDir     = "/u01/oracle/oracle_home/coherence/ext/conf/";
        String   sLogConfig  = "custom-logging.properties";
        String[] asSetValues = {"clusterSize=1", "store.logging.configFile=" + sLogConfig};

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValues, asSetValues);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        String       sSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPod   = getPods(s_k8sCluster, sNamespace, sSelector);

        assertThat(listPod.isEmpty(), is(false));

        Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, sNamespace, listPod.get(0)),
            is(true), Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

        Queue<String> queueLog     = getPodLog(s_k8sCluster, sNamespace, listPod.get(0));
        String        sExpectedCfg = sCfgDir + sLogConfig;
        boolean       fEnvVar      = queueLog.stream().anyMatch(s -> s.equals("COH_LOGGING_CONFIG=" + sExpectedCfg));
        boolean       fSysProp     = queueLog.stream().anyMatch(s -> s.contains(" -Djava.util.logging.config.file=" + sExpectedCfg));

        assertThat("Env variable COH_LOGGING_CONFIG is not correct", fEnvVar, is(true));
        assertThat("System property java.util.logging.config.file is not correct", fSysProp, is(true));
        }

    @Test
    public void shouldUseLoggingConfigurationFromConfigMap() throws Exception
        {
        String sNamespace  = getK8sNamespace();
        String sCfgDir     = "/loggingconfig/";
        String sLogConfig  = "test-logging.properties";
        String sConfigMap  = "coh-log-config";
        URL    url         = Resources.findFileOrResource(sLogConfig, null);
        File   file        = new File(url.toURI());

        // ensure that the ConfigMap does not exist
        s_k8sCluster.kubectlAndWait(Arguments.of("-n",
                                                 sNamespace,
                                                 "delete",
                                                 "configmap",
                                                 sConfigMap,
                                                 "--ignore-not-found=true"));

        // Create the config map containing the logging config
        int nExitCode = s_k8sCluster.kubectlAndWait(Arguments.of("-n",
                                                                 sNamespace,
                                                                 "create",
                                                                 "configmap",
                                                                 sConfigMap,
                                                                 "--from-file=" + file.getAbsolutePath()));

        assertThat("Error creating logging configuration", nExitCode, is(0));

        // install Coherence
        String   sValues     = "values/helm-values-coh-user-artifacts.yaml";
        String[] asSetValues = {"clusterSize=1", "store.logging.configMapName=" + sConfigMap, "store.logging.configFile=" + sLogConfig};

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValues, asSetValues);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        String       sSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPod   = getPods(s_k8sCluster, sNamespace, sSelector);

        assertThat(listPod.isEmpty(), is(false));

        Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, sNamespace, listPod.get(0)),
            is(true), Timeout.after(HELM_TIMEOUT, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

        Queue<String> queueLog     = getPodLog(s_k8sCluster, sNamespace, listPod.get(0));
        String        sExpectedCfg = sCfgDir + sLogConfig;
        boolean       fEnvVar      = queueLog.stream().anyMatch(s -> s.equals("COH_LOGGING_CONFIG=" + sExpectedCfg));
        boolean       fSysProp     = queueLog.stream().anyMatch(s -> s.contains(" -Djava.util.logging.config.file=" + sExpectedCfg));

        assertThat("Env variable COH_LOGGING_CONFIG is not correct", fEnvVar, is(true));
        assertThat("System property java.util.logging.config.file is not correct", fSysProp, is(true));
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
    }
