resource "clumio_policy" "example_backup_windown_timezone" {
  name              = "example-policy-Backup-Window-Timezone"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_ebs_volume_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 30
      }
      rpo_frequency {
        unit  = "days"
        value = 1
      }
    }
    backup_window_tz {
      start_time = "05:00"
    }
    timezone          = "America/Los_Angeles"
  }
}
