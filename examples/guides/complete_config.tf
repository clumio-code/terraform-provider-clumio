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

# Create a Clumio protection group that aggregates S3 buckets with the tag "clumio:example"
resource "clumio_protection_group" "protection_group" {
  name        = "My Clumio Protection Group"
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"clumio\", \"value\":\"example\"}}}"
  object_filter {
    storage_classes = ["S3 Standard", "S3 Standard-IA"]
  }
}

# Create a Clumio policy with support for S3 and EBS
resource "clumio_policy" "policy" {
  name = "Gold"
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
  operations {
    action_setting = "immediate"
    type           = "aws_ebs_volume_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 30
      }
      rpo_frequency {
        unit  = "days"
        value = 1
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

# Create a Clumio policy rule and associate it with the policy
resource "clumio_policy_rule" "rule_1" {
  name           = "First Rule"
  policy_id      = clumio_policy.policy.id
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_tag\":{\"$eq\":{\"key\":\"clumio\", \"value\":\"demo\"}}}"
  before_rule_id = clumio_policy_rule.rule_2.id
}

# Create a second Clumio policy rule, prioritized after "rule_1", and associate it with the policy
resource "clumio_policy_rule" "rule_2" {
  name           = "Second Rule"
  policy_id      = clumio_policy.policy.id
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_tag\":{\"$eq\":{\"key\":\"clumio\", \"value\":\"example\"}}}"
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
