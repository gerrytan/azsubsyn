package plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gerrytan/azsubsyn/internal/config"
	"github.com/gerrytan/azsubsyn/internal/credential"
	"github.com/gerrytan/azsubsyn/internal/pointer"
)

func planRPRegistrations(srcConfig *config.Config, targetConfig *config.Config) (rpRegs []RpRegistration, err error) {
	ctx := context.Background()

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

	targetRPsByNamespace := make(map[string]*armresources.Provider)
	for _, rp := range targetRPs {
		targetRPsByNamespace[pointer.From(rp.Namespace)] = rp
	}

	for _, srcRp := range sourceRPs {
		if strings.EqualFold(pointer.From(srcRp.RegistrationState), "Registered") {
			targetRp, exists := targetRPsByNamespace[pointer.From(srcRp.Namespace)]
			if !exists {
				rpRegs = append(rpRegs, RpRegistration{
					Namespace: pointer.From(srcRp.Namespace),
					Reason:    "NotFoundInTarget",
				})
			} else if !strings.EqualFold(pointer.From(targetRp.RegistrationState), "Registered") {
				rpRegs = append(rpRegs, RpRegistration{
					Namespace: pointer.From(targetRp.Namespace),
					Reason:    "NotRegisteredInTarget",
				})
			}
		}
	}

	return
}

func getResourceProviders(ctx context.Context, config *config.Config) (rps []*armresources.Provider, err error) {
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
			rps = append(rps, provider)
		}
	}

	return
}
