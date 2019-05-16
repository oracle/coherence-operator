/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package util;

import com.oracle.bedrock.Option;
import com.oracle.bedrock.OptionsByType;
import com.oracle.bedrock.deferred.options.RetryFrequency;
import com.oracle.bedrock.options.Timeout;
import com.oracle.bedrock.runtime.k8s.K8sCluster;
import com.oracle.bedrock.util.Duration;

import java.util.Iterator;
import java.util.concurrent.TimeUnit;


/**
 * An extension of a Bedrock {@link K8sCluster} that allows retrying
 * of failed kubectl commands when calling the {@link #kubectlAndWait(Option...)}
 * or {@link #kubectlAndWait(Timeout, Option...)} methods.
 *
 * @author jk  2019.05.16
 */
public class Kubernetes
        extends K8sCluster<Kubernetes>
    {
    /**
     * Execute a kubectl command against the k8s cluster
     * and wait for the command to complete.
     * <p>
     * This implementation allows retrying of failed commands (kubectl
     * returns a non-zero exit code) by specifying a {@link MaxRetries}
     * option.
     *
     * @param timeout  the time to wait for the command to complete
     * @param options  the options to use to run the kubectl command
     *
     * @return  the exit code from the kubectl command
     */
    @Override
    public int kubectlAndWait(Timeout timeout, Option... options)
        {
        OptionsByType      opts           = OptionsByType.of(options);
        Iterator<Duration> retryFrequency = opts.get(RetryFrequency.class).get().iterator();
        int                nMaxRetries    = opts.get(MaxRetries.class).getMaxRetryCount();

        int nExitCode = super.kubectlAndWait(timeout, options);

        while (nExitCode != 0 && nMaxRetries > 0)
            {
            nMaxRetries--;
            long nMillis = retryFrequency.hasNext() ? retryFrequency.next().to(TimeUnit.MILLISECONDS) : 250L;
            System.err.println("Kubectl returned a non-zero exit code (" + nExitCode + ") "
                               + "- retrying command in " + nMillis + " millis. "
                               + "Remaining retry attempts " + nMaxRetries);

            try
                {
                Thread.sleep(nMillis);
                }
            catch (InterruptedException e)
                {
                System.err.println("Kubectl retry thread interrupted, stopping retries");
                break;
                }

            nExitCode = super.kubectlAndWait(timeout, options);
            }

        return nExitCode;
        }
    }
