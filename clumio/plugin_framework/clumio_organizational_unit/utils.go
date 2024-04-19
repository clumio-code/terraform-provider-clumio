// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_organizational_unit Terraform resource.

package clumio_organizational_unit

import (
	"context"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// populateOrganizationalUnitsInDataSourceModel is used to populate the users schema attribute in
// the data source model from the results in the API response.
func populateOrganizationalUnitsInDataSourceModel(ctx context.Context,
	model *clumioOrganizationalUnitDataSourceModel,
	items []*models.OrganizationalUnitWithETag) diag.Diagnostics {

	var diags diag.Diagnostics

	userObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaAssignedRole: types.StringType,
			schemaUserId:       types.StringType,
		},
	}
	objtype := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaId:            types.StringType,
			schemaName:          types.StringType,
			schemaDescription:   types.StringType,
			schemaParentId:      types.StringType,
			schemaDescendantIds: types.SetType{ElemType: types.StringType},
			schemaUsersWithRole: types.SetType{ElemType: userObjType},
		},
	}
	attrVals := make([]attr.Value, 0)
	for _, item := range items {
		attrTypes := make(map[string]attr.Type)
		attrTypes[schemaId] = types.StringType
		attrTypes[schemaName] = types.StringType
		attrTypes[schemaDescription] = types.StringType
		attrTypes[schemaParentId] = types.StringType
		attrTypes[schemaDescendantIds] = types.SetType{ElemType: types.StringType}
		attrTypes[schemaUsersWithRole] = types.SetType{ElemType: userObjType}

		attrValues := make(map[string]attr.Value)
		attrValues[schemaId] = basetypes.NewStringPointerValue(item.Id)
		attrValues[schemaName] = basetypes.NewStringPointerValue(item.Name)
		attrValues[schemaDescription] = basetypes.NewStringPointerValue(item.Description)
		attrValues[schemaParentId] = basetypes.NewStringPointerValue(item.ParentId)
		descIds := make([]string, 0)
		for _, descId := range item.DescendantIds {
			descIds = append(descIds, *descId)
		}
		descendantIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, descIds)
		diags.Append(conversionDiags...)
		attrValues[schemaDescendantIds] = descendantIds

		if item.Users == nil {
			attrValues[schemaUsersWithRole] = basetypes.NewSetNull(userObjType)
		} else {
			userAttrVals := make([]attr.Value, 0)
			for _, user := range item.Users {
				userAttrTypes := make(map[string]attr.Type)
				userAttrTypes[schemaAssignedRole] = types.StringType
				userAttrTypes[schemaUserId] = types.StringType

				userAttrValues := make(map[string]attr.Value)
				userAttrValues[schemaAssignedRole] = basetypes.NewStringPointerValue(
					user.AssignedRole)
				userAttrValues[schemaUserId] = basetypes.NewStringPointerValue(user.UserId)
				userObj, conversionDiags := types.ObjectValue(userAttrTypes, userAttrValues)
				diags.Append(conversionDiags...)
				if diags.HasError() {
					return diags
				}
				userAttrVals = append(attrVals, userObj)
			}
			userSetObj, setDiag := types.SetValue(userObjType, userAttrVals)
			diags.Append(setDiag...)
			if diags.HasError() {
				return diags
			}
			attrValues[schemaUsersWithRole] = userSetObj
		}
		obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
		attrVals = append(attrVals, obj)
	}
	setObj, listdiag := types.SetValue(objtype, attrVals)
	diags.Append(listdiag...)
	model.OrganizationalUnits = setObj

	return diags
}
