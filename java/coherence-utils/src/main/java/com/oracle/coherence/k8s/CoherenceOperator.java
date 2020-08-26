/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;

/**
 * An MBean for use by the Operator.
 *
 * @author Jonathan Knight  2020.08.14
 */
public class CoherenceOperator
        implements CoherenceOperatorMBean {

    private static final String NA = "n/a";

    private String identity = NA;

    /**
     * Create a CoherenceOperator MBean.
     */
    public CoherenceOperator() {
        String id = System.getProperty(PROP_IDENTITY, NA);
        if (!id.isEmpty()) {
            this.identity = id;
        }
    }

    @Override
    public String getIdentity() {
        return identity;
    }

    @Override
    public int getNodeId() {
        return CacheFactory.getCluster().getLocalMember().getId();
    }
}
