// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xhttp_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DataDog/cloudcraft-go/internal/xhttp"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	var (
		testTimeout = 30 * time.Second
		client      = xhttp.NewClient(testTimeout)
	)

	if client == nil {
		t.Fatal("NewClient returned nil, expected *http.Client")
	}

	if client.Timeout != testTimeout {
		t.Errorf("Expected timeout %v, got %v", testTimeout, client.Timeout)
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Transport is not of type *http.Transport")
	}

	if transport.TLSClientConfig.MinVersion != tls.VersionTLS13 {
		t.Errorf("Expected TLS min version %v, got %v", tls.VersionTLS12, transport.TLSClientConfig.MinVersion)
	}

	if transport.MaxIdleConns != xhttp.DefaultMaxIddleConns {
		t.Errorf("Expected MaxIdleConns %d, got %d", xhttp.DefaultMaxIddleConns, transport.MaxIdleConns)
	}

	if transport.MaxIdleConnsPerHost != xhttp.DefaultMaxIddleConnsPerHost {
		t.Errorf("Expected MaxIdleConnsPerHost %d, got %d", xhttp.DefaultMaxIddleConnsPerHost, transport.MaxIdleConnsPerHost)
	}

	if !transport.ForceAttemptHTTP2 {
		t.Error("Expected ForceAttemptHTTP2 to be true")
	}
}

func TestNewClient_CheckRedirect(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.com", http.StatusFound)
	}))
	defer ts.Close()

	var (
		testTimeout = 30 * time.Second
		client      = xhttp.NewClient(testTimeout)
	)

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("Expected status code %d, got %d", http.StatusFound, resp.StatusCode)
	}
}
