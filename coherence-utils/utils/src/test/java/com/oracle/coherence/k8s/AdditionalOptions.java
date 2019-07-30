/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.bedrock.Option;
import com.oracle.bedrock.OptionsByType;
import com.oracle.bedrock.runtime.Application;
import com.oracle.bedrock.runtime.MetaClass;
import com.oracle.bedrock.runtime.Platform;
import com.oracle.bedrock.runtime.Profile;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Iterator;
import java.util.List;
import java.util.stream.Collectors;

/**
 * A Bedrock {@link Profile} and {@link Option} to make
 * it simple to add an additional array of {@link Option}s
 * to a method that has a var-args {@link Option}s parameter.
 * 
 * @author jk
 */
public class AdditionalOptions
        implements Profile, Option, Option.Collector<AdditionalOptions.OptionSet, AdditionalOptions>
    {
    // ----- constructors ---------------------------------------------------

    @OptionsByType.Default
    public AdditionalOptions()
        {
        }

    public AdditionalOptions(AdditionalOptions copy)
        {
        f_list.addAll(copy.f_list);
        }

    // ----- Option.Collector methods ---------------------------------------

    @Override
    public AdditionalOptions with(AdditionalOptions.OptionSet optionSet)
        {
        AdditionalOptions copy = new AdditionalOptions(this);

        copy.f_list.add(optionSet);

        return copy;
        }

    @Override
    public AdditionalOptions without(AdditionalOptions.OptionSet optionSet)
        {
        AdditionalOptions copy;

        if (f_list.contains(optionSet))
            {
            copy = new AdditionalOptions(this);
            copy.f_list.remove(optionSet);

            }
        else
            {
            copy = this;
            }

        return copy;
        }

    @Override
    @SuppressWarnings("unchecked")
    public <O> Iterable<O> getInstancesOf(Class<O> requiredClass)
        {
        return (List<O>) f_list.stream()
                               .flatMap(opts -> Arrays.stream(opts.m_aOption))
                               .filter(requiredClass::isInstance)
                               .collect(Collectors.toList());
        }

    @Override
    public Iterator<OptionSet> iterator()
        {
        return f_list.iterator();
        }

    // ----- Profile methods ------------------------------------------------

    @Override
    public void onLaunching(Platform platform, MetaClass metaClass, OptionsByType optionsByType)
        {
        f_list.forEach(optionSet -> optionSet.addTo(optionsByType));
        }

    @Override
    public void onLaunched(Platform platform, Application application, OptionsByType optionsByType)
        {
        }

    @Override
    public void onClosing(Platform platform, Application application, OptionsByType optionsByType)
        {
        }

    // ----- helper methods -------------------------------------------------

    public static Option of(Option... options)
        {
        return new OptionSet(options);
        }

    // ----- inner class: OptionSet -----------------------------------------

    public static class OptionSet
            implements Option, Option.Collectable
        {
        // ----- constructors -----------------------------------------------

        public OptionSet(Option[] aOption)
            {
            m_aOption = aOption;
            }

        // ----- OptionSet methods ------------------------------------------

        void addTo(OptionsByType optionsByType)
            {
            if (m_aOption != null && m_aOption.length > 0)
                {
                optionsByType.addAll(m_aOption);
                }
            }

        // ----- Option.Collectable methods ---------------------------------

        @Override
        public Class<? extends Collector> getCollectorClass()
            {
            return AdditionalOptions.class;
            }

        // ----- object methods -------------------------------------------------

        @Override
        public boolean equals(Object o)
            {
            if (this == o)
                {
                return true;
                }
            if (o == null || getClass() != o.getClass())
                {
                return false;
                }
            OptionSet optionSet = (OptionSet) o;
            return Arrays.equals(m_aOption, optionSet.m_aOption);
            }

        @Override
        public int hashCode()
            {
            return Arrays.hashCode(m_aOption);
            }

        // ----- data members ---------------------------------------------------

        private final Option[] m_aOption;
        }

    // ----- data members ---------------------------------------------------

    private final List<OptionSet> f_list = new ArrayList<>();
    }
