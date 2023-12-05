// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft_test

import (
	"errors"
	"testing"

	"github.com/DataDog/cloudcraft-go"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *cloudcraft.Config
		want error
	}{
		{
			name: "Valid configuration",
			give: cloudcraft.NewConfig("not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd="),
			want: nil,
		},
		{
			name: "Invalid configuration, missing API key",
			give: &cloudcraft.Config{
				Scheme:  cloudcraft.DefaultScheme,
				Host:    cloudcraft.DefaultHost,
				Port:    cloudcraft.DefaultPort,
				Path:    cloudcraft.DefaultPath,
				Timeout: cloudcraft.DefaultTimeout,
			},
			want: cloudcraft.ErrMissingKey,
		},
		{
			name: "Invalid configuration, invalid API key length",
			give: cloudcraft.NewConfig("short_key"),
			want: cloudcraft.ErrInvalidKey,
		},
		{
			name: "Invalid configuration, invalid endpoint",
			give: &cloudcraft.Config{
				Scheme:  "ftp",
				Host:    cloudcraft.DefaultHost,
				Port:    cloudcraft.DefaultPort,
				Path:    cloudcraft.DefaultPath,
				Timeout: cloudcraft.DefaultTimeout,
				Key:     "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
			},
			want: cloudcraft.ErrInvalidConfig,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := cloudcraft.NewClient(tt.give)

			if !errors.Is(err, tt.want) {
				t.Fatalf("NewClient() error = %v, want %v", err, tt.want)
			}

			if tt.want == nil && client == nil {
				t.Error("Expected non-nil client, got nil")
			}
		})
	}
}

func TestNewClientWithNilConfig(t *testing.T) { //nolint:paralleltest // t.Setenv is not thread-safe
	// Setting environment variables required for NewConfigFromEnv.
	t.Setenv("CLOUDCRAFT_PROTOCOL", cloudcraft.DefaultScheme)
	t.Setenv("CLOUDCRAFT_HOST", cloudcraft.DefaultHost)
	t.Setenv("CLOUDCRAFT_PORT", cloudcraft.DefaultPort)
	t.Setenv("CLOUDCRAFT_PATH", cloudcraft.DefaultPath)
	t.Setenv("CLOUDCRAFT_API_KEY", "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=")
	t.Setenv("CLOUDCRAFT_TIMEOUT", "80s")

	client, err := cloudcraft.NewClient(nil)
	if err != nil {
		t.Fatalf("Unexpected error for nil config: %v", err)
	}

	if client == nil {
		t.Error("Expected non-nil client, got nil")
	}
}
