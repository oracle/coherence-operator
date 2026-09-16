/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package servicemonitor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/utils"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/rest"
)

type serviceMonitorStorage struct {
	utils.Storage
	previous, latest coh.Resources
}

func (s serviceMonitorStorage) GetPrevious() coh.Resources { return s.previous }
func (s serviceMonitorStorage) GetLatest() coh.Resources   { return s.latest }
func (s serviceMonitorStorage) GetHash() (string, bool)    { return "new", true }

func TestServiceMonitorUpdateLabelDrop(t *testing.T) {
	for _, fallback := range []bool{false, true} {
		name := "strategic patch"
		if fallback {
			name = "full spec fallback"
		}
		t.Run(name, func(t *testing.T) {
			enabled, disabled := true, false
			port := coh.NamedPortSpec{Name: "metrics", ServiceMonitor: &coh.ServiceMonitorSpec{Enabled: &enabled}}
			parent := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}
			live := port.CreateServiceMonitor(parent)
			live.TypeMeta = metav1.TypeMeta{APIVersion: "monitoring.coreos.com/v1", Kind: "ServiceMonitor"}
			patches, updates := 0, 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPatch:
					patches++
					if fallback {
						http.Error(w, `{"kind":"Status","apiVersion":"v1","status":"Failure","reason":"UnsupportedMediaType","code":415}`, http.StatusUnsupportedMediaType)
						return
					}
					data, err := io.ReadAll(r.Body)
					if err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					original, err := json.Marshal(live)
					if err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					result, err := strategicpatch.StrategicMergePatch(original, data, monitoring.ServiceMonitor{})
					if err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					// Decode into a fresh object so omitted fields cannot retain old values.
					next := &monitoring.ServiceMonitor{}
					if err = json.Unmarshal(result, next); err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					live = next
				case http.MethodPut:
					updates++
					next := &monitoring.ServiceMonitor{}
					if err := json.NewDecoder(r.Body).Decode(next); err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					live = next
				case http.MethodGet:
				default:
					t.Errorf("unexpected method %s", r.Method)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if err := json.NewEncoder(w).Encode(live); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client, err := monclient.NewForConfig(&rest.Config{Host: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			reconciler := &ReconcileServiceMonitor{monClient: client}
			reconciler.Kind = coh.ResourceTypeServiceMonitor
			reconciler.SetCommonReconciler("test", nil, clients.ClientSet{})
			narrow := []monitoring.RelabelConfig{{Action: "labeldrop", Regex: "(endpoint|service)"}}
			for _, state := range []struct {
				flag  *bool
				rules []monitoring.RelabelConfig
			}{{&disabled, nil}, {&disabled, narrow}, {nil, narrow}} {
				previous := live.DeepCopy()
				port.ServiceMonitor.UseDefaultLabelDrop = state.flag
				port.ServiceMonitor.Relabelings = state.rules
				desired := port.CreateServiceMonitor(parent)
				desired.TypeMeta = live.TypeMeta
				storage := serviceMonitorStorage{previous: coh.Resources{Items: []coh.Resource{{Kind: coh.ResourceTypeServiceMonitor, Name: live.Name, Spec: previous}}}, latest: coh.Resources{Items: []coh.Resource{{Kind: coh.ResourceTypeServiceMonitor, Name: live.Name, Spec: desired}}}}
				if err = reconciler.UpdateServiceMonitor(context.Background(), live.Namespace, live.Name, live.DeepCopy(), storage, logr.Discard()); err != nil {
					t.Fatal(err)
				}
				got, err := client.ServiceMonitors(live.Namespace).Get(context.Background(), live.Name, metav1.GetOptions{})
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got.Spec, desired.Spec) {
					t.Fatalf("live spec does not match desired: got %#v want %#v", got.Spec, desired.Spec)
				}
			}
			if patches != 3 {
				t.Fatalf("expected three patches, got %d", patches)
			}
			if fallback && updates != 3 || !fallback && updates != 0 {
				t.Fatalf("unexpected update count %d", updates)
			}
		})
	}
}
