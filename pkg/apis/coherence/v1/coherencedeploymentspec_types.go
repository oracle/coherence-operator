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

// CoherenceDeploymentSpec defines a deployment in a Coherence cluster. A deployment is one or
// more Pods that perform the same functionality, for example storage members.
// +k8s:openapi-gen=true
type CoherenceDeploymentSpec struct {
	// The image to run.
	ImageSpec `json:",inline"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any
	// of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +listType=map
	// +listMapKey=name
	// +optional
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// Whether or not to auto-mount the Kubernetes API credentials for a service account
	// +optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty"`
	// The timeout to apply to rest requests made back to the operator from Coherence Pods.
	// +optional
	OperatorRequestTimeout *int32 `json:"operatorRequestTimeout,omitempty"`
	// The optional name of the Coherence cluster that this CoherenceDeployment belongs to.
	// If this value is set this deployment will form a cluster with other deployments with
	// the same cluster name. If not set the CoherenceDeployment's name will be used as the
	// cluster name.
	// +optional
	Cluster *string `json:"cluster,omitempty"`
	// The name of the role that this deployment represents in a Coherence cluster.
	// This value will be used to set the Coherence role property for all members of this role
	// +optional
	Role string `json:"role,omitempty"`
	// The desired number of cluster members of this deployment.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Default value is 3.
	// +kubebuilder:validation:Minimum:=0
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// The optional application definition
	// +optional
	Application *ApplicationSpec `json:"application,omitempty"`
	// The optional application definition
	// +optional
	Coherence *CoherenceSpec `json:"coherence,omitempty"`
	// The configuration for the Coherence utils image
	// +optional
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
	// Logging allows configuration of Coherence and java util logging.
	// +optional
	Logging *LoggingSpec `json:"logging,omitempty"`
	// The JVM specific options
	// +optional
	JVM *JVMSpec `json:"jvm,omitempty"`
	// Ports specifies additional port mappings for the Pod and additional Services for those ports
	// +listType=map
	// +listMapKey=name
	// +optional
	Ports []NamedPortSpec `json:"ports,omitempty"`
	// Env is additional environment variable mappings that will be passed to
	// the Coherence container in the Pod.
	// To specify extra variables add them as name value pairs the same as they
	// would be added to a Pod containers spec.
	// +listType=map
	// +listMapKey=name
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
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
	// The configuration to control safe scaling.
	// +optional
	Scaling *ScalingSpec `json:"scaling,omitempty"`
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
	// Annotations are free-form yaml that will be added to the store release as annotations
	// Any annotations should be placed BELOW this annotations: key. For example if we wanted to
	// include annotations for Prometheus it would look like this:
	//
	// annotations:
	//   prometheus.io/scrape: "true"
	//   prometheus.io/port: "2408"
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// The extra labels to add to the all of the Pods in this deployments.
	// Labels here will add to or override those defined for the cluster.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
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
	// The deployments that must be started before this deployment can start.
	// +listType=map
	// +listMapKey=deployment
	// +optional
	StartQuorum []StartQuorum `json:"startQuorum,omitempty"`
	// List of additional initialization containers to add to the deployment's Pod.
	// Init containers cannot be added or removed.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	// +listType=map
	// +listMapKey=name
	AdditionalInitContainers []corev1.Container `json:"additionalInitContainers,omitempty"`
	// List of additional containers to add to the deployment's Pod.
	// Containers cannot be added or removed.
	// +listType=map
	// +listMapKey=name
	AdditionalContainers []corev1.Container `json:"additionalContainers,omitempty"`
}

// Obtain the number of replicas required for a deployment.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replica value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *CoherenceDeploymentSpec) GetReplicas() int32 {
	if in == nil {
		return 0
	}
	if in.Replicas == nil {
		return DefaultReplicas
	}
	return *in.Replicas
}

// Set the number of replicas required for a deployment.
func (in *CoherenceDeploymentSpec) SetReplicas(replicas int32) {
	if in != nil {
		in.Replicas = &replicas
	}
}

func (in *CoherenceDeploymentSpec) GetCoherenceImage() *string {
	if in != nil {
		return in.Image
	}
	return nil
}

// Ensure that the Coherence image is set for the deployment.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (in *CoherenceDeploymentSpec) EnsureCoherenceImage(coherenceImage *string) bool {
	if in.Coherence == nil {
		in.Coherence = &CoherenceSpec{}
	}

	return in.EnsureImage(coherenceImage)
}

func (in *CoherenceDeploymentSpec) GetCoherenceUtilsImage() *string {
	if in != nil && in.CoherenceUtils != nil {
		return in.CoherenceUtils.Image
	}
	return nil
}

// Ensure that the Coherence Utils image is set for the deployment.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (in *CoherenceDeploymentSpec) EnsureCoherenceUtilsImage(utilsImage *string) bool {
	if in.CoherenceUtils == nil {
		in.CoherenceUtils = &ImageSpec{}
	}

	return in.CoherenceUtils.EnsureImage(utilsImage)
}

func (in *CoherenceDeploymentSpec) GetEffectiveScalingPolicy() ScalingPolicy {
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
func (in *CoherenceDeploymentSpec) GetHealthPort() int32 {
	if in == nil || in.HealthPort == nil || *in.HealthPort <= 0 {
		return DefaultHealthPort
	}
	return *in.HealthPort
}

// Returns the ScalingProbe to use for checking Phase HA for the deployment.
// This method will not return nil.
func (in *CoherenceDeploymentSpec) GetScalingProbe() *ScalingProbe {
	if in == nil || in.Scaling == nil || in.Scaling.Probe == nil {
		return in.GetDefaultScalingProbe()
	}
	return in.Scaling.Probe
}

// Obtain a default ScalingProbe
func (in *CoherenceDeploymentSpec) GetDefaultScalingProbe() *ScalingProbe {
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
func (in *CoherenceDeploymentSpec) CreateKubernetesResources(d *CoherenceDeployment, flags *flags.CoherenceOperatorFlags) (Resources, error) {
	var res []Resource

	if in.GetReplicas() <= 0 {
		// replicas is zero so nothing to create
		return Resources{Items: res}, nil
	}

	// Create the fluentd ConfigMap if required
	cm, err := in.Logging.CreateFluentdConfigMap(d)
	if err != nil {
		return Resources{}, err
	}
	if cm != nil {
		res = append(res, *cm)
	}

	// Create the headless WKA Service
	res = append(res, in.CreateWKAService(d))

	// Create the headless Service
	res = append(res, in.CreateHeadlessService(d))

	// Create the StatefulSet
	res = append(res, in.CreateStatefulSet(d, flags))

	// Create the Services for each port (and optionally ServiceMonitors)
	res = append(res, in.CreateServicesForPort(d)...)

	return Resources{Items: res}, nil
}

// Create the Services for each port (and optionally ServiceMonitors)
func (in *CoherenceDeploymentSpec) CreateServicesForPort(deployment *CoherenceDeployment) []Resource {
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
func (in *CoherenceDeploymentSpec) CreatePodSelectorLabels(deployment *CoherenceDeployment) map[string]string {
	selector := deployment.CreateCommonLabels()
	selector[LabelComponent] = LabelComponentCoherencePod
	return selector
}

// Create the headless WKA Service
func (in *CoherenceDeploymentSpec) CreateWKAService(deployment *CoherenceDeployment) Resource {
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
			Name:      deployment.GetWkaServiceName(),
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
func (in *CoherenceDeploymentSpec) CreateHeadlessService(deployment *CoherenceDeployment) Resource {
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
func (in *CoherenceDeploymentSpec) CreateStatefulSet(deployment *CoherenceDeployment, flags *flags.CoherenceOperatorFlags) Resource {
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
	// Add any logging settings
	in.Logging.UpdateStatefulSet(&sts)

	// Add any additional init-containers
	if in.AdditionalInitContainers != nil && len(in.AdditionalInitContainers) > 0 {
		sts.Spec.Template.Spec.InitContainers = append(sts.Spec.Template.Spec.InitContainers, in.AdditionalInitContainers...)
	}

	// Add any additional containers
	if in.AdditionalContainers != nil && len(in.AdditionalContainers) > 0 {
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, in.AdditionalContainers...)
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

func (in *CoherenceDeploymentSpec) GetImagePullSecrets() []corev1.LocalObjectReference {
	var secrets []corev1.LocalObjectReference

	for _, s := range in.ImagePullSecrets {
		secrets = append(secrets, corev1.LocalObjectReference{
			Name: s.Name,
		})
	}

	return secrets
}

// Get the service account name for the cluster.
func (in *CoherenceDeploymentSpec) GetServiceAccountName() string {
	if in != nil && in.ServiceAccountName != DefaultServiceAccount {
		return in.ServiceAccountName
	}
	return ""
}

// Create the Coherence container spec.
func (in *CoherenceDeploymentSpec) CreateCoherenceContainer(deployment *CoherenceDeployment, flags *flags.CoherenceOperatorFlags) corev1.Container {
	var cohImage *string

	if in.Image == nil {
		cohImage = flags.GetCoherenceImage()
	} else {
		cohImage = in.Image
	}

	healthPort := in.GetHealthPort()

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
		VolumeMounts: []corev1.VolumeMount{
			{Name: VolumeNameLogs, MountPath: VolumeMountPathLogs},
			{Name: VolumeNameUtils, MountPath: VolumeMountPathUtils},
			{Name: VolumeNameJVM, MountPath: VolumeMountPathJVM},
		},
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

// Create the default environment variables.
func (in *CoherenceDeploymentSpec) CreateDefaultEnv(deployment *CoherenceDeployment) []corev1.EnvVar {
	return []corev1.EnvVar{
		{Name: EnvVarCohWka, Value: deployment.GetWkaServiceName()},
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
		{
			Name: EnvVarOperatorHost, ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: OperatorConfigName},
					Key:                  OperatorConfigKeyHost,
					Optional:             pointer.BoolPtr(true),
				},
			},
		},
		{Name: EnvVarCohSite, Value: "http://$(OPERATOR_HOST)/site/$(COH_MACHINE_NAME)"},
		{Name: EnvVarCohRack, Value: "http://$(OPERATOR_HOST)/rack/$(COH_MACHINE_NAME)"},
		{Name: EnvVarCohClusterName, Value: deployment.GetCoherenceClusterName()},
		{Name: EnvVarCohRole, Value: deployment.GetRoleName()},
		{Name: EnvVarCohUtilDir, Value: VolumeMountPathUtils},
		{Name: EnvVarOperatorTimeout, Value: Int32PtrToStringWithDefault(in.OperatorRequestTimeout, 120)},
		{Name: EnvVarCohHealthPort, Value: Int32ToString(in.GetHealthPort())},
	}
}

// Create the default Container resources.
func (in *CoherenceDeploymentSpec) CreateDefaultResources() corev1.ResourceRequirements {
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
func (in *CoherenceDeploymentSpec) CreateDefaultReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 30,
		PeriodSeconds:       60,
		TimeoutSeconds:      30,
		SuccessThreshold:    1,
		FailureThreshold:    50,
	}
}

// Update the probe with the default readiness probe action.
func (in *CoherenceDeploymentSpec) UpdateDefaultReadinessProbeAction(probe *corev1.Probe) *corev1.Probe {
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   DefaultReadinessPath,
		Port:   intstr.FromInt(int(DefaultHealthPort)),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}

// Create the default liveness probe.
func (in *CoherenceDeploymentSpec) CreateDefaultLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 60,
		PeriodSeconds:       60,
		TimeoutSeconds:      30,
		SuccessThreshold:    1,
		FailureThreshold:    5,
	}
}

// Update the probe with the default liveness probe action.
func (in *CoherenceDeploymentSpec) UpdateDefaultLivenessProbeAction(probe *corev1.Probe) *corev1.Probe {
	probe.HTTPGet = &corev1.HTTPGetAction{
		Path:   DefaultLivenessPath,
		Port:   intstr.FromInt(int(DefaultHealthPort)),
		Scheme: corev1.URISchemeHTTP,
	}
	return probe
}

// Get the Utils init-container spec.
func (in *CoherenceDeploymentSpec) CreateUtilsContainer(deployment *CoherenceDeployment, flags *flags.CoherenceOperatorFlags) corev1.Container {
	var utilsImage *string
	if in.CoherenceUtils == nil || in.CoherenceUtils.Image == nil {
		utilsImage = flags.GetCoherenceUtilsImage()
	} else {
		utilsImage = in.CoherenceUtils.Image
	}

	c := corev1.Container{
		Name:    ContainerNameUtils,
		Image:   *utilsImage,
		Command: []string{UtilsInitCommand},
		Env: []corev1.EnvVar{
			{Name: "COH_UTIL_DIR", Value: VolumeMountPathUtils},
			{Name: "COH_CLUSTER_NAME", Value: deployment.GetCoherenceClusterName()},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: VolumeNameUtils, MountPath: VolumeMountPathUtils},
		},
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
func (in *CoherenceDeploymentSpec) EnsurePodAffinity(deployment *CoherenceDeployment) *corev1.Affinity {
	if in != nil && in.Affinity != nil {
		return in.Affinity
	}
	// return the default affinity which attempts to spread the Pods for a deployment across fault domains
	return in.CreateDefaultPodAffinity(deployment)
}

// Create the default Pod Affinity to use in a deployment's StatefulSet.
func (in *CoherenceDeploymentSpec) CreateDefaultPodAffinity(deployment *CoherenceDeployment) *corev1.Affinity {
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

func (in *CoherenceDeploymentSpec) GetMetricsPort() int32 {
	if in == nil {
		return 0
	}
	return in.Coherence.GetMetricsPort()
}

func (in *CoherenceDeploymentSpec) GetManagementPort() int32 {
	if in == nil {
		return 0
	}
	return in.Coherence.GetManagementPort()
}
