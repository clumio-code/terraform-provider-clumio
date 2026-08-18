// Copyright 2026. Clumio, Inc.

package clumio_dynamodb_backups

const (
	// Constants used by the datasource model for the clumio_dynamodb_backups Terraform datasource.
	// These values should match the schema tfsdk tags on the datasource model struct in
	// data_source_schema.go.
	schemaId                  = "id"
	schemaTableId             = "table_id"
	schemaType                = "type"
	schemaBeforeTimestamp     = "before_timestamp"
	schemaAfterTimestamp      = "after_timestamp"
	schemaStartTimestamp      = "start_timestamp"
	schemaExpirationTimestamp = "expiration_timestamp"
	schemaBackups             = "backups"

	// Possible values for the "type" attribute of a DynamoDB table backup.
	backupTypeClumioBackup = "clumio_backup"
	backupTypeAwsSnapshot  = "aws_snapshot"
)
