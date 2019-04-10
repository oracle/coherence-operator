/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.concurrent.ThreadFactory;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

import com.squareup.okhttp.Call;
import com.squareup.okhttp.MediaType;
import com.squareup.okhttp.Protocol;
import com.squareup.okhttp.Request;
import com.squareup.okhttp.Response;
import com.squareup.okhttp.ResponseBody;
import io.kubernetes.client.apis.CoreV1Api;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.assertThat;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.ArgumentMatchers.isNull;
import static org.mockito.Mockito.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

/**
 * Unit tests for NamespaceWatcher and CoherenceOperator.NamespaceProcessor.
 *
 * @author sc
 */
public class NamespaceWatcherTest {
    private static final ThreadFactory defaultThreadFactory = Executors.defaultThreadFactory();

    /**
     * Test watching Namspaces.
     *
     * @throws Exception
     */
    @Test
    public void testNamespace() throws Exception
        {
        String[] asNamespaces = new String[]{ "kube-system", "kube-public", "default", "docker" };

        AtomicBoolean    fStopping      = new AtomicBoolean(false);
        List<String>     results        = new ArrayList<>();
        CountDownLatch   countDownLatch = new CountDownLatch(4);
        NamespaceWatcher watcher        = new NamespaceWatcher(fStopping,
                item -> {
                    results.add(item.type + "-" + item.object.getMetadata().getName());
                    countDownLatch.countDown();
                });

        List<String> expectedResult = new ArrayList<>();
        for (String ns : asNamespaces)
            {
            expectedResult.add("ADDED-" + ns);
            }

        CoreV1Api coreV1Api = createMockCoreV1Api(asNamespaces);

        watcher.setApi(coreV1Api);

        watcher.start(defaultThreadFactory);
        countDownLatch.await(DEFAULT_TIMEOUT_SECONDS, TimeUnit.SECONDS);

        assertThat("Results: " + results, results.equals(expectedResult), is(true));
        }

    /**
     * Test NamespaceProcessor with non-empty included namespaces.
     *
     * @throws Exception
     */
    @Test
    public void testConfigMap() throws Exception
        {
        String[] asNamespaces         = new String[]{ "kube-system", "kube-public", "default", "docker", "cohns", "cohns2" };
        String[] asIncludedNamespaces = new String[] { "cohns", "cohns2" };
        String[] asExcludedNamespaces = new String[] {"kube-system", "kube-public", "docker"};

        CoreV1Api coreV1Api = setupTestConfigMap(asNamespaces, asIncludedNamespaces, asExcludedNamespaces);

        verify(coreV1Api, times(1)).createNamespacedConfigMapAsync(eq("cohns"), any(),any(), any());
        verify(coreV1Api, times(1)).createNamespacedConfigMapAsync(eq("cohns2"), any(),any(), any());
        verify(coreV1Api, times(2)).createNamespacedConfigMapAsync(any(), any(),any(), any());
        }

    /**
     * Test NamespaceProcessor without included namespaces.
     *
     * @throws Exception
     */
    @Test
    public void testConfigMapNoIncluded() throws Exception
        {
        String[] asNamespaces         = new String[]{ "kube-system", "kube-public", "default", "docker", "cohns" };
        String[] asIncludedNamespaces = new String[] { };
        String[] asExcludedNamespaces = new String[] {"kube-system", "kube-public", "docker"};

        CoreV1Api coreV1Api = setupTestConfigMap(asNamespaces, asIncludedNamespaces, asExcludedNamespaces);

        verify(coreV1Api, times(1)).createNamespacedConfigMapAsync(eq("cohns"), any(),any(), any());
        verify(coreV1Api, times(1)).createNamespacedConfigMapAsync(eq("default"), any(),any(), any());
        verify(coreV1Api, times(2)).createNamespacedConfigMapAsync(any(), any(),any(), any());
        }

    /**
     * Test NamespaceProcessor with excluded namespaces.
     *
     * @throws Exception
     */
    @Test
    public void testConfigMapExcluded() throws Exception
        {
        String[] asNamespaces         = new String[]{ "kube-system", "kube-public", "default", "docker", "cohns", "internal" };
        String[] asIncludedNamespaces = new String[] { "cohns", "cohns2" };
        String[] asExcludedNamespaces = new String[] {"kube-system", "kube-public", "docker", "internal" };

        CoreV1Api coreV1Api = setupTestConfigMap(asNamespaces, asIncludedNamespaces, asExcludedNamespaces);

        verify(coreV1Api, times(1)).createNamespacedConfigMapAsync(eq("cohns"), any(),any(), any());
        verify(coreV1Api, times(1)).createNamespacedConfigMapAsync(any(), any(),any(), any());
        }

    // ---- helper methods --------------------------------------------------

    /**
     * Set up the CoreV1Api for NamespaceProcessor with given parameters for further testing.
     *
     * @param asNamespaces          an array of namespaces
     * @param asIncludedNamespaces  an array of namespaces included for processing
     * @param asExcludedNamespaces  an array of namespaces excluded for processing
     * @return CoreV1Api for verification
     * @throws Exception
     */
    CoreV1Api setupTestConfigMap(String[] asNamespaces, String[] asIncludedNamespaces, String[] asExcludedNamespaces) throws Exception
        {
        AtomicBoolean  fStopping      = new AtomicBoolean(false);
        CountDownLatch countDownLatch = new CountDownLatch(asNamespaces.length);
        CoreV1Api      coreV1Api      = createMockCoreV1Api(asNamespaces);

        CoherenceOperator.NamespaceProcessor namespaceProcessor =
                new CoherenceOperator.NamespaceProcessor("cohns", asIncludedNamespaces, asExcludedNamespaces);
        namespaceProcessor.setCoreV1Api(coreV1Api);

        NamespaceWatcher watcher = new NamespaceWatcher(fStopping,
                item -> {
                    namespaceProcessor.accept(item);
                    countDownLatch.countDown();
                });

        watcher.setApi(coreV1Api);

        watcher.start(defaultThreadFactory);
        countDownLatch.await(DEFAULT_TIMEOUT_SECONDS, TimeUnit.SECONDS);

        return coreV1Api;
        }

    /**
     * Create a mock CoreV1Api object that support querying namespace.
     *
     * @param asNamespaces  array of namespaces
     * @return
     * @throws Exception
     */
    private CoreV1Api createMockCoreV1Api(String[] asNamespaces) throws Exception
        { 
        CoreV1Api coreV1Api = mock(CoreV1Api.class);
        Call      call      = createMockCall(asNamespaces);

        when(coreV1Api.listNamespaceCall(any(), isNull(), any(), any(), any(), any(),
                any(), any(), eq(Boolean.TRUE), any(), any())).thenReturn(call);

        return coreV1Api;
        }

    /**
     * Create a mock call object which corresponding to cDSize Deployment objects.
     *
     * @param asNamespaces  array of namespaces
     * @return
     * @throws Exception
     */
    private Call createMockCall(String[] asNamespaces) throws Exception
        {
        int    cVersion         = 0;
        String namespaceFormat = "{ \"kind\": \"Namespace\", \"apiVersion\": \"v1\"," +
                "  \"metadata\": { \"selfLink\": \"/api/v1/namespaces/%1$s\", " +
                "\"resourceVersion\": \"%2$d\", \"name\": \"%1$s\" } }";
        String responseFormat  = "{ \"type\": \"ADDED\", \"object\": " + namespaceFormat + "}";

        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < asNamespaces.length; i++)
            {
            if (i > 0)
                {
                sb.append("\n");
                }
            sb.append(String.format(responseFormat, asNamespaces[i], cVersion++));
            }

        Call         call         = mock(Call.class);
        ResponseBody responseBody = ResponseBody.create(MediaType.parse("application/json"), sb.toString());
        Request      request      = new Request.Builder().url(new URL("http://localhost:8080/api/v1/namespaces?watch=true")).build();
        Response     response     = new Response.Builder()
                .protocol(Protocol.HTTP_1_1).request(request)
                .code(200).body(responseBody)
                .build();

        when(call.execute()).thenReturn(response);
        return call;
        }

    // ---- constants -------------------------------------------------------

    /**
     * The default timeout to wait watcher processing.
     */
    private static final int DEFAULT_TIMEOUT_SECONDS = 15;
    }
