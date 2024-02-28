data "clumio_aws_manual_connection_resources" "test_get_resources" {
  account_native_id = "aws_account_id"
  aws_region        = "aws_region"
  asset_types_enabled = {
    ebs   = true
    rds   = true
    ddb   = true
    s3    = true
    mssql = true
  }
}