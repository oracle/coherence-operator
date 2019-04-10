/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.bedrock.runtime.LocalPlatform;
import com.oracle.common.net.SSLSocketProvider;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.HttpsConfigurator;
import com.sun.net.httpserver.HttpsParameters;
import com.sun.net.httpserver.HttpsServer;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.SocketProviderFactory;
import com.tangosol.run.xml.XmlHelper;
import com.tangosol.util.Resources;
import org.junit.rules.ExternalResource;

import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLParameters;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.core.MediaType;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.function.Function;
import java.util.regex.Pattern;

/**
 * A http server stub to use in testing where the response of
 * the server to different requests can be configured and asserted
 * in test methods.
 *
 * @author jk
 */
public class HttpServerStub
        extends ExternalResource
    {
    // ----- ExternalResource methods ---------------------------------------

    @Override
    protected void before()
        {
        start();
        }

    @Override
    protected void after()
        {
        stop();
        }

    // ----- HttpServerStub methods -----------------------------------------

    public HttpServerStub bindToPort(int nPort)
        {
        m_nPort = nPort;

        return this;
        }

    public int getBoundPort()
        {
        if (m_server == null)
            {
            throw new IllegalStateException("Server is not running");
            }

        return m_nPort;
        }

    public HttpServerStub setSocketProviderXml(String sProvider)
        {
        m_sSocketProvider = sProvider;
        
        return this;
        }

    /**
     * Add a GET response function for a path.
     *
     * @param path   the path to use the function for
     * @param sBody  the String to return as the response body
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onGet(String path, String sBody)
        {
        return onGet(path, (p) -> new HttpResponse(sBody));
        }

    /**
     * Add a GET response function for a path.
     *
     * @param path     the path to use the function for
     * @param urlBody  the URL to return as the response body
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onGet(String path, URL urlBody, MediaType mediaType)
        {
        return onGet(path, urlBody, mediaType.toString());
        }

    /**
     * Add a GET response function for a path.
     *
     * @param path     the path to use the function for
     * @param urlBody  the URL to return as the response body
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onGet(String path, URL urlBody, String mediaType)
        {
        return onGet(path, (p) -> fromURL(urlBody, 200, mediaType));
        }

    /**
     * Add a GET response function for a path.
     *
     * @param path     the path to use the function for
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onGet(String path, Function<String, HttpResponse> response)
        {
        f_mapHttpResponse.computeIfAbsent(GET.class, k -> new ConcurrentHashMap<>()).put(path, response);

        return this;
        }


    /**
     * Add a GET response function for a path that matches the specified regex.
     *
     * @param path     the regex to use to match a request path
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onGet(Pattern path, Function<String, HttpResponse> response)
        {
        f_mapRegExHttpResponse.computeIfAbsent(GET.class, k -> new ConcurrentHashMap<>()).put(path, response);

        return this;
        }


    /**
     * Add a PUT response function for a path.
     *
     * @param path     the path to use the function for
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onPut(String path, Function<String, HttpResponse> response)
        {
        f_mapHttpResponse.computeIfAbsent(PUT.class, k -> new ConcurrentHashMap<>()).put(path, response);

        return this;
        }


    /**
     * Add a PUT response function for a path that matches the specified regex.
     *
     * @param path     the regex to use to match a request path
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onPut(Pattern path, Function<String, HttpResponse> response)
        {
        f_mapRegExHttpResponse.computeIfAbsent(PUT.class, k -> new ConcurrentHashMap<>()).put(path, response);

        return this;
        }


    /**
     * Add a POST response function for a path.
     *
     * @param path     the path to use the function for
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onPost(String path, Function<String, HttpResponse> response)
        {
        f_mapHttpResponse.computeIfAbsent(POST.class, k -> new ConcurrentHashMap<>()).put(path, response);

        return this;
        }


    /**
     * Add a POST response function for a path that matches the specified regex.
     *
     * @param path     the regex to use to match a request path
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onPost(Pattern path, Function<String, HttpResponse> response)
        {
        f_mapRegExHttpResponse.computeIfAbsent(POST.class, k -> new ConcurrentHashMap<>()).put(path, response);

        return this;
        }


    /**
     * Add a DELETE response function for a path.
     *
     * @param path     the path to use the function for
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onDelete(String path, Function<String, HttpResponse> response)
        {
        f_mapHttpResponse.computeIfAbsent(DELETE.class, k -> new ConcurrentHashMap<>()).put(path, response);
        return this;
        }


    /**
     * Add a DELETE response function for a path that matches the specified regex.
     *
     * @param path     the regex to use to match a request path
     * @param response the function to use to build a response
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub onDelete(Pattern path, Function<String, HttpResponse> response)
        {
        f_mapRegExHttpResponse.computeIfAbsent(DELETE.class, k -> new ConcurrentHashMap<>()).put(path, response);
        return this;
        }


    /**
     * Clear all the configured responses.
     *
     * @return this {@link HttpServerStub} to allow for method chaining
     */
    public HttpServerStub reset()
        {
        f_mapHttpResponse.forEach((key, value) -> value.clear());
        f_mapHttpResponse.clear();

        f_mapRegExHttpResponse.forEach((key, value) -> value.clear());
        f_mapRegExHttpResponse.clear();

        return this;
        }


    /**
     * Obtain the {@link URI} used to start this server.
     *
     * @return the {@link URI} used to start this server
     */
    public URI getUri()
        {
        return m_uri;
        }


    /**
     * Obtain the port that the server has bound to.
     *
     * @return the port that the server has bound to
     */
    public int getPort()
        {
        return m_nPort;
        }

    /**
     * Start the http server.
     */
    public void start()
        {
        if (m_server != null)
            {
            throw new IllegalStateException("Server is already running");
            }

        try
            {
            if (m_nPort <= 0)
                {
                m_nPort = LocalPlatform.get().getAvailablePorts().next();
                }

            URI uriBase = new URI("http://" + m_sAddress + ":" + m_nPort);

            if ("0.0.0.0".equals(m_sAddress))
                {
                m_uri = new URI("http://localhost:" + m_nPort);
                }
            else
                {
                m_uri = uriBase;
                }


            if (m_sSocketProvider == null)
                {
                m_server = HttpServer.create(new InetSocketAddress(m_sAddress, m_nPort), 0);
                }
            else
                {
                Cluster                                cluster     = CacheFactory.getCluster();
                SocketProviderFactory                  factory     = cluster.getDependencies().getSocketProviderFactory();
                SSLSocketProvider                      provider    = (SSLSocketProvider) factory.getSocketProvider(m_sSocketProvider);
                SSLSocketProvider.Dependencies         deps        = provider.getDependencies();
                boolean                                fClientAuth = deps.isClientAuthenticationRequired();
                final SSLParameters                    sslParams   = deps.getSSLParameters();
                final SSLContext                       context     = deps.getSSLContext();

                m_server = HttpsServer.create(new InetSocketAddress(m_sAddress, m_nPort), 0);

                ((HttpsServer) m_server).setHttpsConfigurator(new HttpsConfigurator(context)
                    {
                    public void configure(HttpsParameters params)
                        {
                        params.setSSLParameters(sslParams);
                        params.setNeedClientAuth(fClientAuth);
                        }
                    });
                }

            m_server.createContext("/", new RequestHandler());
            m_server.setExecutor(Executors.newSingleThreadExecutor());
            m_server.start();

            System.out.println("Http server stub is listening on " + m_server.getAddress());
            }
        catch (IOException | URISyntaxException e)
            {
            throw new RuntimeException("Failed to start server", e);
            }
        }

    /**
     * Stop the http server.
     */
    public void stop()
        {
        if (m_server != null)
            {
            m_server.stop(0);
            }
        m_server = null;
        }

    /**
     * Process the request.
     *
     * @param action the class representing the http action
     * @param path   the request path
     *
     * @return the {@link HttpResponse} to return to the caller
     */
    public HttpResponse getHttpResponse(Class<?> action, String path)
        {
        try
            {
            Map<String, Function<String, HttpResponse>> map
                    = f_mapHttpResponse.computeIfAbsent(action, k -> new ConcurrentHashMap<>());

            Function<String, HttpResponse> function = map.get(path);

            if (function == null && path.endsWith("/") && path.length() > 1)
                {
                function = map.get(path.substring(0, path.length() - 1));
                }
            else if (function == null && !path.endsWith("/"))
                {
                function = map.get(path + "/");
                }

            if (function != null)
                {
                return function.apply(path);
                }
            else
                {
                Map<Pattern, Function<String, HttpResponse>> mapRegex
                        = f_mapRegExHttpResponse.computeIfAbsent(action, k -> new ConcurrentHashMap<>());

                Function<String, HttpResponse> fn = mapRegex.entrySet().stream()
                        .filter(entry -> entry.getKey().matcher(path).matches())
                        .findFirst()
                        .map(Map.Entry::getValue)
                        .orElse(null);

                if (fn != null)
                    {
                    try
                        {
                        HttpResponse response = fn.apply(path);
                        return response;
                        }
                    catch (Throwable e)
                        {
                        e.printStackTrace();
                        }
                    }
                }

            return f_fnNotFound.apply(path);
            }
        catch (Throwable e)
            {
            e.printStackTrace();
            return SERVER_ERROR;
            }
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     */
    public static void main(String[] args)
        {
        URL    urlKeysStore  = Resources.findFileOrResource("certs/groot.jks", null);
        URL    urlTrustStore = Resources.findFileOrResource("certs/truststore-guardians.jks", null);
        String sslProvider   = "ManagementSSLProvider";

        System.setProperty("coherence.override", "ssl-coherence-override.xml");
        System.setProperty(ProbeHttpClient.PROP_HTTP_SOCKET_PROVIDER, sslProvider);
        System.setProperty("test.keystore", urlKeysStore.toExternalForm());
        System.setProperty("test.truststore", urlTrustStore.toExternalForm());

        HttpServerStub stub = new HttpServerStub()
                .setSocketProviderXml(sslProvider);

        stub.onGet("/foo", "hello foo");
        stub.onGet("/bar", "hello bar");

        stub.start();
        }

    private HttpResponse fromURL(URL url, int nStatus, String mediaType)
        {
        try (InputStream in = url.openStream())
            {
            int    cBytes = in.available();
            byte[] ab     = new byte[cBytes];

            in.read(ab);

            return new HttpResponse(ab, nStatus, mediaType);
            }
        catch (IOException e)
            {
            throw new RuntimeException(e);
            }
        }

    // ----- inner class RequestHandler -------------------------------------

    private class RequestHandler
            implements HttpHandler
        {
        @Override
        public void handle(HttpExchange httpExchange)
            {
            String       sPath    = httpExchange.getRequestURI().toString();
            String       sMethod  = httpExchange.getRequestMethod();
            HttpResponse response;

            switch (sMethod)
                {
                case "GET":
                    response = HttpServerStub.this.getHttpResponse(GET.class, sPath);
                    break;
                default:
                    response = NOT_FOUND;
                }

            if (response != null)
                {
                response.send(httpExchange);
                }
            }
        }

    public static class HttpResponse
        {
        public HttpResponse(int statusCode)
            {
            this(null, statusCode, null);
            }
        
        public HttpResponse(String sBody)
            {
            this(sBody, 200);
            }
        
        public HttpResponse(String sBody, int statusCode)
            {
            this(sBody == null ? null : sBody.getBytes(), statusCode, null);
            }

        public HttpResponse(byte[] abBody, int statusCode, String mediaType)
            {
            m_abBody     = abBody;
            m_statusCode = statusCode;
            m_mediaType  = mediaType;
            }

        public void send(HttpExchange httpExchange)
            {
            try
                {
                byte[] ab;
                if (m_abBody != null && m_abBody.length > 0)
                    {
                    ab = m_abBody;
                    }
                else
                    {
                    ab = EMPTY;
                    }

                if (m_mediaType != null)
                    {
                    httpExchange.getResponseHeaders().add("Content-Type", m_mediaType);
                    }
                else
                    {
                    httpExchange.getResponseHeaders().add("Content-Type", MediaType.TEXT_PLAIN);
                    }
                
                httpExchange.sendResponseHeaders(m_statusCode, ab.length);

                try (OutputStream os = httpExchange.getResponseBody())
                    {
                    os.write(ab);
                    }
                }
            catch (Throwable e)
                {
                e.printStackTrace();
                }
            }

        private final byte[] m_abBody;
        
        private final int m_statusCode;

        private final String m_mediaType;
        }
    
    // ----- constants ------------------------------------------------------

    public static final HttpResponse NOT_FOUND = new HttpResponse(404);
    
    public static final HttpResponse SERVER_ERROR = new HttpResponse(500);

    public static final byte[] EMPTY = new byte[0];

    // ----- data members ---------------------------------------------------

    private final Map<Class, Map<String, Function<String, HttpResponse>>> f_mapHttpResponse
            = new ConcurrentHashMap<>();

    private final Map<Class, Map<Pattern, Function<String, HttpResponse>>> f_mapRegExHttpResponse
            = new ConcurrentHashMap<>();

    private final Function<String, HttpResponse> f_fnNotFound
            = (p) -> NOT_FOUND;

    private String m_sAddress = "0.0.0.0";

    private int m_nPort = -1;

    private URI m_uri;

    private HttpServer m_server;

    private String m_sSocketProvider;
    }
