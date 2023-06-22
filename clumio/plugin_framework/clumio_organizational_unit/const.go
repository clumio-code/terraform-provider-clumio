// Copyright 2023. Clumio, Inc.

package clumio_organizational_unit

const (
	schemaId                        = "id"
	schemaName                      = "name"
	schemaDescription               = "description"
	schemaParentId                  = "parent_id"
	schemaChildrenCount             = "children_count"
	schemaConfiguredDatasourceTypes = "configured_datasource_types"
	schemaDescendantIds             = "descendant_ids"
	schemaUserCount                 = "user_count"
	schemaUsers                     = "users"
	schemaUsersWithRole             = "users_with_role"
	schemaUserId                    = "user_id"
	schemaAssignedRole              = "assigned_role"
	http200                         = 200
	http202                         = 202
	errorFmt                        = "Error: %v"
	pollTimeoutInSec                = 3600
	pollIntervalInSec               = 5
)
