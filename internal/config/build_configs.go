package config

import (
	"fmt"
	"os"
)

func BuildConfigs() (srcConfig *Config, targetConfig *Config, err error) {
	srcConfig = &Config{
		ClientID:       os.Getenv("AZSUBSYN_SRC_CLIENT_ID"),
		ClientSecret:   os.Getenv("AZSUBSYN_SRC_CLIENT_SECRET"),
		TenantID:       os.Getenv("AZSUBSYN_SRC_TENANT_ID"),
		SubscriptionID: os.Getenv("AZSUBSYN_SRC_SUBSCRIPTION_ID"),
	}

	targetConfig = &Config{
		ClientID:       os.Getenv("AZSUBSYN_TARGET_CLIENT_ID"),
		ClientSecret:   os.Getenv("AZSUBSYN_TARGET_CLIENT_SECRET"),
		TenantID:       os.Getenv("AZSUBSYN_TARGET_TENANT_ID"),
		SubscriptionID: os.Getenv("AZSUBSYN_TARGET_SUBSCRIPTION_ID"),
	}

	missingEnvVars := []string{}
	missingEnvVars = append(missingEnvVars, checkEnvVar(srcConfig, "SRC")...)
	missingEnvVars = append(missingEnvVars, checkEnvVar(targetConfig, "TARGET")...)
	if len(missingEnvVars) > 0 {
		for _, varName := range missingEnvVars {
			fmt.Printf("  - ❌ Missing var: %s\n", varName)
		}
		return nil, nil, fmt.Errorf("Missing required environment variables")
	}

	return
}

func checkEnvVar(config *Config, name string) []string {
	missing := []string{}

	if config.ClientID == "" {
		missing = append(missing, "AZSUBSYN_"+name+"_CLIENT_ID")
	}
	if config.ClientSecret == "" {
		missing = append(missing, "AZSUBSYN_"+name+"_CLIENT_SECRET")
	}
	if config.TenantID == "" {
		missing = append(missing, "AZSUBSYN_"+name+"_TENANT_ID")
	}
	if config.SubscriptionID == "" {
		missing = append(missing, "AZSUBSYN_"+name+"_SUBSCRIPTION_ID")
	}

	return missing
}
