// Copyright 2026. Clumio, Inc.

package clumio_gcs_protection_group

const (
	// Constants used by the resource and data source models for the clumio_gcs_protection_group
	// Terraform resource and data source. These values should match the schema tfsdk tags on the
	// model structs in schema.go and data_source_schema.go.
	schemaId                = "id"
	schemaName              = "name"
	schemaBucketRule        = "bucket_rule"
	schemaIncludePrefixes   = "include_prefixes"
	schemaExcludePrefixes   = "exclude_prefixes"
	schemaLatestVersionOnly = "latest_version_only"

	// bucket_rule condition fields and operators.
	schemaGcpLabel     = "gcp_label"
	schemaGcpLocation  = "gcp_location"
	schemaGcpProjectId = "gcp_project_id"
	schemaEq           = "eq"
	schemaContains     = "contains"
	schemaIn           = "in"
	schemaAll          = "all"
	schemaNotEq        = "not_eq"
	schemaNotContains  = "not_contains"
	schemaNotIn        = "not_in"
	schemaNotAll       = "not_all"
	schemaKey          = "key"
	schemaValue        = "value"

	// Computed attributes.
	schemaOrganizationalUnitId     = "organizational_unit_id"
	schemaProtectionStatus         = "protection_status"
	schemaProtectionInfo           = "protection_info"
	schemaPolicyId                 = "policy_id"
	schemaBucketCount              = "bucket_count"
	schemaBucketUuids              = "bucket_uuids"
	schemaCreatedTimestamp         = "created_timestamp"
	schemaModifiedTimestamp        = "modified_timestamp"
	schemaVersion                  = "version"
	schemaLabels                   = "labels"
	schemaLastBackupTimestamp      = "last_backup_timestamp"
	schemaLocation                 = "location"
	schemaLocationType             = "location_type"
	schemaManualAddedBucketCount   = "manual_added_bucket_count"
	schemaBucketRuleMatchedCount   = "bucket_rule_matched_bucket_count"
	schemaBucketRuleMatchedUuids   = "bucket_rule_matched_bucket_uuids"
	schemaTotalBackedUpObjectCount = "total_backed_up_object_count"
	schemaTotalBackedUpSizeBytes   = "total_backed_up_size_bytes"
	schemaBackupStatusStats        = "backup_status_stats"
	schemaFailureCount             = "failure_count"
	schemaNoBackupCount            = "no_backup_count"
	schemaPartialSuccessCount      = "partial_success_count"
	schemaSuccessCount             = "success_count"
)
