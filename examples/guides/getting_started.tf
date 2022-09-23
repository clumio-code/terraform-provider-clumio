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

  # Enable protection of all data sources.
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_ec2_mssql_enabled = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}

# Create a Clumio protection group that aggregates S3 buckets with the tag "clumio:example"
resource "clumio_protection_group" "protection_group" {
  name        = "My Clumio Protection Group"
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"clumio\", \"value\":\"example\"}}}"
  object_filter {
    storage_classes = ["S3 Standard", "S3 Standard-IA"]
  }
}

# Create a Clumio policy for protection groups with a 7-day RPO and 3-month retention
resource "clumio_policy" "policy" {
  name = "S3 Gold"
  operations {
    action_setting = "immediate"
    type           = "protection_group_backup"
    slas {
      retention_duration {
        unit  = "months"
        value = 3
      }
      rpo_frequency {
        unit  = "days"
        value = 7
      }
    }
    advanced_settings {
      protection_group_backup {
        backup_tier = "cold"
      }
    }
  }
}

# Assign the policy to the protection group
resource "clumio_policy_assignment" "assignment" {
  entity_id   = clumio_protection_group.protection_group.id
  entity_type = "protection_group"
  policy_id   = clumio_policy.policy.id
}
