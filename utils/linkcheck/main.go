/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sort"
	"strings"
	"time"
)

// checks all the http links in a document

var (
	// log is the logger used by the runner
	log = ctrl.Log.WithName("linkcheck")
)

const (
	ArgFile    = "file"
	ArgExclude = "exclude"
)

func main() {

	if _, err := Execute(); err != nil {
		logf.Log.WithName("runner").Error(err, "Unexpected error while executing command")
		os.Exit(1)
	}
}

// Execution is a holder of details of a command execution
type Execution struct {
	Cmd *cobra.Command
	App string
	V   *viper.Viper
}

// contextKey allows type safe Context Values.
type contextKey int

// The key to obtain an execution from a Context.
var executionKey contextKey

// Execute runs the runner with a given environment.
func Execute() (Execution, error) {
	return ExecuteWithArgs(nil)
}

// ExecuteWithArgs runs the runner with a given environment and argument overrides.
func ExecuteWithArgs(args []string) (Execution, error) {
	cmd, v := NewRootCommand()
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

// NewRootCommand builds the root cobra command that handles our command line tool.
func NewRootCommand() (*cobra.Command, *viper.Viper) {
	v := viper.New()

	// rootCommand is the Cobra root Command to execute
	rootCmd := &cobra.Command{
		Use:   "runner",
		Short: "Run the link checker",
		Long:  "Run the link checker",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	flags := rootCmd.Flags()
	flags.StringArray(ArgFile, nil, "a file or directory to scan")
	flags.StringArray(ArgExclude, nil, "a link prefix to exclude")
	return rootCmd, v
}

func run(cmd *cobra.Command) error {
	flagSet := cmd.Flags()
	files, err := flagSet.GetStringArray(ArgFile)
	if err != nil {
		return err
	}
	excludes, err := flagSet.GetStringArray(ArgExclude)
	if err != nil {
		return err
	}

	excludes = append(excludes, "https://localhost")
	excludes = append(excludes, "http://localhost")
	excludes = append(excludes, "http://127.0.0.1")
	excludes = append(excludes, "https://127.0.0.1")
	excludes = append(excludes, "https://host")
	excludes = append(excludes, "http://host")

	exitCode := checkDocs(files, excludes)
	if exitCode != 0 {
		return fmt.Errorf("link checking failed")
	}
	return nil
}

func checkDocs(paths []string, excludes []string) int {
	exitCode := 0
	for _, path := range paths {
		if strings.HasSuffix(path, "/...") {
			exitCode = checkDirectory(path[0:len(path)-4], excludes)
		} else {
			exitCode = checkDoc(path, excludes)
		}
	}
	return exitCode
}

func checkDirectory(dirName string, excludes []string) int {
	fmt.Printf("Checking directory %s\n", dirName)
	info, err := os.Stat(dirName)
	if err != nil {
		fmt.Printf(err.Error())
		return 1
	}
	if !info.IsDir() {
		fmt.Printf("%s is not a directory", dirName)
		return 1
	}
	return checkFileInfo(dirName, info, excludes)
}

func checkFileInfo(dir string, info os.FileInfo, excludes []string) int {
	if info.IsDir() {
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf(err.Error())
			return 1
		}
		exitCode := 0
		for _, f := range files {
			name := f.Name()
			if !strings.HasPrefix(name, ".") {
				fullName := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), name)
				if f.IsDir() {
					if checkDirectory(fullName, excludes) != 0 {
						exitCode = 1
					}
				} else {
					if checkDoc(fullName, excludes) != 0 {
						exitCode = 1
					}
				}
			}
		}
		return exitCode
	} else {
		return checkDoc(info.Name(), excludes)
	}
}

func checkDoc(path string, excludes []string) int {
	if !strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".html") {
		return 0
	}

	fmt.Printf("Checking file %s\n", path)
	exitCode := 0

	buf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	s := string(buf)
	links, mapFragment, err := parseLinks(s, excludes)
	if err != nil {
		panic(err)
	}
	for _, link := range links {
		fragments := mapFragment[link]
		if checkLink(link, fragments, excludes) != 0 {
			exitCode = 1
		}
	}
	return exitCode
}

func checkLink(link string, fragments []string, excludes []string) int {
	var (
		err      error
		urlToGet *url.URL
	)

	for _, skip := range excludes {
		if strings.HasPrefix(link, skip) {
			return 0
		}
	}

	exitCode := 0

	if strings.Contains(link, "(") {
		fmt.Printf("%s FAILED - not a valid link\n", link)
		return 1
	}

	// Parse URL
	if urlToGet, err = url.Parse(link); err != nil {
		fmt.Println(err)
		return 1
	}

	// Retrieve content of URL
	if check := checkURL(urlToGet, fragments); check != 0 {
		exitCode = 1
	}

	return exitCode
}

func checkURL(urlToGet *url.URL, fragments []string) int {
	var (
		err     error
		resp    *http.Response
		content []byte
	)

	fmt.Printf("%s", urlToGet)

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}

	if resp, err = netClient.Get(urlToGet.String()); err != nil {
		fmt.Printf(" FAILED error: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		fmt.Printf(" FAILED response: %d\n", resp.StatusCode)
		return 1
	}

	count := len(fragments)
	if (count == 1 && fragments[0] != "") || count > 1 {
		// Read the body of the HTTP response
		if content, err = io.ReadAll(resp.Body); err != nil {
			fmt.Printf(" FAILED error: %v\n", err)
			return 1
		}

		fmt.Println(" OK")

		pageContent := string(content)

		for _, fragment := range fragments {
			if fragment != "" {
				fmt.Printf("%s#%s", urlToGet, fragment)
				if !checkFragment(fragment, pageContent) {
					return 1
				}
				fmt.Println(" OK")
			}
		}
	} else {
		fmt.Println(" OK")
	}

	return 0
}

func checkFragment(fragment, pageContent string) bool {
	var headings []string

	headings = append(headings, fmt.Sprintf("id=\"%s\"", fragment))
	headings = append(headings, fmt.Sprintf("id=%s", fragment))
	headings = append(headings, fmt.Sprintf("href=\"#%s\"", fragment))
	headings = append(headings, fmt.Sprintf("href=\\\"#%s\\\"", fragment))

	for _, heading := range headings {
		if strings.Contains(pageContent, heading) {
			return true
		}
	}

	fmt.Println(" FAILED could not find any of the following headings:")
	for _, heading := range headings {
		fmt.Printf("   %s", heading)
	}
	return false
}

func parseLinks(content string, excludes []string) ([]string, map[string][]string, error) {
	var (
		err       error
		matches   [][]string
		links     []string
		findLinks = regexp.MustCompile("<a.*?href=\"(.*?)\"")
	)

	// Retrieve all anchor tag URLs from string
	matches = findLinks.FindAllStringSubmatch(content, -1)

	linkMap := make(map[string][]string)

	for _, val := range matches {
		var linkUrl *url.URL

		// Parse the anchor tag URL
		u := val[1]
		skip := false
		for _, exclude := range excludes {
			if strings.HasPrefix(u, exclude) {
				skip = true
				break
			}
		}

		if skip {
			break
		}

		if linkUrl, err = url.Parse(u); err != nil {
			return links, linkMap, err
		}

		if linkUrl.IsAbs() {
			s := linkUrl.String()
			if linkUrl.Fragment != "" {
				s = s[0:(len(s) - 1 - len(linkUrl.Fragment))]
			}

			l := append(linkMap[s], linkUrl.Fragment)
			linkMap[s] = l
		}

		links = make([]string, 0)
		for link, _ := range linkMap {
			links = append(links, link)
		}
	}

	sort.Strings(links)
	return links, linkMap, err
}
