package xhttp_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/DataDog/cloudcraft-go/internal/xhttp"
)

func TestRetryPolicy_Wait(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		maxRetries      int
		minRetryDelay   time.Duration
		maxRetryDelay   time.Duration
		attempt         int
		contextFunc     func() (context.Context, context.CancelFunc)
		expectedWaitMax time.Duration
		expectedErr     error
	}{
		{
			name:          "Context canceled",
			maxRetries:    3,
			minRetryDelay: 1 * time.Second,
			maxRetryDelay: 30 * time.Second,
			attempt:       1,
			contextFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())

				cancel() // Cancel the context immediately.

				return ctx, cancel
			},
			expectedWaitMax: 1 * time.Millisecond,
			expectedErr:     context.Canceled,
		},
		{
			name:          "Maximum retry delay exceeded",
			maxRetries:    3,
			minRetryDelay: 1 * time.Second,
			maxRetryDelay: 2 * time.Second,
			attempt:       4, // This should result in a calculated delay > maxRetryDelay
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			expectedWaitMax: 3 * time.Second,
			expectedErr:     nil,
		},
		{
			name:          "Minimum retry delay",
			maxRetries:    3,
			minRetryDelay: 1 * time.Second,
			maxRetryDelay: 30 * time.Second,
			attempt:       0, // First attempt
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			expectedWaitMax: 1 * time.Second,
			expectedErr:     nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			policy := xhttp.RetryPolicy{
				IsRetryable:   xhttp.DefaultIsRetryable,
				MaxRetries:    tt.maxRetries,
				MinRetryDelay: tt.minRetryDelay,
				MaxRetryDelay: tt.maxRetryDelay,
			}

			ctx, cancel := tt.contextFunc()
			defer cancel()

			var (
				start           = time.Now()
				err             = policy.Wait(ctx, tt.attempt)
				end             = time.Now()
				expectedEndTime = start.Add(tt.expectedWaitMax)
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Wait() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if end.After(expectedEndTime) {
				t.Errorf("Wait() end time = %v, expectedEndTime %v", end, expectedEndTime)
			}
		})
	}
}

func TestDefaultIsRetryable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		giveResp  *http.Response
		giveError error
		want      bool
	}{
		{
			name:      "Error occurs",
			giveResp:  nil,
			giveError: errors.New("error"),
			want:      true,
		},
		{
			name: "Status code 202",
			giveResp: &http.Response{
				StatusCode: http.StatusAccepted,
			},
			giveError: nil,
			want:      true,
		},
		{
			name: "Status code 408",
			giveResp: &http.Response{
				StatusCode: http.StatusRequestTimeout,
			},
			giveError: nil,
			want:      true,
		},
		{
			name: "Status code 429",
			giveResp: &http.Response{
				StatusCode: http.StatusTooManyRequests,
			},
			giveError: nil,
			want:      true,
		},
		{
			name: "Status code 502",
			giveResp: &http.Response{
				StatusCode: http.StatusBadGateway,
			},
			giveError: nil,
			want:      true,
		},
		{
			name: "Status code 503",
			giveResp: &http.Response{
				StatusCode: http.StatusServiceUnavailable,
			},
			giveError: nil,
			want:      true,
		},
		{
			name: "Status code 504",
			giveResp: &http.Response{
				StatusCode: http.StatusGatewayTimeout,
			},
			giveError: nil,
			want:      true,
		},
		{
			name: "Status code not triggering retry",
			giveResp: &http.Response{
				StatusCode: http.StatusOK,
			},
			giveError: nil,
			want:      false,
		},
		{
			name:      "No error and no response",
			giveResp:  nil,
			giveError: nil,
			want:      false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			want := xhttp.DefaultIsRetryable(tt.giveResp, tt.giveError)
			if want != tt.want {
				t.Errorf("DefaultIsRetryable() = %v, want %v", want, tt.want)
			}
		})
	}
}
