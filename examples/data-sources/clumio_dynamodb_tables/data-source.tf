# Example using table_native_id
data "clumio_dynamodb_tables" "ds_dynamodb_tables" {
  account_native_id = "AWS Account ID"
  aws_region = "AWS Region"
  table_native_id = "DynamoDB table native ID"
}

# Example using table_names
data "clumio_dynamodb_tables" "ds_dynamodb_tables" {
  account_native_id = "AWS Account ID"
  aws_region = "AWS Region"
  name = "DynamoDB table name"
}
