// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// userPath is the path to the user endpoint of the Cloudcraft API.
const userPath string = "user"

// UserService handles communication with the "/user" endpoint of Cloudcraft's
// developer API.
type UserService service

// User represents a Cloudcraft user.
type User struct {
	AccessedAt time.Time      `json:"accessedAt,omitempty"`
	CreatedAt  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `json:"updatedAt,omitempty"`
	Settings   map[string]any `json:"settings,omitempty"`
	ID         string         `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	Email      string         `json:"email,omitempty"`
}

// Me returns the user profile.
//
// [API reference].
//
// [API reference]: https://developers.cloudcraft.co/#a1ac9d21-3d47-4338-b171-8419872f818a
func (s *UserService) Me(ctx context.Context) (*User, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(userPath) + 3)

	endpoint.WriteString(baseURL)
	endpoint.WriteString(userPath)
	endpoint.WriteString("/me")

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	ret, err := s.client.do(req)
	if err != nil {
		return nil, nil, err
	}

	var user *User
	if err := json.Unmarshal(ret.Body, &user); err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return user, ret, nil
}
