/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.security.KeyStore;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.util.Collections;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.Set;
import java.util.function.Supplier;
import java.util.stream.Collectors;

import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLEngine;
import javax.net.ssl.SSLParameters;
import javax.net.ssl.TrustManagerFactory;

import com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.partitionedService.PartitionedCache;
import com.tangosol.coherence.component.util.safeService.safeCacheService.SafeDistributedCacheService;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.Member;
import com.tangosol.net.Service;
import com.tangosol.net.management.MBeanServerProxy;
import com.tangosol.net.management.Registry;
import com.tangosol.util.Filters;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.HttpsConfigurator;
import com.sun.net.httpserver.HttpsParameters;
import com.sun.net.httpserver.HttpsServer;

/**
 * Simple http endpoint for heath checking.
 */
public class OperatorRestServer implements AutoCloseable {
    // ----- constants ------------------------------------------------------

    /**
     * The operator logger to use.
     */
    private static final OperatorLogger LOGGER = OperatorLogger.getLogger();

    /**
     * The system property to use to set the health logging.
     */
    public static final String PROP_HEALTH_LOG = "coherence.k8s.operator.health.logs";

    /**
     * A flag indicating whether debug logging is enabled.
     */
    public static final boolean LOGGING_ENABLED = Boolean.getBoolean(PROP_HEALTH_LOG);

    /**
     * The system property to use to enable the health server.
     */
    public static final String PROP_HEALTH_ENABLED = "coherence.k8s.operator.health.enabled";

    /**
     * The system property to use to set the health port.
     */
    public static final String PROP_HEALTH_PORT = "coherence.k8s.operator.health.port";

    /**
     * The system property to use to determine whether to wait for DCS to start.
     */
    public static final String PROP_WAIT_FOR_DCS = "coherence.k8s.operator.health.wait.dcs";

    /**
     * The system property for the TLS keystore file name.
     */
    public static final String PROP_TLS_KEYSTORE = "coherence.k8s.operator.health.tls.keystore.file";

    /**
     * The system property for the TLS keystore type.
     */
    public static final String PROP_TLS_KEYSTORE_TYPE = "coherence.k8s.operator.health.tls.keystore.type";

    /**
     * The system property for the TLS keystore algorithm.
     */
    public static final String PROP_TLS_KEYSTORE_ALGORITHM = "coherence.k8s.operator.health.tls.keystore.algorithm";

    /**
     * The system property for the TLS keystore password.
     */
    public static final String PROP_TLS_KEYSTORE_PASSWORD = "coherence.k8s.operator.health.tls.keystore.password.plain";

    /**
     * The system property for the TLS keystore password file name.
     */
    public static final String PROP_TLS_KEYSTORE_PASSWORD_FILE = "coherence.k8s.operator.health.tls.keystore.password.file";

    /**
     * The system property for the TLS keystore key password.
     */
    public static final String PROP_TLS_KEY_PASSWORD = "coherence.k8s.operator.health.tls.key.password.plain";

    /**
     * The system property for the TLS keystore key password file name.
     */
    public static final String PROP_TLS_KEY_PASSWORD_FILE = "coherence.k8s.operator.health.tls.key.password.file";

    /**
     * The system property for the TLS trust store file name.
     */
    public static final String PROP_TLS_TRUSTSTORE = "coherence.k8s.operator.health.tls.truststore.file";

    /**
     * The system property for the TLS trust store type.
     */
    public static final String PROP_TLS_TRUSTSTORE_TYPE = "coherence.k8s.operator.health.tls.truststore.type";

    /**
     * The system property for the TLS trust store algorithm.
     */
    public static final String PROP_TLS_TRUSTSTORE_ALGORITHM = "coherence.k8s.operator.health.tls.truststore.algorithm";

    /**
     * The system property for the TLS trust store password.
     */
    public static final String PROP_TLS_TRUSTSTORE_PASSWORD = "coherence.k8s.operator.health.tls.truststore.password.plain";

    /**
     * The system property for the TLS trust store password file name.
     */
    public static final String PROP_TLS_TRUSTSTORE_PASSWORD_FILE = "coherence.k8s.operator.health.tls.truststore.password.file";

    /**
     * The system property for the TLS protocol.
     */
    public static final String PROP_TLS_PROTOCOL = "coherence.k8s.operator.health.tls.protocol";

    /**
     * The system property to indicate whether TLS is 2-way.
     */
    public static final String PROP_TLS_TWO_WAY = "coherence.k8s.operator.health.tls.twoway";

    /**
     * The system property to enable or disable TLS.
     */
    public static final String PROP_INSECURE = "coherence.k8s.operator.health.insecure";

    /**
     * The path to the ready endpoint.
     */
    public static final String PATH_READY = "/ready";

    /**
     * The path to the health endpoint.
     */
    public static final String PATH_HEALTH = "/healthz";

    /**
     * The path to the HA endpoint.
     */
    public static final String PATH_HA = "/ha";

    /**
     * The path to the status endpoint.
     */
    public static final String PATH_STATUS = "/status";

    /**
     * The path to the suspend endpoint.
     */
    public static final String PATH_SUSPEND = "/suspend";

    /**
     * The path to the resume endpoint.
     */
    public static final String PATH_RESUME = "/resume";

    /**
     * The value of the Status HA attribute to signify endangered.
     */
    public static final String STATUS_ENDANGERED = "ENDANGERED";

    /**
     * The error message in an exception due to there being no management member in the cluster.
     */
    public static final String NO_MANAGED_NODES = "None of the nodes are managed";

    /**
     * An empty response body.
     */
    private static final byte[] EMPTY_BODY = new byte[0];

    /**
     * System property to specify service names to be skipped in the StatusHA test.
     */
    public static final String PROP_ALLOW_ENDANGERED = "coherence.k8s.operator.statusha.allowendangered";

    // ----- data members ---------------------------------------------------

    /**
     * The http server.
     */
    private HttpServer httpServer;

    /**
     * The {@link Cluster} supplier.
     */
    private final Supplier<Cluster> clusterSupplier;

    private final Runnable waitForServiceStart;

    /**
     * The {@link java.util.Properties} used to configure the server.
     */
    private final Properties properties;

    /**
     * Flag indicating whether to use TLS.
     */
    private final boolean secure;

    /**
     * Flag indicating whether this application has ever been in a ready state.
     * <p>
     * Because k8s checks the readiness probe a number of time the conditions for the Pod
     * first being ready are different to subsequent checks. This flag is used to determine
     * which checks should be made.
     */
    private volatile boolean hasBeenReady = false;

    // ----- constructors ---------------------------------------------------

    OperatorRestServer(Properties properties) {
        this(CacheFactory::getCluster, OperatorRestServer::waitForDCS, properties);
    }

    OperatorRestServer(Supplier<Cluster> supplier, Runnable waitForServiceStart) {
        this(supplier, waitForServiceStart, System.getProperties());
    }

    OperatorRestServer(Supplier<Cluster> supplier, Runnable waitForServiceStart, Properties properties) {
        this.clusterSupplier = supplier;
        this.waitForServiceStart = waitForServiceStart;
        this.properties = properties;
        this.secure = !Boolean.parseBoolean(properties.getProperty(PROP_INSECURE, "true"));
    }

    // ----- HealthServer methods ------------------------------------------------

    /**
     * Start a http server.
     *
     * @throws Exception if an error occurs
     */
    public synchronized void start() throws Exception {
        boolean enabled = Boolean.parseBoolean(System.getProperty(PROP_HEALTH_ENABLED, "true"));
        if (enabled && httpServer == null) {
            int port = Integer.parseInt(properties.getProperty(PROP_HEALTH_PORT, "6676"));
            HttpServer server = secure ? tlsServer(port) : plainServer(port);

            server.createContext(PATH_READY, this::ready);
            server.createContext(PATH_HEALTH, this::health);
            server.createContext(PATH_HA, this::statusHA);
            server.createContext(PATH_STATUS, this::status);
            server.createContext(PATH_SUSPEND, this::suspend);
            server.createContext(PATH_RESUME, this::resume);

            server.setExecutor(null); // creates a default executor
            server.start();
            System.out.println("CoherenceOperator: listening on " + server.getAddress());

            httpServer = server;
        }
    }

    @Override
    public synchronized void close() {
        if (httpServer != null) {
            httpServer.stop(0);
            httpServer = null;
        }
    }

    public int getPort() {
        return httpServer == null ? -1 : httpServer.getAddress().getPort();
    }

    public boolean isSecure() {
        return secure;
    }

    // ----- helper methods -------------------------------------------------

    private HttpServer plainServer(int port) throws Exception {
        return HttpServer.create(new InetSocketAddress(port), 0);
    }

    private HttpServer tlsServer(int port) throws Exception {
        boolean twoWay = Boolean.parseBoolean(properties.getProperty(PROP_TLS_TWO_WAY, "true"));
        HttpsServer server = HttpsServer.create(new InetSocketAddress(port), 0);
        SSLContext sslContext = createSSLContext(properties);

        server.setHttpsConfigurator(new HttpsConfigurator(sslContext) {
            public void configure(HttpsParameters params) {
                try {
                    // initialise the SSL context
                    params.setSSLParameters(getSSLParameters(twoWay));
                    SSLEngine engine = sslContext.createSSLEngine();
                    params.setCipherSuites(engine.getEnabledCipherSuites());
                    params.setProtocols(engine.getEnabledProtocols());
                    params.setWantClientAuth(twoWay);
                }
                catch (Exception ex) {
                    throw new RuntimeException(ex);
                }
            }
        });

        return server;
    }

    static SSLContext createSSLContext(Properties properties) throws Exception {
        String protocol = properties.getProperty(PROP_TLS_PROTOCOL, "TLS");
        String keystoreFilename = properties.getProperty(PROP_TLS_KEYSTORE);
        char[] keystorePassword = properties.getProperty(PROP_TLS_KEYSTORE_PASSWORD, "").toCharArray();
        String keystorePasswordFile = properties.getProperty(PROP_TLS_KEYSTORE_PASSWORD_FILE);
        char[] keyPassword = properties.getProperty(PROP_TLS_KEY_PASSWORD, "").toCharArray();
        String keyPasswordFile = properties.getProperty(PROP_TLS_KEY_PASSWORD_FILE);
        String keystoreType = properties.getProperty(PROP_TLS_KEYSTORE_TYPE, "JKS");
        String keystoreAlgo = properties.getProperty(PROP_TLS_KEYSTORE_ALGORITHM, "SunX509");
        String truststoreFilename = properties.getProperty(PROP_TLS_TRUSTSTORE);
        char[] truststorePassword = properties.getProperty(PROP_TLS_TRUSTSTORE_PASSWORD, "").toCharArray();
        String truststorePasswordFile = properties.getProperty(PROP_TLS_TRUSTSTORE_PASSWORD_FILE);
        String truststoreType = properties.getProperty(PROP_TLS_TRUSTSTORE_TYPE, "JKS");
        String truststoreAlgo = properties.getProperty(PROP_TLS_TRUSTSTORE_ALGORITHM, "SunX509");

        char[] storepass = FileBasedPasswordProvider.readPassword(keystorePasswordFile, keystorePassword);
        char[] keypass = FileBasedPasswordProvider.readPassword(keyPasswordFile, keyPassword);
        char[] trustpass = FileBasedPasswordProvider.readPassword(truststorePasswordFile, truststorePassword);

        // setup the key manager factory
        URL ksURL = new URL(keystoreFilename);
        KeyStore keystore = KeyStore.getInstance(keystoreType);
        keystore.load(ksURL.openStream(), storepass);
        KeyManagerFactory kmf = KeyManagerFactory.getInstance(keystoreAlgo);
        kmf.init(keystore, keypass);

        // setup the trust manager factory
        URL tsURL = new URL(truststoreFilename);
        KeyStore truststore = KeyStore.getInstance(truststoreType);
        truststore.load(tsURL.openStream(), trustpass);
        TrustManagerFactory tmf = TrustManagerFactory.getInstance(truststoreAlgo);
        tmf.init(truststore);

        // create ssl context
        SSLContext sslContext = SSLContext.getInstance(protocol);

        // setup the HTTPS context and parameters
        sslContext.init(kmf.getKeyManagers(), tmf.getTrustManagers(), null);

        // empty the password arrays
        Arrays.fill(storepass, ' ');
        Arrays.fill(keypass, ' ');
        Arrays.fill(trustpass, ' ');

        return sslContext;
    }

    private SSLParameters getSSLParameters(boolean twoWay) throws NoSuchAlgorithmException {
        SSLContext    ctx         = SSLContext.getDefault();
        SSLEngine     engine = ctx.createSSLEngine();
        SSLParameters params      = ctx.getDefaultSSLParameters();
        String[]      asCiphers   = engine.getEnabledCipherSuites();
        String[]      asProtocols = engine.getEnabledProtocols();

        if (asCiphers != null) {
            params.setCipherSuites(asCiphers);
        }

        if (asProtocols != null) {
            params.setProtocols(asProtocols);
        }

        params.setNeedClientAuth(twoWay);

        return params;
        }

    /**
     * Send a http response.
     *
     * @param t      the {@link HttpExchange} to send the response to
     * @param status the response status
     */
    private static void send(HttpExchange t, int status) {
        try {
            t.sendResponseHeaders(status, 0);
            try (OutputStream os = t.getResponseBody()) {
                os.write(EMPTY_BODY);
            }
        }
        catch (IOException e) {
            e.printStackTrace();
        }
    }

    /**
     * Send a http response.
     *
     * @param t      the {@link HttpExchange} to send the response to
     * @param status the response status
     * @param body   the response body
     */
    private static void send(HttpExchange t, int status, String body) {
        try {
            byte[] bytes = body == null ? EMPTY_BODY : body.getBytes(StandardCharsets.UTF_8);
            t.sendResponseHeaders(status, bytes.length);
            try (OutputStream os = t.getResponseBody()) {
                os.write(bytes);
            }
        }
        catch (IOException e) {
            e.printStackTrace();
        }
    }

    /**
     * Process a ready request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void ready(HttpExchange exchange) {
        try {
            boolean hasCluster = hasClusterMembers();
            int response = 400;
            if (hasBeenReady) {
                response = hasCluster ? 200 : 400;
                LOGGER.debug("CoherenceOperator: Ready check response %d - cluster=%b", response, hasCluster);
            }
            else {
                boolean isHA = hasCluster && isStatusHA();
                boolean isIdle = hasCluster && isPersistenceIdle();
                if (hasCluster && isHA && isIdle) {
                    response = 200;
                    hasBeenReady = true;
                }
                LOGGER.debug("CoherenceOperator: Ready check response %d - cluster=%b HA=%b Idle=%b",
                         response, hasCluster, isHA, isIdle);
            }
            send(exchange, response);
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "Ready check");
        }
    }

    /**
     * Process a health request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void health(HttpExchange exchange) {
        try {
            boolean hasCluster = hasClusterMembers();
            int response = hasClusterMembers() ? 200 : 400;
            LOGGER.debug("CoherenceOperator: Health check response %d - cluster=%b", response, hasCluster);
            send(exchange, response);
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "Health check");
        }
    }

    /**
     * Process a status HA request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void statusHA(HttpExchange exchange) {
        try {
            boolean isHA = isStatusHA();
            boolean isIdle = isPersistenceIdle();
            int response = isHA && isIdle ? 200 : 400;
            if (response == 400 || LOGGING_ENABLED) {
                LOGGER.info("CoherenceOperator: HA check response %d - HA=%b Idle=%b", response, isHA, isIdle);
            }
            send(exchange, response);
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "StatusHA check");
        }
    }

    /**
     * Process a status request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void status(HttpExchange exchange) {
        try {
            String status = getStatusName();
            send(exchange, 200, status.toLowerCase());
        }
        catch (Throwable thrown) {
            handleError(exchange, thrown, "Status check");
        }
    }

    /**
     * Process a suspend request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void suspend(HttpExchange exchange) {
        try {
            String path = exchange.getRequestURI().getPath();
            String name = "";
            String[] parts = path.split("/");
            if (parts.length > 2) {
                name = parts[2].trim();
            }

            Cluster cluster = clusterSupplier.get();
            Registry registry = cluster.getManagement();
            MBeanServerProxy proxy = registry.getMBeanServerProxy();
            Map<Integer, String> identityMap = new HashMap<>();

            Set<String> mbeanNames = proxy.queryNames(":type=KubernetesOperator,nodeId=*", Filters.always());
            for (String mbeanName : mbeanNames) {
                Map<String, Object> attributes = proxy.getAttributes(mbeanName, Filters.always());
                identityMap.put((Integer) attributes.get(CoherenceOperatorMBean.ATTRIBUTE_NODE),
                                (String) attributes.get(CoherenceOperatorMBean.ATTRIBUTE_IDENTITY));
            }

            if (!name.isEmpty()) {
                Service service = cluster.getService(name);
                if (service == null) {
                    send(exchange, 404);
                    return;
                }
                LOGGER.warn("CoherenceOperator: Suspending service %s", name);
                cluster.suspendService(name);
            }
            else {
                LOGGER.warn("CoherenceOperator: Suspending all persistence enabled services");
                Enumeration<String> names = cluster.getServiceNames();
                while (names.hasMoreElements()) {
                    name = names.nextElement();
                    Service svc = cluster.getService(name);
                    LOGGER.debug("CoherenceOperator: Suspending all persistence enabled services - service=%s", name);
                    if (svc instanceof DistributedCacheService && ((DistributedCacheService) svc).isLocalStorageEnabled()) {
                        DistributedCacheService dcs = (DistributedCacheService) svc;

                        if (PersistenceHelper.isActivePersistenceEnabled(dcs)) {
                            long count = dcs.getOwnershipEnabledMembers().stream()
                                    .map(m -> identityMap.get(m.getId()))
                                    .distinct()
                                    .count();

                            if (count == 1) {
                                LOGGER.info("CoherenceOperator: Suspending service %s", name);
                                cluster.suspendService(name);
                            }
                            else {
                                LOGGER.info("CoherenceOperator: Not suspending service %s "
                                        + "- is storage enabled in other deployments", name);
                            }
                        }
                        else {
                            LOGGER.debug("CoherenceOperator: Suspending all persistence enabled services "
                                    + "- service=%s does not have persistence enabled", name);
                        }
                    }
                    else {
                        LOGGER.debug("CoherenceOperator: Suspending all persistence enabled services "
                                + "- service=%s is not a storage enabled DistributedCacheService", name);
                    }
                }
            }
            send(exchange, 200);
        }
        catch (Exception e) {
            CacheFactory.err(e);
            send(exchange, 500);
        }
    }

    /**
     * Process a resume request.
     *
     * @param exchange the {@link HttpExchange} to send the response to
     */
    void resume(HttpExchange exchange) {
        try {
            String path = exchange.getRequestURI().getPath();
            String name = "";
            String[] parts = path.split("/");
            if (parts.length > 2) {
                name = parts[2].trim();
            }

            Cluster cluster = clusterSupplier.get();
            if (!name.isEmpty()) {
                Service service = cluster.getService(name);
                if (service == null) {
                    send(exchange, 404);
                    return;
                }
                LOGGER.warn("CoherenceOperator: Resuming service %s", name);
                cluster.resumeService(name);
            }
            else {
                LOGGER.warn("CoherenceOperator: Resuming all services");
                Enumeration<String> names = cluster.getServiceNames();
                while (names.hasMoreElements()) {
                    cluster.resumeService(names.nextElement());
                }
            }
            send(exchange, 200);
        }
        catch (Exception e) {
            CacheFactory.err(e);
            send(exchange, 500);
        }
    }

    private void handleError(HttpExchange t, Throwable thrown, String action) {
        String msg = thrown.getMessage();
        LOGGER.error("CoherenceOperator: %s failed due to '%s'", action, thrown.getMessage());
        if (msg != null && msg.contains(NO_MANAGED_NODES)) {
            send(t, 400);
        }
        else {
            CacheFactory.err(thrown);
            send(t, 500);
        }
    }

    /**
     * Determine whether there are any members in the cluster.
     *
     * @return {@code true} if the Coherence cluster has members
     */
    protected boolean hasClusterMembers() {
        Cluster cluster = clusterSupplier.get();
        return cluster != null && cluster.isRunning() && !cluster.getMemberSet().isEmpty();
    }

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    boolean isStatusHA() {
        String exclusions = properties.getProperty(PROP_ALLOW_ENDANGERED);
        return isStatusHA(exclusions);
    }

    /**
     * Returns {@code true} if all Coherence services are safe.
     *
     * @param exclusions  the service names that are allowed to be endangered
     *
     * @return {@code true} if all Coherence services are safe
     */
    boolean isStatusHA(String exclusions) {
        try {
            LOGGER.debug("CoherenceOperator: StatusHA check. Waiting for service start...");
            waitForServiceStart.run();
            LOGGER.debug("CoherenceOperator: StatusHA check. services started");

            Set<String> allowEndangered;
            if (exclusions != null) {
                allowEndangered = Arrays.stream(exclusions.split(","))
                        .collect(Collectors.toSet());
            }
            else {
                allowEndangered = Collections.emptySet();
            }

            Cluster cluster = clusterSupplier.get();
            if (cluster != null && cluster.isRunning()) {
                return areCacheServicesHA(cluster, allowEndangered);
            }
            else {
                String reason = cluster == null ? "null" : "not running";
                LOGGER.error("CoherenceOperator: StatusHA check failed - cluster is " + reason);
                return false;
            }
        }
        catch (Exception e) {
            // there is probably no DCS
            LOGGER.error("CoherenceOperator: StatusHA check failed, %s", e.getMessage());
            return false;
        }
    }

    /**
     * Returns {@code true} if all Coherence services are safe.
     *
     * @param cluster          the Coherence cluster
     * @param allowEndangered  the service names that are allowed to be endangered
     *
     * @return {@code true} if all Coherence services are safe
     */
    @SuppressWarnings("unchecked")
    boolean areCacheServicesHA(Cluster cluster, Set<String> allowEndangered) {
        LOGGER.debug("CoherenceOperator: Checking HA: allowEndangered=%s", allowEndangered);

        Enumeration<String> names = cluster.getServiceNames();
        while (names.hasMoreElements()) {
            String name = names.nextElement();

            Service service = cluster.getService(name);
            if (service instanceof DistributedCacheService && ((DistributedCacheService) service).isLocalStorageEnabled()) {
                if (service instanceof SafeDistributedCacheService) {
                    service = ((SafeDistributedCacheService) service).getService();
                }

                PartitionedCache partitionedCache = (PartitionedCache) service;

                if (partitionedCache.isOwnershipEnabled()) {
                    Set<Member> setOwnershipEnabledMembers = partitionedCache.getOwnershipEnabledMembers();
                    int memberCount                        = setOwnershipEnabledMembers.size();

                    if (memberCount == 1) {
                        // storage enabled and only one member, check we own all partitions
                        int partitionCount = partitionedCache.getPartitionCount();
                        int ownedPartitions = partitionedCache.calculateThisOwnership(true);
                        if (ownedPartitions != partitionCount) {
                            LOGGER.debug("CoherenceOperator: StatusHA check failed. "
                                             + "Service %s this member is the only storage enabled member, "
                                             + "but owns only %d of %d partitions",
                                     name, ownedPartitions, partitionCount);
                            return false;
                        }
                    }

                    String sMembersIds = setOwnershipEnabledMembers.stream()
                                                                   .map(Member::getId)
                                                                   .map(String::valueOf)
                                                                   .collect(Collectors.joining(","));
                    String sMembers = String.format("memberCount=%d members=[%s]", memberCount, sMembersIds);

                    String statusHA = partitionedCache.getBackupStrengthName();
                    int backupCount = partitionedCache.getBackupCount();

                    if (memberCount > 1
                        && backupCount > 0
                        && STATUS_ENDANGERED.equals(statusHA)
                        && !allowEndangered.contains(name)) {
                        LOGGER.error("CoherenceOperator: StatusHA check failed. Service %s has HA status of %s, suspended=%b, %s",
                                name, statusHA, partitionedCache.isSuspended(), sMembers);
                        return false;
                    }

                    if (partitionedCache.isDistributionInProgress()) {
                        LOGGER.error("CoherenceOperator: StatusHA check failed. Service %s distribution in progress, %s",
                                name, sMembers);
                        return false;
                    }

                    if (partitionedCache.isRecoveryInProgress()) {
                        LOGGER.error("CoherenceOperator: StatusHA check failed. Service %s recovery in progress, %s",
                                name, sMembers);
                        return false;
                    }

                    if (partitionedCache.isRestoreInProgress()) {
                        LOGGER.error("CoherenceOperator: StatusHA check failed. Service %s restore in progress, %s",
                                name, sMembers);
                        return false;
                    }

                    if (partitionedCache.isTransferInProgress()) {
                        LOGGER.error("CoherenceOperator: StatusHA check failed. Service %s transfer in progress, %s",
                                name, sMembers);
                        return false;
                    }
                }
            }
        }
        return true;
    }

    /**
     * Returns {@code true} if all persistence enabled services are idle.
     *
     * @return {@code true} if all persistence enabled services are idle
     */
    boolean isPersistenceIdle() {
        boolean allIdle = true;
        Cluster cluster = clusterSupplier.get();
        Enumeration<String> names = cluster.getServiceNames();

        while (names.hasMoreElements()) {
            String name = names.nextElement();
            Service service = cluster.getService(name);

            if (service instanceof DistributedCacheService && ((DistributedCacheService) service).isLocalStorageEnabled()) {
                if (service instanceof SafeDistributedCacheService) {
                    service = ((SafeDistributedCacheService) service).getService();
                }
                if (PersistenceHelper.isActive(service)) {
                    LOGGER.debug("CoherenceOperator: Persistence not idle for service %s" + name);
                    allIdle = false;
                }
            }
        }
        return allIdle;
    }

    /**
     * Returns the name of the lowest Coherence cache service Status HA name.
     *
     * @return the name of the lowest Coherence cache service Status HA name
     */
    private String getStatusName() {
        Cluster cluster = clusterSupplier.get();
        int lowestStatus = Integer.MAX_VALUE;
        String status = "n/a";

        if (cluster != null && cluster.isRunning()) {
            Enumeration<String> names = cluster.getServiceNames();
            while (names.hasMoreElements()) {
                String name = names.nextElement();
                Service service = cluster.getService(name);

                if (service instanceof DistributedCacheService && ((DistributedCacheService) service).isLocalStorageEnabled()) {
                    if (service instanceof SafeDistributedCacheService) {
                        service = ((SafeDistributedCacheService) service).getService();
                    }
                    PartitionedCache partitionedCache = (PartitionedCache) service;
                    if (partitionedCache.isOwnershipEnabled()) {
                        int code = partitionedCache.getBackupStrength();
                        if (code < lowestStatus) {
                            status = partitionedCache.getBackupStrengthName();
                        }
                    }
                }
            }
        }
        return status;
    }

    private static void waitForDCS() {
        String s = System.getProperty(PROP_WAIT_FOR_DCS, "false");
        if (Boolean.parseBoolean(s)) {
            DefaultCacheServer dcs = DefaultCacheServer.getInstance();
            dcs.waitForServiceStart();
        }
    }
}
