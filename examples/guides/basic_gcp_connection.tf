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

# Register a new Clumio connection for the GCP project. Clumio currently supports backup of GCP
# resources in us-central1 and us-west1.
resource "clumio_gcp_connection" "connection" {
  project_id  = data.google_project.current.project_id
  description = "My Clumio GCP Connection"
  regions     = ["us-central1", "us-west1"]
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

  # Per-region configuration. Clumio creates the inventory bridge bucket for each region unless
  # an existing bucket is specified via using_custom_inventory_bridge_bucket.
  region_configuration = [
    for region in clumio_gcp_connection.connection.regions : { region = region }
  ]

  # Enable protection of GCS buckets
  is_gcs_enabled = true
}
