// Copyright 2026. Clumio, Inc.

package clumio_gcs_bucket

const (
	// Constants used by the clumio_gcs_bucket Terraform data source schema. These values should
	// match the attribute names used in data_source_schema.go (including nested attributes).
	schemaBucketName           = "bucket_name"
	schemaBucketNames          = "bucket_names"
	schemaGCSBuckets           = "gcs_buckets"
	schemaId                   = "id"
	schemaCreatedTimestamp     = "created_timestamp"
	schemaIsDeleted            = "is_deleted"
	schemaIsVersioningEnabled  = "is_versioning_enabled"
	schemaLastBackupTimestamp  = "last_backup_timestamp"
	schemaLocation             = "location"
	schemaLocationType         = "location_type"
	schemaLocationUuid         = "location_uuid"
	schemaObjectCount          = "object_count"
	schemaOrganizationalUnitId = "organizational_unit_id"
	schemaProjectId            = "project_id"
	schemaProjectUuid          = "project_uuid"
	schemaProtectionGroupCount = "protection_group_count"
	schemaSizeBytes            = "size_bytes"
	schemaLabels               = "labels"
	schemaKey                  = "key"
	schemaValue                = "value"
)
