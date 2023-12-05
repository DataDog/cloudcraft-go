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

const _testAWSDataPath string = "tests/data/aws"

func TestAWSService_List(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "list-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "generic-invalid.json"))
		emptyTestData   = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "list-empty.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    []*cloudcraft.AWSAccount
		wantErr bool
	}{
		{
			name: "Valid AWS account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: []*cloudcraft.AWSAccount{
				{
					ID:         "47830a91-51b7-4708-b9b2-5f3d121fc039",
					Name:       "Go SDK Test",
					RoleARN:    "arn:aws:iam::558791803304:role/cloudcraft",
					ExternalID: "61fc01d6-3e6f-47ab-bc44-53fab97c217a",
					ReadAccess: &[]string{
						"team/5f209338-50a1-495f-90dd-73251dec7329",
						"team/d7cd0211-85a7-45fc-9292-8d5c62cef70a",
					},
					WriteAccess: &[]string{
						"team/5f209338-50a1-495f-90dd-73251dec7329",
						"team/d7cd0211-85a7-45fc-9292-8d5c62cef70a",
					},
					CreatedAt: xtesting.ParseTime(t, "2019-02-19T16:20:34.042Z"),
					UpdatedAt: xtesting.ParseTime(t, "2022-08-05T18:13:05.625Z"),
					CreatorID: "280ccb78-6a06-4e28-adad-8d16d413be50",
					Source:    "aws",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid AWS account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
			},
			context: ctx,
			want:    nil,
			wantErr: true,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			context: ctx,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty AWS account data",
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

			got, _, err := client.AWS.List(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWSService.List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AWSService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSService_Create(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "create-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "generic-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		give    *cloudcraft.AWSAccount
		want    *cloudcraft.AWSAccount
		wantErr bool
	}{
		{
			name: "Valid AWS account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)

				w.Write(validTestData)
			},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				Name:       "Go SDK Test",
				RoleARN:    "arn:aws:iam::558791803304:role/cloudcraft",
				ExternalID: "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
			},
			want: &cloudcraft.AWSAccount{
				ID:          "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
				Name:        "Go SDK Test",
				RoleARN:     "arn:aws:iam::558791803304:role/cloudcraft",
				ExternalID:  "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
				ReadAccess:  nil,
				WriteAccess: nil,
				CreatedAt:   xtesting.ParseTime(t, "2019-02-19T16:20:34.042Z"),
				UpdatedAt:   xtesting.ParseTime(t, "2022-08-05T18:13:05.625Z"),
				CreatorID:   "17d5fe91-9efb-4b1a-90cd-0b885b1d43b9",
			},
			wantErr: false,
		},
		{
			name: "Invalid AWS account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
			},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				Name:       "Go SDK Test",
				RoleARN:    "arn:aws:iam::643880554691j:role/cloudcraft",
				ExternalID: "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "API error response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				Name:       "Go SDK Test",
				RoleARN:    "arn:aws:iam::643880554691j:role/cloudcraft",
				ExternalID: "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			give: &cloudcraft.AWSAccount{
				Name:       "Go SDK Test",
				RoleARN:    "arn:aws:iam::643880554691j:role/cloudcraft",
				ExternalID: "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil AWS account",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty Name",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				Name:       "",
				RoleARN:    "arn:aws:iam::643880554691j:role/cloudcraft",
				ExternalID: "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty RoleARN",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				Name:       "Go SDK Test",
				RoleARN:    "",
				ExternalID: "8a8a745a-d01f-4541-8ab0-e3558e7c6b1c",
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

			got, _, err := client.AWS.Create(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWSAccount.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AWSAccount.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSService_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		give    *cloudcraft.AWSAccount
		want    *cloudcraft.Response
		wantErr bool
	}{
		{
			name: "Valid AWS account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				ID:      "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
				Name:    "My updated AWS account",
				RoleARN: "arn:aws:iam::558791803304:role/cloudcraft",
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
			give: &cloudcraft.AWSAccount{
				ID:      "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
				Name:    "My updated AWS account",
				RoleARN: "arn:aws:iam::558791803304:role/cloudcraft",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil AWS account",
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
			give: &cloudcraft.AWSAccount{
				ID:      "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
				Name:    "My updated AWS account",
				RoleARN: "arn:aws:iam::558791803304:role/cloudcraft",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty ID",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				Name:    "My updated AWS account",
				RoleARN: "arn:aws:iam::558791803304:role/cloudcraft",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty name",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				ID:      "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
				RoleARN: "arn:aws:iam::558791803304:role/cloudcraft",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty role ARN",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: ctx,
			give: &cloudcraft.AWSAccount{
				ID:   "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
				Name: "My updated AWS account",
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

			got, err := client.AWS.Update(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWS().Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AWS().Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSService_Delete(t *testing.T) {
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
			name: "Valid AWS account data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			context: ctx,
			give:    "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
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
			give:    "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil context",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			context: nil,
			give:    "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
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

			got, err := client.AWS.Delete(tt.context, tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWS().Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AWS().Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSService_Snapshot(t *testing.T) {
	t.Parallel()

	var (
		validTestData = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "snapshot-valid.png"))
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
			name: "Valid AWS account snapshot",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			giveRegion: "us-east-1",
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
			giveID:     "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			giveRegion: "us-east-1",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantErr: true,
		},
		{
			name:       "Nil context",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			context:    nil,
			giveID:     "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			giveRegion: "us-east-1",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantErr: true,
		},
		{
			name: "Nil snapshot params",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			giveRegion: "us-east-1",
			giveFormat: "png",
			giveParams: nil,
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
		{
			name: "Empty format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context:    ctx,
			giveID:     "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			giveRegion: "us-east-1",
			giveFormat: "",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantWidth:  1920,
			wantHeight: 1080,
			wantErr:    false,
		},
		{
			name:       "Empty ID",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			context:    ctx,
			giveID:     "",
			giveRegion: "us-east-1",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
			wantErr: true,
		},
		{
			name:       "Empty region",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			context:    ctx,
			giveID:     "fe3e5b29-a0e8-41ca-91e2-02a0441b1d33",
			giveRegion: "",
			giveFormat: "png",
			giveParams: &cloudcraft.SnapshotParams{
				Width:  1920,
				Height: 1080,
			},
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

			got, _, err := client.AWS.Snapshot(tt.context, tt.giveID, tt.giveRegion, tt.giveFormat, tt.giveParams)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWSService.Snapshot() error = %v, wantErr %v", err, tt.wantErr)
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

func TestAWSService_IAMParameters(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "iam-parameters-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "generic-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    *cloudcraft.IAMParams
		wantErr bool
	}{
		{
			name: "Valid IAM parameters data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: &cloudcraft.IAMParams{
				AccountID:     "912185983511",
				ExternalID:    "4414aef4-8f04-4b0b-8425-d73b84dcaa2d",
				AWSConsoleURL: "https://console.aws.amazon.com/iam/home?#/roles$new?step=type&roleType=crossAccount&isThirdParty&accountID=912185983511&externalID=4414aef4-8f04-4b0b-8425-d73b84dcaa2d&roleName=cloudcraft&policies=arn:aws:iam::aws:policy%2FReadOnlyAccess",
			},
			wantErr: false,
		},
		{
			name: "Invalid IAM parameters data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
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

			got, _, err := client.AWS.IAMParameters(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWSService.IAMParameters() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AWSService.IAMParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWSService_IAMPolicy(t *testing.T) {
	t.Parallel()

	var (
		validTestData   = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "iam-policy-valid.json"))
		invalidTestData = xtesting.ReadFile(t, filepath.Join(_testAWSDataPath, "generic-invalid.json"))
		ctx             = context.Background()
	)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    *cloudcraft.IAMPolicy
		wantErr bool
	}{
		{
			name: "Valid IAM policy data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(validTestData)
			},
			context: ctx,
			want: &cloudcraft.IAMPolicy{
				Version: "2012-10-17",
				Statement: []cloudcraft.IAMStatement{
					{
						Effect: "Allow",
						Action: string("apigateway:GET"),
						Resource: []any{
							string("arn:aws:apigateway:*::/apis"),
							string("arn:aws:apigateway:*::/apis/*"),
							string("..."),
						},
					},
					{
						Effect: "Allow",
						Action: []any{
							string("autoscaling:DescribeAutoScalingGroups"),
							string("cassandra:Select"),
							string("..."),
						},
						Resource: string("*"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid IAM policy data",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write(invalidTestData)
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

			got, _, err := client.AWS.IAMPolicy(tt.context)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AWSService.IAMPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AWSService.IAMPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
