/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	// CommandSleep is the argument to sleep for a number of seconds.
	CommandSleep = "sleep"
	// ArgSeconds is the number of seconds to sleep
	ArgSeconds = "seconds"
)

// queryPlusCommand creates the corba "sleep" sub-command
func sleepCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandSleep,
		Short: "Sleep for a number of seconds",
		Long:  "Sleep for a number of seconds",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sleep(cmd, args)
		},
	}

	flagSet := cmd.Flags()
	flagSet.Int64P(ArgSeconds, "t", 600, "The number of seconds to sleep")

	return cmd
}

func sleep(cmd *cobra.Command, _ []string) error {
	flagSet := cmd.Flags()

	var seconds int64
	seconds, err := flagSet.GetInt64(ArgSeconds)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Sleeping for %d seconds", seconds))
	time.Sleep(time.Second * time.Duration(seconds))
	os.Exit(0)
	return nil
}
