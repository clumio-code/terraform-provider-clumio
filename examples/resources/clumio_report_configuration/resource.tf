resource "clumio_report_configuration" "example_report_configuration" {
  name = "example-report-configuration"
  description = "Example Description."
  notification {
	  email_list = ["example-email1@clumio.com", "example-email2@clumio.com"]
  }
  parameter {
    controls {
      asset_backup {
        look_back_period {
          unit = "days"
          value = 7
        }
        minimum_retention_duration {
          unit = "days"
          value = 30
        }
        window_size {
          unit = "days"
          value = 7
        }
      }
	  }
    filters {
      asset {
        groups {
          region = "us-west-2"
          type = "aws"
        }
        tag_op_mode = "equal"
        tags {
          key = "example-key"
          value = "example-value"
        }
      }
      common {
        asset_types = ["aws_ec2_instance"]
        data_sources = ["aws"]
        organizational_units = ["00000000-0000-0000-0000-000000000000"]
      }
    }
  }
  schedule {
    day_of_week = "sunday"
    frequency = "weekly"
    start_time = "15:00"
    timezone = "America/New_York"
  }
}
