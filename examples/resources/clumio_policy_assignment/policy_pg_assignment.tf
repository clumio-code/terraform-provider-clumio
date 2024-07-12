resource "clumio_policy_assignment" "example" {
  entity_id   = clumio_policy.test_policy.id
  entity_type = "protection_group"
  policy_id   = "policy_id"
}

resource "clumio_policy" "test_policy" {
  name = "test-policy"
  operations {
    action_setting = "immediate"
    type = "aws_dynamodb_table_backup"
    slas {
      retention_duration {
        unit = "days"
        value = 3
      }
      rpo_frequency {
        unit = "hours"
        value = 4
      }
    }
  }
  operations {
    action_setting = "immediate"
    type = "protection_group_backup"
    slas {
      retention_duration {
        unit = "months"
        value = 3
      }
      rpo_frequency {
        unit = "days"
        value = 2
      }
    }
    advanced_settings {
      protection_group_backup {
        backup_tier = "cold"
      }
    }
  }
}
