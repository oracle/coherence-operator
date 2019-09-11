/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.tangosol.net.CacheFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;


/**
 * A class that test the state of a Coherence Pod running in Kubernetes.
 *
 * @author jk
 */
public class PodChecker
    {
    // ----- constructors ---------------------------------------------------

    /**
     * Create a {@link PodChecker}.
     */
    public PodChecker()
        {
        f_listProbe = Arrays.asList(new OperatorMBeanProbe(),
                                    new RestProbe(),
                                    new ClusterMemberProbe());
        }

    // ----- public methods -------------------------------------------------

    /**
     * Run the Pod test.
     *
     * @param args  the arguments controlling the test
     *
     * @return  zero if the test passed or non-zero if the test failed
     */
    public int run(String... args)
        {
        String  sType    = args == null || args.length == 0 ? null : args[0];
        boolean fResult = false;

        try
            {
            Type type = (sType == null) ? Type.readiness : Type.valueOf(sType);

            switch (type)
                {
                case readiness:
                    fResult = readiness();
                    break;
                case statusha:
                    fResult = statusHA();
                    break;
                case liveness:
                    fResult = liveness();
                    break;
                default:
                    CacheFactory.log(getClass() + " unrecognised probe type '" + type + "'", CacheFactory.LOG_ERR);
                    return RETURN_CODE_NOT_READY;
                }
            }
        catch (IllegalArgumentException e)
            {
            CacheFactory.log(getClass() + " probe type argument is invalid '" + sType + "'", CacheFactory.LOG_ERR);
            }

        return fResult ? RETURN_CODE_READY : RETURN_CODE_NOT_READY;
        }

    /**
     * Perform the readiness test using the first available {@link Probe}.
     *
     * @return  zero if the Pod is ready
     */
    boolean readiness()
        {
        try (Probe probe = findActiveProbe())
            {
            return probe.isReady();
            }
        }

    /**
     * Perform the StatusHA test using the first available {@link Probe}.
     *
     * @return  zero if the Pod is ready
     */
    boolean statusHA()
        {
        try (Probe probe = findActiveProbe())
            {
            return probe.isStatusHA();
            }
        }

    /**
     * Perform the liveness test using the first available {@link Probe}.
     *
     * @return  zero if the Pod is alive
     */
    boolean liveness()
        {
        try (Probe probe = findActiveProbe())
            {
            return probe.isLive();
            }
        }

    /**
     * Find the first {@link Probe} that is active.
     *
     * @return  the first {@link Probe} that is active
     */
    Probe findActiveProbe()
        {
        return f_listProbe.stream()
                          .filter(Probe::isAvailable)
                          .findFirst()
                          .orElseThrow(() -> new IllegalStateException("No active Probe class available"));
        }

    /**
     * Main Method that takes probe type i.e either readiness or liveness in the
     * list of arguments to the process.
     *
     * @param args  the program arguments
     */
    public static void main(String[] args)
        {
        try
            {
            PodChecker probe     = new PodChecker();
            int        nExitCode = probe.run(args);

            System.exit(nExitCode);
            }
        catch (Throwable t)
            {
            System.err.println("Error executing Pod tests.");
            t.printStackTrace();

            System.exit(RETURN_CODE_NOT_READY);
            }
        }

    // ----- inner enum: Type -----------------------------------------------

    /**
     * The supported types of probe test.
     */
    public enum Type
        {
        /**
         * Perform a readiness test.
         */
        readiness,

        /**
         * Perform a StausHA test.
         */
        statusha,

        /**
         * Perform a liveness test.
         */
        liveness
        }

    // ----- constants ---------------------------------------------------

    /**
     * The result returned if the Pod test passes.
     */
    public static final int RETURN_CODE_READY = 0;

    /**
     * The result returned if the Pod test fails.
     */
    public static final int RETURN_CODE_NOT_READY = 1;

    /**
     * The {@link List} of {@link Probe}s to use to execute the test.
     * <p>
     * When executing a test the first available {@link Probe} in this
     * list will be used.
     */
    private final List<Probe> f_listProbe;
    }
