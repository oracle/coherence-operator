/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import io.kubernetes.client.ProgressRequestBody;
import io.kubernetes.client.ProgressResponseBody;

/**
 * An object which encapsulates common parameters for Kubernetes API calls.
 *
 * @author sc
 */
class CallParams
    {
    /**
     * Returns a boolean indicating whether partially initialized results should be included in the response.
     *
     * @return the current setting of the parameter. Defaults to including everything.
     */
    Boolean getIncludeUninitialized()
        {
        return m_fIncludeUninitialized;
        }

    /**
     * Sets includeUninitialized in Boolean value.
     *
     * @param fIncludeUninitialized  the Boolean includeUninitialized
     */
    void setIncludeUninitialized(Boolean fIncludeUninitialized)
        {
        this.m_fIncludeUninitialized = fIncludeUninitialized;
        }

    /**
     * Returns the limit on the number of updates to send in a single reply.
     *
     * @return the current setting of the parameter. Defaults to 500.
     */
    Integer getLimit()
        {
        return m_nLimit;
        }

    /**
     * Sets limit in Integer.
     *
     * @param nLimit  the number of updates to send
     */
    void setLimit(Integer nLimit)
        {
        this.m_nLimit = nLimit;
        }

    /**
     * Returns the timeout in seconds for the call.
     *
     * @return the current setting. Defaults to 30 seconds.
     */
    Integer getTimeoutSeconds()
        {
        return m_nTimeoutSeconds;
        }

    /**
     * Sets timeout in seconds in Integer.
     *
     * @param nTimeoutSeconds  the timeout for the call.
     */
    void setTimeoutSeconds(Integer nTimeoutSeconds)
        {
        this.m_nTimeoutSeconds = nTimeoutSeconds;
        }

    /**
     * Returns a selector to limit results to those with matching fields.
     *
     * @return the option, if specified. Defaults to null, indicating no record filtering.
     */
    String getFieldSelector()
        {
        return m_sFieldSelector;
        }

    /**
     * Sets field selector String.
     *
     * @param sFieldSelector  the field selector
     */
    void setFieldSelector(String sFieldSelector)
        {
        this.m_sFieldSelector = sFieldSelector;
        }

    /**
     * Returns a selector to limit results to those with matching labels.
     * @return the option, if specified. Defaults to null, indicating no record filtering.
     */
    String getLabelSelector()
        {
        return m_sLabelSelector;
        }

    /**
     * Sets label selector String.
     *
     * @param sLabelSelector  the label selector
     */
    void setLabelSelector(String sLabelSelector)
    {
        this.m_sLabelSelector = sLabelSelector;
    }

    /**
     * Returns the <code>pretty-print/code> option to be sent. If <code>true</code>, then the output is pretty printed.
     * @return the option, if specified. Defaults to null.
     */
    String getPretty()
        {
        return m_sPretty;
        }

    /**
     * Sets pretty String.
     *
     * @param sPretty  the pretty printed option
     */
    void setPretty(String sPretty)
        {
        this.m_sPretty = sPretty;
        }

    /**
     * On a watch call: when specified, shows changes that occur after that particular version of a resource.
     *                  Defaults to changes from the beginning of history.
     * On a list call: when specified, requests values at least as recent as the specified value.
     *                  Defaults to returning the result from remote storage based on quorum-read flag;
     *                  - if it&#39;s 0, then we simply return what we currently have in cache, no guarantee;
     *                  - if set to non zero, then the result is at least as fresh as given version.
     * @return the current setting. Defaults to null.
     */
    String getResourceVersion()
        {
        return m_sResourceVersion;
        }

    /**
     * Sets resource version String.
     *
     * @param sResourceVersion  the resource version
     */
    void setResourceVersion(String sResourceVersion)
        {
        this.m_sResourceVersion = sResourceVersion;
        }

    /**
     * Returns a listener for responses received, to specify on calls.
     * @return the set listener. Defaults to null.
     */
    ProgressResponseBody.ProgressListener getProgressListener()
        {
        return m_progressListener;
        }

    /**
     * Sets progress listener for process response body.
     *
     * @param progressListener  the response progress listener
     */
    void setProgressListener(ProgressResponseBody.ProgressListener progressListener)
        {
        this.m_progressListener = progressListener;
        }

    /**
     * Returns a listener for requests sent, to specify on calls.
     * @return the set listener. Defaults to null.
     */
    ProgressRequestBody.ProgressRequestListener getProgressRequestListener()
        {
        return m_progressRequestListener;
        }

    /**
     * Sets progress request listener of progress request body.
     *
     * @param progressRequestListener  the request progress listener
     */
    void setProgressRequestListener(ProgressRequestBody.ProgressRequestListener progressRequestListener)
        {
        this.m_progressRequestListener = progressRequestListener;
        }

    // ----- constants ------------------------------------------------------

    /**
     * The default maximium number of responses for a Kubernetes client call.
     */
    private static final int DEFAULT_LIMIT           = 500;

    /**
     * The default timeout in seconds for a Kubernetes client call.
     */
    private static final int DEFAULT_TIMEOUT_SECONDS = 30;


    // ----- data members ---------------------------------------------------

    /**
     * A boolean indicating whether partially initialized results should be included in the response.
     */
    private Boolean m_fIncludeUninitialized;

    /**
     * The limit on the number of updates to send in a single reply.
     */
    private Integer m_nLimit = DEFAULT_LIMIT;

    /**
     * The timeout for the call.
     */
    private Integer m_nTimeoutSeconds = DEFAULT_TIMEOUT_SECONDS;

    /**
     * A selector to limit results to those with matching fields.
     */
    private String m_sFieldSelector;

    /**
     * A selector to limit results to those with matching labels.
     */
    private String m_sLabelSelector;

    /**
     * The <code>pretty-print</code> option to be sent. If <code>true</code>, then the output is pretty printed.
     */
    private String m_sPretty;

    /**
     * The last resource version.
     */
    private String m_sResourceVersion;

    /**
     * A listener for responses received, to specify on calls.
     */
    private ProgressResponseBody.ProgressListener m_progressListener;

    /**
     * A listener for requests sent, to specify on calls.
     */
    private ProgressRequestBody.ProgressRequestListener m_progressRequestListener;
    }
