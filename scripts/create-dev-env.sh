#!/bin/bash

ENV_DIR=$PWD/.devenv
mkdir -p $ENV_DIR
cat << EOF > $ENV_DIR/main.tf
terraform {
  required_providers {
    maas = {
      source  = "registry.terraform.io/canonical/maas"
      version = "~> 2.6.0"
    }
  }
}

provider "maas" {
  api_version         = "2.0"
  api_key             = ... # Fill me in from MAAS
  api_url             = ... # Fill me in from MAAS
  installation_method = "snap"
}

# Try me out!
# resource "maas_fabric" "myfabric" {
#   name = "myfabric"
# }

EOF

echo ""
echo "A development directory has been created at $ENV_DIR."
echo "This is a space to get started with during development. Fill in the relevant fields in the provider block from MAAS, and you should be able to start terraforming in MAAS."
echo ""
echo "Happy Terraforming! 🚀"
echo ""
