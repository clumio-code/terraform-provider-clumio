// Copyright 2023. Clumio, Inc.

package clumio_protection_group

const (
	schemaId                   = "id"
	schemaBucketRule           = "bucket_rule"
	schemaDescription          = "description"
	schemaName                 = "name"
	schemaObjectFilter         = "object_filter"
	schemaLatestVersionOnly    = "latest_version_only"
	schemaPrefixFilters        = "prefix_filters"
	schemaExcludedSubPrefixes  = "excluded_sub_prefixes"
	schemaPrefix               = "prefix"
	schemaStorageClasses       = "storage_classes"
	schemaPolicyId             = "policy_id"
	schemaOrganizationalUnitId = "organizational_unit_id"
	schemaProtectionInfo       = "protection_info"
	schemaInheritingEntityId   = "inheriting_entity_id"
	schemaInheritingEntityType = "inheriting_entity_type"
	schemaProtectionStatus     = "protection_status"

	timeoutInSec  = 3600
	intervalInSec = 5

	errorFmt                    = "Error: %v"
	errorProtectionGroupReadFmt = "Error reading Protection Group %v."
)
