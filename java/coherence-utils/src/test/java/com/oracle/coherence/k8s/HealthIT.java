/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URI;

import com.oracle.bedrock.runtime.LocalPlatform;
import com.oracle.bedrock.runtime.coherence.callables.GetAutoStartServiceNames;
import com.oracle.bedrock.runtime.coherence.callables.IsServiceRunning;
import com.oracle.bedrock.runtime.coherence.options.CacheConfig;
import com.oracle.bedrock.runtime.java.JavaApplication;
import com.oracle.bedrock.runtime.java.options.ClassName;
import com.oracle.bedrock.runtime.java.options.SystemProperty;
import com.oracle.bedrock.runtime.options.StabilityPredicate;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.oracle.bedrock.util.Capture;
import org.junit.Test;

import static com.oracle.bedrock.deferred.DeferredHelper.invoking;
import static org.hamcrest.CoreMatchers.is;

public class HealthIT {
    @Test
    public void shouldBeReadySingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config.xml"),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort))) {

            Eventually.assertDeferred(() -> this.isServiceRunning(app), is(true));
            Eventually.assertDeferred(() -> this.isServiceRunning(app), is(true));
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, HealthServer.PATH_READY), is(200));
        }
    }

    @Test
    public void shouldBeReadyMultipleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config.xml"),
                                                    SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, HealthServer.PATH_READY), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, HealthServer.PATH_READY), is(200));
            }
        }
    }

    @Test
    public void shouldBeLiveSingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config.xml"),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort))) {
            Eventually.assertDeferred(() -> this.isServiceRunning(app), is(true));
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, HealthServer.PATH_HEALTH), is(200));
        }
    }

    @Test
    public void shouldBeLiveMultipleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config.xml"),
                                                    SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, HealthServer.PATH_HEALTH), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, HealthServer.PATH_HEALTH), is(200));
            }
        }
    }

    @Test
    public void shouldBeStatusHASingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config.xml"),
                                                   SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort))) {
            Eventually.assertDeferred(() -> this.isServiceRunning(app), is(true));
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, HealthServer.PATH_HA), is(200));
        }
    }

    @Test
    public void shouldBeStatusHAMultipleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config.xml"),
                                                    SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, HealthServer.PATH_HA), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, HealthServer.PATH_HA), is(200));
            }
        }
    }

    @Test
    public void shouldNotBeStatusHAMultipleMemberWithBackupCountTwo() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config.xml"),
                                                    SystemProperty.of("coherence.distributed.backupcount", 2),
                                                    SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of("coherence.distributed.backupcount", 2),
                                                        SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, HealthServer.PATH_HA), is(400));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, HealthServer.PATH_HA), is(400));
            }
        }
    }

    @Test
    public void shouldNotBeStatusHAMultipleMemberWithBackupCountTwoIgnoringService() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config.xml"),
                                                    SystemProperty.of("coherence.distributed.backupcount", 2),
                                                    SystemProperty.of(HealthServer.PROP_ALLOW_ENDANGERED, "PartitionedCache"),
                                                    SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(HealthServer.PROP_ALLOW_ENDANGERED, "PartitionedCache"),
                                                        SystemProperty.of("coherence.distributed.backupcount", 2),
                                                        SystemProperty.of(HealthServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, HealthServer.PATH_HA), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, HealthServer.PATH_HA), is(200));
            }
        }
    }

    // ----- helper methods -------------------------------------------------

    // Must be public - used in Eventually.assertThat
    public int httpRequest(Capture<Integer> httpPort, String path) {
        try {
            URI uri = URI.create("http://127.0.0.1:" + httpPort.get() + path);
            HttpURLConnection connection = (HttpURLConnection) uri.toURL().openConnection();
            connection.setRequestMethod("GET");
            connection.connect();
            return connection.getResponseCode();
        }
        catch (IOException e) {
            System.err.println("ERROR: HTTP Request failed: " + e.getMessage());
            return -1;
        }
    }

    private boolean isServiceRunning(JavaApplication app) {
        try {
            return app.submit(new IsServiceRunning("PartitionedCache")).get();
        }
        catch (Exception e) {
            System.err.println("ERROR: isServiceRunning failed: " + e.getMessage());
            return false;
        }
    }
}
