/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.internal.net.ssl.LegacyXmlSSLSocketProviderDependencies;
import com.tangosol.run.xml.SimpleElement;
import com.tangosol.run.xml.XmlElement;
import com.tangosol.util.Resources;
import org.glassfish.jersey.client.ClientConfig;
import org.glassfish.jersey.client.authentication.HttpAuthenticationFeature;
import org.glassfish.jersey.jackson.JacksonFeature;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLSession;

import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.WebTarget;

import java.net.URI;
import java.net.URL;

/**
 * @author jk  2018.03.30
 */
public class HttpSSLHelper
    {
    // ----- constructors ---------------------------------------------------

    public HttpSSLHelper(int nPort)
        {
        m_nPort = nPort;
        }

    // ----- HttpClientHelper methods ---------------------------------------

    /**
     * Return the HTTP(S) client.
     *
     * @return context path
     */
    public Client getClient()
        {
        return createClientBuilder()
                .build();
        }
    
    /**
     * Return the HTTP(S) client.
     *
     * @param sIdentityStore       the location of the identity keystore
     * @param sIdentityPassword    the optional password for the identity keystore
     * @param sIdentityPrivatePwd  the optional private password for the identity keystore
     * @param sTrustStore          the location of the trust store keystore
     * @param sTrustPassword       the optional password for the trust store keystore
     *
     * @return context path
     */
    public Client getClient(String sIdentityStore,
                            String sIdentityPassword,
                            String sIdentityPrivatePwd,
                            String sTrustStore,
                            String sTrustPassword)
        {
        return createClientBuilder(sIdentityStore,
                                   sIdentityPassword,
                                   sIdentityPrivatePwd,
                                   sTrustStore,
                                   sTrustPassword
        )
                .build();
        }

    /**
     * Return the HTTP(S) client.
     *
     * @param sIdentityStore       the location of the identity keystore
     * @param sIdentityPassword    the optional password for the identity keystore
     * @param sIdentityPrivatePwd  the optional private password for the identity keystore
     * @param sTrustStore          the location of the trust store keystore
     * @param sTrustPassword       the optional password for the trust store keystore
     * @param sUser         the basic auth username
     * @param sCredentials  the basic auth credentials to use
     *
     * @return context path
     */
    public Client getClientWithBasicAuth(String sIdentityStore,
                                         String sIdentityPassword,
                                         String sIdentityPrivatePwd,
                                         String sTrustStore,
                                         String sTrustPassword,
                                         String sUser,
                                         String sCredentials)
        {
        return createClientBuilder(sIdentityStore,
                                   sIdentityPassword,
                                   sIdentityPrivatePwd,
                                   sTrustStore,
                                   sTrustPassword
        )
                .build()
                .register(HttpAuthenticationFeature.basic(sUser, sCredentials));
        }

    /**
     * Create a new {@link ClientBuilder}.
     *
     * @return a new {@link ClientBuilder}.
     */
    public ClientBuilder createClientBuilder()
        {
        return createClientBuilder(null,
                                   null,
                                   null,
                                   null,
                                   null);
        }

    /**
     * Create a new {@link ClientBuilder}.
     *
     * @param sIdentityStore       the location of the identity keystore
     * @param sIdentityPassword    the optional password for the identity keystore
     * @param sIdentityPrivatePwd  the optional private password for the identity keystore
     * @param sTrustStore          the location of the trust store keystore
     * @param sTrustPassword       the optional password for the trust store keystore
     *
     * @return a new {@link ClientBuilder}.
     */
    public ClientBuilder createClientBuilder(String sIdentityStore,
                                             String sIdentityPassword,
                                             String sIdentityPrivatePwd,
                                             String sTrustStore,
                                             String sTrustPassword)
        {
        try
            {
            ClientConfig  clientConfig = new ClientConfig();
            boolean       fHasIdentity = sIdentityStore != null && sIdentityStore.length() > 0;
            boolean       fHasTrust    = sTrustStore    != null && sTrustStore.length() > 0;
            ClientBuilder builder      = ClientBuilder.newBuilder()
                                                      .withConfig(clientConfig)
                                                      .register(JacksonFeature.class);
            if (fHasIdentity || fHasTrust)
                {
                XmlElement xml = new SimpleElement("ssl");

                if (fHasIdentity)
                    {
                    URL        url         = Resources.findFileOrResource(sIdentityStore, null);
                    XmlElement xmlIdentity = xml.addElement("identity-manager");
                    XmlElement xmlKeyStore = xmlIdentity.addElement("key-store");

                    xmlKeyStore.addElement("url").setString(url.toExternalForm());

                    if (sIdentityPrivatePwd != null && sIdentityPrivatePwd.length() > 0)
                        {
                        xmlIdentity.addElement("password").setString(sIdentityPrivatePwd);
                        }

                    if (sIdentityPassword != null && sIdentityPassword.length() > 0)
                        {
                        xmlKeyStore.addElement("password").setString(sIdentityPassword);
                        }
                    }

                if (fHasTrust)
                    {
                    URL        url         = Resources.findFileOrResource(sTrustStore, null);
                    XmlElement xmlIdentity = xml.addElement("trust-manager");
                    XmlElement xmlKeyStore = xmlIdentity.addElement("key-store");

                    xmlKeyStore.addElement("url").setString(url.toExternalForm());

                    if (sTrustPassword != null && sTrustPassword.length() > 0)
                        {
                        xmlKeyStore.addElement("password").setString(sTrustPassword);
                        }
                    }

                LegacyXmlSSLSocketProviderDependencies deps = new LegacyXmlSSLSocketProviderDependencies(xml);

                builder.hostnameVerifier(new NoopHostnameVerifier())
                       .sslContext(deps.getSSLContext());
                }

            return builder;
            }
        catch (Exception e)
            {
            throw new RuntimeException("Error creating client", e);
            }
        }


    /**
     * Obtain the {@link WebTarget} for the specified endpoint.
     *
     * @param client     the http client
     * @param sEndpoint  the target endpoint
     *
     * @return  the {@link WebTarget} for the specified endpoint
     */
    public WebTarget getHttpWebTarget(Client client, String sEndpoint)
        {
        return getWebTarget(client, "http", sEndpoint);
        }

    /**
     * Obtain the {@link WebTarget} for the specified endpoint.
     *
     * @param client     the http client
     * @param sEndpoint  the target endpoint
     *
     * @return  the {@link WebTarget} for the specified endpoint
     */
    public WebTarget getHttpsWebTarget(Client client, String sEndpoint)
        {
        return getWebTarget(client, "https", sEndpoint);
        }

    /**
     * Obtain the {@link WebTarget} for the specified endpoint.
     *
     * @param client     the http client
     * @param sProtocol  the protocol, http or https
     * @param sEndpoint  the target endpoint
     *
     * @return  the {@link WebTarget} for the specified endpoint
     */
    public WebTarget getWebTarget(Client client, String sProtocol, String sEndpoint)
        {
        return client.target(getResourceUrl(sProtocol, sEndpoint));
        }

    /**
     * Return the url of the specified resource
     *
     * @param sResource  test resource
     *
     * @return the resource url
     */
    public URI getResourceUrl(String sProtocol, String sResource)
        {
        String sURL     = String.format("%s://%s:%d%s",
                                       sProtocol,
                                       m_sHostName,
                                       m_nPort,
                                       sResource);

        return URI.create(sURL);
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

    // ----- constants ------------------------------------------------------

    /**
     * The localhost address.
     */
    public static String LOCALHOST = "localhost";

    public static final String STORE_PASSWORD = "password";

    public static final String KEY_PASSWORD = "password";

    public static final String TRUST_PASSWORD = "secret";

    public static final String CERT_STAR_LORD = "certs/star-lord.jks";

    public static final String TRUSTSTORE_GUARDIANS = "certs/truststore-guardians.jks";

    // ----- data members ---------------------------------------------------

    private final String m_sHostName = LOCALHOST;

    private final int m_nPort;
    }
