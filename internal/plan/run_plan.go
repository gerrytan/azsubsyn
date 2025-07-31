package plan

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gerrytan/azdiffit/internal/config"
)

func RunPlan() error {
	if len(os.Args) > 2 {
		printUsage()
		os.Exit(1)
	}

	srcConfig, targetConfig, err := config.BuildConfigs()
	if err != nil {
		return fmt.Errorf("❌ Failed to build configurations: %w", err)
	}

	fmt.Printf("🔄 Creating plan from source and target subscription...\n")
	fmt.Printf("  - Source tenant / sub: %s / %s\n", srcConfig.TenantID, srcConfig.SubscriptionID)
	fmt.Printf("  - Target tenant / sub: %s / %s\n", targetConfig.TenantID, targetConfig.SubscriptionID)

	plan := Plan{}

	fmt.Println("📋 Creating RP registration plan...")

	rpRegs, err := planRPRegistrations(srcConfig, targetConfig)
	if err != nil {
		return fmt.Errorf("❌ Failed to plan RP registrations: %w", err)
	}
	plan.RpRegistrations = rpRegs

	fmt.Println("📋 Creating preview features plan...")

	previewFeatures, err := planPreviewFeatures(srcConfig, targetConfig)
	if err != nil {
		return fmt.Errorf("❌ Failed to plan preview features: %w", err)
	}
	plan.PreviewFeatures = previewFeatures

	jsonData, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("❌ Failed to serialize plan to JSON: %w", err)
	}

	err = os.WriteFile("azdiffit-plan.jsonc", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("❌ Failed to write plan to file: %w", err)
	}

	fmt.Printf("✅ Plan written successfully to azdiffit-plan.jsonc (%d RPs, %d preview features)\n", len(plan.RpRegistrations), len(plan.PreviewFeatures))
	return nil
}

func printUsage() {
	fmt.Println("azdiffit plan - Scan unregistered RPs and preview feature in the target subscription and save the plan to a file")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  azdiffit plan")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Fetch RP and preview features registrations for both source and target subscriptions and creates a")
	fmt.Println("  modification plan to be applied to the target subscription. The plan is saved to `azdiffit-plan.jsonc` file in the working")
	fmt.Println("  directory.")
	fmt.Println()
	fmt.Println("  The modification is always additive, if target subscription already has an RP / feature registered, it won't be turned off.")
	fmt.Println()
	fmt.Println("  The plan file can be modified manually if necessary.")
}
