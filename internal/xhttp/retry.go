package xhttp

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	_backoffFactor float64 = 2.0
	_jitterFactor  float64 = 0.1
)

const (
	// DefaultMinRetryDelay is the default minimum duration to wait before
	// retrying a request.
	DefaultMinRetryDelay time.Duration = 1 * time.Second

	// DefaultMaxRetryDelay is the default maximum duration to wait before
	// retrying a request.
	DefaultMaxRetryDelay time.Duration = 30 * time.Second
)

// RetryPolicy defines a policy for retrying HTTP requests.
type RetryPolicy struct {
	// IsRetryable determines whether a given response and error combination
	// should be retried.
	IsRetryable func(*http.Response, error) bool

	// MaxRetries is the maximum number of times a request will be retried.
	MaxRetries int

	// MinRetryDelay is the minimum duration to wait before retrying a request.
	MinRetryDelay time.Duration

	// MaxRetryDelay is the maximum duration to wait before retrying a request.
	MaxRetryDelay time.Duration
}

// Wait calculates the time to wait before the next retry attempt and blocks
// until it is time to retry the request or the context is canceled. It
// incorporates an exponential backoff strategy with jitter to prevent the
// "thundering herd" problem.
//
// If the context is canceled before the wait is over, Wait returns the
// context's error.
func (p *RetryPolicy) Wait(ctx context.Context, attempt int) error {
	waitTime := float64(p.MinRetryDelay) * math.Pow(_backoffFactor, float64(attempt))

	if time.Duration(waitTime) > p.MaxRetryDelay {
		waitTime = float64(p.MaxRetryDelay)
	}

	jitter := (rand.Float64()*2 - 1) * _jitterFactor * waitTime //nolint:gosec // we don't need cryptographic randomness

	waitTimeWithJitter := time.Duration(waitTime + jitter)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", ctx.Err())
	case <-time.After(waitTimeWithJitter):
		return nil
	}
}

// DefaultIsRetryable defines the default logic to determine if a request should
// be retried.
//
// It returns true if an error occurs or the response status code indicates a
// retry may be successful (e.g., 202, 408, 429, 502, 503, 504).
func DefaultIsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case http.StatusAccepted, // 202
		http.StatusRequestTimeout,     // 408
		http.StatusTooManyRequests,    // 429
		http.StatusBadGateway,         // 502
		http.StatusServiceUnavailable, // 503
		http.StatusGatewayTimeout:     // 504
		return true
	}

	return false
}
