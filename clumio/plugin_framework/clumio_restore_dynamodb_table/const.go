// Copyright 2026. Clumio, Inc.

package clumio_restore_dynamodb_table

const (
	// Constants used by the action model for the clumio_restore_dynamodb_table Terraform action.
	// These values should match the schema tfsdk tags on the action model structs in schema.go.
	schemaSource                  = "source"
	schemaSecurevaultBackup       = "securevault_backup"
	schemaBackupId                = "backup_id"
	schemaContinuousBackup        = "continuous_backup"
	schemaTableId                 = "table_id"
	schemaTimestamp               = "timestamp"
	schemaUseLatestRestorableTime = "use_latest_restorable_time"
	schemaClumioType              = "clumio_type"
	schemaTarget                  = "target"
	schemaEnvironmentId           = "environment_id"
	schemaTableName               = "table_name"

	// Possible values for the "clumio_type" attribute of a continuous backup restore source.
	clumioTypeClumioPitr = "clumio_pitr"
	clumioTypeAwsPitr    = "aws_pitr"
)
