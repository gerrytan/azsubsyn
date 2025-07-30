package plan

import (
	"encoding/json"
	"fmt"
	"os"
)

func RunPlan() error {
	plan := Plan{}

	fmt.Println("📋 Creating RP registration plan...")

	rpRegs, err := planRPRegistrations()
	if err != nil {
		return fmt.Errorf("❌ Failed to plan RP registrations: %w", err)
	}
	plan.RpRegistrations = rpRegs

	jsonData, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("❌ Failed to serialize plan to JSON: %w", err)
	}

	err = os.WriteFile("azdiffit-plan.jsonc", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("❌ Failed to write plan to file: %w", err)
	}

	fmt.Printf("✅ Plan written successfully to azdiffit-plan.jsonc (%d RP registrations)\n", len(plan.RpRegistrations))
	return nil
}
