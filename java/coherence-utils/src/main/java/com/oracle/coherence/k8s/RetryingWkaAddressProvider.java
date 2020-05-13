/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.io.IOException;
import java.net.UnknownHostException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import com.tangosol.net.AddressProvider;
import com.tangosol.net.ConfigurableAddressProvider;
import com.tangosol.util.Base;

/**
 * An AddressProvider that eventually resolves at least one dns host name in provided WKA list, {@link #PROP_WKA_OVERRIDE}.
 * Throws an exception if unable to resolve at least one host name within {@link #wkaDNSResolutionTimeout}.
 */
public class RetryingWkaAddressProvider
        extends ConfigurableAddressProvider {

    /**
     * The name of the System property to use to return a fixed WKA list.
     */
    public static final String PROP_WKA_OVERRIDE = "coherence.wka";

    /**
     * The name of the System property to configure maximum time to attempt to resolve wka addresses.
     * Provides default value for {link #f_WkaDNSResolutionTimeout_ms}.
     * Set this system property to the number of milli-seconds.
     */
    public static final String PROP_WKA_TIMEOUT = RetryingWkaAddressProvider.class.getName() + ".dnsResolutionTimeout";

    /**
     * System property for configuring the frequency to attempt dns resolution of wka addresses.
     * Provides default value for {link #f_WkaDNSReresolveFrequency_ms}.
     * Set this system property to the number of milli-seconds.
     */
    public static final String PROP_WKA_RERESOLVE_FREQUENCY = RetryingWkaAddressProvider.class
            .getName() + ".dnsResolutionFrequency";

    /**
     * WKA DNS Resolution frequency.
     */
    private final long wkaDNSReresolveFrequency;

    /**
     * The timeout value for wks address resolution.
     */
    private final long wkaDNSResolutionTimeout;

    /**
     * Added for testing verification.
     * Count of how many times a DNS resolve of entire WKA address list has been performed.
     */
    private int lastReresolveCount;

    // ----- Constructors ---------------------------------------------------

    /**
     * Construct a {@link RetryingWkaAddressProvider}.
     *
     * @param addressHolders the {@link AddressHolder}s
     * @param fSafe          true if the provider skips unresolved addresses
     * @param timeout        maximum time in milliseconds to attempt to resolve {@link AddressHolder}s
     * @param frequency      frequency in milliseconds to attempt to retry {@link AddressHolder} dns resolution
     */
    public RetryingWkaAddressProvider(Iterable<AddressHolder> addressHolders, boolean fSafe, long frequency, long timeout) {
        super(addressHolders, fSafe);

        wkaDNSReresolveFrequency = frequency;
        wkaDNSResolutionTimeout = timeout;
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
     * @return an {@link AddressProvider}.
     * @throws UnknownHostException if the WKA address is invalid
     */
    public static AddressProvider create() throws UnknownHostException {
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
     * @param sWkaOverride the comma delimited WKA address list
     * @return an {@link AddressProvider}.
     * @throws UnknownHostException if unable to dns resolve at least one address in wka list
     */
    public static AddressProvider create(String sWkaOverride) throws UnknownHostException {
        return create(sWkaOverride,
                      Long.getLong(PROP_WKA_RERESOLVE_FREQUENCY, 2L),
                      Long.getLong(PROP_WKA_TIMEOUT, 6L * 60L * 1000L));
    }

    /**
     * Create an {@link AddressProvider} configured by provided parameters.
     *
     * @param wkaOverride the WKA address to use
     * @param frequency   frequency in milliseconds to retry dns resolution of wka address list
     * @param timeout     timeout in milliseconds to abort retry of dns resolution of wka address list and throw an exception
     * @return {@link AddressProvider}
     * @throws UnknownHostException if unable to dns resolve at least one address in wka list
     */
    static AddressProvider create(String wkaOverride, long frequency, long timeout) throws UnknownHostException {
        if (wkaOverride == null) {
            return new ConfigurableAddressProvider(Collections.emptyList(), true);
        }
        else {
            String[] asAddresses = wkaOverride.split(",");

            List<AddressHolder> list = Arrays.stream(asAddresses)
                    .map(sAddr -> new ConfigurableAddressProvider.AddressHolder(sAddr, 0))
                    .collect(Collectors.toList());

            RetryingWkaAddressProvider provider = new RetryingWkaAddressProvider(list, true, frequency, timeout);
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
     * @param durationFrequency retry wka resolve frequency as duration string
     * @param durationTimeout   timeout wka resolve as duration string
     *
     * @return {@link AddressProvider}
     * @throws UnknownHostException if unable to dns resolve at least one address in wka list
     */
    public static AddressProvider create(String durationFrequency, String durationTimeout) throws UnknownHostException {
        return create(System.getProperty(PROP_WKA_OVERRIDE),
                      Long.parseLong(durationFrequency),
                      Long.parseLong(durationTimeout));
    }

    long getWkaDNSReresolveFrequency() {
        return wkaDNSReresolveFrequency;
    }

    long getWkaDNSResolutionTimeout() {
        return wkaDNSResolutionTimeout;
    }

    int getLastReresolveCount() {
        return lastReresolveCount;
    }

    // ----- helpers --------------------------------------------------------

    /**
     * Attempt resolution of each dns reference in wka until at least on dns reference resolves
     * or throw an {@link IOException} after {@link #wkaDNSResolutionTimeout}.
     *
     * @return this {@link RetryingWkaAddressProvider} instance
     * @throws UnknownHostException if no dns references resolve within {@link #wkaDNSResolutionTimeout}.
     */
    public AddressProvider eventuallyResolve()
            throws UnknownHostException {
        long start = Base.getLastSafeTimeMillis();

        lastReresolveCount = 0;
        while ((Base.getLastSafeTimeMillis() - start) < wkaDNSResolutionTimeout) {
            lastReresolveCount++;
            if (getNextAddress() == null) {
                Base.sleep(wkaDNSReresolveFrequency);
                reset();
            }
            else {
                reset();
                return this;
            }
        }

        throw new UnknownHostException(RetryingWkaAddressProvider.class.getName()
                                               + " failed to resolve configured WKA address(es) "
                                               + System.getProperty(PROP_WKA_OVERRIDE)
                                               + " within " + wkaDNSResolutionTimeout + " milliseconds.");
    }
}
