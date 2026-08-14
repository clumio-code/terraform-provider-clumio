---
page_title: "Sample Configuration"
---

# Sample Configuration
The following sample Terraform configurations highlight several resources in the provider.

## AWS
This configuration includes:
* `clumio_aws_connection` (as well as installation of the [Clumio AWS module](https://registry.terraform.io/modules/clumio-code/aws-template/clumio/latest))
* `clumio_protection_group`
* `clumio_policy`
* `clumio_policy_assignment`
* `clumio_policy_rule`
* `clumio_organizational_unit`
* `clumio_user`

```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.20.0"
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
  default_tags {
    tags = {
      "Vendor" = "Clumio"
    }
  }
}

# Retrieve the effective AWS account ID and region
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Register a new Clumio connection for the effective AWS account ID and region
resource "clumio_aws_connection" "connection" {
  account_native_id = data.aws_caller_identity.current.account_id
  aws_region        = data.aws_region.current.region
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
  role_external_id      = clumio_aws_connection.connection.role_external_id
  aws_account_id        = clumio_aws_connection.connection.account_native_id
  aws_region            = clumio_aws_connection.connection.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled       = true
  is_rds_enabled       = true
  is_dynamodb_enabled  = true
  is_s3_enabled        = true
}

# Create a Clumio protection group that aggregates S3 buckets with the tag "clumio:example"
resource "clumio_protection_group" "protection_group" {
  name        = "My Clumio Protection Group"
  bucket_rule = jsonencode({
    "aws_tag" = {
      "$eq" = {
        "key"   = "Key1"
        "value" = "Value1"
      }
    }
  })
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
  condition = jsonencode({
    "entity_type" : {
      "$eq" : "aws_ebs_volume"
    },
    "aws_tag" : {
      "$eq" : {
        "key"   = "Key2"
        "value" = "Value2"
      }
    }
  })
  before_rule_id = clumio_policy_rule.rule_2.id
}

# Create a second Clumio policy rule, prioritized after "rule_1", and associate it with the policy
resource "clumio_policy_rule" "rule_2" {
  name           = "Second Rule"
  policy_id      = clumio_policy.policy.id
  condition = jsonencode({
    "entity_type" : {
      "$eq" : "aws_ebs_volume"
    },
    "aws_tag" : {
      "$eq" : {
        "key"   = "Key3"
        "value" = "Value3"
      }
    }
  })
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
  full_name = "Foo Bar"
  email     = "foobar@clumio.com"
  access_control_configuration = [
    {
      role_id                 = data.clumio_role.ou_admin.id,
      organizational_unit_ids = [clumio_organizational_unit.ou.id]
    }
  ]
}
```

## GCP
This configuration includes:
* `clumio_gcp_connection` (as well as installation of the [Clumio GCP module](https://registry.terraform.io/modules/clumio-code/gcp-template/clumio/latest))
* `clumio_gcs_protection_group`
* `clumio_policy`
* `clumio_policy_assignment`

```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = ">=0.22.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~>5.0"
    }
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate the Google provider for the GCP project to be protected
provider "google" {
  project = "<gcp_project_id>"
}

# Validate the project ID and retrieve project metadata
data "google_project" "current" {
  project_id = "<gcp_project_id>"
}

# Per-region configuration for the connection. Clumio creates the inventory bridge bucket for
# each region unless an existing bucket is specified via using_custom_inventory_bridge_bucket.
locals {
  region_configuration = [
    {
      region = "us-central1"
    },
    {
      region                               = "us-west1"
      using_custom_inventory_bridge_bucket = "<existing_inventory_bridge_bucket_name>"
    }
  ]
}

# Register a new Clumio connection for the GCP project. Clumio currently supports backup of GCP
# resources in us-central1 and us-west1.
resource "clumio_gcp_connection" "connection" {
  project_id  = data.google_project.current.project_id
  description = "My Clumio GCP Connection"
  regions     = [for r in local.region_configuration : r.region]
}

# Install the Clumio GCP template onto the registered connection
module "clumio_protect_gcp" {
  providers = {
    clumio = clumio
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token                 = clumio_gcp_connection.connection.token
  project_id                   = data.google_project.current.project_id
  clumio_service_account_email = clumio_gcp_connection.connection.clumio_service_account
  region_configuration         = local.region_configuration

  # Enable protection of GCS buckets
  is_gcs_enabled = true
}

# Create a Clumio GCS protection group that aggregates buckets with the label "clumio:example"
resource "clumio_gcs_protection_group" "protection_group" {
  name = "My Clumio GCS Protection Group"
  bucket_rule {
    gcp_label {
      eq = { clumio = "example" }
    }
  }
}

# Create a Clumio policy for GCS protection groups with a 7-day RPO and 3-month retention
resource "clumio_policy" "gcs_policy" {
  name = "GCS Gold"
  operations {
    action_setting = "immediate"
    type           = "gcp_protection_group_backup"
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
      gcp_protection_group_backup {
        backup_tier = "standard"
      }
    }
  }
}

# Assign the policy to the GCS protection group
resource "clumio_policy_assignment" "gcs_assignment" {
  entity_id   = clumio_gcs_protection_group.protection_group.id
  entity_type = "gcp_protection_group"
  policy_id   = clumio_policy.gcs_policy.id
}
```
