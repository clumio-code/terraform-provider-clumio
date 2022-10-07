resource "clumio_policy" "example_1" {
  name              = "example-policy-1"
  activation_status = "activated"
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
}

resource "clumio_policy" "example_2" {
  name              = "example-policy-2"
  activation_status = "activated"
  timezone          = "America/Los_Angeles"
  operations {
    action_setting = "window"
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
      end_time   = "07:00"
    }
  }
}
