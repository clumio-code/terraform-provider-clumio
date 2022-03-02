terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~> 0.2.2"
    }
    aws    = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate two AWS providers with different regions and roles to assume in
# order to access different AWS accounts
provider "aws" {
  region = "us-west-2"
  assume_role {
    role_arn    = "<assume_role_arn_1>"
    external_id = "<assume_role_external_id_1>"
  }
}
provider "aws" {
  alias  = "account_2_east"
  region = "us-east-1"
  assume_role {
    role_arn    = "<assume_role_arn_2>"
    external_id = "<assume_role_external_id_2>"
  }
}

# Retrieve the effective AWS account IDs for the different AWS accounts
data "aws_caller_identity" "account_1" {}
data "aws_caller_identity" "account_2" {
  provider = aws.account_2_east
}

# Register a new Clumio connection on us-west-2 for the first AWS account ID
resource "clumio_aws_connection" "connection_account_1_west" {
  account_native_id           = data.aws_caller_identity.account_1.account_id
  aws_region                  = "us-west-2"
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection Account 1 West"
}

# Register a new Clumio connection on us-east-1 for the second AWS account ID
resource "clumio_aws_connection" "connection_account_2_east" {
  account_native_id           = data.aws_caller_identity.account_2.account_id
  aws_region                  = "us-east-1"
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection Account 2 East"
}

# Install the Clumio Protect template onto the registered connection for the
# first AWS account ID on West
module "clumio_protect_account_1_west" {
  providers             = {
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
  is_ebs_enabled               = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "S3")
}

# Install the Clumio Protect template onto the registered connection for the
# second AWS account ID on East
module "clumio_protect_account_2_east" {
  providers             = {
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
  is_ebs_enabled               = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "S3")
}
