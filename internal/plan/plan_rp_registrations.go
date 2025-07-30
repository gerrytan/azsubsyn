package plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gerrytan/azdiffit/internal/config"
	"github.com/gerrytan/azdiffit/internal/credential"
	"github.com/gerrytan/azdiffit/internal/pointer"
	"github.com/gerrytan/azdiffit/internal/rp"
)

func planRPRegistrations() (rpRegs []RpRegistration, err error) {
	ctx := context.Background()

	srcConfig, targetConfig, err := config.BuildConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to build configurations: %w", err)
	}

	fmt.Println("🔍 Fetching resource providers from source subscription...")
	sourceRPs, err := getResourceProviders(ctx, srcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource providers from source subscription: %w", err)
	}

	fmt.Println("🔍 Fetching resource providers from target subscription...")
	targetRPs, err := getResourceProviders(ctx, targetConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource providers from target subscription: %w", err)
	}

	targetRPsByNamespace := make(map[string]*rp.ResourceProvider)
	for _, rp := range targetRPs {
		targetRPsByNamespace[rp.Namespace] = rp
	}

	for _, srcRp := range sourceRPs {
		if strings.EqualFold(srcRp.RegistrationState, "Registered") {
			targetRp, exists := targetRPsByNamespace[srcRp.Namespace]
			if !exists {
				rpRegs = append(rpRegs, RpRegistration{
					Namespace: srcRp.Namespace,
					Reason:    "NotFoundInTarget",
				})
			} else if !strings.EqualFold(targetRp.RegistrationState, "Registered") {
				rpRegs = append(rpRegs, RpRegistration{
					Namespace: targetRp.Namespace,
					Reason:    "NotRegisteredInTarget",
				})
			}
		}
	}

	return
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
