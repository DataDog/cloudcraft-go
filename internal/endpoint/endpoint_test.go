// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package endpoint_test

import (
	"net/url"
	"testing"

	"github.com/DataDog/cloudcraft-go/internal/endpoint"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		scheme    string
		host      string
		port      string
		path      string
		want      *url.URL
		wantError bool
	}{
		{
			name:      "Valid URL",
			scheme:    "https",
			host:      "example.com",
			port:      "8080",
			path:      "",
			want:      &url.URL{Scheme: "https", Host: "example.com:8080", Path: "/"},
			wantError: false,
		},
		{
			name:      "Invalid URL",
			scheme:    "https",
			host:      "example.com",
			port:      "\n", // Invalid character in port
			path:      "/test",
			want:      nil,
			wantError: true,
		},
		{
			name:      "Invalid URL scheme",
			scheme:    "ftp",
			host:      "example.com",
			port:      "8080",
			path:      "/test",
			want:      nil,
			wantError: true,
		},
		{
			name:      "Missing Scheme",
			scheme:    "",
			host:      "example.com",
			port:      "8080",
			path:      "/test",
			want:      nil,
			wantError: true,
		},
		{
			name:      "Missing Host",
			scheme:    "https",
			host:      "",
			port:      "8080",
			path:      "/test",
			want:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := endpoint.Parse(tt.scheme, tt.host, tt.port, tt.path)

			if (err != nil) != tt.wantError {
				t.Errorf("Expected error? %v, got: %v", tt.wantError, err)
			}

			if tt.want != nil && got.String() != tt.want.String() {
				t.Errorf("Expected URL: %v, got: %v", tt.want, got)
			}
		})
	}
}
