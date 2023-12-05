// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package integration_test

import (
	"context"
	"testing"

	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

func TestUser(t *testing.T) {
	t.Parallel()

	var (
		client = xtesting.SetupLiveClient(t)
		ctx    = context.Background()
	)

	user, _, err := client.User.Me(ctx)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if user == nil {
		t.Fatalf("user is nil")
	}
}
