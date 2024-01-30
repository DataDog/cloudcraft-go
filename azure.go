// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
)

// azureAccountPath is the path to the Azure endpoint of the Cloudcraft API.
const azureAccountPath string = "azure/account"

const (
	// ErrEmptyApplicationID is returned when an Azure account is created with
	// an empty application ID.
	ErrEmptyApplicationID xerrors.Error = "field 'ApplicationID' cannot be empty"

	// ErrEmptyDirectoryID is returned when an Azure account is created with
	// an empty directory ID.
	ErrEmptyDirectoryID xerrors.Error = "field 'DirectoryID' cannot be empty"

	// ErrEmptySubscriptionID is returned when an Azure account is created with
	// an empty subscription ID.
	ErrEmptySubscriptionID xerrors.Error = "field 'SubscriptionID' cannot be empty"

	// ErrEmptyClientSecret is returned when an Azure account is created with
	// an empty client secret.
	ErrEmptyClientSecret xerrors.Error = "field 'ClientSecret' cannot be empty"
)

// AzureService handles communication with the "/azure" endpoint of Cloudcraft's
// developer API.
type AzureService service

// AzureAccount represents an Azure account registered with Cloudcraft.
type AzureAccount struct {
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
	ReadAccess     *[]string `json:"readAccess,omitempty"`
	WriteAccess    *[]string `json:"writeAccess,omitempty"`
	CustomerID     *string   `json:"CustomerId,omitempty"`
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	ApplicationID  string    `json:"applicationId,omitempty"`
	DirectoryID    string    `json:"directoryId,omitempty"`
	SubscriptionID string    `json:"subscriptionId,omitempty"`
	ClientSecret   string    `json:"clientSecret,omitempty"`
	CreatorID      string    `json:"CreatorId,omitempty"`
	Hint           string    `json:"hint,omitempty"`
	Source         string    `json:"source,omitempty"`
}

// List returns a list of Azure accounts linked with Cloudcraft.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#29470635-2970-4205-8256-85c5835b92a1
func (s *AzureService) List(ctx context.Context) ([]*AzureAccount, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(azureAccountPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(azureAccountPath)

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	var result map[string][]*AzureAccount
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	accounts, ok := result["accounts"]
	if !ok {
		return nil, resp, fmt.Errorf("%w", ErrAccountsKey)
	}

	return accounts, resp, nil
}

// Create registers a new Azure account with Cloudcraft.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#09a9a67d-c807-45c1-b8a8-f5a6df08da12
func (s *AzureService) Create(ctx context.Context, account *AzureAccount) (*AzureAccount, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if account == nil {
		return nil, nil, ErrNilAccount
	}

	if account.Name == "" {
		return nil, nil, ErrEmptyAccountName
	}

	if account.ApplicationID == "" {
		return nil, nil, ErrEmptyApplicationID
	}

	if account.DirectoryID == "" {
		return nil, nil, ErrEmptyDirectoryID
	}

	if account.SubscriptionID == "" {
		return nil, nil, ErrEmptySubscriptionID
	}

	if account.ClientSecret == "" {
		return nil, nil, ErrEmptyClientSecret
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(azureAccountPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(azureAccountPath)

	payload, err := json.Marshal(account)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	req, err := s.client.request(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	var result *AzureAccount
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	return result, resp, nil
}

// Update updates an AWS account registered in Cloudcraft.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#d04fdf78-ea33-4846-a8b2-bb5e693e8f64
func (s *AzureService) Update(ctx context.Context, account *AzureAccount) (*Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if account == nil {
		return nil, ErrNilAccount
	}

	if account.ID == "" {
		return nil, ErrEmptyAccountID
	}

	if account.Name == "" {
		return nil, ErrEmptyAccountName
	}

	if account.ApplicationID == "" {
		return nil, ErrEmptyApplicationID
	}

	if account.DirectoryID == "" {
		return nil, ErrEmptyDirectoryID
	}

	if account.SubscriptionID == "" {
		return nil, ErrEmptySubscriptionID
	}

	if account.ClientSecret == "" {
		return nil, ErrEmptyClientSecret
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(azureAccountPath) + len(account.ID) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(azureAccountPath)
	endpoint.WriteByte('/')
	endpoint.WriteString(account.ID)

	payload, err := json.Marshal(account)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	req, err := s.client.request(ctx, http.MethodPut, endpoint.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return resp, fmt.Errorf("%w", err)
	}

	return resp, nil
}

// Delete deletes a registered AWS account from Cloudcraft by ID.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#c4788665-d408-4535-8aa1-bf27dfb064aa
func (s *AzureService) Delete(ctx context.Context, id string) (*Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if id == "" {
		return nil, ErrEmptyAccountID
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(azureAccountPath) + len(id) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(azureAccountPath)
	endpoint.WriteByte('/')
	endpoint.WriteString(id)

	req, err := s.client.request(ctx, http.MethodDelete, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return resp, fmt.Errorf("%w", err)
	}

	return resp, nil
}

// Snapshot scans and render a region of an Azure account into a blueprint in
// JSON, SVG, PNG, PDF or MxGraph format.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#e687cfa9-f181-4eaf-bf76-f167235fa4fe
func (s *AzureService) Snapshot(
	ctx context.Context,
	id, region, format string,
	params *SnapshotParams,
) ([]byte, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if id == "" {
		return nil, nil, ErrEmptyAccountID
	}

	if region == "" {
		return nil, nil, ErrEmptyRegion
	}

	if format == "" {
		format = DefaultSnapshotFormat
	}

	if params == nil {
		params = &SnapshotParams{
			Width:  DefaultSnapshotWidth,
			Height: DefaultSnapshotHeight,
		}
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(azureAccountPath) + len(id) + len(region) + len(format) + 3)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(azureAccountPath)
	endpoint.WriteByte('/')
	endpoint.WriteString(id)
	endpoint.WriteByte('/')
	endpoint.WriteString(region)
	endpoint.WriteByte('/')
	endpoint.WriteString(format)

	u, err := url.Parse(endpoint.String())
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	u.RawQuery = params.query().Encode()

	req, err := s.client.request(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	return resp.Body, resp, nil
}
