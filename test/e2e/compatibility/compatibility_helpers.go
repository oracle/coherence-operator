/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"strings"
)

/*
func findContainer(name string, d *appsv1.Deployment) *corev1.Container {
	for _, c := range d.Spec.Template.Spec.Containers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func findEnvVar(name string, c *corev1.Container) *corev1.EnvVar {
	for _, e := range c.Env {
		if e.Name == name {
			return &e
		}
	}
	return nil
}

func helmInstall(args ...string) (*HelmInstallResult, error) {
	chart, err := helper.FindOperatorHelmChartDir()
	if err != nil {
		return nil, err
	}

	ns := helper.GetTestNamespace()

	args = append([]string{"install", "--dry-run", "-o", "json"}, args...)
	args = append(args, "--namespace", ns, "operator", chart)

	cmd := exec.Command("helm", args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	manifest := m["manifest"]

	return parseHelmManifest(fmt.Sprintf("%v", manifest))
}

func parseHelmManifest(manifest string) (*HelmInstallResult, error) {
	resources := make(map[schema.GroupVersionResource]map[string]runtime.Object)
	s := scheme.Scheme
	decoder := scheme.Codecs.UniversalDecoder()

	parts := strings.Split(manifest, "\n---\n")
	list := make([]runtime.Object, len(parts))
	ownerRefs := make([]metav1.OwnerReference, 0)

	index := 0
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			u := unstructured.Unstructured{}
			err := yaml.Unmarshal([]byte(trimmed), &u)
			if err != nil {
				return nil, err
			}
			gvr, _ := meta.UnsafeGuessKindToResource(u.GroupVersionKind())

			// remove owner references
			u.SetOwnerReferences(ownerRefs)
			data, err := yaml.Marshal(u.Object)
			if err != nil {
				return nil, err
			}

			o, err := s.New(u.GroupVersionKind())
			if err != nil {
				return nil, err
			}
			_, _, err = decoder.Decode(data, nil, o)
			if err != nil {
				return nil, err
			}

			m, ok := resources[gvr]
			if !ok {
				m = make(map[string]runtime.Object)
			}
			list[index] = o
			index++
			m[u.GetName()] = o
			resources[gvr] = m
		}
	}

	ordered := list[0:index]
	return &HelmInstallResult{resources: resources, ordered: ordered, decoder: decoder}, nil
}
*/

type HelmInstallResult struct {
	resources map[schema.GroupVersionResource]map[string]runtime.Object
	ordered   []runtime.Object
	decoder   runtime.Decoder
}

type HelmInstallResultFilter func(runtime.Object) bool

func (h *HelmInstallResult) ToString(filter HelmInstallResultFilter, w io.Writer) error {
	var sep = []byte("\n---\n")

	for _, res := range h.ordered {
		if filter == nil || filter(res) {
			_, err := w.Write(sep)
			if err != nil {
				return err
			}

			d, err := yaml.Marshal(res)
			if err != nil {
				return err
			}
			_, err = w.Write(d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *HelmInstallResult) Size() int {
	if h == nil {
		return 0
	}
	return len(h.ordered)
}
func (h *HelmInstallResult) Get(name string, o runtime.Object) error {
	if h == nil {
		return fmt.Errorf("resource '%s' not found", name)
	}

	gvr, err := h.getGVRFromObject(o, scheme.Scheme)
	if err != nil {
		return err
	}

	if h.resources == nil {
		return k8serr.NewNotFound(gvr.GroupResource(), name)
	}

	m, ok := h.resources[gvr]
	if !ok {
		return k8serr.NewNotFound(gvr.GroupResource(), name)
	}

	item, ok := m[name]
	if !ok {
		return k8serr.NewNotFound(gvr.GroupResource(), name)
	}

	j, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, _, err = h.decoder.Decode(j, nil, o)
	if err != nil {
		return err
	}

	return nil
}

func (h *HelmInstallResult) List(list runtime.Object) error {
	if h == nil || h.resources == nil {
		return nil
	}

	gvk, err := getGVKFromList(list, scheme.Scheme)
	if err != nil {
		return err
	}

	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	m, ok := h.resources[gvr]
	if ok {
		items := make([]runtime.Object, len(m))
		i := 0
		for _, o := range m {
			items[i] = o.DeepCopyObject()
			i++
		}

		if err := meta.SetList(list, items); err != nil {
			return err
		}
	}

	return nil
}

func (h *HelmInstallResult) getGVRFromObject(obj runtime.Object, scheme *runtime.Scheme) (schema.GroupVersionResource, error) {
	gvk, err := apiutil.GVKForObject(obj, scheme)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	return gvr, nil
}

func getGVKFromList(list runtime.Object, scheme *runtime.Scheme) (schema.GroupVersionKind, error) {
	gvk, err := apiutil.GVKForObject(list, scheme)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	if gvk.Kind == "List" {
		return schema.GroupVersionKind{}, fmt.Errorf("cannot derive GVK for generic List type %T (kind %q)", list, gvk)
	}

	if !strings.HasSuffix(gvk.Kind, "List") {
		return schema.GroupVersionKind{}, fmt.Errorf("non-list type %T (kind %q) passed as output", list, gvk)
	}
	// we need the non-list GVK, so chop off the "List" from the end of the kind
	gvk.Kind = gvk.Kind[:len(gvk.Kind)-4]
	return gvk, nil
}
