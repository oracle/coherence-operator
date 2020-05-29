/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"github.com/oracle/coherence-operator/pkg/flags"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"strconv"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoherenceResourceSpec defines the specification of a Coherence resource. A Coherence resource is
//typically one or more Pods that perform the same functionality, for example storage members.
// +k8s:openapi-gen=true
type CoherenceResourceSpec struct {
	// The name of the image.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image *string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any
	// of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +listType=map
	// +listMapKey=name
	// +optional
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// The desired number of cluster members of this deployment.
	// This is a pointer to distinguish between explicit zero and not specified.
	// If not specified a default value of 3 will be used.
	// This field cannot be negative.
	// +kubebuilder:validation:Minimum:=0
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// The optional name of the Coherence cluster that this Coherence resource belongs to.
	// If this value is set the Pods controlled by this Coherence resource will form a cluster
	// with other Pods controlled by Coherence resources with the same cluster name.
	// If not set the Coherence resource's name will be used as the cluster name.
	// +optional
	Cluster *string `json:"cluster,omitempty"`
	// The name of the role that this deployment represents in a Coherence cluster.
	// This value will be used to set the Coherence role property for all members of this role
	// +optional
	Role string `json:"role,omitempty"`
	// The optional settings specific to Coherence functionality.
	// +optional
	Coherence *CoherenceSpec `json:"coherence,omitempty"`
	// The optional application specific settings.
	// +optional
	Application *ApplicationSpec `json:"application,omitempty"`
	// The JVM specific options
	// +optional
	JVM *JVMSpec `json:"jvm,omitempty"`
	// Ports specifies additional port mappings for the Pod and additional Services for those ports.
	// +listType=map
	// +listMapKey=name
	// +optional
	Ports []NamedPortSpec `json:"ports,omitempty"`
	// The configuration to control safe scaling.
	// +optional
	Scaling *ScalingSpec `json:"scaling,omitempty"`
	// StartQuorum controls the start-up order of this Coherence resource
	// in relation to other Coherence resources.
	// +listType=map
	// +listMapKey=deployment
	// +optional
	StartQuorum []StartQuorum `json:"startQuorum,omitempty"`
	// Env is additional environment variable mappings that will be passed to
	// the Coherence container in the Pod.
	// To specify extra variables add them as name value pairs the same as they
	// would be added to a Pod containers spec.
	// +listType=map
	// +listMapKey=name
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// The extra labels to add to the all of the Pods in this deployments.
	// Labels here will add to or override those defined for the cluster.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are free-form yaml that will be added to the store release as annotations
	// Any annotations should be placed BELOW this annotations: key. For example if we wanted to
	// include annotations for Prometheus it would look like this:
	//
	// annotations:
	//   prometheus.io/scrape: "true"
	//   prometheus.io/port: "2408"
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// List of additional initialization containers to add to the deployment's Pod.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	// +listType=map
	// +listMapKey=name
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
	// List of additional side-car containers to add to the deployment's Pod.
	// +listType=map
	// +listMapKey=name
	SideCars []corev1.Container `json:"sideCars,omitempty"`
	// A list of ConfigMaps to add as volumes.
	// Each entry in the list will be added as a ConfigMap Volume to the deployment's
	// Pods and as a VolumeMount to all of the containers and init-containers in the Pod.
	// +coh:doc=misc_pod_settings/050_configmap_volumes.adoc,Add ConfigMap Volumes
	// +listType=map
	// +listMapKey=name
	ConfigMapVolumes []ConfigMapVolumeSpec `json:"configMapVolumes,omitempty"`
	// A list of Secrets to add as volumes.
	// Each entry in the list will be added as a Secret Volume to the deployment's
	// Pods and as a VolumeMount to all of the containers and init-containers in the Pod.
	// +coh:doc=misc_pod_settings/020_secret_volumes.adoc,Add Secret Volumes
	// +listType=map
	// +listMapKey=name
	SecretVolumes []SecretVolumeSpec `json:"secretVolumes,omitempty"`
	// Volumes defines extra volume mappings that will be added to the Coherence Pod.
	//   The content of this yaml should match the normal k8s volumes section of a Pod definition
	//   as described in https://kubernetes.io/docs/concepts/storage/volumes/
	// +listType=map
	// +listMapKey=name
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// VolumeClaimTemplates defines extra PVC mappings that will be added to the Coherence Pod.
	//   The content of this yaml should match the normal k8s volumeClaimTemplates section of a Pod definition
	//   as described in https://kubernetes.io/docs/concepts/storage/persistent-volumes/
	// +listType=atomic
	// +optional
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	// VolumeMounts defines extra volume mounts to map to the additional volumes or PVCs declared above
	//   in store.volumes and store.volumeClaimTemplates
	// +listType=map
	// +listMapKey=name
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// The port that the health check endpoint will bind to.
	// +optional
	HealthPort *int32 `json:"healthPort,omitempty"`
	// The readiness probe config to be used for the Pods in this deployment.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	ReadinessProbe *ReadinessProbeSpec `json:"readinessProbe,omitempty"`
	// The liveness probe config to be used for the Pods in this deployment.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	LivenessProbe *ReadinessProbeSpec `json:"livenessProbe,omitempty"`
	// Resources is the optional resource requests and limits for the containers
	//  ref: http://kubernetes.io/docs/user-guide/compute-resources/
	//
	// By default the cpu requests is set to zero and the cpu limit set to 32. This
	// is because it appears that K8s defaults cpu to one and since Java 10 the JVM
	// now correctly picks up cgroup cpu limits then the JVM will only see one cpu.
	// By setting resources.requests.cpu=0 and resources.limits.cpu=32 it ensures that
	// the JVM will see the either the number of cpus on the host if this is <= 32 or
	// the JVM will see 32 cpus if the host has > 32 cpus. The limit is set to zero
	// so that there is no hard-limit applied.
	//
	// No default memory limits are applied.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Affinity controls Pod scheduling preferences.
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// NodeSelector is the Node labels for pod assignment
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Tolerations is for nodes that have taints on them.
	//   Useful if you want to dedicate nodes to just run the coherence container
	// For example:
	//   tolerations:
	//   - key: "key"
	//     operator: "Equal"
	//     value: "value"
	//     effect: "NoSchedule"
	//
	//   ref: https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
	// +listType=map
	// +listMapKey=key
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// SecurityContext is the PodSecurityContext that will be added to all of the Pods in this deployment.
	// See: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	// Share a single process namespace between all of the containers in a pod. When this is set containers will
	// be able to view and signal processes from other containers in the same pod, and the first process in each
	// container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set.
	// Optional: Default to false.
	// +optional
	ShareProcessNamespace *bool `json:"shareProcessNamespace,omitempty"`
	// Use the host's ipc namespace. Optional: Default to false.
	// +optional
	HostIPC *bool `json:"hostIPC,omitempty"`
	// Configure various networks and DNS settings for Pods in this rolw.
	// +optional
	Network *NetworkSpec `json:"network,omitempty"`
	// The configuration for the Coherence utils image
	// +optional
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// Whether or not to auto-mount the Kubernetes API credentials for a service account
	// +optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty"`
	// The timeout to apply to REST requests made back to the Operator from Coherence Pods.
	// These requests are typically to obtain site and rack information for the Pod.
	// +optional
	OperatorRequestTimeout *int32 `json:"operatorRequestTimeout,omitempty"`
}

// Obtain the number of replicas required for a deployment.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replica value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *CoherenceResourceSpec) GetReplicas() int32 {
	if in == nil {
		return 0
	}
	if in.Replicas == nil {
		return DefaultReplicas
	}
	return *in.Replicas
}

// Set the number of replicas required for a deployment.
func (in *CoherenceResourceSpec) SetReplicas(replicas int32) {
	if in != nil {
		in.Replicas = &replicas
	}
}

func (in *CoherenceResourceSpec) GetCoherenceImage() *string {
	if in != nil {
		return in.Image
	}
	return nil
}

// Ensure that the Coherence image is set for the deployment.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (in *CoherenceResourceSpec) EnsureCoherenceImage(coherenceImage *string) bool {
	if in.Image == nil {
		in.Image = coherenceImage
		return true
	}
	return false
}

func (in *CoherenceResourceSpec) GetCoherenceUtilsImage() *string {
	if in != nil && in.CoherenceUtils != nil {
		return in.CoherenceUtils.Image
	}
	return nil
}

// Ensure that the Coherence Utils image is set for the deployment.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (in *CoherenceResourceSpec) EnsureCoherenceUtilsImage(utilsImage *string) bool {
	if in.CoherenceUtils == nil {
		in.CoherenceUtils = &ImageSpec{}
	}

	return in.CoherenceUtils.EnsureImage(utilsImage)
}

func (in *CoherenceResourceSpec) GetEffectiveScalingPolicy() ScalingPolicy {
	if in == nil {
		return SafeScaling
	}

	var policy ScalingPolicy

	if in.Scaling == nil || in.Scaling.Policy == nil {
		// the scaling policy is not set the look at the storage enabled flag
		if in.Coherence == nil || in.Coherence.StorageEnabled == nil || *in.Coherence.StorageEnabled {
			// storage enabled is either not set or is true so do safe scaling
			policy = ParallelUpSafeDownScaling
		} else {
			// storage enabled is false so do parallel scaling
			policy = ParallelScaling
		}
	} else {
		// scaling policy is set so use it
		policy = *in.Scaling.Policy
	}

	return policy
}

// Returns the port that the health check endpoint will bind to.
func (in *CoherenceResourceSpec) GetHealthPort() int32 {
	if in == nil || in.HealthPort == nil || *in.HealthPort <= 0 {
		return DefaultHealthPort
	}
	return *in.HealthPort
}

// Returns the ScalingProbe to use for checking Phase HA for the deployment.
// This method will not return nil.
func (in *CoherenceResourceSpec) GetScalingProbe() *ScalingProbe {
	if in == nil || in.Scaling == nil || in.Scaling.Probe == nil {
		return in.GetDefaultScalingProbe()
	}
	return in.Scaling.Probe
}

// Obtain a default ScalingProbe
func (in *CoherenceResourceSpec) GetDefaultScalingProbe() *ScalingProbe {
	timeout := 10

	defaultStatusHA := ScalingProbe{
		TimeoutSeconds: &timeout,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/ha",
				Port: intstr.FromString(PortNameHealth),
			},
		},
	}

	return defaultStatusHA.DeepCopy()
}

// Create the Kubernetes resources that should be deployed for this deployment.
// The order of the resources in the returned array is the order that they should be
// created or updated in Kubernetes.
func (in *CoherenceResourceSpec) CreateKubernetesResources(d *Coherence, flags *flags.CoherenceOperatorFlags) (Resources, error) {
	var res []Resource

	if in.GetReplicas() <= 0 {
		// replicas is zero so nothing to create
		return Resources{Items: res}, nil
	}

	// Create the headless WKA Service if this deployment is a WKA member
	if in.Coherence.RequiresWKAService() {
		res = append(res, in.CreateWKAService(d))
	}

	// Create the headless Service
	res = append(res, in.CreateHeadlessService(d))

	// Create the StatefulSet
	res = append(res, in.CreateStatefulSet(d, flags))

	// Create the Services for each port (and optionally ServiceMonitors)
	res = append(res, in.CreateServicesForPort(d)...)

	return Resources{Items: res}, nil
}

// Create the Services for each port (and optionally ServiceMonitors)
func (in *CoherenceResourceSpec) CreateServicesForPort(deployment *Coherence) []Resource {
	var resources []Resource

	if in == nil || in.Ports == nil || len(in.Ports) == 0 {
		return resources
	}

	// Create the Service and ServiceMonitor for each port
	for _, p := range in.Ports {
		service := p.CreateService(deployment)
		if service != nil {
			resources = append(resources, Resource{
				Kind: ResourceTypeService,
				Name: service.GetName(),
				Spec: service,
			})
		}
		sm := p.CreateServiceMonitor(deployment)
		if sm != nil {
			resources = append(resources, Resource{
				Kind: ResourceTypeServiceMonitor,
				Name: sm.GetName(),
				Spec: sm,
			})
		}
	}

	return resources
}

// Create the selector that can be used to match this deployments Pods, for example by Services or StatefulSets.
func (in *CoherenceResourceSpec) CreatePodSelectorLabels(deployment *Coherence) map[string]string {
	selector := deployment.CreateCommonLabels()
	selector[LabelComponent] = LabelComponentCoherencePod
	return selector
}

// Create the headless WKA Service
func (in *CoherenceResourceSpec) CreateWKAService(deployment *Coherence) Resource {
	labels := deployment.CreateCommonLabels()
	labels[LabelComponent] = LabelComponentWKA

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[LabelCoherenceCluster] = deployment.GetCoherenceClusterName()
	selector[LabelComponent] = LabelComponentCoherencePod
	selector[LabelCoherenceWKAMember] = "true"

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deployment.Namespace,
			Name:      deployment.Name + WKAServiceNameSuffix,
			Labels:    labels,
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			// Pods must be part of the WKA service even if not ready
			PublishNotReadyAddresses: true,
			Ports: []corev1.ServicePort{
				{
					Name:       PortNameCoherence,
					Protocol:   corev1.ProtocolTCP,
					Port:       7,
					TargetPort: intstr.FromInt(7),
				},
			},
			Selector: selector,
		},
	}

	return Resource{
		Kind: ResourceTypeService,
		Name: svc.GetName(),
		Spec: svc,
	}
}

// Create the headless Service for the deployment's StatefulSet.
func (in *CoherenceResourceSpec) CreateHeadlessService(deployment *Coherence) Resource {
	// The labels for the service
	svcLabels := deployment.CreateCommonLabels()
	svcLabels[LabelComponent] = LabelComponentCoherenceHeadless

	// The selector for the service
	selector := in.CreatePodSelectorLabels(deployment)

	// Create the Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deployment.GetNamespace(),
			Name:      deployment.GetName(),
			Labels:    svcLabels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP:                "None",
			PublishNotReadyAddresses: true,
			Selector:                 selector,
			Ports: []corev1.ServicePort{
				{
					Name:       PortNameCoherence,
					Protocol:   corev1.ProtocolTCP,
					Port:       7,
					TargetPort: intstr.FromInt(7),
				},
			},
		},
	}

	return Resource{
		Kind: ResourceTypeService,
		Name: svc.GetName(),
		Spec: svc,
	}
}

// Create the deployment's StatefulSet.
func (in *CoherenceResourceSpec) CreateStatefulSet(deployment *Coherence, flags *flags.CoherenceOperatorFlags) Resource {
	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deployment.GetNamespace(),
			Name:      deployment.GetName(),
			Labels:    deployment.CreateCommonLabels(),
		},
	}

	// Create the PodSpec labels
	podLabels := in.CreatePodSelectorLabels(deployment)
	// Add the WKA member label
	podLabels[LabelCoherenceWKAMember] = strconv.FormatBool(in.Coherence.IsWKAMember())
	// Add any labels specified for the deployment
	for k, v := range in.Labels {
		podLabels[k] = v
	}

	replicas := in.GetReplicas()

	cohContainer := in.CreateCoherenceContainer(deployment, flags)

	// Add additional ports
	for _, p := range in.Ports {
		cohContainer.Ports = append(cohContainer.Ports, p.CreatePort(deployment))
	}

	// append any additional VolumeMounts
	cohContainer.VolumeMounts = append(cohContainer.VolumeMounts, in.VolumeMounts...)

	// Add the component label
	sts.Labels[LabelComponent] = LabelComponentCoherenceStatefulSet
	sts.Spec = appsv1.StatefulSetSpec{
		Replicas:            &replicas,
		PodManagementPolicy: appsv1.ParallelPodManagement,
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
		RevisionHistoryLimit: pointer.Int32Ptr(5),
		ServiceName:          deployment.GetName(),
		Selector: &metav1.LabelSelector{
			MatchLabels: in.CreatePodSelectorLabels(deployment),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      podLabels,
				Annotations: in.Annotations,
			},
			Spec: corev1.PodSpec{
				ImagePullSecrets:             in.GetImagePullSecrets(),
				ServiceAccountName:           in.GetServiceAccountName(),
				AutomountServiceAccountToken: in.AutomountServiceAccountToken,
				SecurityContext:              in.SecurityContext,
				ShareProcessNamespace:        in.ShareProcessNamespace,
				HostIPC:                      notNilBool(in.HostIPC),
				Tolerations:                  in.Tolerations,
				Affinity:                     in.EnsurePodAffinity(deployment),
				NodeSelector:                 in.NodeSelector,
				InitContainers: []corev1.Container{
					in.CreateUtilsContainer(deployment, flags),
				},
				Containers: []corev1.Container{cohContainer},
				Volumes: []corev1.Volume{
					{Name: VolumeNameLogs, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
					{Name: VolumeNameUtils, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				},
			},
		},
	}

	// Add any network settings
	in.Network.UpdateStatefulSet(&sts)
	// Add any JVM settings
	in.JVM.UpdateStatefulSet(&sts)
	// Add any Coherence settings
	in.Coherence.UpdateStatefulSet(deployment, &sts)

	// Add any additional init-containers and any additional containers
	in.ProcesssideCars(deployment, &sts)

	// Add any ConfigMap Volumes
	for _, cmv := range in.ConfigMapVolumes {
		cmv.AddVolumes(&sts)
	}

	// Add any Secret Volumes
	for _, sv := range in.SecretVolumes {
		sv.AddVolumes(&sts)
	}

	// append any additional Volumes
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, in.Volumes...)
	// append any additional PVCs
	sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, in.VolumeClaimTemplates...)

	return Resource{
		Kind: ResourceTypeStatefulSet,
		Name: sts.GetName(),
		Spec: &sts,
	}
}

func (in *CoherenceResourceSpec) GetImagePullSecrets() []corev1.LocalObjectReference {
	var secrets []corev1.LocalObjectReference

	for _, s := range in.ImagePullSecrets {
		secrets = append(secrets, corev1.LocalObjectReference{
			Name: s.Name,
		})
	}

	return secrets
}

// Get the service account name for the cluster.
func (in *CoherenceResourceSpec) GetServiceAccountName() string {
	if in != nil && in.ServiceAccountName != DefaultServiceAccount {
		return in.ServiceAccountName
	}
	return ""
}

// Create the Coherence container spec.
func (in *CoherenceResourceSpec) CreateCoherenceContainer(deployment *Coherence, flags *flags.CoherenceOperatorFlags) corev1.Container {
	var cohImage *string

	if in.Image == nil {
		cohImage = flags.GetCoherenceImage()
	} else {
		cohImage = in.Image
	}

	healthPort := in.GetHealthPort()
	vm := in.CreateCommonVolumeMounts()

	c := corev1.Container{
		Name:    ContainerNameCoherence,
		Image:   *cohImage,
		Command: []string{"/utils/runner", "server"},
		Env:     in.Env,
		Ports: []corev1.ContainerPort{
			{
				Name:          PortNameCoherence,
				ContainerPort: 7,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          PortNameHealth,
				ContainerPort: healthPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts: vm,
	}

	if in.ImagePullPolicy != nil {
		c.ImagePullPolicy = *in.ImagePullPolicy
	}

	c.Env = append(c.Env, in.CreateDefaultEnv(deployment)...)

	in.Application.UpdateCoherenceContainer(&c)

	if in.Resources != nil {
		// set the container resources if specified
		c.Resources = *in.Resources
	} else {
		// No resources specified so default to 32 cores
		c.Resources = in.CreateDefaultResources()
	}

	c.ReadinessProbe = in.CreateDefaultReadinessProbe()
	in.ReadinessProbe.UpdateProbeSpec(healthPort, DefaultReadinessPath, c.ReadinessProbe)

	c.LivenessProbe = in.CreateDefaultLivenessProbe()
	in.LivenessProbe.UpdateProbeSpec(healthPort, DefaultLivenessPath, c.LivenessProbe)

	return c
}

// Create the common VolumeMounts added all containers.
func (in *CoherenceResourceSpec) CreateCommonVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{Name: VolumeNameLogs, MountPath: VolumeMountPathLogs},
		{Name: VolumeNameUtils, MountPath: VolumeMountPathUtils},
		{Name: VolumeNameJVM, MountPath: VolumeMountPathJVM},
	}
}

// Create the common environment variables added all.
func (in *CoherenceResourceSpec) CreateCommonEnv(deployment *Coherence) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: EnvVarCohMachineName, ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		},
		{
			Name: EnvVarCohMemberName, ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: EnvVarCohPodUID, ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.uid",
				},
			},
		},
		{Name: EnvVarCohClusterName, Value: deployment.GetCoherenceClusterName()},
		{Name: EnvVarCohRole, Value: deployment.GetRoleName()},
	}
}

// Create the default environment variables for the Coherence container.
func (in *CoherenceResourceSpec) CreateDefaultEnv(deployment *Coherence) []corev1.EnvVar {
	return append(in.CreateCommonEnv(deployment),
		corev1.EnvVar{Name: EnvVarCohWka, Value: deployment.Spec.Coherence.GetWKA(deployment.Name)},
		corev1.EnvVar{
			Name: EnvVarOperatorHost, ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: OperatorConfigName},
					Key:                  OperatorConfigKeyHost,
					Optional:             pointer.BoolPtr(true),
				},
			},
		},
		corev1.EnvVar{Name: EnvVarCohSite, Value: "http://$(OPERATOR_HOST)/site/$(COH_MACHINE_NAME)"},
		corev1.EnvVar{Name: EnvVarCohRack, Value: "http://$(OPERATOR_HOST)/rack/$(COH_MACHINE_NAME)"},
		corev1.EnvVar{Name: EnvVarCohUtilDir, Value: VolumeMountPathUtils},
		corev1.EnvVar{Name: EnvVarOperatorTimeout, Value: Int32PtrToStringWithDefault(in.OperatorRequestTimeout, 120)},
		corev1.EnvVar{Name: EnvVarCohHealthPort, Value: Int32ToString(in.GetHealthPort())},
	)
}

// Create the default Container resources.
func (in *CoherenceResourceSpec) CreateDefaultResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("32"),
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("0"),
		},
	}
}

// Create the default readiness probe.
func (in *CoherenceResourceSpec) CreateDefaultReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 30,
		PeriodSeconds:       60,
		TimeoutSeconds:      30,
		SuccessThreshold:    1,
		FailureThreshold:    50,
	}
}

// Update the probe with the default readiness probe action.
func (in *CoherenceResourceSpec) UpdateDefaultReadinessProbeAction(probe *corev1.Probe) *corev1.Probe {
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   DefaultReadinessPath,
		Port:   intstr.FromInt(int(DefaultHealthPort)),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}

// Create the default liveness probe.
func (in *CoherenceResourceSpec) CreateDefaultLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 60,
		PeriodSeconds:       60,
		TimeoutSeconds:      30,
		SuccessThreshold:    1,
		FailureThreshold:    5,
	}
}

// Update the probe with the default liveness probe action.
func (in *CoherenceResourceSpec) UpdateDefaultLivenessProbeAction(probe *corev1.Probe) *corev1.Probe {
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   DefaultLivenessPath,
		Port:   intstr.FromInt(int(DefaultHealthPort)),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}

// Get the Utils init-container spec.
func (in *CoherenceResourceSpec) CreateUtilsContainer(deployment *Coherence, flags *flags.CoherenceOperatorFlags) corev1.Container {
	var utilsImage *string
	if in.CoherenceUtils == nil || in.CoherenceUtils.Image == nil {
		utilsImage = flags.GetCoherenceUtilsImage()
	} else {
		utilsImage = in.CoherenceUtils.Image
	}

	vm := in.CreateCommonVolumeMounts()

	c := corev1.Container{
		Name:    ContainerNameUtils,
		Image:   *utilsImage,
		Command: []string{UtilsInitCommand},
		Env: []corev1.EnvVar{
			{Name: "COH_UTIL_DIR", Value: VolumeMountPathUtils},
			{Name: "COH_CLUSTER_NAME", Value: deployment.GetCoherenceClusterName()},
		},
		VolumeMounts: vm,
	}

	// set the image pull policy if set for the deployment
	if in.CoherenceUtils != nil && in.CoherenceUtils.ImagePullPolicy != nil {
		c.ImagePullPolicy = *in.CoherenceUtils.ImagePullPolicy
	}

	// set the persistence volume mounts if required
	in.Coherence.AddPersistenceVolumeMounts(&c)

	return c
}

// Get the Pod Affinity either from that configured for the cluster or the default affinity.
func (in *CoherenceResourceSpec) EnsurePodAffinity(deployment *Coherence) *corev1.Affinity {
	if in != nil && in.Affinity != nil {
		return in.Affinity
	}
	// return the default affinity which attempts to spread the Pods for a deployment across fault domains
	return in.CreateDefaultPodAffinity(deployment)
}

// Create the default Pod Affinity to use in a deployment's StatefulSet.
func (in *CoherenceResourceSpec) CreateDefaultPodAffinity(deployment *Coherence) *corev1.Affinity {
	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 1,
					PodAffinityTerm: corev1.PodAffinityTerm{
						TopologyKey: AffinityTopologyKey,
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      LabelCoherenceCluster,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{deployment.GetCoherenceClusterName()},
								},
								{
									Key:      LabelCoherenceDeployment,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{deployment.Name},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (in *CoherenceResourceSpec) GetMetricsPort() int32 {
	if in == nil {
		return 0
	}
	return in.Coherence.GetMetricsPort()
}

func (in *CoherenceResourceSpec) GetManagementPort() int32 {
	if in == nil {
		return 0
	}
	return in.Coherence.GetManagementPort()
}

// Add any additional init-containers or additional containers to the StatefulSet.
// This will add any common environment variables to te container too, unless those variable names
// have already been specified in the container spec
func (in *CoherenceResourceSpec) ProcesssideCars(deployment *Coherence, sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}

	for _, c := range in.InitContainers {
		in.processAdditionalContainer(deployment, &c)
		sts.Spec.Template.Spec.InitContainers = append(sts.Spec.Template.Spec.InitContainers, c)
	}

	for _, c := range in.SideCars {
		in.processAdditionalContainer(deployment, &c)
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, c)
	}
}

func (in *CoherenceResourceSpec) processAdditionalContainer(deployment *Coherence, c *corev1.Container) {
	in.appendCommonEnvVars(deployment, c)
	in.appendCommonVolumeMounts(c)
}

func (in *CoherenceResourceSpec) appendCommonEnvVars(deployment *Coherence, c *corev1.Container) {
	envVars := c.Env
	for _, toAdd := range in.CreateCommonEnv(deployment) {
		envVars = in.appendEnvVarIfMissing(envVars, toAdd)
	}
	c.Env = envVars
}

func (in *CoherenceResourceSpec) appendEnvVarIfMissing(envVars []corev1.EnvVar, toAdd corev1.EnvVar) []corev1.EnvVar {
	for _, ev := range envVars {
		if ev.Name == toAdd.Name {
			return envVars
		}
	}
	return append(envVars, toAdd)
}

func (in *CoherenceResourceSpec) appendCommonVolumeMounts(c *corev1.Container) {
	mounts := c.VolumeMounts
	for _, toAdd := range in.CreateCommonVolumeMounts() {
		mounts = in.appendVolumeMountIfMissing(mounts, toAdd)
	}
	c.VolumeMounts = mounts
}

func (in *CoherenceResourceSpec) appendVolumeMountIfMissing(mounts []corev1.VolumeMount, toAdd corev1.VolumeMount) []corev1.VolumeMount {
	for _, m := range mounts {
		if m.Name == toAdd.Name {
			return mounts
		}
	}
	return append(mounts, toAdd)
}
