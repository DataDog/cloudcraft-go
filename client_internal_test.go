// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSnapshotParams_Query(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give SnapshotParams
		want url.Values
	}{
		{
			name: "Empty parameters",
			want: url.Values{},
		},
		{
			name: "All parameters set",
			give: SnapshotParams{
				PaperSize:   "A4",
				Projection:  "top",
				Theme:       "dark",
				Filter:      []string{"instance", "database"},
				Exclude:     []string{"network"},
				Label:       true,
				Autoconnect: true,
				Grid:        true,
				Transparent: true,
				Landscape:   true,
				Scale:       2.0,
				Width:       1920,
				Height:      1080,
			},
			want: url.Values{
				"paperSize":   []string{"A4"},
				"projection":  []string{"top"},
				"theme":       []string{"dark"},
				"filter":      []string{"instance,database"},
				"exclude":     []string{"network"},
				"label":       []string{"true"},
				"autoconnect": []string{"true"},
				"grid":        []string{"true"},
				"transparent": []string{"true"},
				"landscape":   []string{"true"},
				"scale":       []string{strconv.FormatFloat(2.0, 'f', -1, 32)},
				"width":       []string{"1920"},
				"height":      []string{"1080"},
			},
		},
		{
			name: "Only a few parameters set",
			give: SnapshotParams{
				Theme:     "light",
				Landscape: true,
			},
			want: url.Values{
				"theme":     []string{"light"},
				"landscape": []string{"true"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.give.query(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SnapshotParams.query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		context context.Context
		want    *Response
		wantErr bool
	}{
		{
			name: "Valid response",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				w.Write([]byte(`Hello, World!`))
			},
			context: ctx,
			want: &Response{
				Header: http.Header{
					"Content-Length": []string{"13"},
					"Content-Type":   []string{"text/plain; charset=utf-8"},
					"Date": []string{
						time.Now().In(time.UTC).Format(http.TimeFormat),
					},
				},
				Body:   []uint8{'H', 'e', 'l', 'l', 'o', ',', ' ', 'W', 'o', 'r', 'l', 'd', '!'},
				Status: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Context timeout",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				time.Sleep(100 * time.Millisecond)

				w.WriteHeader(http.StatusOK)

				w.Write([]byte(`Delayed response`))
			},
			context: func() context.Context {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 50*time.Millisecond)

				t.Cleanup(cancel)

				return ctxWithTimeout
			}(),
			wantErr: true,
		},
		{
			name: "Invalid HTTP status code",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
			context: ctx,
			wantErr: true,
		},
		{
			name: "HTTP Client Do error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				conn, _, _ := w.(http.Hijacker).Hijack() //nolint:forcetypeassert // should be fine for testing
				conn.Close()
			},
			context: ctx,
			wantErr: true,
		},
		{
			name: "Response Body Read Error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)

				// Instead of writing to the response writer, we'll set a custom
				// body that fails on reading.
				hijacker, ok := w.(http.Hijacker)
				if !ok {
					t.Fatal("ResponseWriter does not support Hijacker interface")
				}

				conn, _, err := hijacker.Hijack()
				if err != nil {
					t.Fatal("Hijack failed:", err)
				}

				_, _ = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 25\r\nContent-Type: text/plain\r\n\r\n"))
				conn.Close()
			},
			context: ctx,
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
				t.Fatalf("failed to parse mock server URL: %v", err)
			}

			cfg := &Config{
				Scheme: endpoint.Scheme,
				Host:   endpoint.Hostname(),
				Port:   endpoint.Port(),
				Path:   DefaultPath,
				Key:    "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
			}

			client, err := NewClient(cfg)
			if err != nil {
				t.Fatalf("failed to create client for mock tests: %v", err)
			}

			req, err := client.request(tt.context, http.MethodGet, endpoint.String(), http.NoBody)
			if err != nil {
				t.Fatalf("Request() error = %v", err)
			}

			got, err := client.do(req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Scheme: DefaultScheme,
		Host:   DefaultHost,
		Path:   DefaultPath,
		Key:    "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client for mock tests: %v", err)
	}

	tests := []struct {
		name    string
		method  string
		uri     string
		want    *http.Request
		wantErr bool
	}{
		{
			name:   "Valid request",
			method: http.MethodGet,
			uri:    "https://example.com",
			want: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid request",
			method:  http.MethodGet,
			uri:     "://example.com",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := client.request(context.Background(), tt.method, tt.uri, http.NoBody)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Request() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got == nil && tt.wantErr {
				return
			}

			if got.Method != tt.want.Method {
				t.Fatalf("Request().Method = %v, want %v", got.Method, tt.want.Method)
			}

			if got.URL.String() != tt.want.URL.String() {
				t.Fatalf("Request().URL = %v, want %v", got.URL, tt.want.URL)
			}
		})
	}
}

func TestDo_Retries_GET(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		giveResponses      []int
		giveRetries        int
		giveDelay          time.Duration
		giveTimeoutContext time.Duration
		wantAttempts       int
		wantFinalResponse  int
		wantErr            bool
	}{
		{
			name: "Success on first attempt",
			giveResponses: []int{
				http.StatusOK,
			},
			giveRetries:        3,
			giveDelay:          0,
			giveTimeoutContext: 5 * time.Second,
			wantAttempts:       1,
			wantFinalResponse:  http.StatusOK,
			wantErr:            false,
		},
		{
			name: "Success on third attempt",
			giveResponses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusOK,
			},
			giveRetries:        3,
			giveDelay:          0,
			giveTimeoutContext: 5 * time.Second,
			wantAttempts:       3,
			wantFinalResponse:  http.StatusOK,
			wantErr:            false,
		},
		{
			name: "Failure by exceeding max retries",
			giveResponses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
			giveRetries:        2,
			giveDelay:          0,
			giveTimeoutContext: 5 * time.Second,
			wantAttempts:       3,
			wantFinalResponse:  http.StatusServiceUnavailable,
			wantErr:            true,
		},
		{
			name: "Failure by context cancellation",
			giveResponses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
			giveRetries:        3,
			giveDelay:          2 * time.Second,
			giveTimeoutContext: 1 * time.Second,
			wantAttempts:       0,
			wantFinalResponse:  http.StatusInternalServerError,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				mu       sync.Mutex
				requests int
				srv      = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					time.Sleep(tt.giveDelay)

					respIndex := requests

					if respIndex >= len(tt.giveResponses) {
						respIndex = len(tt.giveResponses) - 1
					}

					w.WriteHeader(tt.giveResponses[respIndex])

					mu.Lock()
					requests++
					mu.Unlock()
				}))
			)

			defer srv.Close()

			endpoint, err := url.Parse(srv.URL)
			if err != nil {
				t.Fatalf("failed to parse mock server URL: %v", err)
			}

			cfg := &Config{
				Scheme:     endpoint.Scheme,
				Host:       endpoint.Hostname(),
				Port:       endpoint.Port(),
				Path:       DefaultPath,
				MaxRetries: tt.giveRetries,
				Key:        "not-a-real-key-oRbwhd5RTvWsPJ89ZkASHU13qcyd=",
			}

			client, err := NewClient(cfg)
			if err != nil {
				t.Fatalf("failed to create client for mock tests: %v", err)
			}

			client.retryPolicy.MinRetryDelay = 1 * time.Millisecond
			client.retryPolicy.MaxRetryDelay = 1 * time.Millisecond
			client.retryPolicy.IsRetryable = func(resp *http.Response, _ error) bool {
				return resp.StatusCode >= 500
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.giveTimeoutContext)
			defer cancel()

			req, err := client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
			if err != nil {
				t.Fatalf("failed to create mock request: %v", err)
			}

			got, err := client.do(req)
			if tt.wantErr != (err != nil) {
				t.Fatalf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got.Status != tt.wantFinalResponse {
				t.Fatalf("unexpected final response status: got %d, want %d", got.Status, tt.wantFinalResponse)
			}

			mu.Lock()
			gotRequests := requests
			mu.Unlock()

			if gotRequests != tt.wantAttempts {
				t.Fatalf("unexpected number of attempts: got %d, want %d", gotRequests, tt.wantAttempts)
			}
		})
	}
}
