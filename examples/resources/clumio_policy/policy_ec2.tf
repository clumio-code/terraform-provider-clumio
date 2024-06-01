resource "clumio_policy" "example_ec2" {
  name              = "example-policy-EC2"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_ec2_instance_backup"
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
      aws_ec2_instance_backup {
        backup_tier = "standard"
      }
    }
  }
}
