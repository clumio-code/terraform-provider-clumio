terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.2.2"
    }
    aws    = {}
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
  account_native_id           = data.aws_caller_identity.current.account_id
  aws_region                  = data.aws_region.current.name
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection"
}

# Install the Clumio Protect template onto the registered connection
module "clumio_protect" {
  providers             = {
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
  is_ebs_enabled               = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "S3")
}

# Create a Clumio policy with a 7-day RPO and 14-day retention
resource "clumio_policy" "policy" {
  name = "Gold"
  operations {
    action_setting = "immediate"
    type           = "aws_ebs_volume_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 14
      }
      rpo_frequency {
        unit  = "days"
        value = 7
      }
    }
  }
}

# Create a Clumio policy rule and associate it with the policy
resource "clumio_policy_rule" "rule" {
  name           = "Tag-Based Rule"
  policy_id      = clumio_policy.policy.id
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_tag\":{\"$eq\":{\"key\":\"random-test-123\", \"value\":\"random-test-123\"}}}"
  before_rule_id = ""
}

# Retrive the role for OU Admin
data "clumio_role" "ou_admin" {
  name = "Organizational Unit Admin"
}

# Create a new OU
resource "clumio_organizational_unit" "ou" {
  name = "My OU"
}

# Create a user for the OU
resource "clumio_user" "user" {
  full_name               = "Foo Bar"
  email                   = "foobar@clumio.com"
  assigned_role           = data.clumio_role.ou_admin.id
  organizational_unit_ids = [clumio_organizational_unit.ou.id]
}
