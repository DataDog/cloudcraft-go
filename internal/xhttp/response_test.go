// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xhttp_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/DataDog/cloudcraft-go/internal/xhttp"
)

type errReader struct{}

func (*errReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

type customReadCloser struct {
	data *bytes.Buffer
}

func (c *customReadCloser) Read(p []byte) (n int, err error) {
	return c.data.Read(p)
}

func (*customReadCloser) Close() error {
	return errors.New("mock close error")
}

func TestDrainResponseBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		resp    *http.Response
		wantErr bool
	}{
		{
			name: "valid response body",
			resp: &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte("valid response body"))),
			},
			wantErr: false,
		},
		{
			name: "empty response body",
			resp: &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(""))),
			},
			wantErr: false,
		},
		{
			name: "error response body",
			resp: &http.Response{
				Body: io.NopCloser(&errReader{}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := xhttp.DrainResponseBody(tt.resp)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDrainResponseBody_ErrorClose(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &customReadCloser{data: bytes.NewBufferString("test")},
	}

	err := xhttp.DrainResponseBody(resp)
	if err == nil {
		t.Error("expected error, got nil")
	}

	want := fmt.Errorf("%w: %w", xhttp.ErrCannotCloseResponse, errors.New("mock close error"))
	if err.Error() != want.Error() {
		t.Errorf("got: %v, want: %v", err, want)
	}
}
