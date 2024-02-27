// Copyright 2023. Clumio, Inc.

package clumio_policy_assignment

const (
	// Constants used by the resource model for the clumio_policy_assignment Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId                   = "id"
	schemaEntityId             = "entity_id"
	schemaEntityType           = "entity_type"
	schemaPolicyId             = "policy_id"
	schemaOrganizationalUnitId = "organizational_unit_id"

	entityTypeProtectionGroup = "protection_group"
	protectionGroupBackup     = "protection_group_backup"

	timeoutInSec  = 3600
	intervalInSec = 5
)

var (
	actionAssign   = "assign"
	actionUnassign = "unassign"
	policyIdEmpty  = ""
)
