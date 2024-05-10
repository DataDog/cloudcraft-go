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

const _testTeamDataPath string = "tests/data/team"

func TestTeamService_List(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testTeamDataPath, "list-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testTeamDataPath, "list-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    []*cloudcraft.Team
		wantErr bool
	}{
		{
			name: "Valid team data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: []*cloudcraft.Team{
				{
					ID:                  "2d2f9b4d-53ac-463f-9c7a-97249a22a7fa",
					Name:                "Team Cloudcraft SDKs",
					Visible:             true,
					CrossOrganizational: false,
					CustomerID:          "e6d4ccd3-551e-425f-b415-637f9005ed2a",
					UpdatedAt:           xtesting.ParseTime(t, "2022-10-10T17:01:11.214Z"),
					CreatedAt:           xtesting.ParseTime(t, "2022-10-10T17:01:11.214Z"),
					ExternalSharing:     true,
					Role:                "owner",
					Members: []cloudcraft.Members{
						{
							ID:         "aa54884d-5c7b-4586-8139-7a6e39231cc9",
							Role:       "owner",
							UserID:     nil,
							Name:       nil,
							Email:      "hello@example.com",
							MFAEnabled: true,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid team data",
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

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			endpoint, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			client := xtesting.SetupMockClient(t, endpoint)

			got, _, err := client.Team.List(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AzureTeam.List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AzureTeam.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
