package credential

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gerrytan/azsubsyn/internal/config"
)

func BuildCredential(config *config.Config) (cred *azidentity.ClientSecretCredential, err error) {
	return azidentity.NewClientSecretCredential(
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
}
