terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.11.0"
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
  default_tags {
    tags = {
      "Vendor" = "Clumio"
    }
  }
}

provider "aws" {
  alias  = "east"
  region = "us-east-1"
  default_tags {
    tags = {
      "Vendor" = "Clumio"
    }
  }
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
  role_external_id      = clumio_aws_connection.connection_west.role_external_id
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
  role_external_id      = clumio_aws_connection.connection_east.role_external_id
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
