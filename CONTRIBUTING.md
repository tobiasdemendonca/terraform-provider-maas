# Contributing to the Terraform MAAS Provider

Thank you for your interest in contributing to the Terraform MAAS Provider! We appreciate your help in making this project better. This document provides information on how to set up your development environment and best practices for contributing to the project.


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.4.x
- [Go](https://golang.org/doc/install) >= 1.23
- A MAAS installation running. See the [maas-dev-setup](https://github.com/canonical/maas-dev-setup) repository for more information on a development setup.
- [CLA](https://ubuntu.com/legal/contributors) signed with the email used with git and GitHub.

## Branching strategy

This project follows a fork-based development model with a single long-running master branch. All contributions should be made via pull requests (PRs) from forked repositories.

1. Fork the Repository:
    1. Go to the repository on GitHub and click "Fork" in the top-right corner.
    1. Clone your fork locally:
       ```bash
       git clone <https-or-ssh-url-to-your-fork>
       cd terraform-provider-maas
       ```
    1. Add the upstream repository (the original repo) as a second remote:
       ```bash
       git remote add upstream <https-or-ssh-url-to-original>
       ```
1. Create a Feature Branch:
    ```bash
    git checkout -b feat/feature-name
    ```
1. Keep Your Branch Up to Date:
    1. Before working, sync your branch with the latest changes from master:
       ```bash
       git fetch upstream
       git checkout master
       git merge upstream/master
       ```
    1. Then, rebase or merge your feature branch if necessary:
        ```bash
        git checkout feat/feature-name
        git rebase master
        ```
1. Commit and Push Changes:
    1. Follow commit message guidelines (e.g., fix: correct typo in readme).
    1. Push your branch to your forked repository:
        ```bash
        git push origin feat/feature-name
        ```
1. Submit a Pull Request:
    1. Go to the your forked repository on GitHub.
    1. Click "New Pull Request". Select your feature branch to merge from your forked repo, into the master branch of the original repo.
1. Address Review Feedback. Once approved, a maintainer will merge your PR. 🎉

## Commit messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification. Conventional Commits defines the following structure for the Git commit message:

```bash
<type>[scope][!]: <description>

[body]

[footer(s)]
```

Where
- `type` is the kind of the change (e.g. feature, bug fix, documentation change, refactor).
- `scope` may be used to provide additional contextual information (e.g. which system component is affected). If scope is provided, it’s enclosed in parentheses.
- `!` MUST be added if commit introduces a breaking change.
- `description` is a brief summary of a change (try to keep it short, so overall title no more than 72 characters).
- `footer` is detailed information about the change (e.g. breaking change, related bugs, etc.).


## Create a development environment
### Setup
1. Run `make build` to build the provider binary locally, located in `./bin`.
1. Run `make create-dev-overrides` and follow any output instructions. More info [here](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers).
1. Run `make dev-env` to create a development directory and follow any output instructions.
1. In your dev-env directory, use these commands to get started:
   1. Run `terraform fmt` to format the `main.tf` file.
   1. Run `terraform init` to initialize the provider.
   1. Run `terraform plan` to see the changes that will be applied.
   1. Run `terraform apply` to apply the changes. These should be reflected in the MAAS environment.
   1. Run `terraform destroy` to destroy the resources.

### Workflow
Assuming you have already setup dev-overrides:
1. Make a change to the provider.
1. Rebuild the provider binary locally with `make build`.
1. In your dev-env directory, you can immediately run `terraform apply` with your new changes.

## Testing
Tests are written as advised in the [Terraform docs](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing). They are split into unit tests and acceptance tests, with the latter creating real resources in the MAAS environment. Therefore, you will need to ensure MAAS is running locally for these tests to pass.

To run the tests:
1. Ensure MAAS_API_KEY and MAAS_API_URL environment variables are set in your shell, relevant to your running MAAS installation.

2. Run the tests:
    - Run the unit tests:
        ```bash
        make test
        ```
    - Run both the unit tests and all Terraform acceptance tests:
        ```bash
        make testacc
        ```
        > [!NOTE]
        > You may need to set specific environment variables for some tests to pass, for example machine ids. Add these to your `env.sh` file before sourcing it again, if required:
        > ```bash
        > export TF_ACC_NETWORK_INTERFACE_MACHINE=<system_id>   # b68rn4
        > export TF_ACC_TAG_MACHINES=<system_id>                # b68rn4
        > export TF_ACC_VM_HOST_ID=<system_id>                  # maas-host
        > export TF_ACC_BLOCK_DEVICE_MACHINE=<system_id>        # b68rn4
        > export TF_ACC_RACK_CONTROLLER_HOSTNAME=<name>         # maas-dev
        > ```
    - Run a specific acceptance test:
        ```bash
        make testacc TESTARGS="-run=TestAcc<resource_name>_basic"
        ```
    - Run a group of acceptance tests with names matching the regex:
        ```bash
        make testacc TESTARGS="-run=TestAcc<resource_name>"
        ```

## Getting Help

Check for existing issues [here](https://github.com/canonical/terraform-provider-maas/issues), or open a new one for bugs and feature requests.

## Release Process

Releases are handled by the maintainers, see [RELEASING.md](RELEASING.md).

## Additional Resources

- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)
- [Go Documentation](https://golang.org/doc/)
- [MAAS API Documentation](https://maas.io/docs/api)
