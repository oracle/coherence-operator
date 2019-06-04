/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.runtime.options.Arguments;
import com.tangosol.util.Resources;
import org.junit.Rule;
import org.junit.rules.ExpectedException;

import javax.ws.rs.client.Client;
import java.io.File;
import java.net.URL;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * @author jk  2019.02.21
 */
public abstract class BaseHttpsTest
        extends BaseHelmChartTest
    {
    // ----- helper methods -------------------------------------------------

    /**
     * Setup the test using the specified values file.
     *
     * @param sValuesFile        the values file to use when installing the Coherence chart
     * @param sCoherenceVersion  the Coherence Docker image tag (version)
     * @param nPortToForward     the container port to forward
     *
     * @throws Exception  if there is an error setting up the test
     */
    protected static void setup(String sValuesFile, String sCoherenceVersion, int nPortToForward) throws Exception
        {
        // create the SSL k8s secret (delete it first to ensure the secret is correct)
        String sSecret         = "ssl-secret";
        String sNamespace      = getK8sNamespace();
        String sCoherenceImage = COHERENCE_IMAGE_PREFIX + sCoherenceVersion;
        String sImageName      = "coherenceImage.name=" + sCoherenceImage;

        createSecret(sSecret, sNamespace);

        m_sRelease = installCoherence(s_k8sCluster, sNamespace, sValuesFile, sImageName);

        assertCoherence(s_k8sCluster, sNamespace, m_sRelease);

        // port forward to the management port on the Pod
        s_appPortForward = portForwardCoherencePod(s_k8sCluster, sNamespace, m_sRelease, nPortToForward);

        PortMapping portMapping = s_appPortForward.get(PortMapping.class);
        int         nPort       = portMapping.getPort().getActualPort();

        createClients(nPort);
        }

    /**
     * Clean up the test.
     *
     * @param clsTest  the test class
     */
    protected static void cleanupDeployment(Class clsTest)
        {
        if (m_sRelease != null)
            {
            deleteCoherence(clsTest, s_k8sCluster, getK8sNamespace(), m_sRelease, false);
            }
        }

    /**
     * Create the k8s SSL secret.
     *
     * @param sSecret     the name of the secret
     * @param sNamespace  the namespace to create the secret in
     *
     * @throws Exception  if there is an error creating the secret
     */
    protected static void createSecret(String sSecret, String sNamespace) throws Exception
        {
        URL  urlKeyStore    = Resources.findFileOrResource(HttpTestHelper.CERT_ICARUS, null);
        File fileKeyStore   = new File(urlKeyStore.toURI());
        URL  urlTrustStore  = Resources.findFileOrResource(HttpTestHelper.TRUSTSTORE_GUARDIANS, null);
        File fileTrustStore = new File(urlTrustStore.toURI());

        s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "delete", "secret", sSecret));

        int nExitCode = s_k8sCluster.kubectlAndWait(Arguments.of("-n", sNamespace, "create", "secret", "generic", sSecret,
                                                                 "--from-file", fileKeyStore.getCanonicalPath(),
                                                                 "--from-file", fileTrustStore.getCanonicalPath(),
                                                                 "--from-literal", "keypassword.txt=" + HttpTestHelper.KEY_PASSWORD,
                                                                 "--from-literal", "storepassword.txt=" + HttpTestHelper.STORE_PASSWORD,
                                                                 "--from-literal", "trustpassword.txt=" + HttpTestHelper.TRUST_PASSWORD));

        assertThat("Failed to create SSL secret", nExitCode, is(0));
        }

    /**
     * Create the http clients to use in the tests.
     *
     * @param nPort  the port that the http server is listening on
     */
    protected static void createClients(int nPort)
        {
        s_clientHelper        = new HttpTestHelper(nPort);
        s_clientStarLord      = s_clientHelper.getClient(HttpTestHelper.CERT_STAR_LORD,
                                                         HttpTestHelper.STORE_PASSWORD,
                                                         HttpTestHelper.KEY_PASSWORD,
                                                         HttpTestHelper.TRUSTSTORE_GUARDIANS,
                                                         HttpTestHelper.TRUST_PASSWORD);
        s_clientYondu         = s_clientHelper.getClient(HttpTestHelper.CERT_YONDU,
                                                         HttpTestHelper.STORE_PASSWORD,
                                                         HttpTestHelper.KEY_PASSWORD,
                                                         HttpTestHelper.TRUSTSTORE_GUARDIANS,
                                                         HttpTestHelper.TRUST_PASSWORD);
        s_clientNoServerTrust = s_clientHelper.getClient(HttpTestHelper.CERT_STAR_LORD,
                                                         HttpTestHelper.STORE_PASSWORD,
                                                         HttpTestHelper.KEY_PASSWORD,
                                                         HttpTestHelper.TRUSTSTORE_RAVAGERS,
                                                         HttpTestHelper.TRUST_PASSWORD);
        }

    // ----- data members ---------------------------------------------------

    /**
     * The Management over ReST endpoint.
     */
    protected static final String URL_MANAGEMENT = "/management/coherence/cluster";
    
    /**
     * The metrics endpoint.
     */
    protected static final String URL_METRICS = "/metrics";

    /**
     * A JAX-RS client with valid certificates.
     */
    protected static Client s_clientStarLord;

    /**
     * A JAX-RS client with invalid certificates.
     */
    protected static Client s_clientYondu;

    /**
     * The trust store that does not trust the server.
     */
    protected static Client s_clientNoServerTrust;

    /**
     * A JUnit rule to assert exceptions.
     */
    @Rule
    public final ExpectedException m_exception = ExpectedException.none();

    /**
     * The k8s cluster to use to install the charts.
     */
    protected static K8sCluster s_k8sCluster = getDefaultCluster();

    /**
     * The name of the deployed Coherence Helm release.
     */
    protected static String m_sRelease;

    /**
     * The kubectl port-forward process.
     */
    protected static Application s_appPortForward;

    /**
     * The {@link HttpTestHelper} utility class.
     */
    protected static HttpTestHelper s_clientHelper;

    /**
     * The Docker registry name to use to pull Coherence images.
     */
    public static final String DOCKER_REGISTRY = System.getProperty("docker.repo");

    /**
     * The version (tag) for the latest Coherence image version.
     */
    public static final String COHERENCE_VERSION = System.getProperty("coherence.docker.version");

    /**
     * The base Coherence image name without a tag.
     */
    public static final String COHERENCE_IMAGE_PREFIX = DOCKER_REGISTRY + "coherence:";
    }
