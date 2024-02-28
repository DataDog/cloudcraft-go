// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/DataDog/cloudcraft-go/internal/endpoint"
	"github.com/DataDog/cloudcraft-go/internal/meta"
	"github.com/DataDog/cloudcraft-go/internal/xerrors"
	"github.com/DataDog/cloudcraft-go/internal/xhttp"
)

const (
	// ErrInvalidConfig is returned when a Client is created with an invalid
	// Config.
	ErrInvalidConfig xerrors.Error = "invalid config"

	// ErrRequestFailed is returned when a request to the Cloudcraft API fails
	// for unknown reasons.
	ErrRequestFailed xerrors.Error = "request failed with status code"

	// ErrMaxRetriesExceeded is returned when the maximum number of retries is
	// exceeded for HTTP requests.
	ErrMaxRetriesExceeded xerrors.Error = "maximum number of retries exceeded"
)

type (
	// Service is a common struct that can be reused instead of allocating a new
	// one for each service on the heap.
	service struct {
		client *Client
	}

	// Client is a client for the Cloudcraft API.
	Client struct {
		// httpClient is the underlying HTTP client used by the API client.
		httpClient *http.Client

		// retryPolicy specifies the policy used to retry failed requests.
		retryPolicy *xhttp.RetryPolicy

		// cfg specifies the configuration used by the API client.
		cfg *Config

		// Cloudcraft API service fields.
		Azure     *AzureService
		AWS       *AWSService
		Blueprint *BlueprintService
		User      *UserService

		// common specifies a common service shared by all services.
		common service
	}
)

// NewClient returns a new Client given a Config. If Config is nil, NewClient
// will try to look up the configuration from the environment.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = NewConfigFromEnv()
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	baseURL, err := endpoint.Parse(cfg.Scheme, cfg.Host, cfg.Port, cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	cfg.endpoint = baseURL

	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = DefaultMaxRetries
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultTimeout
	}

	client := &Client{
		httpClient: xhttp.NewClient(cfg.Timeout),
		retryPolicy: &xhttp.RetryPolicy{
			IsRetryable:   xhttp.DefaultIsRetryable,
			MaxRetries:    cfg.MaxRetries,
			MinRetryDelay: xhttp.DefaultMinRetryDelay,
			MaxRetryDelay: xhttp.DefaultMaxRetryDelay,
		},
		cfg: cfg,
	}

	client.common.client = client
	client.Azure = (*AzureService)(&client.common)
	client.AWS = (*AWSService)(&client.common)
	client.Blueprint = (*BlueprintService)(&client.common)
	client.User = (*UserService)(&client.common)

	return client, nil
}

// SnapshotParams represents query parameters used to customize an Azure or AWS
// account snapshot.
type SnapshotParams struct {
	PaperSize   string
	Projection  string
	Theme       string
	Filter      []string
	Exclude     []string
	Label       bool
	Autoconnect bool
	Grid        bool
	Transparent bool
	Landscape   bool
	Scale       float32
	Width       int
	Height      int
}

// query builds a query string from fields with non-zero values and returns it
// as url.Values.
func (p *SnapshotParams) query() url.Values {
	values := url.Values{}

	if p.PaperSize != "" {
		values.Set("paperSize", p.PaperSize)
	}

	if p.Projection != "" {
		values.Set("projection", p.Projection)
	}

	if p.Theme != "" {
		values.Set("theme", p.Theme)
	}

	if len(p.Filter) > 0 {
		values.Set("filter", strings.Join(p.Filter, ","))
	}

	if len(p.Exclude) > 0 {
		values.Set("exclude", strings.Join(p.Exclude, ","))
	}

	if p.Label {
		values.Set("label", "true")
	}

	if p.Autoconnect {
		values.Set("autoconnect", "true")
	}

	if p.Grid {
		values.Set("grid", "true")
	}

	if p.Transparent {
		values.Set("transparent", "true")
	}

	if p.Landscape {
		values.Set("landscape", "true")
	}

	if p.Scale != 0 {
		values.Set("scale", strconv.FormatFloat(float64(p.Scale), 'f', -1, 32))
	}

	if p.Width != 0 {
		values.Set("width", strconv.Itoa(p.Width))
	}

	if p.Height != 0 {
		values.Set("height", strconv.Itoa(p.Height))
	}

	return values
}

// Response represents a response from the Cloudcraft API.
type Response struct {
	// Header contains the response headers.
	Header http.Header

	// Body contains the response body as a byte slice.
	Body []byte

	// Status is the HTTP status code of the response.
	Status int
}

// do performs an HTTP request using the underlying HTTP client.
func (c *Client) do(req *http.Request) (*Response, error) { //nolint:gocyclo // Necessary complexity.
	var (
		attempt int
		resp    *http.Response
		err     error
		body    *bytes.Buffer
	)

	if req.Body != nil {
		body = bytes.NewBuffer(make([]byte, 0))

		_, err = io.Copy(body, req.Body)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		req.Body = io.NopCloser(body)

		if err = req.Body.Close(); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	for attempt = 0; attempt <= c.retryPolicy.MaxRetries; attempt++ {
		if body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body.Bytes()))
		}

		resp, err = c.httpClient.Do(req)
		if err != nil || !c.retryPolicy.IsRetryable(resp, err) {
			break
		}

		if resp != nil {
			if err = xhttp.DrainResponseBody(resp); err != nil {
				_ = resp.Body.Close()
			}
		}

		waitErr := c.retryPolicy.Wait(req.Context(), attempt)
		if waitErr != nil {
			return nil, fmt.Errorf("%w", waitErr)
		}
	}

	if resp == nil && attempt >= c.retryPolicy.MaxRetries {
		return nil, fmt.Errorf("%w: %d", ErrMaxRetriesExceeded, attempt)
	}

	if err != nil {
		select {
		case <-req.Context().Done():
			return nil, fmt.Errorf("%w", req.Context().Err())
		default:
			return nil, fmt.Errorf("%w", err)
		}
	}

	defer func() {
		if err = xhttp.DrainResponseBody(resp); err != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode > http.StatusNoContent {
		return nil, fmt.Errorf("%w: %d", ErrRequestFailed, resp.StatusCode)
	}

	var buffer *bytes.Buffer

	if resp.ContentLength > 0 {
		buffer = bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0))
	}

	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Response{
		Header: resp.Header,
		Body:   buffer.Bytes(),
		Status: resp.StatusCode,
	}, nil
}

// request is a convenience function for creating an HTTP request.
func (c *Client) request(
	ctx context.Context,
	method, uri string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.Key)
	req.Header.Set("User-Agent", meta.UserAgent)

	return req, nil
}
