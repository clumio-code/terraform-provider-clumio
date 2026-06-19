terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = ">=0.21.0"
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

# Instantiate two Google providers, one per GCP project to be protected
provider "google" {
  alias   = "project_a"
  project = "<gcp_project_a_id>"
}

provider "google" {
  alias   = "project_b"
  project = "<gcp_project_b_id>"
}

# Validate each project ID and retrieve project metadata
data "google_project" "project_a" {
  provider   = google.project_a
  project_id = "<gcp_project_a_id>"
}

data "google_project" "project_b" {
  provider   = google.project_b
  project_id = "<gcp_project_b_id>"
}

# Register a Clumio connection for the first GCP project. Clumio currently supports backup of GCP
# resources in us-central1 and us-west1.
resource "clumio_gcp_connection" "project_a" {
  project_id  = data.google_project.project_a.project_id
  description = "My Clumio GCP Connection (Project A)"
  regions     = ["us-central1", "us-west1"]
}

# Register a Clumio connection for the second GCP project
resource "clumio_gcp_connection" "project_b" {
  project_id  = data.google_project.project_b.project_id
  description = "My Clumio GCP Connection (Project B)"
  regions     = ["us-central1", "us-west1"]
}

# Install the Clumio GCP template onto the first project
module "clumio_protect_project_a" {
  providers = {
    clumio = clumio
    google = google.project_a
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token                          = clumio_gcp_connection.project_a.token
  project_id                            = data.google_project.project_a.project_id
  regions                               = clumio_gcp_connection.project_a.regions
  clumio_service_account_email          = clumio_gcp_connection.project_a.clumio_service_account
  create_clumio_inventory_bridge_bucket = true

  # Enable protection of GCS buckets
  is_gcs_enabled = true
}

# Install the Clumio GCP template onto the second project
module "clumio_protect_project_b" {
  providers = {
    clumio = clumio
    google = google.project_b
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token                          = clumio_gcp_connection.project_b.token
  project_id                            = data.google_project.project_b.project_id
  regions                               = clumio_gcp_connection.project_b.regions
  clumio_service_account_email          = clumio_gcp_connection.project_b.clumio_service_account
  create_clumio_inventory_bridge_bucket = true

  # Enable protection of GCS buckets
  is_gcs_enabled = true
}
