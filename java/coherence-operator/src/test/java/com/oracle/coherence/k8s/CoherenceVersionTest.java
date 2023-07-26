/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import org.junit.jupiter.api.Test;

import static org.hamcrest.CoreMatchers.is;

import static org.hamcrest.MatcherAssert.assertThat;

/**
 * CoherenceVersion tests.
 */
public class CoherenceVersionTest {
    @Test
    public void shouldBeGreater() {
        assertThat(CoherenceVersion.versionCheck("1", "0"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1", "0"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1", "1.0"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1", "1.1.0"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1", "1.1.1.0"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1.1", "1.1.1.1.0"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1.1.1", "1.1.1.1.1.0"), is(true));
        assertThat(CoherenceVersion.versionCheck("2.1", "1.2"), is(true));
        assertThat(CoherenceVersion.versionCheck("2.1-some-text", "1.2"), is(true));
    }

    @Test
    public void shouldBeEqual() {
        assertThat(CoherenceVersion.versionCheck("1", "1"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1", "1"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1", "1.1"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1", "1.1.1"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1", "1.1.1.1"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1.1", "1.1.1.1.1"), is(true));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1.1.1", "1.1.1.1.1.1"), is(true));
    }

    @Test
    public void shouldBeLess() {
        assertThat(CoherenceVersion.versionCheck("1", "2"), is(false));
        assertThat(CoherenceVersion.versionCheck("1.1", "2"), is(false));
        assertThat(CoherenceVersion.versionCheck("1.1", "1.2"), is(false));
        assertThat(CoherenceVersion.versionCheck("1.1.1", "1.1.2"), is(false));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1", "1.1.1.2"), is(false));
        assertThat(CoherenceVersion.versionCheck("1.1.1.1.1", "1.1.1.1.2"), is(false));
        assertThat(CoherenceVersion.versionCheck("1.2", "2.1"), is(false));
    }

    @Test
    public void shouldWorkWithInterimBuild() throws Exception {
        assertThat(CoherenceVersion.versionCheck("14.1.1.0.15 (101966-Int)", "14.1.1.0.0"), is(true));
        assertThat(CoherenceVersion.versionCheck("14.1.1.0.15 (101966-Int)", "22.06.0"), is(false));
    }

}
