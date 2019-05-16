/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package util;

import com.oracle.bedrock.OptionsByType;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;


public class MaxRetriesTest
    {
    @Test
    public void shouldHaveZeroRetries()
        {
        MaxRetries maxRetries = MaxRetries.none();
        assertThat(maxRetries.getMaxRetryCount(), is(0));
        }

    @Test
    public void shouldHaveSpecifiedRetries()
        {
        MaxRetries maxRetries = MaxRetries.of(19);
        assertThat(maxRetries.getMaxRetryCount(), is(19));
        }

    @Test
    public void shouldHaveDefaultValue()
        {
        OptionsByType options    = OptionsByType.empty();
        MaxRetries    maxRetries = options.get(MaxRetries.class);
        assertThat(maxRetries, is(notNullValue()));
        assertThat(maxRetries.getMaxRetryCount(), is(0));
        }

    @Test
    public void shouldHaveSpecifiedOption()
        {
        OptionsByType options    = OptionsByType.of(MaxRetries.of(10));
        MaxRetries    maxRetries = options.get(MaxRetries.class);
        assertThat(maxRetries, is(notNullValue()));
        assertThat(maxRetries.getMaxRetryCount(), is(10));
        }
    }
