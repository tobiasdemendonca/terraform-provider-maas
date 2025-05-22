# Terraform provider for MAAS

This repository contains the source code for the Terraform provider for MAAS, which allows you to manage [MAAS](https://maas.io/) (Metal as a Service) resources using Terraform.

## Quick links

- [Our latest provider documentation (Canonical/MAAS)](https://registry.terraform.io/providers/canonical/maas/latest/docs) - the best place    for detailed information about what the provider does and how to use it.
- [Provider documentation for v2.2.0 and below (MAAS/MAAS)](https://registry.terraform.io/providers/maas/maas/latest/docs) - no longer maintained.
- [Development Guide](DEVELOPMENT.md)
- [Release process](RELEASING.md)
- [Changelog](CHANGELOG.md)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.4.x
- [Go](https://golang.org/doc/install) >= 1.23
- A MAAS installation

## Getting started

### Usage

The provider is available from the [Terraform Registry](https://registry.terraform.io/providers/canonical/maas/latest), and should be available to use when used in your Terraform configuration, without the need for any additional steps.

```hcl
terraform {
  required_providers {
      maas = {
        source = "canonical/maas"
        version = "~>2.0"
      }
    }
  }

provider "maas" {
  api_version = "2.0"
  api_key     = "YOUR MAAS API KEY"
  api_url     = "http://<MAAS_SERVER>[:MAAS_PORT]/MAAS"
  installation_method = "snap"
}

# <your MAAS Terraform configuration here>
```
### Provider configuration

The provider accepts the following config options:

- **api_key**: [MAAS API key](https://maas.io/docs/snap/3.0/cli/maas-cli#heading--log-in-required).
- **api_url**: URL for the MAAS API server (eg: <http://127.0.0.1:5240/MAAS>).
- **api_version**: MAAS API version used. It is optional and it defaults to `2.0`.
- **installation_method**: MAAS Installation method used. Optional, defaults to `snap`. Valid options: `snap`, and `deb`.


### Build from source

If you want to build the provider from source and (optionally) install it:

1. Clone the repository
2. Enter the repository directory
3. Build the provider with:

    ```sh
    make build
    ```

4. (Optional): Install the freshly built provider with:

    ```sh
    make install
    ```

## Contributing

If you're interested in contributing to the provider, please see the [Development Guide](DEVELOPMENT.md) for where to start.

## Testing

Unit tests run with every pull request and merge to master. The end to end tests run on a nightly basis against a hosted MAAS deployment, results can be found [here](https://raw.githubusercontent.com/canonical/maas-terraform-e2e-tests/main/results.json?token=GHSAT0AAAAAAB3FX6R5C67Q4LH7ADOO5O3IY4ODCNA) and are checked on each PR, with a warning if failed.

## :warning: Repository ownership and provider name change

The Terraform Provider for MAAS repository now lives under the [Canonical GitHub organisation](https://github.com/canonical) with a new name `github.com/canonical/terraform-provider-maas`.

Ensure you are pointing at the new provider name inside your Terraform module(s), which is `canonical/maas`:

1. Manually update the list of required providers in your Terraform module(s):

    ```diff
    terraform {
      required_providers {
        maas = {
    -     source  = "maas/maas"
    +     source  = "canonical/maas"
          version = "~>2.0"
        }
      }
    }
    ```

2. Upgrade your provider dependencies to add the `canonical/maas` provider info:

    ```bash
    terraform init -upgrade
    ```

3. Replace the provider reference in your state:

    ```bash
    terraform state replace-provider maas/maas canonical/maas
    ```

4. Upgrade your provider dependencies to remove the `maas/maas` provider info:

    ```bash
    terraform init -upgrade
    ```

References:

- <https://developer.hashicorp.com/terraform/language/files/dependency-lock#dependency-on-a-new-provider>
- <https://developer.hashicorp.com/terraform/language/files/dependency-lock#providers-that-are-no-longer-required>
- <https://developer.hashicorp.com/terraform/cli/commands/state/replace-provider>

---

## Additional resources

- [MAAS documentation](https://maas.io/docs)
- [MAAS Launchpad](https://launchpad.net/maas) and [MAAS GitHub mirror](https://github.com/canonical/maas)
- [Terraform documentation](https://www.terraform.io/docs)
