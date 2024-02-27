// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_organizational_unit Terraform resource.

package clumio_organizational_unit

import (
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// getUsersFromHTTPRes parses "users" field from HTTP response and returns list of user ids and roles
func getUsersFromHTTPRes(users []*models.UserWithRole) ([]*string, []userWithRole) {
	userSlice, userWithRoleSlice := make([]*string, len(users)), make([]userWithRole, len(users))
	for idx, user := range users {
		userSlice[idx] = user.UserId
		userWithRoleSlice[idx] = userWithRole{
			UserId:       types.StringPointerValue(user.UserId),
			AssignedRole: types.StringPointerValue(user.AssignedRole),
		}
	}
	return userSlice, userWithRoleSlice
}
