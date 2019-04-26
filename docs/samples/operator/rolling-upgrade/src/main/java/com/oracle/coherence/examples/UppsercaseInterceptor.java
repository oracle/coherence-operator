package com.oracle.coherence.examples;

import com.tangosol.net.events.EventInterceptor;
import com.tangosol.net.events.annotation.Interceptor;
import com.tangosol.net.events.partition.cache.EntryEvent;
import com.tangosol.util.BinaryEntry;

import java.io.Serializable;

import static com.tangosol.net.events.partition.cache.EntryEvent.Type.INSERTING;
import static com.tangosol.net.events.partition.cache.EntryEvent.Type.UPDATING;

@Interceptor(identifier = "Mutator", entryEvents = {INSERTING, UPDATING})
public class UppsercaseInterceptor
        implements EventInterceptor<EntryEvent<?, ?>>, Serializable
    {

    // ----- constructors ---------------------------------------------------

    /**
     * Construct a MutatingInterceptor that will register for all mutable events.
     */
    public UppsercaseInterceptor()
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
