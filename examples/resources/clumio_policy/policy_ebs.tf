resource "clumio_policy" "example_ebs" {
  name              = "example-policy-EBS"
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
    advanced_settings {
      aws_ebs_volume_backup {
        backup_tier = "standard"
      }
    }
  }
}
