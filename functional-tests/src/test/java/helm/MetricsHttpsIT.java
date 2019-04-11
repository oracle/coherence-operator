/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.runtime.Application;
import org.junit.*;
import util.AssumingCoherenceVersion;

import javax.net.ssl.SSLException;
import javax.ws.rs.client.Client;
import javax.ws.rs.client.WebTarget;
import javax.ws.rs.core.Response;

import static org.hamcrest.CoreMatchers.instanceOf;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * @author jk  2019.02.21
 */
@Ignore
public class MetricsHttpsIT
        extends BaseHttpsTest
    {
    // ----- test lifecycle -------------------------------------------------

    @BeforeClass
    public static void setup() throws Exception
        {
        BaseHttpsTest.setup("values/helm-values-ssl-metrics.yaml", COHERENCE_VERSION, 9095);
        }

    @AfterClass
    public static void cleanupDeployment()
        {
        BaseHttpsTest.cleanupDeployment(MetricsHttpsIT.class);
        }

    // ----- test methods ---------------------------------------------------

    @Test
    public void shouldGetWithTrustedClient()
        {
        WebTarget webTarget = s_clientHelper.getHttpsWebTarget(s_clientStarLord, URL_METRICS);
        Response  response  = webTarget.request().get();
        int       status    = response.getStatus();

        assertThat(status, is(Response.Status.OK.getStatusCode()));
        }

    @Test
    public void shouldNotGetIfServerDoesNotTrustClient()
        {
        WebTarget webTarget = s_clientHelper.getHttpsWebTarget(s_clientYondu, URL_METRICS);

        // should get an error trying due to the client's CA cert not being in the servers's trust store
        m_exception.expect(javax.ws.rs.ProcessingException.class);
        m_exception.expectCause(instanceOf(SSLException.class));

        webTarget.request().get();
        }

    @Test
    public void shouldNotGetIfClientDoesNotTrustServer()
        {
        WebTarget webTarget = s_clientHelper.getHttpsWebTarget(s_clientNoServerTrust, URL_METRICS);

        // should get an error trying due to the server's CA cert not being in the client's trust store
        m_exception.expect(javax.ws.rs.ProcessingException.class);
        m_exception.expectCause(instanceOf(SSLException.class));

        webTarget.request().get();
        }

    @Test
    public void shouldConnectToManagementEndpointAsPlainHttp() throws Exception
        {
        String sNamespace = getK8sNamespace();

        try (Application app = portForwardCoherencePod(s_k8sCluster, sNamespace, m_sRelease, 30000))
            {
            PortMapping    portMapping  = app.get(PortMapping.class);
            int            nPort        = portMapping.getPort().getActualPort();
            HttpTestHelper clientHelper = new HttpTestHelper(nPort);
            Client         client       = clientHelper.getClient();
            WebTarget      webTarget    = clientHelper.getHttpWebTarget(client, URL_MANAGEMENT);
            Response       response     = webTarget.request().get();
            int            status       = response.getStatus();

            assertThat(status, is(Response.Status.OK.getStatusCode()));
            }
        }

// Disabled as actually authenticating the client cert's username is not supported
//    @Test
//    public void shouldNotGetSecuredResourceWithUnauthorisedClient()
//        {
//        WebTarget webTarget = s_clientHelper.getHttpsWebTarget(s_clientGroot, URL_METRICS);
//        Response  response  = webTarget.request().get();
//        int       status    = response.getStatus();
//
//        assertThat(status, is(Response.Status.UNAUTHORIZED.getStatusCode()));
//        }

    /**
     * Prometheus metrics are only available from Coherence 12.2.1.4.0 and greater.
     */
    @ClassRule
    public static AssumingCoherenceVersion assumingCoherenceVersion = new AssumingCoherenceVersion(COHERENCE_VERSION, "12.2.1.4.0");
    }
