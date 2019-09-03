/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

// Initialise the canary test in the role being scaled.
func StartCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryStart", http.MethodPut)
}

// Invoke the canary test in the role being scaled.
func CheckCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryCheck", http.MethodGet)
}

// Make a canary ReST PUT call to Pod zero of the role.
func canary(namespace, clusterName, roleName, endpoint, method string) error {
	podName := fmt.Sprintf("%s-%s-0", clusterName, roleName)
	f := framework.Global

	pod, err := f.KubeClient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	forwarder, ports, err := StartPortForwarderForPod(pod)
	if err != nil {
		return err
	}

	defer forwarder.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d/%s", ports["rest"], endpoint)
	client := &http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(""))
	if err != nil {
		return err
	}

	request.ContentLength = 0
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected http response %d but received %d from '%s'", http.StatusOK, resp.StatusCode, url)
	}

	return nil
}
