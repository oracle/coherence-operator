/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

var k8sCfg *rest.Config
var k8sClient client.Client
var kubeClient kubernetes.Interface
var testEnv *envtest.Environment
var k8sManager ctrl.Manager
var testLogger logr.Logger

func GetK8sClient() client.Client {
	return k8sClient
}

func TestMain(m *testing.M) {
	testLogger = zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout))

	// run the tests
	if err := beforeSuite(); err != nil {
		testLogger.Error(err, "error running before suite tasks")
		os.Exit(1)
	}
	exitCode := m.Run()
	if err := afterSuite(); err != nil {
		testLogger.Error(err, "error running after suite tasks")
		exitCode = 1
	}

	os.Exit(exitCode)
}

func beforeSuite() error {
	logf.SetLogger(testLogger)

	// We need a real cluster for these tests
	useCluster := true

	testLogger.Info("bootstrapping test environment")
	testEnv = &envtest.Environment{
		UseExistingCluster:       &useCluster,
		AttachControlPlaneOutput: true,
	}

	var err error
	k8sCfg, err = testEnv.Start()
	if err != nil {
		return err
	}

	err = corev1.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}
	err = coh.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	testEnv.CRDs = []runtime.Object{
		&corev1.CustomResourceDefinition{},
	}

	k8sManager, err = ctrl.NewManager(k8sCfg, ctrl.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}

	// Create the Coherence controller
	err = (&controllers.CoherenceReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("Coherence.controller"),
	}).SetupWithManager(k8sManager)
	if err != nil {
		return err
	}

	// Start the manager, which will start the controller
	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
	}()

	k8sClient = k8sManager.GetClient()
	kubeClient, err = kubernetes.NewForConfig(k8sCfg)

	return err
}

func afterSuite() error {
	testLogger.Info("tearing down the test environment")
	Cleanup()
	return testEnv.Stop()
}

func Cleanup() {
	ns := helper.GetTestNamespace()
	err := k8sClient.DeleteAllOf(context.Background(), &coh.Coherence{}, client.InNamespace(ns))
	if err != nil {
		testLogger.Info("error tearing down the test environment: " + err.Error())
	}
}
