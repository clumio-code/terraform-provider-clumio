resource "clumio_policy" "example_s3_backtrack" {
  name              = "example-policy-S3-Backtrack"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_s3_backtrack"
    slas {
      retention_duration {
        unit  = "months"
        value = 3
      }
      rpo_frequency {
        unit  = "days"
        value = 1
      }
    }
  }
}
