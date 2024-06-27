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

resource "clumio_policy_rule" "example_1" {
  name           = "example-policy-rule-1"
  # Using the clumio_policy resource to get the id of the policy.
  policy_id      = clumio_policy.example_ebs.id
  before_rule_id = clumio_policy_rule.example_2.id
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_account_native_id\":{\"$in\":[\"aws_account_id_1\", \"aws_account_id_2\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"aws_tag_key\", \"value\":\"aws_tag_value\"}}}"
}
