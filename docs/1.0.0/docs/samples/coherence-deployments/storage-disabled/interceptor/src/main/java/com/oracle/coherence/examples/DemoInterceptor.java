/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.NamedCache;
import com.tangosol.net.events.EventInterceptor;
import com.tangosol.net.events.application.LifecycleEvent;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * Storage-disabled client (as {@link EventInterceptor}) to insert data into cluster.
 *
 * @author tam 2019.05.21
 */
public class DemoInterceptor implements EventInterceptor<LifecycleEvent>
{
    @Override
    public void onEvent(LifecycleEvent event)
    {
        if (event.getType() == LifecycleEvent.Type.ACTIVATED)
        {
        // This Runnable simulates a storage disabled client    
        Runnable runnableTask = () ->
            {
            int                         i  = 0;
            NamedCache<Integer, String> nc = CacheFactory.getCache("interceptor-cache");

            while (true)
                {
                String value = format.format(new Date(System.currentTimeMillis()));
                nc.put(i, value);
                CacheFactory.log("Inserted key=" + i++ + ", value=" +value);
                try
                    {
                    Thread.sleep(1000L);
                    }
                catch (InterruptedException e)
                    {
                    }
                }
            };

        executor.submit(runnableTask);
        }
    }

    private ExecutorService  executor = Executors.newFixedThreadPool(10);
    private SimpleDateFormat format   = new SimpleDateFormat("HH:mm:ss");
}
