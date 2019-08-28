/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;
import com.tangosol.util.Resources;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.ClassRule;
import org.junit.Test;
import util.AssumingCoherenceVersion;

import javax.ws.rs.core.MediaType;
import java.net.URL;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * @author jk
 */
public class RestProbeHttpsIT
    {
    @BeforeClass
    public static void setup()
        {
        URL    urlKeysStore  = Resources.findFileOrResource("certs/groot.jks", null);
        URL    urlTrustStore = Resources.findFileOrResource("certs/truststore-guardians.jks", null);
        URL    urlJson       = Resources.findFileOrResource("json/cluster.json", null);
        String sProvider     = "ManagementSSLProvider";

        System.setProperty("coherence.override", "ssl-coherence-override.xml");
        System.setProperty(ProbeHttpClient.PROP_HTTP_SOCKET_PROVIDER, sProvider);
        System.setProperty("test.keystore", urlKeysStore.toExternalForm());
        System.setProperty("test.truststore", urlTrustStore.toExternalForm());


        s_httpServer = new HttpServerStub()
                .setSocketProviderXml(sProvider)
                .onGet(RestProbe.PATH_CLUSTER, urlJson, MediaType.APPLICATION_JSON);

        s_httpServer.start();

        System.setProperty(ProbeHttpClient.PROP_HTTP_PORT, String.valueOf(s_httpServer.getBoundPort()));
        }

    @AfterClass
    public static void cleanup()
        {
        if (s_httpServer != null)
            {
            s_httpServer.stop();
            }
        }

    
    @Test
    public void shouldConnectWithSSL()
        {
        RestProbe probe = new RestProbe();

        assertThat(probe.isAvailable(), is(true));
        }

    // ----- data members ---------------------------------------------------

    public static HttpServerStub s_httpServer;

    /**
     * The full Coherence image name to use for tests.
     */
    public static final String COHERENCE_IMAGE = CacheFactory.VERSION;

    /**
     * FileBasedPasswordProvider is only available from Coherence 12.2.1.4.0 and greater.
     */
    @ClassRule
    public static AssumingCoherenceVersion assumingCoherenceVersion = new AssumingCoherenceVersion(COHERENCE_IMAGE, "12.2.1.4.0");
    }
