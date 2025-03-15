/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package run_details

import (
	"fmt"
	"github.com/go-logr/logr"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func NewRunDetails(v *viper.Viper, log logr.Logger) *RunDetails {
	var err error

	skipSiteVar := v.GetString(v1.EnvVarCohSkipSite)
	skipSite := strings.ToLower(skipSiteVar) != "true"

	details := &RunDetails{
		env:           v,
		CoherenceHome: v.GetString(v1.EnvVarCoherenceHome),
		UtilsDir:      v.GetString(v1.EnvVarCohUtilDir),
		JavaHome:      v.GetString(v1.EnvVarJavaHome),
		AppType:       strings.ToLower(v.GetString(v1.EnvVarAppType)),
		Dir:           v.GetString(v1.EnvVarCohAppDir),
		MainClass:     v1.DefaultMain,
		GetSite:       skipSite,
		log:           log,
	}

	// add any Classpath items
	details.AddClasspath(v.GetString(v1.EnvVarJvmExtraClasspath))
	details.AddClasspath(v.GetString(v1.EnvVarJavaClasspath))

	cpFile := fmt.Sprintf(v1.FileNamePattern, details.UtilsDir, os.PathSeparator, v1.OperatorClasspathFile)
	if _, err = os.Stat(cpFile); err == nil {
		details.ClassPathFile = cpFile
	}

	argFile := fmt.Sprintf(v1.FileNamePattern, details.UtilsDir, os.PathSeparator, v1.OperatorJvmArgsFile)
	if _, err = os.Stat(argFile); err == nil {
		details.JvmArgsFile = argFile
	}

	return details
}

// RunDetails contains the information to run an application.
type RunDetails struct {
	Command           string
	CoherenceHome     string
	JavaHome          string
	UtilsDir          string
	Dir               string
	GetSite           bool
	UseOperatorHealth bool
	AppType           string
	Classpath         string
	MainClass         string
	InnerMainClass    string
	MainArgs          []string
	BuildPacks        *bool
	ExtraEnv          []string
	ClassPathFile     string
	JvmArgsFile       string
	args              []string
	vmOptions         []string
	memoryArgs        []string
	diagnosticArgs    []string
	env               *viper.Viper
	log               logr.Logger
}

func (in *RunDetails) GetAllArgs() []string {
	var args []string
	args = append(args, in.vmOptions...)
	args = append(args, in.memoryArgs...)
	args = append(args, in.diagnosticArgs...)
	args = append(args, in.args...)
	return args
}

func (in *RunDetails) GetArguments() []string {
	return in.args
}

func (in *RunDetails) GetMemoryOptions() []string {
	return in.memoryArgs
}

func (in *RunDetails) GetDiagnosticOptions() []string {
	return in.diagnosticArgs
}

func (in *RunDetails) GetVMOptions() []string {
	return in.vmOptions
}

// IsSpringBoot returns true if this is a Spring Boot application
func (in *RunDetails) IsSpringBoot() bool {
	if in.env == nil {
		return false
	}
	return in.AppType == v1.AppTypeSpring2 || in.AppType == v1.AppTypeSpring3
}

// Getenv returns the value for the specified environment variable, or empty string if not set.
func (in *RunDetails) Getenv(name string) string {
	if in.env == nil {
		return ""
	}
	return in.env.GetString(name)
}

// ExpandEnv replaces ${var} or $(var) or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func (in *RunDetails) ExpandEnv(s string) string {
	return in.Expand(s, in.Getenv)
}

// Expand replaces ${var} or $(var) or $var in the string based on the mapping function.
// For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).
func (in *RunDetails) Expand(s string, mapping func(string) string) string {
	var buf []byte
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, w := in.GetShellName(s[j+1:])
			switch {
			case name == "" && w > 0:
				// Encountered invalid syntax; eat the
				// characters.
				break
			case name == "":
				// Valid syntax, but $ was not followed by a
				// name. Leave the dollar character untouched.
				buf = append(buf, s[j])
			default:
				buf = append(buf, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

// GetShellName returns the name that begins the string and the number of bytes
// consumed to extract it. If the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name.
func (in *RunDetails) GetShellName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && in.isShellSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "${}"
				}
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; eat "${"
	case s[0] == '(':
		if len(s) > 2 && in.isShellSpecialVar(s[1]) && s[2] == ')' {
			return s[1:2], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == ')' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "$()"
				}
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; eat "$("
	case in.isShellSpecialVar(s[0]):
		return s[0:1], 1
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(s) && in.isAlphaNum(s[i]); i++ {
		// empty ??
	}
	return s[:i], i
}

// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func (in *RunDetails) isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func (in *RunDetails) isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func (in *RunDetails) GetenvOrDefault(name string, defaultValue string) string {
	if in.env != nil && in.env.IsSet(name) {
		return in.env.GetString(name)
	}
	return defaultValue
}

func (in *RunDetails) LookupEnv(name string) (string, bool) {
	if in.env != nil && in.env.IsSet(name) {
		return in.env.GetString(name), true
	}
	return "", false
}

func (in *RunDetails) GetenvWithPrefix(prefix, name string) string {
	return in.Getenv(prefix + name)
}

func (in *RunDetails) Setenv(key, value string) {
	if in.env == nil {
		in.env = viper.New()
	}
	in.env.Set(key, value)
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
	return in.GetCommandWithPrefix("", "")
}

func (in *RunDetails) GetCommandWithPrefix(propPrefix, jvmPrefix string) []string {
	var cmd []string
	cp := in.GetClasspath()
	if cp != "" {
		cmd = append(cmd, v1.JvmOptClassPath, cp)
	}
	if propPrefix == "" && jvmPrefix == "" {
		cmd = append(cmd, in.GetAllArgs()...)
	} else {
		for _, arg := range in.GetAllArgs() {
			if strings.HasPrefix(arg, "-D") && propPrefix != "" {
				cmd = append(cmd, propPrefix+arg)
			} else if jvmPrefix != "" {
				cmd = append(cmd, jvmPrefix+arg)
			}
		}
	}
	return cmd
}

func (in *RunDetails) GetSpringBootCommand() []string {
	return append(in.GetSpringBootArgs(), in.MainClass)
}

func (in *RunDetails) GetSpringBootArgs() []string {
	var cmd []string
	// Are we using a Spring Boot fat jar
	if jar, _ := in.LookupEnv(v1.EnvVarSpringBootFatJar); jar != "" {
		cmd = append(cmd, v1.JvmOptClassPath, jar)
	}
	return append(cmd, in.GetAllArgs()...)
}

func (in *RunDetails) AddArgs(args ...string) {
	for _, a := range args {
		in.AddArg(a)
	}
}

func (in *RunDetails) AddArg(arg string) {
	if arg != "" {
		in.args = append(in.args, in.ExpandEnv(arg))
	}
}

func (in *RunDetails) AddSystemPropertyArg(propName, value string) {
	if propName != "" {
		arg := fmt.Sprintf(v1.SystemPropertyPattern, propName, value)
		in.args = append(in.args, in.ExpandEnv(arg))
	}
}

func (in *RunDetails) AddVMOption(arg string) {
	if arg != "" {
		in.vmOptions = append(in.vmOptions, in.ExpandEnv(arg))
	}
}

func (in *RunDetails) AddMemoryOption(arg string) {
	if arg != "" {
		in.memoryArgs = append(in.memoryArgs, in.ExpandEnv(arg))
	}
}

func (in *RunDetails) AddDiagnosticOption(arg string) {
	if arg != "" {
		in.diagnosticArgs = append(in.diagnosticArgs, in.ExpandEnv(arg))
	}
}

func (in *RunDetails) AddToFrontOfClasspath(path string) {
	if path != "" {
		if in.Classpath == "" {
			in.Classpath = in.ExpandEnv(path)
		} else {
			in.Classpath = in.ExpandEnv(path) + ":" + in.Classpath
		}
	}
}

// addJarsToClasspath adds all jars in the specified directory to the classpath
func (in *RunDetails) AddJarsToClasspath(dir string) {
	path := in.ExpandEnv(dir)
	if _, err := os.Stat(path); err == nil {
		var jars []string
		_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			name := info.Name()
			if !info.IsDir() && (strings.HasSuffix(name, ".jar") || strings.HasSuffix(name, ".JAR")) {
				jars = append(jars, path)
			}
			return nil
		})

		for _, jar := range jars {
			in.AddClasspath(jar)
		}
	}
}

func (in *RunDetails) AddClasspathIfExists(path string) {
	if _, err := os.Stat(path); err == nil {
		in.AddClasspath(path)
	}
}

func (in *RunDetails) AddClasspath(path string) {
	if path != "" {
		if in.Classpath == "" {
			in.Classpath = in.ExpandEnv(path)
		} else {
			in.Classpath += ":" + in.ExpandEnv(path)
		}
	}
}

func (in *RunDetails) GetClasspath() string {
	cp := in.Classpath
	// if ${COHERENCE_HOME} exists add coherence.jar to the classpath
	if in.CoherenceHome != "" {
		if _, err := os.Stat(in.CoherenceHome); err == nil {
			if _, err := os.Stat(in.CoherenceHome + "/conf"); err == nil {
				cp = cp + ":" + in.CoherenceHome + "/conf"
			}
			if _, err := os.Stat(in.CoherenceHome + "/lib/coherence.jar"); err == nil {
				cp = cp + ":" + in.CoherenceHome + "/lib/coherence.jar"
			}
		}
	}
	return cp
}

func (in *RunDetails) AddArgFromEnvVar(name, property string) {
	value := in.Getenv(name)
	if value != "" {
		s := fmt.Sprintf("%s=%s", property, value)
		in.AddArg(s)
	}
}

func (in *RunDetails) AddSystemPropertyFromEnvVar(name, property string) {
	value := in.Getenv(name)
	if value != "" {
		in.AddSystemPropertyArg(property, value)
	}
}

func (in *RunDetails) SetSystemPropertyFromEnvVarOrDefault(name, property, dflt string) {
	value := in.Getenv(name)
	if value != "" {
		in.AddSystemPropertyArg(property, value)
	} else {
		in.AddSystemPropertyArg(property, dflt)
	}
}

func (in *RunDetails) GetJavaExecutable() string {
	if in.JavaHome != "" {
		return in.JavaHome + "/bin/java"
	}
	return "java"
}

func (in *RunDetails) GetJShellExecutable() string {
	if in.JavaHome != "" {
		return in.JavaHome + "/bin/jshell"
	}
	return "jshell"
}

// IsBuildPacks determines whether to run the application with the Cloud Native Buildpack launcher
func (in *RunDetails) IsBuildPacks() bool {
	if in.BuildPacks == nil {
		var bp bool
		detect := strings.ToLower(in.env.GetString(v1.EnvVarCnbpEnabled))
		switch detect {
		case "true":
			in.log.Info("Detecting Cloud Native Buildpacks", "envVar", v1.EnvVarCnbpEnabled, "value", detect)
			bp = true
		case "false":
			in.log.Info("Detecting Cloud Native Buildpacks", "envVar", v1.EnvVarCnbpEnabled, "value", detect)
			bp = false
		default:
			in.log.Info("Auto-detecting Cloud Native Buildpacks")
			// else auto detect
			// look for the CNB API environment variable
			_, ok := os.LookupEnv("CNB_PLATFORM_API")
			in.log.Info(fmt.Sprintf("Auto-detecting Cloud Native Buildpacks: CNB_PLATFORM_API found=%t", ok))
			// look for the CNB launcher
			launcher := GetBuildpackLauncher()
			_, err := os.Stat(launcher)
			in.log.Info(fmt.Sprintf("Auto-detecting Cloud Native Buildpacks: CNB Launcher '%s' found=%t\n", launcher, err == nil))
			bp = ok && err == nil
		}
		in.BuildPacks = &bp
	}
	return *in.BuildPacks
}

func GetBuildpackLauncher() string {
	if launcher, ok := os.LookupEnv(v1.EnvVarCnbpLauncher); ok {
		return launcher
	}
	return v1.DefaultCnbpLauncher
}
