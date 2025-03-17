/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
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
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/runner/run_details"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"k8s.io/apimachinery/pkg/api/resource"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
	"time"
)

// The code that actually starts the process in the Coherence container.

const (
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
func NewRootCommand(env map[string]string, v *viper.Viper) *cobra.Command {
	operator.SetViper(v)

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

	rootCmd.PersistentFlags().Bool(operator.FlagDryRun, false, "Just print information about the commands that would execute")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s.yaml)", defaultConfig))
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	rootCmd.AddCommand(initCommand(env))
	rootCmd.AddCommand(configCommand(env))
	rootCmd.AddCommand(serverCommand())
	rootCmd.AddCommand(consoleCommand(v))
	rootCmd.AddCommand(queryPlusCommand(v))
	rootCmd.AddCommand(statusCommand())
	rootCmd.AddCommand(readyCommand())
	rootCmd.AddCommand(nodeCommand())
	rootCmd.AddCommand(operatorCommand(v))
	rootCmd.AddCommand(networkTestCommand())
	rootCmd.AddCommand(jShellCommand(v))
	rootCmd.AddCommand(sleepCommand(v))

	return rootCmd
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
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
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
	parent := cmd.Parent()
	if parent != nil {
		_ = v.BindPFlags(cmd.Parent().Flags())
	}
	_ = v.BindPFlags(cmd.PersistentFlags())
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
	return ExecuteWithArgsAndViper(nil, nil, viper.GetViper())
}

// ExecuteWithArgsAndNewViper runs the runner with a given environment and argument overrides.
func ExecuteWithArgsAndNewViper(env map[string]string, args []string) (Execution, error) {
	v := viper.New()
	for key, value := range env {
		v.SetDefault(key, value)
	}
	return ExecuteWithArgsAndViper(env, args, v)
}

// ExecuteWithArgsAndViper runs the runner with a given environment and argument overrides.
func ExecuteWithArgsAndViper(env map[string]string, args []string, v *viper.Viper) (Execution, error) {
	cmd := NewRootCommand(env, v)

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
type RunFunction func(*run_details.RunDetails, *cobra.Command)

// MaybeRunFunction is a function to maybe run a command depending on the return bool
type MaybeRunFunction func(*run_details.RunDetails, *cobra.Command) (bool, error)

// always is a wrapper around a RunFunction to turn it into a MaybeFunction that always runs
type always struct {
	Fn RunFunction
}

// run will wrap a RunFunction and always return true
func (in always) run(details *run_details.RunDetails, cmd *cobra.Command) (bool, error) {
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

	details := run_details.NewRunDetails(e.V, log)
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

			dryRun := operator.IsDryRun()
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

// configure the command details.
func configureCommand(details *run_details.RunDetails) error {
	var err error

	// Set standard system properties
	details.AddSystemPropertyArg(v1.SysPropCoherenceTTL, "0")
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohWka, v1.SysPropCoherenceWKA)
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohMachineName, v1.SysPropCoherenceMachine)
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohMemberName, v1.SysPropCoherenceMember)
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohClusterName, v1.SysPropCoherenceCluster)
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohCacheConfig, v1.SysPropCoherenceCacheConfig)
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohIdentity, v1.SysPropOperatorIdentity)
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohForceExit, v1.SysPropOperatorForceExit)
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohHealthPort, v1.SysPropOperatorHealthPort, fmt.Sprintf("%d", v1.DefaultHealthPort))
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohMgmtPrefix+v1.EnvVarCohPortSuffix, v1.SysPropCoherenceManagementHttpPort, fmt.Sprintf("%d", v1.DefaultManagementPort))
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohMetricsPrefix+v1.EnvVarCohPortSuffix, v1.SysPropCoherenceMetricsHttpPort, fmt.Sprintf("%d", v1.DefaultMetricsPort))

	details.AddVMOption(v1.JvmOptUnlockDiagnosticVMOptions)

	// Configure the classpath to support images created with the JIB Maven plugin
	// This is enabled by default unless the image is a buildpacks image, or we
	// are running a Spring Boot application.
	if !details.IsBuildPacks() && !details.IsSpringBoot() && details.IsEnvTrueOrBlank(v1.EnvVarJvmClasspathJib) {
		appDir := details.GetenvOrDefault(v1.EnvVarCohAppDir, "/app")
		cpFile := filepath.Join(appDir, "jib-classpath-file")
		fi, e := os.Stat(cpFile)
		if e != nil && appDir != "/app" {
			// try in /app
			cpFile = "/app/jib-classpath-file"
			fi, e = os.Stat(cpFile)
		}

		if e == nil && (fi.Size() != 0) {
			clsPath, _ := readFirstLineFromFile(cpFile)
			if len(clsPath) != 0 {
				details.AddClasspath(clsPath)
			}
		} else {
			details.AddClasspathIfExists(appDir + "/resources")
			details.AddClasspathIfExists(appDir + "/classes")
			details.AddJarsToClasspath(appDir + "/classpath")
			details.AddJarsToClasspath(appDir + "/libs")
		}
	}

	// Add the Operator Utils jar to the classpath
	details.AddClasspath(details.UtilsDir + v1.OperatorJarFileSuffix)
	details.AddClasspathIfExists(details.UtilsDir + v1.OperatorConfigDirSuffix)

	// Configure Coherence persistence
	mode := details.GetenvOrDefault(v1.EnvVarCohPersistenceMode, "on-demand")
	details.AddSystemPropertyArg(v1.SysPropCoherencePersistenceMode, mode)

	persistence := details.Getenv(v1.EnvVarCohPersistenceDir)
	if persistence != "" {
		details.AddSystemPropertyArg(v1.SysPropCoherencePersistenceBaseDir, persistence)
	}

	snapshots := details.Getenv(v1.EnvVarCohSnapshotDir)
	if snapshots != "" {
		details.AddSystemPropertyArg(v1.SysPropCoherencePersistenceSnapshotDir, snapshots)
	}

	// Set the Coherence site and rack values
	configureSiteAndRack(details)

	// Set the Coherence log level
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohLogLevel, v1.SysPropCoherenceLogLevel)

	// Disable IPMonitor
	ipMon := details.Getenv(v1.EnvVarEnableIPMonitor)
	if ipMon != "TRUE" {
		details.AddSystemPropertyArg(v1.SysPropCoherenceIpMonitor, "0")
	}

	details.AddSystemPropertyArg(v1.SysPropCoherenceOverride, "k8s-coherence-override.xml")
	details.AddSystemPropertyFromEnvVar(v1.EnvVarCohOverride, v1.SysPropOperatorOverride)

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
		details.AddArg("-Dcoherence.operator.statusha.allowendangered=" + allowEndangered)
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
			return err
		}
		if err = os.MkdirAll(jvmDir+"/jfr", os.ModePerm); err != nil {
			return err
		}
		if err = os.MkdirAll(jvmDir+"/heap-dumps", os.ModePerm); err != nil {
			return err
		}
	}

	details.AddArg(fmt.Sprintf("-Dcoherence.operator.diagnostics.dir=%s", jvmDir))
	details.AddVMOption(fmt.Sprintf("-XX:HeapDumpPath=%s/heap-dumps/%s-%s.hprof", jvmDir, member, podUID))

	// set the flag that allows the operator to resume suspended services on start-up
	if !details.IsEnvTrueOrBlank(v1.EnvVarOperatorAllowResume) {
		details.AddArg("-Dcoherence.operator.can.resume.services=false")
	} else {
		details.AddArg("-Dcoherence.operator.can.resume.services=true")
	}

	if svc := details.Getenv(v1.EnvVarOperatorResumeServices); svc != "" {
		details.AddArg("-Dcoherence.operator.resume.services=base64:" + svc)
	}

	gc := strings.ToLower(details.Getenv(v1.EnvVarJvmGcCollector))
	switch {
	case gc == "" || gc == "g1":
		details.AddMemoryOption("-XX:+UseG1GC")
	case gc == "cms":
		details.AddMemoryOption("-XX:+UseConcMarkSweepGC")
	case gc == "parallel":
		details.AddMemoryOption("-XX:+UseParallelGC")
	}

	maxRAM := details.Getenv(v1.EnvVarJvmMaxRAM)
	if maxRAM != "" {
		details.AddMemoryOption("-XX:MaxRAM=" + maxRAM)
	}

	heap := details.Getenv(v1.EnvVarJvmMemoryHeap)
	if heap != "" {
		// if heap is set use it
		details.AddMemoryOption("-XX:InitialHeapSize=" + heap)
		details.AddMemoryOption("-XX:MaxHeapSize=" + heap)
	} else {
		// if heap is not set check whether the individual heap values are set
		initialHeap := details.Getenv(v1.EnvVarJvmMemoryInitialHeap)
		if initialHeap != "" {
			details.AddMemoryOption("-XX:InitialHeapSize=" + initialHeap)
		}
		maxHeap := details.Getenv(v1.EnvVarJvmMemoryMaxHeap)
		if maxHeap != "" {
			details.AddMemoryOption("-XX:MaxHeapSize=" + maxHeap)
		}
	}

	percentageHeap := details.Getenv(v1.EnvVarJvmRAMPercentage)
	if percentageHeap != "" {
		// the heap percentage is set so use it
		q, err := resource.ParseQuantity(percentageHeap)
		if err == nil {
			d := q.AsDec()
			details.AddMemoryOption("-XX:InitialRAMPercentage=" + d.String())
			details.AddMemoryOption("-XX:MinRAMPercentage=" + d.String())
			details.AddMemoryOption("-XX:MaxRAMPercentage=" + d.String())
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
				details.AddMemoryOption("-XX:InitialRAMPercentage=" + d.String())
			} else {
				log.Info("ERROR: InitialRAMPercentage is not a valid resource.Quantity", "Value", initial, "Error", err.Error())
				os.Exit(1)
			}
		}

		maxRam := details.Getenv(v1.EnvVarJvmMaxRAMPercentage)
		if maxRam != "" {
			q, err := resource.ParseQuantity(maxRam)
			if err == nil {
				d := q.AsDec()
				details.AddMemoryOption("-XX:MaxRAMPercentage=" + d.String())
			} else {
				log.Info("ERROR: MaxRAMPercentage is not a valid resource.Quantity", "Value", maxRam, "Error", err.Error())
				os.Exit(1)
			}
		}

		minRam := details.Getenv(v1.EnvVarJvmMinRAMPercentage)
		if minRam != "" {
			q, err := resource.ParseQuantity(minRam)
			if err == nil {
				d := q.AsDec()
				details.AddMemoryOption("-XX:MinRAMPercentage=" + d.String())
			} else {
				log.Info("ERROR: MinRAMPercentage is not a valid resource.Quantity", "Value", minRam, "Error", err.Error())
				os.Exit(1)
			}
		}
	}

	direct := details.Getenv(v1.EnvVarJvmMemoryDirect)
	if direct != "" {
		details.AddMemoryOption("-XX:MaxDirectMemorySize=" + direct)
	}
	stack := details.Getenv(v1.EnvVarJvmMemoryStack)
	if stack != "" {
		details.AddMemoryOption("-Xss" + stack)
	}
	meta := details.Getenv(v1.EnvVarJvmMemoryMeta)
	if meta != "" {
		details.AddMemoryOption("-XX:MetaspaceSize=" + meta)
		details.AddMemoryOption("-XX:MaxMetaspaceSize=" + meta)
	}
	track := details.GetenvOrDefault(v1.EnvVarJvmMemoryNativeTracking, "summary")
	if track != "" {
		details.AddDiagnosticOption("-XX:NativeMemoryTracking=" + track)
		details.AddDiagnosticOption("-XX:+PrintNMTStatistics")
	}

	// Configure debugging
	debugArgs := ""
	if details.IsEnvTrue(v1.EnvVarJvmDebugEnabled) {
		var suspend string
		if details.IsEnvTrue(v1.EnvVarJvmDebugSuspended) {
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

	details.AddVMOption(fmt.Sprintf("-XX:ErrorFile=%s/hs-err-%s-%s.log", jvmDir, member, podUID))

	if details.IsEnvTrueOrBlank(v1.EnvVarJvmOomHeapDump) {
		details.AddVMOption("-XX:+HeapDumpOnOutOfMemoryError")
	}

	if details.IsEnvTrueOrBlank(v1.EnvVarJvmOomExit) {
		details.AddVMOption("-XX:+ExitOnOutOfMemoryError")
	}

	// Use JVM container support
	if details.IsEnvTrueOrBlank(v1.EnvVarJvmUseContainerLimits) {
		details.AddVMOption("-XX:+UseContainerSupport")
	}

	details.AddArgs(debugArgs)

	gcArgs := details.Getenv(v1.EnvVarJvmGcArgs)
	if gcArgs != "" {
		details.AddArgs(strings.Split(gcArgs, " ")...)
	}

	jvmArgs := details.Getenv(v1.EnvVarJvmArgs)
	if jvmArgs != "" {
		details.AddArgs(strings.Split(jvmArgs, " ")...)
	}

	extraJvmArgs := operator.GetExtraJvmArgs()
	if extraJvmArgs != nil {
		details.AddArgs(extraJvmArgs...)
	}

	return nil
}

// create the process to execute.
func createCommand(details *run_details.RunDetails) (string, *exec.Cmd, error) {
	var err error
	var cmd *exec.Cmd
	var app string

	switch {
	case details.AppType == v1.AppTypeNone || details.AppType == v1.AppTypeJava:
		app = "Java"
		cmd, err = createJavaCommand(details.GetJavaExecutable(), details)
	case details.IsSpringBoot():
		app = "SpringBoot"
		cmd, err = createSpringBootCommand(details.GetJavaExecutable(), details)
	case details.AppType == v1.AppTypeHelidon:
		app = "Java"
		cmd, err = createJavaCommand(details.GetJavaExecutable(), details)
	case details.AppType == v1.AppTypeCoherence:
		app = "Java"
		cmd, err = createJavaCommand(details.GetJavaExecutable(), details)
	case details.AppType == v1.AppTypeJShell:
		app = "JShell"
		cmd, err = createJShellCommand(details.GetJShellExecutable(), details)
	case details.AppType == v1.AppTypeOperator:
		app = "Operator"
		cmd, err = createOperatorCommand(details)
	default:
		app = "Graal (" + details.AppType + ")"
		cmd, err = createGraalCommand(details)
	}

	extraEnv := operator.GetExtraEnvVars()
	if cmd != nil && extraEnv != nil {
		cmd.Env = append(cmd.Env, extraEnv...)
	}

	return app, cmd, err
}

func createJavaCommand(javaCmd string, details *run_details.RunDetails) (*exec.Cmd, error) {
	args := details.GetCommand()
	args = append(args, details.MainClass)
	return _createJavaCommand(javaCmd, details, args)
}

func createJShellCommand(jshellCmd string, details *run_details.RunDetails) (*exec.Cmd, error) {
	args := details.GetCommandWithPrefix("-R", "-J")
	return _createJavaCommand(jshellCmd, details, args)
}

func readFirstLineFromFile(path string) (string, error) {
	file, err := os.Open(maybeStripFileScheme(path))
	if err != nil {
		return "", err
	}
	defer closeFile(file, log)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if len(text) == 0 {
		return "", nil
	}
	return text[0], nil
}

func createSpringBootCommand(javaCmd string, details *run_details.RunDetails) (*exec.Cmd, error) {
	if details.IsBuildPacks() {
		if details.AppType == v1.AppTypeSpring2 {
			return _createBuildPackCommand(details, v1.SpringBootMain2, details.GetSpringBootArgs())
		}
		return _createBuildPackCommand(details, v1.SpringBootMain3, details.GetSpringBootArgs())
	}
	args := details.GetSpringBootCommand()
	return _createJavaCommand(javaCmd, details, args)
}

func _createJavaCommand(javaCmd string, details *run_details.RunDetails, args []string) (*exec.Cmd, error) {
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

func createOperatorCommand(details *run_details.RunDetails) (*exec.Cmd, error) {
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

func _createBuildPackCommand(_ *run_details.RunDetails, className string, args []string) (*exec.Cmd, error) {
	launcher := run_details.GetBuildpackLauncher()

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

func createGraalCommand(details *run_details.RunDetails) (*exec.Cmd, error) {
	ex := details.AppType
	args := []string{"--polyglot", "--jvm"}
	args = append(args, details.GetCommand()...)
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
func configureSiteAndRack(details *run_details.RunDetails) {
	var err error
	if !details.GetSite {
		return
	}

	log.Info("Configuring Coherence site and rack")

	site := details.Getenv(v1.EnvVarCoherenceSite)
	if site == "" {
		siteLocation := details.ExpandEnv(details.Getenv(v1.EnvVarCohSite))
		log.Info("Configuring Coherence site", "url", siteLocation)
		if siteLocation != "" {
			switch {
			case strings.ToLower(siteLocation) == "http://":
				site = ""
			case strings.HasPrefix(siteLocation, "http://"):
				// do http get
				site = httpGetWithBackoff(siteLocation, details)
			case strings.HasPrefix(siteLocation, "https://"):
				// https not supported
				log.Info("Cannot read site URI, https is not supported", "URI", siteLocation)
			default:
				site, err = readFirstLineFromFile(siteLocation)
				if err != nil {
					log.Error(err, "error reading site info", "Location", siteLocation)
				}
			}
		}

		if site != "" {
			details.AddSystemPropertyArg(v1.SysPropCoherenceSite, site)
		}
	} else {
		expanded := details.ExpandEnv(site)
		if expanded != site {
			log.Info("Coherence site property set from expanded "+v1.EnvVarCoherenceSite+" environment variable", v1.EnvVarCoherenceSite, site, "Site", expanded)
			site = expanded
			if strings.TrimSpace(site) != "" {
				details.AddSystemPropertyArg(v1.SysPropCoherenceSite, site)
			}
		} else {
			details.AddSystemPropertyArg(v1.SysPropCoherenceSite, site)
		}
	}

	rack := details.Getenv(v1.EnvVarCoherenceRack)
	if rack == "" {
		rackLocation := details.ExpandEnv(details.Getenv(v1.EnvVarCohRack))
		log.Info("Configuring Coherence rack", "url", rackLocation)
		if rackLocation != "" {
			switch {
			case strings.ToLower(rackLocation) == "http://":
				rack = ""
			case strings.HasPrefix(rackLocation, "http://"):
				// do http get
				rack = httpGetWithBackoff(rackLocation, details)
			case strings.HasPrefix(rackLocation, "https://"):
				// https not supported
				log.Info("Cannot read rack URI, https is not supported", "URI", rackLocation)
			default:
				rack, err = readFirstLineFromFile(rackLocation)
				if err != nil {
					log.Error(err, "error reading site info", "Location", rackLocation)
				}
			}
		}

		if rack != "" {
			details.AddSystemPropertyArg(v1.SysPropCoherenceRack, rack)
		} else if site != "" {
			details.AddSystemPropertyArg(v1.SysPropCoherenceRack, site)
		}
	} else {
		expanded := details.ExpandEnv(rack)
		if expanded != rack {
			log.Info("Coherence site property set from expanded "+v1.EnvVarCoherenceRack+" environment variable", v1.EnvVarCoherenceRack, rack, "Rack", expanded)
			rack = expanded
			if len(rack) == 0 {
				// if the expanded COHERENCE_RACK value is blank then set rack to site as
				// the rack cannot be blank if site is set
				rack = site
			}
			if strings.TrimSpace(rack) != "" {
				details.AddSystemPropertyArg(v1.SysPropCoherenceRack, rack)
			}
		} else {
			details.AddSystemPropertyArg(v1.SysPropCoherenceRack, rack)
		}
	}
}

func maybeStripFileScheme(uri string) string {
	if strings.HasPrefix(uri, "file://") {
		return strings.TrimPrefix(uri, "file://")
	}
	return uri
}

// httpGetWithBackoff does a http get for the specified url with retry back-off for errors.
func httpGetWithBackoff(url string, details *run_details.RunDetails) string {
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
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}
		log.Info("http get backoff", "url", url, "backoff", backoff.String(), "status", strconv.Itoa(status), "error", errorMsg)
		time.Sleep(backoff)
	}

	// now just retry using the final back-off value for a maximum of five more attempts...
	for i := 0; i < 5; i++ {
		s, status, err := httpGet(url, client)
		if err == nil && status == http.StatusOK {
			return s
		}
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}
		log.Info("http get backoff", "url", url, "backoff", backoff.String(), "status", strconv.Itoa(status), "error", errorMsg)
		time.Sleep(backoff)
	}

	log.Info("Unable to perform get request within backoff limit", "url", url)
	return ""
}

// Do a http get for the specified url and return the response body for
// a 200 response or empty string for a non-200 response or error.
func httpGet(urlString string, client http.Client) (string, int, error) {
	log.Info("Performing http get", "url", urlString)

	u, err := url.Parse(urlString)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrapf(err, "failed to parse URL %s", urlString)
	}

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrapf(err, "failed to create request for URL %s", urlString)
	}

	req.Host = u.Host

	h := http.Header{}
	h.Set("Host", u.Host)
	h.Set("User-Agent", fmt.Sprintf("coherence-operator-runner/%s", operator.GetVersion()))
	req.Header = h

	resp, err := client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrapf(err, "failed to get URL %s", urlString)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, errors.Wrapf(err, "failed to read response body from URL %s", urlString)
	}

	s := string(body)

	if resp.StatusCode != http.StatusOK {
		log.Info("Did not receive a 200 response from URL", "Status", resp.Status, "Body", s)
	} else {
		log.Info("Received 200 response", "Body", s)
	}

	return s, resp.StatusCode, nil
}

func checkCoherenceVersion(v string, details *run_details.RunDetails) bool {
	log.Info("Performing Coherence version check", "version", v)

	if details.IsEnvTrue(v1.EnvVarCohSkipVersionCheck) {
		log.Info("Skipping Coherence version check", "envVar", v1.EnvVarCohSkipVersionCheck, "value", details.Getenv(v1.EnvVarCohSkipVersionCheck))
		return true
	}

	// Get the classpath to use (we need Coherence jar)
	cp := details.GetClasspath()

	var exe string
	var cmd *exec.Cmd
	var args []string

	if details.IsBuildPacks() {
		// This is a build-packs image so use the Build-packs launcher to run Java
		exe = run_details.GetBuildpackLauncher()
		args = []string{exe}
	} else {
		// this should be a normal image with Java available
		exe = details.GetJavaExecutable()
	}

	if details.IsSpringBoot() {
		// This is a Spring Boot App so Coherence jar is embedded in the Spring Boot application
		cp := strings.ReplaceAll(cp, ":", ",")
		args = append(args, "-Dloader.path="+cp,
			"-Dcoherence.operator.springboot.listener=false",
			"-Dloader.main=com.oracle.coherence.k8s.CoherenceVersion")

		if jar, _ := details.LookupEnv(v1.EnvVarSpringBootFatJar); jar != "" {
			// This is a fat jar Spring boot app so put the fat jar on the classpath
			args = append(args, v1.JvmOptClassPath, jar)
		}

		if details.AppType == v1.AppTypeSpring2 {
			// we are running SpringBoot 2.x
			args = append(args, v1.SpringBootMain2, v)
		} else {
			// we are running SpringBoot 3.x
			args = append(args, v1.SpringBootMain3, v)
		}
	} else {
		// We can use normal Java
		args = append(args, v1.JvmOptClassPath, cp,
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
	if exitError, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code == 99 which means the version is lower than requested
		if exitError.ExitCode() == 99 {
			log.Info("Executed Coherence version check, version is lower than expected", "version", v)
			return false
		}
	}
	// command exited with some other error, assume the version is good
	log.Info("Coherence version check failed, assuming version is valid", "version", v, "error", err.Error())
	return true
}

func cohPost2206(details *run_details.RunDetails) {
	if details.UseOperatorHealth {
		details.AddArg("-Dcoherence.operator.health.enabled=true")
	} else {
		useOperator := details.GetenvOrDefault(v1.EnvVarUseOperatorHealthCheck, "false")
		if strings.EqualFold("true", useOperator) {
			details.AddSystemPropertyArg(v1.SysPropOperatorHealthEnabled, "true")
		} else {
			details.AddSystemPropertyArg(v1.SysPropOperatorHealthEnabled, "false")
			details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohHealthPort, v1.SysPropCoherenceHealthHttpPort, fmt.Sprintf("%d", v1.DefaultHealthPort))
		}
	}
}

func addManagementSSL(details *run_details.RunDetails) {
	addSSL(v1.EnvVarCohMgmtPrefix, v1.PortNameManagement, details)
}

func addMetricsSSL(details *run_details.RunDetails) {
	addSSL(v1.EnvVarCohMetricsPrefix, v1.PortNameMetrics, details)
}

func addSSL(prefix, prop string, details *run_details.RunDetails) {
	var urlPrefix string

	sslCerts := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLCerts)
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

	if details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLEnabled) != "" {
		details.AddArg("-Dcoherence." + prop + ".http.provider=ManagementSSLProvider")
	}

	ks := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStore)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.keystore=" + urlPrefix + ks)
	}
	kspw := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreCredFile)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.keystore.password=" + urlPrefix + kspw)
	}
	kpw := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyCredFile)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.key.password=" + urlPrefix + kpw)
	}
	kalg := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreAlgo)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.keystore.algorithm=" + urlPrefix + kalg)
	}
	kprov := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreProvider)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.keystore.provider=" + urlPrefix + kprov)
	}
	ktyp := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLKeyStoreType)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.keystore.type=" + urlPrefix + ktyp)
	}

	ts := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStore)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.truststore=" + urlPrefix + ts)
	}
	tspw := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreCredFile)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.truststore.password=" + urlPrefix + tspw)
	}
	talg := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreAlgo)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.truststore.algorithm=" + urlPrefix + talg)
	}
	tprov := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreProvider)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.truststore.provider=" + urlPrefix + tprov)
	}
	ttyp := details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLTrustStoreType)
	if ks != "" {
		details.AddArg("-Dcoherence." + prop + ".security.truststore.type=" + urlPrefix + ttyp)
	}

	if details.GetenvWithPrefix(prefix, v1.EnvVarSuffixSSLRequireClientCert) != "" {
		details.AddArg("-Dcoherence." + prop + ".http.auth=cert")
	}
}

func closeFile(f *os.File, log logr.Logger) {
	err := f.Close()
	if err != nil {
		log.Error(err, "error closing file "+f.Name())
	}
}

func addEnvVarFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSlice(
		operator.FlagEnvVar,
		nil,
		"Additional environment variables to pass to the process",
	)
}

func addJvmArgFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSlice(
		operator.FlagJvmArg,
		nil,
		"AdditionalJVM args to pass to the process",
	)
}

func setupFlags(cmd *cobra.Command, v *viper.Viper) {
	// enable using dashed notation in flags and underscores in env
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		setupLog.Error(err, "binding flags")
		os.Exit(1)
	}

	v.AutomaticEnv()
}
