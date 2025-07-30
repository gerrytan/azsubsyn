# azdiffit

A CLI tool to setup a target Azure subscription based on a source subscription (which can be on a separate tenant). The
setup includes ensuring the target subscription has all Resource Providers (RPs) registered, all preview features
registered, and enough quotas.

## Installation

Install using go toolchain: `go install github.com/gerrytan/azdiffit@latest` or alternatively download precompiled
binaries for your OS from the Release section and place it in /usr/local/bin or other places registered in your PATH
variable.

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

The steps to create the service principal for source and target subscriptions are almost identical. You need to have [Azure
CLI](https://learn.microsoft.com/cli/azure/install-azure-cli) installed:

1. Ensure you're logged in to the correct tenant and subscription: logout, login, set subscription and check session as
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

`azdiffit plan` will fetch RP registrations, Preview features and Quotas for both source and target subscriptions,
create a modification plan to be applied to the target subscription, and save it to `azdiffit-plan.jsonc` in the working
directory.

The plan file has following format:

```jsonc
{
  "rpRegistrations": [
    { "namespace": "Microsoft.Foo", "reason": "NotRegisteredInTarget" },
    { "namespace": "Microsoft.Bar", "reason": "NotFoundInTarget" }
  ]
}

```

The modification is always additive, if target subscription already has an RP registered / more than required quota, it
won't be turned off / reduced.

The plan file can be modified manually if necessary.

### Apply

`azdiffit apply -f azdiffit-plan.jsonc` will execute the modification plan as per the supplied plan file.
