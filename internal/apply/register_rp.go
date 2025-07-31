package apply

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gerrytan/azdiffit/internal/config"
	"github.com/gerrytan/azdiffit/internal/credential"
	"github.com/gerrytan/azdiffit/internal/plan"
	"github.com/gerrytan/azdiffit/internal/pointer"
)

func registerRPs(config *config.Config, rpRegistrations []plan.RpRegistration) error {
	if len(rpRegistrations) == 0 {
		fmt.Printf("ℹ️  No resource provider registrations required\n")
		return nil
	}

	cred, err := credential.BuildCredential(config)
	if err != nil {
		return fmt.Errorf("failed to build credentials: %w", err)
	}

	providersClient, err := armresources.NewProvidersClient(config.SubscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create providers client: %w", err)
	}

	for _, rpReg := range rpRegistrations {
		fmt.Printf("  - Registering RP: %s (Reason: %s)\n", rpReg.Namespace, rpReg.Reason)

		_, err := providersClient.Register(context.Background(), rpReg.Namespace, &armresources.ProvidersClientRegisterOptions{
			Properties: &armresources.ProviderRegistrationRequest{
				ThirdPartyProviderConsent: &armresources.ProviderConsentDefinition{
					ConsentToAuthorization: pointer.To(true),
				},
			},
		})
		if err != nil {
			fmt.Printf("   ❌ Failed to register RP %s: %s\n", rpReg.Namespace, err)
		}

	}

	return nil
}
