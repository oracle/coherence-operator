/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package healthmutator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// This gate requires an isolated API server, never the user's current kube context.
func TestHealthActionAdmission(t *testing.T) {
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		t.Skip("KUBEBUILDER_ASSETS is required for the admission certification gate")
	}
	environment := &envtest.Environment{CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")}, ErrorIfCRDPathMissing: true}
	config, err := environment.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := environment.Stop(); err != nil {
			t.Error(err)
		}
	})
	scheme := runtime.NewScheme()
	if err := coh.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	c, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := c.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "health-test"}}); err != nil {
		t.Fatal(err)
	}
	valid := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "valid", Namespace: "health-test"}}
	valid.Spec.SuspendProbe = &coh.Probe{HTTP: &coh.HTTPAction{Method: "PUT", HTTPEndpoint: coh.HTTPEndpoint{Path: "/suspend", Port: intstr.FromInt(6676)}}}
	if err := c.Create(ctx, valid); err != nil {
		t.Fatal(err)
	}
	stored := &coh.Coherence{}
	if err := c.Get(ctx, client.ObjectKeyFromObject(valid), stored); err != nil {
		t.Fatal(err)
	}
	if stored.Spec.SuspendProbe.HTTP == nil || stored.Spec.SuspendProbe.HTTP.Method != "PUT" {
		t.Fatal("HTTP action was pruned")
	}
	invalid := valid.DeepCopy()
	invalid.ResourceVersion = ""
	invalid.Name = "invalid"
	invalid.Spec.SuspendProbe.Exec = &corev1.ExecAction{Command: []string{"true"}}
	if err := c.Create(ctx, invalid); err == nil {
		t.Fatal("admission accepted multiple actions")
	}
	selector := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "health-auth-v1"}, Key: "value"}
	authenticated := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "authenticated-pod-probe", Namespace: "health-test"}}
	authenticated.Spec.SuspendProbe = &coh.Probe{HTTP: &coh.HTTPAction{
		HTTPEndpoint: coh.HTTPEndpoint{Path: "/suspend", Port: intstr.FromInt(6676), Scheme: corev1.URISchemeHTTPS},
		Method:       "PUT",
		BasicAuth:    &coh.HTTPBasicAuth{Username: selector, Password: selector},
	}}
	if err := c.Create(ctx, authenticated); err != nil {
		t.Fatalf("admission rejected an authenticated target-Pod probe: %v", err)
	}
	customHost := authenticated.DeepCopy()
	customHost.ResourceVersion = ""
	customHost.Name = "authenticated-custom-host-probe"
	customHost.Spec.SuspendProbe.HTTP.Host = "health.example"
	if err := c.Create(ctx, customHost); err == nil {
		t.Fatal("admission accepted Basic authentication with a custom host")
	}
	legacyEmpty := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "legacy-empty-probe", Namespace: "health-test"}}
	legacyEmpty.Spec.SuspendProbe = &coh.Probe{}
	if err := c.Create(ctx, legacyEmpty); err != nil {
		t.Fatalf("admission rejected a legacy empty probe: %v", err)
	}
	legacyEmpty.Annotations = map[string]string{"compatibility-test": "updated"}
	if err := c.Update(ctx, legacyEmpty); err != nil {
		t.Fatalf("admission rejected an update containing a legacy empty probe: %v", err)
	}
	legacyMultiple := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "legacy-multiple-probe-actions", Namespace: "health-test"}}
	legacyMultiple.Spec.SuspendProbe = &coh.Probe{ProbeHandler: corev1.ProbeHandler{
		Exec:    &corev1.ExecAction{Command: []string{"true"}},
		HTTPGet: &corev1.HTTPGetAction{Path: "/suspend", Port: intstr.FromInt(6676)},
	}}
	if err := c.Create(ctx, legacyMultiple); err != nil {
		t.Fatalf("admission rejected legacy multiple probe actions: %v", err)
	}
	invalid = valid.DeepCopy()
	invalid.ResourceVersion = ""
	invalid.Name = "incomplete-secure"
	invalid.Spec.SuspendProbe = nil
	invalid.Spec.HealthMutatorConnection = &coh.HealthMutatorConnection{Secure: true}
	if err := c.Create(ctx, invalid); err == nil {
		t.Fatal("admission accepted incomplete secure configuration")
	}
	managed := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "managed", Namespace: "health-test"}}
	managed.Spec.HealthMutatorConnection = &coh.HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &coh.HTTPBasicAuth{Username: selector, Password: selector},
		TLS:       &coh.HTTPClientTLS{CASecret: corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "health-ca-v1"}, Key: "ca.pem"}, ServerName: "health.example"},
		JAAS:      &coh.HealthJAAS{},
		ServerTLS: &coh.SSLSpec{
			Enabled:              ptr.To(true),
			RequireClientCert:    ptr.To(false),
			Secrets:              ptr.To("health-stores-v1"),
			KeyStore:             ptr.To("server.jks"),
			KeyStorePasswordFile: ptr.To("store-password"),
			KeyPasswordFile:      ptr.To("store-password"),
		},
	}
	if err := c.Create(ctx, managed); err != nil {
		t.Fatal(err)
	}
	storedManaged := &coh.Coherence{}
	if err := c.Get(ctx, client.ObjectKeyFromObject(managed), storedManaged); err != nil {
		t.Fatal(err)
	}
	connection := storedManaged.Spec.HealthMutatorConnection
	if connection == nil || connection.BasicAuth == nil || connection.TLS == nil || connection.ServerTLS == nil || connection.JAAS == nil || connection.JAAS.Mode != "packaged" {
		t.Fatal("complete healthMutatorConnection was pruned or not defaulted")
	}
	if connection.ServerTLS.TrustStore != nil || connection.ServerTLS.TrustStorePasswordFile != nil {
		t.Fatal("admission added a server truststore to one-way TLS")
	}
	legacy := managed.DeepCopy()
	legacy.ResourceVersion = ""
	legacy.Name = "managed-with-legacy-truststore"
	legacy.Spec.HealthMutatorConnection.ServerTLS.TrustStore = ptr.To("trust.jks")
	legacy.Spec.HealthMutatorConnection.ServerTLS.TrustStorePasswordFile = ptr.To("store-password")
	if err := c.Create(ctx, legacy); err != nil {
		t.Fatalf("admission rejected previously valid truststore fields: %v", err)
	}
	conflict := managed.DeepCopy()
	conflict.ResourceVersion = ""
	conflict.Name = "managed-conflict"
	conflict.Spec.SuspendProbe = valid.Spec.SuspendProbe.DeepCopy()
	if err := c.Create(ctx, conflict); err == nil {
		t.Fatal("admission accepted managed security with a custom suspend probe")
	}
	job := &coh.CoherenceJob{ObjectMeta: metav1.ObjectMeta{Name: "job", Namespace: "health-test"}}
	job.Spec.ReadyAction = &coh.CoherenceJobProbe{Probe: *valid.Spec.SuspendProbe.DeepCopy()}
	if err := c.Create(ctx, job); err != nil {
		t.Fatal(err)
	}
	storedJob := &coh.CoherenceJob{}
	if err := c.Get(ctx, client.ObjectKeyFromObject(job), storedJob); err != nil {
		t.Fatal(err)
	}
	if storedJob.Spec.ReadyAction == nil || storedJob.Spec.ReadyAction.HTTP == nil {
		t.Fatal("Job HTTP action was pruned")
	}
	rawJob := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": coh.GroupVersion.String(),
		"kind":       "CoherenceJob",
		"metadata": map[string]interface{}{
			"name":      "job-with-managed-field",
			"namespace": "health-test",
		},
		"spec": map[string]interface{}{
			"healthMutatorConnection": map[string]interface{}{"secure": false},
		},
	}}
	if err := c.Create(ctx, rawJob); err != nil {
		t.Fatal(err)
	}
	storedRawJob := &unstructured.Unstructured{}
	storedRawJob.SetAPIVersion(coh.GroupVersion.String())
	storedRawJob.SetKind("CoherenceJob")
	if err := c.Get(ctx, client.ObjectKeyFromObject(rawJob), storedRawJob); err != nil {
		t.Fatal(err)
	}
	if _, found, err := unstructured.NestedFieldNoCopy(storedRawJob.Object, "spec", "healthMutatorConnection"); err != nil || found {
		t.Fatalf("managed StatefulSet field survived CoherenceJob admission: found=%t err=%v", found, err)
	}

	t.Run("list-map key defaults", func(t *testing.T) {
		resource := newUnstructuredResource("Coherence", "list-map-defaults", map[string]interface{}{
			"network": map[string]interface{}{
				"dnsConfig": map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{"value": "2"},
						map[string]interface{}{"name": "ndots", "value": "1"},
					},
				},
			},
			"actions": []interface{}{
				map[string]interface{}{
					"name": "pull-secret-defaulting",
					"job": map[string]interface{}{
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"restartPolicy": "Never",
									"containers": []interface{}{
										map[string]interface{}{"name": "test", "image": "example.invalid/test:latest"},
									},
									"imagePullSecrets": []interface{}{
										map[string]interface{}{},
										map[string]interface{}{"name": "registry-credentials"},
									},
								},
							},
						},
					},
				},
			},
		})
		storedResource := createAndGetUnstructured(t, ctx, c, resource)
		assertNestedListMapNames(t, storedResource.Object, []string{"", "ndots"}, "spec", "network", "dnsConfig", "options")
		storedActions, found, err := unstructured.NestedSlice(storedResource.Object, "spec", "actions")
		if err != nil || !found || len(storedActions) != 1 {
			t.Fatalf("actions were not stored: count=%d found=%t err=%v", len(storedActions), found, err)
		}
		storedAction, ok := storedActions[0].(map[string]interface{})
		if !ok {
			t.Fatalf("stored action has type %T", storedActions[0])
		}
		assertNestedListMapNames(t, storedAction, []string{"", "registry-credentials"}, "job", "spec", "template", "spec", "imagePullSecrets")

		jobResource := newUnstructuredResource("CoherenceJob", "job-list-map-defaults", map[string]interface{}{
			"network": map[string]interface{}{
				"dnsConfig": map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{"value": "2"},
						map[string]interface{}{"name": "ndots", "value": "1"},
					},
				},
			},
		})
		storedJobResource := createAndGetUnstructured(t, ctx, c, jobResource)
		assertNestedListMapNames(t, storedJobResource.Object, []string{"", "ndots"}, "spec", "network", "dnsConfig", "options")
	})

	t.Run("actions limit", func(t *testing.T) {
		actions := make([]coh.Action, 128)
		for i := range actions {
			actions[i].Name = fmt.Sprintf("action-%03d", i)
		}
		atLimit := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "actions-at-limit", Namespace: "health-test"}}
		atLimit.Spec.Actions = actions
		if err := c.Create(ctx, atLimit); err != nil {
			t.Fatalf("admission rejected 128 actions: %v", err)
		}
		storedAtLimit := &coh.Coherence{}
		if err := c.Get(ctx, client.ObjectKeyFromObject(atLimit), storedAtLimit); err != nil {
			t.Fatal(err)
		}
		if len(storedAtLimit.Spec.Actions) != 128 {
			t.Fatalf("stored %d actions, want 128", len(storedAtLimit.Spec.Actions))
		}

		overLimit := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "actions-over-limit", Namespace: "health-test"}}
		overLimit.Spec.Actions = append(append([]coh.Action{}, actions...), coh.Action{Name: "action-128"})
		if err := c.Create(ctx, overLimit); err == nil {
			t.Fatal("admission accepted 129 actions")
		}
	})

	if os.Getenv("HEALTH_MUTATOR_RUNNER_SUITE") == "true" {
		kubeconfig := clientcmdapi.Config{Clusters: map[string]*clientcmdapi.Cluster{"test": {Server: config.Host, CertificateAuthorityData: config.CAData}}, AuthInfos: map[string]*clientcmdapi.AuthInfo{"test": {ClientCertificateData: config.CertData, ClientKeyData: config.KeyData}}, Contexts: map[string]*clientcmdapi.Context{"test": {Cluster: "test", AuthInfo: "test"}}, CurrentContext: "test"}
		path := filepath.Join(t.TempDir(), "kubeconfig")
		if err := clientcmd.WriteToFile(kubeconfig, path); err != nil {
			t.Fatal(err)
		}
		command := exec.Command("go", "test", "./pkg/runner")
		command.Dir = filepath.Join("..", "..")
		command.Env = append(os.Environ(), "KUBECONFIG="+path)
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("runner suite failed: %s", output)
		}
	}

}

func newUnstructuredResource(kind, name string, spec map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": coh.GroupVersion.String(),
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": "health-test",
		},
		"spec": spec,
	}}
}

func createAndGetUnstructured(t *testing.T, ctx context.Context, c client.Client, resource *unstructured.Unstructured) *unstructured.Unstructured {
	t.Helper()
	if err := c.Create(ctx, resource); err != nil {
		t.Fatal(err)
	}
	stored := &unstructured.Unstructured{}
	stored.SetAPIVersion(resource.GetAPIVersion())
	stored.SetKind(resource.GetKind())
	if err := c.Get(ctx, client.ObjectKeyFromObject(resource), stored); err != nil {
		t.Fatal(err)
	}
	return stored
}

func assertNestedListMapNames(t *testing.T, object map[string]interface{}, expected []string, fields ...string) {
	t.Helper()
	items, found, err := unstructured.NestedSlice(object, fields...)
	if err != nil || !found {
		t.Fatalf("list-map field %v was not stored: found=%t err=%v", fields, found, err)
	}
	if len(items) != len(expected) {
		t.Fatalf("list-map field %v has %d items, want %d", fields, len(items), len(expected))
	}
	for i, item := range items {
		entry, ok := item.(map[string]interface{})
		if !ok {
			t.Fatalf("list-map field %v item %d has type %T", fields, i, item)
		}
		name, found, err := unstructured.NestedString(entry, "name")
		if err != nil || !found || name != expected[i] {
			t.Fatalf("list-map field %v item %d name=%q found=%t err=%v, want %q", fields, i, name, found, err, expected[i])
		}
	}
}
