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
  condition = jsonencode({
    "entity_type" : {
      "$eq" : "aws_ebs_volume"
    },
    "aws_account_native_id" : {
      "$in" : ["123456789012", "234567890123"]
    },
    "aws_tag" : {
      "$eq" : {
        "key" : "Key1",
        "value" : "Value1"
      }
    }
  })
}
