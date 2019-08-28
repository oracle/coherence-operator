/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.common.net.SSLSettings;
import com.oracle.common.net.SocketProvider;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.SocketProviderFactory;
import org.glassfish.jersey.client.ClientConfig;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.WebTarget;
import javax.ws.rs.core.Response;
import java.io.Closeable;
import java.net.URI;

/**
 * A simple utility class that handles correctly configuring a http or https client.
 * <p>
 * It is expected the the probe JVM is started with the same parameters as the co-located
 * Coherence server so that the http client can be configured with the same settings,
 * in particular SSL certificates.
 *
 * @author jk
 */
public class ProbeHttpClient
        implements Closeable
    {
    // ----- constructors ---------------------------------------------------

    public ProbeHttpClient()
        {
        this(System.getProperty(PROP_HTTP_PORT, String.valueOf(DEFAULT_MANAGEMENT_PORT)),
             System.getProperty(PROP_HTTP_SOCKET_PROVIDER));
        }

    /**
     * Create a {@link ProbeHttpClient}.
     */
    public ProbeHttpClient(int nPort)
        {
        this(String.valueOf(nPort), null);
        }

    /**
     * Create a {@link ProbeHttpClient}.
     */
    public ProbeHttpClient(int nPort, String sSocketProvider)
        {
        this(String.valueOf(nPort), sSocketProvider);
        }

    /**
     * Create a {@link ProbeHttpClient}.
     */
    public ProbeHttpClient(String sPort, String sSocketProvider)
        {
        ClientConfig  clientConfig = new ClientConfig();
        ClientBuilder builder      = ClientBuilder.newBuilder()
                                                  .withConfig(clientConfig)
                                                  .hostnameVerifier(new NoopHostnameVerifier());

        try
            {
            m_nPort = Integer.parseInt(sPort);
            }
        catch (NumberFormatException e)
            {
            m_nPort = DEFAULT_MANAGEMENT_PORT;
            }

        if (sSocketProvider != null && sSocketProvider.length() > 0)
            {
            Cluster               cluster  = CacheFactory.getCluster();
            SocketProviderFactory factory  = cluster.getDependencies().getSocketProviderFactory();
            SocketProvider        provider = factory.getSocketProvider(sSocketProvider);
            SSLSettings           settings = factory.getSSLSettings(provider);

            if (settings != null)
                {
                SSLContext context = settings.getSSLContext();

                m_sProtocol = "https";
                builder     = builder.sslContext(context);
                }
            }

        m_client = builder.build();
        }

    // ----- ProbeHttpClient methods ----------------------------------------

    /**
     * Obtain a {@link WebTarget} for the specified request path.
     *
     * @param sPath  the request path.
     *
     * @return  the {@link WebTarget} created from the path
     */
    public WebTarget getWebTarget(String sPath)
        {
        return m_client.target(getResourceURL(sPath));
        }

    /**
     * Perform a plain http get request for the specified path.
     *
     * @param sPath  the request path
     *
     * @return  the {@link Response} from invoking the request
     */
    public Response get(String sPath)
        {
        return getWebTarget(sPath).request().get();
        }

    /**
     * Obtain a request {@link URI} for the specified path.
     *
     * @param sPath  the request path
     *
     * @return  a {@link URI} made up of this client's configured protocol,
     *          host and port with the request path appended
     */
    public URI getResourceURL(String sPath)
        {
        String sSep = sPath.charAt(0) == '/' ? "" : "/";
        String sURL = String.format("%s://%s:%d%s%s", m_sProtocol, HOST_NAME, m_nPort, sSep, sPath);

        return URI.create(sURL);
        }

    // ----- Closeable interface --------------------------------------------

    @Override
    public void close()
        {
        m_client.close();
        }

    // ----- inner class NoopHostnameVerifier -------------------------------

    /**
     * A no-op implementation of a {@link HostnameVerifier}.
     */
    private class NoopHostnameVerifier
            implements HostnameVerifier
        {
        @Override
        public boolean verify(String s, SSLSession sslSession)
            {
            return true;
            }
        }

    // ----- data members ---------------------------------------------------

    /**
     * The System property used in Coherence to obtain the optional socket
     * provider to use for the management server.
     */
    public static final String PROP_HTTP_SOCKET_PROVIDER = "coherence.management.http.provider";

    /**
     * The System property used by Coherence to set the port that the management server will bind to.
     */
    public static final String PROP_HTTP_PORT = "coherence.management.http.port";

    /**
     * The default value of the management port used by Coherence.
     */
    public static final int DEFAULT_MANAGEMENT_PORT = 30000;

    /**
     * The host name to use for http requests.
     */
    public static final String HOST_NAME = "127.0.0.1";

    /**
     * The protocol to use for http requests. This will be http unless the server is
     * using SSL in which case it will be https.
     */
    private String m_sProtocol = "http";

    /**
     * The port to use for http requests.
     */
    private int m_nPort;

    /**
     * The {@link Client} to use for requests.
     */
    private Client m_client;
    }
