/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.bedrock.Option;
import com.oracle.bedrock.runtime.LocalPlatform;
import com.oracle.bedrock.runtime.coherence.CoherenceCluster;
import com.oracle.bedrock.runtime.coherence.CoherenceClusterBuilder;
import com.oracle.bedrock.runtime.coherence.CoherenceClusterMember;
import com.oracle.bedrock.runtime.coherence.options.ClusterName;
import com.oracle.bedrock.runtime.coherence.options.LocalHost;
import com.oracle.bedrock.runtime.coherence.options.WellKnownAddress;
import com.oracle.bedrock.runtime.java.JavaApplication;
import com.oracle.bedrock.runtime.java.options.ClassName;
import com.oracle.bedrock.runtime.java.options.IPv4Preferred;
import com.oracle.bedrock.runtime.java.options.SystemProperty;
import com.oracle.bedrock.runtime.options.Argument;
import com.oracle.bedrock.runtime.options.DisplayName;
import com.oracle.bedrock.testsupport.junit.TestLogs;
import com.oracle.bedrock.util.Capture;
import org.junit.Rule;
import org.junit.Test;

import java.util.UUID;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * @author jk
 */
public class ClusterMemberProbeIT
    {
    @Test
    public void shouldBeLiveIfCoherenceStartedFirst()
        {
        try (CoherenceCluster cluster = startCluster(1))
            {
            try (JavaApplication appPodCheck = runPodChecker(PodChecker.Type.liveness))
                {
                int nExitCode = appPodCheck.waitFor();
                assertThat(nExitCode, is(0));
                }
            }
        }

    @Test
    public void shouldNotBeLiveIfCoherenceStartedSecond()
        {
        try (JavaApplication appPodCheck = runBlockingApp())
            {
            try (CoherenceCluster cluster = startCluster(1))
                {
                appPodCheck.submit(() -> PodChecker.main(new String[]{PodChecker.Type.liveness.name()}));

                int nExitCode = appPodCheck.waitFor();
                assertThat(nExitCode, is(1));
                }
            }
        }

    @Test
    public void shouldBeReadyWhenSingleMember()
        {
        try (CoherenceCluster cluster = startCluster(1))
            {
            try (JavaApplication appPodCheck = runPodChecker(PodChecker.Type.readiness))
                {
                int nExitCode = appPodCheck.waitFor();
                assertThat(nExitCode, is(0));
                }
            }
        }

    @Test
    public void shouldBeReadyWhenMultipleMembers()
        {
        try (CoherenceCluster cluster = startCluster(2))
            {
            try (JavaApplication appPodCheck = runPodChecker(PodChecker.Type.readiness))
                {
                int nExitCode = appPodCheck.waitFor();
                assertThat(nExitCode, is(0));
                }
            }
        }

    @Test
    public void shouldBeReadyWhenMultipleMembersWithZeroBackups()
        {
        try (CoherenceCluster cluster = startCluster(2))
            {
            try (JavaApplication appPodCheck = runPodChecker(PodChecker.Type.readiness))
                {
                int nExitCode = appPodCheck.waitFor();
                assertThat(nExitCode, is(0));
                }
            }
        }

    // ----- helper methods -------------------------------------------------

    private JavaApplication runPodChecker(PodChecker.Type type)
        {
        return runJavaApp(ClassName.of(PodChecker.class), Argument.of(type));
        }

    private JavaApplication runBlockingApp()
        {
        return runJavaApp(ClassName.of(BlockingApp.class));
        }

    private JavaApplication runJavaApp(Option... aOpt)
        {
        return LocalPlatform.get().launch(JavaApplication.class,
                                          getCommonOptions(),
                                          AdditionalOptions.of(aOpt));

        }

    private CoherenceCluster startCluster(int cMembers, Option... options) 
        {
        CoherenceClusterBuilder builder = new CoherenceClusterBuilder();

        builder.include(cMembers,
                        CoherenceClusterMember.class,
                        DisplayName.of("DCS"),
                        getCommonOptions(),
                        AdditionalOptions.of(options));

        return builder.build();
        }

    private Option getCommonOptions()
        {
        return AdditionalOptions.of(m_testLogs.builder(),
                                    ClusterName.of(m_sCluster),
                                    LocalHost.only(),
                                    SystemProperty.of(ProbeHttpClient.PROP_HTTP_PORT, m_httpPort),
                                    WellKnownAddress.of("127.0.0.1", 0),
                                    IPv4Preferred.yes());
        }

    // ----- data members ---------------------------------------------------

    private String m_sCluster = UUID.randomUUID().toString();

    private final Capture<Integer> m_httpPort = new Capture<>(LocalPlatform.get().getAvailablePorts());

    @Rule
    public TestLogs m_testLogs = new TestLogs();
    }
