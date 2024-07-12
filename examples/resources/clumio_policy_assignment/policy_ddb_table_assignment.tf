# Example using table_native_id
data "clumio_dynamodb_tables" "ds_dynamodb_tables" {
  account_native_id = "AWS Account ID"
  aws_region = "AWS Region"
  table_native_id = "DynamoDB table native ID"
}

resource "clumio_policy_assignment" "example" {
  entity_id   = clumio_dynamodb_tables.ds_dynamodb_tables.dynamodb_tables[0].id
  entity_type = "aws_dynamodb_table"
  policy_id   = "policy_id"
}
