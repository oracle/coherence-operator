/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.net.HttpURLConnection;
import java.net.SocketException;
import java.net.URI;
import java.util.Properties;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.SSLSocketFactory;

import com.oracle.bedrock.runtime.LocalPlatform;
import org.junit.Test;

import static com.oracle.coherence.k8s.OperatorRestServer.PROP_HEALTH_LOG;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_HEALTH_PORT;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_INSECURE;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_TLS_KEYSTORE;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_TLS_KEYSTORE_PASSWORD_FILE;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_TLS_KEY_PASSWORD_FILE;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_TLS_TRUSTSTORE;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_TLS_TRUSTSTORE_PASSWORD_FILE;
import static com.oracle.coherence.k8s.OperatorRestServer.PROP_TLS_TWO_WAY;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.Assert.fail;

public class OperatorRestHttpsIT {

    public static final String CERTS = System.getProperty("test.certs.location", "../build/_output/certs");

    private LocalPlatform platform = LocalPlatform.get();

    @Test
    public void shouldBeInsecure() throws Exception {
        Properties properties = getProperties();
        properties.setProperty(PROP_INSECURE, "true");

        try (OperatorRestServer server = new OperatorRestServer(properties)) {
            server.start();
            HttpURLConnection connection = readyRequestHttp(server);
            assertThat(connection.getResponseCode(), is(400));
        }
    }

    @Test
    public void shouldUseTLS() throws Exception {
        Properties properties = getProperties();

        try (OperatorRestServer server = new OperatorRestServer(properties)) {
            server.start();
            HttpURLConnection connection = readyRequestTLS(server);
            assertThat(connection.getResponseCode(), is(400));
        }
    }

    @Test
    public void shouldUseTLSNotAllowingHttpClient() throws Exception {
        Properties properties = getProperties();

        try (OperatorRestServer server = new OperatorRestServer(properties)) {
            server.start();
            HttpURLConnection connection = readyRequestHttp(server);
            try {
                connection.getResponseCode();
                fail("Should not have allowed connection");
            } catch (SocketException e) {
                // expected
            }
        }
    }

    private HttpsURLConnection readyRequestTLS(OperatorRestServer server) throws Exception {
        int port = server.getPort();
        URI uri = URI.create("https://127.0.0.1:" + port + "/ready");
        HttpsURLConnection connection = (HttpsURLConnection) uri.toURL().openConnection();
        connection.setSSLSocketFactory(getFactory());
        connection.setHostnameVerifier(new HostnameVerifier() {
            @Override
            public boolean verify(String s, SSLSession sslSession) {
                return true;
            }
        });
        return connection;
    }

    private HttpURLConnection readyRequestHttp(OperatorRestServer server) throws Exception {
        int port = server.getPort();
        URI uri = URI.create("http://127.0.0.1:" + port + "/ready");
        return (HttpURLConnection) uri.toURL().openConnection();
    }

    private Properties getProperties() {
        Properties properties = new Properties();

        Integer port = platform.getAvailablePorts().next();
        properties.setProperty(PROP_HEALTH_PORT, String.valueOf(port));

        properties.setProperty(PROP_INSECURE, "false");
        properties.setProperty(PROP_TLS_KEYSTORE, "file:" + CERTS + "/icarus.jks");
        properties.setProperty(PROP_TLS_KEYSTORE_PASSWORD_FILE, "file:" + CERTS + "/storepassword.txt");
        properties.setProperty(PROP_TLS_KEY_PASSWORD_FILE, "file:" + CERTS + "/keypassword.txt");
        properties.setProperty(PROP_TLS_TRUSTSTORE, "file:" + CERTS + "/truststore-all.jks");
        properties.setProperty(PROP_TLS_TRUSTSTORE_PASSWORD_FILE, "file:" + CERTS + "/trustpassword.txt");
        properties.setProperty(PROP_TLS_TWO_WAY, "true");
        properties.setProperty(PROP_HEALTH_LOG, "true");
        return properties;
    }

    private Properties getClientProperties() {
        Properties properties = new Properties();

        properties.setProperty(PROP_INSECURE, "false");
        properties.setProperty(PROP_TLS_KEYSTORE, "file:" + CERTS + "/groot.jks");
        properties.setProperty(PROP_TLS_KEYSTORE_PASSWORD_FILE, "file:" + CERTS + "/storepassword.txt");
        properties.setProperty(PROP_TLS_KEY_PASSWORD_FILE, "file:" + CERTS + "/keypassword.txt");
        properties.setProperty(PROP_TLS_TRUSTSTORE, "file:" + CERTS + "/truststore-all.jks");
        properties.setProperty(PROP_TLS_TRUSTSTORE_PASSWORD_FILE, "file:" + CERTS + "/trustpassword.txt");
        properties.setProperty(PROP_TLS_TWO_WAY, "true");
        properties.setProperty(PROP_HEALTH_LOG, "true");
        return properties;
    }

    private SSLSocketFactory getFactory() throws Exception {
        SSLContext context = OperatorRestServer.createSSLContext(getClientProperties());
        return context.getSocketFactory();
    }
}
