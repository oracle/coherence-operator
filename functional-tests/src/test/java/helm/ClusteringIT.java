/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.console.CapturingApplicationConsole;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.options.Arguments;
import com.oracle.bedrock.runtime.options.Console;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.util.Resources;
import org.junit.*;
import org.junit.rules.TemporaryFolder;

import java.io.File;
import java.io.PrintWriter;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.junit.Assert.fail;

/**
 * @author jk  2019.02.13
 */
public class ClusteringIT
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
        m_listRelease.stream()
                .filter(Objects::nonNull)
                .forEach(sRelease -> deleteCoherence(s_k8sCluster, getK8sNamespace(), sRelease, false));

        if (m_headlessService != null)
            {
            s_k8sCluster.kubectlAndWait(Arguments.of("-n", getK8sNamespace(), "delete", "svc", m_headlessService));
            }
        }

    // ----- test methods ---------------------------------------------------

    @Test
    public void shouldFormClusterWithHelmRelease() throws Exception
        {
        String sNamespace      = getK8sNamespace();
        String sValClusterSize = "clusterSize=2";
        String sValClusterName = "cluster=foo";
        String sValEnableJmx   = "store.jmx.enabled=true";
        String sValuesFile     = "values/helm-values-coh.yaml";
        String sReleaseOne     = installCoherence(s_k8sCluster, sNamespace, sValuesFile, sValClusterName, sValClusterSize, sValEnableJmx);

        m_listRelease.add(sReleaseOne);

        assertCoherence(s_k8sCluster, sNamespace, sReleaseOne);
        assertCoherenceJMX(s_k8sCluster, sNamespace, sReleaseOne);

        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(s_k8sCluster, sNamespace, sReleaseOne),
            greaterThanOrEqualTo(3), RetryFrequency.every(10, TimeUnit.SECONDS), Timeout.after(60, TimeUnit.SECONDS));

        String[] asSetValues = {sValClusterName, sValClusterSize, "store.wkaRelease=" + sReleaseOne};
        String   sReleaseTwo = installCoherence(s_k8sCluster, sNamespace, null, asSetValues);

        m_listRelease.add(sReleaseTwo);

        assertCoherence(s_k8sCluster, sNamespace, sReleaseTwo);

        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(s_k8sCluster, sNamespace, sReleaseOne),
            greaterThanOrEqualTo(5), RetryFrequency.every(10, TimeUnit.SECONDS), Timeout.after(60, TimeUnit.SECONDS));
        }

    @Test
    public void shouldFormClusterUsingExistingHeadlessService() throws Exception
        {
        String   sNamespace      = getK8sNamespace();
        String   sClusterName    = "foo";
        String   sWkaService     = "foo-wka-service";
        String   sValClusterSize = "clusterSize=2";
        String   sValClusterName = "cluster=" + sClusterName;
        String   sWka            = "store.wka=" + sWkaService;
        String   sValEnableJmx   = "store.jmx.enabled=true";
        String[] asSetArgsOne    = {sValClusterName, sValClusterSize, sWka, sValEnableJmx};
        String   sValuesFile     = "values/helm-values-coh.yaml";

        createHeadlessService(sNamespace, sWkaService, sClusterName);

        String sReleaseOne = installCoherence(s_k8sCluster, sNamespace, sValuesFile, asSetArgsOne);

        m_listRelease.add(sReleaseOne);

        assertCoherence(s_k8sCluster, sNamespace, sReleaseOne);
        assertCoherenceJMX(s_k8sCluster, sNamespace, sReleaseOne);

        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(s_k8sCluster, sNamespace, sReleaseOne),
            greaterThanOrEqualTo(3), RetryFrequency.every(10, TimeUnit.SECONDS), Timeout.after(60, TimeUnit.SECONDS));

        List<String> listPod           = getPods(s_k8sCluster, sNamespace, getCoherencePodSelector(sReleaseOne));
        List<String> listPodsInService = getPodsInService(sNamespace, sWkaService);

        assertThat(listPodsInService, is(listPod));

        String[] asSetArgsTwo = {sValClusterName, sValClusterSize, sWka};
        String sReleaseTwo = installCoherence(s_k8sCluster, sNamespace, sValuesFile, asSetArgsTwo);

        m_listRelease.add(sReleaseTwo);

        assertCoherence(s_k8sCluster, sNamespace, sReleaseTwo);

        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(s_k8sCluster, sNamespace, sReleaseOne),
            greaterThanOrEqualTo(5), RetryFrequency.every(10, TimeUnit.SECONDS), Timeout.after(60, TimeUnit.SECONDS));

        listPod.addAll(getPods(s_k8sCluster, sNamespace, getCoherencePodSelector(sReleaseTwo)));
        Collections.sort(listPod);

        listPodsInService = getPodsInService(sNamespace, sWkaService);

        assertThat(listPodsInService, is(listPod));
        s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "delete", "svc", sWkaService));
        }

    // Integration test for RetryingWkaAddressProvider that waits for wka dns entry to come up
    @Test
    public void shouldFormClusterUsingDeferredHeadlessService() throws Exception
        {
        String   sNamespace      = getK8sNamespace();
        String   sClusterName    = "bar";
        String   sWkaService     = "bar-wka-service";
        String   sValClusterSize = "clusterSize=2";
        String   sValClusterName = "cluster=" + sClusterName;
        String   sWka            = "store.wka=" + sWkaService;
        String   sValEnableJmx   = "store.jmx.enabled=true";
        String[] asSetArgsOne    = {sValClusterName, sValClusterSize, sWka, sValEnableJmx};
        String   sValuesFile     = "values/helm-values-coh.yaml";

        // rather than start headless service referenced in wka first,
        // ensure that cluster can form despite delay in headless service creation.
        // simulates behavior when coherence headless service does not start immediately waiting for volume claims.
        ScheduledExecutorService scheduler = Executors.newSingleThreadScheduledExecutor();
        Runnable task = new Runnable()
            {
            public void run()
                {
                try
                    {
                    createHeadlessService(sNamespace, sWkaService, sClusterName);
                    }
                catch (Throwable t)
                    {
                    fail("delayed start of wka service bar-wka-service failed with exception " + t);
                    t.printStackTrace();
                    }
                }
            };

        scheduler.schedule(task, 45, TimeUnit.SECONDS);
        scheduler.shutdown();

        String sReleaseOne = installCoherence(s_k8sCluster, sNamespace, sValuesFile, asSetArgsOne);

        m_listRelease.add(sReleaseOne);

        assertCoherence(s_k8sCluster, sNamespace, sReleaseOne);
        assertCoherenceJMX(s_k8sCluster, sNamespace, sReleaseOne);

        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(s_k8sCluster, sNamespace, sReleaseOne),
            greaterThanOrEqualTo(3), RetryFrequency.every(10, TimeUnit.SECONDS), Timeout.after(60, TimeUnit.SECONDS));

        List<String> listPod           = getPods(s_k8sCluster, sNamespace, getCoherencePodSelector(sReleaseOne));
        List<String> listPodsInService = getPodsInService(sNamespace, sWkaService);

        assertThat(listPodsInService, is(listPod));

        String[] asSetArgsTwo = {sValClusterName, sValClusterSize, sWka};
        String sReleaseTwo = installCoherence(s_k8sCluster, sNamespace, sValuesFile, asSetArgsTwo);

        m_listRelease.add(sReleaseTwo);

        assertCoherence(s_k8sCluster, sNamespace, sReleaseTwo);

        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(s_k8sCluster, sNamespace, sReleaseOne),
            greaterThanOrEqualTo(5), RetryFrequency.every(10, TimeUnit.SECONDS), Timeout.after(60, TimeUnit.SECONDS));

        listPod.addAll(getPods(s_k8sCluster, sNamespace, getCoherencePodSelector(sReleaseTwo)));
        Collections.sort(listPod);

        listPodsInService = getPodsInService(sNamespace, sWkaService);

        assertThat(listPodsInService, is(listPod));
        }


    // ----- helper methods -------------------------------------------------

    @SuppressWarnings("unchecked")
    private List<String> getPodsInService(String sNamespace, String sServiceName) throws Exception
        {
        CapturingApplicationConsole console  = new CapturingApplicationConsole();
        List<String>                listPods = new ArrayList<>();

        int nExitCode = s_k8sCluster.kubectlAndWait(Console.of(console),
                                Arguments.of("-n", sNamespace, "get", "endpoints", sServiceName, "-o", "json"));

        if (nExitCode != 0)
            {
            HelmUtils.logConsoleOutput("kubectl-get-endpoints", console);
            fail("Kubectl get endpoints failed");
            }

        String    sJson = String.join("\n", console.getCapturedOutputLines());
        Map       map   = MAPPER.readValue(sJson, LinkedHashMap.class);
        List<Map> list  = (List) map.get("subsets");

        if (list != null)
            {
            for (Map mapSubset : list)
                {
                List<Map> listAddresses = (List) mapSubset.get("addresses");

                if (listAddresses != null)
                    {
                    for (Map mapAddress : listAddresses)
                        {
                        Map    mapTarget = (Map) mapAddress.get("targetRef");
                        String sPod      = (String) mapTarget.get("name");

                        if (sPod != null)
                            {
                            listPods.add(sPod);
                            }
                        }
                    }
                }
            }

        Collections.sort(listPods);

        return listPods;
        }

    private void createHeadlessService(String sNamespace, String sServiceName, String sClusterName) throws Exception
        {
        URL          url   = Resources.findFileOrResource("headlessService.yaml", null);
        List<String> lines = Files.readAllLines(Paths.get(url.toURI()));
        File         file  = s_temporaryFolder.newFile();

        try (PrintWriter writer = new PrintWriter(file))
            {
            lines.stream()
                 .map(s -> s.replace("%serviceName%", sServiceName))
                 .map(s -> s.replace("%clusterName%", sClusterName))
                 .forEach(writer::println);
            }

        // delete the service first to make sure that it does not exist
        s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "delete", "svc", sServiceName));

        // create the service
        int nExitCode = s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "create", "-f", file.getCanonicalPath()));

        assertThat("Creation of headless service failed", nExitCode, is(0));

        m_headlessService = sServiceName;
        }

    private Application portForwardManagement(String sNamespace, String sRelease) throws Exception
        {
        return portForwardCoherencePod(s_k8sCluster, sNamespace, sRelease, 30000);
        }

    private Map query(Application appPortForward, String sPath) throws Exception
        {
        PortMapping       portMapping = appPortForward.get(PortMapping.class);
        int               nPort       = portMapping.getPort().getActualPort();
        String            sSep        = sPath.startsWith("/") ? "" : "/";
        String            sURL        = "http://127.0.0.1:" + nPort + sSep + sPath;
        URI               uri         = URI.create(sURL);
        HttpURLConnection connection  = (HttpURLConnection) uri.toURL().openConnection();

        assertThat(connection.getResponseCode(), is(200));

        return MAPPER.readValue(connection.getInputStream(), LinkedHashMap.class);
        }

    // ----- data members ---------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The name of the deployed Coherence Helm release.
     */
    private final List<String> m_listRelease = new ArrayList<>();

    /**
     * The name of the headless service used by current test.
     */
    private String             m_headlessService;

    /**
     * The mapper to parse json.
     */
    private static final ObjectMapper MAPPER = new ObjectMapper();

    /**
     * JUnit temporary folder rule.
     */
    @ClassRule
    public static TemporaryFolder s_temporaryFolder = new TemporaryFolder();
    }
