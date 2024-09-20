/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package nodes

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/oracle/coherence-operator/pkg/operator"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetExactLabelForNode is a GET request that returns the node label on a k8s node
func GetExactLabelForNode(ctx context.Context, c kubernetes.Interface, name, label string, log logr.Logger) (string, error) {
	var prefix []string
	var labels []string
	labels = append(labels, label)
	value, _, err := GetLabelForNode(ctx, c, name, labels, prefix, log)
	return value, err
}

// GetLabelForNode is a GET request that returns the node label on a k8s node
func GetLabelForNode(ctx context.Context, c kubernetes.Interface, name string, labels, prefixLabels []string, log logr.Logger) (string, string, error) {
	var value string
	labelUsed := "<None>"
	var prefixUsed = "<None>"
	var err error

	if operator.IsNodeLookupEnabled() {
		var node *corev1.Node
		node, err = c.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err == nil {
			var ok bool

			prefixValue := ""
			for _, label := range prefixLabels {
				if prefix, ok := node.Labels[label]; ok && prefix != "" {
					labelUsed = label
					prefixValue = prefix + "-"
					break
				}
			}

			for _, label := range labels {
				if value, ok = node.Labels[label]; ok && value != "" {
					labelUsed = label
					value = prefixValue + value
					break
				}
			}
		} else {
			if apierrors.IsNotFound(err) {
				log.Info("GET query for node labels - NotFound", "node", name, "label", labelUsed, "prefix", prefixUsed, "value", value)
			} else {
				log.Error(err, "GET query for node labels - Error", "node", name, "label", labelUsed, "prefix", prefixUsed, "value", value)
			}
			value = ""
			labelUsed = ""
		}
	} else {
		log.Info("Node labels lookup disabled", "node", name, "label", labelUsed, "prefix", prefixUsed, "value", value)
		value = ""
	}
	return value, labelUsed, err
}
