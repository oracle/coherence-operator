/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.bedrock.runtime.LocalPlatform;
import com.oracle.bedrock.runtime.java.JavaApplication;
import com.oracle.bedrock.runtime.java.options.ClassName;
import com.oracle.bedrock.runtime.java.options.SystemProperty;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.oracle.bedrock.util.Capture;
import org.junit.Test;

import java.net.HttpURLConnection;
import java.net.URI;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static org.hamcrest.CoreMatchers.is;

/**
 * @author jk
 */
public class HealthIT
    {
    @Test
    public void shouldBeReadySingleMember() throws Exception
        {
        LocalPlatform    platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort)))
            {
            Eventually.assertThat(invoking(this).httpRequest(httpPort, HealthServer.PATH_READY), is(200));
            }
        }

    @Test
    public void shouldBeReadyMultipleMember() throws Exception
        {
        LocalPlatform    platform  = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1)))
            {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                       ClassName.of(Main.class),
                                                       SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2)))
                {
                Eventually.assertThat(invoking(this).httpRequest(httpPort1, HealthServer.PATH_READY), is(200));
                Eventually.assertThat(invoking(this).httpRequest(httpPort2, HealthServer.PATH_READY), is(200));
                }
            }
        }

    @Test
    public void shouldBeLiveSingleMember() throws Exception
        {
        LocalPlatform    platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort)))
            {
            Eventually.assertThat(invoking(this).httpRequest(httpPort, HealthServer.PATH_HEALTH), is(200));
            }
        }

    @Test
    public void shouldBeLiveMultipleMember() throws Exception
        {
        LocalPlatform    platform  = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1)))
            {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                       ClassName.of(Main.class),
                                                       SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2)))
                {
                Eventually.assertThat(invoking(this).httpRequest(httpPort1, HealthServer.PATH_HEALTH), is(200));
                Eventually.assertThat(invoking(this).httpRequest(httpPort2, HealthServer.PATH_HEALTH), is(200));
                }
            }
        }

    @Test
    public void shouldBeStatusHASingleMember() throws Exception
        {
        LocalPlatform    platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort)))
            {
            Eventually.assertThat(invoking(this).httpRequest(httpPort, HealthServer.PATH_HA), is(200));
            }
        }

    @Test
    public void shouldBeStatusHAMultipleMember() throws Exception
        {
        LocalPlatform    platform  = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1)))
            {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                       ClassName.of(Main.class),
                                                       SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2)))
                {
                Eventually.assertThat(invoking(this).httpRequest(httpPort1, HealthServer.PATH_HA), is(200));
                Eventually.assertThat(invoking(this).httpRequest(httpPort2, HealthServer.PATH_HA), is(200));
                }
            }
        }

    // ----- helper methods -------------------------------------------------

    // Must be public - used in Eventually.assertThat
    public int httpRequest(Capture<Integer> httpPort, String path) throws Exception
        {
        URI uri = URI.create("http://127.0.0.1:" + httpPort.get() + path);
        HttpURLConnection connection = (HttpURLConnection) uri.toURL().openConnection();
        connection.setRequestMethod("GET");
        connection.connect();
        return connection.getResponseCode();
        }
    }
