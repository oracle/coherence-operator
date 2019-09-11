/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.sun.tools.attach.VirtualMachine;
import com.sun.tools.attach.VirtualMachineDescriptor;

import javax.management.MBeanServerConnection;
import javax.management.MalformedObjectNameException;
import javax.management.ObjectName;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;

import java.io.IOException;
import java.util.List;
import java.util.Properties;

/**
 * A {@link Probe} that connects using local JMX to a JVM
 * to perform its checks.
 *
 * @author jk
 */
public class OperatorMBeanProbe
        implements Probe
    {
    // ----- Probe methods --------------------------------------------------

    @Override
    public boolean isAvailable()
        {
        return hasOperatorMBean();
        }

    @Override
    public boolean isReady()
        {
        try
            {
            Object o = invokeMBeanMethod("ready", NO_ARGS, NO_SIG);
            return o instanceof Boolean && (Boolean) o;
            }
        catch (Exception e)
            {
            return false;
            }
        }

    @Override
    public boolean isLive()
        {
        try
            {
            Object o = invokeMBeanMethod("healthy", NO_ARGS, NO_SIG);
            return o instanceof Boolean && (Boolean) o;
            }
        catch (Exception e)
            {
            return false;
            }
        }

    @Override
    public boolean isStatusHA()
        {
        try
            {
            Object o = invokeMBeanMethod("statusHA", NO_ARGS, NO_SIG);
            return o instanceof Boolean && (Boolean) o;
            }
        catch (Exception e)
            {
            return false;
            }
        }

    @Override
    public synchronized void close()
        {
        if (m_connector != null)
            {
            try
                {
                m_connector.close();
                }
            catch (IOException e)
                {
                e.printStackTrace();
                }
            m_connector = null;
            }
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Obtain a connection to the server JVM.
     *
     * @return the {@link JMXConnector} to use or null if no connection is possible
     */
    public synchronized JMXConnector ensureConnection()
        {
        if (m_connector == null)
            {
            try
                {
                List<VirtualMachineDescriptor> vms = VirtualMachine.list();
                for (VirtualMachineDescriptor desc : vms)
                    {
                    if (!desc.displayName().contains(Main.class.getCanonicalName()))
                        {
                        continue;
                        }

                    VirtualMachine vm;
                    try
                        {
                        vm = VirtualMachine.attach(desc);
                        }
                    catch (Throwable e)
                        {
                        continue;
                        }

                    try
                        {
                        Properties props    = vm.getAgentProperties();
                        String     sAddress = props.getProperty(PROP_ADDRESS);

                        if (sAddress != null)
                            {
                            m_connector = JMXConnectorFactory.connect(new JMXServiceURL(sAddress));
                            break;
                            }
                        }
                    catch (Exception e)
                        {
                        e.printStackTrace();
                        }
                    }
                }
            catch (Throwable t)
                {
                t.printStackTrace();
                }
            }

        return m_connector;
        }

    /**
     * Determine whether the health MBean is present on the JVM.
     *
     * @return {@code true} if the health MBean is present
     */
    boolean hasOperatorMBean()
        {
        try
            {
            JMXConnector connector = ensureConnection();
            if (connector != null)
                {
                MBeanServerConnection serverConnection  = connector.getMBeanServerConnection();
                return serverConnection.isRegistered(new ObjectName(Main.HealthObjectName));
                }
            else
                {
                System.err.println("Error checking for Operator MBean JMXConnector is null");
                return false;
                }
            }
        catch (IOException | MalformedObjectNameException e)
            {
            System.err.println("Error checking for Operator MBean " + Main.HealthObjectName);
            e.printStackTrace();
            return false;
            }
        }

    /**
     * Invoke a method on the health MBean.
     *
     * @param sMethodName  the method name
     * @param aoParam      the method parameters
     * @param asSig        the method signature
     *
     * @return             the method return value
     *
     * @throws Exception if the method call fails
     */
    Object invokeMBeanMethod(String sMethodName, Object[] aoParam, String[] asSig) throws Exception
        {
        JMXConnector connector = ensureConnection();
        if (connector != null)
            {
            MBeanServerConnection serverConnection  = connector.getMBeanServerConnection();
            return serverConnection.invoke(new ObjectName(Main.HealthObjectName), sMethodName, aoParam, asSig);
            }
        else
            {
            System.err.println("Error invoking Operator MBean method " + sMethodName + " JMXConnector is null");
            return null;
            }
        }

    // ----- constants ------------------------------------------------------

    private static final String PROP_ADDRESS = "com.sun.management.jmxremote.localConnectorAddress";

    private static final Object[] NO_ARGS = new Object[0];

    private static final String[] NO_SIG = new String[0];

    // ----- data members ---------------------------------------------------

    /**
     * The current {@link JMXConnector} instance to use.
     */
    private JMXConnector m_connector;
    }
