/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.net.events.EventInterceptor;
import com.tangosol.net.events.annotation.Interceptor;
import com.tangosol.net.events.partition.cache.EntryEvent;
import com.tangosol.util.BinaryEntry;

import java.io.Serializable;

import static com.tangosol.net.events.partition.cache.EntryEvent.Type.INSERTING;
import static com.tangosol.net.events.partition.cache.EntryEvent.Type.UPDATING;

/**
 * Interceptor that converts all data cache values to uppercase.
 *
 * @author tam  2019.04.29
 */
@Interceptor(identifier = "Mutator", entryEvents = {INSERTING, UPDATING})
public class UppercaseInterceptor
        implements EventInterceptor<EntryEvent<?, ?>>, Serializable
    {

    // ----- constructors ---------------------------------------------------

    /**
     * Construct a MutatingInterceptor that will register for all mutable events.
     */
    public UppercaseInterceptor()
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
                entry.setValue(sValue.toUpperCase());
                }
            }
        }
    }
