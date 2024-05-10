package cloudcraft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// teamPath is the path to the team endpoint of the Cloudcraft API.
const teamPath string = "team"

// TeamService handles communication with the "/team" endpoint of Cloudcraft's
// developer API.
type TeamService service

// Team represents a team in Cloudcraft.
type Team struct {
	UpdatedAt           time.Time `json:"updatedAt,omitempty"`
	CreatedAt           time.Time `json:"createdAt,omitempty"`
	ID                  string    `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	CustomerID          string    `json:"customerId,omitempty"`
	Role                string    `json:"role,omitempty"`
	Members             []Members `json:"members,omitempty"`
	Visible             bool      `json:"visible,omitempty"`
	CrossOrganizational bool      `json:"crossOrganizational,omitempty"`
	ExternalSharing     bool      `json:"externalSharing,omitempty"`
}

// Members represents a list of members in a team.
type Members struct {
	ID         string  `json:"id,omitempty"`
	Role       string  `json:"role,omitempty"`
	UserID     *string `json:"userId,omitempty"`
	Name       *string `json:"name,omitempty"`
	Email      string  `json:"email,omitempty"`
	MFAEnabled bool    `json:"mfaEnabled,omitempty"`
}

// List returns a list of teams.
//
// [API Reference].
//
// [API Reference]: https://docs.datadoghq.com/cloudcraft/api/teams/#list-teams
func (s *TeamService) List(ctx context.Context) ([]*Team, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	var (
		baseURL  = s.client.cfg.endpoint.String()
		endpoint strings.Builder
	)

	endpoint.Grow(len(baseURL) + len(teamPath))

	endpoint.WriteString(baseURL)
	endpoint.WriteString(teamPath)

	req, err := s.client.request(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	resp, err := s.client.do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	var result []*Team
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, resp, fmt.Errorf("%w", err)
	}

	return result, resp, nil
}
