/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package util;


import com.oracle.bedrock.Option;
import com.oracle.bedrock.runtime.Application;
import org.junit.Rule;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;


public class KubernetesTest
    {
    @Test
    public void shouldNotRetrySuccess()
        {
        Kubernetes  kubernetes  = spy(new Kubernetes());
        Application application = mock(Application.class);
        Option[]    options     = new Option[0];

        doReturn(application).when(kubernetes).kubectl(any());
        when(application.waitFor(any())).thenReturn(0);

        int nExitCode = kubernetes.kubectlAndWait(options);

        assertThat(nExitCode, is(0));
        verify(kubernetes, times(1)).kubectl(any());
        }

    @Test
    public void shouldRetryFiveTimesByDefault()
        {
        Kubernetes  kubernetes  = spy(new Kubernetes());
        Application application = mock(Application.class);
        Option[]    options     = new Option[0];

        doReturn(application).when(kubernetes).kubectl(any());
        when(application.waitFor(any())).thenReturn(1);

        int nExitCode = kubernetes.kubectlAndWait(options);

        assertThat(nExitCode, is(1));
        verify(kubernetes, times(6)).kubectl(any());
        }

    @Test
    public void shouldRetryUntilSuccess()
        {
        Kubernetes  kubernetes  = spy(new Kubernetes());
        Application application = mock(Application.class);
        Option[]    options     = new Option[]{MaxRetries.of(5)};

        doReturn(application).when(kubernetes).kubectl(any());
        when(application.waitFor(any())).thenReturn(1, 1, 0);

        int nExitCode = kubernetes.kubectlAndWait(options);

        assertThat(nExitCode, is(0));
        verify(kubernetes, times(3)).kubectl(any());
        }

    @Test
    public void shouldRetryMaxRetriesTimes()
        {
        Kubernetes  kubernetes  = spy(new Kubernetes());
        Application application = mock(Application.class);
        Option[]    options     = new Option[]{MaxRetries.of(3)};

        doReturn(application).when(kubernetes).kubectl(any());
        when(application.waitFor(any())).thenReturn(1);

        int nExitCode = kubernetes.kubectlAndWait(options);

        assertThat(nExitCode, is(1));
        verify(kubernetes, times(4)).kubectl(any());
        }

    @Test
    public void shouldNotRetryIfMaxRetriesIsNone()
        {
        Kubernetes  kubernetes  = spy(new Kubernetes());
        Application application = mock(Application.class);
        Option[]    options     = new Option[]{MaxRetries.none()};

        doReturn(application).when(kubernetes).kubectl(any());
        when(application.waitFor(any())).thenReturn(1);

        int nExitCode = kubernetes.kubectlAndWait(options);

        assertThat(nExitCode, is(1));
        verify(kubernetes, times(1)).kubectl(any());
        }
    }
