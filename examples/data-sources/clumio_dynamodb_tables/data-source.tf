# Example using table_native_id
data "clumio_dynamodb_tables" "by_native_id" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
  table_native_id   = "DynamoDB table native ID"
}

# Example using name
data "clumio_dynamodb_tables" "by_name" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
  name              = "DynamoDB table name"
}
