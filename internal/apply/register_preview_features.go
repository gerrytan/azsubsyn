package apply

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armfeatures"
	"github.com/gerrytan/azdiffit/internal/config"
	"github.com/gerrytan/azdiffit/internal/credential"
	"github.com/gerrytan/azdiffit/internal/plan"
)

func registerPreviewFeatures(config *config.Config, previewFeatures []plan.PreviewFeature) error {
	if len(previewFeatures) == 0 {
		fmt.Printf("ℹ️  No preview feature registrations required\n")
		return nil
	}

	cred, err := credential.BuildCredential(config)
	if err != nil {
		return fmt.Errorf("failed to build credentials: %w", err)
	}

	client, err := armfeatures.NewClient(config.SubscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, feature := range previewFeatures {
		fmt.Printf("  - Registering Preview Feature: %s/%s (Reason: %s)\n", feature.Namespace, feature.Key, feature.Reason)

		_, err := client.Register(context.Background(), feature.Namespace, feature.Key, nil)
		if err != nil {
			fmt.Printf("   ❌ Failed to register preview feature %s/%s: %s\n", feature.Namespace, feature.Key, err)
		}
	}

	return nil
}
