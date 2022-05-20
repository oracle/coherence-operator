/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package webhook

import (
	"bytes"
	"context"
	"fmt"
	"github.com/oracle/coherence-operator/controllers/predicates"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/pkg/certs"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const (
	// The name of this controller. This is used in events, log messages, etc.
	controllerName = "controllers.Certs"
)

// blank assignment to verify that CertReconciler implements reconcile.Reconciler.
// If the reconcile.Reconciler API was to change then we'd get a compile error here.
var _ reconcile.Reconciler = &CertReconciler{}

type CertReconciler struct {
	reconciler.CommonReconciler
	Clientset     clients.ClientSet
	rotateBefore  time.Duration
	hookInstaller *HookInstaller
}

func (r *CertReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	r.SetCommonReconciler(controllerName, mgr)
	r.rotateBefore = operator.GetCACertRotateBefore()

	// determine how webhook certs will be managed
	switch {
	case operator.ShouldUseCertManager():
		r.hookInstaller = &HookInstaller{Clients: r.Clientset}
		if err := r.hookInstaller.InstallWithCertManager(ctx); err != nil {
			return errors.Wrap(err, " unable to install cert-manager resources")
		}
		// if in dev-mode write the certs to local cert files
		if err := r.writeLocalCerts(ctx); err != nil {
			return err
		}
	case operator.ShouldUseSelfSignedCerts():
		// do an initial reconcile to make sure certs and web hooks are configured
		if err := r.ReconcileResources(ctx); err != nil {
			return errors.Wrap(err, " unable to setup and fill the webhook certificates")
		}
	default:
		// certificates are manually managed
		if err := r.writeLocalCerts(ctx); err != nil {
			return err
		}
		// don't use this controller for manual certs so just return
		return nil
	}

	// set-up this controller
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		Named("coherence-webhook-server-cert").
		WithEventFilter(&predicates.NamedPredicate{
			Namespace: operator.GetNamespace(),
			Name:      viper.GetString(operator.FlagWebhookSecret),
		}).
		Complete(r)
}

func (r *CertReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	if err := r.ReconcileResources(ctx); err != nil {
		return reconcile.Result{}, err
	}

	namespace := operator.GetNamespace()
	secretName := request.Name
	secret, err := r.Clientset.KubeClient.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return reconcile.Result{}, err
	}

	if operator.ShouldUseCertManager() {
		// for cert-manager certs we don't need to do anything except update the local certs
		if operator.IsDevMode() {
			if err := r.writeLocalCertsFromSecret(secret); err != nil {
				r.GetLog().Error(err, "error writing local certs")
			}
		}
		return reconcile.Result{}, nil
	}

	serverCA := certs.BuildCAFromSecret(*secret)
	if serverCA == nil {
		return reconcile.Result{}, fmt.Errorf("cannot find CA in webhook secret %s/%s", request.Namespace, secretName)
	}

	return reconcile.Result{
		RequeueAfter: certs.RotateIn(time.Now(), serverCA.Cert.NotAfter, r.rotateBefore),
	}, nil
}

// ReconcileResources reconciles the certificates used by the webhook client and the webhook server.
// It also returns the duration after which a certificate rotation should be scheduled.
func (r *CertReconciler) ReconcileResources(ctx context.Context) error {
	var err error
	secretName := viper.GetString(operator.FlagWebhookSecret)
	namespace := operator.GetNamespace()
	updateSecret := true

	secClient := r.Clientset.KubeClient.CoreV1().Secrets(namespace)
	secret, err := secClient.Get(ctx, secretName, metav1.GetOptions{})

	if err != nil {
		if operator.IsDevMode() && kerrors.IsNotFound(err) {
			// if in dev mode and we're using self-signed certs we can deal with the secret not being there
			secret = baseWebhookSecret(namespace)
			updateSecret = false
		} else {
			return errors.Wrap(err, fmt.Sprintf("failed to get webhook certificate secret %s", secretName))
		}
	}

	serverCA := certs.BuildCAFromSecret(*secret)

	// check if we need to renew the certificates used in the resources
	if serverCA.ShouldRenew(r.rotateBefore) {
		r.GetLog().Info("Creating new mutating webhook certificates",
			"secret_namespace", secret.Namespace,
			"secret_name", secret.Name,
		)

		ca, err := certs.CreateSelfSignedCA()
		if err != nil {
			return errors.Wrap(err, "unable to set up webhook CA")
		}
		// update the cert secret
		ca.PopulateSecret(secret)

		if updateSecret {
			if _, err := secClient.Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
				return err
			}
		} else {
			if _, err := secClient.Create(ctx, secret, metav1.CreateOptions{}); err != nil {
				return err
			}
		}
		// refresh the local copy of the certs
		serverCA = certs.BuildCAFromSecret(*secret)
	}

	if operator.IsDevMode() {
		// In dev mode (i.e. outside of k8s) write the certs to the local cert dir
		if err = r.writeLocalCertsFromSecret(secret); err != nil {
			return err
		}
	}

	if r.shouldRenewWebhookConfigs(ctx, serverCA) {
		m := createMutatingWebhookWithCABundle(operator.GetNamespace(), secret.Data[operator.CertFileName])
		if err = installMutatingWebhook(ctx, r.Clientset, m); err != nil {
			return err
		}
		v := createValidatingWebhookWithCABundle(operator.GetNamespace(), secret.Data[operator.CertFileName])
		if err = installValidatingWebhook(ctx, r.Clientset, v); err != nil {
			return err
		}
	}

	return nil
}

func (r *CertReconciler) writeLocalCerts(ctx context.Context) error {
	if !operator.IsDevMode() {
		return nil
	}
	secretName := viper.GetString(operator.FlagWebhookSecret)
	namespace := operator.GetNamespace()
	secret, err := r.Clientset.KubeClient.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get secret %s", secretName))
	}
	return r.writeLocalCertsFromSecret(secret)
}

func (r *CertReconciler) writeLocalCertsFromSecret(secret *corev1.Secret) error {
	certDir := operator.GetWebhookCertDir()
	_, err := os.Stat(certDir)
	if err != nil {
		if err = os.MkdirAll(certDir, os.ModePerm); err != nil {
			return errors.Wrap(err, "creating local cert directory")
		}
	}

	for name, data := range secret.Data {
		fn := filepath.Join(certDir, name)
		if err = ioutil.WriteFile(fn, data, os.ModePerm); err != nil {
			return errors.Wrap(err, "writing local cert file")
		}
	}
	return nil
}

func (r *CertReconciler) shouldRenewWebhookConfigs(ctx context.Context, ca *certs.CA) bool {
	// Read the current certificate used by the server
	mCfg, err := r.Clientset.KubeClient.AdmissionregistrationV1().
		MutatingWebhookConfigurations().Get(ctx, viper.GetString(operator.FlagMutatingWebhookName), metav1.GetOptions{})
	if err != nil {
		// probably does not exists so needs creating
		return true
	}

	expectedType := viper.GetString(operator.FlagCertType)
	certType, found := mCfg.Annotations[certTypeAnnotation]
	if !found || certType != expectedType {
		return true
	}

	vCfg, err := r.Clientset.KubeClient.AdmissionregistrationV1().
		ValidatingWebhookConfigurations().Get(ctx, viper.GetString(operator.FlagValidatingWebhookName), metav1.GetOptions{})
	if err != nil {
		// probably does not exists so needs creating
		return true
	}

	certType, found = vCfg.Annotations[certTypeAnnotation]
	if !found || certType != expectedType {
		return true
	}

	// Read the certificate in the mutating webhook configuration
	for _, webhook := range mCfg.Webhooks {
		caBytes := webhook.ClientConfig.CABundle
		if len(caBytes) == 0 || !bytes.Equal(caBytes, ca.Cert.Raw) {
			return true
		}
	}

	// Read the certificate in the mutating webhook configuration
	for _, webhook := range vCfg.Webhooks {
		caBytes := webhook.ClientConfig.CABundle
		if len(caBytes) == 0 || !bytes.Equal(caBytes, ca.Cert.Raw) {
			return true
		}
	}

	return false
}

func (r *CertReconciler) Cleanup() {
	r.GetLog().Info("cleaning up")
	if r.hookInstaller != nil {
		if err := r.hookInstaller.uninstallWebHook(); err != nil {
			r.GetLog().Error(err, "error cleaning up")
		}
	}
}
