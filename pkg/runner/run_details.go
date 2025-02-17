/*
 * Copyright (c) 2021, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func NewRunDetails(v *viper.Viper) *RunDetails {
	skipSiteVar := v.GetString(v1.EnvVarCohSkipSite)
	skipSite := strings.ToLower(skipSiteVar) != "true"

	details := &RunDetails{
		env:           v,
		CoherenceHome: v.GetString(v1.EnvVarCoherenceHome),
		UtilsDir:      v.GetString(v1.EnvVarCohUtilDir),
		JavaHome:      v.GetString(v1.EnvVarJavaHome),
		AppType:       strings.ToLower(v.GetString(v1.EnvVarAppType)),
		Dir:           v.GetString(v1.EnvVarCohAppDir),
		MainClass:     DefaultMain,
		GetSite:       skipSite,
	}

	// add any Classpath items
	details.addClasspath(v.GetString(v1.EnvVarJvmExtraClasspath))
	details.addClasspath(v.GetString(v1.EnvVarJavaClasspath))

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
	Args              []string
	MainClass         string
	MainArgs          []string
	BuildPacks        *bool
	ExtraEnv          []string
	env               *viper.Viper
}

// IsSpringBoot returns true if this is a Spring Boot application
func (in *RunDetails) IsSpringBoot() bool {
	if in.env == nil {
		return false
	}
	return in.AppType == AppTypeSpring2 || in.AppType == AppTypeSpring3
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
			name, w := in.getShellName(s[j+1:])
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

// getShellName returns the name that begins the string and the number of bytes
// consumed to extract it. If the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name.
func (in *RunDetails) getShellName(s string) (string, int) {
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

func (in *RunDetails) getenvOrDefault(name string, defaultValue string) string {
	if in.env != nil && in.env.IsSet(name) {
		return in.env.GetString(name)
	}
	return defaultValue
}

func (in *RunDetails) lookupEnv(name string) (string, bool) {
	if in.env != nil && in.env.IsSet(name) {
		return in.env.GetString(name), true
	}
	return "", false
}

func (in *RunDetails) getenvWithPrefix(prefix, name string) string {
	return in.Getenv(prefix + name)
}

func (in *RunDetails) setenv(key, value string) {
	if in.env == nil {
		in.env = viper.New()
	}
	in.env.Set(key, value)
}

func (in *RunDetails) unsetenv(key string) {
	if in.env == nil {
		in.env = viper.New()
	}
	in.env.Set(key, nil)
}

func (in *RunDetails) isEnvTrue(name string) bool {
	value := in.Getenv(name)
	return strings.ToLower(value) == "true"
}

func (in *RunDetails) isEnvTrueOrBlank(name string) bool {
	value := in.Getenv(name)
	return value == "" || strings.ToLower(value) == "true"
}

func (in *RunDetails) getCommand() []string {
	var cmd []string
	cp := in.getClasspath()
	if cp != "" {
		cmd = append(cmd, "-cp", cp)
	}
	cmd = append(cmd, in.Args...)
	return cmd
}

func (in *RunDetails) getSpringBootCommand() []string {
	return append(in.getSpringBootArgs(), in.MainClass)
}

func (in *RunDetails) getSpringBootArgs() []string {
	var cmd []string
	cp := strings.ReplaceAll(in.getClasspath(), ":", ",")
	if cp != "" {
		cmd = append(cmd, "-Dloader.path="+cp)
	}

	// Are we using a Spring Boot fat jar
	if jar, _ := in.lookupEnv(v1.EnvVarSpringBootFatJar); jar != "" {
		cmd = append(cmd, "-cp", jar)
	}
	cmd = append(cmd, in.Args...)

	return cmd
}

/*
func (in *RunDetails) getGraalCommand() []string {
	cmd := in.getCommand()
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
*/

func (in *RunDetails) addArgs(args ...string) {
	for _, a := range args {
		in.addArg(a)
	}
}

func (in *RunDetails) addArg(arg string) {
	if arg != "" {
		in.Args = append(in.Args, in.ExpandEnv(arg))
	}
}

func (in *RunDetails) addToFrontOfClasspath(path string) {
	if path != "" {
		if in.Classpath == "" {
			in.Classpath = in.ExpandEnv(path)
		} else {
			in.Classpath = in.ExpandEnv(path) + ":" + in.Classpath
		}
	}
}

// addJarsToClasspath adds all jars in the specified directory to the classpath
func (in *RunDetails) addJarsToClasspath(dir string) {
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

		sort.Strings(jars)
		for _, jar := range jars {
			in.addClasspath(jar)
		}
	}
}

func (in *RunDetails) addClasspathIfExists(path string) {
	if _, err := os.Stat(path); err == nil {
		in.addClasspath(path)
	}
}

func (in *RunDetails) addClasspath(path string) {
	if path != "" {
		if in.Classpath == "" {
			in.Classpath = in.ExpandEnv(path)
		} else {
			in.Classpath += ":" + in.ExpandEnv(path)
		}
	}
}

func (in *RunDetails) getClasspath() string {
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

func (in *RunDetails) addArgFromEnvVar(name, property string) {
	value := in.Getenv(name)
	if value != "" {
		s := fmt.Sprintf("%s=%s", property, value)
		in.Args = append(in.Args, s)
	}
}

func (in *RunDetails) setSystemPropertyFromEnvVarOrDefault(name, property, dflt string) {
	value := in.Getenv(name)
	var s string
	if value != "" {
		s = fmt.Sprintf("%s=%s", property, value)
	} else {
		s = fmt.Sprintf("%s=%s", property, dflt)
	}
	in.Args = append(in.Args, s)
}

func (in *RunDetails) getJavaExecutable() string {
	if in.JavaHome != "" {
		return in.JavaHome + "/bin/java"
	}
	return "java"
}

// isBuildPacks determines whether to run the application with the Cloud Native Buildpack launcher
func (in *RunDetails) isBuildPacks() bool {
	if in.BuildPacks == nil {
		var bp bool
		detect := strings.ToLower(in.env.GetString(v1.EnvVarCnbpEnabled))
		switch detect {
		case "true":
			log.Info("Detecting Cloud Native Buildpacks", "envVar", v1.EnvVarCnbpEnabled, "value", detect)
			bp = true
		case "false":
			log.Info("Detecting Cloud Native Buildpacks", "envVar", v1.EnvVarCnbpEnabled, "value", detect)
			bp = false
		default:
			log.Info("Auto-detecting Cloud Native Buildpacks")
			// else auto detect
			// look for the CNB API environment variable
			_, ok := os.LookupEnv("CNB_PLATFORM_API")
			log.Info(fmt.Sprintf("Auto-detecting Cloud Native Buildpacks: CNB_PLATFORM_API found=%t", ok))
			// look for the CNB launcher
			launcher := getBuildpackLauncher()
			_, err := os.Stat(launcher)
			log.Info(fmt.Sprintf("Auto-detecting Cloud Native Buildpacks: CNB Launcher '%s' found=%t\n", launcher, err == nil))
			bp = ok && err == nil
		}
		in.BuildPacks = &bp
	}
	return *in.BuildPacks
}
