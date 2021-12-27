resource "clumio_post_process_aws_connection" "example" {
  # example configuration here
  token                              = "mytoken"
  role_external_id                   = "role_external_id"
  account_id                         = "account_id"
  region                             = "region"
  role_arn                           = "role_arn"
  config_version                     = "1"
  discover_version                   = "3"
  protect_config_version             = "18"
  protect_ebs_version                = "19"
  protect_rds_version                = "18"
  protect_ec2_mssql_version          = "1"
  protect_warm_tier_version          = "2"
  protect_warm_tier_dynamodb_version = "2"
  protect_s3_version                 = "1"
}