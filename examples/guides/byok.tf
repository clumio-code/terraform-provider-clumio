terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.4.0"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate the AWS provider
provider "aws" {
  region = "us-west-2"
}

# Retrieve the effective AWS account ID and region
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Register a new Clumio wallet for the effective AWS account ID and region
resource "clumio_wallet" "wallet" {
  account_native_id = data.aws_caller_identity.current.account_id
}

# Use a random string as an external ID for the role that manages the BYOK
resource "random_uuid" "external_id" {
}

# Install the Clumio BYOK template onto the registered wallet
module "clumio_byok" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source            = "clumio-code/byok-template/clumio"
  account_native_id = clumio_wallet.wallet.account_native_id
  token             = clumio_wallet.wallet.token
  clumio_account_id = clumio_wallet.wallet.clumio_account_id
  external_id       = var.external_id != "" ? var.external_id : random_uuid.external_id.id
  existing_cmk_id   = var.existing_cmk_id
}
