/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import com.squareup.okhttp.Call;
import io.kubernetes.client.ApiClient;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.models.V1ObjectMeta;
import io.kubernetes.client.util.Config;
import io.kubernetes.client.util.Watch;

import java.io.IOException;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.util.concurrent.ThreadFactory;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.Consumer;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * This abstract class drives the watch support for a specific type of object.
 * It runs in a separate thread to drive watching asynchronously to
 * the main thread.
 *
 * @param <T> the type of the object to be watched.
 * @param <A> the type of the Api object to be used.
 *
 * @author sc
 */
abstract class AbstractWatcher<T, A>
    {
    // ----- constructors ----------------------------------------------------

    /**
     * Constructs a watcher.
     *
     * @param sNamespace                namespace where null or empty String means "all namespaces"
     * @param fStopping                 an AtomicBoolean to determine when to stop the watcher
     * @param consumer                  a consumer to consume the watch events
     * @param clzWatchType              the class type of T
     * @param clzApiType                the class type of A
     *
     * @exception IllegalAccessException  if the class or its nullary constructor is not accessible
     * @exception InstantiationException  if this Class represents an abstract class, an interface, an array class,
     *     a primitive type, or void; or if the class has no nullary constructor; or if the instantiation fails for
     *     some other reason
     */
    AbstractWatcher(String sNamespace, AtomicBoolean fStopping,
                    Consumer<Watch.Response<T>> consumer, Class<T> clzWatchType, Class<A> clzApiType)
            throws IllegalAccessException, InstantiationException
        {
        this.m_sNamespace               = sNamespace;
        this.m_fStopping                = fStopping;
        this.m_consumer                 = consumer;
        this.m_clzWatchType             = clzWatchType;

        this.m_api = clzApiType.newInstance();
        }

    // ----- abstract methods ------------------------------------------------

    /**
     * Creates a Call object for a given namespace returned by api.{ListOperation}Call(...) method.
     * Make sure watch flag is set in the call.
     *
     * @return a call object
     * @throws ApiException  Kubernetes API client exception
     */
    abstract Call createCall(CallParams callParams) throws ApiException;

    /**
     * Retrieves meta data of the given object.
     *
     * @param obj
     * @return the meta data of the given object
     */
    abstract V1ObjectMeta getMetadata(T obj);

    // ----- methods ---------------------------------------------------------

    /**
     * Kicks off the watcher processing that runs in a separate thread.
     *
     * @param threadFactory  the threadFactory
     */
    void start(ThreadFactory threadFactory)
        {
        m_thread = threadFactory.newThread(this::doWatch);
        m_thread.start();
        }

    /**
     * Wait for the thread death.
     */
    void waitForDeath()
        {
        try
            {
            if (m_thread != null)
                {
                m_thread.join();
                }
            }
        catch (InterruptedException ignore)
            {
            // ignoring
            }
        }

    /**
     * Checks whether it is null or empty which corresponds to all namespaces.
     *
     * @return <code>true</code> if it is for all namespaces.
     */
    boolean isAllNamespaces(String sNamespace)
        {
        return sNamespace == null || sNamespace.isEmpty();
        }

    // ---- helper methods ---------------------------------------------------

    /**
     * Watches for events until stopping.
     */
    private void doWatch() {
        try
            {
            ApiClient client = Config.defaultClient();

            while (!m_fStopping.get())
                {
                try
                    {
                    CallParams callParams = new CallParams();
                    callParams.setResourceVersion(m_sLastResourceVersion);
                    watchCall(createCall(callParams), client);
                    }
                catch (RuntimeException | ApiException ignore)
                    {
                    LOGGER.finest(() -> "Ignore watcher["+ this.getClass().getSimpleName() + "@" + m_sNamespace +
                            "] fails: " + ignore);
                    }
                }
            }
        catch(RuntimeException | IOException e)
            {
            LOGGER.log(Level.WARNING, e,
                    () -> "Watcher[" + this.getClass().getSimpleName() + "@" + m_sNamespace + "] fails: " + e.getMessage());
            m_fStopping.set(true);
            }
        }

    /**
     * Watches for a call.
     *
     * @param call  call the call object returned by api.{ListOperation}Call(...) method.
     *              Make sure watch flag is set in the call.
     * @param client  the Kubernetes API client
     * @throws ApiException  Kubernetes API client exception
     * @throws IOException  IOException while watching
     */
    private void watchCall(Call call, ApiClient client) throws ApiException, IOException
        {
        Type paramType = getParameterizedType(m_clzWatchType);

        try (Watch<T> watch = Watch.createWatch(client, call, paramType))
            {
            for (Watch.Response<T> item : watch)
                {
                if (item != null)
                    {
                    if (item.object != null)
                        {
                        m_sLastResourceVersion = getMetadata(item.object).getResourceVersion();
                        }
                    m_consumer.accept(item);
                    }
                }
            }
        }

    /**
     * Returns a ParameterizedType of Watch.Response&lt;responseBodyType&gt;
     * @param responseBodyType  the body generic type
     * @return  the parameterized type
     */
    private static Type getParameterizedType(Type responseBodyType)
        {
        return new ParameterizedType()
            {
            @Override
            public Type[] getActualTypeArguments()
                {
                return new Type[] { responseBodyType };
                }

            @Override
            public Type getRawType()
                {
                return Watch.Response.class;
                }

            @Override
            public Type getOwnerType()
                {
                return Watch.class;
                }
            };
        }

    /**
     * Sets api object.
     *
     * @param api  the api object
     */
    void setApi(A api)
        {
        this.m_api = api;
        }

    // ----- data members ---------------------------------------------------

    /**
     * Class Logger.
     */
    static final Logger LOGGER = Logger.getLogger("Operator");

    /**
     * The namespace where null means "all namespaces".
     */
    final String m_sNamespace;

    /**
     * The AtomicBoolean indicates whether the operator is stopping.
     */
    final AtomicBoolean m_fStopping;

    /**
     * The Consumer for Watch.Response event.
     */
    final Consumer<Watch.Response<T>> m_consumer;

    /**
     * The Class type
     */
    final Class<T> m_clzWatchType;

    /**
     * The last resource versions being watched.
     */
    private String m_sLastResourceVersion;

    /**
     * The Thread for running this Watcher.
     */
    private Thread m_thread = null;

    /**
     * The api object to access Kubernetes T info.
     */
    A m_api;
    }
