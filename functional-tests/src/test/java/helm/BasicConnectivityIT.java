/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.tangosol.net.ConfigurableCacheFactory;
import com.tangosol.net.ExtensibleConfigurableCacheFactory;
import org.junit.After;
import org.junit.Ignore;
import org.junit.Test;

import java.net.HttpURLConnection;
import java.net.URI;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.junit.Assume.assumeThat;

/**
 * Test basic connectivity to the various ports exposed by the Helm chart.
 *
 * @author jk  2019.02.11
 */
public class BasicConnectivityIT
        extends BaseHelmChartTest
    {
    @After
    public void cleanUpCoherence()
        {
        if (m_sRelease != null)
            {
            deleteCoherence(m_k8sCluster, m_sNamespace, m_sRelease, false);
            }
        }

    @Test
    public void shouldConnectToManagementPort() throws Exception
        {
        assumeThat(versionCheck("12.2.1.4.0"), is(true));

        m_sRelease = installCoherence(m_k8sCluster, m_sNamespace, DEFAULT_VALUES_YAML);

        assertCoherence(m_k8sCluster, m_sNamespace, m_sRelease);

        try (Application application = portForwardManagement())
            {
            PortMapping portMapping = application.get(PortMapping.class);
            int         nPort       = portMapping.getPort().getActualPort();
            URI         uri         = URI.create("http://127.0.0.1:" + nPort + "/management/coherence/cluster");

            Eventually.assertThat(invoking(this).canConnect(uri), is(true));
            }
        }

    @Test
    @Ignore
    public void shouldConnectToMetricsPort() throws Exception
        {
        assumeThat(versionCheck("12.2.1.4.0"), is(true));

        m_sRelease = installCoherence(m_k8sCluster, m_sNamespace, DEFAULT_VALUES_YAML);

        assertCoherence(m_k8sCluster, m_sNamespace, m_sRelease);

        try (Application application = portForwardMetrics())
            {
            PortMapping portMapping = application.get(PortMapping.class);
            int         nPort       = portMapping.getPort().getActualPort();
            URI         uri         = URI.create("http://127.0.0.1:" + nPort + "/metrics");

            Eventually.assertThat(invoking(this).canConnect(uri), is(true));
            }
        }

    @Test
    @Ignore
    public void shouldConnectToExtendPort() throws Exception
        {
        m_sRelease = installCoherence(m_k8sCluster, m_sNamespace, DEFAULT_VALUES_YAML);

        assertCoherence(m_k8sCluster, m_sNamespace, m_sRelease);

        try (Application application = portForwardExtend())
            {
            PortMapping              portMapping = application.get(PortMapping.class);
            int                      nPort       = portMapping.getPort().getActualPort();
            ConfigurableCacheFactory ccf         = getCacheFactory(nPort);


            Eventually.assertThat(invoking(ccf).ensureCache("foo", null), is(notNullValue()));
            }
        }

    @Test
    @Ignore
    public void shouldConnectToJmxPort() throws Exception
        {
        m_sRelease = installCoherence(m_k8sCluster,
                                      m_sNamespace,
                                      DEFAULT_VALUES_YAML,
                                      "store.jmx.enabled=true");

        assertCoherence(m_k8sCluster, m_sNamespace, m_sRelease);

        // assert that the JMX nodes have been deployed in their own k8s Deployment
        assertCoherenceJMX(m_k8sCluster, m_sNamespace, m_sRelease);

        // Custer size should be at least two - one Coherence storage member and one JMX member. Cluster probe could be a third member.
        Eventually.assertThat(invoking(this).getClusterSizeViaJMX(), greaterThanOrEqualTo(2));

        // get local member Id via JMX
        int    nNode      = jmxQuery(m_k8sCluster, m_sNamespace, m_sRelease, "Coherence:type=Cluster", "LocalMemberId");
        String sNodeMBean = "Coherence:type=Node,nodeId=" + nNode;

        // invoke the "reportNodeState" MBean operation
        String sState = jmxInvoke(m_k8sCluster, m_sNamespace, m_sRelease,
                                  sNodeMBean, "reportNodeState",
                                  NO_JMX_PARAMS, EMPTY_JMX_SIGNATURE);

        assertThat(sState, is(notNullValue()));
        }

    // ----- helper methods -------------------------------------------------

    public Integer getClusterSizeViaJMX() throws Exception
        {
        String sClusterMBean = "Coherence:type=Cluster";
        String sAttribute = "ClusterSize";
        return jmxQuery(m_k8sCluster, m_sNamespace, m_sRelease, sClusterMBean, sAttribute);
        }

    public boolean canConnect(URI uri) throws Exception
        {
        HttpURLConnection connection = (HttpURLConnection) uri.toURL().openConnection();

        return connection.getResponseCode() == 200;
        }

    private Application portForwardExtend() throws Exception
        {
        return portForwardCoherencePod(m_k8sCluster, m_sNamespace, m_sRelease, 20000);
        }

    private Application portForwardManagement() throws Exception
        {
        return portForwardCoherencePod(m_k8sCluster, m_sNamespace, m_sRelease, 30000);
        }

    private Application portForwardMetrics() throws Exception
        {
        return portForwardCoherencePod(m_k8sCluster, m_sNamespace, m_sRelease, 9095);
        }

    private ConfigurableCacheFactory getCacheFactory(int nPort)
        {
        System.setProperty("coherence.serializer", "java");
        System.setProperty("coherence.extend.address", "127.0.0.1");
        System.setProperty("coherence.extend.port", String.valueOf(nPort));

        ExtensibleConfigurableCacheFactory.Dependencies deps
                = ExtensibleConfigurableCacheFactory.DependenciesHelper.newInstance("test-client-cache-config.xml");

        return new ExtensibleConfigurableCacheFactory(deps);
        }

    // ----- constants ------------------------------------------------------

    public static final String DEFAULT_VALUES_YAML = "values/helm-values-coh.yaml";

    // ----- data members ---------------------------------------------------

    private K8sCluster m_k8sCluster = getDefaultCluster();

    private String m_sNamespace = getK8sNamespace();

    private String m_sRelease;
    }
