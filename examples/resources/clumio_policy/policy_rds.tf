resource "clumio_policy" "example_rds" {
  name              = "example-policy-RDS"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_rds_resource_granular_backup"
    slas {
      retention_duration {
        unit  = "months"
        value = 12
      }
      rpo_frequency {
        unit  = "months"
        value = 1
      }
    }
    advanced_settings {
      aws_rds_resource_granular_backup {
        backup_tier = "frozen"
      }
    }
  }
}
