/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controller_test

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	cc "github.com/oracle/coherence-operator/controllers"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/pkg/fakes"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

const (
	testCoherenceImage = "oracle/coherence-ce:1.2.3"
	testUtilsImage     = "oracle/operator:1.2.3-utils"
)

func FindContainer(name string, sts *appsv1.StatefulSet) (corev1.Container, bool) {
	for _, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == name {
			return c, true
		}
	}
	return corev1.Container{}, false
}

//func FindInitContainer(name string, sts *appsv1.StatefulSet) (corev1.Container, bool) {
//	for _, c := range sts.Spec.Template.Spec.InitContainers {
//		if c.Name == name {
//			return c, true
//		}
//	}
//	return corev1.Container{}, false
//}

func FindContainerPort(container corev1.Container, name string) (corev1.ContainerPort, bool) {
	for _, port := range container.Ports {
		if port.Name == name {
			return port, true
		}
	}
	return corev1.ContainerPort{}, false
}

func toCoherence(mgr *fakes.FakeManager, obj runtime.Object) (*coh.Coherence, error) {
	c := &coh.Coherence{}
	err := mgr.Scheme.Convert(obj, c, context.TODO())
	return c, err
}

func toSecret(mgr *fakes.FakeManager, obj runtime.Object) (*corev1.Secret, error) {
	s := &corev1.Secret{}
	err := mgr.Scheme.Convert(obj, s, context.TODO())
	return s, err
}

func toService(mgr *fakes.FakeManager, obj runtime.Object) (*corev1.Service, error) {
	s := &corev1.Service{}
	err := mgr.Scheme.Convert(obj, s, context.TODO())
	return s, err
}

func toStatefulSet(mgr *fakes.FakeManager, obj runtime.Object) (*appsv1.StatefulSet, error) {
	s := &appsv1.StatefulSet{}
	err := mgr.Scheme.Convert(obj, s, context.TODO())
	return s, err
}

// Run the specified deployment through a fake reconciler
func Reconcile(t *testing.T, d coh.Coherence) ([]runtime.Object, *fakes.FakeManager) {
	g := NewGomegaWithT(t)

	_, _ = coh.EnsureHashLabel(&d)

	chain, err := NewFakeReconcileChain()
	g.Expect(err).NotTo(HaveOccurred())
	results, err := chain.ReconcileDeployments(d)
	g.Expect(err).NotTo(HaveOccurred())

	// should be one result for the deployment
	result, found := results[d.Name]
	g.Expect(found).To(BeTrue(), "No result found for deployment "+d.Name)
	// result should not be re-queued
	g.Expect(result.Requeue).To(BeFalse(), "Result for deployment "+d.Name+" was re-queued")

	mgr := chain.GetManager()
	resources := mgr.Client.GetCreates()
	finalCoh := &coh.Coherence{}

	// Set the first resource to be the final Coherence resource state
	err = mgr.Client.Get(context.TODO(), d.GetNamespacedName(), finalCoh)
	g.Expect(err).NotTo(HaveOccurred())
	resources[0] = finalCoh

	return resources, mgr
}

type FakeReconcileChain interface {
	ReconcileDeploymentsFromYaml(yamlFile string) ([]coh.Coherence, map[string]reconcile.Result, error)
	ReconcileDeployments(deployments ...coh.Coherence) (map[string]reconcile.Result, error)
	ReconcileExisting(names ...apitypes.NamespacedName) (map[string]reconcile.Result, error)
	GetManager() *fakes.FakeManager
	GetNamespace() string
}

// NewFakeReconcileChain creates a FakeReconcileChain to reconcile clusters.
// This chain effectively reconciles the CoherenceCluster using the Cluster controller
// then each role that was created is reconciled by the Role controller.
func NewFakeReconcileChain() (FakeReconcileChain, error) {
	mgr, err := fakes.NewFakeManager()
	if err != nil {
		return nil, err
	}

	viper.Set(operator.FlagCoherenceImage, testCoherenceImage)
	viper.Set(operator.FlagUtilsImage, testUtilsImage)

	r := &cc.CoherenceReconciler{
		Log: ctrl.Log.WithName("controllers").WithName("Coherence"),
	}
	if err = r.SetupWithManager(mgr); err != nil {
		return nil, err
	}
	r.SetPatchType(apitypes.StrategicMergePatchType)

	fh := &fakeReconcileChain{
		mgr: mgr,
		r:   r,
		ns:  "test-namespace",
	}

	return fh, nil
}

type fakeReconcileChain struct {
	mgr *fakes.FakeManager
	r   *cc.CoherenceReconciler
	ns  string
}

func (in *fakeReconcileChain) GetManager() *fakes.FakeManager {
	if in == nil {
		return nil
	}
	return in.mgr
}

func (in *fakeReconcileChain) GetNamespace() string {
	if in == nil {
		return ""
	}
	return in.ns
}

func (in *fakeReconcileChain) ReconcileDeploymentsFromYaml(yamlFile string) ([]coh.Coherence, map[string]reconcile.Result, error) {
	deployments, err := helper.NewCoherenceFromYaml(in.ns, yamlFile)
	if err != nil {
		return nil, nil, err
	}
	results, err := in.ReconcileDeployments(deployments...)
	return deployments, results, err
}

func (in *fakeReconcileChain) ReconcileDeployments(deployments ...coh.Coherence) (map[string]reconcile.Result, error) {
	var err error
	var names []apitypes.NamespacedName

	for _, d := range deployments {
		err = in.mgr.Client.Get(context.TODO(), d.GetNamespacedName(), &coh.Coherence{})
		if errors.IsNotFound(err) {
			if err = in.mgr.Client.Create(context.TODO(), &d); err != nil {
				return nil, err
			}
		} else {
			if err = in.mgr.Client.Update(context.TODO(), &d); err != nil {
				return nil, err
			}
		}
		names = append(names, d.GetNamespacedName())
	}
	return in.ReconcileExisting(names...)
}

func (in *fakeReconcileChain) ReconcileExisting(names ...apitypes.NamespacedName) (map[string]reconcile.Result, error) {
	results := make(map[string]reconcile.Result)
	ctx := context.TODO()

	for _, name := range names {
		request := reconcile.Request{NamespacedName: name}
		c := coh.Coherence{}
		err := in.r.GetClient().Get(ctx, request.NamespacedName, &c)
		if err != nil {
			return results, err
		}

		if utils.StringArrayDoesNotContain(c.GetFinalizers(), coh.CoherenceFinalizer) {
			// there is no finalizer, so we need to do a call first that will just add the finalizer
			_, err := in.r.Reconcile(ctx, request)
			if err != nil {
				return results, err
			}
		}
		result, err := in.r.Reconcile(ctx, request)
		if err != nil {
			return results, err
		}
		results[name.Name] = result
	}

	return results, nil
}

func AssertStatefulSetCreationEvent(t *testing.T, roleName string, mgr *fakes.FakeManager) {
	g := NewGomegaWithT(t)

	event, found := mgr.NextEvent()
	g.Expect(found).To(BeTrue())
	g.Expect(event.Type).To(Equal(corev1.EventTypeNormal))
	g.Expect(event.Reason).To(Equal(reconciler.EventReasonCreated))
	msg := fmt.Sprintf(statefulset.CreateMessage, roleName)
	g.Expect(event.Message).To(Equal(msg))
}

func AssertNoRemainingEvents(t *testing.T, mgr *fakes.FakeManager) {
	g := NewGomegaWithT(t)
	evt, found := mgr.NextEvent()
	g.Expect(found).To(BeFalse(), fmt.Sprintf("Found unexpected event: %v", evt))
}
