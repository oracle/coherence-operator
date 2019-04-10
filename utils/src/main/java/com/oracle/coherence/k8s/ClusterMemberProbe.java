/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.Member;
import com.tangosol.net.management.MBeanServerProxy;
import com.tangosol.net.management.Registry;
import com.tangosol.util.filter.AlwaysFilter;

import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.function.Supplier;

/**
 * A readiness/liveness probe that joins the Coherence cluster
 * and uses the {@link MBeanServerProxy}.
 *
 * @author jk
 */
public class ClusterMemberProbe
        extends Probe
    {
    // ----- constructors ---------------------------------------------------

    /**
     * Create a {@link ClusterMemberProbe}.
     */
    public ClusterMemberProbe()
        {
        this(CacheFactory::ensureCluster);
        }

    /**
     * Create a {@link ClusterMemberProbe}.
     *
     * @param supplier  the Coherence {@link Cluster} {@link Supplier}
     */
    ClusterMemberProbe(Supplier<Cluster> supplier)
        {
        f_supplierCluster = supplier;
        }

    // ----- Probe methods --------------------------------------------------

    @Override
    public boolean isAvailable()
        {
        return true;
        }

    @Override
    @SuppressWarnings("unchecked")
    public boolean isClusterMember()
        {
        CacheFactory.log(getClass().getSimpleName() + " isClusterMember()");

        Cluster     cluster       = f_supplierCluster.get();
        Set<Member> setMember     = cluster.getMemberSet();
        Member      memberLocal   = cluster.getLocalMember();
        String      sAddressLocal = memberLocal.getAddress().getHostAddress();
        int         nPIDLocal     = Integer.parseInt(memberLocal.getProcessName());

        // determine whether there is at least one other cluster member on this host with a lower PID than this process
        long cMember = setMember.stream()
                                .filter(m -> m.getAddress().getHostAddress().equals(sAddressLocal))
                                .filter(m -> Integer.parseInt(m.getProcessName()) < nPIDLocal)
                                .count();

        return cMember >= 1;
        }

    @Override
    protected Set<String> getPartitionAssignmentMBeans()
        {
        return getMBeanServerProxy().queryNames(MBEAN_PARTITION_ASSIGNMENT, null);
        }

    @Override
    protected Map<String, Object> getMBeanAttributes(String sMBean, String[] asAttributes)
        {
        Map<String, Object> mapAttrValue = new HashMap<>();

        for (String attribute: asAttributes)
            {
            mapAttrValue.put(attribute, getMBeanServerProxy().getAttribute(sMBean, attribute));
            }

        return mapAttrValue;
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Obtain the {@link MBeanServerProxy} to use to query Coherence MBeans.
     *
     * @return  the {@link MBeanServerProxy} to use to query Coherence MBeans
     */
    private MBeanServerProxy getMBeanServerProxy()
        {
        if (m_proxy == null)
            {
            Cluster  cluster  = f_supplierCluster.get();
            Registry registry = cluster.getManagement();

            m_proxy = registry.getMBeanServerProxy();

            m_proxy.isMBeanRegistered("foo");
            }

        return m_proxy;
        }

    // ----- constants ------------------------------------------------------

    /**
     * The MBean name of the PartitionAssignment MBean.
     */
    public static final String MBEAN_PARTITION_ASSIGNMENT = Registry.PARTITION_ASSIGNMENT_TYPE
            + ",service=*,responsibility=DistributionCoordinator";

    // ----- data members ---------------------------------------------------

    /**
     * The Coherence {@link Cluster}.
     */
    private final Supplier<Cluster> f_supplierCluster;

    /**
     * The Coherence MBean server proxy.
     */
    private MBeanServerProxy m_proxy;
    }
