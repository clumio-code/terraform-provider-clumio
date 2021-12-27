resource "clumio_aws_connection" "example" {
  account_native_id           = "aws_account_id"
  aws_region                  = "aws_region"
  description                 = "description"
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  services_enabled            = ["discover", "protect"]
  organizational_unit_id      = "organizational_unit_id"
}