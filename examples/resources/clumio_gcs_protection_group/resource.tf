resource "clumio_gcs_protection_group" "example" {
  name = "my-gcp-protection-group"

  bucket_rule {
    gcp_project_id {
      eq = "my-project"
    }
    gcp_label {
      # Single-label and match-all operators are maps of key => value.
      eq  = { team = "data" }
      all = { managed = "true", tier = "gold" }
      # in/not_in are repeated blocks so the same key can carry multiple values.
      in {
        key   = "env"
        value = "prod"
      }
      in {
        key   = "env"
        value = "staging"
      }
    }
  }

  include_prefixes    = ["data/"]
  latest_version_only = true
}
