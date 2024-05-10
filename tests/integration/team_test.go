// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package integration_test

import (
	"context"
	"testing"

	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

func TestTeam(t *testing.T) {
	t.Parallel()

	var (
		client = xtesting.SetupLiveClient(t)
		ctx    = context.Background()
	)

	teams, _, err := client.Team.List(ctx)
	if err != nil {
		t.Fatalf("failed to get teams: %v", err)
	}

	if len(teams) == 0 {
		t.Fatalf("no teams found")
	}
}
