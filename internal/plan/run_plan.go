package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gerrytan/azdiffit/internal/config"
	"github.com/gerrytan/azdiffit/internal/credential"
	"github.com/gerrytan/azdiffit/internal/pointer"
	"github.com/gerrytan/azdiffit/internal/rp"
)

func RunPlan() error {
	plan := Plan{}

	fmt.Println("📋 Creating RP registration plan...")

	srcConfig, targetConfig, err := config.BuildConfigs()
	if err != nil {
		return fmt.Errorf("❌ Failed to build configurations: %w", err)
	}

	ctx := context.Background()

	fmt.Println("🔍 Fetching resource providers from source subscription...")
	sourceRPs, err := getResourceProviders(ctx, srcConfig)
	if err != nil {
		return fmt.Errorf("❌ Failed to get resource providers from source subscription: %w", err)
	}

	fmt.Println("🔍 Fetching resource providers from target subscription...")
	targetRPs, err := getResourceProviders(ctx, targetConfig)
	if err != nil {
		return fmt.Errorf("❌ Failed to get resource providers from target subscription: %w", err)
	}

	targetRPsByNamespace := make(map[string]*rp.ResourceProvider)
	for _, rp := range targetRPs {
		targetRPsByNamespace[rp.Namespace] = rp
	}

	for _, srcRp := range sourceRPs {
		if strings.EqualFold(srcRp.RegistrationState, "Registered") {
			targetRp, exists := targetRPsByNamespace[srcRp.Namespace]
			if !exists {
				fmt.Printf("Will be registered in target: Namespace: %s, Reason: NotFoundInTarget\n", srcRp.Namespace)
				plan.RpRegistrations = append(plan.RpRegistrations, RpRegistration{
					Namespace: srcRp.Namespace,
					Reason:    "NotFoundInTarget",
				})
				continue
			}

			if !strings.EqualFold(targetRp.RegistrationState, "Registered") {
				fmt.Printf("Will be registered in target: Namespace: %s, Reason: NotRegisteredInTarget\n", targetRp.Namespace)
				plan.RpRegistrations = append(plan.RpRegistrations, RpRegistration{
					Namespace: targetRp.Namespace,
					Reason:    "NotRegisteredInTarget",
				})
			}
		}
	}

	jsonData, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("❌ Failed to serialize plan to JSON: %w", err)
	}

	err = os.WriteFile("azdiffit-plan.jsonc", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("❌ Failed to write plan to file: %w", err)
	}

	fmt.Printf("✅ Plan written successfully to azdiffit-plan.jsonc (%d resource provider registrations)\n", len(plan.RpRegistrations))
	return nil
}

func getResourceProviders(ctx context.Context, config *config.Config) (rps []*rp.ResourceProvider, err error) {
	cred, err := credential.BuildCredential(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	client, err := armresources.NewProvidersClient(config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource providers client: %w", err)
	}

	pager := client.NewListPager(&armresources.ProvidersClientListOptions{
		Expand: nil,
	})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource providers page: %w", err)
		}

		for _, provider := range page.Value {
			rps = append(rps, &rp.ResourceProvider{
				Id:                 pointer.From(provider.ID),
				Namespace:          pointer.From(provider.Namespace),
				RegistrationState:  pointer.From(provider.RegistrationState),
				RegistrationPolicy: pointer.From(provider.RegistrationPolicy),
			})
		}
	}

	return
}
