// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package integration_test

import (
	"bytes"
	"context"
	"image/png"
	"testing"

	"github.com/DataDog/cloudcraft-go"
	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

const (
	_envAzureApplicationID  string = "CLOUDCRAFT_TEST_AZURE_APPLICATION_ID"
	_envAzureDirectoryID    string = "CLOUDCRAFT_TEST_AZURE_DIRECTORY_ID"
	_envAzureSubscriptionID string = "CLOUDCRAFT_TEST_AZURE_SUBSCRIPTION_ID"
	_envAzureClientSecret   string = "CLOUDCRAFT_TEST_AZURE_CLIENT_SECRET"
)

func TestAzure(t *testing.T) {
	t.Parallel()

	var (
		client         = xtesting.SetupLiveClient(t)
		ctx            = context.Background()
		appID          = xtesting.GetEnv(t, _envAzureApplicationID)
		directoryID    = xtesting.GetEnv(t, _envAzureDirectoryID)
		subscriptionID = xtesting.GetEnv(t, _envAzureSubscriptionID)
		clientSecret   = xtesting.GetEnv(t, _envAzureClientSecret)
	)

	give := &cloudcraft.AzureAccount{
		Name:           xtesting.UniqueName(t),
		ApplicationID:  appID,
		DirectoryID:    directoryID,
		SubscriptionID: subscriptionID,
		ClientSecret:   clientSecret,
	}

	account, _, err := client.Azure.Create(ctx, give)
	if err != nil {
		t.Fatalf("failed to create Azure account: %v", err)
	}

	if account == nil {
		t.Fatalf("Azure account is nil")
	}

	accounts, _, err := client.Azure.List(ctx)
	if err != nil {
		t.Fatalf("failed to list Azure accounts: %v", err)
	}

	if len(accounts) == 0 {
		t.Fatalf("no Azure accounts found")
	}

	var (
		accountID   = account.ID
		accountName = account.Name
	)

	give = &cloudcraft.AzureAccount{
		ID:             accountID,
		Name:           accountName + " (Updated)",
		ApplicationID:  appID,
		DirectoryID:    directoryID,
		SubscriptionID: subscriptionID,
		ClientSecret:   clientSecret,
	}

	_, err = client.Azure.Update(ctx, give)
	if err != nil {
		t.Fatalf("failed to update Azure account: %v", err)
	}

	snapshot, _, err := client.Azure.Snapshot(ctx, accountID, "brazilsouth", "", nil)
	if err != nil {
		t.Fatalf("failed to snapshot Azure account: %v", err)
	}

	snapshotData, err := png.DecodeConfig(bytes.NewReader(snapshot))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	if snapshotData.Width != 1920 || snapshotData.Height != 1080 {
		t.Fatalf("unexpected snapshot size: %dx%d", snapshotData.Width, snapshotData.Height)
	}

	_, err = client.Azure.Delete(ctx, accountID)
	if err != nil {
		t.Fatalf("failed to delete Azure account: %v", err)
	}
}
