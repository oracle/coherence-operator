/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"k8s.io/apimachinery/pkg/api/resource"
	"net/http"
	"os"
	"os/exec"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
	"time"
)

// The code that actually starts the process in the Coherence container.

const (
	// DefaultMain is an indicator to run the default main class.
	DefaultMain = "$DEFAULT$"
	// HelidonMain is the default Helidon main class name.
	HelidonMain = "io.helidon.microprofile.server.Main"
	// ServerMain is the default server main class name.
	ServerMain = "com.oracle.coherence.k8s.Main"
	// SpringBootMain is the default Spring Boot main class name.
	SpringBootMain = "org.springframework.boot.loader.PropertiesLauncher"

	// AppTypeNone is the argument to specify no application type.
	AppTypeNone = ""
	// AppTypeJava is the argument to specify a Java application.
	AppTypeJava = "java"
	// AppTypeCoherence is the argument to specify a Coherence application.
	AppTypeCoherence = "coherence"
	// AppTypeHelidon is the argument to specify a Helidon application.
	AppTypeHelidon = "helidon"
	// AppTypeSpring is the argument to specify a Spring application.
	AppTypeSpring = "spring"
	// AppTypeOperator is the argument to specify running an Operator command.
	AppTypeOperator = "operator"

	// defaultConfig is the root name of the default configuration file
	defaultConfig = ".coherence-runner"
)

var (
	// An alternative configuration file to use instead of program arguments
	cfgFile string

	// backoffSchedule is a sequence of back-off times for re-trying http requests.
	backoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}

	// log is the logger used by the runner
	log = ctrl.Log.WithName("runner")

	// dryRun is true if the execution is a dry-run
	dryRun = false
)

// contextKey allows type safe Context Values.
type contextKey int

// The key to obtain an execution from a Context.
var executionKey contextKey

// Execution is a holder of details of a command execution
type Execution struct {
	Cmd   *cobra.Command
	App   string
	OsCmd *exec.Cmd
	V     *viper.Viper
}

// NewRootCommand builds the root cobra command that handles our command line tool.
func NewRootCommand(env map[string]string) (*cobra.Command, *viper.Viper) {
	v := viper.New()

	// rootCommand is the Cobra root Command to execute
	rootCmd := &cobra.Command{
		Use:   "runner",
		Short: "Start the Coherence operator runner",
		Long:  "runner starts the Coherence Operator runner",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd, v, env)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Just print information about the commands that would execute")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s.yaml)", defaultConfig))
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	rootCmd.AddCommand(initCommand(env))
	rootCmd.AddCommand(serverCommand())
	rootCmd.AddCommand(consoleCommand())
	rootCmd.AddCommand(queryPlusCommand())
	rootCmd.AddCommand(statusCommand())
	rootCmd.AddCommand(readyCommand())
	rootCmd.AddCommand(nodeCommand())
	rootCmd.AddCommand(operatorCommand())
	rootCmd.AddCommand(networkTestCommand())
	rootCmd.AddCommand(jShellCommand())
	rootCmd.AddCommand(sleepCommand())

	return rootCmd, v
}

func initializeConfig(cmd *cobra.Command, v *viper.Viper, env map[string]string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".coherence" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfig)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	// v.SetEnvPrefix(EnvPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind any environment overrides
	for key, value := range env {
		v.Set(key, value)
	}

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// bindFlags binds each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = v.BindEnv(f.Name, envVarSuffix)
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// Execute runs the runner with a given environment.
func Execute() (Execution, error) {
	return ExecuteWithArgs(nil, nil)
}

// ExecuteWithArgs runs the runner with a given environment and argument overrides.
func ExecuteWithArgs(env map[string]string, args []string) (Execution, error) {
	cmd, v := NewRootCommand(env)
	if len(args) > 0 {
		cmd.SetArgs(args)
	}

	e := Execution{
		Cmd: cmd,
		V:   v,
	}

	ctx := context.WithValue(context.Background(), executionKey, &e)
	err := cmd.ExecuteContext(ctx)
	return e, err
}

// RunFunction is a function to run a command
type RunFunction func(*RunDetails, *cobra.Command)

// MaybeRunFunction is a function to maybe run a command depending on the return bool
type MaybeRunFunction func(*RunDetails, *cobra.Command) (bool, error)

// always is a wrapper around a RunFunction to turn it into a MaybeFunction that always runs
type always struct {
	Fn RunFunction
}

// run will wrap a RunFunction and always return true
func (in always) run(details *RunDetails, cmd *cobra.Command) (bool, error) {
	in.Fn(details, cmd)
	return true, nil
}

// run executes the required command.
func run(cmd *cobra.Command, fn RunFunction) error {
	a := always{Fn: fn}
	return maybeRun(cmd, a.run)
}

// maybeRun executes the required command.
func maybeRun(cmd *cobra.Command, fn MaybeRunFunction) error {
	var err error
	e := fromContext(cmd.Context())

	details := NewRunDetails(e.V)
	runCommand, err := fn(details, cmd)
	if err != nil {
		return err
	}

	if runCommand {
		e.App, e.OsCmd, err = createCommand(details)

		if err != nil {
			return err
		}

		if e.OsCmd != nil {
			b := new(bytes.Buffer)
			sep := ""
			for _, value := range e.OsCmd.Env {
				_, _ = fmt.Fprintf(b, "%s%s", sep, value)
				sep = ", "
			}

			log.Info("Executing command", "dryRun", dryRun, "application", e.App,
				"path", e.OsCmd.Path, "args", strings.Join(e.OsCmd.Args, " "), "env", b.String())

			if !dryRun {
				return e.OsCmd.Run()
			}
		}
	}
	return nil
}

// fromContext obtains the current execution from the specified context.
func fromContext(ctx context.Context) *Execution {
	e, ok := ctx.Value(executionKey).(*Execution)
	if ok {
		return e
	}
	return &Execution{}
}

// create the process to execute.
func createCommand(details *RunDetails) (string, *exec.Cmd, error) {
	var err error

	//ns := os.Getenv(v1.EnvVarCoherenceNamespace)
	hostname := os.Getenv(v1.EnvVarCohMemberName)
	//subdomain := os.Getenv(v1.EnvVarCoherenceSubDomain)
	//localhost := hostname + "." + subdomain + "." + ns
	details.addArg("-Dcoherence.localhost=" + hostname)

	// Set standard system properties
	details.addArgFromEnvVar(v1.EnvVarCohWka, "-Dcoherence.wka")
	details.addArgFromEnvVar(v1.EnvVarCohMachineName, "-Dcoherence.machine")
	details.addArgFromEnvVar(v1.EnvVarCohMemberName, "-Dcoherence.member")
	details.addArgFromEnvVar(v1.EnvVarCohClusterName, "-Dcoherence.cluster")
	details.addArgFromEnvVar(v1.EnvVarCohCacheConfig, "-Dcoherence.cacheconfig")
	details.addArgFromEnvVar(v1.EnvVarCohIdentity, "-Dcoherence.k8s.operator.identity")
	details.addArgFromEnvVar(v1.EnvVarCohForceExit, "-Dcoherence.k8s.operator.force.exit")
	details.setSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohHealthPort, "-Dcoherence.k8s.operator.health.port", fmt.Sprintf("%d", v1.DefaultHealthPort))
	details.setSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohMgmtPrefix+v1.EnvVarCohPortSuffix, "-Dcoherence.management.http.port", fmt.Sprintf("%d", v1.DefaultManagementPort))
	details.setSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohMetricsPrefix+v1.EnvVarCohPortSuffix, "-Dcoherence.metrics.http.port", fmt.Sprintf("%d", v1.DefaultMetricsPort))

	details.addArg("-XX:+UnlockDiagnosticVMOptions")

	// Configure the classpath to support images created with the JIB Maven plugin
	// This is enabled by default unless the image is a buildpacks image, or we
	// are running a Spring Boot application.
	if !details.isBuildPacks() && details.AppType != AppTypeSpring && details.isEnvTrueOrBlank(v1.EnvVarJvmClasspathJib) {
		appDir := details.getenvOrDefault(v1.EnvVarCohAppDir, "/app")
		fi, e := os.Stat(appDir + "/jib-classpath-file")
		if e == nil && (fi.Size() != 0) {
			clsPath := readFirstLineFromFile(appDir+"/jib-classpath-file", fi)
			if len(clsPath) != 0 {
				details.addClasspath(clsPath)
			}
		} else {
			details.addClasspathIfExists(appDir + "/resources")
			details.addClasspathIfExists(appDir + "/classes")
			details.addJarsToClasspath(appDir + "/classpath")
			details.addJarsToClasspath(appDir + "/libs")
		}
	}

	// Add the Operator Utils jar to the classpath
	details.addClasspath(details.UtilsDir + "/lib/coherence-operator.jar")
	details.addClasspath(details.UtilsDir + "/config")

	// Configure Coherence persistence
	mode := details.getenvOrDefault(v1.EnvVarCohPersistenceMode, "on-demand")
	details.addArg("-Dcoherence.distributed.persistence-mode=" + mode)

	persistence := details.Getenv(v1.EnvVarCohPersistenceDir)
	if persistence != "" {
		details.addArg("-Dcoherence.distributed.persistence.base.dir=" + persistence)
	}

	snapshots := details.Getenv(v1.EnvVarCohSnapshotDir)
	if snapshots != "" {
		details.addArg("-Dcoherence.distributed.persistence.snapshot.dir=" + snapshots)
	}

	// Set the Coherence site and rack values
	configureSiteAndRack(details)

	// Set the Coherence log level
	details.addArgFromEnvVar(v1.EnvVarCohLogLevel, "-Dcoherence.log.level")

	// Disable IPMonitor
	ipMon := details.Getenv(v1.EnvVarEnableIPMonitor)
	if ipMon != "TRUE" {
		details.addArg("-Dcoherence.ipmonitor.pingtimeout=0")
	}

	// Do the Coherence version specific configuration
	if ok := checkCoherenceVersion("12.2.1.4.0", details); ok {
		// is at least 12.2.1.4
		cohPost12214(details)
	} else {
		// is at pre-12.2.1.4
		cohPre12214(details)
	}

	post2206 := checkCoherenceVersion("14.1.1.2206.0", details)
	if post2206 {
		// at least CE 22.06
		cohPost2206(details)
	} else {
		post2006 := checkCoherenceVersion("14.1.1.2006.0", details)
		if !post2006 {
			// pre CE 20.06 - could be 14.1.1.2206
			if post14112206 := checkCoherenceVersion("14.1.1.2206.0", details); post14112206 {
				// at least 14.1.1.2206
				cohPost2206(details)
			}
		}
	}

	addManagementSSL(details)
	addMetricsSSL(details)

	// Get the Coherence member name
	member := details.Getenv(v1.EnvVarCohMemberName)
	if member == "" {
		member = "unknown"
	}

	allowEndangered := details.Getenv(v1.EnvVarCohAllowEndangered)
	if allowEndangered != "" {
		details.addArg("-Dcoherence.k8s.operator.statusha.allowendangered=" + allowEndangered)
	}

	// Get the K8s Pod UID
	podUID := details.Getenv(v1.EnvVarCohPodUID)
	if podUID == "" {
		podUID = "unknown"
	}

	// Configure the /jvm directory to hold heap dumps, jfr files etc. if the jvm root dir exists.
	jvmDir := v1.VolumeMountPathJVM + "/" + member + "/" + podUID
	if _, err = os.Stat(v1.VolumeMountPathJVM); err == nil {
		if err = os.MkdirAll(jvmDir, os.ModePerm); err != nil {
			return "", nil, err
		}
		if err = os.MkdirAll(jvmDir+"/jfr", os.ModePerm); err != nil {
			return "", nil, err
		}
		if err = os.MkdirAll(jvmDir+"/heap-dumps", os.ModePerm); err != nil {
			return "", nil, err
		}
	}

	details.addArg(fmt.Sprintf("-Dcoherence.k8s.operator.diagnostics.dir=%s", jvmDir))
	details.addArg(fmt.Sprintf("-XX:HeapDumpPath=%s/heap-dumps/%s-%s.hprof", jvmDir, member, podUID))

	// set the flag that allows the operator to resume suspended services on start-up
	if !details.isEnvTrueOrBlank(v1.EnvVarOperatorAllowResume) {
		details.addArg("-Dcoherence.k8s.operator.can.resume.services=false")
	} else {
		details.addArg("-Dcoherence.k8s.operator.can.resume.services=true")
	}

	if svc := details.Getenv(v1.EnvVarOperatorResumeServices); svc != "" {
		details.addArg("-Dcoherence.k8s.operator.resume.services=base64:" + svc)
	}

	gc := strings.ToLower(details.Getenv(v1.EnvVarJvmGcCollector))
	switch {
	case gc == "" || gc == "g1":
		details.addArg("-XX:+UseG1GC")
	case gc == "cms":
		details.addArg("-XX:+UseConcMarkSweepGC")
	case gc == "parallel":
		details.addArg("-XX:+UseParallelGC")
	}

	maxRAM := details.Getenv(v1.EnvVarJvmMaxRAM)
	if maxRAM != "" {
		details.addArg("-XX:MaxRAM=" + maxRAM)
	}

	heap := details.Getenv(v1.EnvVarJvmMemoryHeap)
	if heap != "" {
		// if heap is set use it
		details.addArg("-XX:InitialHeapSize=" + heap)
		details.addArg("-XX:MaxHeapSize=" + heap)
	} else {
		// if heap is not set check whether the individual heap values are set
		initialHeap := details.Getenv(v1.EnvVarJvmMemoryInitialHeap)
		if initialHeap != "" {
			details.addArg("-XX:InitialHeapSize=" + initialHeap)
		}
		maxHeap := details.Getenv(v1.EnvVarJvmMemoryMaxHeap)
		if maxHeap != "" {
			details.addArg("-XX:MaxHeapSize=" + maxHeap)
		}
	}

	percentageHeap := details.Getenv(v1.EnvVarJvmRAMPercentage)
	if percentageHeap != "" {
		// the heap percentage is set so use it
		q, err := resource.ParseQuantity(percentageHeap)
		if err == nil {
			d := q.AsDec()
			details.addArg("-XX:InitialRAMPercentage=" + d.String())
			details.addArg("-XX:MinRAMPercentage=" + d.String())
			details.addArg("-XX:MaxRAMPercentage=" + d.String())
		} else {
			log.Info("ERROR: Heap Percentage is not a valid resource.Quantity", "Value", percentageHeap, "Error", err.Error())
			os.Exit(1)
		}
	} else {
		// if heap is not set check whether the individual heap percentage values are set
		initial := details.Getenv(v1.EnvVarJvmInitialRAMPercentage)
		if initial != "" {
			q, err := resource.ParseQuantity(initial)
			if err == nil {
				d := q.AsDec()
				details.addArg("-XX:InitialRAMPercentage=" + d.String())
			} else {
				log.Info("ERROR: InitialRAMPercentage is not a valid resource.Quantity", "Value", initial, "Error", err.Error())
				os.Exit(1)
			}
		}

		max := details.Getenv(v1.EnvVarJvmMaxRAMPercentage)
		if max != "" {
			q, err := resource.ParseQuantity(max)
			if err == nil {
				d := q.AsDec()
				details.addArg("-XX:MaxRAMPercentage=" + d.String())
			} else {
				log.Info("ERROR: MaxRAMPercentage is not a valid resource.Quantity", "Value", max, "Error", err.Error())
				os.Exit(1)
			}
		}

		min := details.Getenv(v1.EnvVarJvmMinRAMPercentage)
		if min != "" {
			q, err := resource.ParseQuantity(min)
			if err == nil {
				d := q.AsDec()
				details.addArg("-XX:MinRAMPercentage=" + d.String())
			} else {
				log.Info("ERROR: MinRAMPercentage is not a valid resource.Quantity", "Value", min, "Error", err.Error())
				os.Exit(1)
			}
		}
	}

	direct := details.Getenv(v1.EnvVarJvmMemoryDirect)
	if direct != "" {
		details.addArg("-XX:MaxDirectMemorySize=" + direct)
	}
	stack := details.Getenv(v1.EnvVarJvmMemoryStack)
	if stack != "" {
		details.addArg("-Xss" + stack)
	}
	meta := details.Getenv(v1.EnvVarJvmMemoryMeta)
	if meta != "" {
		details.addArg("-XX:MetaspaceSize=" + meta)
		details.addArg("-XX:MaxMetaspaceSize=" + meta)
	}
	track := details.getenvOrDefault(v1.EnvVarJvmMemoryNativeTracking, "summary")
	if track != "" {
		details.addArg("-XX:NativeMemoryTracking=" + track)
		details.addArg("-XX:+PrintNMTStatistics")
	}

	// Configure debugging
	debugArgs := ""
	if details.isEnvTrue(v1.EnvVarJvmDebugEnabled) {
		var suspend string
		if details.isEnvTrue(v1.EnvVarJvmDebugSuspended) {
			suspend = "y"
		} else {
			suspend = "n"
		}

		port := details.Getenv(v1.EnvVarJvmDebugPort)
		if port == "" {
			port = fmt.Sprintf("%d", v1.DefaultDebugPort)
		}

		attach := details.Getenv(v1.EnvVarJvmDebugAttach)
		if attach == "" {
			debugArgs = fmt.Sprintf("-agentlib:jdwp=transport=dt_socket,server=y,suspend=%s,address=*:%s", suspend, port)
		} else {
			debugArgs = fmt.Sprintf("-agentlib:jdwp=transport=dt_socket,server=n,address=%s,suspend=%s,timeout=10000", attach, suspend)
		}
	}

	details.addArg("-Dcoherence.ttl=0")

	details.addArg(fmt.Sprintf("-XX:ErrorFile=%s/hs-err-%s-%s.log", jvmDir, member, podUID))

	if details.isEnvTrueOrBlank(v1.EnvVarJvmOomHeapDump) {
		details.addArg("-XX:+HeapDumpOnOutOfMemoryError")
	}

	if details.isEnvTrueOrBlank(v1.EnvVarJvmOomExit) {
		details.addArg("-XX:+ExitOnOutOfMemoryError")
	}

	// Use JVM container support
	if details.isEnvTrueOrBlank(v1.EnvVarJvmUseContainerLimits) {
		details.addArg("-XX:+UseContainerSupport")
	}

	details.addArgs(debugArgs)

	gcArgs := details.Getenv(v1.EnvVarJvmGcArgs)
	if gcArgs != "" {
		details.addArgs(strings.Split(gcArgs, " ")...)
	}

	jvmArgs := details.Getenv(v1.EnvVarJvmArgs)
	if jvmArgs != "" {
		details.addArgs(strings.Split(jvmArgs, " ")...)
	}

	var cmd *exec.Cmd
	var app string
	switch {
	case details.AppType == AppTypeNone || details.AppType == AppTypeJava:
		app = "Java"
		cmd, err = createJavaCommand(details.getJavaExecutable(), details)
	case details.AppType == AppTypeSpring:
		app = "SpringBoot"
		cmd, err = createSpringBootCommand(details.getJavaExecutable(), details)
	case details.AppType == AppTypeHelidon:
		app = "Java"
		cmd, err = createJavaCommand(details.getJavaExecutable(), details)
	case details.AppType == AppTypeCoherence:
		app = "Java"
		cmd, err = createJavaCommand(details.getJavaExecutable(), details)
	case details.AppType == AppTypeOperator:
		app = "Operator"
		cmd, err = createOperatorCommand(details)
	default:
		app = "Graal (" + details.AppType + ")"
		cmd, err = createGraalCommand(details)
	}

	return app, cmd, err
}

func createJavaCommand(javaCmd string, details *RunDetails) (*exec.Cmd, error) {
	args := details.getCommand()
	args = append(args, details.MainClass)
	return _createJavaCommand(javaCmd, details, args)
}

func readFirstLineFromFile(fqfn string, fi os.FileInfo) string {
	log.Info(fmt.Sprintf("%s size=%d", fi.Name(), fi.Size()))
	file, _ := os.Open(fqfn)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	closeFile(file, log)

	if len(text) == 0 {
		return ""
	}
	log.Info(fmt.Sprintf("First Line of the %s:\n%s\n", fi.Name(), text[0]))
	return text[0]
}

func createSpringBootCommand(javaCmd string, details *RunDetails) (*exec.Cmd, error) {
	if details.isBuildPacks() {
		return _createBuildPackCommand(details, SpringBootMain, details.getSpringBootArgs())
	}
	args := details.getSpringBootCommand()
	return _createJavaCommand(javaCmd, details, args)
}

func _createJavaCommand(javaCmd string, details *RunDetails, args []string) (*exec.Cmd, error) {
	args = append(args, details.MainArgs...)
	cmd := exec.Command(javaCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if details.Dir != "" {
		_, err := os.Stat(details.Dir)
		if err != nil {
			return nil, errors.Wrapf(err, "Working directory %s does not exists or is not a directory", details.Dir)
		}
		cmd.Dir = details.Dir
	}

	return cmd, nil
}

func createOperatorCommand(details *RunDetails) (*exec.Cmd, error) {
	executable := os.Args[0]
	args := details.MainArgs[1:]
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if details.Dir != "" {
		_, err := os.Stat(details.Dir)
		if err != nil {
			return nil, errors.Wrapf(err, "Working directory %s does not exists or is not a directory", details.Dir)
		}
		cmd.Dir = details.Dir
	}

	return cmd, nil
}

func _createBuildPackCommand(_ *RunDetails, className string, args []string) (*exec.Cmd, error) {
	launcher := getBuildpackLauncher()

	// Create the JVM arguments file
	argsFile, err := os.CreateTemp("", "jvm-args")
	if err != nil {
		return nil, err
	}
	defer closeFile(argsFile, log)

	// write the JVM args to the file
	data := strings.Join(args, "\n")
	if _, err := argsFile.WriteString(data); err != nil {
		return nil, err
	}
	log.Info("Created JVM Arguments file", "filename", argsFile.Name(), "data", data)

	cmd := exec.Command(launcher, "java", "@"+argsFile.Name(), className)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd, nil
}

func getBuildpackLauncher() string {
	if launcher, ok := os.LookupEnv(v1.EnvVarCnbpLauncher); ok {
		return launcher
	}
	return v1.DefaultCnbpLauncher
}

func createGraalCommand(details *RunDetails) (*exec.Cmd, error) {
	ex := details.AppType
	args := []string{"--polyglot", "--jvm"}
	args = append(args, details.getCommand()...)
	args = append(args, details.MainClass)
	args = append(args, details.MainArgs...)

	cmd := exec.Command(ex, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if details.Dir != "" {
		_, err := os.Stat(details.Dir)
		if err != nil {
			return nil, errors.Wrapf(err, "Working directory %s does not exists or is not a directory", details.Dir)
		}
		cmd.Dir = details.Dir
	}

	return cmd, nil
}

// Set the Coherence site and rack values
func configureSiteAndRack(details *RunDetails) {
	log.Info("Configuring Coherence site and rack")
	if !details.GetSite {
		return
	}

	var site string

	siteLocation := details.ExpandEnv(details.Getenv(v1.EnvVarCohSite))
	log.Info("Configuring Coherence site", "url", siteLocation)
	if siteLocation != "" {
		switch {
		case strings.ToLower(siteLocation) == "http://":
			site = ""
		case strings.HasPrefix(siteLocation, "http://"):
			// do http get
			site = httpGetWithBackoff(siteLocation, details)
		default:
			st, err := os.Stat(siteLocation)
			if err == nil && !st.IsDir() {
				d, err := os.ReadFile(siteLocation)
				if err != nil {
					site = string(d)
				}
			}
		}
	}

	if site != "" {
		details.addArg("-Dcoherence.site=" + site)
	}

	var rack string

	rackLocation := details.ExpandEnv(details.Getenv(v1.EnvVarCohRack))
	log.Info("Configuring Coherence rack", "url", rackLocation)
	if rackLocation != "" {
		switch {
		case strings.ToLower(rackLocation) == "http://":
			rack = ""
		case strings.HasPrefix(rackLocation, "http://"):
			// do http get
			rack = httpGetWithBackoff(rackLocation, details)
		default:
			st, err := os.Stat(rackLocation)
			if err == nil && !st.IsDir() {
				d, err := os.ReadFile(rackLocation)
				if err != nil {
					rack = string(d)
				}
			}
		}
	}

	if rack != "" {
		details.addArg("-Dcoherence.rack=" + rack)
	} else if site != "" {
		details.addArg("-Dcoherence.rack=" + site)
	}
}

// httpGetWithBackoff does a http get for the specified url with retry back-off for errors.
func httpGetWithBackoff(url string, details *RunDetails) string {
	var backoff time.Duration
	timeout := 120

	val := details.Getenv(v1.EnvVarOperatorTimeout)
	if val != "" {
		t, err := strconv.Atoi(val)
		if err == nil {
			timeout = t
		} else {
			log.Info("Invalid value set for GET request timeout, using default of 120\n", "envVar", v1.EnvVarOperatorTimeout, "value", val)
		}
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	for _, backoff = range backoffSchedule {
		s, status, err := httpGet(url, client)
		if err == nil && status == http.StatusOK {
			return s
		}
		log.Info("http get backoff", "url", url, "backoff", backoff.String(), "status", strconv.Itoa(status), "error", err.Error())
		time.Sleep(backoff)
	}

	// now just retry using the final back-off value for a maximum of five more attempts...
	for i := 0; i < 5; i++ {
		s, status, err := httpGet(url, client)
		if err == nil && status == http.StatusOK {
			return s
		}
		log.Info("http get backoff", "url", url, "backoff", backoff.String(), "status", strconv.Itoa(status), "error", err.Error())
		time.Sleep(backoff)
	}

	log.Info("Unable to perform get request within backoff limit", "url", url)
	return ""
}

// Do a http get for the specified url and return the response body for
// a 200 response or empty string for a non-200 response or error.
func httpGet(url string, client http.Client) (string, int, error) {
	log.Info("Performing http get", "url", url)

	resp, err := client.Get(url)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrapf(err, "failed to get URL %s", url)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, errors.Wrapf(err, "failed to read response body from URL %s", url)
	}

	s := string(body)

	if resp.StatusCode != http.StatusOK {
		log.Info("Did not receive a 200 response from URL", "Status", resp.Status, "Body", s)
	} else {
		log.Info("Received 200 response", "Body", s)
	}

	return s, resp.StatusCode, nil
}

func checkCoherenceVersion(v string, details *RunDetails) bool {
	log.Info("Performing Coherence version check", "version", v)

	if details.isEnvTrue(v1.EnvVarCohSkipVersionCheck) {
		log.Info("Skipping Coherence version check", "envVar", v1.EnvVarCohSkipVersionCheck, "value", details.Getenv(v1.EnvVarCohSkipVersionCheck))
		return true
	}

	// Get the classpath to use (we need Coherence jar)
	cp := details.UtilsDir + "/lib/coherence-operator.jar" + ":" + details.getClasspath()

	var exe string
	var cmd *exec.Cmd
	var args []string

	if details.isBuildPacks() {
		// This is a buildpacks image so use the Buildpacks launcher to run Java
		exe = getBuildpackLauncher()
		args = []string{"java"}
	} else {
		// this should be a normal image with Java available
		exe = details.getJavaExecutable()
	}

	if details.AppType == AppTypeSpring {
		// This is a Spring Boot App so Coherence jar is embedded in the Spring Boot application
		args = append(args, "-Dloader.path="+cp,
			"-Dcoherence.operator.springboot.listener=false",
			"-Dloader.main=com.oracle.coherence.k8s.CoherenceVersion",
			"org.springframework.boot.loader.PropertiesLauncher", v)

		if jar, _ := details.lookupEnv(v1.EnvVarSpringBootFatJar); jar != "" {
			// This is a fat jar Spring boot app so put the fat jar on the classpath
			args = append(args, "-cp", jar)
		}
	} else {
		// We can use normal Java
		args = append(args, "-cp", cp,
			"-Dcoherence.operator.springboot.listener=false",
			"com.oracle.coherence.k8s.CoherenceVersion", v)
	}

	cmd = exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Info("Executing version check", "command", strings.Join(cmd.Args, " "))
	// execute the command
	err := cmd.Run()
	if err == nil {
		// command exited with exit code 0
		log.Info("Executed Coherence version check, version is greater than or equal to expected", "version", v)
		return true
	}
	if _, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0
		log.Info("Executed Coherence version check, version is lower than expected", "version", v)
		return false
	}
	// command exited with some other error
	log.Error(err, "Coherence version check failed")
	return false
}

func cohPre12214(details *RunDetails) {
	details.addArg("-Dcoherence.override=k8s-coherence-nossl-override.xml")
	details.addArgFromEnvVar(v1.EnvVarCohOverride, "-Dcoherence.k8s.override")
}

func cohPost12214(details *RunDetails) {
	details.addArg("-Dcoherence.override=k8s-coherence-override.xml")
	details.addArgFromEnvVar(v1.EnvVarCohOverride, "-Dcoherence.k8s.override")
}

func cohPost2206(details *RunDetails) {
	useOperator, found := os.LookupEnv(v1.EnvVarUseOperatorHealthCheck)
	if found && strings.EqualFold("true", useOperator) {
		details.addArg("-Dcoherence.k8s.operator.health.enabled=true")
	} else {
		details.addArg("-Dcoherence.k8s.operator.health.enabled=false")
		details.setSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohHealthPort, "-Dcoherence.health.http.port", fmt.Sprintf("%d", v1.DefaultHealthPort))
	}
}

func addManagementSSL(details *RunDetails) {
	addSSL(v1.EnvVarCohMgmtPrefix, v1.PortNameManagement, details)
}

func addMetricsSSL(details *RunDetails) {
	addSSL(v1.EnvVarCohMetricsPrefix, v1.PortNameMetrics, details)
}

func addSSL(prefix, prop string, details *RunDetails) {
	var urlPrefix string

	sslCerts := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLCerts)
	if sslCerts != "" {
		if !strings.HasSuffix(sslCerts, "/") {
			sslCerts += "/"
		}
		if strings.HasSuffix(sslCerts, "file:") {
			urlPrefix = sslCerts
		} else {
			urlPrefix = "file:" + sslCerts
		}
	} else {
		urlPrefix = "file:"
	}

	if details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLEnabled) != "" {
		details.addArg("-Dcoherence." + prop + ".http.provider=ManagementSSLProvider")
	}

	ks := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStore)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.keystore=" + urlPrefix + ks)
	}
	kspw := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreCredFile)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.keystore.password=" + urlPrefix + kspw)
	}
	kpw := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyCredFile)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.key.password=" + urlPrefix + kpw)
	}
	kalg := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreAlgo)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.keystore.algorithm=" + urlPrefix + kalg)
	}
	kprov := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreProvider)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.keystore.provider=" + urlPrefix + kprov)
	}
	ktyp := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreType)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.keystore.type=" + urlPrefix + ktyp)
	}

	ts := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStore)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.truststore=" + urlPrefix + ts)
	}
	tspw := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreCredFile)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.truststore.password=" + urlPrefix + tspw)
	}
	talg := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreAlgo)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.truststore.algorithm=" + urlPrefix + talg)
	}
	tprov := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreProvider)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.truststore.provider=" + urlPrefix + tprov)
	}
	ttyp := details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreType)
	if ks != "" {
		details.addArg("-Dcoherence." + prop + ".security.truststore.type=" + urlPrefix + ttyp)
	}

	if details.getenvWithPrefix(prefix, v1.EnvVarSuffixSSLRequireClientCert) != "" {
		details.addArg("-Dcoherence." + prop + ".http.auth=cert")
	}
}

func closeFile(f *os.File, log logr.Logger) {
	err := f.Close()
	if err != nil {
		log.Error(err, "error closing file "+f.Name())
	}
}
