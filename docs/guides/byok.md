---
page_title: "Using a Wallet and BYOK"
---

# Using a Wallet and BYOK
- [Prerequisites](#prerequisites)
- [Sample configuration](#sample-configuration)

<a name="prerequisites"></a>
## Prerequisites
At least one AWS account must be connected prior to using a Wallet and BYOK.

<a name="sample-configuration"></a>
## Sample configuration
This sample Terraform configuration highlights the creation of a Wallet resource and the
installation of the BYOK module.

```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.3.0"
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

# Install the Clumio BYOK template onto the registered wallet
module "clumio_byok" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                          = "clumio-code/byok-template/clumio"
  other_regions                   = setsubtract(clumio_wallet.wallet.supported_regions, toset([data.aws_region.current.name]))
  account_native_id               = clumio_wallet.wallet.account_native_id
  clumio_control_plane_account_id = clumio_wallet.wallet.clumio_control_plane_account_id
  clumio_account_id               = clumio_wallet.wallet.clumio_account_id
  token                           = clumio_wallet.wallet.token
}
```
