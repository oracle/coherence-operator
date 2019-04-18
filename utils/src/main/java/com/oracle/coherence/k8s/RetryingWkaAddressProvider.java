/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.common.util.Duration;
import com.oracle.common.util.Duration.Magnitude;
import com.tangosol.net.AddressProvider;
import com.tangosol.net.CacheFactory;
import com.tangosol.net.ConfigurableAddressProvider;
import com.tangosol.util.Base;
import com.tangosol.util.WrapperException;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.UnknownHostException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

/**
 * An AddressProvider that eventually resolves at least one dns host name in provided WKA list, {@link #PROP_WKA_OVERRIDE}.
 * Throws an exception if unable to resolve at least one host name within {@link #f_WkaDNSResolutionTimeout_ms}.
 */
public class RetryingWkaAddressProvider
    extends ConfigurableAddressProvider
    {
    // ----- Constructors ---------------------------------------------------

    /**
     * Construct a {@link RetryingWkaAddressProvider}
     *
     * @param addressHolders  the {@link AddressHolder}s
     * @param fSafe           true if the provider skips unresolved addresses
     * @param timeout_ms      maximum time in milliseconds to attempt to resolve {@link AddressHolder}s
     * @param frequency_ms    frequency in milliseconds to attempt to retry {@link AddressHolder} dns resolution
     */
    public RetryingWkaAddressProvider(Iterable<AddressHolder> addressHolders, boolean fSafe, long frequency_ms, long timeout_ms)
        {
        super(addressHolders, fSafe);

        f_WkaDNSReresolveFrequency_ms = frequency_ms;
        f_WkaDNSResolutionTimeout_ms  = timeout_ms;
        }

    /**
     * Create an {@link AddressProvider} using System properties for the WKA override.
     * <p>
     * The {@link #PROP_WKA_OVERRIDE} property is used to obtain a WKA override value, if
     * this value is set a {@link ConfigurableAddressProvider} will be returned using the
     * addresses from the {@link #PROP_WKA_OVERRIDE} property.
     * <p>
     * If {@link #PROP_WKA_OVERRIDE} is not set an empty
     * {@link ConfigurableAddressProvider} will be returned.
     *
     *
     * @return an {@link AddressProvider}.
     */
    public static AddressProvider create()
        throws UnknownHostException
        {
        return create(System.getProperty(PROP_WKA_OVERRIDE));
        }

    /**
     * Create an {@link AddressProvider} using a comma delimited list of addresses and configured with
     * system properties {@link #PROP_WKA_RERESOLVE_FREQUENCY} and {@link #PROP_WKA_TIMEOUT}.
     * <p>
     * The sWkaOverride parameter is used to obtain a WKA override value, if
     * this parameter is not null a {@link ConfigurableAddressProvider} will be
     * returned using the addresses from the sWkaOverride parameter.
     * <p>
     * If the sWkaOverride parameter is null then an empty
     * {@link ConfigurableAddressProvider} will be returned.
     *
     * @param sWkaOverride  the comma delimited WKA address list
     *
     * @return an {@link AddressProvider}.
     *
     * @throws UnknownHostException if unable to dns resolve at least one address in wka list
     */
    public static AddressProvider create(String sWkaOverride)
        throws UnknownHostException
        {
        return create(sWkaOverride, new Duration(System.getProperty(PROP_WKA_RERESOLVE_FREQUENCY, "2s")).as(Duration.Magnitude.MILLI),
            new Duration(System.getProperty(PROP_WKA_TIMEOUT, "6m")).as(Duration.Magnitude.MILLI));
        }

    /**
     * Create an {@link AddressProvider} configured by provided parameters.
     *
     * @param frequency_ms frequency in milliseconds to retry dns resolution of wka address list
     * @param timeout_ms   timeout in milliseconds to abort retry of dns resolution of wka address list and throw an exception
     *
     * @return {@link AddressProvider}
     *
     * @throws UnknownHostException if unable to dns resolve at least one address in wka list
     */
    public static AddressProvider create(String sWkaOverride, long frequency_ms, long timeout_ms)
        throws UnknownHostException
        {
        if (sWkaOverride == null)
            {
            return new ConfigurableAddressProvider(Collections.emptyList(), true);
            }
        else
            {
            String[] asAddresses = sWkaOverride.split(",");

            List<AddressHolder> list = Arrays.stream(asAddresses)
                .map(sAddr -> new ConfigurableAddressProvider.AddressHolder(sAddr, 0))
                .collect(Collectors.toList());

            RetryingWkaAddressProvider provider = new RetryingWkaAddressProvider(list, true, frequency_ms, timeout_ms);
            return provider.eventuallyResolve();
            }
        }

    /**
     * Create an {@link AddressProvider} using System property for the WKA override
     * <p>
     * The {@link #PROP_WKA_OVERRIDE} property is used to obtain a WKA override value, if
     * this value is set a {@link ConfigurableAddressProvider} will be returned using the
     * addresses from the {@link #PROP_WKA_OVERRIDE} property.
     * <p>
     * If {@link #PROP_WKA_OVERRIDE} is not set an empty
     * {@link ConfigurableAddressProvider} will be returned.
     *
     * @param sDurationFrequency retry wka resolve frequency as duration string, see format documented in {@see Duration(String)}
     * @param sDurationTimeout   timeout wka resolve as duration string, examples are "120s", "2m" and "2000ms", all durations of 2 minutes.
     *
     * @return {@link AddressProvider}
     *
     * @throws UnknownHostException if unable to dns resolve at least one address in wka list
     */
    public static AddressProvider create(String sDurationFrequency, String sDurationTimeout)
        throws UnknownHostException
        {
        return create(System.getProperty(PROP_WKA_OVERRIDE),
            new Duration(sDurationFrequency).as(Magnitude.MILLI), new Duration(sDurationTimeout).as(Magnitude.MILLI));
        }

    // ----- helpers --------------------------------------------------------

    /**
     * Attempt resolution of each dns reference in wka until at least on dns reference resolves
     * or throw an {@link IOException} after {@link #f_WkaDNSResolutionTimeout_ms}.
     *
     * @return this {@link RetryingWkaAddressProvider} instance
     *
     * @throws WrapperException if no dns references resolve within {@link #f_WkaDNSResolutionTimeout_ms}.
     */
    public AddressProvider eventuallyResolve()
        throws UnknownHostException
        {
        long ldtStart = Base.getLastSafeTimeMillis();

        m_nLastReresolveCount = 0;
        while ((Base.getLastSafeTimeMillis() - ldtStart) < f_WkaDNSResolutionTimeout_ms)
            {
            m_nLastReresolveCount++;

            InetSocketAddress addr = getNextAddress();

            if (addr == null)
                {
                Base.sleep(f_WkaDNSReresolveFrequency_ms);
                reset();
                }
            else
                {
                reset();
                CacheFactory.log("RetryingWkaAddressProvider: resolved " + addr.getHostName() + " in " + (Base.getLastSafeTimeMillis() - ldtStart) + " ms",
                    Base.LOG_INFO);
                return this;
                }
            }

        throw new UnknownHostException(RetryingWkaAddressProvider.class.getName() +
            " failed to resolve configured WKA address(es) " + System.getProperty(PROP_WKA_OVERRIDE) +
            " within " + f_WkaDNSResolutionTimeout_ms + " milliseconds.");
        }

    /**
     * The name of the System property to use to return a fixed WKA list.
     */
    public static final String PROP_WKA_OVERRIDE = "coherence.wka";

    /**
     * The name of the System property to configure maximum time to attempt to resolve wka addresses.
     * Provides default value for {link #f_WkaDNSResolutionTimeout_ms}. Set this system property
     * to the format of the string parameter described in {@link Duration(String)}.  Example settings
     * are "10s", "20m", "60000ms" for 10 seconds, 20 minutes and 60,000 milliseconds respectively.
     */
    public static final String PROP_WKA_TIMEOUT = RetryingWkaAddressProvider.class.getName() + ".dnsResolutionTimeout";

    /**
     * System property for configuring the frequency to attempt dns resolution of wka addresses.
     * Provides default value for {link #f_WkaDNSReresolveFrequency_ms}. Set this system property
     * to the format of the string parameter described in {@link Duration(String)}.  Example settings
     * are "10s", "20m", "60000ms" for 10 seconds, 20 minutes and 60,000 milliseconds respectively.
     */
    public static final String PROP_WKA_RERESOLVE_FREQUENCY = RetryingWkaAddressProvider.class.getName() + ".dnsResolutionFrequency";

    /**
     * WKA DNS Resolution frequency.
     */
    public final long f_WkaDNSReresolveFrequency_ms;

    /**
     * The timeout value for wks address resolution.
     */
    public final long f_WkaDNSResolutionTimeout_ms;

    /**
     * Added for testing verification.
     * Count of how many times a DNS resolve of entire WKA address list has been performed.
     */
    int m_nLastReresolveCount;
    }
