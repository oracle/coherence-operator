/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package util;

import com.oracle.coherence.k8s.CoherenceVersion;
import org.junit.AssumptionViolatedException;
import org.junit.rules.TestRule;
import org.junit.runner.Description;
import org.junit.runners.model.Statement;

/**
 * A {@TestRule} specifying what Coherence versions a test can be run against.
 * The current Coherence version being tested must be equal or greater than the
 * specified minimal Coherence version.
 */
public class AssumingCoherenceVersion
    implements TestRule
    {
    // ----- constructors ---------------------------------------------------

    /**
     * Specify minimal Coherence Version for a test class by this {@link TestRule}.
     *
     * @param sCoherenceVersion         current Coherence version being tested
     * @param sMinimalCoherenceVersion  minimal Coherence version
     */
    public AssumingCoherenceVersion(String sCoherenceVersion, String sMinimalCoherenceVersion)
        {
        f_sCoherenceVersion        = sCoherenceVersion;
        f_sMinimalCoherenceVersion = sMinimalCoherenceVersion;
        }

    // ----- TestRule interface ---------------------------------------------

    @Override
    public Statement apply(Statement base, Description description)
        {
        return new Statement()
            {
            @Override
            public void evaluate() throws Throwable
                {
                if (!CoherenceVersion.versionCheck(f_sCoherenceVersion, f_sMinimalCoherenceVersion))
                    {

                    throw new AssumptionViolatedException("Specified Coherence Version " + f_sCoherenceVersion + " does not meet required Coherence version " +
                        f_sMinimalCoherenceVersion + " or greater.");
                    }
                else
                    {
                    base.evaluate();
                    }
                }
            };
        }

    // ----- Constants ------------------------------------------------------

    /**
     * Current Coherence version being tested by a test class with this {@link TestRule}.
     */
    private final String f_sCoherenceVersion;

    /**
     * Minimal Coherence version required for a test class with this {@link TestRule}.
     */
    private final String f_sMinimalCoherenceVersion;
    }
