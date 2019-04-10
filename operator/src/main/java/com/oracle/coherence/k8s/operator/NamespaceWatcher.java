/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import com.squareup.okhttp.Call;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.apis.CoreV1Api;
import io.kubernetes.client.models.V1Namespace;
import io.kubernetes.client.models.V1ObjectMeta;
import io.kubernetes.client.util.Watch;

import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.Consumer;

/**
 * This class drives the watch support for Namespace.
 * It runs in a separate thread to drive watching asynchronously to
 * the main thread.
 *
 * @author sc
 */
class NamespaceWatcher extends AbstractWatcher<V1Namespace, CoreV1Api>
    {
    // ----- constructors ----------------------------------------------------

    /**
     * Constructs a namespace watcher.
     *
     * @param fStopping   an AtomicBoolean to determine when to stop the watcher
     * @param consumer    a consumer to consume the watch events
     *
     * @exception IllegalAccessException  if the class or its nullary constructor is not accessible
     * @exception InstantiationException  if this Class represents an abstract class, an interface, an array class,
     *     a primitive type, or void; or if the class has no nullary constructor; or if the instantiation fails for
     *     some other reason
     */
    NamespaceWatcher(AtomicBoolean fStopping,
                     Consumer<Watch.Response<V1Namespace>> consumer)
            throws IllegalAccessException, InstantiationException
        {
        super(null, fStopping, consumer, V1Namespace.class, CoreV1Api.class);
        }

    // ---- methods ----------------------------------------------------------

    @Override
    Call createCall(CallParams callParams) throws ApiException
        {
        return m_api.listNamespaceCall(callParams.getPretty(), null, callParams.getFieldSelector(),
                callParams.getIncludeUninitialized(), callParams.getLabelSelector(), callParams.getLimit(),
                callParams.getResourceVersion(), callParams.getTimeoutSeconds(),
                Boolean.TRUE, callParams.getProgressListener(), callParams.getProgressRequestListener());
        }

    @Override
    V1ObjectMeta getMetadata(V1Namespace obj)
        {
        return obj.getMetadata();
        }
    }
