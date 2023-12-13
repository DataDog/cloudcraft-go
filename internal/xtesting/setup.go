// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xtesting

import (
	"net/url"
	"testing"

	"github.com/DataDog/cloudcraft-go"
)

const _envAPIKey string = "CLOUDCRAFT_TEST_API_KEY"

// SetupMockClient sets up a test API client for unit tests against a mock
// version of the Cloudcraft API.
func SetupMockClient(t *testing.T, endpoint *url.URL) *cloudcraft.Client {
	t.Helper()

	cfg := &cloudcraft.Config{
		Scheme: endpoint.Scheme,
		Host:   endpoint.Hostname(),
		Port:   endpoint.Port(),
		Path:   cloudcraft.DefaultPath,
		Key:    "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
	}

	client, err := cloudcraft.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client for mock tests: %v", err)
	}

	return client
}

// SetupLiveClient sets up a test API client for unit tests against the live
// Cloudcraft API.
//
// The following environment variables are required:
// - CLOUDCRAFT_TEST_API_KEY
//
// If any of these variables are not set, SetupLiveClient will fail the test.
func SetupLiveClient(t *testing.T) *cloudcraft.Client {
	t.Helper()

	key := GetEnv(t, _envAPIKey)

	cfg := &cloudcraft.Config{
		Scheme: cloudcraft.DefaultScheme,
		Host:   cloudcraft.DefaultHost,
		Port:   cloudcraft.DefaultPort,
		Path:   cloudcraft.DefaultPath,
		Key:    key,
	}

	client, err := cloudcraft.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client for live tests: %v", err)
	}

	return client
}
