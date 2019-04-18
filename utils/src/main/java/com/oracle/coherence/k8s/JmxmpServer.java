/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import com.oracle.common.base.Blocking;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.management.MBeanServerFinder;
import com.tangosol.util.Base;

import javax.management.MBeanServer;
import javax.management.remote.JMXConnectorServer;
import javax.management.remote.JMXConnectorServerFactory;
import javax.management.remote.JMXServiceURL;
import java.io.IOException;
import java.lang.management.ManagementFactory;

/**
 * An implementation of a Coherence {@link MBeanServerFinder}
 * that creates a {@link JMXConnectorServer} server that uses
 * JMXMP as its transport rather than RMI. This allows JMX
 * to be visible from inside a container that uses NAT'ing.
 *
 * @author jk
 */
public class JmxmpServer
        implements MBeanServerFinder
    {
    // ----- constructors ---------------------------------------------------

    /**
     * Create a {@link JmxmpServer} that binds to the any local address.
     */
    public JmxmpServer()
        {
        this("0.0.0.0");
        }

    /**
     * Create a {@link JmxmpServer} that binds to the specified address.
     *
     * @param sAddress the address to listen on
     */
    public JmxmpServer(String sAddress)
        {
        this.f_sAddress = sAddress;
        }

    // ----- MBeanServerFinder methods --------------------------------------

    @Override
    public MBeanServer findMBeanServer(String s)
        {
        return ensureServer(f_sAddress).getMBeanServer();
        }

    @Override
    public JMXServiceURL findJMXServiceUrl(String s)
        {
        return s_jmxServiceURL;
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Obtain the JMXMP protocol {@link JMXConnectorServer} instance, creating the instance of the connector server if
     * one does not already exist.
     *
     * @param address the address to listen on
     *
     * @return the JMXMP protocol {@link JMXConnectorServer} instance.
     */
    private static synchronized JMXConnectorServer ensureServer(String address)
        {
        try
            {
            if (s_connectorServer == null)
                {
                MBeanServer server = ManagementFactory.getPlatformMBeanServer();
                int nPort = Integer.getInteger("coherence.jmxmp.port", 9000);

                s_jmxServiceURL = new JMXServiceURL("jmxmp", address, nPort);
                s_connectorServer = JMXConnectorServerFactory.newJMXConnectorServer(s_jmxServiceURL, null, server);

                s_connectorServer.start();

                CacheFactory.log("Started JMXMP connector " + s_connectorServer.getAddress());
                }

            return s_connectorServer;
            }
        catch (IOException e)
            {
            throw Base.ensureRuntimeException(e);
            }
        }

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     *
     * @throws Exception if there is a program error
     */
    public static void main(String[] args)
        {
        Cluster cluster = CacheFactory.ensureCluster();

        while (cluster.isRunning())
            {
            try
                {
                Blocking.sleep(1000);
                }
            catch (InterruptedException e)
                {
                // we don't know what effect setting the interrupt flag will
                // have on the stop() methods below; since this is the main
                // thread it would be safer not to set the flag
                // Thread.currentThread().interrupt();
                break;
                }
            }
        }

    // ----- data members ---------------------------------------------------

    /**
     * The JMXServiceURL for the MBeanConnector used by the Coherence JMX framework.
     */
    private static JMXServiceURL s_jmxServiceURL;

    /**
     * The {@link JMXConnectorServer} using the JMXMP protocol.
     */
    private static JMXConnectorServer s_connectorServer;

    private final String f_sAddress;
    }
