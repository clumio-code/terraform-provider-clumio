// Copyright 2022. Clumio, Inc.

package clumio_policy_assignment

const (
	schemaEntityId             = "entity_id"
	schemaEntityType           = "entity_type"
	schemaPolicyId             = "policy_id"
	schemaOrganizationalUnitId = "organizational_unit_id"

	entityTypeProtectionGroup = "protection_group"
)

var (
	actionAssign       = "assign"
	actionUnassign     = "unassign"
	policyIdEmpty      = ""
	validEntityTypeMap = map[string]struct{}{"protection_group": struct{}{}}
)
