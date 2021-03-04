/*
 * Copyright (c) 2019, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Collections;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import java.util.Set;
import java.util.function.Supplier;
import java.util.stream.Collectors;

import javax.management.MalformedObjectNameException;
import javax.management.ObjectName;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.DistributedCacheService;
import com.tangosol.net.Service;
import com.tangosol.net.management.MBeanHelper;
import com.tangosol.net.management.MBeanServerProxy;
import com.tangosol.net.management.Registry;
import com.tangosol.util.Filters;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

/**
 * Simple http endpoint for heath checking.
 */
public class OperatorRestServer {
    // ----- constants ------------------------------------------------------

    /**
     * The system property to use to set the health logging.
     */
    public static final String PROP_HEALTH_LOG = "coherence.k8s.operator.health.logs";

    /**
     * A flag indicating whether debug logging is enabled.
     */
    public static final boolean LOGGING_ENABLED = Boolean.getBoolean(PROP_HEALTH_LOG);

    /**
     * The system property to use to set the health port.
     */
    public static final String PROP_HEALTH_PORT = "coherence.k8s.operator.health.port";

    /**
     * The system property to use to determine whether to wait for DCS to start.
     */
    public static final String PROP_WAIT_FOR_DCS = "coherence.k8s.operator.health.wait.dcs";

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
     * The MBean name of the PartitionAssignment MBean.
     */
    public static final String MBEAN_PARTITION_ASSIGNMENT = Registry.PARTITION_ASSIGNMENT_TYPE
            + ",service=*,responsibility=DistributionCoordinator";

    /**
     * The MBean name of the Service MBean.
     */
    public static final String MBEAN_SERVICE = "%s:" + Registry.SERVICE_TYPE
            + ",name=%s,nodeId=%d";

    /**
     * Service MBean Attributes required to compute HAStatus.
     *
     * @see #isServiceStatusHA(String, java.util.Map)
     */
    public static final String[] SERVICE_STATUS_HA_ATTRIBUTES =
            {
                    "HAStatus",
                    "HAStatusCode",
                    "BackupCount",
                    "ServiceNodeCount",
                    "RemainingDistributionCount"
            };

    /**
     * The MBean name of the Persistence Coordinator MBean.
     */
    public static final String MBEAN_PERSISTENCE_COORDINATOR = Registry.PERSISTENCE_SNAPSHOT_TYPE
            + ",service=*,responsibility=PersistenceCoordinator";

    /**
     * The MBean attribute to check the idle state of the persistence coordinator.
     */
    public static final String[] PERSISTENCE_IDLE_ATTRIBUTES = new String[] {"Idle"};

    /**
     * The MBean attribute to check the state of a partitioned cache service.
     */
    public static final String[] CACHE_SERVICE_ATTRIBUTES = new String[] {"StorageEnabled", "MemberCount",
            "OwnedPartitionsPrimary", "PartitionsAll"};

    /**
     * The value of the Status HA attribute to signify endangered.
     */
    public static final String STATUS_ENDANGERED = "ENDANGERED";

    /**
     * The name of the HA status MBean attribute.
     */
    public static final String ATTRIB_HASTATUS = "hastatus";

    /**
     * The name of the Remaining Distribution Count MBean attribute.
     */
    public static final String ATTRIB_REMAINING_DISTRIBUTION_COUNT = "remainingdistributioncount";

    /**
     * The name of the HA status code MBean attribute.
     */
    public static final String ATTRIB_HASTATUS_CODE = "hastatuscode";

    /**
     * The name of the backup count MBean attribute.
     */
    public static final String ATTRIB_BACKUPS = "backupcount";

    /**
     * The name of the service node count MBean attribute.
     */
    public static final String ATTRIB_NODE_COUNT = "servicenodecount";

    /**
     * The name of the persistence coordinator idle state MBean attribute.
     */
    public static final String ATTRIB_IDLE = "idle";

    /**
     * The name of the persistence coordinator idle state MBean attribute.
     */
    public static final String ATTRIB_STORAGE_ENABLED = "storageenabled";

    /**
     * The name of the member count MBean attribute.
     */
    public static final String ATTRIB_MEMBER_COUNT = "membercount";

    /**
     * The name of the owned primary partitions MBean attribute.
     */
    public static final String ATTRIB_OWNED_PARTITIONS_PRIMARY = "ownedpartitionsprimary";

    /**
     * The name of the partition count MBean attribute.
     */
    public static final String ATTRIB_PARTITIONS_ALL = "partitionsall";

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
     * Flag indicating whether this application has ever been in a ready state.
     * <p>
     * Because k8s checks the readiness probe a number of time the conditions for the Pod
     * first being ready are different to subsequent checks. This flag is used to determine
     * which checks should be made.
     */
    private volatile boolean hasBeenReady = false;

    // ----- constructors ---------------------------------------------------

    OperatorRestServer() {
        this(CacheFactory::getCluster, OperatorRestServer::waitForDCS);
    }

    OperatorRestServer(Supplier<Cluster> supplier, Runnable waitForServiceStart) {
        this.clusterSupplier = supplier;
        this.waitForServiceStart = waitForServiceStart;
    }

    // ----- HealthServer methods ------------------------------------------------

    /**
     * Start a http server.
     *
     * @throws IOException if an error occurs
     */
    public synchronized void start() throws IOException {
        if (httpServer == null) {
            int port = Integer.getInteger(PROP_HEALTH_PORT, 6676);
            HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);

            server.createContext(PATH_READY, this::ready);
            server.createContext(PATH_HEALTH, this::health);
            server.createContext(PATH_HA, this::statusHA);
            server.createContext(PATH_STATUS, this::status);
            server.createContext(PATH_SUSPEND, this::suspend);
            server.createContext(PATH_RESUME, this::resume);

            server.setExecutor(null); // creates a default executor
            server.start();

            log("CoherenceOperator REST server is listening on http://localhost:%d", server.getAddress().getPort());

            httpServer = server;
        }
    }

    // ----- helper methods -------------------------------------------------

    /**
     * Send a http response.
     *
     * @param t      the {@link HttpExchange} to send the response to
     * @param status the response status
     */
    private static void send(HttpExchange t, int status) {
        try {
            t.sendResponseHeaders(status, 0);
            OutputStream os = t.getResponseBody();
            os.write(EMPTY_BODY);
            os.close();
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
            OutputStream os = t.getResponseBody();
            os.write(bytes);
            os.close();
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
                logDebug("CoherenceOperator: Ready check response %d - cluster=%b", response, hasCluster);
            }
            else {
                boolean isHA = isStatusHA();
                boolean isIdle = isPersistenceIdle();
                if (hasCluster && isHA && isIdle) {
                    response = 200;
                    hasBeenReady = true;
                }
                logDebug("CoherenceOperator: Ready check response %d - cluster=%b HA=%b Idle=%b",
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
            logDebug("CoherenceOperator: Health check response %d - cluster=%b", response, hasCluster);
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
            if (response == 400) {
                logDebug("CoherenceOperator: HA check response %d - HA=%b Idle=%b", response, isHA, isIdle);
            }
            else {
                log("CoherenceOperator: HA check response %d - HA=%b Idle=%b", response, isHA, isIdle);
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
    @SuppressWarnings("unchecked")
    void suspend(HttpExchange exchange) {
        try {
            String   path  = exchange.getRequestURI().getPath();
            String   name  = "";
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
                warn("CoherenceOperator: Suspending service %s", name);
                cluster.suspendService(name);
            }
            else {
                warn("CoherenceOperator: Suspending all services");
                Enumeration<String> names = cluster.getServiceNames();
                while (names.hasMoreElements()) {
                    name = names.nextElement();
                    Service svc = cluster.getService(name);
                    if (svc instanceof DistributedCacheService && ((DistributedCacheService) svc).isLocalStorageEnabled()) {
                        DistributedCacheService dcs = (DistributedCacheService) svc;
                        long count = dcs.getOwnershipEnabledMembers().stream()
                                        .map(m -> identityMap.get(m.getId()))
                                        .distinct()
                                        .count();
                        if (count == 1) {
                            log("CoherenceOperator: Suspending service %s", name);
                            cluster.suspendService(name);
                        }
                        else {
                            log("CoherenceOperator: Not suspending service %s - is storage enabled in other deployments", name);
                        }
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
            String   path  = exchange.getRequestURI().getPath();
            String   name  = "";
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
                warn("CoherenceOperator: Resuming service %s", name);
                cluster.resumeService(name);
            }
            else {
                warn("CoherenceOperator: Resuming all services");
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
        err("CoherenceOperator: %s failed due to '%s'", action, thrown.getMessage());
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
    private boolean hasClusterMembers() {
        Cluster cluster = clusterSupplier.get();
        return cluster != null && cluster.isRunning() && !cluster.getMemberSet().isEmpty();
    }

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    boolean isStatusHA() {
        String exclusions = System.getProperty(PROP_ALLOW_ENDANGERED);
        return isStatusHA(exclusions);
    }

    boolean isStatusHA(String exclusions) {
        try {
            waitForServiceStart.run();

            Set<String> allowEndangered = null;
            if (exclusions != null) {
                allowEndangered = Arrays.stream(exclusions.split(","))
                    .map(this::quoteMBeanName)
                    .map(s -> ",service=" + s + ",")
                    .collect(Collectors.toSet());
            }

            Cluster cluster = clusterSupplier.get();
            if (cluster != null && cluster.isRunning()) {
                int id = cluster.getLocalMember().getId();
                for (String mBean : getPartitionAssignmentMBeans()) {
                    if (allowEndangered != null && allowEndangered.stream().anyMatch(mBean::contains)) {
                        // this service is allowed to be endangered so skip it.
                        continue;
                    }
                    Map<String, Object> attributes = getMBeanServiceStatusHAAttributes(mBean);
                    if (!isServiceStatusHA(mBean, attributes)) {
                        return false;
                    }
                    if (!isCacheServiceSafe(mBean, id)) {
                        return false;
                    }
                }
                return true;
            }
            else {
                err("CoherenceOperator: StatusHA check failed - cluster is null");
                return false;
            }
        }
        catch (Exception e) {
            // there is probably no DCS
            err("CoherenceOperator: StatusHA check failed, %s", e.getMessage());
            return false;
        }
    }

    boolean isPersistenceIdle() {
        boolean allIdle = true;

        for (String mBean : getPersistenceCoordinatorMBeans()) {
            Map<String, Object> attributes = getMBeanAttributes(mBean, PERSISTENCE_IDLE_ATTRIBUTES);
            Boolean isIdle = (Boolean) attributes.get(ATTRIB_IDLE);
            if (!isIdle) {
                log("CoherenceOperator: Persistence not idle for MBean %s" + mBean);
                allIdle = false;
            }
        }

        return allIdle;
    }

    /**
     * Check that a given service is safe.
     * <p>
     * If the service has more than one member this method returns {@code true}.
     * If there is only a single storage enabled member then verify that the single member
     * owns all of the partitions. This ensures that when using active persistence the
     * member is not still creating stores.
     *
     * @param mBean     the name of the paersistence manager MBean
     * @param memberId  this member Id
     *
     * @return true if the service is safe
     * @throws MalformedObjectNameException if there is an error creating the MBean name
     */
    private boolean isCacheServiceSafe(String mBean, int memberId) throws MalformedObjectNameException {
        ObjectName objectName = ObjectName.getInstance(mBean);
        String domain = objectName.getDomain();
        String serviceName = objectName.getKeyProperty("service");
        String serviceMBean = String.format(MBEAN_SERVICE, domain, serviceName, memberId);
        Map<String, Object> attributes = getMBeanAttributes(serviceMBean, CACHE_SERVICE_ATTRIBUTES);
        Boolean storageEnabled = (Boolean) attributes.get(ATTRIB_STORAGE_ENABLED);
        Integer memberCount = (Integer) attributes.get(ATTRIB_MEMBER_COUNT);
        Integer ownedPartitions = (Integer) attributes.get(ATTRIB_OWNED_PARTITIONS_PRIMARY);
        Integer partitionCount = (Integer) attributes.get(ATTRIB_PARTITIONS_ALL);
        boolean safe = true;

        if (storageEnabled != null && storageEnabled && memberCount != null && memberCount == 1) {
            // storage enabled and only one member, check we own all partitions
            safe = ownedPartitions != null && partitionCount != null && ownedPartitions.intValue() == partitionCount.intValue();
            if (!safe) {
                log("CoherenceOperator: Partitioned Cache Service MBean %s is not safe - %s", serviceMBean, attributes);
            }
        }

        if (safe) {
            logDebug("CoherenceOperator: Partitioned Cache Service MBean %s is safe - %s", serviceMBean, attributes);
        }
        return safe;
    }

    private String quoteMBeanName(String sMBean) {
        if (MBeanHelper.isQuoteRequired(sMBean)) {
            return MBeanHelper.quote(sMBean);
        }
        return sMBean;
    }

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    private String getStatusName() {
        Cluster cluster = clusterSupplier.get();
        int lowestStatus = Integer.MAX_VALUE;
        String status = "n/a";

        if (cluster != null && cluster.isRunning()) {
            for (String mBean : getPartitionAssignmentMBeans()) {
                Map<String, Object> attributes = getMBeanServiceStatusHAAttributes(mBean);
                // convert the attribute case as MBeanProxy or REST return them with different cases
                Map<String, Object> map = attributes.entrySet()
                        .stream()
                        .collect(Collectors.toMap(e -> e.getKey().toLowerCase(), Map.Entry::getValue));

                Integer code = (Integer) map.get(ATTRIB_HASTATUS_CODE);
                if (code != null && code < lowestStatus) {
                    status = (String) map.get(ATTRIB_HASTATUS);
                }
            }
        }

        return status;
    }

    private static void waitForDCS() {
        String s = System.getProperty(PROP_WAIT_FOR_DCS, "false");
        if (Boolean.parseBoolean(s)) {
            DefaultCacheServer dcs = DefaultCacheServer.getInstance();
            // Wait for service start to ensure that we will get back any partition cache MBeans
            dcs.waitForServiceStart();
        }
    }

    /**
     * Determine whether the Status HA state of the specified service is endangered.
     * <p>
     * If the service only has a single member then it will always be endangered but
     * this method will return {@code false}.
     *
     * @param mBean         the name of the MBean being checked
     * @param attributes    the MBean attributes to use to determine whether the service is HA
     * @return {@code true} if the service is endangered
     */
    private boolean isServiceStatusHA(String mBean, Map<String, Object> attributes) {
        boolean statusHA = true;

        Number nodeCount = (Number) attributes.get(ATTRIB_NODE_COUNT);
        Number backupCount = (Number) attributes.get(ATTRIB_BACKUPS);
        Object status = attributes.get(ATTRIB_HASTATUS);
        Integer remainingDistributionCount = (Integer) attributes.get(ATTRIB_REMAINING_DISTRIBUTION_COUNT);

        if (remainingDistributionCount > 0) {
            // still re-distributing
            statusHA = false;
        }
        else if (nodeCount != null && nodeCount.intValue() > 1 && backupCount != null && backupCount.intValue() > 0) {
            // more than one node with backup > 1 check status is not endangered
            statusHA = !Objects.equals(STATUS_ENDANGERED, status);
        }

        if (!statusHA) {
            log("CoherenceOperator: StatusHA check failed for MBean %s - %s", mBean, attributes);
        }
        else {
            logDebug("CoherenceOperator: StatusHA check passed for MBean %s - %s", mBean, attributes);
        }
        return statusHA;
    }

    /**
     * Obtain the {@link MBeanServerProxy} to use to query Coherence MBeans.
     *
     * @return the {@link MBeanServerProxy} to use to query Coherence MBeans
     */
    private Optional<MBeanServerProxy> getMBeanServerProxy() {
        Cluster cluster = clusterSupplier.get();
        if (cluster != null && cluster.isRunning()) {
            Registry registry = cluster.getManagement();
            if (registry != null) {
                return Optional.ofNullable(registry.getMBeanServerProxy());
            }
        }
        return Optional.empty();
    }

    private Set<String> getPartitionAssignmentMBeans() {
        return getMBeanServerProxy()
                .map(p -> p.queryNames(MBEAN_PARTITION_ASSIGNMENT, null))
                .orElse(Collections.emptySet());
    }

    private Set<String> getPersistenceCoordinatorMBeans() {
        return getMBeanServerProxy()
                .map(p -> p.queryNames(MBEAN_PERSISTENCE_COORDINATOR, null))
                .orElse(Collections.emptySet());
    }

    private void logDebug(String message, Object... args) {
        logDebug(String.format(message, args));
    }

    private void logDebug(String message) {
        if (LOGGING_ENABLED) {
            CacheFactory.log(message, CacheFactory.LOG_DEBUG);
        }
    }

    private void log(String message, Object... args) {
        log(String.format(message, args));
    }

    private void log(String message) {
        CacheFactory.log(message, CacheFactory.LOG_INFO);
    }

    private void warn(String message, Object... args) {
        warn(String.format(message, args));
    }

    private void warn(String message) {
        CacheFactory.log(message, CacheFactory.LOG_WARN);
    }

    private void err(String message, Object... args) {
        err(String.format(message, args));
    }

    private void err(String message) {
        CacheFactory.log(message, CacheFactory.LOG_ERR);
    }

    /**
     * Return attribute and values needed to compute service HA Status.
     *
     * @param mBeanName the MBean name
     * @return attribute/value pairs of the specified MBean needed to compute service HAStatus.
     */
    private Map<String, Object> getMBeanServiceStatusHAAttributes(String mBeanName) {
        return getMBeanAttributes(mBeanName, SERVICE_STATUS_HA_ATTRIBUTES);
    }

    private Map<String, Object> getMBeanAttributes(String sMBean, String[] asAttributes) {
        Map<String, Object> mapAttrValue = new HashMap<>();
        Optional<MBeanServerProxy> optional = getMBeanServerProxy();
        if (optional.isPresent()) {
            MBeanServerProxy proxy = optional.get();
            for (String attribute : asAttributes) {
                mapAttrValue.put(attribute.toLowerCase(), proxy.getAttribute(sMBean, attribute));
            }
        }

        return mapAttrValue;
    }
}
