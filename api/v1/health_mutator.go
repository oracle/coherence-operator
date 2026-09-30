/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package v1

import (
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
)

// HealthMutatorAnnotation records connection references for each Pod generation.
const HealthMutatorAnnotation = "coherence.oracle.com/health-mutator-connection"

// HealthMutatorEnv carries non-secret managed configuration to the runner.
const HealthMutatorEnv = "COHERENCE_HEALTH_MUTATOR_CONNECTION"

// HealthMutatorConnection configures authenticated health operations for StatefulSets.
// +kubebuilder:validation:XValidation:rule="!self.secure || (has(self.basicAuth) && has(self.tls) && has(self.serverTLS))",message="secure mode requires credentials, client TLS and server TLS"
// +kubebuilder:validation:XValidation:rule="self.secure || (!has(self.basicAuth) && !has(self.tls) && !has(self.serverTLS) && !has(self.jaas))",message="security settings require secure mode"
// +kubebuilder:validation:XValidation:rule="!self.secure || (has(self.serverTLS.enabled) && self.serverTLS.enabled && (!has(self.serverTLS.requireClientCert) || !self.serverTLS.requireClientCert))",message="secure mode requires one-way server TLS"
// +kubebuilder:validation:XValidation:rule="!self.secure || (has(self.tls.serverName) && size(self.tls.serverName) > 0)",message="secure mode requires the verified server name"
// +kubebuilder:validation:XValidation:rule="!self.secure || (has(self.serverTLS.keyStore) && size(self.serverTLS.keyStore) > 0 && has(self.serverTLS.keyStorePasswordFile) && size(self.serverTLS.keyStorePasswordFile) > 0 && has(self.serverTLS.keyPasswordFile) && size(self.serverTLS.keyPasswordFile) > 0)",message="secure mode requires a server keystore and password files"
type HealthMutatorConnection struct {
	// Secure enables authenticated PUT over verified TLS.
	// +kubebuilder:default=false
	Secure bool `json:"secure"`
	// BasicAuth supplies credentials for authenticated health requests.
	// +optional
	BasicAuth *HTTPBasicAuth `json:"basicAuth,omitempty"`
	// TLS configures CA trust and server name verification for health requests.
	// +optional
	TLS *HTTPClientTLS `json:"tls,omitempty"`
	// ServerTLS configures the member's one-way TLS identity. TrustStore fields
	// are accepted for compatibility but are not used by HealthSSLProvider.
	// +optional
	ServerTLS *SSLSpec `json:"serverTLS,omitempty"`
	// JAAS selects how the member's CoherenceREST login context is configured.
	// +optional
	JAAS *HealthJAAS `json:"jaas,omitempty"`
}

// HealthJAAS selects the packaged validator or an application-owned JAAS configuration.
type HealthJAAS struct {
	// Mode is "packaged" to generate the packaged JAAS configuration or "existing" to use the
	// application's configuration. It defaults to "packaged".
	// +optional
	// +kubebuilder:default=packaged
	// +kubebuilder:validation:Enum=packaged;existing
	Mode string `json:"mode,omitempty"`
}

// Validate checks the complete managed connection before use.
func (in *HealthMutatorConnection) Validate() error {
	if in == nil {
		return nil
	}
	if !in.Secure {
		if in.BasicAuth != nil || in.TLS != nil || in.ServerTLS != nil || in.JAAS != nil {
			return fmt.Errorf("security settings require secure mode")
		}
		return nil
	}
	if in.BasicAuth == nil || in.TLS == nil || in.TLS.ServerName == "" || in.ServerTLS == nil || in.ServerTLS.Enabled == nil || !*in.ServerTLS.Enabled || (in.ServerTLS.RequireClientCert != nil && *in.ServerTLS.RequireClientCert) {
		return fmt.Errorf("secure health mutators require credentials and one-way TLS with a verified server name")
	}
	a := &Probe{HTTP: &HTTPAction{HTTPEndpoint: HTTPEndpoint{Scheme: corev1.URISchemeHTTPS}, BasicAuth: in.BasicAuth, TLS: in.TLS}}
	if err := a.ValidateAction(); err != nil {
		return err
	}
	if in.JAAS != nil && in.JAAS.Mode != "" && in.JAAS.Mode != "packaged" && in.JAAS.Mode != "existing" {
		return fmt.Errorf("unsupported health JAAS mode")
	}
	for _, value := range []*string{in.ServerTLS.KeyStore, in.ServerTLS.KeyStorePasswordFile, in.ServerTLS.KeyPasswordFile} {
		if value == nil || *value == "" {
			return fmt.Errorf("server keystore and password files are required")
		}
	}
	return nil
}

// ConfigureHealthMutators supplies only references and read-only mounts to the member container.
func (in *CoherenceStatefulSetResourceSpec) ConfigureHealthMutators(pod *corev1.PodTemplateSpec) {
	connection := in.HealthMutatorConnection
	if connection == nil || !connection.Secure {
		return
	}
	data, _ := json.Marshal(connection)
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations[HealthMutatorAnnotation] = string(data)
	var containers []*corev1.Container
	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == ContainerNameCoherence {
			containers = append(containers, &pod.Spec.Containers[i])
		}
	}
	for i := range pod.Spec.InitContainers {
		if pod.Spec.InitContainers[i].Name == ContainerNameOperatorConfig {
			containers = append(containers, &pod.Spec.InitContainers[i])
		}
	}
	mounted := map[string]bool{}
	for _, container := range containers {
		container.Env = append(container.Env, corev1.EnvVar{Name: HealthMutatorEnv, Value: string(data)})
		mount := func(name, secret, path string, items []corev1.KeyToPath) {
			if !mounted[name] {
				pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{Name: name, VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: secret, Items: items}}})
				mounted[name] = true
			}
			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{Name: name, MountPath: path, ReadOnly: true})
		}
		if connection.BasicAuth != nil {
			mount("health-mutator-username", connection.BasicAuth.Username.Name, "/coherence/health-auth/username", []corev1.KeyToPath{{Key: connection.BasicAuth.Username.Key, Path: "value"}})
			mount("health-mutator-password", connection.BasicAuth.Password.Name, "/coherence/health-auth/password", []corev1.KeyToPath{{Key: connection.BasicAuth.Password.Key, Path: "value"}})
		}
		if connection.ServerTLS != nil && connection.ServerTLS.Secrets != nil {
			mount("health-mutator-tls", *connection.ServerTLS.Secrets, "/coherence/health-tls", nil)
		}
		probes := []struct {
			probe     *corev1.Probe
			generated bool
		}{
			{container.ReadinessProbe, isGeneratedHealthProbe(in.ReadinessProbe)},
			{container.LivenessProbe, isGeneratedHealthProbe(in.LivenessProbe)},
			{container.StartupProbe, in.StartupProbe != nil && isGeneratedHealthProbe(in.StartupProbe)},
		}
		for _, item := range probes {
			if item.generated && item.probe != nil && item.probe.HTTPGet != nil && item.probe.HTTPGet.Port.IntValue() == int(in.GetHealthPort()) {
				item.probe.HTTPGet.Scheme = corev1.URISchemeHTTPS
			}
		}
	}
}

func isGeneratedHealthProbe(probe *ReadinessProbeSpec) bool {
	return probe == nil || (probe.Exec == nil && probe.HTTPGet == nil && probe.TCPSocket == nil)
}
