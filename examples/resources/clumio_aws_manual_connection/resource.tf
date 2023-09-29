resource "clumio_aws_manual_connection" "test_update_resources" {

  account_id = "aws_account_id"
  aws_region = "aws_region"
  assets_enabled = {
    ebs = true
    rds = true
    ddb = true
    s3 = true
    mssql = true
  }
  resources = {
    clumio_iam_role_arn = "clumio_iam_role_arn"
    clumio_event_pub_arn = "clumio_event_pub_arn"
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
        ssm_notification_role_arn = "ssm_notification_role_arn"
        ec2_ssm_instance_profile_arn = "ec2_ssm_instance_profile_arn"
      }
    }
  }
}
