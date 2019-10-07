/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

// Initialise the canary test in the role being scaled.
func StartCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryStart", http.MethodPut)
}

// Invoke the canary test in the role being scaled.
func CheckCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryCheck", http.MethodGet)
}

// Clear the canary test in the role being scaled.
func ClearCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryClear", http.MethodPost)
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

	var resp *http.Response
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		request, err := http.NewRequest(method, url, strings.NewReader(""))
		if err == nil {
			request.ContentLength = 0
			resp, err = client.Do(request)
			if err == nil {
				break;
			}
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := ""
		if err == nil {
			bodyString = string(bodyBytes)
		}
		return fmt.Errorf("expected http response %d but received %d from '%s' with body '%s'",
			http.StatusOK, resp.StatusCode, url, bodyString)
	}

	return nil
}
