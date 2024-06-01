resource "clumio_policy" "example_mssql-ec2" {
  name              = "example-policy-MSSQL-EC2"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "ec2_mssql_database_backup"
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
    advanced_settings {
      ec2_mssql_database_backup {
        alternative_replica = "sync_secondary"
        preferred_replica = "primary"
      }
    }
  }
  operations {
    action_setting = "immediate"
    type = "ec2_mssql_log_backup"
    slas {
      retention_duration {
        unit = "days"
        value = 5
      }
      rpo_frequency {
        unit = "minutes"
        value = 15
      }
    }
    advanced_settings {
      ec2_mssql_log_backup {
        alternative_replica = "sync_secondary"
        preferred_replica = "primary"
      }
    }
  }
}
