resource "clumio_policy" "example_dynamodb" {
  name              = "example-policy-DynamoDB"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_dynamodb_table_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 7
      }
      rpo_frequency {
        unit  = "hours"
        value = 12
      }
    }
  }
}
