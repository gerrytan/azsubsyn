package credcheck

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/gerrytan/azdiffit/internal/credential"
)

func RunCredCheck() error {
	fmt.Println("🔍 Checking credentials and connectivity...")
	fmt.Println()

	srcConfig := credential.CredentialConfig{
		ClientID:       os.Getenv("AZDIFFIT_SRC_CLIENT_ID"),
		ClientSecret:   os.Getenv("AZDIFFIT_SRC_CLIENT_SECRET"),
		TenantID:       os.Getenv("AZDIFFIT_SRC_TENANT_ID"),
		SubscriptionID: os.Getenv("AZDIFFIT_SRC_SUBSCRIPTION_ID"),
	}

	targetConfig := credential.CredentialConfig{
		ClientID:       os.Getenv("AZDIFFIT_TARGET_CLIENT_ID"),
		ClientSecret:   os.Getenv("AZDIFFIT_TARGET_CLIENT_SECRET"),
		TenantID:       os.Getenv("AZDIFFIT_TARGET_TENANT_ID"),
		SubscriptionID: os.Getenv("AZDIFFIT_TARGET_SUBSCRIPTION_ID"),
	}

	allEnvVarsSet := true
	if !checkEnvVars(&srcConfig, "SRC") {
		allEnvVarsSet = false
	}
	if !checkEnvVars(&targetConfig, "TARGET") {
		allEnvVarsSet = false
	}

	if !allEnvVarsSet {
		return fmt.Errorf("❌ Missing required environment variables. Please set all required variables and try again")
	}

	fmt.Println("✅ All required environment variables are set")
	fmt.Println()

	ctx := context.Background()

	fmt.Println("🔐 Testing authentication and connectivity...")

	if err := testSubscriptionAccess(ctx, &srcConfig, "source"); err != nil {
		return fmt.Errorf("❌ Source subscription access failed: %w", err)
	}

	if err := testSubscriptionAccess(ctx, &targetConfig, "target"); err != nil {
		return fmt.Errorf("❌ Target subscription access failed: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ All credentials are valid and subscriptions are accessible!")
	fmt.Println("🎉 You're ready to run 'azdiffit plan' command")

	return nil
}

func checkEnvVars(config *credential.CredentialConfig, name string) bool {
	missing := []string{}

	if config.ClientID == "" {
		missing = append(missing, fmt.Sprintf("AZDIFFIT_%s_CLIENT_ID", name))
	}
	if config.ClientSecret == "" {
		missing = append(missing, fmt.Sprintf("AZDIFFIT_%s_CLIENT_SECRET", name))
	}
	if config.TenantID == "" {
		missing = append(missing, fmt.Sprintf("AZDIFFIT_%s_TENANT_ID", name))
	}
	if config.SubscriptionID == "" {
		missing = append(missing, fmt.Sprintf("AZDIFFIT_%s_SUBSCRIPTION_ID", name))
	}

	if len(missing) > 0 {
		fmt.Printf("❌ Missing %s subscription environment variables:\n", name)
		for _, envVar := range missing {
			fmt.Printf("   - %s\n", envVar)
		}
		fmt.Println()
		return false
	}

	fmt.Printf("✅ %s subscription environment variables are set\n", name)
	return true
}

func testSubscriptionAccess(ctx context.Context, config *credential.CredentialConfig, kind string) error {
	cred, err := azidentity.NewClientSecretCredential(
		config.TenantID,
		config.ClientID,
		config.ClientSecret,
		&azidentity.ClientSecretCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Retry: policy.RetryOptions{
					MaxRetries: 3,
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create credential for %s subscription: %w", kind, err)
	}

	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions client for %s: %w", kind, err)
	}

	sub, err := client.Get(ctx, config.SubscriptionID, nil)
	if err != nil {
		return fmt.Errorf("failed to access %s subscription %s: %w", kind, config.SubscriptionID, err)
	}

	if sub.Subscription.DisplayName == nil {
		return fmt.Errorf("%s subscription %s returned invalid response", kind, config.SubscriptionID)
	}

	fmt.Printf("✅ %s subscription access successful\n", kind)
	fmt.Printf("   - Subscription: %s (%s)\n", *sub.Subscription.DisplayName, config.SubscriptionID)
	if sub.Subscription.TenantID != nil {
		fmt.Printf("   - Tenant: %s\n", *sub.Subscription.TenantID)
	}
	if sub.Subscription.State != nil {
		fmt.Printf("   - State: %s\n", *sub.Subscription.State)
	}

	return nil
}
