/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package v1

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

// ValidateAction validates the new Operator HTTP action while preserving the
// historical behavior of legacy probe actions.
func (in *Probe) ValidateAction() error {
	if in == nil {
		return fmt.Errorf("probe action is required")
	}
	a := in.HTTP
	if a == nil {
		return nil
	}
	if in.Exec != nil || in.HTTPGet != nil || in.TCPSocket != nil || in.GRPC != nil {
		return fmt.Errorf("HTTP cannot be combined with a legacy probe action")
	}
	if a.Method != "" && a.Method != "GET" && a.Method != "PUT" {
		return fmt.Errorf("HTTP method must be GET or PUT")
	}
	if a.Scheme != "" && a.Scheme != corev1.URISchemeHTTP && a.Scheme != corev1.URISchemeHTTPS {
		return fmt.Errorf("unsupported HTTP scheme")
	}
	if strings.ContainsAny(a.Host, "/@?#") || (!strings.HasPrefix(a.Path, "/") && a.Path != "") || strings.HasPrefix(a.Path, "//") {
		return fmt.Errorf("invalid HTTP endpoint")
	}
	for _, h := range a.HTTPHeaders {
		if strings.EqualFold(h.Name, "Authorization") && a.Scheme != corev1.URISchemeHTTPS {
			return fmt.Errorf("authorization requires HTTPS")
		}
	}
	valid := func(s corev1.SecretKeySelector) bool {
		return s.Name != "" && s.Key != "" && (s.Optional == nil || !*s.Optional)
	}
	if a.BasicAuth != nil {
		if a.Host != "" {
			return fmt.Errorf("basic authentication requires the target Pod address; host must be empty")
		}
		if a.Scheme != corev1.URISchemeHTTPS {
			return fmt.Errorf("basic authentication requires HTTPS")
		}
		if !valid(a.BasicAuth.Username) || !valid(a.BasicAuth.Password) {
			return fmt.Errorf("required credential Secret selectors are incomplete")
		}
		for _, h := range a.HTTPHeaders {
			if strings.EqualFold(h.Name, "Authorization") {
				return fmt.Errorf("authorization conflicts with basicAuth")
			}
		}
	}
	if a.TLS != nil && (a.Scheme != corev1.URISchemeHTTPS || !valid(a.TLS.CASecret)) {
		return fmt.Errorf("TLS requires HTTPS and a required CA Secret selector")
	}
	return nil
}
