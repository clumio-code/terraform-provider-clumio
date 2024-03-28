// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio Role SDK API to perform read operation and
// set the attributes from the response of the API in the data source model.

package clumio_role

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// readRole invokes the API to read the roles and from the response populates the
// attributes of the role.
func (r *clumioRoleDataSource) readRole(
	_ context.Context, state *clumioRoleDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to read the list of roles.
	res, apiErr := r.roles.ListRoles()
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Find the expected role from the list of fetched roles via its name.
	var expectedRole *models.RoleWithETag
	for _, roleItem := range res.Embedded.Items {
		if *roleItem.Name == state.Name.ValueString() {
			expectedRole = roleItem
			break
		}
	}
	// Throw error if role with provided name was not found.
	if expectedRole == nil {
		summary := "Role not found"
		detail := fmt.Sprintf(
			"Couldn't find a role with the provided name %s", state.Name.ValueString())
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response for the expected role into a schema and update the state.
	state.Id = types.StringPointerValue(expectedRole.Id)
	state.Name = types.StringPointerValue(expectedRole.Name)
	state.Description = types.StringPointerValue(expectedRole.Description)
	state.UserCount = types.Int64PointerValue(expectedRole.UserCount)

	// Go through all permissions inside the expected role API response and map it to
	// permissionModel. Then update the state with it.
	permissions := make([]*permissionModel, len(expectedRole.Permissions))
	for ind, permission := range expectedRole.Permissions {
		permissionModel := &permissionModel{}
		permissionModel.Description = types.StringPointerValue(permission.Description)
		permissionModel.Id = types.StringPointerValue(permission.Id)
		permissionModel.Name = types.StringPointerValue(permission.Name)
		permissions[ind] = permissionModel
	}
	state.Permissions = permissions

	return diags
}
