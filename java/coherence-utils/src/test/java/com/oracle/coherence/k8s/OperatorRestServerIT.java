/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URI;

import com.tangosol.coherence.component.util.SafeService;
import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.partitionedService.PartitionedCache;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.Service;

import com.oracle.bedrock.runtime.LocalPlatform;
import com.oracle.bedrock.runtime.coherence.callables.IsServiceRunning;
import com.oracle.bedrock.runtime.coherence.options.CacheConfig;
import com.oracle.bedrock.runtime.coherence.options.LocalStorage;
import com.oracle.bedrock.runtime.concurrent.RemoteCallable;
import com.oracle.bedrock.runtime.java.JavaApplication;
import com.oracle.bedrock.runtime.java.options.ClassName;
import com.oracle.bedrock.runtime.java.options.SystemProperty;
import com.oracle.bedrock.testsupport.deferred.Eventually;
import com.oracle.bedrock.util.Capture;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

public class OperatorRestServerIT {
    @Test
    public void shouldBeReadySingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config.xml"),
                                                   SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort))) {

            Eventually.assertDeferred(() -> this.isServiceOneRunning(app), is(true));
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_READY), is(200));
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
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_READY), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, OperatorRestServer.PATH_READY), is(200));
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
                                                   SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort))) {
            Eventually.assertDeferred(() -> this.isServiceOneRunning(app), is(true));
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_HEALTH), is(200));
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
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_HEALTH), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, OperatorRestServer.PATH_HEALTH), is(200));
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
                                                   SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort))) {
            Eventually.assertDeferred(() -> this.isServiceOneRunning(app), is(true));
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_HA), is(200));
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
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_HA), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, OperatorRestServer.PATH_HA), is(200));
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
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of("coherence.distributed.backupcount", 2),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_HA), is(400));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, OperatorRestServer.PATH_HA), is(400));
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
                                                    SystemProperty.of(OperatorRestServer.PROP_ALLOW_ENDANGERED, "PartitionedCacheOne"),
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config.xml"),
                                                        SystemProperty.of(OperatorRestServer.PROP_ALLOW_ENDANGERED, "PartitionedCacheOne"),
                                                        SystemProperty.of("coherence.distributed.backupcount", 2),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_HA), is(200));
                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, OperatorRestServer.PATH_HA), is(200));
            }
        }
    }

    @Test
    public void shouldSuspendAllServicesSingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config-two.xml"),
                                                   SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort))) {

            Eventually.assertDeferred(() -> this.isServiceOneRunning(app), is(true));
            Eventually.assertDeferred(() -> this.isServiceTwoRunning(app), is(true));

            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_SUSPEND), is(200));
            Eventually.assertDeferred(() -> this.isServiceOneSuspended(app), is(true));
            Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app), is(true));

            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_RESUME), is(200));
            Eventually.assertDeferred(() -> this.isServiceOneSuspended(app), is(false));
            Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app), is(false));
        }
    }

    @Test
    public void shouldSuspendSpecifiedServicesSingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config-two.xml"),
                                                   SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort))) {

            Eventually.assertDeferred(() -> this.isServiceOneRunning(app), is(true));
            Eventually.assertDeferred(() -> this.isServiceTwoRunning(app), is(true));

            String path = OperatorRestServer.PATH_SUSPEND + "/PartitionedCacheOne";
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, path), is(200));
            Eventually.assertDeferred(() -> this.isServiceOneSuspended(app), is(true));
            assertThat(isServiceTwoSuspended(app), is(false));

            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_RESUME), is(200));
            Eventually.assertDeferred(() -> this.isServiceOneSuspended(app), is(false));
            Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app), is(false));
        }
    }

    @Test
    public void shouldResumeSpecifiedServicesSingleMember() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app = platform.launch(JavaApplication.class,
                                                   ClassName.of(Main.class),
                                                   CacheConfig.of("test-cache-config-two.xml"),
                                                   SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort))) {

            Eventually.assertDeferred(() -> this.isServiceOneRunning(app), is(true));
            Eventually.assertDeferred(() -> this.isServiceTwoRunning(app), is(true));

            Eventually.assertDeferred(() -> this.httpRequest(httpPort, OperatorRestServer.PATH_SUSPEND), is(200));
            Eventually.assertDeferred(() -> this.isServiceOneSuspended(app), is(true));
            Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app), is(true));

            String path = OperatorRestServer.PATH_RESUME + "/PartitionedCacheOne";
            Eventually.assertDeferred(() -> this.httpRequest(httpPort, path), is(200));
            Eventually.assertDeferred(() -> this.isServiceOneSuspended(app), is(false));
            assertThat(isServiceTwoSuspended(app), is(true));
        }
    }

    @Test
    public void shouldSuspendAllServicesMultipleMembers() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config-two.xml"),
                                                    SystemProperty.of(CoherenceOperatorMBean.PROP_IDENTITY, "foo"),
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config-two.xml"),
                                                        SystemProperty.of(CoherenceOperatorMBean.PROP_IDENTITY, "bar"),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));

                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_SUSPEND), is(200));
                Eventually.assertDeferred(() -> this.isServiceOneSuspended(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneSuspended(app2), is(true));
                Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app2), is(true));

                Eventually.assertDeferred(() -> this.httpRequest(httpPort2, OperatorRestServer.PATH_RESUME), is(200));
                Eventually.assertDeferred(() -> this.isServiceOneSuspended(app1), is(false));
                Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app1), is(false));
                Eventually.assertDeferred(() -> this.isServiceOneSuspended(app2), is(false));
                Eventually.assertDeferred(() -> this.isServiceTwoSuspended(app2), is(false));
            }
        }
    }

    @Test
    public void shouldNotSuspendServicesWithDifferentIdentities() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config-two.xml"),
                                                    SystemProperty.of(CoherenceOperatorMBean.PROP_IDENTITY, "foo"),
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config-two.xml"),
                                                        SystemProperty.of(CoherenceOperatorMBean.PROP_IDENTITY, "bar"),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));

                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_SUSPEND), is(200));
                assertThat(isServiceOneSuspended(app1), is(false));
                assertThat(isServiceTwoSuspended(app1), is(false));
                assertThat(isServiceOneSuspended(app2), is(false));
                assertThat(isServiceTwoSuspended(app2), is(false));
            }
        }
    }

    @Test
    public void shouldNotSuspendStorageDisabledServices() throws Exception {
        LocalPlatform platform = LocalPlatform.get();
        Capture<Integer> httpPort1 = new Capture<>(platform.getAvailablePorts());
        Capture<Integer> httpPort2 = new Capture<>(platform.getAvailablePorts());

        try (JavaApplication app1 = platform.launch(JavaApplication.class,
                                                    ClassName.of(Main.class),
                                                    CacheConfig.of("test-cache-config-two.xml"),
                                                    LocalStorage.disabled(),
                                                    SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort1))) {
            try (JavaApplication app2 = platform.launch(JavaApplication.class,
                                                        ClassName.of(Main.class),
                                                        CacheConfig.of("test-cache-config-two.xml"),
                                                        LocalStorage.disabled(),
                                                        SystemProperty.of(OperatorRestServer.PROP_HEALTH_PORT, httpPort2))) {
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app1), is(true));
                Eventually.assertDeferred(() -> this.isServiceOneRunning(app2), is(true));

                Eventually.assertDeferred(() -> this.httpRequest(httpPort1, OperatorRestServer.PATH_SUSPEND), is(200));
                assertThat(isServiceOneSuspended(app1), is(false));
                assertThat(isServiceTwoSuspended(app1), is(false));
                assertThat(isServiceOneSuspended(app2), is(false));
                assertThat(isServiceTwoSuspended(app2), is(false));
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

    private boolean isServiceOneRunning(JavaApplication app) {
        try {
            return app.submit(new IsServiceRunning("PartitionedCacheOne")).get();
        }
        catch (Exception e) {
            System.err.println("ERROR: isServiceRunning failed: " + e.getMessage());
            return false;
        }
    }

    private boolean isServiceTwoRunning(JavaApplication app) {
        try {
            return app.submit(new IsServiceRunning("PartitionedCacheOne")).get();
        }
        catch (Exception e) {
            System.err.println("ERROR: isServiceRunning failed: " + e.getMessage());
            return false;
        }
    }

    private boolean isServiceOneSuspended(JavaApplication app) {
        return isServiceSuspended(app, "PartitionedCacheOne");
    }

    private boolean isServiceTwoSuspended(JavaApplication app) {
        return isServiceSuspended(app, "PartitionedCacheTwo");
    }

    private boolean isServiceSuspended(JavaApplication app, String svc) {
        try {
            return app.submit(new IsServiceSuspended(svc)).get();
        }
        catch (Exception e) {
            System.err.println("ERROR: isServiceSuspended failed: " + e.getMessage());
            return false;
        }
    }

    // ----- inner class: IsServiceSuspended --------------------------------

    public static class IsServiceSuspended
            implements RemoteCallable<Boolean> {
        /**
         * The name of the service.
         */
        private String serviceName;

        /**
         * Constructs an {@link IsServiceSuspended}
         *
         * @param serviceName the name of the service
         */
        public IsServiceSuspended(String serviceName) {
            this.serviceName = serviceName;
        }

        @Override
        public Boolean call() {
            Cluster cluster = CacheFactory.getCluster();
            Service service = cluster.getService(serviceName);
            if (service instanceof SafeService) {
                service = ((SafeService) service).getService();
            }
            PartitionedCache partitionedCache = (PartitionedCache) service;
            return partitionedCache.isSuspended();
        }
    }

}
