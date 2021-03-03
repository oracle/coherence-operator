/*
 * Copyright (c) 2019, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"context"
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

// Initialise the canary test in the deployment being scaled.
func StartCanary(ctx TestContext, namespace, deploymentName string) error {
	return canary(ctx, namespace, deploymentName, "canaryStart", http.MethodPut)
}

// Invoke the canary test in the deployment being scaled.
func CheckCanary(ctx TestContext, namespace, deploymentName string) error {
	return canary(ctx, namespace, deploymentName, "canaryCheck", http.MethodGet)
}

// Clear the canary test in the deployment being scaled.
func ClearCanary(ctx TestContext, namespace, deploymentName string) error {
	return canary(ctx, namespace, deploymentName, "canaryClear", http.MethodPost)
}

// Make a canary ReST PUT call to Pod zero of the deployment.
func canary(ctx TestContext, namespace, deploymentName, endpoint, method string) error {
	podName := fmt.Sprintf("%s-0", deploymentName)

	pod, err := ctx.KubeClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	forwarder, ports, err := StartPortForwarderForPod(pod)
	if err != nil {
		return err
	}

	defer forwarder.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d/%s", ports["rest"], endpoint)
	client := &http.Client{Timeout: time.Minute * 1}

	var resp *http.Response
	var request *http.Request
	// try a max of 5 times

	for i := 0; i < 5; i++ {
		request, err = http.NewRequest(method, url, strings.NewReader(""))
		if err == nil {
			request.ContentLength = 0
			resp, err = client.Do(request)
			if err == nil {
				break
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
