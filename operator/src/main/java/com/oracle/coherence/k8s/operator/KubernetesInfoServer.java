/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpServer;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.apis.CoreV1Api;
import io.kubernetes.client.models.V1Node;
import io.kubernetes.client.models.V1ObjectMeta;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.Map;
import java.util.logging.Logger;

/**
 * A Simple Http server returns info from Kubernetes.
 * Only zone information is returned in this moment.
 *
 * @author sc
 */
public class KubernetesInfoServer
    {
    // ----- constructors ----------------------------------------------------

    /**
     * Constructs a Kubernetes Info Server.
     *
     * @param port  the http port to listen the request.
     *
     * @exception IOException  if there is an issue to create a HttpServer bounding to given port
     */
    public KubernetesInfoServer(int port) throws IOException
        {
                    m_api        = new CoreV1Api();
                    m_httpServer = HttpServer.create(new InetSocketAddress(port), 0);
        HttpContext context      = m_httpServer.createContext("/zone");

        context.setHandler((httpExchange) -> {
            try
                {
                String sRequestURI = httpExchange.getRequestURI().toString();
                int    nStatusCode = 200;
                String sPayload    = null;

                if (sRequestURI.length() >= 6)
                    {
                    String sNodeName = sRequestURI.substring(6); // trim prefix "/zone/"
                        try
                            {
                            sPayload = getZone(sNodeName);
                            }
                        catch(Throwable throwable)
                            {
                            nStatusCode = 500;
                            LOGGER.warning("Exception in getting zone[" + sNodeName + "]: " + throwable);
                            }
                    }

                if (sPayload == null)
                    {
                    sPayload = "";
                    }
                if (sPayload.length() == 0 && nStatusCode == 200)
                    {
                    nStatusCode = 404;
                    }

                httpExchange.sendResponseHeaders(nStatusCode, sPayload.getBytes().length);
                OutputStream output = httpExchange.getResponseBody();
                output.write(sPayload.getBytes());
                output.flush();
                }
            finally
                {
                httpExchange.close();
                }
            });
        }

    // ---- methods ----------------------------------------------------------

    /**
     * Start the server.
     */
    public void start()
        {
        m_httpServer.start();
        }

    /**
     * Stop the server.
     *
     * @param delay  the maximum time in seconds to wait until exchanges have finished.
     */
    void stop(int delay)
        {
        m_httpServer.stop(delay);
        }

    /**
     * Sets api object.
     *
     * @param api  the api object
     */
    void setApi(CoreV1Api api)
    {
        this.m_api = api;
    }

    /**
     * Retrieve zone for the given node name.
     *
     * @param sNodeName
     * @return
     * @throws ApiException
     */
    private String getZone(String sNodeName) throws ApiException
        {
        String zone = "";
        V1Node node = m_api.readNode(sNodeName, null, Boolean.TRUE, Boolean.TRUE);
        if (node != null)
            {
            V1ObjectMeta meta = node.getMetadata();
            if (meta != null)
                {
                Map<String, String> labels = meta.getLabels();
                if (labels != null)
                    {
                    zone = labels.get("failure-domain.beta.kubernetes.io/zone");
                    }
                }
            }
        return zone;
        }

    // ----- data members ---------------------------------------------------
    /**
     * Class Logger.
     */
    static final Logger LOGGER = Logger.getLogger("Operator");

    /**
     * The Core V1 api object to access Kubernetes info.
     */
    private CoreV1Api m_api;

    /**
     * The HttpServer to serve the Http request.
     */
    private HttpServer m_httpServer;
    }
