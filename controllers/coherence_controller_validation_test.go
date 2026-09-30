/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"strings"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestValidateHealthMutatorLauncherConfiguration(t *testing.T) {
	selector := corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: "health-secret"},
		Key:                  "value",
	}
	secureConnection := &coh.HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &coh.HTTPBasicAuth{Username: selector, Password: selector},
		TLS:       &coh.HTTPClientTLS{CASecret: selector, ServerName: "health.example"},
		ServerTLS: &coh.SSLSpec{
			Enabled:              ptr.To(true),
			KeyStore:             ptr.To("server.p12"),
			KeyStorePasswordFile: ptr.To("store-password"),
			KeyPasswordFile:      ptr.To("key-password"),
		},
	}

	tests := []struct {
		name        string
		connection  *coh.HealthMutatorConnection
		application *coh.ApplicationSpec
		wantError   string
	}{
		{
			name:       "secure normal operator launcher",
			connection: secureConnection,
		},
		{
			name:       "secure image entry point with default JDK options",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				UseImageEntryPoint: ptr.To(true),
			},
		},
		{
			name:       "secure image entry point with explicit JDK options",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				UseImageEntryPoint: ptr.To(true),
				UseJdkJavaOptions:  ptr.To(true),
			},
		},
		{
			name:       "secure image entry point without JVM argument delivery",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				UseImageEntryPoint: ptr.To(true),
				UseJdkJavaOptions:  ptr.To(false),
			},
			wantError: "requires application.useJdkJavaOptions=true",
		},
		{
			name:       "secure image entry point with empty alternate JVM options",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				UseImageEntryPoint:      ptr.To(true),
				UseJdkJavaOptions:       ptr.To(false),
				AlternateJdkJavaOptions: ptr.To(""),
			},
			wantError: "non-empty application.alternateJdkJavaOptions",
		},
		{
			name:       "secure image entry point with alternate JVM options",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				UseImageEntryPoint:      ptr.To(true),
				UseJdkJavaOptions:       ptr.To(false),
				AlternateJdkJavaOptions: ptr.To("ALT_JAVA_OPTS"),
			},
		},
		{
			name:       "secure custom entry point without image entry point delivery",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				EntryPoint: []string{"/application/start"},
			},
			wantError: "application.entryPoint requires application.useImageEntryPoint=true",
		},
		{
			name:       "secure custom entry point with default JDK options delivery",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				EntryPoint:         []string{"/application/start"},
				UseImageEntryPoint: ptr.To(true),
			},
		},
		{
			name:       "secure custom entry point with alternate JVM options delivery",
			connection: secureConnection,
			application: &coh.ApplicationSpec{
				EntryPoint:              []string{"/application/start"},
				UseImageEntryPoint:      ptr.To(true),
				UseJdkJavaOptions:       ptr.To(false),
				AlternateJdkJavaOptions: ptr.To("ALT_JAVA_OPTS"),
			},
		},
		{
			name: "omitted managed health leaves image settings unchanged",
			application: &coh.ApplicationSpec{
				UseImageEntryPoint: ptr.To(true),
				UseJdkJavaOptions:  ptr.To(false),
			},
		},
		{
			name:       "insecure managed health leaves custom entry point unchanged",
			connection: &coh.HealthMutatorConnection{Secure: false},
			application: &coh.ApplicationSpec{
				EntryPoint: []string{"/application/start"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployment := &coh.Coherence{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
			}
			deployment.Spec.HealthMutatorConnection = tt.connection
			deployment.Spec.Application = tt.application

			err := validateHealthMutatorConfiguration(deployment)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("validateHealthMutatorConfiguration() returned an error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("validateHealthMutatorConfiguration() error = %v, want containing %q", err, tt.wantError)
			}
		})
	}
}
