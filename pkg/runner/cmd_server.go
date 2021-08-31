/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	"os"
	"strconv"
	"strings"
)

const (
	// CommandServer is the argument to launch a server.
	CommandServer = "server"
)

// serverCommand creates the corba "server" sub-command
func serverCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandServer,
		Short: "Start a Coherence server",
		Long:  "Starts a Coherence server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, server)
		},
	}
}

// Configure the runner to run a Coherence Server
func server(details *RunDetails, _ *cobra.Command) {
	details.Command = CommandServer
	details.MainClass = ServerMain

	// If the main class environment variable is set then use that
	// otherwise run Coherence DCS.
	mc, found := details.lookupEnv(v1.EnvVarAppMainClass)
	switch {
	case found && details.AppType != AppTypeSpring:
		// we have a main class specified, and we're not a Spring Boot app
		details.MainArgs = []string{mc}
	case found && details.AppType == AppTypeSpring:
		// we have a main class and the app is Spring Boot
		// the main is PropertiesLauncher,
		details.MainClass = SpringBootMain
		// the specified main class is set as a Spring loader property
		details.addArg("-Dloader.main=" + mc)
	case !found && details.AppType == AppTypeSpring:
		// the app type is Spring Boot so main is PropertiesLauncher
		details.MainClass = SpringBootMain
	case !found && details.AppType == AppTypeCoherence:
		// the app type is Coherence so main is DCS
		details.MainArgs = []string{DCS}
	case !found && details.AppType == AppTypeHelidon:
		// the app type is Helidon so main is the Helidon CDI starter
		details.MainArgs = []string{HelidonMain}
	default:
		// no main or app type specified, use DCS
		details.MainArgs = []string{DCS}
	}

	// Check for any main class arguments
	ma, found := details.lookupEnv(v1.EnvVarAppMainArgs)
	if found {
		if ma != "" {
			for _, arg := range strings.Split(ma, " ") {
				details.MainArgs = append(details.MainArgs, details.ExpandEnv(arg))
			}
		}
	}

	// Configure the Coherence member's role
	details.setSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohRole, "-Dcoherence.role", "storage")
	// Configure whether this member is storage enabled
	details.addArgFromEnvVar(v1.EnvVarCohStorage, "-Dcoherence.distributed.localstorage")

	// Configure Coherence Tracing
	ratio := details.Getenv(v1.EnvVarCohTracingRatio)
	if ratio != "" {
		q, err := resource.ParseQuantity(ratio)
		if err == nil {
			d := q.AsDec()
			details.addArg("-Dcoherence.tracing.ratio=" + d.String())
		} else {
			fmt.Printf("ERROR: Coherence tracing ratio \"%s\" is invalid - %s\n", ratio, err.Error())
			os.Exit(1)
		}
	}

	// Configure whether Coherence management is enabled
	hasMgmt := details.isEnvTrue(v1.EnvVarCohMgmtPrefix + v1.EnvVarCohEnabledSuffix)
	log.Info("Coherence Management over REST", "enabled", strconv.FormatBool(hasMgmt), "envVar", v1.EnvVarCohMgmtPrefix+v1.EnvVarCohEnabledSuffix)
	if hasMgmt {
		fmt.Println("INFO: Configuring Coherence Management over REST")
		details.addArg("-Dcoherence.management.http=all")
		if details.CoherenceHome != "" {
			// If management is enabled and the COHERENCE_HOME environment variable is set
			// then $COHERENCE_HOME/lib/coherence-management.jar will be added to the classpath
			details.addClasspath(details.CoherenceHome + "/lib/coherence-management.jar")
		}
	}

	// Configure whether Coherence metrics is enabled
	hasMetrics := details.isEnvTrue(v1.EnvVarCohMetricsPrefix + v1.EnvVarCohEnabledSuffix)
	log.Info("Coherence Metrics", "enabled", strconv.FormatBool(hasMetrics), "envVar", v1.EnvVarCohMetricsPrefix+v1.EnvVarCohEnabledSuffix)
	if hasMetrics {
		details.addArg("-Dcoherence.metrics.http.enabled=true")
		fmt.Println("INFO: Configuring Coherence Metrics")
		if details.CoherenceHome != "" {
			// If metrics is enabled and the COHERENCE_HOME environment variable is set
			// then $COHERENCE_HOME/lib/coherence-metrics.jar will be added to the classpath
			details.addClasspath(details.CoherenceHome + "/lib/coherence-metrics.jar")
		}
	}

	// Configure whether to add third-party modules to the classpath if management over rest
	// or metrics are enabled and the directory pointed to by the DEPENDENCY_MODULES environment
	// variable exists.
	if hasMgmt || hasMetrics {
		dm := details.Getenv(v1.EnvVarCohDependencyModules)
		if dm != "" {
			stat, err := os.Stat(dm)
			if err == nil && stat.IsDir() {
				// dependency modules directory exists
				details.addClasspath(dm + "/*")
			}
		}
	}

	if details.isEnvTrueOrBlank(v1.EnvVarJvmShowSettings) {
		details.addArg("-XshowSettings:all")
		details.addArg("-XX:+PrintCommandLineFlags")
		details.addArg("-XX:+PrintFlagsFinal")
	}

	// Add GC logging parameters if required
	if details.isEnvTrue(v1.EnvVarJvmGcLogging) {
		details.addArg("-verbose:gc")
		details.addArg("-XX:+PrintGCDetails")
		details.addArg("-XX:+PrintGCTimeStamps")
		details.addArg("-XX:+PrintHeapAtGC")
		details.addArg("-XX:+PrintTenuringDistribution")
		details.addArg("-XX:+PrintGCApplicationStoppedTime")
		details.addArg("-XX:+PrintGCApplicationConcurrentTime")
	}
}
