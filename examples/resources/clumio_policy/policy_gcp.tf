resource "clumio_policy" "example_gcp_protection_group" {
  name              = "example-policy-GCP-Protection-Group"
  activation_status = "activated"
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
        value = 1
      }
    }
    advanced_settings {
      gcp_protection_group_backup {
        backup_tier = "standard"
      }
    }
  }
}
