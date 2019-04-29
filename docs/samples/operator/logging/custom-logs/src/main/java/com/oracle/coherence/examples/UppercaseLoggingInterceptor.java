/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.Cluster;
import com.tangosol.net.events.EventInterceptor;
import com.tangosol.net.events.annotation.Interceptor;
import com.tangosol.net.events.partition.cache.EntryEvent;
import com.tangosol.util.BinaryEntry;

import java.io.Serializable;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.logging.Level;
import java.util.logging.Logger;

import static com.tangosol.net.events.partition.cache.EntryEvent.Type.INSERTING;
import static com.tangosol.net.events.partition.cache.EntryEvent.Type.UPDATING;

/**
 * Interceptor that converts all data cache values to uppercase and uses
 * the custom 'sample' logger to log messages.
 *
 * @author tam  2019.04.29
 */
@Interceptor(identifier = "Mutator", entryEvents = {INSERTING, UPDATING})
public class UppercaseLoggingInterceptor
        implements EventInterceptor<EntryEvent<?, ?>>, Serializable
    {

    // ----- constructors ---------------------------------------------------

    /**
     * Construct a MutatingInterceptor that will register for all mutable events.
     */
    public UppercaseLoggingInterceptor()
        {
        super();
        }

    // ----- EventInterceptor methods ---------------------------------------

    /**
     * {@inheritDoc}
     */
    public void onEvent(EntryEvent<?, ?> entryEvent)
        {
        for (BinaryEntry entry : entryEvent.getEntrySet())
            {
            String sValue = (String) entry.getValue();
            if (entryEvent.getType() == INSERTING || entryEvent.getType() == UPDATING)
                {
                log(Level.INFO, "Before, key=" + entry.getKey() + ", value=" + entry.getValue());
                entry.setValue(sValue.toUpperCase());
                log(Level.INFO,"Changed key=" + entry.getKey() + " to value=" + entry.getValue());
                }
            }
        }

    // ----- static ---------------------------------------------------------

    static void log(Level level, String msg)
        {
        Cluster cluster = CacheFactory.getCluster();
        String  member  = cluster.getLocalMember().getMemberName();
        s_logger.info(getLogTimestamp() + " Cloud 1.0" + " <" + level.getName() + "> " +
            "(cluster=" + cluster.getClusterName() + ", member=" + member + ", thread=" + Thread.currentThread().getName() + "): " +
            msg);
        }

    static private String getLogTimestamp()
        {
        SimpleDateFormat f = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        return f.format(new Date(System.currentTimeMillis()));
        }

    /**
     * Custom logger.
     */
    private final static Logger s_logger = Logger.getLogger("cloud");
    }
