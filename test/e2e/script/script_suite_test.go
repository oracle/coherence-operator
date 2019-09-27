/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package script

import (
	"context"
	"encoding/json"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
	"k8s.io/utils/pointer"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	framework.MainEntry(m)
}

type AppData struct {
	Args []string          `json:"args,omitempty"`
	Env  map[string]string `json:"env,omitempty"`
}

func (in *AppData) GetEnv(name string) string {
	if in == nil {
		return ""
	}
	return in.Env[name]
}

func (in *AppData) GetSystemProperty(name string) string {
	if in == nil || len(in.Args) == 0 {
		return ""
	}

	prefix := fmt.Sprintf("-D%s=", name)
	for _, arg := range in.Args {
		if strings.HasPrefix(arg, prefix) {
			return arg[len(prefix):]
		}
	}

	return ""
}

func (in *AppData) FindJvmOption(prefix string) []string {
	var opts []string
	if in == nil || len(in.Args) == 0 {
		return opts
	}

	for _, arg := range in.Args {
		if strings.HasPrefix(arg, prefix) {
			opts = append(opts, arg)
		}
	}

	return opts
}

func RunScript(t *testing.T, role v1.CoherenceRoleSpec) (*AppData, *v1.CoherenceCluster, error) {
	var err error

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()

	// Fix to only one replica
	role.SetReplicas(1)
	app := role.Application
	if app == nil {
		app = &v1.ApplicationSpec{}
	}
	app.Type = pointer.StringPtr("op-test")
	role.Application = app

	// Fix the readiness probe to speed up the ready check
	probe := role.ReadinessProbe
	if probe == nil {
		probe = &v1.ReadinessProbeSpec{}
	}
	probe.InitialDelaySeconds = pointer.Int32Ptr(2)
	probe.PeriodSeconds = pointer.Int32Ptr(1)
	role.ReadinessProbe = probe

	// generate a unique cluster name
	name := fmt.Sprintf("test-%d", time.Now().UnixNano()/1000000)

	cluster, err := helper.NewCoherenceCluster(ns)
	if err != nil {
		return nil, nil, err
	}

	roleSpec := role.DeepCopyWithDefaults(&cluster.Spec.CoherenceRoleSpec)
	cluster.Spec.CoherenceRoleSpec = *roleSpec
	cluster.SetName(name)

	f := framework.Global
	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	if err != nil {
		return nil, nil, err
	}

	_, err = helper.WaitForStatefulSetForRole(f.KubeClient, ns, &cluster, cluster.Spec.CoherenceRoleSpec, time.Second*5, time.Minute*2, t)
	if err != nil {
		return nil, nil, err
	}

	pods, err := helper.ListCoherencePodsForCluster(f.KubeClient, ns, cluster.Name)
	if err != nil {
		return nil, nil, err
	}

	if len(pods) == 0 {
		return nil, nil, fmt.Errorf("no pods found for cluster %s", cluster.Name)
	}

	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	if err != nil {
		return nil, nil, err
	}
	defer pf.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d/", ports["health"])
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("http request returned status %d", resp.StatusCode)
	}

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	data := &AppData{}
	err = json.Unmarshal(j, data)
	return data, &cluster, err
}
