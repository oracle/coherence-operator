/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;

import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * A readiness/liveness probe.
 *
 * @author jk
 */
public abstract class Probe
        implements AutoCloseable
    {
    /**
     * Perform the readiness test.
     *
     * @return  {@code true} if the readiness test passes.
     */
    public boolean isReady()
        {
        return isClusterMember() && isStatusHA();
        }

    /**
     * Perform the liveness test.
     *
     * @return  {@code true} if the liveness test passes.
     */
    public boolean isLive()
        {
        return isClusterMember();
        }

    @Override
    public void close()
        {
        }

    /**
     * Determine whether the services are all HA.
     *
     * @return  {@code true} if the status is HA
     */
    public boolean isStatusHA()
        {
        boolean fStatusHA = false;

        try
            {
            fStatusHA = getPartitionAssignmentMBeans()
                                .stream()
                                .map(this::getMBeanServiceStatusHAAttributes)
                                .filter(Objects::nonNull)
                                .allMatch(this::isServiceStatusHA);
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
    protected boolean isServiceStatusHA(Map<String, Object> mapAttributes)
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
     * Determine whether this probe can be used.
     *
     * @return  {@code true} if this probe can be used.
     */
    protected abstract boolean isAvailable();

    /**
     * Determine whether the Coherence server is a cluster member.
     *
     * @return  {@code true} if the server is a cluster member
     */
    protected abstract boolean isClusterMember();

    /**
     * Obtain the set of partition assignment MBeans.
     * 
     * @return  the set of partition assignment MBeans
     */
    protected abstract Set<String> getPartitionAssignmentMBeans();

    /**
     * Obtain the values of several attributes of a named MBean.
     *
     * @param sMBean        the MBean name
     * @param asAttributes  the list of attributes to get values for
     *
     * @return the specified attribute/value pairs of the specified MBean
     */
    protected abstract Map<String, Object> getMBeanAttributes(String sMBean, String[] asAttributes);

    /**
     * Return attribute and values needed to compute service HA Status.
     *
     * @param sMBean the MBean name
     *
     * @return attribute/value pairs of the specified MBean needed to compute service HAStatus.
     */
    protected Map<String, Object> getMBeanServiceStatusHAAttributes(String sMBean)
        {
        return getMBeanAttributes(sMBean, SERVICE_STATUS_HA_ATTRIBUTES);
        }

    // ----- constants ------------------------------------------------------

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
    }
