resource "clumio_policy" "example_iceberg" {
  name              = "example-policy-Iceberg"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_iceberg_table_backup"
    slas {
      retention_duration {
        unit  = "months"
        value = 12
      }
      rpo_frequency {
        unit  = "days"
        value = 7
      }
    }
    advanced_settings {
      aws_iceberg_table_backup {
        backup_tier                    = "standard"
        backup_last_snapshot_only      = true
        backup_compacted_snapshot_only = false
      }
    }
  }
}
