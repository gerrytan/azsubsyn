# azdiffit

A CLI tool to ensure a target Azure subscription has all RPs (resource providers) and preview features registered
compared to source (which can be on a different tenant).

## Prerequisites and installation

You need to have [go](https://go.dev/doc/install) and [Azure CLI](https://learn.microsoft.com/cli/azure/install-azure-cli) installed.

Install via go toolchain: `go install github.com/gerrytan/azdiffit@latest`. The binary will be available in
`$GOPATH/bin/azdiffit`.

## Usage

### Credentials setup

This tool uses a service principal with password based authentication. You need 'Owner' role in both source and target
subscription to setup the principals.

The following environment variables need to be set:

```bash
export AZDIFFIT_SRC_CLIENT_ID="" # See instruction below
export AZDIFFIT_SRC_CLIENT_SECRET="" # See instruction below
export AZDIFFIT_SRC_TENANT_ID="12345678-1234-1234-1234-123456789abc"
export AZDIFFIT_SRC_SUBSCRIPTION_ID="12345678-1234-1234-1234-123456789abc"
export AZDIFFIT_TARGET_CLIENT_ID="" # See instruction below
export AZDIFFIT_TARGET_CLIENT_SECRET="" # See instruction below
export AZDIFFIT_TARGET_TENANT_ID="12345678-1234-1234-1234-123456789abc"
export AZDIFFIT_TARGET_SUBSCRIPTION_ID="12345678-1234-1234-1234-123456789abc"
```

The steps to create the service principal for source and target subscriptions are almost identical:

1. Ensure you're logged in to the correct tenant and subscription. Logout, login, set subscription and check session as
   required:

    ```bash
    az logout
    az login -t mytenant.onmicrosoft.com
    az account set --subscription "My subscription name"
    az account show
    ```

1. Create the service principal:

    ```bash
    az ad sp create-for-rbac \
      --name myname-azdiffit-20250701 \
      --role Contributor \
      --scopes /subscriptions/12345678-1234-1234-1234-123456789abc \
      --years 1
    ```

1. Upon successful creation you will get a response like this:

    ```json
    {
      "appId": "12345678-1234-1234-1234-123456789abc",
      "displayName": "myuser-azdiffit-20250701",
      "password": "abcdefghijkl",
      "tenant": "12345678-1234-1234-1234-123456789abc"
    }
    ```

    Use the `"appId"` value as `_CLIENT_ID` env var and `"password"` as `_CLIENT_SECRET`

1. Once you've exported the environment variables, check using `azdiffit credcheck`

### Plan

`azdiffit plan` will fetch RP and preview features registrations for both source and target subscriptions and creates a
modification plan to be applied to the target subscription. The plan is saved to the `azdiffit-plan.jsonc` file in the
working directory.

The plan file has following format:

```jsonc
{
  "rpRegistrations": [
    { "namespace": "Microsoft.CertificateRegistration", "reason": "NotRegisteredInTarget" },
    { "namespace": "Microsoft.VideoIndexer", "reason": "NotFoundInTarget" }
  ],
  "previewFeatures": [
    {
      "key": "AllowMultiplePeeringLinksBetweenVnets",
      "namespace": "Microsoft.Network",
      "reason": "NotRegisteredInTarget"
    },
    {
      "key": "locationCapability",
      "namespace": "Microsoft.DBforPostgreSQL",
      "reason": "NotFoundInTarget"
    }
  ]
}

```

The modification is always additive, if target subscription already has an RP / feature registered, it won't be turned
off.

The plan file can be modified manually if necessary.

### Apply

`azdiffit apply azdiffit-plan.jsonc` will execute the modification plan as per the supplied file.
