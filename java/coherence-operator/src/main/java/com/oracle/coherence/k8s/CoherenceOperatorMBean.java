/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

/**
 * An MBean for use by the Operator.
 *
 * @author Jonathan Knight  2020.08.14
 */
public interface CoherenceOperatorMBean {
    /**
     * The ObjectName for the MBean.
     */
    String OBJECT_NAME = "type=KubernetesOperator";

    /**
     * The System property that is used to set this members identity.
     */
    String PROP_IDENTITY = "coherence.operator.identity";

    /**
     * The name of the Identity MBean attribute.
     */
    String ATTRIBUTE_IDENTITY = "Identity";

    /**
     * The name of the Node ID MBean attribute.
     */
    String ATTRIBUTE_NODE = "NodeId";

    /**
     * Get the identity of the deployment.
     * <p>
     * This will typically be made up of the deployment name and namespace.
     *
     * @return the identity of the deployment
     */
    String getIdentity();

    /**
     * Get the Coherence node id.
     *
     * @return the Coherence node id
     */
    int getNodeId();
}
