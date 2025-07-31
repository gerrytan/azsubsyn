package credential

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/gerrytan/azdiffit/internal/config"
)

func RunCredCheck() error {
	if len(os.Args) > 2 {
		printUsage()
		os.Exit(1)
	}

	fmt.Println("🔍 Checking environment variables...")

	srcConfig, targetConfig, err := config.BuildConfigs()
	if err != nil {
		return fmt.Errorf("❌ Failed to build configurations: %w", err)
	}

	ctx := context.Background()

	fmt.Println("🔐 Testing authentication and connectivity...")

	if err := testSubscriptionAccess(ctx, srcConfig, "source"); err != nil {
		return fmt.Errorf("❌ Source subscription access failed: %w", err)
	}

	if err := testSubscriptionAccess(ctx, targetConfig, "target"); err != nil {
		return fmt.Errorf("❌ Target subscription access failed: %w", err)
	}

	fmt.Println("✅ All credentials are valid and subscriptions are accessible!")
	fmt.Println("🎉 You're ready to run 'azdiffit plan' command")

	return nil
}

func testSubscriptionAccess(ctx context.Context, config *config.Config, kind string) error {
	cred, err := BuildCredential(config)
	if err != nil {
		return fmt.Errorf("failed to build %s credential: %w", kind, err)
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

func printUsage() {
	fmt.Println("azdiffit credcheck - Check credentials and connectivity to both source and target subscriptions")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  azdiffit credcheck")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  This command checks if the environment variables are set up correctly,")
	fmt.Println("  tests authentication and connectivity to both source and target Azure subscriptions.")
}
