/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
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
 * A {@link TestRule} specifying what Coherence versions a test can be run against.
 * The current Coherence version being tested must be equal or greater than the
 * specified minimal Coherence version.
 */
public class AssumingCoherenceVersion
        implements TestRule {

    // ----- Constants ------------------------------------------------------

    /**
     * Current Coherence version being tested by a test class with this {@link TestRule}.
     */
    private final String coherenceVersion;

    /**
     * Minimal Coherence version required for a test class with this {@link TestRule}.
     */
    private final String minimalCoherenceVersion;

    // ----- constructors ---------------------------------------------------

    /**
     * Specify minimal Coherence Version for a test class by this {@link TestRule}.
     *
     * @param coherenceVersion        current Coherence version being tested
     * @param minimalCoherenceVersion minimal Coherence version
     */
    public AssumingCoherenceVersion(String coherenceVersion, String minimalCoherenceVersion) {
        if (coherenceVersion.contains(":")) {
            this.coherenceVersion = coherenceVersion.substring(coherenceVersion.indexOf(":") + 1);
        } else {
            this.coherenceVersion = coherenceVersion;
        }
        this.minimalCoherenceVersion = minimalCoherenceVersion;
    }

    // ----- TestRule interface ---------------------------------------------

    @Override
    public Statement apply(Statement base, Description description) {
        return new Statement() {
            @Override
            public void evaluate() throws Throwable {
                if (!CoherenceVersion.versionCheck(coherenceVersion, minimalCoherenceVersion)) {

                    throw new AssumptionViolatedException("Specified Coherence Version " + coherenceVersion + " does not "
                                                                  + "meet required Coherence version " +
                                                                  minimalCoherenceVersion + " or greater.");
                } else {
                    base.evaluate();
                }
            }
        };
    }
}
