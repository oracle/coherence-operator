/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/tebeka/go2xunit/lib"
)

const (
	// Version is the current version
	Version = "1.4.10"
)

// getInput return input io.File from file name, if file name is - it will
// return os.Stdin
func getInput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

// getInput return output io.File from file name, if file name is - it will
// return os.Stdout
func getOutput(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdout, nil
	}

	return os.Create(filename)
}

// getIO returns input and output streams from file names
func getIO(inFile, outFile string) (*os.File, io.Writer, error) {
	input, err := getInput(inFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for reading: %s", inFile, err)
	}

	output, err := getOutput(outFile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %s for writing: %s", outFile, err)
	}

	return input, output, nil
}

func main() {
	if args.showVersion {
		fmt.Printf("go2xunit %s\n", Version)
		os.Exit(0)
	}

	// No time ... prefix for error messages
	log.SetFlags(0)

	if err := validateArgs(); err != nil {
		log.Fatalf("error: %s", err)
	}

	input, output, err := getIO(args.inFile, args.outFile)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// We'd like the test time to be the time of the generated file
	var testTime time.Time
	stat, err := input.Stat()
	if err != nil {
		testTime = time.Now()
	} else {
		testTime = stat.ModTime()
	}

	var parse func(rd io.Reader, suiteName string) (lib.Suites, error)

	if args.isGocheck {
		parse = lib.ParseGocheck
	} else {
		parse = lib.ParseGotest
	}

	suites, err := parse(input, args.suitePrefix)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	if len(suites) == 0 {
		log.Fatalf("error: no tests found")
		os.Exit(1)
	}

	xmlTemplate := lib.XUnitTemplate
	if args.xunitnetOut {
		xmlTemplate = lib.XUnitNetTemplate
	} else if args.bambooOut || (len(suites) > 1) {
		xmlTemplate = lib.XMLMultiTemplate
	}

	lib.WriteXML(suites, output, xmlTemplate, testTime)
	if args.fail && suites.HasFailures() {
		os.Exit(1)
	}
}

var args struct {
	inFile      string
	outFile     string
	fail        bool
	showVersion bool
	bambooOut   bool
	xunitnetOut bool
	isGocheck   bool
	suitePrefix string
}

func init() {
	flag.StringVar(&args.inFile, "input", "", "input file (default to stdin)")
	flag.StringVar(&args.outFile, "output", "", "output file (default to stdout)")
	flag.BoolVar(&args.fail, "fail", false, "fail (non zero exit) if any test failed")
	flag.BoolVar(&args.showVersion, "version", false, "print version and exit")
	flag.BoolVar(&args.bambooOut, "bamboo", false,
		"xml compatible with Atlassian's Bamboo")
	flag.BoolVar(&args.xunitnetOut, "xunitnet", false, "xml compatible with xunit.net")
	flag.BoolVar(&args.isGocheck, "gocheck", false, "parse gocheck output")
	flag.BoolVar(&lib.Options.FailOnRace, "fail-on-race", false,
		"mark test as failing if it exposes a data race")
	flag.StringVar(&args.suitePrefix, "suite-name-prefix", "",
		"prefix to include before all suite names")

	flag.Parse()
}

// validateArgs validates command line arguments
func validateArgs() error {
	if flag.NArg() > 0 {
		return fmt.Errorf("%s does not take parameters (did you mean -input?)", os.Args[0])
	}

	if args.bambooOut && args.xunitnetOut {
		return fmt.Errorf("-bamboo and -xunitnet are mutually exclusive")
	}

	return nil
}
