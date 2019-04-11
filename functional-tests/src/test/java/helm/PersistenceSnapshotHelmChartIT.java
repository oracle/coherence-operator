/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.deferred.options.InitialDelay;
import com.oracle.bedrock.deferred.options.MaximumRetryDelay;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.util.AssertionException;
import org.junit.After;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Ignore;
import org.junit.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.TimeUnit;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static helm.HelmUtils.getK8sObject;
import static helm.HelmUtils.getPods;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.hasItem;
import static org.hamcrest.CoreMatchers.not;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Test persistence and snapshot aspects of the Helm chart values file.
 * <p>
 * This test depends on certain artifacts produced by the Maven build
 * so although this test can be run from an IDE it requires that at least
 * a Maven build with at least the package phase being run first.
 *
 * @author sc
 */
public class PersistenceSnapshotHelmChartIT
        extends BaseHelmChartTest {

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
        }

    @After
    public void cleanUpCoherence()
        {
        if (m_sRelease != null)
            {
            String sJmxSelector       = getCoherenceJmxPodSelector(m_sRelease);
            String sCoherenceSelector = getCoherencePodSelector(m_sRelease);

            //capturePodLogs(PersistenceSnapshotHelmChartIT.class, s_k8sCluster, sJmxSelector);
            //capturePodLogs(PersistenceSnapshotHelmChartIT.class, s_k8sCluster, sCoherenceSelector, COHERENCE_CONTAINER_NAME);
            deleteCoherence(s_k8sCluster, getK8sNamespace(), m_sRelease, m_fStatefulSet);
            }
        }

    // ----- test methods ---------------------------------------------------

    /**
     * Test active persistence: persistence enabled, snapshot enabled.
     *
     * @throws Exception if the test fails
     */
    @Test
    public void testActivePersistence() throws Exception
        {
        testPersistence("values/helm-values-coh-active.yaml", true, true);

        String       sNamespace         = getK8sNamespace();
        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPvcs           = getPodPvcs(sNamespace, sCoherenceSelector);

        assertThat(listPvcs.size(), is(2));

        Set<String> setExpectedPvcs = new HashSet<>();
        setExpectedPvcs.add("persistence-volume");
        setExpectedPvcs.add("snapshot-volume");

        assertThat(new HashSet<>(listPvcs), is (setExpectedPvcs));

        List<String> listEmptyDirs = getPodEmptyDirs(sNamespace, sCoherenceSelector);

        assertThat(listEmptyDirs, not(hasItem("persistence-volume")));
        assertThat(listEmptyDirs, not(hasItem("snapshot-volume")));
        }

    /**
     * Test active persistence with volume specified: persistence enabled, snapshot disabled.
     *
     * @throws Exception if the test fails
     */

    @Test
    public void testActivePersistenceVolume() throws Exception
        {
        testPersistence("values/helm-values-coh-active-vol.yaml", true, false);

        String       sNamespace         = getK8sNamespace();
        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPvcs           = getPodPvcs(sNamespace, sCoherenceSelector);

        assertThat(listPvcs.size(), is(0));

        List<String> listEmptyDirs = getPodEmptyDirs(sNamespace, sCoherenceSelector);

        assertThat(listEmptyDirs, hasItem("persistence-volume"));
        assertThat(listEmptyDirs, not(hasItem("snapshot-volume")));
        }

    /**
     * Test snapshot without persistence but with pvc: persistence disabled, snapshot enabled.
     *
     * @throws Exception if the test fails
     */
    @Test
    public void testSnapshotWithoutPersistenceWithPvc() throws Exception
        {
        testPersistence("values/helm-values-coh-snapshot.yaml", false, true);

        String       sNamespace         = getK8sNamespace();
        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPvcs           = getPodPvcs(sNamespace, sCoherenceSelector);

        assertThat(listPvcs.size(), is(1));

        List<String> listEmptyDirs = getPodEmptyDirs(sNamespace, sCoherenceSelector);

        assertThat(listEmptyDirs, not(hasItem("persistence-volume")));
        assertThat(listEmptyDirs, not(hasItem("snapshot-volume")));
        }

    /**
     * Test snapshot without persistence, pvc: persistence disabled, snapshot disabled.
     *
     * @throws Exception if the test fails
     */
    @Test
    public void testSnapshotWithoutPersistenceWithoutPvc() throws Exception
        {
        testPersistence("values/helm-values-coh.yaml", false, false);

        String       sNamespace         = getK8sNamespace();
        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPvcs           = getPodPvcs(sNamespace, sCoherenceSelector);

        assertThat(listPvcs.size(), is(0));

        List<String> listEmptyDirs = getPodEmptyDirs(sNamespace, sCoherenceSelector);

        assertThat(listEmptyDirs, not(hasItem("persistence-volume")));
        assertThat(listEmptyDirs, not(hasItem("snapshot-volume")));
        }

    /**
     * Test snapshot without persistence but with volume: persistence disabled, snapshot disabled.
     *
     * @throws Exception if the test fails
     */
    @Test
    @Ignore
    public void testSnapshotWithoutPersistenceWithVol() throws Exception
        {
        testPersistence("values/helm-values-coh-snapshot-vol.yaml", false, false);

        String       sNamespace         = getK8sNamespace();
        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPvcs           = getPodPvcs(sNamespace, sCoherenceSelector);

        assertThat(listPvcs.size(), is(0));

        List<String> listEmptyDirs = getPodEmptyDirs(sNamespace, sCoherenceSelector);

        assertThat(listEmptyDirs, not(hasItem("persistence-volume")));
        assertThat(listEmptyDirs, hasItem("snapshot-volume"));
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Assert that Coherence is up. If active persistence is on, then check log file.
     * And test snapshots.
     *
     * @param sHelmValues   the Helm Chart values yaml file
     * @param fPersistence  whether the persistence is enabled in the Helm Chart
     * @param fSnapshot     whether the snapshot is enabled in the Helm Chart
     *
     * @throws Exception if the test fails
     */
    private void testPersistence(String sHelmValues, boolean fPersistence, boolean fSnapshot) throws Exception
        {
        String sNamespace     = getK8sNamespace();

               m_fStatefulSet = fPersistence || fSnapshot;
               m_sRelease     = installCoherence(s_k8sCluster, sNamespace, sHelmValues, "store.jmx.enabled=true"
               );

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);
        assertCoherenceService(sNamespace, m_sRelease);
        assertCoherenceJMX(s_k8sCluster, sNamespace, m_sRelease);

        String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
        List<String> listPods           = getPods(s_k8sCluster, sNamespace, sCoherenceSelector);

        assertThat(listPods.size(), is(1));
        String sCoherencePod = listPods.get(0);

        Eventually.assertThat(invoking(this).hasDefaultCacheServerStarted(s_k8sCluster, sNamespace, sCoherencePod), is(true),
                              Timeout.after(60, TimeUnit.SECONDS), InitialDelay.of(3, TimeUnit.SECONDS));

        if (fPersistence)
            {
            dumpPodLog(s_k8sCluster, sNamespace, listPods.get(0));
            Queue<String> sLogs = getPodLog(s_k8sCluster, sNamespace, listPods.get(0));

            assertThat(sLogs.stream().anyMatch(l -> l.contains("Creating persistence active directory")), is(true));
            }

        testSnapshot(sNamespace);
        }

    private void testSnapshot(String sNamespace)
        throws Exception
        {
        assertPersistenceManagerIdle(sNamespace, m_sRelease);

        // create a snapshot called test-snapshot for service PartitionedCache
        List<String> listCreateSnapshot = jmxInvokePersistenceManager(sNamespace, m_sRelease, "createSnapshot", "test-snapshot");
        if (!listCreateSnapshot.stream().anyMatch(l -> l.contains("test-snapshot")))
            {
            String       sCoherenceSelector = getCoherencePodSelector(m_sRelease);
            List<String> listPods           = getPods(s_k8sCluster, sNamespace, sCoherenceSelector);

            System.out.println("CreateSnapshot failed to create a test-snapshot, dump coherence log for failures");
            dumpPodLog(s_k8sCluster, sNamespace, listPods.get(0));
            Eventually.assertThat(invoking(this).getPersistenceManagerSnapshotsViaJMX(sNamespace, m_sRelease), is(true), Timeout.after(60, TimeUnit.SECONDS));
            }

        // recover a snapshot named test-snapshot
        List<String> listRecoverSnapshot = jmxInvokePersistenceManager(sNamespace, m_sRelease, "recoverSnapshot", "test-snapshot");
        assertThat(listRecoverSnapshot.stream().anyMatch(l -> l.contains("test-snapshot")), is(true));

        // Remove snapshot
        List<String> listRemoveSnapshot  = jmxInvokePersistenceManager(sNamespace, m_sRelease, "removeSnapshot", "test-snapshot");
        assertThat(listRemoveSnapshot.size(), is(0));

        assertPersistenceManagerIdle(sNamespace, m_sRelease);
        }

    public Boolean getPersistenceManagerIdleViaJMX(String sNamespace, String sRelease) throws Exception
        {
        return jmxQuery(s_k8sCluster, sNamespace, sRelease, PERSISTENCE_MANAGER_MBEAN_FOR_SVC_PARTIONED_CACHE, "Idle");
        }

    public List<String> getPersistenceManagerSnapshotsViaJMX(String sNamespace, String sRelease) throws Exception
        {
        return new ArrayList<String>(Arrays.asList(jmxQuery(s_k8sCluster, sNamespace, sRelease, PERSISTENCE_MANAGER_MBEAN_FOR_SVC_PARTIONED_CACHE, "Snapshots")));
        }

    /**
     * Jmx invoke on PersistenceManager MBean that returns snapshots available after command completes.
     * The invocation is made to be synchronous by waiting on idle before returning.
     *
     * @param sNamespace    find jmx server in this namespace
     * @param sRelease      find jmx server for this release
     * @param sOperation    operations are createSnapshot, recoverSnapshot and removeSnapshot
     * @param sSnapshotName operation parameter that is snapshot name
     *
     * @return List of snapshot names after the operation completes
     *
     * @throws Exception
     */
    public List<String> jmxInvokePersistenceManager(String sNamespace, String sRelease, String sOperation, String sSnapshotName)
        throws Exception
        {
        Object[] aoParameters = new Object[]{sSnapshotName};;
        String[] aaSignatures = new String[]{"java.lang.String"};

        for (int i=0; i < 3; i++)
            {
            try
                {
                jmxInvoke(s_k8sCluster, sNamespace, sRelease, PERSISTENCE_MANAGER_MBEAN_FOR_SVC_PARTIONED_CACHE,
                    sOperation, aoParameters, aaSignatures);
                break;
                }
            catch (AssertionException ae)
                {
                // ignore failed port forward assertion and retry
                }
            }

        assertPersistenceManagerIdle(sNamespace, m_sRelease);

        for (int i=0; i < 3; i++)
            {
            try
                {
                return getPersistenceManagerSnapshotsViaJMX(sNamespace, sRelease);
                }
            catch (AssertionException ae)
                {
                // ignore failed port forward assertion and retry
                }
            }
        return null;
        }

    private void assertPersistenceManagerIdle(String sNamespace, String sRelease)
        throws Exception
        {
        Eventually.assertThat(invoking(this).getPersistenceManagerIdleViaJMX(sNamespace, sRelease), is(true),
                InitialDelay.of(3, TimeUnit.SECONDS),
                MaximumRetryDelay.of(RETRY_FREQUENCEY_SECONDS, TimeUnit.SECONDS),
                RetryFrequency.every(RETRY_FREQUENCEY_SECONDS, TimeUnit.SECONDS));
        }

    private void assertCoherenceService(String sNamespace, String sRelease)
        {
        String       sCoherenceServiceSelector = getCoherenceServiceSelector(sRelease);
        List<String> listSvcs                  = getK8sObject(s_k8sCluster, "svc", sNamespace, sCoherenceServiceSelector);
        assertThat(listSvcs.size(), is(1));
        assertThat(listSvcs.get(0), is(sRelease + "-coherence"));
        }

    private List<String> getPodPvcs(String sNamespace, String sCoherenceSelector)
        {
        return getK8sObject(s_k8sCluster, "pods", sNamespace, sCoherenceSelector,
                "{.items[*].spec.volumes[?(@.persistentVolumeClaim.claimName)].name}");
        }

    private List<String> getPodEmptyDirs(String sNamespace, String sCoherenceSelector)
        {
        return getK8sObject(s_k8sCluster, "pods", sNamespace, sCoherenceSelector,
                "{.items[*].spec.volumes[?(@.emptyDir)].name}");
        }

    // ----- constants ------------------------------------------------------

    /**
     * PersistenceManager MBean ObjectName for a PartitionedCache service.
     */
    private static final String PERSISTENCE_MANAGER_MBEAN_FOR_SVC_PARTIONED_CACHE =
        "Coherence:type=Persistence,service=PartitionedCache,responsibility=PersistenceCoordinator";

    // ----- data members ---------------------------------------------------

    /**
     * The k8s cluster to use to install the charts.
     */
    private static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The retry frequency in seconds.
     */
    private static final int RETRY_FREQUENCEY_SECONDS = 5;

    /**
     * The name of the deployed Coherence Helm release.
     */
    private String m_sRelease;

    /**
     * Indicate whether it is a StatefulSet.
     */
    private boolean m_fStatefulSet;
}
