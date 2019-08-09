/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.NamedCache;
import com.tangosol.net.partition.SimplePartitionKey;
import io.helidon.common.http.Http;
import io.helidon.webserver.Routing;
import io.helidon.webserver.ServerConfiguration;
import io.helidon.webserver.ServerRequest;
import io.helidon.webserver.ServerResponse;
import io.helidon.webserver.WebServer;

/**
 * @author jk  2019.08.09
 */
public class RestServer
    {
    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     */
    public static void main(String[] args)
        {
        ServerConfiguration configuration = ServerConfiguration.builder()
                                                               .port(8080)
                                                               .build();
        Routing routing = Routing.builder()
                                 .put("/suspend", RestServer::suspend)
                                 .put("/resume", RestServer::resume)
                                 .put("/canaryStart", RestServer::canaryStart)
                                 .put("/canaryCheck", RestServer::canaryCheck)
                                 .build();

        WebServer.create(configuration, routing)
                .start()
                .thenAccept(s -> {
                    System.out.println("ReST server is UP! http://localhost:" + s.port());
                    s.whenShutdown().thenRun(() -> System.out.println("ReST server is DOWN. Good bye!"));
                })
                .exceptionally(t -> {
                    System.err.println("ReST server startup failed: " + t.getMessage());
                    t.printStackTrace(System.err);
                    return null;
                });

        DefaultCacheServer.main(args);
        }

    static void suspend(ServerRequest req, ServerResponse res)
        {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.suspendService("PartitionedCache");
        res.send("OK");
        }

    static void resume(ServerRequest req, ServerResponse res)
        {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.resumeService("PartitionedCache");
        res.send("OK");
        }

    @SuppressWarnings("unchecked")
    static void canaryStart(ServerRequest req, ServerResponse res)
        {
        NamedCache              cache   = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int                     nPart   = service.getPartitionCount();

        for (int i=0; i<nPart; i++)
            {
            SimplePartitionKey key = SimplePartitionKey.getPartitionKey(i);
            cache.put(key, "data");
            }

        res.send("OK");
        }

    static void canaryCheck(ServerRequest req, ServerResponse res)
        {
        NamedCache              cache   = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int                     nPart   = service.getPartitionCount();
        int                     nSize   = cache.size();

        if (nSize == nPart)
            {
            res.status(Http.Status.OK_200).send("OK " + nSize + " entries");
            }
        else
            {
            res.status(Http.Status.BAD_REQUEST_400).send("Expected " + nPart + " entries but there are only " + nSize);
            }
        }
    }
