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
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
)

// blueprintPath is the path to the blueprint endpoint of the Cloudcraft API.
const blueprintPath string = "blueprint"

const (
	// ErrNilBlueprint is returned when you try to send a request without a
	// blueprint.
	ErrNilBlueprint xerrors.Error = "blueprint cannot be nil"

	// ErrBlueprintKey is returned when the response from the API does to a List
	// call is not a list of blueprints.
	ErrBlueprintKey xerrors.Error = "key 'blueprints' not found in the response"

	// ErrMissingID is returned when you try to send a request without the ID of
	// a blueprint.
	ErrMissingBlueprintID xerrors.Error = "missing blueprint ID"
)

const (
	// DefaultImageExportFormat is the default format used to export blueprint
	// images.
	DefaultImageExportFormat string = "png"

	// DefaultImageExportWidth is the default width used to export blueprint
	// images.
	DefaultImageExportWidth int = 1920

	// DefaultImageExportHeight is the default height used to export blueprint
	// images.
	DefaultImageExportHeight int = 1080

	// DefaultBudgetExportFormat is the default format used to export a
	// blueprint's budget.
	DefaultBudgetExportFormat string = "csv"

	// DefaultBudgetExportCurrency is the default currency used to export a
	// blueprint's budget.
	DefaultBudgetExportCurrency string = "USD"

	// DefaultBudgetExportPeriod is the default period used to export a blueprint's
	// budget.
	DefaultBudgetExportPeriod string = "m"
)

// BlueprintService handles communication with the "/blueprint" endpoint of
// Cloudcraft's developer API.
type BlueprintService service

// Blueprint represents a blueprint in Cloudcraft.
type Blueprint struct {
	CustomerID       *string        `json:"CustomerId,omitempty"`
	ReadAccess       *[]string      `json:"readAccess,omitempty"`
	WriteAccess      *[]string      `json:"writeAccess,omitempty"`
	Tags             *[]string      `json:"tags,omitempty"`
	Data             *BlueprintData `json:"data,omitempty"`
	CreatedAt        time.Time      `json:"createdAt,omitempty"`
	UpdatedAt        time.Time      `json:"updatedAt,omitempty"`
	ID               string         `json:"id,omitempty"`
	Name             string         `json:"name,omitempty"`
	CreatorID        string         `json:"CreatorId,omitempty"`
	CurrentVersionID string         `json:"CurrentVersionId,omitempty"`
	LastUserID       string         `json:"LastUserId,omitempty"`
}

// BlueprintData represents a collection of data that makes up a blueprint.
type BlueprintData struct {
	LiveAccount    *LiveAccount     `json:"liveAccount,omitempty"`
	Theme          *Theme           `json:"theme,omitempty"`
	LiveOptions    *LiveOptions     `json:"liveOptions,omitempty"`
	Name           string           `json:"name,omitempty"`
	Projection     string           `json:"projection,omitempty"`
	LinkKey        string           `json:"linkKey,omitempty"`
	Grid           string           `json:"grid,omitempty"`
	Images         []map[string]any `json:"images,omitempty"`
	Groups         []map[string]any `json:"groups,omitempty"`
	Nodes          []map[string]any `json:"nodes,omitempty"`
	Icons          []map[string]any `json:"icons,omitempty"`
	Surfaces       []map[string]any `json:"surfaces,omitempty"`
	Connectors     []map[string]any `json:"connectors,omitempty"`
	Edges          []map[string]any `json:"edges,omitempty"`
	Text           []map[string]any `json:"text,omitempty"`
	DisabledLayers []string         `json:"disabledLayers,omitempty"`
	Version        int              `json:"version,omitempty"`
	ShareDocs      bool             `json:"shareDocs,omitempty"`
}

// Theme represents the color scheme of a blueprint.
type Theme struct {
	Base string `json:"base,omitempty"`
}

// LiveAccount represents the AWS account that a blueprint is connected to.
type LiveAccount struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// LiveOptions represents options for a blueprint's live view.
type LiveOptions struct {
	ExcludedTypes      []string `json:"excludedTypes,omitempty"`
	AutoLabel          bool     `json:"autoLabel,omitempty"`
	AutoConnect        bool     `json:"autoConnect,omitempty"`
	UpdatesEnabled     bool     `json:"updatesEnabled,omitempty"`
	UpdateAllOnScan    bool     `json:"updateAllOnScan,omitempty"`
	UpdateGroupsOnScan bool     `json:"updateGroupsOnScan,omitempty"`
	UpdateNodeOnSelect bool     `json:"updateNodeOnSelect,omitempty"`
}

// ImageExportParams represents optional query parameters that can be used to
// customize an image export.
type ImageExportParams struct {
	PaperSize   string
	Grid        bool
	Transparent bool
	Landscape   bool
	Scale       float32
	Width       int
	Height      int
}

// query builds a query string from fields with non-zero values and returns it
// as url.Values.
func (p *ImageExportParams) query() url.Values {
	values := make(url.Values)

	if p.PaperSize != "" {
		values["paperSize"] = []string{p.PaperSize}
	}

	if p.Grid {
		values["grid"] = []string{"true"}
	}

	if p.Transparent {
		values["transparent"] = []string{"true"}
	}

	if p.Landscape {
		values["landscape"] = []string{"true"}
	}

	if p.Scale != 0 {
		scaleStr := strconv.FormatFloat(float64(p.Scale), 'f', -1, 32)

		values["scale"] = []string{scaleStr}
	}

	if p.Width != 0 {
		values["width"] = []string{strconv.Itoa(p.Width)}
	}

	if p.Height != 0 {
		values["height"] = []string{strconv.Itoa(p.Height)}
	}

	return values
}

// BudgetExportParams represents optional query parameters that can be used to
// customize an a budget export.
type BudgetExportParams struct {
	Currency string
	Period   string
	Rate     string
}

// query builds a query string from fields with non-zero values and returns it
// as url.Values.
func (p *BudgetExportParams) query() url.Values {
	values := url.Values{}

	if p.Currency != "" {
		values.Set("currency", p.Currency)
	}

	if p.Period != "" {
		values.Set("period", p.Period)
	}

	if p.Rate != "" {
		values.Set("rate", p.Rate)
	}

	return values
}

// List returns a list of blueprints.
//
// [API Reference].
//
// [API Reference]: https://developers.cloudcraft.co/#19d9d681-b3b7-4950-a0e0-aeb518101714
func (s *BlueprintService) List(ctx context.Context) ([]*Blueprint, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	ret, err := s.client.do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	var result map[string][]*Blueprint
	if err := json.Unmarshal(ret.Body, &result); err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	blueprints, ok := result["blueprints"]
	if !ok {
		return nil, nil, fmt.Errorf("%w", ErrBlueprintKey)
	}

	return blueprints, ret, nil
}

// Get retrieves a blueprint by its ID.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#dfc05b6e-a851-46aa-8019-c839eae7d695
func (s *BlueprintService) Get(ctx context.Context, id string) (*Blueprint, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if id == "" {
		return nil, nil, ErrMissingBlueprintID
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath) + len(id) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)
	endpoint.WriteString("/" + id)

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	var result *Blueprint
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return result, resp, nil
}

// Create creates a new blueprint.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#d72c9b37-9f03-4c24-98d0-92971493780f
func (s *BlueprintService) Create(ctx context.Context, blueprint *Blueprint) (*Blueprint, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if blueprint == nil {
		return nil, nil, ErrNilBlueprint
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)

	payload, err := json.Marshal(blueprint)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	req, err := s.client.request(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	var result *Blueprint
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return result, resp, nil
}

// Update updates an existing blueprint.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#7139bd5a-cf80-4bff-b2da-be0d35250b8f
func (s *BlueprintService) Update(ctx context.Context, blueprint *Blueprint, etag string) (*Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if blueprint == nil {
		return nil, ErrNilBlueprint
	}

	if blueprint.ID == "" {
		return nil, ErrMissingBlueprintID
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath) + len(blueprint.ID) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)
	endpoint.WriteString("/" + blueprint.ID)

	payload, err := json.Marshal(blueprint)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	req, err := s.client.request(ctx, http.MethodPut, endpoint.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if etag != "" {
		req.Header.Set("If-Match", etag)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

// Delete deletes a blueprint by ID.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#38e2767f-7b42-4573-85ba-6137b61fe0ef
func (s *BlueprintService) Delete(ctx context.Context, id string) (*Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if id == "" {
		return nil, ErrMissingBlueprintID
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath) + len(id) + 1)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)
	endpoint.WriteString("/" + id)

	req, err := s.client.request(ctx, http.MethodDelete, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

// ExportImage renders a blueprint for export in SVG, PNG, PDF or MxGraph format.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#8ad8ffa1-4a34-44e1-8795-4a851fc2fa58
func (s *BlueprintService) ExportImage(
	ctx context.Context,
	id string,
	format string,
	params *ImageExportParams,
) ([]byte, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if id == "" {
		return nil, nil, ErrMissingBlueprintID
	}

	if format == "" {
		format = DefaultImageExportFormat
	}

	if params == nil {
		params = &ImageExportParams{
			Width:  DefaultImageExportWidth,
			Height: DefaultImageExportHeight,
		}
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath) + len(id) + len(format) + 2)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)
	endpoint.WriteString("/" + id)
	endpoint.WriteString("/" + format)

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
		return nil, nil, fmt.Errorf("%w", err)
	}

	return resp.Body, resp, nil
}

// ExportBudget exports a blueprint's budget in CSV or XLSX format.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#4280d5b3-c9a1-423f-8074-0499447dd8d6
func (s *BlueprintService) ExportBudget(
	ctx context.Context,
	id string,
	format string,
	params *BudgetExportParams,
) ([]byte, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if id == "" {
		return nil, nil, ErrMissingBlueprintID
	}

	if format == "" {
		format = DefaultBudgetExportFormat
	}

	if params == nil {
		params = &BudgetExportParams{
			Currency: DefaultBudgetExportCurrency,
			Period:   DefaultBudgetExportPeriod,
		}
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(blueprintPath) + len(id) + len(format) + 9)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(blueprintPath)
	endpoint.WriteString("/" + id)
	endpoint.WriteString("/budget/" + format)

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
		return nil, nil, fmt.Errorf("%w", err)
	}

	return resp.Body, resp, nil
}
