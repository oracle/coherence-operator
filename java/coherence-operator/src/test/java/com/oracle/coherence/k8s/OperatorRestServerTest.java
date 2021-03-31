/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.util.Collections;
import java.util.Set;
import java.util.function.Supplier;

import com.tangosol.net.Cluster;
import com.tangosol.net.Member;

import com.sun.net.httpserver.HttpExchange;
import org.junit.Before;
import org.junit.Test;

import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

public class OperatorRestServerTest {

    private Cluster cluster;
    private Supplier<Cluster> clusterSupplier;
    private HttpExchange exchange;

    private Runnable waitForDCS = () -> {};

    @Before
    public void setup() {
        Set<Member> memberSet = Collections.singleton(mock(Member.class));
        cluster = mock(Cluster.class);

        when(cluster.isRunning()).thenReturn(true);
        when(cluster.getMemberSet()).thenReturn(memberSet);
        clusterSupplier = () -> cluster;

        exchange = mock(HttpExchange.class);
        OutputStream out = new ByteArrayOutputStream();
        when(exchange.getResponseBody()).thenReturn(out);
    }

    @Test
    public void shouldHandleReadyCheckError() throws IOException {
        when(cluster.getManagement()).thenThrow(new RuntimeException("Computer says No!"));

        OperatorRestServer server = new OperatorRestServer(clusterSupplier, waitForDCS);

        server.ready(exchange);

        verify(exchange).sendResponseHeaders(500, 0);
    }

    @Test
    public void shouldHandleReadyCheckWhenNoManagedNodes() throws IOException {
        when(cluster.getManagement()).thenThrow(new RuntimeException(OperatorRestServer.NO_MANAGED_NODES));

        OperatorRestServer server = new OperatorRestServer(clusterSupplier, waitForDCS);

        server.ready(exchange);

        verify(exchange).sendResponseHeaders(400, 0);
    }

    @Test
    public void shouldHandleStatusCheckError() throws IOException {
        when(cluster.getManagement()).thenThrow(new RuntimeException("Computer says No!"));

        OperatorRestServer server = new OperatorRestServer(clusterSupplier, waitForDCS);

        server.statusHA(exchange);

        verify(exchange).sendResponseHeaders(500, 0);
    }

    @Test
    public void shouldHandleStatusCheckWhenNoManagedNodes() throws IOException {
        when(cluster.getManagement()).thenThrow(new RuntimeException(OperatorRestServer.NO_MANAGED_NODES));

        OperatorRestServer server = new OperatorRestServer(clusterSupplier, waitForDCS);

        server.statusHA(exchange);

        verify(exchange).sendResponseHeaders(400, 0);
    }

    @Test
    public void shouldHandleStatusCheckAndWaitForDCS() throws IOException {
        Runnable dcs = mock(Runnable.class);

        when(cluster.getManagement()).thenThrow(new RuntimeException(OperatorRestServer.NO_MANAGED_NODES));

        OperatorRestServer server = new OperatorRestServer(clusterSupplier, dcs);

        server.statusHA(exchange);

        verify(exchange).sendResponseHeaders(400, 0);
        verify(dcs).run();
    }

    @Test
    public void shouldHandleStatusCheckAndWaitForDCSFailure() throws IOException {
        Runnable dcs = mock(Runnable.class);

        doThrow(new IllegalStateException("Oops!")).when(dcs).run();

        OperatorRestServer server = new OperatorRestServer(clusterSupplier, dcs);

        server.statusHA(exchange);

        verify(exchange).sendResponseHeaders(400, 0);
        verify(dcs).run();
    }
}
