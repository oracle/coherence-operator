/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;

import java.util.function.Supplier;

/**
 * The health MBean implementation.
 *
 * @author jk
 */
public class Health
        implements HealthMBean
    {
    // ----- constructors ---------------------------------------------------

    Health()
        {
        this(CacheFactory::ensureCluster);
        }

    Health(Supplier<Cluster> supplier)
        {
        f_supplier = supplier;
        f_probe    = new ClusterMemberProbe(supplier);
        }

    // ----- HealthMBean methods --------------------------------------------

    @Override
    public boolean ready()
        {
        return hasClusterMembers() && statusHA();
        }

    @Override
    public boolean statusHA()
        {
        return f_probe.isStatusHA();
        }

    @Override
    public boolean healthy()
        {
        return hasClusterMembers();
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Determine whether there are any members in the cluster.
     *
     * @return  {@code true} if the Coherence cluster has members
     */
    private boolean hasClusterMembers()
        {
        return !f_supplier.get().getMemberSet().isEmpty();
        }

    // ----- data members ---------------------------------------------------

    /**
     * The {@link Cluster} supplier.
     */
    private final Supplier<Cluster> f_supplier;

    /**
     * The probe to use to obtain StatusHA.
     */
    private final ClusterMemberProbe f_probe;
    }
