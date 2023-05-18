/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import java.io.IOException;
import java.io.OutputStream;
import java.lang.reflect.Method;
import java.net.InetSocketAddress;
import java.util.HashMap;
import java.util.Map;
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
     * @throws Exception if the server fails to start
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
            server.createContext("/shutdown", RestServer::shutdown);

            server.setExecutor(null); // creates a default executor
            server.start();

            System.out.println("REST server is UP! http://localhost:" + server.getAddress().getPort());
        }
        catch (Throwable thrown) {
            System.err.println("Failed to start http server");
            thrown.printStackTrace();
        }

        Class<?> clsMain = Class.forName(getMainClass());
        Method method = clsMain.getMethod("main", args.getClass());
        method.invoke(null, (Object) args);
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

    @SuppressWarnings({"unchecked", "rawtypes"})
    private static void canaryStart(HttpExchange t) throws IOException {
        NamedCache cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();

        for (int i = 0; i < nPart; i++) {
            SimplePartitionKey key = SimplePartitionKey.getPartitionKey(i);
            cache.put(key, "data");
        }

        System.out.println("CanaryCheck: Loaded canary cache " + cache.getCacheName() + " with " + cache.size() + " entries");
        send(t, 200, "OK");
    }

    @SuppressWarnings({"rawtypes"})
    private static void canaryCheck(HttpExchange t) throws IOException {
        NamedCache cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();
        int nSize = cache.size();

        if (nSize == nPart) {
            System.out.println("CanaryCheck: Passed found " + nSize + " partitions in " + cache.getCacheName());
            send(t, 200, "OK " + nSize + " entries");
        }
        else {
            System.out.println("CanaryCheck: Failed found " + nPart + " of " + nSize + " partitions in " + cache.getCacheName());
            send(t, 400, "Expected " + nPart + " entries but there are only " + nSize);
        }
    }

    @SuppressWarnings({"rawtypes"})
    private static void canaryClear(HttpExchange t) throws IOException {
        NamedCache cache = CacheFactory.getCache("canary");
        System.out.println("CanaryCheck: Clearing canary cache " + cache.getCacheName());
        cache.clear();
        send(t, 200, "OK");
    }

    private static void shutdown(HttpExchange t) throws IOException {
        send(t, 202, "OK");
        Map<String, String> map = queryToMap(t);
        int exitCode = 0;
        String s = map.get("exitCode");
        if (s != null && !s.isEmpty()) {
            try {
                exitCode = Integer.parseInt(s);
            }
            catch (Exception e) {
                // ignored
            }
        }
        System.exit(exitCode);
    }

    private static String getMainClass() {
        try {
            return Coherence.class.getCanonicalName();
        }
        catch (Throwable e) {
            return DefaultCacheServer.class.getCanonicalName();
        }
    }

    public static Map<String, String> queryToMap(HttpExchange exchange) {
        String query = exchange.getRequestURI().getQuery();
        Map<String, String> result = new HashMap<>();
        for (String param : query.split("&")) {
            String[] pair = param.split("=");
            if (pair.length > 1) {
                result.put(pair[0], pair[1]);
            }
            else {
                result.put(pair[0], "");
            }
        }
        return result;
    }

}
