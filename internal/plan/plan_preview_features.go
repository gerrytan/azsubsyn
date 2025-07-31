package plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armfeatures"
	"github.com/gerrytan/azsubsyn/internal/config"
	"github.com/gerrytan/azsubsyn/internal/credential"
	"github.com/gerrytan/azsubsyn/internal/pointer"
)

func planPreviewFeatures(srcConfig *config.Config, targetConfig *config.Config) (prFeats []PreviewFeature, err error) {
	ctx := context.Background()

	fmt.Println("🔍 Fetching preview features from source subscription...")
	sourceFeatures, err := getPreviewFeatures(ctx, srcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get preview features from source subscription: %w", err)
	}

	fmt.Println("🔍 Fetching preview features from target subscription...")
	targetFeatures, err := getPreviewFeatures(ctx, targetConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get preview features from target subscription: %w", err)
	}

	// example name: "Microsoft.DevAI/Dev"
	targetFeaturesByName := make(map[string]*armfeatures.FeatureResult)
	for _, feat := range targetFeatures {
		targetFeaturesByName[pointer.From(feat.Name)] = feat
	}

	for _, srcFeature := range sourceFeatures {
		if strings.EqualFold(getState(srcFeature), "Registered") {
			targetFeature, exists := targetFeaturesByName[pointer.From(srcFeature.Name)]
			if !exists {
				srcKey, srcNamespace := parseKeyAndNamespace(srcFeature.Name)
				prFeats = append(prFeats, PreviewFeature{
					Key:       srcKey,
					Namespace: srcNamespace,
					Reason:    "NotFoundInTarget",
				})
			} else if !strings.EqualFold(getState(targetFeature), "Registered") {
				targetKey, targetNamespace := parseKeyAndNamespace(targetFeature.Name)
				prFeats = append(prFeats, PreviewFeature{
					Key:       targetKey,
					Namespace: targetNamespace,
					Reason:    "NotRegisteredInTarget",
				})
			}
		}
	}

	return
}

func getPreviewFeatures(ctx context.Context, config *config.Config) (features []*armfeatures.FeatureResult, err error) {
	cred, err := credential.BuildCredential(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	client, err := armfeatures.NewClient(config.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create features client: %w", err)
	}

	pager := client.NewListAllPager(&armfeatures.ClientListAllOptions{})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get features page: %w", err)
		}

		for _, feat := range page.Value {
			features = append(features, feat)
		}
	}

	return
}

func getState(f *armfeatures.FeatureResult) string {
	if f.Properties == nil {
		return ""
	}
	return pointer.From(f.Properties.State)
}

func parseKeyAndNamespace(name *string) (key, namespace string) {
	if name == nil || *name == "" {
		panic("empty feature name")
	}

	parts := strings.SplitN(*name, "/", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("bad format: %q", *name))
	}

	return parts[1], parts[0]
}
