resource "clumio_policy" "policy" {
  name = "S3 Continuous"
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
        value = 1
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
    type           = "aws_s3_continuous_backup"
    slas {
      # Use the same retention as PG backup
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
    advanced_settings {
      protection_group_continuous_backup {
        disable_eventbridge_notification = true
      }
    }
  }
}