/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples.testing;

import com.oracle.bedrock.Option;
import com.oracle.bedrock.OptionsByType;

/**
 * An {@link Option} that signifies the number of times
 * to retry a failed operation.
 * <p>
 * The default value for this option is zero (no retry attempts).
 *
 * @author jk
 */
public class MaxRetries
        implements Option
    {
    // ----- constructors ---------------------------------------------------

    private MaxRetries(int nCount)
        {
        m_nCount = nCount;
        }

    // ----- MaxRetries methods ---------------------------------------------

    /**
     * Obtain the maximum number of retries to attempt.
     *
     * @return  the maximum number of retries to attempt
     */
    public int getMaxRetryCount()
        {
        return m_nCount;
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Obtain a {@link MaxRetries} option that will
     * disable retries.
     *
     * @return  a {@link MaxRetries} option that will
     *          disable retries
     */
    @OptionsByType.Default
    public static MaxRetries none()
        {
        return new MaxRetries(0);
        }

    /**
     * Obtain a {@link MaxRetries} option with the specified
     * retry count.
     *
     * @param nCount the number of retries to attempt
     *
     * @return  a {@link MaxRetries} option with the specified
     *          retry count
     */
    public static MaxRetries of(int nCount)
        {
        return new MaxRetries(nCount);
        }

    // ----- data members ---------------------------------------------------

    /**
     * The maximum number of retries.
     */
    private final int m_nCount;
    }
