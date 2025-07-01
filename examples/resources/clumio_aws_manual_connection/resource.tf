resource "clumio_aws_manual_connection" "test_update_resources" {

  account_id = "123456789012" # Replace with your actual AWS account ID.
  aws_region = "us-west-2" # Replace with your actual AWS region.
  assets_enabled = {
    ebs   = true
    rds   = true
    ddb   = true
    s3    = true
    mssql = false  # Note that "mssql" is only available on legacy connections.
  }
  resources = {
    clumio_iam_role_arn     = "clumio_iam_role_arn"
    clumio_event_pub_arn    = "clumio_event_pub_arn"
    clumio_support_role_arn = "clumio_support_role_arn"
    event_rules = {
      cloudtrail_rule_arn = "cloudtrail_rule_arn"
      cloudwatch_rule_arn = "cloudwatch_rule_arn"
    }

    service_roles = {
      s3 = {
        continuous_backups_role_arn = "continuous_backups_role_arn"
      }
      mssql = {
        ssm_notification_role_arn    = "ssm_notification_role_arn"
        ec2_ssm_instance_profile_arn = "ec2_ssm_instance_profile_arn"
      }
    }
  }
}
