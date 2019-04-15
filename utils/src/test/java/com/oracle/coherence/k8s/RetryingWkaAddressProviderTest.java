/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.common.util.Duration;

import com.tangosol.net.ConfigurableAddressProvider.AddressHolder;
import com.tangosol.util.Base;
import com.tangosol.util.WrapperException;

import org.junit.Test;

import java.net.UnknownHostException;
import java.util.Arrays;
import java.util.Iterator;

import static com.oracle.coherence.k8s.RetryingWkaAddressProvider.*;
import static com.oracle.common.util.Duration.Magnitude.MILLI;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.Matchers.is;
import static org.hamcrest.Matchers.lessThanOrEqualTo;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

public class RetryingWkaAddressProviderTest
    {
    @Test(expected = UnknownHostException.class)
    public void testShouldTimeoutOnNonExistentDnsReference()
        throws UnknownHostException
        {
        final long FREQUENCY_MS = 500;
        final long TIMEOUT_MS   = 5000;

        long ldtStart = Base.getLastSafeTimeMillis();
        Iterable<AddressHolder> holders = new Iterable<AddressHolder>()
            {
            @Override
            public Iterator<AddressHolder> iterator()
                {
                return Arrays.asList(new AddressHolder("NonExiStentHoStName12345678", 0)).iterator();
                }

            };
        RetryingWkaAddressProvider provider = new RetryingWkaAddressProvider(holders, true, FREQUENCY_MS, TIMEOUT_MS);
        try
            {
            provider.eventuallyResolve();
            fail("should throw exception");
            }
        finally
            {
            long ldtTimeoutDuration = Base.getLastSafeTimeMillis() - ldtStart;

            assertThat("validate reresolve happened over specified minimum timeout",
                ldtTimeoutDuration, greaterThanOrEqualTo(TIMEOUT_MS));
            assertThat("validate lower bound of reresolve count" ,
                provider.m_nLastReresolveCount, greaterThanOrEqualTo(3));
            assertThat("validate upper bound of reresolve count" ,
                provider.m_nLastReresolveCount, lessThanOrEqualTo((int)(TIMEOUT_MS / FREQUENCY_MS)));
            }
        }

    @Test
    public void testShouldRevolveImmediately()
        throws UnknownHostException
        {
        final long FREQUENCY_MS = 3022;
        final long TIMEOUT_MS   = 13456;

        RetryingWkaAddressProvider provider =
            (RetryingWkaAddressProvider) create("127.0.0.1", FREQUENCY_MS, TIMEOUT_MS);

        assertNotNull("confirm wka resolved", provider.getNextAddress());
        assertThat("validate configured frequency of reresolve", provider.f_WkaDNSReresolveFrequency_ms, is(FREQUENCY_MS));
        assertThat("validate configured frequency of reresolve", provider.f_WkaDNSResolutionTimeout_ms, is(TIMEOUT_MS));
        }

    @Test
    public void shouldConfigureBySystemProperties()
        throws UnknownHostException
        {
        String sTimeout   = "4m";
        String sFrequency = "22s";

        System.setProperty(PROP_WKA_OVERRIDE, "127.0.0.1");
        System.setProperty(PROP_WKA_TIMEOUT, sTimeout);
        System.setProperty(PROP_WKA_RERESOLVE_FREQUENCY, sFrequency);

        try
            {
            RetryingWkaAddressProvider provider = (RetryingWkaAddressProvider) create();
            assertNotNull("confirm wka resolved", provider.getNextAddress());
            assertThat("validate configured frequency of dns resolve",
                provider.f_WkaDNSReresolveFrequency_ms, is(new Duration(sFrequency).as(MILLI)));
            assertThat("validate configured max time to attempt to resolve wka dns addresses",
                provider.f_WkaDNSResolutionTimeout_ms, is(new Duration(sTimeout).as(MILLI)));
            }
        finally
            {
            System.clearProperty(PROP_WKA_OVERRIDE);
            System.clearProperty(PROP_WKA_TIMEOUT);
            System.clearProperty(PROP_WKA_RERESOLVE_FREQUENCY);
            }
        }

    @Test
    public void testCreateWithDurationParameters()
        throws UnknownHostException
        {
        final String FREQUENCY_DURATION = "2s";
        final String TIMEOUT_DURATION = "6s";

        System.setProperty(PROP_WKA_OVERRIDE, "127.0.0.1");

        try
            {
            RetryingWkaAddressProvider provider =
                (RetryingWkaAddressProvider) create(FREQUENCY_DURATION, TIMEOUT_DURATION);

            assertNotNull("confirm wka resolved", provider.getNextAddress());
            assertThat("validate configured frequency of reresolve",
                provider.f_WkaDNSReresolveFrequency_ms, is(new Duration(FREQUENCY_DURATION).as(MILLI)));
            assertThat("validate configured frequency of reresolve",
                provider.f_WkaDNSResolutionTimeout_ms, is(new Duration(TIMEOUT_DURATION).as(MILLI)));
            }
        finally
            {
            System.clearProperty(PROP_WKA_OVERRIDE);
            }
        }
    }
