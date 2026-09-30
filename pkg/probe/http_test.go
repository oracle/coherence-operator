/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package probe

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type fakeHTTPClient struct {
	response *http.Response
	err      error
}

func (f fakeHTTPClient) Do(*http.Request) (*http.Response, error) {
	return f.response, f.err
}

type timeoutError struct{}

func (timeoutError) Error() string { return "timed out" }
func (timeoutError) Timeout() bool { return true }

func TestDoHTTPProbePreservesFailureClassification(t *testing.T) {
	target, err := url.Parse("http://127.0.0.1/test")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("rejected response", func(t *testing.T) {
		client := fakeHTTPClient{response: &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       io.NopCloser(strings.NewReader("unavailable")),
		}}
		result, _, err := DoHTTPProbe(target, nil, client)
		failure := &RequestError{}
		if result != Failure || !errors.As(err, &failure) || failure.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("result=%q err=%v", result, err)
		}
	})

	t.Run("transport timeout", func(t *testing.T) {
		result, _, err := DoHTTPProbe(target, nil, fakeHTTPClient{err: timeoutError{}})
		failure := &RequestError{}
		if result != Failure || !errors.As(err, &failure) || failure.Kind != "transport timeout" {
			t.Fatalf("result=%q err=%v", result, err)
		}
	})
}
