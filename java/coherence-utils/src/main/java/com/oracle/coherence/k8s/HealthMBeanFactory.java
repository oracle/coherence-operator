/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.util.Base;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;

/**
 * A factory to produce the health check MBean.
 *
 * @author jk
 */
public class HealthMBeanFactory
    {
    public static HealthMBean registerHealthMBean()
        {
        try
            {
            if (s_health == null)
                {
                s_health = new Health();

                MBeanServer server = ManagementFactory.getPlatformMBeanServer();
                server.registerMBean(s_health,  new ObjectName(HealthMBean.HEALTH_OBJECT_NAME));
                }

            return s_health;
            }
        catch (Throwable t)
            {
            throw Base.ensureRuntimeException(t);
            }
        }

    // ----- data members ---------------------------------------------------

    private static Health s_health;
    }
