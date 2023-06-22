---
page_title: "Using Wallet and the BYOK Module"
---

# Using Wallet and the BYOK Module
- [Preparation](#preparation)
- [Sample Configuration](#sample)

The following is an example of how to instantiate a Clumio wallet and install the
[Clumio BYOK module](https://registry.terraform.io/modules/clumio-code/byok-template/clumio/latest)
to an AWS account. NOTE that the AWS account used does not have to be the same as an AWS account to
be protected. However, at least one AWS account must be connected with the
[Clumio AWS module](https://registry.terraform.io/modules/clumio-code/aws-template/clumio/latest)
prior to setting up a Clumio wallet and BYOK. The subsequent steps assume that such an AWS account
has already been setup.

<a name="preparation"></a>
## Preparation
Please see the "Getting Started" guide for notes about setting up a Clumio API key as well as
setting up AWS environment variables.

<a name="sample"></a>
## Sample Configuration
This sample configuration highlights the creation of a Clumio wallet and the installation of the
[Clumio BYOK module](https://registry.terraform.io/modules/clumio-code/byok-template/clumio/latest).
NOTE that if desired, an existing Multi-Region AWS CMK ID can be given.

```terraform
variable "existing_cmk_id" {
  description = "An existing CMK (if any) to use."
  type        = string
  default     = ""
}

terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.5.1"
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
  external_id       = random_uuid.external_id.id
  existing_cmk_id   = var.existing_cmk_id
}
```
