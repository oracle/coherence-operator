/*
 * Copyright (c) 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

const (
	// CommandNode is the argument to launch a Node label reader.
	CommandNode = "node"

	// ArgNode is the name of the node to query
	ArgNode = "node-name"
	// ArgDir is the directory to write files to
	ArgDir = "dir"
	// ArkKubeConfig is the location of the k8s config
	ArkKubeConfig = "kubeconfig"
)

// nodeCommand reads node labels into files in a directory
func nodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNode,
		Short: "Read node labels into files in a directory",
		Long:  "Read node labels into files in a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeNodeQuery(cmd)
		},
	}

	path, err := os.Getwd()
	if err != nil {
		log.Error(err, "could not obtain working directory")
	}

	fmt.Println(path) // for example /home/user
	flagSet := cmd.Flags()
	flagSet.String(ArgNode, "", "The name of the Kubernetes node to obtain labels for")
	flagSet.String(ArgDir, path, "The directory to write the label files to")

	if home := homedir.HomeDir(); home != "" {
		flagSet.String(ArkKubeConfig, filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flagSet.String(ArkKubeConfig, "", "absolute path to the kubeconfig file")
	}

	return cmd
}

func executeNodeQuery(cmd *cobra.Command) error {
	flagSet := cmd.Flags()

	nodeName, err := flagSet.GetString(ArgNode)
	if err != nil {
		return err
	}

	if nodeName == "" {
		return fmt.Errorf("no %s argument has been set", ArgNode)
	}

	kubeConfig, err := flagSet.GetString(ArkKubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot get Kubernetes config file")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return errors.Wrap(err, "cannot get Kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "cannot get Kubernetes client")
	}

	// Get node object
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "cannot get Kubernetes Node")
	}

	sep := os.PathSeparator
	dir, err := flagSet.GetString(ArgDir)
	if err != nil {
		return errors.Wrap(err, "cannot get output directory")
	}

	for label, value := range node.Labels {
		fileName := fmt.Sprintf("%s%c%s", dir, sep, label)
		fileDir := filepath.Dir(fileName)

		err = os.MkdirAll(fileDir, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to directory file %s", fileDir)
		}

		f, err := os.Create(fileName)
		if err != nil {
			return errors.Wrapf(err, "failed to create file %s", fileName)
		}

		_, err = f.WriteString(value)
		if err != nil {
			return errors.Wrapf(err, "failed to writing to file %s", fileName)
		}
	}

	return nil
}
