resource "clumio_policy" "example" {
  name                   = "example-policy"
  organizational_unit_id = "organizational_unit_id"
  activation_status      = "activated"
  operations {
    action_setting = "window"
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
