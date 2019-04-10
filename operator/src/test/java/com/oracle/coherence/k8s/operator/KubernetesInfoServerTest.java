/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import com.oracle.bedrock.runtime.LocalPlatform;
import io.kubernetes.client.apis.CoreV1Api;
import io.kubernetes.client.models.V1Node;
import io.kubernetes.client.models.V1ObjectMeta;
import org.junit.Test;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Unit tests for KubernetesInfoServer.
 *
 * @author  sw
 */
public class KubernetesInfoServerTest {
    @Test
    public void testBasic() throws Exception
        {
        int nPort = LocalPlatform.get().getAvailablePorts().next();
        KubernetesInfoServer kubernetesInfoServer = new KubernetesInfoServer(nPort);

        CoreV1Api    api  = mock(CoreV1Api.class);

        V1Node       node = new V1Node();
        V1ObjectMeta meta = new V1ObjectMeta();
        meta.putLabelsItem("failure-domain.beta.kubernetes.io/zone", "myzone");
        node.setMetadata(meta);
        when(api.readNode("mynode", null, Boolean.TRUE, Boolean.TRUE)).thenReturn(node);

        kubernetesInfoServer.setApi(api);
        kubernetesInfoServer.start();

        verifyZone(nPort, "", 404, "");

        verifyZone(nPort, "/mynode", 200, "myzone");

        verifyZone(nPort, "/nosuchnode", 404, "");

        kubernetesInfoServer.stop(1);
        }

    private void verifyZone(int nPort, String sPath, int nExpectedStatus, String sExpectedZone) throws Exception
        {
        HttpURLConnection connection = (HttpURLConnection) new URL("http://localhost:" + nPort + "/zone" + sPath).openConnection();
        connection.setRequestMethod("GET");
        connection.connect();
        int nStatus = connection.getResponseCode();
        assertEquals(nExpectedStatus, nStatus);

        InputStream inputStream = (nStatus >= 400) ? connection.getErrorStream() : connection.getInputStream();
        String sZone = new BufferedReader(new InputStreamReader(inputStream)).readLine();
        if (sExpectedZone.length() == 0)
            {
            assertTrue(sZone == null || sZone.length() == 0);
            }
        else
            {
            assertEquals(sExpectedZone, sZone);
            }
        }
}
