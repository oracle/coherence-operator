/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package com.oracle.coherence.k8s;

import java.io.IOException;
import java.lang.management.ManagementFactory;

import javax.management.MBeanServer;
import javax.management.remote.JMXConnectorServer;
import javax.management.remote.JMXConnectorServerFactory;
import javax.management.remote.JMXServiceURL;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.management.MBeanServerFinder;
import com.tangosol.util.Base;

/**
 * An implementation of a Coherence {@link MBeanServerFinder}
 * that creates a {@link JMXConnectorServer} server that uses
 * JMXMP as its transport rather than RMI. This allows JMX
 * to be visible from inside a container that uses NAT'ing.
 */
public class JmxmpServer
        implements MBeanServerFinder {
    // ----- data members ---------------------------------------------------

    /**
     * The JMXServiceURL for the MBeanConnector used by the Coherence JMX framework.
     */
    private static JMXServiceURL jmxServiceURL;

    /**
     * The {@link JMXConnectorServer} using the JMXMP protocol.
     */
    private static JMXConnectorServer connectorServer;

    private final String address;

    // ----- constructors ---------------------------------------------------

    /**
     * Create a {@link JmxmpServer} that binds to the any local address.
     */
    public JmxmpServer() {
        this("0.0.0.0");
    }

    /**
     * Create a {@link JmxmpServer} that binds to the specified address.
     *
     * @param address the address to listen on
     */
    public JmxmpServer(String address) {
        this.address = address;
    }

    // ----- MBeanServerFinder methods --------------------------------------

    @Override
    public MBeanServer findMBeanServer(String s) {
        return ensureServer(address).getMBeanServer();
    }

    @Override
    public JMXServiceURL findJMXServiceUrl(String s) {
        return jmxServiceURL;
    }

    // ----- helper methods -------------------------------------------------

    /**
     * Obtain the JMXMP protocol {@link JMXConnectorServer} instance, creating the instance of the connector server if
     * one does not already exist.
     *
     * @param address the address to listen on
     * @return the JMXMP protocol {@link JMXConnectorServer} instance.
     */
    private static synchronized JMXConnectorServer ensureServer(String address) {
        try {
            if (connectorServer == null) {
                MBeanServer server = ManagementFactory.getPlatformMBeanServer();
                int port = Integer.getInteger("coherence.jmxmp.port", 9000);

                jmxServiceURL = new JMXServiceURL("jmxmp", address, port);
                connectorServer = JMXConnectorServerFactory.newJMXConnectorServer(jmxServiceURL, null, server);

                connectorServer.start();

                CacheFactory.log("Started JMXMP connector " + connectorServer.getAddress());
            }

            return connectorServer;
        }
        catch (IOException e) {
            throw Base.ensureRuntimeException(e);
        }
    }

    /**
     * Program entry point.
     *
     * @param args the program command line arguments
     */
    public static void main(String[] args) {
        Cluster cluster = CacheFactory.ensureCluster();

        while (cluster.isRunning()) {
            try {
                Thread.sleep(1000);
            }
            catch (InterruptedException e) {
                // we don't know what effect setting the interrupt flag will
                // have on the stop() methods below; since this is the main
                // thread it would be safer not to set the flag
                // Thread.currentThread().interrupt();
                break;
            }
        }
    }
}
