/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package webhook

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

type HookInstaller struct {
	Clients               clients.ClientSet
	certManagerGroup      string
	certManagerAPIVersion string
	issuer                *unstructured.Unstructured
	certificate           *unstructured.Unstructured
}

type certManagerVersion struct {
	group    string
	versions []string
}

const (
	admissionAPI        = "admissionregistration.k8s.io/v1"
	certManagerCertName = "coherence-webhook-server-certificate"
	certTypeAnnotation  = "operator.coherence.oracle.com/cert-type"
)

var (
	log = logf.Log.WithName(controllerName)

	// Cert-Manager APIs that we can detect
	certManagerAPIs = []certManagerVersion{
		{group: "cert-manager.io", versions: []string{"v1"}},
	}
)

func (k *HookInstaller) uninstallWebHook() error {
	log.Info("Uninstall webhook resources")

	// We only clean up cert-manager resource here.
	// We specifically DO NOT clean up the web-hook resources because we do not
	// want mutations of Coherence resources to go through whilst the operator is not
	// running as these may result in invalid configurations.

	if k.certificate != nil {
		log.Info("deleting cert-manager certificate " + k.certificate.GetName())
		if err := k.uninstallUnstructured(k.certificate); err != nil {
			log.Error(err, "error deleting cert-manager Certificate "+k.certificate.GetName())
		}
	}
	if k.issuer != nil {
		log.Info("deleting cert-manager issuer " + k.issuer.GetName())
		if err := k.uninstallUnstructured(k.issuer); err != nil {
			log.Error(err, "error deleting cert-manager Issuer "+k.issuer.GetName())
		}
	}

	return nil
}

func (k *HookInstaller) InstallWithCertManager(ctx context.Context) error {
	var err error
	if err = k.validateCertManagerInstallation(); err != nil {
		return err
	}
	// install the cert-manager Issuer
	if err = k.installUnstructured(ctx, k.issuer); err != nil {
		return err
	}
	// install the cert-manager Certificate
	if err = k.installUnstructured(ctx, k.certificate); err != nil {
		return err
	}
	// Install the webhooks
	ns := operator.GetNamespace()
	m := createMutatingWebhookWithCertManager(ns, k.certManagerGroup)
	if err = installMutatingWebhook(ctx, k.Clients, m); err != nil {
		return err
	}
	v := createValidatingWebhookWithCertManager(ns, k.certManagerGroup)

	err = installValidatingWebhook(ctx, k.Clients, v)
	return err
}

func baseWebhookSecret(ns string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      viper.GetString(operator.FlagWebhookSecret),
			Namespace: ns,
		},
		Type: "kubernetes.io/tls",
	}
}

func (k *HookInstaller) detectCertManagerVersion() error {
	for _, api := range certManagerAPIs {
		group, version, err := k.detectCertManagerCRD(api)
		if err != nil {
			return err
		}
		if group != "" && version != "" {
			log.Info(fmt.Sprintf("Detected cert-manager CRDs %s/%s", k.certManagerGroup, k.certManagerAPIVersion))

			if !contains(api.versions, version) {
				return errors.Wrap(err, fmt.Sprintf("Detected cert-manager CRDs with version %s, only versions %v are fully supported. Certificates for webhooks may not work.", version, api.versions))
			}

			k.certManagerGroup = group
			k.certManagerAPIVersion = version

			log.Info(fmt.Sprintf("Detected cert-manager %s/%s", group, version))
			return nil
		}
	}
	return fmt.Errorf("failed to detect any valid cert-manager CRDs. Make sure cert-manager is installed")
}

func (k *HookInstaller) detectCertManagerCRD(api certManagerVersion) (string, string, error) {
	testCRD := fmt.Sprintf("certificates.%s", api.group)
	log.Info(fmt.Sprintf("Try to retrieve cert-manager CRD %s", testCRD))
	crd, err := k.Clients.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), testCRD, metav1.GetOptions{})
	if err == nil {
		// crd.Spec.Versions[0] must be the one that is stored and served, we should use that one
		log.Info(fmt.Sprintf("Got CRD. Group: %s, Version: %s", api.group, crd.Spec.Versions[0].Name))
		return api.group, crd.Spec.Versions[0].Name, nil
	}
	if !kerrors.IsNotFound(err) {
		return "", "", fmt.Errorf("failed to detect cert manager CRD %s: %v", testCRD, err)
	}
	return "", "", nil
}

func (k *HookInstaller) validateCertManagerInstallation() error {
	if err := k.detectCertManagerVersion(); err != nil {
		return err
	}

	certificateCRD := fmt.Sprintf("certificates.%s", k.certManagerGroup)
	if err := k.validateCrdVersion(certificateCRD, k.certManagerAPIVersion); err != nil {
		return err
	}
	issuerCRD := fmt.Sprintf("issuers.%s", k.certManagerGroup)
	if err := k.validateCrdVersion(issuerCRD, k.certManagerAPIVersion); err != nil {
		return err
	}

	// Initialize the custom resources that we're going to install
	k.certificate = certificate(operator.GetNamespace(), k.certManagerGroup, k.certManagerAPIVersion)
	k.issuer = issuer(operator.GetNamespace(), k.certManagerGroup, k.certManagerAPIVersion)

	return nil
}

func (k *HookInstaller) validateCrdVersion(crdName string, expectedVersion string) error {
	certCRD, err := k.Clients.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), crdName, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return errors.Wrap(err, fmt.Sprintf("failed to find CRD '%s': %s", crdName, err))
		}
		return err
	}
	crdVersion := certCRD.Spec.Versions[0].Name

	if crdVersion != expectedVersion {
		return errors.Wrap(err, fmt.Sprintf("invalid CRD version found for '%s': %s instead of %s", crdName, crdVersion, expectedVersion))
	}
	log.Info(fmt.Sprintf("CRD %s is installed with version %s", crdName, crdVersion))
	return nil
}

func (k *HookInstaller) installUnstructured(ctx context.Context, item *unstructured.Unstructured) error {
	var err error

	gvk := item.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: fmt.Sprintf("%ss", strings.ToLower(gvk.Kind)), // since we know what kinds are we dealing with here, this is OK
	}

	client := k.Clients.DynamicClient.Resource(gvr).Namespace(item.GetNamespace())

	_, err = client.Get(ctx, item.GetName(), metav1.GetOptions{})
	if err != nil {
		// issuer does not exist, so create it
		_, err = client.Create(context.TODO(), item, metav1.CreateOptions{})
		if kerrors.IsAlreadyExists(err) {
			log.Info(fmt.Sprintf("resource %s/%s %s already registered and will not be overwritten", item.GetAPIVersion(), item.GetKind(), item.GetName()))
		} else if err != nil {
			return fmt.Errorf("error when creating resource %s/%s. %v", item.GetName(), item.GetNamespace(), err)
		}
		log.Info(fmt.Sprintf("created resource %s/%s %s", item.GetAPIVersion(), item.GetKind(), item.GetName()))
	} else {
		log.Info(fmt.Sprintf("resource %s/%s %s already registered and will not be overwritten", item.GetAPIVersion(), item.GetKind(), item.GetName()))
	}

	return nil
}

func (k *HookInstaller) uninstallUnstructured(item *unstructured.Unstructured) error {
	gvk := item.GroupVersionKind()
	err := k.Clients.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: fmt.Sprintf("%ss", strings.ToLower(gvk.Kind)), // since we know what kinds are we dealing with here, this is OK
	}).Namespace(item.GetNamespace()).Delete(context.TODO(), item.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("error when deleting resource %s/%s. %v", item.GetName(), item.GetNamespace(), err)
	}
	return nil
}

func installMutatingWebhook(ctx context.Context, c clients.ClientSet, webhook admissionv1.MutatingWebhookConfiguration) error {
	log.Info(fmt.Sprintf("installing webhook %s/%s", webhook.Namespace, webhook.Name))
	cl := c.KubeClient.AdmissionregistrationV1()
	existing, err := cl.MutatingWebhookConfigurations().Get(context.TODO(), webhook.GetName(), metav1.GetOptions{})
	exists := err == nil

	if exists && existing != nil {
		existing.Webhooks = webhook.Webhooks
		existing.Annotations = webhook.Annotations
		_, err = cl.MutatingWebhookConfigurations().Update(ctx, existing, metav1.UpdateOptions{})
	} else {
		_, err = cl.MutatingWebhookConfigurations().Create(ctx, &webhook, metav1.CreateOptions{})
	}
	return err
}

func installValidatingWebhook(ctx context.Context, c clients.ClientSet, webhook admissionv1.ValidatingWebhookConfiguration) error {
	log.Info(fmt.Sprintf("installing webhook %s/%s", webhook.Namespace, webhook.Name))
	cl := c.KubeClient.AdmissionregistrationV1()
	existing, err := cl.ValidatingWebhookConfigurations().Get(context.TODO(), webhook.GetName(), metav1.GetOptions{})
	exists := err == nil

	if exists && existing != nil {
		existing.Webhooks = webhook.Webhooks
		existing.Annotations = webhook.Annotations
		_, err = cl.ValidatingWebhookConfigurations().Update(ctx, existing, metav1.UpdateOptions{})
	} else {
		_, err = cl.ValidatingWebhookConfigurations().Create(ctx, &webhook, metav1.CreateOptions{})
	}
	return err
}

func createMutatingWebhookWithCABundle(ns string, caData []byte) admissionv1.MutatingWebhookConfiguration {
	cfg := createMutatingWebhookConfiguration(ns)
	for i := range cfg.Webhooks {
		cfg.Webhooks[i].ClientConfig.CABundle = caData
	}
	return cfg
}

func createValidatingWebhookWithCABundle(ns string, caData []byte) admissionv1.ValidatingWebhookConfiguration {
	cfg := createValidatingWebhookConfiguration(ns)
	for i := range cfg.Webhooks {
		cfg.Webhooks[i].ClientConfig.CABundle = caData
	}
	return cfg
}

func createMutatingWebhookWithCertManager(ns string, certManagerGroup string) admissionv1.MutatingWebhookConfiguration {
	cfg := createMutatingWebhookConfiguration(ns)
	injectCaAnnotationName := fmt.Sprintf("%s/inject-ca-from", certManagerGroup)
	cfg.Annotations[injectCaAnnotationName] = fmt.Sprintf("%s/%s", ns, certManagerCertName)
	return cfg
}

func createValidatingWebhookWithCertManager(ns string, certManagerGroup string) admissionv1.ValidatingWebhookConfiguration {
	cfg := createValidatingWebhookConfiguration(ns)
	injectCaAnnotationName := fmt.Sprintf("%s/inject-ca-from", certManagerGroup)
	cfg.Annotations[injectCaAnnotationName] = fmt.Sprintf("%s/%s", ns, certManagerCertName)
	return cfg
}

func createMutatingWebhookConfiguration(ns string) admissionv1.MutatingWebhookConfiguration {
	namespacedScope := admissionv1.NamespacedScope
	failedType := admissionv1.Fail
	equivalentType := admissionv1.Equivalent
	noSideEffects := admissionv1.SideEffectClassNone
	path := coh.MutatingWebHookPath
	clientConfig := createWebhookClientConfig(ns, path)

	return admissionv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: viper.GetString(operator.FlagMutatingWebhookName),
			Annotations: map[string]string{
				certTypeAnnotation: viper.GetString(operator.FlagCertType),
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "MutatingWebhookConfiguration",
			APIVersion: admissionAPI,
		},
		Webhooks: []admissionv1.MutatingWebhook{
			{
				Name:                    "coherence.oracle.com",
				AdmissionReviewVersions: []string{"v1", "v1beta1"},
				Rules: []admissionv1.RuleWithOperations{
					{
						Operations: []admissionv1.OperationType{"CREATE", "UPDATE"},
						Rule: admissionv1.Rule{
							APIGroups:   []string{"coherence.oracle.com"},
							APIVersions: []string{"v1"},
							Resources:   []string{"coherence"},
							Scope:       &namespacedScope,
						},
					},
				},
				FailurePolicy: &failedType, // this means that the request to update instance would fail, if webhook is not up
				MatchPolicy:   &equivalentType,
				SideEffects:   &noSideEffects,
				ClientConfig:  clientConfig,
			},
		},
	}
}

func createValidatingWebhookConfiguration(ns string) admissionv1.ValidatingWebhookConfiguration {
	namespacedScope := admissionv1.NamespacedScope
	failedType := admissionv1.Fail
	equivalentType := admissionv1.Equivalent
	noSideEffects := admissionv1.SideEffectClassNone
	path := coh.ValidatingWebHookPath
	clientConfig := createWebhookClientConfig(ns, path)

	return admissionv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: viper.GetString(operator.FlagValidatingWebhookName),
			Annotations: map[string]string{
				certTypeAnnotation: viper.GetString(operator.FlagCertType),
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ValidatingWebhookConfiguration",
			APIVersion: admissionAPI,
		},
		Webhooks: []admissionv1.ValidatingWebhook{
			{
				Name:                    "coherence.oracle.com",
				AdmissionReviewVersions: []string{"v1", "v1beta1"},
				Rules: []admissionv1.RuleWithOperations{
					{
						Operations: []admissionv1.OperationType{"CREATE", "UPDATE"},
						Rule: admissionv1.Rule{
							APIGroups:   []string{"coherence.oracle.com"},
							APIVersions: []string{"v1"},
							Resources:   []string{"coherence"},
							Scope:       &namespacedScope,
						},
					},
				},
				FailurePolicy: &failedType, // this means that the request to update instance would fail, if webhook is not up
				MatchPolicy:   &equivalentType,
				SideEffects:   &noSideEffects,
				ClientConfig:  clientConfig,
			},
		},
	}
}

func createWebhookClientConfig(ns, path string) admissionv1.WebhookClientConfig {

	var clientConfig admissionv1.WebhookClientConfig
	if operator.IsDevMode() {
		hn := operator.GetWebhookServiceDNSNames()[0]
		url := fmt.Sprintf("https://%s:9443%s", hn, path)
		clientConfig = admissionv1.WebhookClientConfig{
			URL: &url,
		}
	} else {
		clientConfig = admissionv1.WebhookClientConfig{
			Service: &admissionv1.ServiceReference{
				Name:      viper.GetString(operator.FlagWebhookService),
				Namespace: ns,
				Path:      &path,
			},
		}
	}
	return clientConfig
}

func issuer(ns string, group string, apiVersion string) *unstructured.Unstructured {
	apiString := fmt.Sprintf("%s/%s", group, apiVersion)
	certIssuer := viper.GetString(operator.FlagCertIssuer)
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiString,
			"kind":       "Issuer",
			"metadata": map[string]interface{}{
				"name":      certIssuer,
				"namespace": ns,
			},
			"spec": map[string]interface{}{
				"selfSigned": map[string]interface{}{},
			},
		},
	}
}

func certificate(ns string, group string, apiVersion string) *unstructured.Unstructured {
	apiString := fmt.Sprintf("%s/%s", group, apiVersion)
	name := viper.GetString(operator.FlagWebhookService)
	certIssuer := viper.GetString(operator.FlagCertIssuer)
	dns := operator.GetWebhookServiceDNSNames()
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiString,
			"kind":       "Certificate",
			"metadata": map[string]interface{}{
				"name":      certManagerCertName,
				"namespace": ns,
			},
			"spec": map[string]interface{}{
				"commonName": fmt.Sprintf("%s.%s.svc", name, ns),
				"dnsNames":   dns,
				"issuerRef": map[string]interface{}{
					"kind": "Issuer",
					"name": certIssuer,
				},
				"secretName": viper.GetString(operator.FlagWebhookSecret),
			},
		},
	}
}

// Contains returns true if an element is present in a iteratee.
func contains(in interface{}, elem interface{}) bool {
	inValue := reflect.ValueOf(in)
	elemValue := reflect.ValueOf(elem)
	inType := inValue.Type()

	switch inType.Kind() {
	case reflect.String:
		return strings.Contains(inValue.String(), elemValue.String())
	case reflect.Map:
		for _, key := range inValue.MapKeys() {
			if equal(key.Interface(), elem) {
				return true
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < inValue.Len(); i++ {
			if equal(inValue.Index(i).Interface(), elem) {
				return true
			}
		}
	default:
		panic(fmt.Sprintf("Type %s is not supported by Contains, supported types are String, Map, Slice, Array", inType.String()))
	}

	return false
}

func equal(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	return reflect.DeepEqual(expected, actual)
}
