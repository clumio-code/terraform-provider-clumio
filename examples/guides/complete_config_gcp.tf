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
