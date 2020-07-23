/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package legacy

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/pointer"
	"os"
	y2 "sigs.k8s.io/yaml"
	"strings"
)

func Convert(fileName string, out io.Writer) error {
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrap(err, "Failed to read file "+fileName)
	}

	cc := &CoherenceCluster{}
	decoder := yaml.NewYAMLToJSONDecoder(strings.NewReader(string(data)))
	err = decoder.Decode(cc)
	if err != nil {
		return errors.Wrap(err, "Failed to parse CoherenceCluster from file "+fileName)
	}

	roles := cc.GetRoles()
	sep := false
	for i := range roles {
		role := roles[i]
		if sep {
			_, err = fmt.Fprintln(out, "---")
			if err != nil {
				return err
			}
		} else {
			sep = true
		}

		c, err := convertRole(&role, cc)
		if err != nil {
			return err
		}
		bytes, err := y2.Marshal(c)
		if err != nil {
			return err
		}

		fmt.Print(string(bytes))
	}

	return nil
}

func convertRole(r *CoherenceRoleSpec, cc *CoherenceCluster) (*v1.Coherence, error) {
	clusterName := cc.Name
	c := &v1.Coherence{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Coherence",
			APIVersion: "coherence.oracle.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        r.GetFullRoleName(cc),
			Namespace:   cc.Namespace,
			Labels:      cc.Labels,
			Annotations: cc.Annotations,
		},
		Spec: v1.CoherenceResourceSpec{
			Image:                        doImage(r),
			ImagePullPolicy:              doImagePullPolicy(r),
			ImagePullSecrets:             doPullSecrets(cc),
			Replicas:                     r.Replicas,
			Cluster:                      &clusterName,
			Role:                         r.GetRoleName(),
			Coherence:                    doCoherence(r),
			Application:                  doApplication(r),
			JVM:                          doJVM(r),
			Ports:                        doPorts(r),
			Scaling:                      doScaling(r),
			StartQuorum:                  doStartQuorum(r, cc),
			Env:                          r.Env,
			Labels:                       r.Labels,
			Annotations:                  r.Annotations,
			Volumes:                      r.Volumes,
			VolumeClaimTemplates:         r.VolumeClaimTemplates,
			VolumeMounts:                 r.VolumeMounts,
			HealthPort:                   r.HealthPort,
			ReadinessProbe:               doReadiness(r),
			LivenessProbe:                doLiveness(r),
			Resources:                    r.Resources,
			Affinity:                     r.Affinity,
			NodeSelector:                 r.NodeSelector,
			Tolerations:                  r.Tolerations,
			SecurityContext:              r.SecurityContext,
			ShareProcessNamespace:        r.ShareProcessNamespace,
			HostIPC:                      r.HostIPC,
			Network:                      doNetwork(r),
			ServiceAccountName:           cc.Spec.ServiceAccountName,
			AutomountServiceAccountToken: cc.Spec.AutomountServiceAccountToken,
			OperatorRequestTimeout:       cc.Spec.OperatorRequestTimeout,
			InitContainers:               nil,
			SideCars:                     nil,
			ConfigMapVolumes:             nil,
			SecretVolumes:                nil,
			CoherenceUtils:               nil,
		},
	}

	return c, nil
}

func doCoherence(r *CoherenceRoleSpec) *v1.CoherenceSpec {
	if r.Coherence == nil {
		return nil
	}

	return &v1.CoherenceSpec{
		CacheConfig:    r.Coherence.CacheConfig,
		OverrideConfig: r.Coherence.OverrideConfig,
		StorageEnabled: r.Coherence.StorageEnabled,
		LogLevel:       r.Coherence.LogLevel,
		ExcludeFromWKA: r.Coherence.ExcludeFromWKA,
		Persistence:    doPersistence(r.Coherence.Persistence, r.Coherence.Snapshot),
		Management:     doPortSpecWithSSL(r.Coherence.Management),
		Metrics:        doPortSpecWithSSL(r.Coherence.Metrics),
	}
}

func doPersistence(p *PersistentStorageSpec, s *PersistentStorageSpec) *v1.PersistenceSpec {
	if p == nil && s == nil {
		return nil
	}

	spec := &v1.PersistenceSpec{}

	if p != nil {
		if p.Enabled != nil && *p.Enabled {
			spec.Mode = pointer.StringPtr("active")
		}
		spec.PersistentVolumeClaim = p.PersistentVolumeClaim
		spec.Volume = p.Volume
	}

	if s != nil && s.Enabled != nil && *s.Enabled {
		spec.Snapshots = &v1.PersistentStorageSpec{
			PersistentVolumeClaim: s.PersistentVolumeClaim,
			Volume:                s.Volume,
		}
	}

	return spec
}

func doPortSpecWithSSL(p *PortSpecWithSSL) *v1.PortSpecWithSSL {
	if p == nil {
		return nil
	}

	spec := &v1.PortSpecWithSSL{
		Enabled: p.Enabled,
		Port:    p.Port,
	}

	if p.SSL != nil {
		spec.SSL = &v1.SSLSpec{
			Enabled:                p.SSL.Enabled,
			Secrets:                p.SSL.Secrets,
			KeyStore:               p.SSL.KeyStore,
			KeyStorePasswordFile:   p.SSL.KeyStorePasswordFile,
			KeyPasswordFile:        p.SSL.KeyPasswordFile,
			KeyStoreAlgorithm:      p.SSL.KeyStoreAlgorithm,
			KeyStoreProvider:       p.SSL.KeyStoreProvider,
			KeyStoreType:           p.SSL.KeyStoreType,
			TrustStore:             p.SSL.TrustStore,
			TrustStorePasswordFile: p.SSL.TrustStorePasswordFile,
			TrustStoreAlgorithm:    p.SSL.TrustStoreAlgorithm,
			TrustStoreProvider:     p.SSL.TrustStoreProvider,
			TrustStoreType:         p.SSL.TrustStoreType,
			RequireClientCert:      p.SSL.RequireClientCert,
		}
	}

	return spec
}

func doJVM(r *CoherenceRoleSpec) *v1.JVMSpec {
	if r.JVM == nil {
		return nil
	}

	jvm := &v1.JVMSpec{
		Args:               r.JVM.Args,
		UseContainerLimits: r.JVM.UseContainerLimits,
		DiagnosticsVolume:  r.JVM.DiagnosticsVolume,
		Jmxmp:              nil,
	}

	if r.JVM.Memory != nil {
		jvm.Memory = &v1.JvmMemorySpec{
			HeapSize:             r.JVM.Memory.HeapSize,
			StackSize:            r.JVM.Memory.StackSize,
			MetaspaceSize:        r.JVM.Memory.MetaspaceSize,
			DirectMemorySize:     r.JVM.Memory.DirectMemorySize,
			NativeMemoryTracking: r.JVM.Memory.NativeMemoryTracking,
		}

		if r.JVM.Memory.OnOutOfMemory != nil {
			jvm.Memory.OnOutOfMemory = &v1.JvmOutOfMemorySpec{
				Exit:     r.JVM.Memory.OnOutOfMemory.Exit,
				HeapDump: r.JVM.Memory.OnOutOfMemory.HeapDump,
			}
		}
	}

	if r.JVM.Jmxmp != nil {
		jvm.Jmxmp = &v1.JvmJmxmpSpec{
			Enabled: r.JVM.Jmxmp.Enabled,
			Port:    r.JVM.Jmxmp.Port,
		}
	}

	if r.JVM.Gc != nil {
		jvm.Gc = &v1.JvmGarbageCollectorSpec{
			Collector: r.JVM.Gc.Collector,
			Args:      r.JVM.Gc.Args,
			Logging:   r.JVM.Gc.Logging,
		}
	}

	if r.JVM.Debug != nil {
		jvm.Debug = &v1.JvmDebugSpec{
			Enabled: r.JVM.Debug.Enabled,
			Suspend: r.JVM.Debug.Suspend,
			Attach:  r.JVM.Debug.Attach,
			Port:    r.JVM.Debug.Port,
		}
	}

	return jvm
}

func doPorts(r *CoherenceRoleSpec) []v1.NamedPortSpec {
	var ports []v1.NamedPortSpec
	for _, port := range r.Ports {
		nps := v1.NamedPortSpec{
			Name: port.Name,
			Port: port.Port,
		}

		if port.Protocol != nil {
			p := corev1.Protocol(*port.Protocol)
			nps.Protocol = &p
		}

		if port.Service != nil {
			nps.Service = &v1.ServiceSpec{
				Enabled:                  port.Service.Enabled,
				Name:                     port.Service.Name,
				Port:                     port.Service.Port,
				Type:                     port.Service.Type,
				LoadBalancerIP:           port.Service.LoadBalancerIP,
				Labels:                   port.Service.Labels,
				Annotations:              port.Service.Annotations,
				SessionAffinity:          port.Service.SessionAffinity,
				LoadBalancerSourceRanges: port.Service.LoadBalancerSourceRanges,
				ExternalName:             port.Service.ExternalName,
				ExternalTrafficPolicy:    port.Service.ExternalTrafficPolicy,
				HealthCheckNodePort:      port.Service.HealthCheckNodePort,
				PublishNotReadyAddresses: port.Service.PublishNotReadyAddresses,
				SessionAffinityConfig:    port.Service.SessionAffinityConfig,
			}
		}

		ports = append(ports, nps)
	}

	return ports
}

func doImage(r *CoherenceRoleSpec) *string {
	switch {
	case r.Application != nil && r.Application.Image != nil:
		return r.Application.Image
	case r.Coherence != nil && r.Coherence.Image != nil:
		return r.Coherence.Image
	default:
		return nil
	}
}

func doImagePullPolicy(r *CoherenceRoleSpec) *corev1.PullPolicy {
	switch {
	case r.Application != nil && r.Application.ImagePullPolicy != nil:
		return r.Application.ImagePullPolicy
	case r.Coherence != nil && r.Coherence.ImagePullPolicy != nil:
		return r.Coherence.ImagePullPolicy
	default:
		return nil
	}
}

func doApplication(r *CoherenceRoleSpec) *v1.ApplicationSpec {
	if r.Application == nil {
		return nil
	}

	if r.Application.Main != nil || len(r.Application.Args) > 0 {
		return &v1.ApplicationSpec{
			Main: r.Application.Main,
			Args: r.Application.Args,
		}
	}
	return nil
}

func doScaling(r *CoherenceRoleSpec) *v1.ScalingSpec {
	if r.Scaling == nil {
		return nil
	}

	scaling := &v1.ScalingSpec{}
	if r.Scaling.Policy != nil {
		policy := v1.ScalingPolicy(*r.Scaling.Policy)
		scaling.Policy = &policy
	}

	probe := &v1.ScalingProbe{}
	if r.Scaling.Probe != nil {
		probe.TimeoutSeconds = r.Scaling.Probe.TimeoutSeconds
		probe.Exec = r.Scaling.Probe.Exec
		probe.HTTPGet = r.Scaling.Probe.HTTPGet
		probe.TCPSocket = r.Scaling.Probe.TCPSocket
	}

	return scaling
}

func doPullSecrets(cc *CoherenceCluster) []v1.LocalObjectReference {
	var pullSecrets []v1.LocalObjectReference
	for _, ps := range cc.Spec.ImagePullSecrets {
		pullSecrets = append(pullSecrets, v1.LocalObjectReference{
			Name: ps.Name,
		})
	}
	return pullSecrets
}

func doStartQuorum(r *CoherenceRoleSpec, cc *CoherenceCluster) []v1.StartQuorum {
	if r.StartQuorum == nil {
		return nil
	}

	var quorum []v1.StartQuorum
	for _, sq := range r.StartQuorum {
		quorum = append(quorum, v1.StartQuorum{
			Deployment: cc.GetFullRoleName(sq.Role),
			PodCount:   sq.PodCount,
		})
	}

	return quorum
}

func doReadiness(r *CoherenceRoleSpec) *v1.ReadinessProbeSpec {
	if r.ReadinessProbe == nil {
		return nil
	}

	rp := r.ReadinessProbe

	return &v1.ReadinessProbeSpec{
		ProbeHandler: v1.ProbeHandler{
			Exec:      rp.Exec,
			HTTPGet:   rp.HTTPGet,
			TCPSocket: rp.TCPSocket,
		},
		InitialDelaySeconds: rp.InitialDelaySeconds,
		TimeoutSeconds:      rp.TimeoutSeconds,
		PeriodSeconds:       rp.PeriodSeconds,
		SuccessThreshold:    rp.SuccessThreshold,
		FailureThreshold:    rp.FailureThreshold,
	}
}

func doLiveness(r *CoherenceRoleSpec) *v1.ReadinessProbeSpec {
	if r.LivenessProbe == nil {
		return nil
	}

	lp := r.LivenessProbe

	return &v1.ReadinessProbeSpec{
		ProbeHandler: v1.ProbeHandler{
			Exec:      lp.Exec,
			HTTPGet:   lp.HTTPGet,
			TCPSocket: lp.TCPSocket,
		},
		InitialDelaySeconds: lp.InitialDelaySeconds,
		TimeoutSeconds:      lp.TimeoutSeconds,
		PeriodSeconds:       lp.PeriodSeconds,
		SuccessThreshold:    lp.SuccessThreshold,
		FailureThreshold:    lp.FailureThreshold,
	}
}

func doNetwork(r *CoherenceRoleSpec) *v1.NetworkSpec {
	if r.Network == nil {
		return nil
	}

	n := r.Network

	spec := &v1.NetworkSpec{
		HostAliases: n.HostAliases,
		HostNetwork: n.HostNetwork,
		Hostname:    n.Hostname,
	}

	if n.DNSConfig != nil {
		spec.DNSConfig = &v1.PodDNSConfig{
			Nameservers: n.DNSConfig.Nameservers,
			Searches:    n.DNSConfig.Searches,
			Options:     n.DNSConfig.Options,
		}
	}

	if n.DNSPolicy != nil {
		policy := corev1.DNSPolicy(*n.DNSPolicy)
		spec.DNSPolicy = &policy
	}

	return spec
}
