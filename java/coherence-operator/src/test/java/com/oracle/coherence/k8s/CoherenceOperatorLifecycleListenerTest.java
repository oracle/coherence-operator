/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.Map;

import org.junit.jupiter.api.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.collection.IsMapContaining.hasEntry;

public class CoherenceOperatorLifecycleListenerTest {

    @Test
    public void shouldHaveNullResumeMapForNullServicesString() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap(null);
        assertThat(resumeMap, is(nullValue()));
    }

    @Test
    public void shouldHaveNullResumeMapEmptyServicesString() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("");
        assertThat(resumeMap, is(nullValue()));
    }

    @Test
    public void shouldHaveNullResumeMapBlankServicesString() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("  ");
        assertThat(resumeMap, is(nullValue()));
    }

    @Test
    public void shouldHaveSingleServiceEnabled() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("\"foo\"=true");
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(1));
        assertThat(resumeMap, hasEntry("foo", true));
    }

    @Test
    public void shouldHaveSingleServiceDisabled() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("\"foo\"=false");
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(1));
        assertThat(resumeMap, hasEntry("foo", false));
    }

    @Test
    public void shouldHaveEscapedQuoteInServiceName() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("\"\\\"foo\"=true");
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(1));
        assertThat(resumeMap, hasEntry("\"foo", true));
    }

    @Test
    public void shouldHaveEqualsInServiceName() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("\"=foo\"=true");
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(1));
        assertThat(resumeMap, hasEntry("=foo", true));
    }

    @Test
    public void shouldHaveCommaInServiceName() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("\",foo\"=true");
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(1));
        assertThat(resumeMap, hasEntry(",foo", true));
    }

    @Test
    public void shouldHaveMultipleServices() {
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap("\"foo\"=true,\"bar\"=false");
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(2));
        assertThat(resumeMap, hasEntry("foo", true));
        assertThat(resumeMap, hasEntry("bar", false));
    }

    @Test
    public void shouldSkipEmptyNames() {
        String s = ",\"foo\"=true,,\"bar\"=false,";
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap(s);
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(2));
        assertThat(resumeMap, hasEntry("foo", true));
        assertThat(resumeMap, hasEntry("bar", false));
    }

    @Test
    public void shouldDecodeBase64() {
        String s = "\"foo\"=true,\"bar\"=false";
        String encoded = "base64:" + Base64.getEncoder().encodeToString(s.getBytes(StandardCharsets.UTF_8));
        Map<String, Boolean> resumeMap = CoherenceOperatorLifecycleListener.getResumeMap(encoded);
        assertThat(resumeMap, is(notNullValue()));
        assertThat(resumeMap.size(), is(2));
        assertThat(resumeMap, hasEntry("foo", true));
        assertThat(resumeMap, hasEntry("bar", false));
    }
}
