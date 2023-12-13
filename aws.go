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

// awsAccountPath is the path to the AWS endpoint of the Cloudcraft API.
const awsAccountPath string = "aws/account"

const (
	// ErrEmptyRoleARN is returned when the AWS account's role ARN is empty.
	ErrEmptyRoleARN xerrors.Error = "role ARN cannot be empty"
)

const (
	// DefaultSnapshotFormat is the default format used for account snapshots.
	DefaultSnapshotFormat string = "png"

	// DefaultSnapshotWidth is the default width used for account snapshots.
	DefaultSnapshotWidth int = 1920

	// DefaultSnapshotHeight is the default height used for account snapshots.
	DefaultSnapshotHeight int = 1080
)

// AWSService handles communication with the "/aws" endpoint of Cloudcraft's
// developer API.
type AWSService service

// AWSAccount represents an AWS account registered with Cloudcraft.
type AWSAccount struct {
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ReadAccess  *[]string `json:"readAccess,omitempty"`
	WriteAccess *[]string `json:"writeAccess,omitempty"`
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	RoleARN     string    `json:"roleArn,omitempty"`
	ExternalID  string    `json:"externalId,omitempty"`
	CreatorID   string    `json:"CreatorId,omitempty"`
	Source      string    `json:"source,omitempty"`
}

// IAMParams represents the AWS IAM role parameters used by Cloudcraft.
type IAMParams struct {
	AccountID     string `json:"accountId,omitempty"`
	ExternalID    string `json:"externalId,omitempty"`
	AWSConsoleURL string `json:"awsConsoleUrl,omitempty"`
}

// IAMPolicy represents the AWS IAM policy used by Cloudcraft.
type IAMPolicy struct {
	Version   string         `json:"Version,omitempty"`
	Statement []IAMStatement `json:"Statement,omitempty"`
}

// IAMStatement represents an AWS IAM policy statement.
type IAMStatement struct {
	Action   any    `json:"Action,omitempty"`
	Resource any    `json:"Resource,omitempty"`
	Effect   string `json:"Effect,omitempty"`
}

// List lists your AWS accounts linked with Cloudcraft.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#a83b30f1-8949-4c68-9944-2e2ab2710670
func (s *AWSService) List(ctx context.Context) ([]*AWSAccount, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(awsAccountPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	var result map[string][]*AWSAccount
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	accounts, ok := result["accounts"]
	if !ok {
		return nil, resp, fmt.Errorf("%w", ErrAccountsKey)
	}

	return accounts, resp, nil
}

// Create registers a new AWS account with Cloudcraft.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#51c4726e-ce1a-4e16-8b3f-f15dcee0aebe
func (s *AWSService) Create(ctx context.Context, account *AWSAccount) (*AWSAccount, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if account == nil {
		return nil, nil, ErrNilAccount
	}

	if account.Name == "" {
		return nil, nil, ErrEmptyAccountName
	}

	if account.RoleARN == "" {
		return nil, nil, ErrEmptyRoleARN
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(awsAccountPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)

	payload, err := json.Marshal(account)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	req, err := s.client.request(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	var result *AWSAccount
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
func (s *AWSService) Update(ctx context.Context, account *AWSAccount) (*Response, error) {
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

	if account.RoleARN == "" {
		return nil, ErrEmptyRoleARN
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(awsAccountPath) + len(account.ID) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)
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

	resp, err := s.client.do(ctx, req)
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
func (s *AWSService) Delete(ctx context.Context, id string) (*Response, error) {
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

	endpoint.Grow(len(baseURL) + len(awsAccountPath) + len(id) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)
	endpoint.WriteByte('/')
	endpoint.WriteString(id)

	req, err := s.client.request(ctx, http.MethodDelete, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return resp, fmt.Errorf("%w", err)
	}

	return resp, nil
}

// Snapshot scans and render a region of an AWS account into a blueprint in
// JSON, SVG, PNG, PDF or MxGraph format.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#13e7daaf-e22a-42c6-b6bc-e34a24f05e60
func (s *AWSService) Snapshot(
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

	endpoint.Grow(len(baseURL) + len(awsAccountPath) + len(id) + len(region) + len(format) + 3)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)
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

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	return resp.Body, resp, nil
}

// IAMParameters list all parameters required for registering a new IAM Role in
// AWS for use with Cloudcraft.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#aa18999e-f6da-4628-96bd-49d5a286b928
func (s *AWSService) IAMParameters(ctx context.Context) (*IAMParams, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(awsAccountPath) + len("/iamParameters"))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)
	endpoint.WriteString("/iamParameters")

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	var result *IAMParams
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	return result, resp, nil
}

// IAMPolicy lists all permissions required for registering a new IAM Role in AWS for use with Cloudcraft.
//
// [API reference].
//
// [API reference]: https://help.cloudcraft.co/article/64-minimal-iam-policy
func (s *AWSService) IAMPolicy(ctx context.Context) (*IAMPolicy, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(awsAccountPath) + len("/iamParameters/policy/minimal"))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(awsAccountPath)
	endpoint.WriteString("/iamParameters/policy/minimal")

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	var result *IAMPolicy
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	return result, resp, nil
}
