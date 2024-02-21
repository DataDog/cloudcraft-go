package integration_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/DataDog/cloudcraft-go"
	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

func TestRetries(t *testing.T) {
	t.Parallel()

	var (
		client = xtesting.SetupLiveClient(t)
		ctx    = context.Background()
	)

	give := &cloudcraft.Blueprint{
		Data: &cloudcraft.BlueprintData{
			Name: xtesting.UniqueName(t),
		},
	}

	blueprint, _, err := client.Blueprint.Create(ctx, give)
	if err != nil {
		t.Fatalf("failed to create blueprint: %v", err)
	}

	if blueprint == nil {
		t.Fatalf("blueprint is nil")
	}

	var (
		blueprintID   = blueprint.ID
		blueprintName = blueprint.Name
	)

	for i := 0; i < 10; i++ {
		iStr := strconv.Itoa(i)

		give = &cloudcraft.Blueprint{
			ID:   blueprintID,
			Name: blueprintName + " (" + iStr + ")",
			Data: &cloudcraft.BlueprintData{
				Name: blueprintName + " (" + iStr + ")",
				Nodes: []map[string]any{
					{
						"id":           "98172baa-a059-4b04-832d-8a7f5d14b595",
						"type":         "ec2",
						"region":       "us-east-1",
						"platform":     "linux",
						"instanceType": "m5",
						"instanceSize": "large",
					},
				},
			},
		}

		_, err = client.Blueprint.Update(ctx, give, "")
		if err != nil {
			t.Fatalf("failed to update blueprint: %v", err)
		}
	}

	blueprint, _, err = client.Blueprint.Get(ctx, blueprintID)
	if err != nil {
		t.Fatalf("failed to get blueprint: %v", err)
	}

	if blueprint.Name != give.Name {
		t.Fatalf("blueprint name not updated; possibly due to error in retry logic: wanted %q, got %q", give.Name, blueprint.Name)
	}

	_, err = client.Blueprint.Delete(ctx, blueprintID)
	if err != nil {
		t.Fatalf("failed to delete blueprint: %v", err)
	}
}
