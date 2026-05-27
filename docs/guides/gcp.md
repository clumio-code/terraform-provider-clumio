---
page_title: "Using Connections and the GCP Module"
---

# Using Connections and the GCP Module

> ⚠️ **Beta Resource**
>
> The Clumio GCP integration is currently in beta and available only to select customers.
> Behavior, schema, and APIs may change in future releases.

- [Preparation](#preparation)
- [Basic, One Connection](#basic)
- [Multi-Project, Two Connections](#multi-project)
- [Custom WIF Names](#custom-wif)
- [Troubleshooting](#troubleshooting)

The following are examples of various ways to instantiate Clumio connections and install the
[Clumio GCP module](https://registry.terraform.io/modules/clumio-code/gcp-template/clumio/latest) to
one or more GCP projects to be protected.

<a name="preparation"></a>
## Preparation
Please see the "Getting Started" guide for notes about setting up a Clumio API key. In addition,
ensure that the [Google provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
is configured with credentials whose principal can create Workload Identity Pools, service accounts,
and IAM bindings in the target GCP project. For details on configuring the Google provider, see:
https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference

<a name="basic"></a>
## Basic, One Connection
The following configuration sets up a single Clumio connection and installs the
[Clumio GCP module](https://registry.terraform.io/modules/clumio-code/gcp-template/clumio/latest) to
the GCP project to be protected.

```terraform
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
  is_gcs_enabled            = true
}
```

<a name="multi-project"></a>
## Multi-Project, Two Connections
The following configuration sets up two Clumio connections to two different GCP projects, using
Google provider aliases to scope each project's resources independently. The
[Clumio GCP module](https://registry.terraform.io/modules/clumio-code/gcp-template/clumio/latest) is
installed onto each project independently.

```terraform
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
```

<a name="custom-wif"></a>
## Custom WIF Names
By default, the Clumio GCP module creates Workload Identity Pool and Provider resources with the IDs
`clumio-aws-pool` and `clumio-aws-provider`. GCP **soft-deletes** Workload Identity Pools for 30 days
after deletion, during which time the pool's ID remains reserved. If you re-onboard a project within
that window, or the default IDs collide with existing resources in the project, override the pool and
provider IDs via the module's optional inputs:

```shell
module "clumio_protect_gcp" {
  providers = {
    clumio = clumio
  }
  source = "clumio-code/gcp-template/clumio"

  clumio_token              = clumio_gcp_connection.connection.token
  project_id                = data.google_project.current.project_id
  clumio_control_plane_id   = clumio_gcp_connection.connection.clumio_control_plane_id
  clumio_control_plane_role = clumio_gcp_connection.connection.clumio_control_plane_role

  # Override default WIF identifiers to avoid collisions with existing or soft-deleted pools
  clumio_wif_pool_id     = "clumio-aws-pool-v2"
  clumio_wif_provider_id = "clumio-aws-provider-v2"

  is_gcs_enabled = true
}
```

<a name="troubleshooting"></a>
## Troubleshooting

### Workload Identity Pool or Provider already exists

**Symptom:** `terraform apply` returns an error similar to:

```
Error: Error creating WorkloadIdentityPool: googleapi: Error 409:
Requested entity already exists, alreadyExists
```

**Cause:** Workload Identity Pool and Provider IDs must be unique within a GCP project. GCP
soft-deletes Workload Identity Pools for 30 days after deletion, during which time the pool's ID is
still reserved. A pool in a soft-deleted state will cause this conflict.

**Resolution:** Override the module's default WIF identifiers using the `clumio_wif_pool_id` and
`clumio_wif_provider_id` inputs (see [Custom WIF Names](#custom-wif)).

### Insufficient GCP permissions

**Symptom:**

```
Error: googleapi: Error 403: Permission denied on resource ...
```

**Resolution:** Ensure the GCP identity running Terraform has permission to create Workload Identity
Pools, service accounts, custom IAM roles, project-level IAM bindings, Pub/Sub topics, and Cloud
Asset Inventory feeds in the target project.
