// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft_test

import (
	"testing"
	"time"

	"github.com/DataDog/cloudcraft-go"
)

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    cloudcraft.Config
		wantErr bool
	}{
		{
			name: "Valid configuration",
			give: cloudcraft.Config{
				Scheme:  "https",
				Host:    "api.example.com",
				Port:    "443",
				Path:    "/",
				Key:     "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
				Timeout: time.Second * 80,
			},
			wantErr: false,
		},
		{
			name: "Missing scheme",
			give: cloudcraft.Config{
				Host: "api.example.com",
				Key:  "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
			},
			wantErr: true,
		},
		{
			name: "Missing host",
			give: cloudcraft.Config{
				Scheme: "https",
				Key:    "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
			},
			wantErr: true,
		},
		{
			name: "Missing key",
			give: cloudcraft.Config{
				Scheme: "https",
				Host:   "api.example.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid key length",
			give: cloudcraft.Config{
				Scheme: "https",
				Host:   "api.example.com",
				Key:    "shortkey",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
