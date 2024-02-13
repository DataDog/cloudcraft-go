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

const _testBlueprintDataPath string = "tests/data/blueprint"

func TestBlueprintService_List(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "list-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "generic-invalid.json"))
		emptyTestData   = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "list-empty.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    []*cloudcraft.Blueprint
		wantErr bool
	}{
		{
			name: "Valid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: []*cloudcraft.Blueprint{
				{
					ID:               "60ec30c9-741f-4acb-b5f5-794934987802",
					Name:             "Web App Reference Architecture",
					Tags:             nil,
					ReadAccess:       nil,
					WriteAccess:      nil,
					CreatedAt:        xtesting.ParseTime(t, "2023-04-01T21:02:10.781Z"),
					UpdatedAt:        xtesting.ParseTime(t, "2023-04-01T21:02:10.781Z"),
					CreatorID:        "9e52d877-4dab-4aa6-95be-c7ba5d685689",
					CurrentVersionID: "60ec30c9-741f-4acb-b5f5-794934987802",
					LastUserID:       "9e52d877-4dab-4aa6-95be-c7ba5d685689",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid blueprint data",
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
			name: "Empty blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(emptyTestData)
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

			got, _, err := client.Blueprint.List(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Blueprint.List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Blueprint.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlueprintService_Get(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "get-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "generic-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		id      string
		want    *cloudcraft.Blueprint
		wantErr bool
	}{
		{
			name: "Valid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			id:      "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			want: &cloudcraft.Blueprint{
				ID:          "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
				Name:        "Test blueprint",
				Tags:        &[]string{},
				ReadAccess:  nil,
				WriteAccess: nil,
				CreatedAt:   xtesting.ParseTime(t, "2023-11-09T23:19:29.611Z"),
				UpdatedAt:   xtesting.ParseTime(t, "2023-11-09T23:19:41.018Z"),
				CreatorID:   "9e52d877-4dab-4aa6-95be-c7ba5d685689",
				CustomerID:  nil,
				Data: &cloudcraft.BlueprintData{
					Grid:  "infinite",
					Name:  "Test blueprint",
					Text:  []map[string]any{},
					Edges: []map[string]any{},
					Icons: []map[string]any{},
					Nodes: []map[string]any{
						{
							"id":           "d801fe26-1f73-49a5-bbe9-23c5fb0888e0",
							"type":         "ec2",
							"mapPos":       []any{float64(-2), float64(11)},
							"region":       "us-east-1",
							"platform":     "linux",
							"transparent":  false,
							"instanceSize": "large",
							"instanceType": "m5",
						},
					},
					Theme: &cloudcraft.Theme{
						Base: "light",
					},
					Groups:     []map[string]any{},
					Images:     []map[string]any{},
					Version:    4,
					Surfaces:   []map[string]any{},
					ShareDocs:  false,
					Connectors: []map[string]any{},
					Projection: "isometric",
					LiveOptions: &cloudcraft.LiveOptions{
						AutoLabel:   true,
						AutoConnect: true,
						ExcludedTypes: []string{
							"ebs",
							"dxconnection",
							"natgateway",
							"internetgateway",
							"vpngateway",
							"customergateway",
						},
						UpdatesEnabled:     true,
						UpdateAllOnScan:    true,
						UpdateGroupsOnScan: true,
						UpdateNodeOnSelect: true,
					},
					DisabledLayers: []string{},
				},
				LastUserID: "9e52d877-4dab-4aa6-95be-c7ba5d685689",
			},
			wantErr: false,
		},
		{
			name: "Invalid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
			},
			context: ctx,
			id:      "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			want:    nil,
			wantErr: true,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context: ctx,
			id:      "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: nil,
			id:      "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Missing ID",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: ctx,
			id:      "",
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

			got, _, err := client.Blueprint.Get(tt.context, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Blueprint.Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Blueprint.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlueprintService_Create(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "create-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "generic-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		give    *cloudcraft.Blueprint
		want    *cloudcraft.Blueprint
		wantErr bool
	}{
		{
			name: "Valid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)

				w.Write(validTestData)
			},
			context: ctx,
			give: &cloudcraft.Blueprint{
				Name: "My new blueprint",
			},
			want: &cloudcraft.Blueprint{
				ID:          "31c014b0-279a-4662-9fd4-3f104a2c4f84",
				Name:        "My new blueprint",
				CreatorID:   "9e52d877-4dab-4aa6-95be-c7ba5d685689",
				Tags:        nil,
				ReadAccess:  nil,
				WriteAccess: nil,
				UpdatedAt:   xtesting.ParseTime(t, "2023-11-14T22:00:39.332Z"),
				CreatedAt:   xtesting.ParseTime(t, "2023-11-14T22:00:39.332Z"),
				CustomerID:  nil,
				Data: &cloudcraft.BlueprintData{
					Name:     "My new blueprint",
					Surfaces: []map[string]any{},
					Version:  4,
				},
				LastUserID: "9e52d877-4dab-4aa6-95be-c7ba5d685689",
			},
			wantErr: false,
		},
		{
			name: "Invalid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)

				w.Write(invalidTestData)
			},
			context: ctx,
			give: &cloudcraft.Blueprint{
				Name: "My new blueprint",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context: ctx,
			give: &cloudcraft.Blueprint{
				Name: "My new blueprint",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: nil,
			give: &cloudcraft.Blueprint{
				Name: "My new blueprint",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil blueprint",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: ctx,
			give:    nil,
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

			got, _, err := client.Blueprint.Create(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Blueprint.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Blueprint.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlueprintService_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		context  context.Context
		give     *cloudcraft.Blueprint
		giveEtag string
		want     *cloudcraft.Response
		wantErr  bool
	}{
		{
			name: "Valid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give: &cloudcraft.Blueprint{
				ID:   "31c014b0-279a-4662-9fd4-3f104a2c4f84",
				Name: "My updated blueprint",
			},
			giveEtag: `W/"31c014b0-279a-4662-9fd4-3f104a2c4f84"`,
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
			name: "Valid blueprint data without etag",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give: &cloudcraft.Blueprint{
				ID:   "31c014b0-279a-4662-9fd4-3f104a2c4f84",
				Name: "My updated blueprint",
			},
			giveEtag: "",
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
			name:    "Nil context",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: nil,
			give: &cloudcraft.Blueprint{
				ID:   "31c014b0-279a-4662-9fd4-3f104a2c4f84",
				Name: "My updated blueprint",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context: ctx,
			give: &cloudcraft.Blueprint{
				ID:   "31c014b0-279a-4662-9fd4-3f104a2c4f84",
				Name: "My updated blueprint",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil blueprint",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: ctx,
			give:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Missing blueprint ID",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: ctx,
			give: &cloudcraft.Blueprint{
				Name: "My updated blueprint",
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

			got, err := client.Blueprint.Update(tt.context, tt.give, tt.giveEtag)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Blueprint.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Blueprint.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlueprintService_Delete(t *testing.T) {
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
			name: "Valid blueprint data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give:    "31c014b0-279a-4662-9fd4-3f104a2c4f84",
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
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context: ctx,
			give:    "31c014b0-279a-4662-9fd4-3f104a2c4f84",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
			context: nil,
			give:    "31c014b0-279a-4662-9fd4-3f104a2c4f84",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Missing blueprint ID",
			handler: func(_ http.ResponseWriter, _ *http.Request) {},
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

			got, err := client.Blueprint.Delete(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Blueprint.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Blueprint.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlueprintService_ExportImages(t *testing.T) {
	t.Parallel()

	var (
		validTestData = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "export-image-valid.png"))
		ctx           = context.Background()
	)

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		context    context.Context
		giveID     string
		giveFormat string
		giveParams *cloudcraft.ImageExportParams
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name: "Valid blueprint export",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "png",
			giveParams: &cloudcraft.ImageExportParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "png",
			giveParams: &cloudcraft.ImageExportParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name:       "Nil context",
			handler:    func(_ http.ResponseWriter, _ *http.Request) {},
			context:    nil,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "png",
			giveParams: &cloudcraft.ImageExportParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name: "Nil image params",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "png",
			giveParams: nil,
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
		{
			name:       "Missing blueprint ID",
			handler:    func(_ http.ResponseWriter, _ *http.Request) {},
			context:    ctx,
			giveID:     "",
			giveFormat: "png",
			giveParams: &cloudcraft.ImageExportParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    true,
		},
		{
			name: "Missing image format",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "",
			giveParams: &cloudcraft.ImageExportParams{
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

			got, _, err := client.Blueprint.ExportImage(tt.context, tt.giveID, tt.giveFormat, tt.giveParams)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Blueprint.Export() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				gotPNG, err := png.DecodeConfig(bytes.NewReader(got))
				if err != nil {
					t.Fatal(err)
				}

				if gotPNG.Width != tt.wantWidth {
					t.Fatalf("Blueprint.Export() width = %v, want %v", gotPNG.Width, tt.wantWidth)
				}

				if gotPNG.Height != tt.wantHeight {
					t.Fatalf("Blueprint.Export() height = %v, want %v", gotPNG.Height, tt.wantHeight)
				}

				if !bytes.Equal(got, validTestData) {
					t.Fatalf("Blueprint.Export() = %v, want %v", got, validTestData)
				}
			}
		})
	}
}

func TestBlueprintService_ExportBudget(t *testing.T) {
	t.Parallel()

	var (
		validTestData = xtesting.ReadFile(t, filepath.Join(_testBlueprintDataPath, "export-budget-valid.csv"))
		ctx           = context.Background()
	)

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		context    context.Context
		giveID     string
		giveFormat string
		giveParams *cloudcraft.BudgetExportParams
		wantSize   int
		wantErr    bool
	}{
		{
			name: "Valid budget data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "csv",
			giveParams: &cloudcraft.BudgetExportParams{
				Currency: "USD",
				Period:   "month",
				Rate:     "monthly",
			},
			wantSize: 308,
			wantErr:  false,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "csv",
			giveParams: &cloudcraft.BudgetExportParams{
				Currency: "USD",
				Period:   "month",
				Rate:     "monthly",
			},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:       "Nil context",
			handler:    func(_ http.ResponseWriter, _ *http.Request) {},
			context:    nil,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "csv",
			giveParams: &cloudcraft.BudgetExportParams{
				Currency: "USD",
				Period:   "month",
				Rate:     "monthly",
			},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name: "Nil budget params",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "csv",
			giveParams: nil,
			wantSize:   308,
			wantErr:    false,
		},
		{
			name:       "Missing blueprint ID",
			handler:    func(_ http.ResponseWriter, _ *http.Request) {},
			context:    ctx,
			giveID:     "",
			giveFormat: "csv",
			giveParams: &cloudcraft.BudgetExportParams{
				Currency: "USD",
				Period:   "month",
				Rate:     "monthly",
			},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name: "Missing budget format",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "0f1a4e20-a887-4467-a37b-1bc7a3deb9a9",
			giveFormat: "",
			giveParams: &cloudcraft.BudgetExportParams{
				Currency: "USD",
				Period:   "month",
				Rate:     "monthly",
			},
			wantSize: 308,
			wantErr:  false,
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

			got, _, err := client.Blueprint.ExportBudget(tt.context, tt.giveID, tt.giveFormat, tt.giveParams)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BlueprintService.ExportBudget() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.wantSize > 0 && len(got) != tt.wantSize {
				t.Fatalf("BlueprintService.ExportBudget() length = %v, want %v", len(got), tt.wantSize)
			}

			if !tt.wantErr && tt.wantSize > 0 && !bytes.Equal(got, validTestData) {
				t.Fatalf("BlueprintService.ExportBudget() data differs from valid test data")
			}
		})
	}
}
