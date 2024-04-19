// Copyright 2023. Clumio, Inc.

package clumio_s3_bucket

const (
	// Constants used by the resource model for the clumio_s3_bucket Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId                            = "id"
	schemaName                          = "name"
	schemaBucketNames                   = "bucket_names"
	schemaRegion                        = "aws_region"
	schemaAccountNativeId               = "account_native_id"
	schemaProtectionGroupCount          = "protection_group_count"
	schemaEventBridgeEnabled            = "event_bridge_enabled"
	schemaLastBackupTimestamp           = "last_backup_timestamp"
	schemaLastContinuousBackupTimestamp = "last_continuous_backup_timestamp"
	schemaS3Buckets                     = "s3_buckets"
)
