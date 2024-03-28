// Copyright 2023. Clumio, Inc.

package clumio_user

const (
	// Constants used by the resource model for the clumio_user Terraform resource. These values
	// should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId                         = "id"
	schemaEmail                      = "email"
	schemaFullName                   = "full_name"
	schemaAssignedRole               = "assigned_role"
	schemaOrganizationalUnitIds      = "organizational_unit_ids"
	schemaAccessControlConfiguration = "access_control_configuration"
	schemaRoleId                     = "role_id"
	schemaInviter                    = "inviter"
	schemaIsConfirmed                = "is_confirmed"
	schemaIsEnabled                  = "is_enabled"
	schemaLastActivityTimestamp      = "last_activity_timestamp"
	schemaOrganizationalUnitCount    = "organizational_unit_count"

	// Common error messages used by the resource.
	invalidUserMsg = "Invalid user id."
	invalidUserFmt = "Invalid user id: %v"
	createErrorFmt = "Unable to create %s"
)
