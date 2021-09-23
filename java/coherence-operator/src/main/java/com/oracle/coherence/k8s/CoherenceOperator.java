/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.PrintStream;
import java.net.URL;
import java.util.Properties;

import com.tangosol.net.CacheFactory;

/**
 * An MBean for use by the Operator.
 *
 * @author Jonathan Knight  2020.08.14
 */
public class CoherenceOperator
        implements CoherenceOperatorMBean {

    private static final String NA = "n/a";

    private static Properties properties;

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

    /**
     * Returns the operator version.
     *
     * @return the operator version
     */
    public static String getVersion() {
        return ensureProperties().getProperty("version", NA);
    }

    /**
     * Print the Operator banner.
     *
     * @param out  the {@link PrintStream} to print he banner on
     */
    public static void printBanner(PrintStream out) {
        out.printf("CoherenceOperator: Java Runner version %s\n", getVersion());
    }

    private static synchronized Properties ensureProperties() {
        if (properties == null) {
            Properties props = new Properties();
            try {
                URL url = CoherenceOperator.class.getResource("/META-INF/operator.properties");
                if (url != null) {
                    props.load(url.openStream());
                }
            }
            catch (Throwable t) {
                t.printStackTrace();
            }
            properties = props;
        }
        return properties;
    }
}
