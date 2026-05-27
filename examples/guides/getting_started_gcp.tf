terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.20.0"
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

# Register a new Clumio connection for the GCP project
resource "clumio_gcp_connection" "connection" {
  project_id  = data.google_project.current.project_id
  description = "My Clumio GCP Connection"
}

# Install the Clumio GCP template onto the registered connection
module "clumio_protect_gcp" {
  providers = {
    clumio = clumio
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token              = clumio_gcp_connection.connection.token
  project_id                = data.google_project.current.project_id
  clumio_control_plane_id   = clumio_gcp_connection.connection.clumio_control_plane_id
  clumio_control_plane_role = clumio_gcp_connection.connection.clumio_control_plane_role

  # Enable protection of GCS buckets
  is_gcs_enabled = true
}
