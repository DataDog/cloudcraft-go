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

const _envAWSRoleARN string = "CLOUDCRAFT_TEST_AWS_ROLE_ARN"

func TestAWS(t *testing.T) {
	t.Parallel()

	var (
		client = xtesting.SetupLiveClient(t)
		arn    = xtesting.GetEnv(t, _envAWSRoleARN)
		ctx    = context.Background()
	)

	give := &cloudcraft.AWSAccount{
		Name:    xtesting.UniqueName(t),
		RoleARN: arn,
	}

	account, _, err := client.AWS.Create(ctx, give)
	if err != nil {
		t.Fatalf("failed to create AWS account: %v", err)
	}

	if account == nil {
		t.Fatalf("AWS account is nil")
	}

	accounts, _, err := client.AWS.List(ctx)
	if err != nil {
		t.Fatalf("failed to list AWS accounts: %v", err)
	}

	if len(accounts) == 0 {
		t.Fatalf("no AWS accounts found")
	}

	var (
		accountID   = account.ID
		accountName = account.Name
	)

	give = &cloudcraft.AWSAccount{
		ID:      accountID,
		Name:    accountName + " (Updated)",
		RoleARN: arn,
	}

	_, err = client.AWS.Update(ctx, give)
	if err != nil {
		t.Fatalf("failed to update AWS account: %v", err)
	}

	snapshot, _, err := client.AWS.Snapshot(ctx, accountID, "us-east-1", "", nil)
	if err != nil {
		t.Fatalf("failed to snapshot AWS account: %v", err)
	}

	snapshotData, err := png.DecodeConfig(bytes.NewReader(snapshot))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	if snapshotData.Width != 1920 || snapshotData.Height != 1080 {
		t.Fatalf("unexpected snapshot size: %dx%d", snapshotData.Width, snapshotData.Height)
	}

	_, err = client.AWS.Delete(ctx, accountID)
	if err != nil {
		t.Fatalf("failed to delete AWS account: %v", err)
	}

	parameters, _, err := client.AWS.IAMParameters(ctx)
	if err != nil {
		t.Fatalf("failed to get IAM parameters: %v", err)
	}

	if parameters == nil {
		t.Fatalf("IAM parameters are nil")
	}

	policies, _, err := client.AWS.IAMPolicy(ctx)
	if err != nil {
		t.Fatalf("failed to get IAM policies: %v", err)
	}

	if policies == nil {
		t.Fatalf("IAM policies are nil")
	}
}
