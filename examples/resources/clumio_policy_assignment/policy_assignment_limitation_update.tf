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
}

resource "clumio_policy_assignment" "example" {
  entity_id   = "entity_id"
  entity_type = "protection_group"
  policy_id   = "policy_id"
}
