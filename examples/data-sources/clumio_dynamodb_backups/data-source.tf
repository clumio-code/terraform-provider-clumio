# Example retrieving all backups of a DynamoDB table. The backups are sorted newest first, so
# backups[0].id is the ID of the latest backup.
data "clumio_dynamodb_backups" "example_backups" {
  table_id = "Clumio-assigned ID of the DynamoDB table"
}

# Example retrieving the SecureVault backups of a DynamoDB table taken at or before a given time
data "clumio_dynamodb_backups" "example_backups_before" {
  table_id         = "Clumio-assigned ID of the DynamoDB table"
  type             = "clumio_backup"
  before_timestamp = "2026-07-20T00:00:00Z"
}
