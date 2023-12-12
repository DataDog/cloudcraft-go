package integration_test

import (
	"bytes"
	"context"
	"image/png"
	"testing"

	"github.com/DataDog/cloudcraft-go"
	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

func TestBlueprint(t *testing.T) {
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

	blueprints, _, err := client.Blueprint.List(ctx)
	if err != nil {
		t.Fatalf("failed to list blueprints: %v", err)
	}

	if len(blueprints) == 0 {
		t.Fatalf("no blueprints found")
	}

	var (
		blueprintID   = blueprint.ID
		blueprintName = blueprint.Name
	)

	give = &cloudcraft.Blueprint{
		ID:   blueprintID,
		Name: blueprintName + _updatedSuffix,
		Data: &cloudcraft.BlueprintData{
			Name: blueprintName + _updatedSuffix,
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

	blueprint, _, err = client.Blueprint.Get(ctx, blueprintID)
	if err != nil {
		t.Fatalf("failed to get blueprint: %v", err)
	}

	if blueprint.Name != give.Name {
		t.Fatalf("blueprint name not updated")
	}

	image, _, err := client.Blueprint.ExportImage(ctx, blueprintID, "png", nil)
	if err != nil {
		t.Fatalf("failed to export blueprint image: %v", err)
	}

	imageData, err := png.DecodeConfig(bytes.NewReader(image))
	if err != nil {
		t.Fatalf("failed to decode blueprint image: %v", err)
	}

	if imageData.Width != 1920 || imageData.Height != 1080 {
		t.Fatalf("unexpected image size: %dx%d", imageData.Width, imageData.Height)
	}

	_, err = client.Blueprint.Delete(ctx, blueprintID)
	if err != nil {
		t.Fatalf("failed to delete blueprint: %v", err)
	}
}
