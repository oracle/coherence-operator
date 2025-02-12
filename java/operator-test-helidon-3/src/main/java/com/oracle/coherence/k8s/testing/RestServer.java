/*
 * Copyright (c) 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing;

import java.util.Collections;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.NamedCache;
import com.tangosol.net.partition.SimplePartitionKey;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.json.Json;
import jakarta.json.JsonArrayBuilder;
import jakarta.json.JsonBuilderFactory;
import jakarta.json.JsonObject;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.PUT;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

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

    /**
     * @return always returns {@code "ok"}
     */
    @GET
    @Path("ready")
    @Produces(MediaType.TEXT_PLAIN)
    public String ready() {
        return "ok";
    }

    /**
     * Returns the JVM environment variables.
     *
     * @return the JVM environment variables
     */
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

    /**
     * Returns the JVM system properties.
     *
     * @return the JVM system properties
     */
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

    /**
     * Suspend the canary cache service.
     *
     * @return always returns {@code "ok"}
     */
    @PUT
    @Path("suspend")
    @Produces(MediaType.TEXT_PLAIN)
    public String suspend() {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.suspendService("PartitionedCache");
        return "ok";
    }

    /**
     * Resume the canary cache service.
     *
     * @return always returns {@code "ok"}
     */
    @PUT
    @Path("resume")
    @Produces(MediaType.TEXT_PLAIN)
    public String resume() {
        Cluster cluster = CacheFactory.ensureCluster();
        cluster.resumeService("PartitionedCache");
        return "ok";
    }

    /**
     * Initialise the canary cache.
     *
     * @return always returns {@code "ok"}
     */
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

    /**
     * Check the canary cache.
     *
     * @return the number of entries in the canary cache
     */
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

    /**
     * Clear the canary cache.
     *
     * @return always returns {@code "ok"}
     */
    @POST
    @Path("canaryClear")
    @Produces(MediaType.TEXT_PLAIN)
    public String canaryClear() {
        NamedCache<?, ?> cache = CacheFactory.getCache("canary");
        cache.clear();
        return "ok";
    }
}
