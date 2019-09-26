package com.oracle.coherence.k8s;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.management.MBeanServerProxy;
import com.tangosol.net.management.Registry;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.function.Supplier;
import java.util.stream.Collectors;

/**
 * Simple http endpoint for heath checking.
 *
 * @author jk
 */
public class HealthServer
    {
    // ----- constructors ---------------------------------------------------

    HealthServer()
        {
        this(CacheFactory::getCluster);
        }

    HealthServer(Supplier<Cluster> supplier)
        {
        f_supplier = supplier;
        }

    // ----- HealthServer methods ------------------------------------------------

    /**
     * Start a http server.
     *
     * @throws IOException  if an error occurs
     */
    public synchronized void start() throws IOException
        {
        if(m_server == null)
            {
            int nPort = Integer.getInteger(PROP_HEALTH_PORT, 1002);
            HttpServer server = HttpServer.create(new InetSocketAddress(nPort), 0);

            server.createContext(PATH_READY, this::ready);
            server.createContext(PATH_HEALTH, this::health);
            server.createContext(PATH_HA, this::statusHA);

            server.setExecutor(null); // creates a default executor
            server.start();

            System.out.println("ReST server is UP! http://localhost:" + server.getAddress().getPort());

            m_server = server;
            }
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Send a http response.
     * @param t       the {@link HttpExchange} to send the response to
     * @param status  the response status
     */
    private static void send(HttpExchange t, int status)
        {
        try
            {
            t.sendResponseHeaders(status, 0);
            OutputStream os = t.getResponseBody();
            os.write(new byte[0]);
            os.close();
            }
        catch (IOException e)
            {
            e.printStackTrace();
            }
        }

    /**
     * Process a ready request.
     * @param t  the {@link HttpExchange} to send the response to
     */
    private void ready(HttpExchange t)
        {
        int nResponse = hasClusterMembers() && isStatusHA() ? 200 : 400;
        send(t, nResponse);
        }

    /**
     * Process a health request.
     * @param t  the {@link HttpExchange} to send the response to
     */
    private void health(HttpExchange t)
        {
        int nResponse = hasClusterMembers() ? 200 : 400;
        send(t, nResponse);
        }

    /**
     * Process a status HA request.
     * @param t  the {@link HttpExchange} to send the response to
     */
    private void statusHA(HttpExchange t)
        {
        int nResponse = isStatusHA() ? 200 : 400;
        send(t, nResponse);
        }

    /**
     * Determine whether there are any members in the cluster.
     *
     * @return  {@code true} if the Coherence cluster has members
     */
    private boolean hasClusterMembers()
        {
        Cluster cluster = f_supplier.get();
        return cluster != null && cluster.isRunning() && !cluster.getMemberSet().isEmpty();
        }

    /**
     * Returns {@code true} if the JVM is StatusHA.
     *
     * @return {@code true} if the JVM is StatusHA
     */
    boolean isStatusHA()
        {
        boolean fStatusHA = false;

        try
            {
            Cluster cluster = f_supplier.get();

            if (cluster != null && cluster.isRunning())
                {
                fStatusHA = getPartitionAssignmentMBeans()
                                    .stream()
                                    .map(this::getMBeanServiceStatusHAAttributes)
                                    .filter(Objects::nonNull)
                                    .allMatch(this::isServiceStatusHA);
                }
            }
        catch (Throwable t)
            {
            CacheFactory.log(t);
            }

        return fStatusHA;
        }

    /**
     * Determine whether the Status HA state of the specified service is endangered.
     * <p>
     * If the service only has a single member then it will always be endangered but
     * this method will return {@code false}.
     *
     * @param mapAttributes  the MBean attributes to use to determine whether the service is HA
     *
     * @return  {@code true} if the service is endangered
     */
    private boolean isServiceStatusHA(Map<String, Object> mapAttributes)
        {
        boolean fStatusHA = true;

        // convert the attribute case as MBeanProxy or ReST return them with different cases
        Map map = mapAttributes.entrySet()
                        .stream()
                        .collect(Collectors.toMap(e -> e.getKey().toLowerCase(), Map.Entry::getValue));

        Number cNode   = (Number) map.get(ATTRIB_NODE_COUNT);
        Number cBackup = (Number) map.get(ATTRIB_BACKUPS);


        if (cNode != null && cNode.intValue() > 1 && cBackup != null && cBackup.intValue() > 0)
            {
            fStatusHA = !Objects.equals(STATUS_ENDANGERED, map.get(ATTRIB_HASTATUS));
            }

        return fStatusHA;
        }

    /**
     * Obtain the {@link MBeanServerProxy} to use to query Coherence MBeans.
     *
     * @return  the {@link MBeanServerProxy} to use to query Coherence MBeans
     */
    private MBeanServerProxy getMBeanServerProxy()
        {
        Cluster  cluster  = f_supplier.get();
        Registry registry = cluster.getManagement();

        return registry.getMBeanServerProxy();
        }

    private Set<String> getPartitionAssignmentMBeans()
        {
        return getMBeanServerProxy().queryNames(MBEAN_PARTITION_ASSIGNMENT, null);
        }

    /**
     * Return attribute and values needed to compute service HA Status.
     *
     * @param sMBean the MBean name
     *
     * @return attribute/value pairs of the specified MBean needed to compute service HAStatus.
     */
    private Map<String, Object> getMBeanServiceStatusHAAttributes(String sMBean)
        {
        return getMBeanAttributes(sMBean, SERVICE_STATUS_HA_ATTRIBUTES);
        }

    private Map<String, Object> getMBeanAttributes(String sMBean, String[] asAttributes)
        {
        Map<String, Object> mapAttrValue = new HashMap<>();

        for (String attribute: asAttributes)
            {
            mapAttrValue.put(attribute, getMBeanServerProxy().getAttribute(sMBean, attribute));
            }

        return mapAttrValue;
        }

    // ----- constants ------------------------------------------------------

    /**
     * The system property to use to set the health port.
     */
    public static final String PROP_HEALTH_PORT = "coherence.health.port";

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
     * The MBean name of the PartitionAssignment MBean.
     */
    public static final String MBEAN_PARTITION_ASSIGNMENT = Registry.PARTITION_ASSIGNMENT_TYPE
            + ",service=*,responsibility=DistributionCoordinator";

    /**
     * Service MBean Attributes required to compute HAStatus.
     *
     * @see #isServiceStatusHA(Map)
     */
    public static final String[] SERVICE_STATUS_HA_ATTRIBUTES =
        {
            "HAStatus",
            "BackupCount",
            "ServiceNodeCount"
        };

    /**
     * The value of the Status HA attribute to signify endangered.
     */
    public static final String STATUS_ENDANGERED = "ENDANGERED";

    /**
     * The name of the HA status MBean attribute.
     */
    public static final String ATTRIB_HASTATUS = "hastatus";

    /**
     * The name of the backup count MBean attribute.
     */
    public static final String ATTRIB_BACKUPS = "backupcount";

    /**
     * The name of the service node count MBean attribute.
     */
    public static final String ATTRIB_NODE_COUNT = "servicenodecount";

    // ----- data members ---------------------------------------------------

    /**
     * The http server.
     */
    private HttpServer m_server;

    /**
     * The {@link Cluster} supplier.
     */
    private final Supplier<Cluster> f_supplier;
    }
