// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/DataDog/cloudcraft-go"
	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

const _testUserDataPath string = "tests/data/user"

func TestUserService_Me(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testUserDataPath, "me-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testUserDataPath, "me-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    *cloudcraft.User
		wantErr bool
	}{
		{
			name: "Valid user data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: &cloudcraft.User{
				ID:    "b92570ba-8969-4e41-b6a3-3d672b44f9f5",
				Name:  "Go SDK",
				Email: "hi@example.com",
				Settings: map[string]any{
					"currency":  "USD",
					"firstTime": false,
				},
				CreatedAt:  xtesting.ParseTime(t, "2022-10-10T16:52:40.771Z"),
				UpdatedAt:  xtesting.ParseTime(t, "2023-11-08T14:44:28.872Z"),
				AccessedAt: xtesting.ParseTime(t, "2023-11-08T14:44:28.872Z"),
			},
			wantErr: false,
		},
		{
			name: "Invalid user data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
			},
			context: ctx,
			want:    nil,
			wantErr: true,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context: ctx,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: nil,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.handler)
			defer server.Close()

			endpoint, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, _, err := client.User.Me(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UserService.Me() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("UserService.Me() = %v, want %v", got, tt.want)
			}
		})
	}
}
