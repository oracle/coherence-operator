/*
 * Copyright (c) 2019, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import java.util.Collections;

import javax.enterprise.context.ApplicationScoped;
import javax.json.Json;
import javax.json.JsonArrayBuilder;
import javax.json.JsonBuilderFactory;
import javax.json.JsonObject;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.NamedCache;
import com.tangosol.net.partition.SimplePartitionKey;

/**
 * A simple JAX-RS service that is deployed into a Coherence cluster
 * and can be used to perform various tests.
 *
 * @author jk  2021.03.12
 */
@Path("/")
@ApplicationScoped
public class RestServer {

    private static final JsonBuilderFactory JSON = Json.createBuilderFactory(Collections.emptyMap());

    @GET
    @Path("ready")
    @Produces(MediaType.TEXT_PLAIN)
    public String ready() {
        return "ok";
    }

    @GET
    @Path("env")
    @Produces(MediaType.APPLICATION_JSON)
    public JsonObject env() {
        JsonArrayBuilder arrayBuilder = JSON.createArrayBuilder();
        System.getenv()
                .entrySet()
                .stream()
                .map(e -> String.format("{\"%s\":\"%s\"}", e.getKey(), e.getValue()))
                .forEach(arrayBuilder::add);

        return JSON.createObjectBuilder()
                .add("env", arrayBuilder)
                .build();
    }

    @GET
    @Path("props")
    @Produces(MediaType.APPLICATION_JSON)
    public JsonObject props()  {
        JsonArrayBuilder arrayBuilder = JSON.createArrayBuilder();
        System.getProperties()
                .entrySet()
                .stream()
                .map(e -> String.format("{\"%s\":\"%s\"}", e.getKey(), e.getValue()))
                .forEach(arrayBuilder::add);

        return JSON.createObjectBuilder()
                .add("env", arrayBuilder)
                .build();
    }


    @PUT
    @Path("suspend")
    @Produces(MediaType.TEXT_PLAIN)
    public String suspend() {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.suspendService("PartitionedCache");
        return "ok";
    }

    @PUT
    @Path("resume")
    @Produces(MediaType.TEXT_PLAIN)
    public String resume() {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.resumeService("PartitionedCache");
        return "ok";
    }

    @PUT
    @Path("canaryStart")
    @Produces(MediaType.TEXT_PLAIN)
    public String canaryStart() {
        NamedCache<SimplePartitionKey, String> cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();

        for (int i = 0; i < nPart; i++) {
            SimplePartitionKey key = SimplePartitionKey.getPartitionKey(i);
            cache.put(key, "data");
        }
        return "ok";
    }

    @GET
    @Path("canaryCheck")
    @Produces(MediaType.APPLICATION_JSON)
    public JsonObject canaryCheck() {
        NamedCache<?, ?> cache = CacheFactory.getCache("canary");
        DistributedCacheService service = (DistributedCacheService) cache.getCacheService();
        int nPart = service.getPartitionCount();
        int nSize = cache.size();

        if (nSize == nPart) {
            return JSON.createObjectBuilder()
                    .add("entries", nSize)
                    .build();
        }
        else {
            throw new IllegalStateException("Data loss " + nSize + " of " + nPart + " partitions");
        }
    }

    @POST
    @Path("canaryClear")
    @Produces(MediaType.TEXT_PLAIN)
    public String canaryClear() {
        NamedCache<?, ?> cache = CacheFactory.getCache("canary");
        cache.clear();
        return "ok";
    }
}
