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
