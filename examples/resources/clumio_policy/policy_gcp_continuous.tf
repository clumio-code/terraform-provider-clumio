resource "clumio_policy" "gcp_continuous" {
  name = "GCS Continuous"
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
  operations {
    action_setting = "immediate"
    type           = "gcp_continuous_backup"
    slas {
      # Use the same retention as protection group backup.
      retention_duration {
        unit  = "months"
        value = 3
      }
      # RPO can be set to minutely or hourly intervals.
      rpo_frequency {
        unit  = "minutes"
        value = 15
      }
    }
  }
}
