resource "clumio_policy" "example" {
  name                   = "example-policy"
  organizational_unit_id = "organizational_unit_id"
  activation_status      = "activated"
  operations {
    action_setting = "window"
    type           = "aws_ebs_volume_backup"
    backup_window {
      start_time = "08:00"
      end_time   = "20:00"
    }
    slas {
      retention_duration {
        unit  = "days"
        value = 1
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
