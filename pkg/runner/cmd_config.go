/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"bytes"
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
)

var configLog = ctrl.Log.WithName("config")

// configCommand creates the corba "config" sub-command
func configCommand(env map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   v1.RunnerConfig,
		Short: "Create the Operator JVM args files a Coherence server",
		Long:  "Create the Operator JVM args files a Coherence server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return maybeRun(cmd, createsFiles)
		},
	}

	utilDir, found := env[v1.EnvVarCohUtilDir]
	if !found || utilDir == "" {
		utilDir = v1.VolumeMountPathUtils
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgUtilsDir, utilDir, "The utils files root directory")

	return cmd
}

// createsFiles will create the various config files
func createsFiles(details *RunDetails, _ *cobra.Command) (bool, error) {
	populateMainClass(details)
	populateServerDetails(details)
	err := configureCommand(details)
	if err != nil {
		return false, errors.Wrap(err, "failed to configure server command")
	}
	if err := createClassPathFile(details); err != nil {
		return false, err
	}
	if err := createArgsFile(details); err != nil {
		return false, err
	}
	if err := createMainClassFile(details); err != nil {
		return false, err
	}
	if err := createSpringBootFile(details); err != nil {
		return false, err
	}
	if err := createCliConfig(details); err != nil {
		return false, err
	}
	return false, nil
}

// createClassPathFile will create the class path files for a Coherence Pod - typically this is run from an init-container
func createClassPathFile(details *RunDetails) error {
	cpFile := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorClasspathFile)
	classpath := details.getClasspath()
	err := os.WriteFile(cpFile, []byte(classpath), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to write coherence classpath file")
	}
	configLog.Info("Created class path file", "FileName", cpFile, "ClassPath", classpath)
	return nil
}

// createArgsFile will create the JVM args files for a Coherence Pod - typically this is run from an init-container
func createArgsFile(details *RunDetails) error {
	args := details.GetArgs()
	argFileName := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorJvmArgsFile)

	var buffer bytes.Buffer
	for _, arg := range args {
		buffer.WriteString(arg + "\n")
	}
	if err := os.WriteFile(argFileName, buffer.Bytes(), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write JVM args file "+argFileName)
	}

	configLog.Info("Created JVM args file", "FileName", argFileName, "Args", buffer.String())
	return nil
}

// createSpringBootFile will create the SpringBoot JVM args files for a Coherence Pod - typically this is run from an init-container
func createSpringBootFile(details *RunDetails) error {
	argsFile := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorSpringBootArgsFile)
	cp := strings.ReplaceAll(details.getClasspath(), ":", ",")

	var args string
	if details.InnerMainClass == "" || details.InnerMainClass == v1.DefaultMain {
		args = fmt.Sprintf("-Dloader.path=%s", cp)
	} else {
		args = fmt.Sprintf("-Dloader.path=%s\n-Dloader.main=%s", cp, details.InnerMainClass)
	}

	err := os.WriteFile(argsFile, []byte(args), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to write coherence classpath file")
	}
	configLog.Info("Created SpringBoot args file", "FileName", argsFile, "Args", args)
	return nil
}

// createMainClassFile will create the file containing the main class name for a Coherence Pod - typically this is run from an init-container
func createMainClassFile(details *RunDetails) error {
	fileName := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorMainClassFile)

	var s string
	if details.InnerMainClass == "" || details.IsSpringBoot() {
		s = details.MainClass
	} else {
		s = fmt.Sprintf("%s\n%s", details.MainClass, details.InnerMainClass)
	}

	if err := os.WriteFile(fileName, []byte(s), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write coherence classpath file")
	}
	configLog.Info("Created main class file", "FileName", fileName, "MainClass", details.InnerMainClass)
	return nil
}

func createCliConfig(details *RunDetails) error {
	home := details.getenvOrDefault(v1.EnvVarCohCtlHome, details.UtilsDir)
	fileName := fmt.Sprintf("%s%c%s", home, os.PathSeparator, "cohctl.yaml")

	cluster := details.Getenv(v1.EnvVarCohClusterName)
	port := details.Getenv(v1.EnvVarCohMgmtPrefix + v1.EnvVarCohPortSuffix)
	if port == "" {
		port = fmt.Sprintf("%d", v1.DefaultManagementPort)
	}
	protocol := details.Getenv(v1.EnvVarCohCliProtocol)
	if protocol == "" {
		protocol = "http"
	}

	var buffer bytes.Buffer
	buffer.WriteString("clusters:\n")
	buffer.WriteString("    - name: default\n")
	buffer.WriteString("      discoverytype: manual\n")
	buffer.WriteString("      connectiontype: " + protocol + "\n")
	buffer.WriteString("      connectionurl: " + protocol + "://127.0.0.1:" + port + "/management/coherence/cluster\n")
	buffer.WriteString("      nameservicediscovery: \"\"\n")
	buffer.WriteString("      clusterversion: \"\"\n")
	buffer.WriteString("      clustername: \"" + cluster + "\"\n")
	buffer.WriteString("      clustertype: Standalone\n")
	buffer.WriteString("      manuallycreated: true\n")
	buffer.WriteString("      baseclasspath: \"\"\n")
	buffer.WriteString("      additionalclasspath: \"\"\n")
	buffer.WriteString("      arguments: \"\"\n")
	buffer.WriteString("      managementport: " + port + "\n")
	buffer.WriteString("      persistencemode: \"\"\n")
	buffer.WriteString("      loggingdestination: \"\"\n")
	buffer.WriteString("      managementavailable: false\n")
	buffer.WriteString("color: \"on\"\n")
	buffer.WriteString("currentcontext: default\n")
	buffer.WriteString("debug: false\n")
	buffer.WriteString("defaultbytesformat: m\n")
	buffer.WriteString("ignoreinvalidcerts: false\n")
	buffer.WriteString("requesttimeout: 30\n")
	if err := os.WriteFile(fileName, buffer.Bytes(), os.ModePerm); err != nil {
		configLog.Error(err, "Failed to write coherence CLI config file", "FileName", fileName)
		return nil
	}

	configLog.Info("Created CLI config file", "FileName", fileName, "Config", buffer.String())
	return nil
}

// Configure the main class
func populateMainClass(details *RunDetails) {
	details.MainClass = ServerMain

	// If the main class environment variable is set then use that
	// otherwise run Coherence DefaultMain.
	mc, found := details.lookupEnv(v1.EnvVarAppMainClass)

	if !found || mc == "" {
		// no custom mani set so check for a JIB main class file
		appDir := details.getenvOrDefault(v1.EnvVarCohAppDir, "/app")
		jibMainClassFileName := filepath.Join(appDir, "jib-main-class-file")
		fi, err := os.Stat(jibMainClassFileName)
		if err != nil && appDir != "/app" {
			// try /app dir
			jibMainClassFileName = "/app/jib-main-class-file"
			fi, err = os.Stat(jibMainClassFileName)
		}

		if err == nil && (fi.Size() != 0) {
			mainCls, _ := readFirstLineFromFile(jibMainClassFileName)
			if len(mainCls) != 0 {
				mc = mainCls
				found = true
			}
		}
	}

	isSpring := details.IsSpringBoot()
	switch {
	case found && !isSpring:
		// we have a main class specified, and we're not a Spring Boot app
		details.InnerMainClass = mc
	case found && details.AppType == AppTypeSpring2:
		// we have a main class and the app is Spring Boot 2.x
		// the main is PropertiesLauncher,
		details.MainClass = SpringBootMain2
		// the specified main class is set as a Spring loader property
		details.InnerMainClass = mc
	case found && details.AppType == AppTypeSpring3:
		// we have a main class and the app is Spring Boot 3.x
		// the main is PropertiesLauncher,
		details.MainClass = SpringBootMain3
		// the specified main class is set as a Spring loader property
		details.InnerMainClass = mc
	case !found && details.AppType == AppTypeSpring2:
		// the app type is Spring Boot 2.x so main is PropertiesLauncher
		details.MainClass = SpringBootMain2
	case !found && details.AppType == AppTypeSpring3:
		// the app type is Spring Boot 3.x so main is PropertiesLauncher
		details.MainClass = SpringBootMain3
	case !found && details.AppType == AppTypeCoherence:
		// the app type is Coherence so main is DefaultMain
		details.InnerMainClass = DefaultMain
	case !found && details.AppType == AppTypeHelidon:
		// the app type is Helidon so main is the Helidon CDI starter
		details.InnerMainClass = HelidonMain
	default:
		// no main or app type specified, use DefaultMain
		details.InnerMainClass = DefaultMain
	}
}

// Configure the runner to run a Coherence Server
func populateServerDetails(details *RunDetails) {
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
			// This is for legacy 14.1.1.0 and 12.2.1.4 images
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
			// This is for legacy 14.1.1.0 and 12.2.1.4 images
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
