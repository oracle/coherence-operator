package v1

import (
	"errors"
	"github.com/ghodss/yaml"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// CoherenceClusterSpec defines the desired state of CoherenceCluster
// +k8s:openapi-gen=true
type CoherenceClusterSpec struct {
	// The secrets to be used when pulling images. Secrets must be manually created in the target namespace.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	// +optional
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// This spec is either the spec of a single role cluster or is used as the
	// default values applied to roles in Roles array.
	CoherenceRoleSpec `json:",inline"`
	// Roles is the list of different roles in the cluster
	// There must be at least one role in a cluster.
	// +optional
	Roles []CoherenceRoleSpec `json:"roles,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceCluster is the Schema for the coherenceclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CoherenceCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceClusterSpec   `json:"spec,omitempty"`
	Status CoherenceClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceClusterList contains a list of CoherenceCluster
type CoherenceClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceCluster `json:"items"`
}

// CoherenceClusterStatus defines the observed state of CoherenceCluster
// +k8s:openapi-gen=true
type CoherenceClusterStatus struct {
}

func init() {
	SchemeBuilder.Register(&CoherenceCluster{}, &CoherenceClusterList{})
}

func (c *CoherenceCluster) GetWkaServiceName() string {
	if c == nil {
		return ""
	}
	return c.Name + WKAServiceNameSuffix
}

// Obtain the CoherenceRoleSpec for the specified role name
func (c *CoherenceCluster) GetRole(name string) CoherenceRoleSpec {
	if len(c.Spec.Roles) > 0 {
		for _, role := range c.Spec.Roles {
			if role.GetRoleName() == name {
				return role
			}
		}
	} else if name == c.Spec.CoherenceRoleSpec.GetRoleName() {
		return c.Spec.CoherenceRoleSpec
	}
	return CoherenceRoleSpec{}
}

// Set the CoherenceRoleSpec
func (c *CoherenceCluster) SetRole(spec CoherenceRoleSpec) {
	name := spec.GetRoleName()
	if len(c.Spec.Roles) > 0 {
		for index, role := range c.Spec.Roles {
			if role.GetRoleName() == name {
				c.Spec.Roles[index] = spec
				break
			}
		}
	} else if name == c.Spec.CoherenceRoleSpec.GetRoleName() {
		c.Spec.CoherenceRoleSpec = spec
	}
}

// Load this CoherenceCluster from the specified yaml file
func (c *CoherenceCluster) FromYaml(files ...string) error {
	return c.loadYaml(files...)
}

// NewCoherenceClusterFromYaml creates a new CoherenceCluster from a yaml file.
func NewCoherenceClusterFromYaml(namespace string, file ...string) (CoherenceCluster, error) {
	c := CoherenceCluster{}
	err := c.loadYaml(file...)

	if namespace != "" {
		c.SetNamespace(namespace)
	}

	return c, err
}

func (c *CoherenceCluster) loadYaml(files ...string) error {
	if c == nil || files == nil {
		return nil
	}

	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			if !strings.HasPrefix(file, "/") {
				// the file does not exist so try relative to the caller's file location.
				_, caller, _, ok := runtime.Caller(2)
				if ok {
					dir := filepath.Dir(caller)
					file = dir + string(os.PathSeparator) + file
					_, e := os.Stat(file)
					if e != nil {
						return errors.New(err.Error() + "\n" + e.Error())
					}
				}
			} else {
				// file does not exist
				return err
			}
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.New("Failed to read file " + file + " caused by " + err.Error())
		}

		err = yaml.Unmarshal(data, c)
		if err != nil {
			return errors.New("Failed to parse yaml file " + file + " caused by " + err.Error())
		}
	}

	return nil
}
