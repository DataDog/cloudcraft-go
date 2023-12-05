// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft_test

import (
	"bytes"
	"context"
	"image/png"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/DataDog/cloudcraft-go"
	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

const _testAzureDataPath string = "tests/data/azure"

func TestAzureService_List(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testAzureDataPath, "list-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testAzureDataPath, "generic-invalid.json"))
		emptyTestData   = xtesting.ReadFile(t, filepath.Join(_testAzureDataPath, "list-empty.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    []*cloudcraft.AzureAccount
		wantErr bool
	}{
		{
			name: "Valid Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: []*cloudcraft.AzureAccount{
				{
					ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
					Name:           "Go SDK Test",
					ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
					DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
					SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
					ReadAccess:     &[]string{},
					WriteAccess:    &[]string{},
					CreatedAt:      xtesting.ParseTime(t, "2023-03-15T20:42:52.704Z"),
					UpdatedAt:      xtesting.ParseTime(t, "2023-03-15T20:43:10.171Z"),
					CreatorID:      "6935c7da-cdfb-4885-902c-25aa00720ab4",
					Hint:           "3RK",
					Source:         "azure",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
			},
			context: ctx,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(emptyTestData)
			},
			context: ctx,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: nil,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			endpoint, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, _, err := client.Azure.List(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AzureService.List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AzureService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAzureService_Create(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testAzureDataPath, "create-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testAzureDataPath, "generic-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		give    *cloudcraft.AzureAccount
		want    *cloudcraft.AzureAccount
		wantErr bool
	}{
		{
			name: "Valid Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)

				w.Write(validTestData)
			},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want: &cloudcraft.AzureAccount{
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ReadAccess:     nil,
				WriteAccess:    nil,
				CreatorID:      "6935c7da-cdfb-4885-902c-25aa00720ab4",
				UpdatedAt:      xtesting.ParseTime(t, "2023-11-20T22:11:43.688Z"),
				CreatedAt:      xtesting.ParseTime(t, "2023-11-20T22:11:43.688Z"),
				CustomerID:     nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)

				w.Write(invalidTestData)
			},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "API error response",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: nil,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil Azure account",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty name",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty application ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty directory ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty subscription ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty client secret",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			endpoint, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, _, err := client.Azure.Create(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAzureService_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		give    *cloudcraft.AzureAccount
		want    *cloudcraft.Response
		wantErr bool
	}{
		{
			name: "Valid Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want: &cloudcraft.Response{
				Header: http.Header{
					"Date": []string{
						time.Now().In(time.UTC).Format(http.TimeFormat),
					},
				},
				Body:   []uint8{},
				Status: http.StatusNoContent,
			},
			wantErr: false,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil Azure account",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: nil,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty name",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty Application ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty Directory ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty Subscription ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "",
				ClientSecret:   "tV>0}(,[g91|V5mV|:>~rC841E7}[~n9~Wt4;H%II4",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty Client Secret",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AzureAccount{
				ID:             "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
				Name:           "Go SDK Test",
				ApplicationID:  "3a64bc23-5dd6-4624-8ce8-fe3e61b41579",
				DirectoryID:    "5d7ef62e-c8bb-41fc-9a55-9a2c30701027",
				SubscriptionID: "db0297eb-ad6c-4e63-86b0-c1acb6a16570",
				ClientSecret:   "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			endpoint, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, err := client.Azure.Update(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Azure.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Azure.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAzureService_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		give    string
		want    *cloudcraft.Response
		wantErr bool
	}{
		{
			name: "Valid Azure account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give:    "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			want: &cloudcraft.Response{
				Header: http.Header{
					"Date": []string{
						time.Now().In(time.UTC).Format(http.TimeFormat),
					},
				},
				Body:   []uint8{},
				Status: http.StatusNoContent,
			},
			wantErr: false,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			context: ctx,
			give:    "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			give:    "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give:    "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			endpoint, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, err := client.Azure.Delete(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Azure.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Azure.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAzureService_Snapshot(t *testing.T) {
	t.Parallel()

	var (
		validTestData = xtesting.ReadFile(t, filepath.Join(_testAzureDataPath, "snapshot-valid.png"))
		ctx           = context.Background()
	)

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		context    context.Context
		giveID     string
		giveRegion string
		giveFormat string
		giveParams *cloudcraft.SnapshotParams
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name: "Valid Azure account snapshot",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			giveRegion: "centralus",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			context:    ctx,
			giveID:     "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			giveRegion: "centralus",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name:       "Nil context",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			giveID:     "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			giveRegion: "centralus",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name: "Nil snapshot params",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			giveRegion: "centralus",
			giveFormat: "png",
			giveParams: nil,
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
		{
			name:       "Empty ID",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			context:    ctx,
			giveID:     "",
			giveRegion: "centralus",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name:       "Empty region",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			context:    ctx,
			giveID:     "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			giveRegion: "",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name: "Empty format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "4349ccdb-a2fd-4a89-a07b-48e3e330670b",
			giveRegion: "centralus",
			giveFormat: "",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			endpoint, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, _, err := client.Azure.Snapshot(tt.context, tt.giveID, tt.giveRegion, tt.giveFormat, tt.giveParams)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AzureService.Snapshot() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				gotPNG, err := png.DecodeConfig(bytes.NewReader(got))
				if err != nil {
					t.Fatal(err)
				}

				if gotPNG.Width != tt.wantWidth {
					t.Fatalf("Azure.Snapshot() width = %v, want %v", gotPNG.Width, tt.wantWidth)
				}

				if gotPNG.Height != tt.wantHeight {
					t.Fatalf("Azure.Snapshot() height = %v, want %v", gotPNG.Height, tt.wantHeight)
				}

				if !bytes.Equal(got, validTestData) {
					t.Fatalf("Azure.Snapshot() = %v, want %v", got, validTestData)
				}
			}
		})
	}
}
