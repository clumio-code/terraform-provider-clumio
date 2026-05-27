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

# Register a Clumio connection for the first GCP project
resource "clumio_gcp_connection" "project_a" {
  project_id  = data.google_project.project_a.project_id
  description = "My Clumio GCP Connection (Project A)"
}

# Register a Clumio connection for the second GCP project
resource "clumio_gcp_connection" "project_b" {
  project_id  = data.google_project.project_b.project_id
  description = "My Clumio GCP Connection (Project B)"
}

# Install the Clumio GCP template onto the first project
module "clumio_protect_project_a" {
  providers = {
    clumio = clumio
    google = google.project_a
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token              = clumio_gcp_connection.project_a.token
  project_id                = data.google_project.project_a.project_id
  clumio_control_plane_id   = clumio_gcp_connection.project_a.clumio_control_plane_id
  clumio_control_plane_role = clumio_gcp_connection.project_a.clumio_control_plane_role
  is_gcs_enabled            = true
}

# Install the Clumio GCP template onto the second project
module "clumio_protect_project_b" {
  providers = {
    clumio = clumio
    google = google.project_b
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token              = clumio_gcp_connection.project_b.token
  project_id                = data.google_project.project_b.project_id
  clumio_control_plane_id   = clumio_gcp_connection.project_b.clumio_control_plane_id
  clumio_control_plane_role = clumio_gcp_connection.project_b.clumio_control_plane_role
  is_gcs_enabled            = true
}
