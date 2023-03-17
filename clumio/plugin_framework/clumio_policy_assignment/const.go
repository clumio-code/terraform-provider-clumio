// Copyright 2023. Clumio, Inc.

package clumio_policy_assignment

const (
	schemaId                   = "id"
	schemaEntityId             = "entity_id"
	schemaEntityType           = "entity_type"
	schemaPolicyId             = "policy_id"
	schemaOrganizationalUnitId = "organizational_unit_id"

	entityTypeProtectionGroup = "protection_group"
	protectionGroupBackup     = "protection_group_backup"

	timeoutInSec  = 3600
	intervalInSec = 5

	errorFmt = "Error: %v"
)

var (
	actionAssign   = "assign"
	actionUnassign = "unassign"
	policyIdEmpty  = ""
)
