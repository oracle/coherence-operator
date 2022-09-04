/*
 * Copyright (c) 2019, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import java.io.IOException;
import java.io.OutputStream;
import java.lang.reflect.Method;
import java.net.InetSocketAddress;
import java.util.stream.Collectors;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.Coherence;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.NamedCache;
import com.tangosol.net.partition.SimplePartitionKey;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

/**
 * A simple Http server that is deployed into a Coherence cluster
 * and can be used to perform various tests.
 *
 * @author jk  2019.08.09
 */
public class RestServer {

    /**
     * Private constructor.
     */
    private RestServer() {
    }

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     *
     * @throws java.lang.Exception if the process fails to run
     */
    public static void main(String[] args) throws Exception {
        try {
            int port = Integer.parseInt(System.getProperty("test.rest.port", "8080"));
            HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);

            server.createContext("/ready", RestServer::ready);
            server.createContext("/env", RestServer::env);
            server.createContext("/props", RestServer::props);
            server.createContext("/suspend", RestServer::suspend);
            server.createContext("/resume", RestServer::resume);
            server.createContext("/canaryStart", RestServer::canaryStart);
            server.createContext("/canaryCheck", RestServer::canaryCheck);
            server.createContext("/canaryClear", RestServer::canaryClear);

            server.setExecutor(null); // creates a default executor
            server.start();

            System.out.println("REST server is UP! http://localhost:" + server.getAddress().getPort());
        }
        catch (Throwable thrown) {
            System.err.println("Failed to start http server");
            thrown.printStackTrace();
        }

        String main = System.getProperty("coherence.k8s.testing.main", getMainClass());
        Class<?> cls = Class.forName(main);
        Method method = cls.getMethod("main", String[].class);
        method.invoke(method, (Object) args);
    }

    private static String getMainClass() {
        try {
            return Coherence.class.getCanonicalName();
        }
        catch (Throwable e) {
            return DefaultCacheServer.class.getCanonicalName();
        }
    }

    private static void send(HttpExchange t, int status, String body) throws IOException {
        t.sendResponseHeaders(status, body.length());
        OutputStream os = t.getResponseBody();
        os.write(body.getBytes());
        os.close();
    }

    private static void ready(HttpExchange t) throws IOException {
        send(t, 200, "OK");
    }

    private static void env(HttpExchange t) throws IOException {
        String data = System.getenv()
                .entrySet()
                .stream()
                .map(e -> String.format("{\"%s\":\"%s\"}", e.getKey(), e.getValue()))
                .collect(Collectors.joining(",\n"));

        send(t, 200, "[" + data + "]");
    }

    private static void props(HttpExchange t) throws IOException {
        String data = System.getProperties()
                .entrySet()
                .stream()
                .map(e -> String.format("{\"%s\":\"%s\"}", e.getKey(), e.getValue()))
                .collect(Collectors.joining(",\n"));

        send(t, 200, "[" + data + "]");
    }

    private static void suspend(HttpExchange t) throws IOException {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.suspendService("PartitionedCache");
        send(t, 200, "OK");
    }

    private static void resume(HttpExchange t) throws IOException {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.resumeService("PartitionedCache");
        send(t, 200, "OK");
    }

    @SuppressWarnings("unchecked")
    private static void canaryStart(HttpExchange t) throws IOException {
        NamedCache cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();

        for (int i = 0; i < nPart; i++) {
            SimplePartitionKey key = SimplePartitionKey.getPartitionKey(i);
            cache.put(key, "data");
        }

        send(t, 200, "OK");
    }

    private static void canaryCheck(HttpExchange t) throws IOException {
        NamedCache cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();
        int nSize = cache.size();

        if (nSize == nPart) {
            send(t, 200, "OK " + nSize + " entries");
        }
        else {
            send(t, 400, "Expected " + nPart + " entries but there are only " + nSize);
        }
    }

    private static void canaryClear(HttpExchange t) throws IOException {
        NamedCache cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();

        cache.clear();

        send(t, 200, "OK");
    }
}
