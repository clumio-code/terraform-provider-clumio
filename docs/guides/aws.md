---
page_title: "Using Connections and the AWS Module"
---

# Using Connections and the AWS Module
- [Preparation](#preparation)
- [Basic, One Connection](#basic)
- [Cross-Region, Two Connections](#cross-region)
- [Cross-Account, Role Assumption, Two Connections](#cross-account)

The following are examples of various ways to instantiate Clumio connections and install the
[Clumio AWS module](https://registry.terraform.io/modules/clumio-code/aws-template/clumio/latest) to
one or more AWS accounts and regions to be protected.

<a name="preparation"></a>
## Preparation
Please see the "Getting Started" guide for notes about setting up a Clumio API key as well as
setting up AWS environment variables.

<a name="basic"></a>
## Basic, One Connection
The following configuration sets up a single Clumio connection and installs the
[Clumio AWS module](https://registry.terraform.io/modules/clumio-code/aws-template/clumio/latest) to
the AWS account and region to be protected.

```terraform
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

# Register a new Clumio connection for the effective AWS account ID and region
resource "clumio_aws_connection" "connection" {
  account_native_id = data.aws_caller_identity.current.account_id
  aws_region        = data.aws_region.current.name
  description       = "My Clumio Connection"
}

# Install the Clumio AWS template onto the registered connection
module "clumio_protect" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection.token
  role_external_id      = "my_external_id"
  aws_account_id        = clumio_aws_connection.connection.account_native_id
  aws_region            = clumio_aws_connection.connection.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_ec2_mssql_enabled = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}
```

<a name="cross-region"></a>
## Cross-Region, Two Connections
The following configuration sets up two Clumio connections, one on us-west-2 and another on
us-east-1. The [Clumio AWS module](https://registry.terraform.io/modules/clumio-code/aws-template/clumio/latest)
is installed onto both regions of the same AWS account.

```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~> 0.4.0"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate two AWS providers with different regions
provider "aws" {
  region = "us-west-2"
}

provider "aws" {
  alias  = "east"
  region = "us-east-1"
}

# Retrieve the effective AWS account ID
data "aws_caller_identity" "current" {}

# Register a new Clumio connection on us-west-2 for the effective AWS account ID
resource "clumio_aws_connection" "connection_west" {
  account_native_id = data.aws_caller_identity.current.account_id
  aws_region        = "us-west-2"
  description       = "My Clumio Connection West"
}

# Register a new Clumio connection on us-east-1 for the effective AWS account ID
resource "clumio_aws_connection" "connection_east" {
  account_native_id = data.aws_caller_identity.current.account_id
  aws_region        = "us-east-1"
  description       = "My Clumio Connection East"
}

# Install the Clumio AWS template onto the registered connection for West
module "clumio_protect_west" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_west.token
  role_external_id      = "my_external_id_west"
  aws_account_id        = clumio_aws_connection.connection_west.account_native_id
  aws_region            = clumio_aws_connection.connection_west.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_west.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_ec2_mssql_enabled = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}

# Install the Clumio AWS template onto the registered connection for East
module "clumio_protect_east" {
  providers = {
    clumio = clumio
    aws    = aws.east
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_east.token
  role_external_id      = "my_external_id_east"
  aws_account_id        = clumio_aws_connection.connection_east.account_native_id
  aws_region            = clumio_aws_connection.connection_east.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_east.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_ec2_mssql_enabled = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}
```

<a name="cross-account"></a>
## Cross-Account, Role-Assumption, Two Connections
The following configuration sets up two Clumio connections to two different AWS accounts. The
[Clumio AWS module](https://registry.terraform.io/modules/clumio-code/aws-template/clumio/latest)
is subsequently installed onto us-west-2 for one of the accounts and us-east-1 for the other.

In addition, IAM role assumption is used to provision AWS resources onto one of the accounts. See
[Assuming an IAM Role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#assuming-an-iam-role)
for additional details.

```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~> 0.4.0"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate two AWS providers with different regions. One of the providers is additionally
# associated with a different AWS account and thus requires role assumption
provider "aws" {
  region = "us-west-2"
}

provider "aws" {
  alias  = "account_2_east"
  region = "us-east-1"
  assume_role {
    role_arn    = "<assume_role_arn>"
    external_id = "<assume_role_external_id>"
  }
}

# Retrieve the effective AWS account IDs for the different AWS accounts
data "aws_caller_identity" "account_1" {}

data "aws_caller_identity" "account_2" {
  provider = aws.account_2_east
}

# Register a new Clumio connection on us-west-2 for the first AWS account ID
resource "clumio_aws_connection" "connection_account_1_west" {
  account_native_id = data.aws_caller_identity.account_1.account_id
  aws_region        = "us-west-2"
  description       = "My Clumio Connection Account 1 West"
}

# Register a new Clumio connection on us-east-1 for the second AWS account ID
resource "clumio_aws_connection" "connection_account_2_east" {
  account_native_id = data.aws_caller_identity.account_2.account_id
  aws_region        = "us-east-1"
  description       = "My Clumio Connection Account 2 East"
}

# Install the Clumio AWS template onto the registered connection for the first AWS account ID
# on West
module "clumio_protect_account_1_west" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_account_1_west.token
  role_external_id      = "my_external_id_account_1_west"
  aws_account_id        = clumio_aws_connection.connection_account_1_west.account_native_id
  aws_region            = clumio_aws_connection.connection_account_1_west.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_account_1_west.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_ec2_mssql_enabled = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}

# Install the Clumio AWS template onto the registered connection for the second AWS account ID
# on East
module "clumio_protect_account_2_east" {
  providers = {
    clumio = clumio
    aws    = aws.account_2_east
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_account_2_east.token
  role_external_id      = "my_external_id_account_2_east"
  aws_account_id        = clumio_aws_connection.connection_account_2_east.account_native_id
  aws_region            = clumio_aws_connection.connection_account_2_east.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_account_2_east.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_ec2_mssql_enabled = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}
```
