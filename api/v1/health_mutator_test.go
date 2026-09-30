/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package v1

import (
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"reflect"
	"strings"
	"testing"
)

func TestDisabledManagedHealthDoesNotChangePodMetadata(t *testing.T) {
	for name, connection := range map[string]*HealthMutatorConnection{
		"omitted": nil,
		"false":   {Secure: false},
	} {
		t.Run(name, func(t *testing.T) {
			pod := &corev1.PodTemplateSpec{}
			spec := &CoherenceStatefulSetResourceSpec{HealthMutatorConnection: connection}
			spec.ConfigureHealthMutators(pod)
			if pod.Annotations != nil {
				t.Fatalf("disabled managed Health allocated annotations: %#v", pod.Annotations)
			}
		})
	}

	annotations := map[string]string{"example.com/owner": "application"}
	pod := &corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Annotations: annotations}}
	spec := &CoherenceStatefulSetResourceSpec{HealthMutatorConnection: &HealthMutatorConnection{Secure: false}}
	spec.ConfigureHealthMutators(pod)
	if !reflect.DeepEqual(pod.Annotations, map[string]string{"example.com/owner": "application"}) {
		t.Fatalf("disabled managed Health changed unrelated annotations: %#v", pod.Annotations)
	}
}

func TestManagedHealthStatefulSetTemplateTransitions(t *testing.T) {
	deployment := &Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "health-transition"},
		Spec: CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: CoherenceResourceSpec{StartupProbe: &ReadinessProbeSpec{}},
		},
	}

	omitted := deployment.Spec.CreateStatefulSet(deployment).Spec.Template
	assertDisabledManagedHealthTemplate(t, &omitted)

	deployment.Spec.HealthMutatorConnection = &HealthMutatorConnection{Secure: false}
	disabled := deployment.Spec.CreateStatefulSet(deployment).Spec.Template
	assertDisabledManagedHealthTemplate(t, &disabled)
	if !reflect.DeepEqual(omitted, disabled) {
		t.Fatal("omitted and secure:false managed Health generated different Pod templates")
	}

	secureConnection := validManagedHealthConnection()
	if err := secureConnection.Validate(); err != nil {
		t.Fatalf("test secure connection is invalid: %v", err)
	}
	deployment.Spec.HealthMutatorConnection = secureConnection
	secure := deployment.Spec.CreateStatefulSet(deployment).Spec.Template
	assertSecureManagedHealthTemplate(t, &secure)
	if reflect.DeepEqual(disabled, secure) {
		t.Fatal("enabling managed Health did not change the Pod template")
	}

	deployment.Spec.HealthMutatorConnection = &HealthMutatorConnection{Secure: false}
	disabledAfterSecure := deployment.Spec.CreateStatefulSet(deployment).Spec.Template
	assertDisabledManagedHealthTemplate(t, &disabledAfterSecure)
	if !reflect.DeepEqual(disabled, disabledAfterSecure) {
		t.Fatal("regenerating after disabling managed Health did not restore the disabled template")
	}

	deployment.Spec.HealthMutatorConnection = validManagedHealthConnection()
	secureAfterDisabled := deployment.Spec.CreateStatefulSet(deployment).Spec.Template
	if !reflect.DeepEqual(secure, secureAfterDisabled) {
		t.Fatal("re-enabling managed Health did not restore the secure Pod template")
	}
}

func validManagedHealthConnection() *HealthMutatorConnection {
	username := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "health-auth"}, Key: "username"}
	password := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "health-auth"}, Key: "password"}
	ca := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "health-ca"}, Key: "ca.crt"}
	return &HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &HTTPBasicAuth{Username: username, Password: password},
		TLS:       &HTTPClientTLS{CASecret: ca, ServerName: "health.example"},
		ServerTLS: &SSLSpec{
			Enabled:              ptr.To(true),
			Secrets:              ptr.To("health-server-tls"),
			KeyStore:             ptr.To("server.jks"),
			KeyStorePasswordFile: ptr.To("store-password"),
			KeyPasswordFile:      ptr.To("key-password"),
		},
	}
}

func assertDisabledManagedHealthTemplate(t *testing.T, pod *corev1.PodTemplateSpec) {
	t.Helper()
	if _, found := pod.Annotations[HealthMutatorAnnotation]; found {
		t.Fatal("disabled managed Health generated a connection annotation")
	}
	assertNoManagedHealthWiring(t, &pod.Spec)
	coherence := FindContainerInPodTemplate(ContainerNameCoherence, pod)
	if coherence == nil {
		t.Fatal("missing Coherence container")
	}
	for name, probe := range map[string]*corev1.Probe{
		"readiness": coherence.ReadinessProbe,
		"liveness":  coherence.LivenessProbe,
		"startup":   coherence.StartupProbe,
	} {
		if probe == nil || probe.HTTPGet == nil || probe.HTTPGet.Scheme != corev1.URISchemeHTTP {
			t.Fatalf("disabled managed Health did not retain the generated HTTP %s probe", name)
		}
	}
}

func assertNoManagedHealthWiring(t *testing.T, pod *corev1.PodSpec) {
	t.Helper()
	for _, volume := range pod.Volumes {
		if strings.HasPrefix(volume.Name, "health-mutator-") {
			t.Fatalf("disabled managed Health generated volume %q", volume.Name)
		}
	}
	containers := append(append([]corev1.Container{}, pod.InitContainers...), pod.Containers...)
	for _, container := range containers {
		for _, env := range container.Env {
			if env.Name == HealthMutatorEnv {
				t.Fatalf("disabled managed Health configured %s on container %q", HealthMutatorEnv, container.Name)
			}
		}
		for _, mount := range container.VolumeMounts {
			if strings.HasPrefix(mount.Name, "health-mutator-") {
				t.Fatalf("disabled managed Health generated mount %q on container %q", mount.Name, container.Name)
			}
		}
	}
}

func assertSecureManagedHealthTemplate(t *testing.T, pod *corev1.PodTemplateSpec) {
	t.Helper()
	if value := pod.Annotations[HealthMutatorAnnotation]; !strings.Contains(value, `"secure":true`) {
		t.Fatalf("secure managed Health annotation is missing its connection: %q", value)
	}
	for _, name := range []string{"health-mutator-username", "health-mutator-password", "health-mutator-tls"} {
		found := false
		for _, volume := range pod.Spec.Volumes {
			if volume.Name == name {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("secure managed Health did not generate volume %q", name)
		}
	}
	for _, container := range []*corev1.Container{
		FindContainerInPodTemplate(ContainerNameCoherence, pod),
		FindInitContainerInPodTemplate(ContainerNameOperatorConfig, pod),
	} {
		if container == nil {
			t.Fatal("missing managed Health target container")
		}
		if !hasEnv(container.Env, HealthMutatorEnv) {
			t.Fatalf("secure managed Health did not configure %s on container %q", HealthMutatorEnv, container.Name)
		}
		for _, name := range []string{"health-mutator-username", "health-mutator-password", "health-mutator-tls"} {
			if !hasVolumeMount(container.VolumeMounts, name) {
				t.Fatalf("secure managed Health did not mount %q on container %q", name, container.Name)
			}
		}
	}
	coherence := FindContainerInPodTemplate(ContainerNameCoherence, pod)
	for name, probe := range map[string]*corev1.Probe{
		"readiness": coherence.ReadinessProbe,
		"liveness":  coherence.LivenessProbe,
		"startup":   coherence.StartupProbe,
	} {
		if probe == nil || probe.HTTPGet == nil || probe.HTTPGet.Scheme != corev1.URISchemeHTTPS {
			t.Fatalf("secure managed Health did not convert the generated %s probe to HTTPS", name)
		}
	}
}

func hasEnv(env []corev1.EnvVar, name string) bool {
	for _, value := range env {
		if value.Name == name {
			return true
		}
	}
	return false
}

func hasVolumeMount(mounts []corev1.VolumeMount, name string) bool {
	for _, mount := range mounts {
		if mount.Name == name {
			return true
		}
	}
	return false
}

func TestManagedHealthWiringIsStatefulSetOnly(t *testing.T) {
	selector := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "auth-v1"}, Key: "value"}
	connection := &HealthMutatorConnection{Secure: true, BasicAuth: &HTTPBasicAuth{Username: selector, Password: selector}, TLS: &HTTPClientTLS{CASecret: selector, ServerName: "health.example"}, ServerTLS: &SSLSpec{Enabled: ptr.To(true), Secrets: ptr.To("stores-v1"), KeyStore: ptr.To("server.jks"), KeyStorePasswordFile: ptr.To("password"), KeyPasswordFile: ptr.To("password")}}
	if err := connection.Validate(); err != nil {
		t.Fatal(err)
	}
	spec := &CoherenceStatefulSetResourceSpec{HealthMutatorConnection: connection}
	pod := &corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: ContainerNameCoherence}, {Name: "sidecar"}}}}
	pod.Spec.InitContainers = []corev1.Container{{Name: ContainerNameOperatorConfig}, {Name: ContainerNameOperatorInit}}
	spec.ConfigureHealthMutators(pod)
	if len(pod.Spec.Volumes) != 3 || len(pod.Spec.InitContainers[0].VolumeMounts) != 3 || len(pod.Spec.InitContainers[1].VolumeMounts) != 0 {
		t.Fatal("configuration init container wiring is incorrect")
	}
	if len(pod.Spec.Containers[0].VolumeMounts) != 3 || len(pod.Spec.Containers[1].VolumeMounts) != 0 {
		t.Fatal("credential mounts must be confined to the member")
	}
	if !strings.Contains(pod.Annotations[HealthMutatorAnnotation], "auth-v1") {
		t.Fatal("missing active generation references")
	}
	data, _ := json.Marshal(&CoherenceJobResourceSpec{})
	if strings.Contains(string(data), "healthMutatorConnection") {
		t.Fatal("managed fields leaked into Job")
	}
	defaultProbe := spec.GetDefaultSuspendProbe()
	if defaultProbe.HTTP == nil || defaultProbe.HTTP.Method != "PUT" || defaultProbe.HTTP.Path != "/suspend" {
		t.Fatal("invalid generated suspend request")
	}
}

func TestManagedHealthServerTLSRequiresOnlyIdentityMaterial(t *testing.T) {
	selector := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "auth-v1"}, Key: "value"}
	connection := &HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &HTTPBasicAuth{Username: selector, Password: selector},
		TLS:       &HTTPClientTLS{CASecret: selector, ServerName: "health.example"},
		ServerTLS: &SSLSpec{
			Enabled:              ptr.To(true),
			KeyStore:             ptr.To("server.jks"),
			KeyStorePasswordFile: ptr.To("store-password"),
			KeyPasswordFile:      ptr.To("key-password"),
		},
	}
	if err := connection.Validate(); err != nil {
		t.Fatalf("one-way TLS without a truststore was rejected: %v", err)
	}
	for name, clear := range map[string]func(*SSLSpec){
		"keystore":          func(ssl *SSLSpec) { ssl.KeyStore = nil },
		"keystore password": func(ssl *SSLSpec) { ssl.KeyStorePasswordFile = nil },
		"key password":      func(ssl *SSLSpec) { ssl.KeyPasswordFile = nil },
	} {
		t.Run(name, func(t *testing.T) {
			invalid := connection.DeepCopy()
			clear(invalid.ServerTLS)
			if err := invalid.Validate(); err == nil {
				t.Fatal("missing server identity material was accepted")
			}
		})
	}
	invalid := connection.DeepCopy()
	invalid.ServerTLS.RequireClientCert = ptr.To(true)
	if err := invalid.Validate(); err == nil {
		t.Fatal("client certificate requirement was accepted")
	}
	connection.ServerTLS.TrustStore = ptr.To("legacy-trust.jks")
	connection.ServerTLS.TrustStorePasswordFile = ptr.To("legacy-password")
	if err := connection.Validate(); err != nil {
		t.Fatalf("previously valid truststore fields were rejected: %v", err)
	}
}

func TestManagedHealthUsesHTTPSForGeneratedCustomPortProbes(t *testing.T) {
	selector := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "auth-v1"}, Key: "value"}
	connection := &HealthMutatorConnection{Secure: true, BasicAuth: &HTTPBasicAuth{Username: selector, Password: selector}, TLS: &HTTPClientTLS{CASecret: selector, ServerName: "health.example"}, ServerTLS: &SSLSpec{Enabled: ptr.To(true), Secrets: ptr.To("stores-v1"), KeyStore: ptr.To("server.jks"), KeyStorePasswordFile: ptr.To("password"), KeyPasswordFile: ptr.To("password")}}
	healthPort := int32(8089)
	spec := &CoherenceStatefulSetResourceSpec{
		HealthMutatorConnection: connection,
		CoherenceResourceSpec: CoherenceResourceSpec{
			HealthPort:   &healthPort,
			StartupProbe: &ReadinessProbeSpec{},
		},
	}
	container := spec.CreateCoherenceContainer(&Coherence{})
	pod := &corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{container}}}
	spec.ConfigureHealthMutators(pod)
	for name, probe := range map[string]*corev1.Probe{
		"readiness": pod.Spec.Containers[0].ReadinessProbe,
		"liveness":  pod.Spec.Containers[0].LivenessProbe,
		"startup":   pod.Spec.Containers[0].StartupProbe,
	} {
		if probe == nil || probe.HTTPGet == nil || probe.HTTPGet.Port.IntValue() != int(healthPort) || probe.HTTPGet.Scheme != corev1.URISchemeHTTPS {
			t.Fatalf("%s probe did not use HTTPS on custom health port", name)
		}
	}
}

func TestManagedHealthPreservesCustomKubeletProbe(t *testing.T) {
	healthPort := int32(8089)
	custom := &corev1.HTTPGetAction{Path: "/custom", Port: intstr.FromInt32(healthPort), Scheme: corev1.URISchemeHTTP}
	spec := &CoherenceStatefulSetResourceSpec{
		HealthMutatorConnection: &HealthMutatorConnection{Secure: true},
		CoherenceResourceSpec: CoherenceResourceSpec{
			HealthPort:     &healthPort,
			ReadinessProbe: &ReadinessProbeSpec{ProbeHandler: ProbeHandler{HTTPGet: custom}},
		},
	}
	container := spec.CreateCoherenceContainer(&Coherence{})
	pod := &corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{container}}}
	spec.ConfigureHealthMutators(pod)
	if pod.Spec.Containers[0].ReadinessProbe.HTTPGet.Scheme != corev1.URISchemeHTTP {
		t.Fatal("managed security changed an explicit custom readiness probe")
	}
}
