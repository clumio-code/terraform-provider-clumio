data "clumio_aws_manual_connection_resources" "test_get_resources" {
  account_native_id = "123456789012" # Replace with your actual AWS account ID.
  aws_region        = "us-west-2" # Replace with your actual AWS region.
  asset_types_enabled = {
    ebs   = true
    rds   = true
    ddb   = true
    s3    = true
    mssql = false  # Note that "mssql" is only available on legacy connections.
  }
}