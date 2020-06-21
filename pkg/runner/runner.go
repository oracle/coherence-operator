/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/resource"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// The code that actually starts the process in the Coherence container.

const (
	DCS         = "com.tangosol.net.DefaultCacheServer"
	HelidonMain = "io.helidon.microprofile.server.Main"
	ServerMain  = "com.oracle.coherence.k8s.Main"

	CommandServer      = "server"
	CommandConsole     = "console"
	CommandQueryPlus   = "queryplus"
	CommandMBeanServer = "mbeanserver"
	CommandVersion     = "version"

	EnvVarBuildInfo = "OPERATOR_BUILD_INFO"
	EnvVarVersion   = "OPERATOR_VERSION"
	EnvVarGitCommit = "OPERATOR_GIT_COMMIT"
	EnvVarBuildDate = "OPERATOR_BUILD_DATE"

	AppTypeNone      = ""
	AppTypeJava      = "java"
	AppTypeCoherence = "coherence"
	AppTypeHelidon   = "helidon"
	AppTypeTest      = "op-test"
)

// Run the Coherence process using the specified args and environment variables.
func Run(args []string, env map[string]string) error {
	app, cmd, err := DryRun(args, env)
	if err != nil {
		return err
	}
	if cmd != nil {
		fmt.Printf("\nINFO: Starting the Coherence %s process using:\n", app)
		fmt.Printf("INFO: %s %s\n\n", cmd.Path, strings.Join(cmd.Args, " "))
		return cmd.Run()
	}
	return nil
}

// Build the command to the Coherence process using the specified args and environment variables
// but do not actually start it.
func DryRun(args []string, env map[string]string) (string, *exec.Cmd, error) {
	skipSite := env[v1.EnvVarCohSkipSite]
	details := &RunDetails{
		OsArgs:        args,
		Env:           env,
		CoherenceHome: env[v1.EnvVarCoherenceHome],
		UtilsDir:      env[v1.EnvVarCohUtilDir],
		JavaHome:      env[v1.EnvVarJavaHome],
		AppType:       strings.ToLower(env[v1.EnvVarAppType]),
		Dir:           env[v1.EnvVarCohAppDir],
		MainClass:     DCS,
		GetSite:       strings.ToLower(skipSite) != "true",
	}

	printHeader(details)

	// add any Classpath items
	details.AddClasspath(env[v1.EnvVarJvmExtraClasspath])
	details.AddClasspath(env[v1.EnvVarJavaClasspath])

	if len(details.OsArgs) == 1 {
		details.Command = CommandServer
	} else {
		switch details.OsArgs[1] {
		case CommandServer:
			server(details)
		case CommandConsole:
			console(details)
		case CommandQueryPlus:
			queryPlus(details)
		case CommandMBeanServer:
			mbeanServer(details)
		case CommandVersion:
			return "", nil, nil
		default:
			usage()
			return "", nil, fmt.Errorf("invalid command %s", details.OsArgs[1])
		}
	}

	return start(details)
}

// Configure the runner to run a Coherence Server
func server(details *RunDetails) {
	details.Command = CommandServer
	details.MainClass = ServerMain

	// If the main class environment variable is set then use that
	// otherwise run Coherence DCS.
	mc, found := details.LookupEnv(v1.EnvVarAppMainClass)
	switch {
	case found:
		details.MainArgs = []string{mc}
	case !found && details.AppType == AppTypeCoherence:
		details.MainArgs = []string{DCS}
	case !found && details.AppType == AppTypeHelidon:
		details.MainArgs = []string{HelidonMain}
	default:
		details.MainArgs = []string{DCS}
	}

	// Check for any main class arguments
	ma, found := details.LookupEnv(v1.EnvVarAppMainArgs)
	if found {
		if ma != "" {
			details.MainArgs = append(details.MainArgs, strings.Split(ma, " ")...)
		}
	}

	// Configure the Coherence member's role
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohRole, "-Dcoherence.role", "storage")
	// Configure whether this member is storage enabled
	details.AddArgFromEnvVar(v1.EnvVarCohStorage, "-Dcoherence.distributed.localstorage")

	// Configure Coherence Tracing
	ratio := details.Getenv(v1.EnvVarCohTracingRatio)
	if ratio != "" {
		q, err := resource.ParseQuantity(ratio)
		if err == nil {
			d := q.AsDec()
			details.AddArg("-Dcoherence.tracing.ratio=" + d.String())
		} else {
			fmt.Printf("ERROR: Coherence tracing ratio \"%s\" is invalid - %s\n", ratio, err.Error())
			os.Exit(1)
		}
	}

	// Configure whether Coherence management is enabled
	hasMgmt := details.IsEnvTrue(v1.EnvVarCohMgmtPrefix + v1.EnvVarCohEnabledSuffix)
	fmt.Printf("INFO: Coherence Management over REST (%s%s=%t)\n", v1.EnvVarCohMgmtPrefix, v1.EnvVarCohEnabledSuffix, hasMgmt)
	if hasMgmt {
		fmt.Println("INFO: Configuring Coherence Management over REST")
		details.AddArg("-Dcoherence.management.http=all")
		if details.CoherenceHome != "" {
			// If management is enabled and the COHERENCE_HOME environment variable is set
			// then $COHERENCE_HOME/lib/coherence-management.jar will be added to the classpath
			details.AddClasspath(details.CoherenceHome + "/lib/coherence-management.jar")
		}
	}

	// Configure whether Coherence metrics is enabled
	hasMetrics := details.IsEnvTrue(v1.EnvVarCohMetricsPrefix + v1.EnvVarCohEnabledSuffix)
	fmt.Printf("INFO: Coherence Metrics (%s%s=%t)\n", v1.EnvVarCohMetricsPrefix, v1.EnvVarCohEnabledSuffix, hasMgmt)
	if hasMetrics {
		details.AddArg("-Dcoherence.metrics.http.enabled=true")
		fmt.Println("INFO: Configuring Coherence Metrics")
		if details.CoherenceHome != "" {
			// If metrics is enabled and the COHERENCE_HOME environment variable is set
			// then $COHERENCE_HOME/lib/coherence-metrics.jar will be added to the classpath
			details.AddClasspath(details.CoherenceHome + "/lib/coherence-metrics.jar")
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
				details.AddClasspath(dm + "/*")
			}
		}
	}

	details.AddArg("-XshowSettings:all")
	details.AddArg("-XX:+PrintCommandLineFlags")
	details.AddArg("-XX:+PrintFlagsFinal")

	// Add GC logging parameters if required
	if details.IsEnvTrue(v1.EnvVarJvmGcLogging) {
		details.AddArg("-verbose:gc")
		details.AddArg("-XX:+PrintGCDetails")
		details.AddArg("-XX:+PrintGCTimeStamps")
		details.AddArg("-XX:+PrintHeapAtGC")
		details.AddArg("-XX:+PrintTenuringDistribution")
		details.AddArg("-XX:+PrintGCApplicationStoppedTime")
		details.AddArg("-XX:+PrintGCApplicationConcurrentTime")
	}
}

// Configure the runner to run a Coherence CacheFactory console
func console(details *RunDetails) {
	details.Command = CommandConsole
	details.AppType = AppTypeJava
	details.MainClass = "com.tangosol.net.CacheFactory"
	details.AddArg("-Dcoherence.distributed.localstorage=false")
	details.Setenv(v1.EnvVarCohRole, "console")
	details.Unsetenv(v1.EnvVarJvmMemoryHeap)
	if len(details.OsArgs) > 2 {
		details.MainArgs = details.OsArgs[2:]
	}
}

// Configure the runner to run a Coherence Query Plus console
func queryPlus(details *RunDetails) {
	details.Command = CommandQueryPlus
	details.AppType = AppTypeJava
	details.MainClass = "com.tangosol.coherence.dslquery.QueryPlus"
	if len(details.OsArgs) > 2 {
		details.MainArgs = details.OsArgs[2:]
	}
	details.AddArg("-Dcoherence.distributed.localstorage=false")
	details.Setenv(v1.EnvVarCohRole, "queryPlus")
	details.Unsetenv(v1.EnvVarJvmMemoryHeap)
}

// Configure the runner to run a JMXMP MBean server
func mbeanServer(details *RunDetails) {
	details.Command = CommandMBeanServer
	details.AppType = AppTypeJava
	details.AddClasspath(details.UtilsDir + "/lib/*")
	details.MainClass = "com.oracle.coherence.k8s.JmxmpServer"
	details.MainArgs = []string{}
	details.Setenv(v1.EnvVarJvmJmxmpEnabled, "true")
	details.Setenv(v1.EnvVarCohRole, "MBeanServer")
	details.AddArg("-Dcoherence.distributed.localstorage=false")
	details.AddArg("-Dcoherence.management=all")
	details.AddArg("-Dcoherence.management.remote=true")
	details.AddArg("-Dcom.sun.management.jmxremote.ssl=false")
	details.AddArg("-Dcom.sun.management.jmxremote.authenticate=false")
}

// Start the required process
func start(details *RunDetails) (string, *exec.Cmd, error) {
	// Set standard system properties
	details.AddArgFromEnvVar(v1.EnvVarCohWka, "-Dcoherence.wka")
	details.AddArgFromEnvVar(v1.EnvVarCohMachineName, "-Dcoherence.machine")
	details.AddArgFromEnvVar(v1.EnvVarCohMemberName, "-Dcoherence.member")
	details.AddArgFromEnvVar(v1.EnvVarCohClusterName, "-Dcoherence.cluster")
	details.AddArgFromEnvVar(v1.EnvVarCohCacheConfig, "-Dcoherence.cacheconfig")
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohHealthPort, "-Dcoherence.health.port", fmt.Sprintf("%d", v1.DefaultHealthPort))
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohMgmtPrefix+v1.EnvVarCohPortSuffix, "-Dcoherence.management.http.port", fmt.Sprintf("%d", v1.DefaultManagementPort))
	details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarCohMetricsPrefix+v1.EnvVarCohPortSuffix, "-Dcoherence.metrics.http.port", fmt.Sprintf("%d", v1.DefaultMetricsPort))

	details.AddArg("-XX:+UnlockDiagnosticVMOptions")

	// Configure the classpath to support images created with the JIB Maven plugin
	// This is enabled by default.
	if details.IsEnvTrueOrBlank(v1.EnvVarJvmClasspathJib) {
		details.AddClasspath("/app/resources")
		details.AddClasspath("/app/classes")
		details.AddClasspath("/app/classpath/*")
		details.AddClasspath("/app/libs/*")
	}

	// Add the Operator Utils jar to the classpath
	details.AddClasspath(details.UtilsDir + "/lib/coherence-utils.jar")
	details.AddClasspath(details.UtilsDir + "/config")

	// Configure Coherence persistence
	mode := details.GetenvOrDefault(v1.EnvVarCohPersistenceMode, "on-demand")
	details.AddArg("-Dcoherence.distributed.persistence-mode=" + mode)

	persistence := details.Getenv(v1.EnvVarCohPersistenceDir)
	if persistence != "" {
		details.AddArg("-Dcoherence.distributed.persistence.base.dir=" + persistence)
	}

	snapshots := details.Getenv(v1.EnvVarCohSnapshotDir)
	if snapshots != "" {
		details.AddArg("-Dcoherence.distributed.persistence.snapshot.dir=" + snapshots)
	}

	// Set the Coherence site and rack values
	configureSiteAndRack(details)

	// Set the Coherence log level
	details.AddArgFromEnvVar(v1.EnvVarCohLogLevel, "-Dcoherence.log.level")

	// Do the Coherence version specific configuration
	ok, err := checkCoherenceVersion("12.2.1.4.0", details)
	if err != nil {
		return "", nil, err
	}
	if ok {
		// is at least 12.2.1.4
		cohPost12214(details)
	} else {
		// is at pre-12.2.1.4
		cohPre12214(details)
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

	// Configure the /jvm directory to hold heap dumps, jfr files etc if the jvm root dir exists.
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

	details.AddArg(fmt.Sprintf("-XX:HeapDumpPath=%s/heap-dumps/%s-%s.hprof", jvmDir, member, podUID))

	if details.IsEnvTrue(v1.EnvVarJvmJmxmpEnabled) {
		details.AddClasspath(details.UtilsDir + "/lib/opendmk_jmxremote_optional_jar.jar")
		details.AddArg("-Dcoherence.management.serverfactory=com.oracle.coherence.k8s.JmxmpServer")

		details.SetSystemPropertyFromEnvVarOrDefault(v1.EnvVarJvmJmxmpPort, "-Dcoherence.jmxmp.port", fmt.Sprintf("%d", v1.DefaultJmxmpPort))
	}

	gc := strings.ToLower(details.Getenv(v1.EnvVarJvmGcCollector))
	switch {
	case gc == "" || gc == "g1":
		details.AddArg("-XX:+UseG1GC")
	case gc == "cms":
		details.AddArg("-XX:+UseConcMarkSweepGC")
	case gc == "parallel":
		details.AddArg("-XX:+UseParallelGC")
	}

	maxRAM := details.Getenv(v1.EnvVarJvmMaxRAM)
	if maxRAM != "" {
		details.AddArg("-XX:MaxRAM=" + maxRAM)
	}

	heap := details.Getenv(v1.EnvVarJvmMemoryHeap)
	if heap != "" {
		details.AddArg("-Xms" + heap)
		details.AddArg("-Xmx" + heap)
	}

	initial := details.Getenv(v1.EnvVarJvmInitialRAMPercentage)
	if initial != "" {
		q, err := resource.ParseQuantity(initial)
		if err == nil {
			d := q.AsDec()
			details.AddArg("-XX:InitialRAMPercentage=" + d.String())
		} else {
			fmt.Printf("ERROR: InitialRAMPercentage \"%s\" not a valid resource.Quantity - %s\n", initial, err.Error())
			os.Exit(1)
		}
	}

	max := details.Getenv(v1.EnvVarJvmMaxRAMPercentage)
	if max != "" {
		q, err := resource.ParseQuantity(max)
		if err == nil {
			d := q.AsDec()
			details.AddArg("-XX:MaxRAMPercentage=" + d.String())
		} else {
			fmt.Printf("ERROR: MaxRAMPercentage \"%s\" not a valid resource.Quantity - %s\n", max, err.Error())
			os.Exit(1)
		}
	}

	min := details.Getenv(v1.EnvVarJvmMinRAMPercentage)
	if min != "" {
		q, err := resource.ParseQuantity(min)
		if err == nil {
			d := q.AsDec()
			details.AddArg("-XX:MinRAMPercentage=" + d.String())
		} else {
			fmt.Printf("ERROR: MinRAMPercentage \"%s\" not a valid resource.Quantity - %s\n", min, err.Error())
			os.Exit(1)
		}
	}

	direct := details.Getenv(v1.EnvVarJvmMemoryDirect)
	if direct != "" {
		details.AddArg("-XX:MaxDirectMemorySize=" + direct)
	}
	stack := details.Getenv(v1.EnvVarJvmMemoryStack)
	if stack != "" {
		details.AddArg("-Xss" + stack)
	}
	meta := details.Getenv(v1.EnvVarJvmMemoryMeta)
	if meta != "" {
		details.AddArg("-XX:MetaspaceSize=" + meta)
		details.AddArg("-XX:MaxMetaspaceSize=" + meta)
	}
	track := details.GetenvOrDefault(v1.EnvVarJvmMemoryNativeTracking, "summary")
	if track != "" {
		details.AddArg("-XX:NativeMemoryTracking=" + track)
		details.AddArg("-XX:+PrintNMTStatistics")
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

	details.AddArg("-Dcoherence.ttl=0")

	details.AddArg(fmt.Sprintf("-XX:ErrorFile=%s/hs-err-%s-%s.log", jvmDir, member, podUID))

	if details.IsEnvTrueOrBlank(v1.EnvVarJvmOomHeapDump) {
		details.AddArg("-XX:+HeapDumpOnOutOfMemoryError")
	}

	if details.IsEnvTrueOrBlank(v1.EnvVarJvmOomExit) {
		details.AddArg("-XX:+ExitOnOutOfMemoryError")
	}

	// Use JVM container support
	if details.IsEnvTrueOrBlank(v1.EnvVarJvmUseContainerLimits) {
		details.AddArg("-XX:+UseContainerSupport")
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

	var cmd *exec.Cmd
	var app string
	switch {
	case details.AppType == AppTypeNone || details.AppType == AppTypeJava:
		app = "Java"
		cmd, err = createJavaCommand(details.GetJava(), details)
	case details.AppType == AppTypeHelidon:
		app = "Java"
		cmd, err = createJavaCommand(details.GetJava(), details)
	case details.AppType == AppTypeCoherence:
		app = "Java"
		cmd, err = createJavaCommand(details.GetJava(), details)
	case details.AppType == AppTypeTest:
		app = "Java"
		cmd, err = createJavaCommand(v1.TestCommand, details)
	default:
		app = "Graal (" + details.AppType + ")"
		cmd, err = runGraal(details)
	}

	return app, cmd, err
}

func createJavaCommand(javaCmd string, details *RunDetails) (*exec.Cmd, error) {
	args := details.GetCommand()
	args = append(args, details.MainClass)
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

func runGraal(details *RunDetails) (*exec.Cmd, error) {
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
func configureSiteAndRack(details *RunDetails) {
	fmt.Println("INFO: Configuring Coherence site and rack")
	if !details.GetSite {
		return
	}

	var site string
	var rack string
	siteLocation := details.Getenv(v1.EnvVarCohSite)
	fmt.Printf("INFO: Configuring Coherence site from '%s'\n", siteLocation)
	rackLocation := details.Getenv(v1.EnvVarCohRack)
	fmt.Printf("INFO: Configuring Coherence rack from '%s'\n", rackLocation)

	if siteLocation != "" {
		switch {
		case strings.ToLower(siteLocation) == "http://":
			site = ""
		case strings.HasPrefix(siteLocation, "http://"):
			// do http get
			site = httpGet(siteLocation, details)
		default:
			st, err := os.Stat(siteLocation)
			if err == nil && !st.IsDir() {
				bytes, err := ioutil.ReadFile(siteLocation)
				if err != nil {
					site = string(bytes)
				}
			}
		}
	}

	if rackLocation != "" {
		switch {
		case strings.ToLower(rackLocation) == "http://":
			rack = ""
		case strings.HasPrefix(rackLocation, "http://"):
			// do http get
			rack = httpGet(rackLocation, details)
		default:
			st, err := os.Stat(rackLocation)
			if err == nil && !st.IsDir() {
				bytes, err := ioutil.ReadFile(rackLocation)
				if err != nil {
					rack = string(bytes)
				}
			}
		}
	}

	if site != "" {
		details.AddArg("-Dcoherence.site=" + site)
	}

	if rack != "" {
		details.AddArg("-Dcoherence.rack=" + rack)
	} else if site != "" {
		details.AddArg("-Dcoherence.rack=" + site)
	}
}

// Do a http get for the specified url and return the response body for
// a 200 response or empty string for a non-200 response or error.
func httpGet(url string, details *RunDetails) string {
	fmt.Printf("INFO: Performing http get from '%s'\n", url)
	timeout := 120

	val := details.Getenv(v1.EnvVarOperatorTimeout)
	if val != "" {
		t, err := strconv.Atoi(val)
		if err == nil {
			timeout = t
		} else {
			fmt.Printf("ERROR: Invalid value set for %s '%s' using default of 120\n", v1.EnvVarOperatorTimeout, val)
		}
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("ERROR: failed to get url %s - %s\n", url, err.Error())
		return ""
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("ERROR: filed to get 200 response from %s - %s\n", url, resp.Status)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: filed to read response body from %s - %s\n", url, resp.Status)
		return ""
	}

	s := string(body)
	fmt.Printf("INFO: Get response from '%s' was '%s'\n", url, s)
	return s
}

func checkCoherenceVersion(v string, details *RunDetails) (bool, error) {
	fmt.Printf("INFO: Checking for Coherence version %s\n", v)

	if details.IsEnvTrue(v1.EnvVarCohSkipVersionCheck) {
		fmt.Printf("INFO: Skipping Coherence version check %s=%s\n", v1.EnvVarCohSkipVersionCheck, details.Getenv(v1.EnvVarCohSkipVersionCheck))
		return true, nil
	}

	// Get the classpath to use (we need Coherence jar)
	cp := fmt.Sprintf("%s/lib/*:%s", details.UtilsDir, details.GetClasspath())

	// Configure the os command
	java := details.GetJava()
	cmd := exec.Command(java, "-cp", cp, "com.oracle.coherence.k8s.CoherenceVersion", v)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("INFO: Command: %s %s\n", cmd.Path, strings.Join(cmd.Args, " "))
	// execute the command
	err := cmd.Run()
	if err == nil {
		// command exited with exit code 0
		fmt.Printf("INFO: Coherence version is at least %s\n", v)
		return true, nil
	}
	if _, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0
		fmt.Printf("INFO: Coherence version is lower than %s\n", v)
		return false, nil
	}
	// command exited with some other error
	fmt.Printf("ERROR: Coherence version check failed %s\n", err.Error())
	return false, err
}

func cohPre12214(details *RunDetails) {
	details.AddArg("-Dcoherence.override=k8s-coherence-nossl-override.xml")
	details.AddArgFromEnvVar(v1.EnvVarCohOverride, "-Dcoherence.k8s.override")
}

func cohPost12214(details *RunDetails) {
	details.AddArg("-Dcoherence.override=k8s-coherence-override.xml")
	details.AddArgFromEnvVar(v1.EnvVarCohOverride, "-Dcoherence.k8s.override")
}

func addManagementSSL(details *RunDetails) {
	addSSL(v1.EnvVarCohMgmtPrefix, v1.PortNameManagement, details)
}

func addMetricsSSL(details *RunDetails) {
	addSSL(v1.EnvVarCohMetricsPrefix, v1.PortNameMetrics, details)
}

func addSSL(prefix, prop string, details *RunDetails) {
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

func usage() {
	message := `Runner Usage:
  server         Start a Coherence server
  console        Start a Coherence console
  queryplus      Start a Coherence Query Plus console
  mbeanserver    Start a Coherence MBean server`
	fmt.Println(message)
}

func printHeader(details *RunDetails) {
	fmt.Println("INFO: Coherence Operator Utils Runner")
	fmt.Printf("INFO:   Version: %s\n", details.Getenv(EnvVarVersion))
	fmt.Printf("INFO:   Commit:  %s\n", details.Getenv(EnvVarGitCommit))
	fmt.Printf("INFO:   Built:   %s\n", details.Getenv(EnvVarBuildDate))

	fmt.Println("INFO: Args:")
	for _, a := range details.OsArgs {
		fmt.Println("INFO:     " + a)
	}

	keys := make([]string, 0, len(details.Env))
	for k := range details.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("INFO: Env:")
	for _, k := range keys {
		fmt.Printf("INFO:     %s=%s\n", k, details.Env[k])
	}

	fmt.Println("INFO:")
}

// ---- RunDetails struct --------------------------------------------------------------------------

type RunDetails struct {
	Command       string
	Env           map[string]string
	OsArgs        []string
	CoherenceHome string
	JavaHome      string
	UtilsDir      string
	Dir           string
	GetSite       bool
	AppType       string
	Classpath     string
	Args          []string
	MainClass     string
	MainArgs      []string
}

func (in *RunDetails) Getenv(name string) string {
	return in.Env[name]
}

func (in *RunDetails) GetenvOrDefault(name string, defaultValue string) string {
	v, ok := in.Env[name]
	if ok && v != "" {
		return v
	}
	return defaultValue
}

func (in *RunDetails) LookupEnv(name string) (string, bool) {
	v, ok := in.Env[name]
	return v, ok
}

func (in *RunDetails) GetenvWithPrefix(prefix, name string) string {
	return in.Getenv(prefix + name)
}

func (in *RunDetails) Setenv(key, value string) {
	in.Env[key] = value
}

func (in *RunDetails) Unsetenv(key string) {
	delete(in.Env, key)
}

func (in *RunDetails) IsEnvTrue(name string) bool {
	value := in.Getenv(name)
	return strings.ToLower(value) == "true"
}

func (in *RunDetails) IsEnvTrueOrBlank(name string) bool {
	value := in.Getenv(name)
	return value == "" || strings.ToLower(value) == "true"
}

func (in *RunDetails) GetCommand() []string {
	var cmd []string
	cp := in.GetClasspath()
	if cp != "" {
		cmd = append(cmd, "-cp", cp)
	}
	cmd = append(cmd, in.Args...)
	return cmd
}

func (in *RunDetails) GetGraalCommand() []string {
	cmd := in.GetCommand()
	for i, c := range cmd {
		switch {
		case c == "-cp":
			cmd[i] = "--vm.cp"
		case strings.HasPrefix(c, "-D"):
			cmd[i] = "--vm." + c[1:]
		case strings.HasPrefix(c, "-XX"):
			cmd[i] = "--vm." + c[1:]
		case strings.HasPrefix(c, "-Xms"):
			cmd[i] = "--vm." + c[1:]
		case strings.HasPrefix(c, "-Xmx"):
			cmd[i] = "--vm." + c[1:]
		case strings.HasPrefix(c, "-Xss"):
			cmd[i] = "--vm." + c[1:]
		}
	}

	return cmd
}

func (in *RunDetails) AddArgs(args ...string) {
	for _, a := range args {
		in.AddArg(a)
	}
}

func (in *RunDetails) AddArg(arg string) {
	if arg != "" {
		in.Args = append(in.Args, arg)
	}
}

func (in *RunDetails) AddToFrontOfClasspath(path string) {
	if path != "" {
		if in.Classpath == "" {
			in.Classpath = path
		} else {
			in.Classpath = path + ":" + in.Classpath
		}
	}
}

func (in *RunDetails) AddClasspath(path string) {
	if path != "" {
		if in.Classpath == "" {
			in.Classpath = path
		} else {
			in.Classpath += ":" + path
		}
	}
}

func (in *RunDetails) GetClasspath() string {
	cp := in.Classpath
	// if ${COHERENCE_HOME} exists add coherence.jar to the classpath
	if in.CoherenceHome != "" {
		cp = cp + ":" + in.CoherenceHome + "/conf"
		cp = cp + ":" + in.CoherenceHome + "/lib/coherence.jar"
	}
	return cp
}

func (in *RunDetails) AddArgFromEnvVar(name, property string) {
	value := in.Getenv(name)
	if value != "" {
		s := fmt.Sprintf("%s=%s", property, value)
		in.Args = append(in.Args, s)
	}
}

func (in *RunDetails) SetSystemPropertyFromEnvVarOrDefault(name, property, dflt string) {
	value := in.Getenv(name)
	var s string
	if value != "" {
		s = fmt.Sprintf("%s=%s", property, value)
	} else {
		s = fmt.Sprintf("%s=%s", property, dflt)
	}
	in.Args = append(in.Args, s)
}

func (in *RunDetails) GetJava() string {
	if in.JavaHome != "" {
		return in.JavaHome + "/bin/java"
	}
	return "java"
}
