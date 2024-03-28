// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the organizational unit SDK APIs to perform CRUD operations
// and set the attributes from the response of the API in the resource model.

package clumio_organizational_unit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createOrganizationalUnit invokes the API to create the organizational unit and from the response
// populates the computed attributes of the organizational unit.
func (r *clumioOrganizationalUnitResource) createOrganizationalUnit(
	ctx context.Context, plan *clumioOrganizationalUnitResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Convert the schema to a Clumio API request to create an organizational unit.
	request := &models.CreateOrganizationalUnitV2Request{
		Name:        plan.Name.ValueStringPointer(),
		ParentId:    plan.ParentId.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
	}

	// Call the Clumio API to create the organizational unit.
	res, apiErr := r.sdkOrgUnits.CreateOrganizationalUnit(nil, request)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s", r.name)
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

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	var id types.String
	var parentIdString types.String
	var childrenCount types.Int64
	var userCount types.Int64
	var configuredDatasourceTypes []*string
	var userSlice []*string
	var userWithRoleSlice []userWithRole
	var descendantIdSlice []*string

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	// API call can result in either 200 or 201 status code. Relevant data is returned inside a field
	// name mapped to the status code. Data is extracted from the correct field according to the
	// status code. Else return an empty response error.
	if res.StatusCode == http.StatusOK && res.Http200 != nil {
		id = types.StringPointerValue(res.Http200.Id)
		parentIdString = types.StringPointerValue(res.Http200.ParentId)
		childrenCount = types.Int64PointerValue(res.Http200.ChildrenCount)
		userCount = types.Int64PointerValue(res.Http200.UserCount)
		configuredDatasourceTypes = res.Http200.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http200.Users)
		descendantIdSlice = res.Http200.DescendantIds
	} else if res.StatusCode == http.StatusAccepted && res.Http202 != nil {
		id = types.StringPointerValue(res.Http202.Id)
		parentIdString = types.StringPointerValue(res.Http202.ParentId)
		childrenCount = types.Int64PointerValue(res.Http202.ChildrenCount)
		userCount = types.Int64PointerValue(res.Http202.UserCount)
		configuredDatasourceTypes = res.Http202.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http202.Users)
		descendantIdSlice = res.Http202.DescendantIds
	} else {
		summary := "Empty response returned."
		detail := "CreateOrganizationalUnit returned empty response returned which is not expected."
		diags.AddError(summary, detail)
		return diags
	}
	plan.Id = id
	plan.ParentId = parentIdString
	plan.ChildrenCount = childrenCount
	plan.UserCount = userCount

	configuredDataTypes, conversionDiags := types.ListValueFrom(
		ctx, types.StringType, configuredDatasourceTypes)
	diags.Append(conversionDiags...)
	plan.ConfiguredDatasourceTypes = configuredDataTypes

	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	diags.Append(conversionDiags...)
	plan.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	diags.Append(conversionDiags...)
	plan.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, descendantIdSlice)
	diags.Append(conversionDiags...)
	plan.DescendantIds = descendantIds
	return diags
}

// readOrganizationalUnit invokes the API to read the organizational unit and from the response
// populates the attributes of the organizational unit. If the organizational unit has been removed
// externally, the function returns "true" to indicate to the caller that the resource no longer exists.
func (r *clumioOrganizationalUnitResource) readOrganizationalUnit(
	ctx context.Context, state *clumioOrganizationalUnitResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	// Call the Clumio API to read the orgainizational unit.
	res, apiErr := r.sdkOrgUnits.ReadOrganizationalUnit(state.Id.ValueString(), nil)
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state", r.name, state.Id.ValueString())
			tflog.Warn(ctx, summary)
			remove = true
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.Id.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return remove, diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return false, diags
	}

	// Convert the Clumio API response back to a schema and update the state. In addition to
	// computed fields, all fields are populated from the API response in case any values have been
	// changed externally. ID is not updated however given that it is the field used to query the
	// resource from the backend.
	state.Name = types.StringPointerValue(res.Name)

	// Since the Description field is optional, it should only be populated if it initially
	// contained a non-null value or if there is a specific value that needs to be assigned.
	description := types.StringPointerValue(res.Description)
	if !state.Description.IsNull() || description.ValueString() != "" {
		state.Description = description
	}
	state.ParentId = types.StringPointerValue(res.ParentId)
	state.ChildrenCount = types.Int64PointerValue(res.ChildrenCount)
	state.UserCount = types.Int64PointerValue(res.UserCount)

	configuredDataTypes, conversionDiags := types.ListValueFrom(
		ctx, types.StringType, res.ConfiguredDatasourceTypes)
	diags.Append(conversionDiags...)
	state.ConfiguredDatasourceTypes = configuredDataTypes

	userSlice, userWithRoleSlice := getUsersFromHTTPRes(res.Users)
	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	diags.Append(conversionDiags...)
	state.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	diags.Append(conversionDiags...)
	state.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.DescendantIds)
	diags.Append(conversionDiags...)
	state.DescendantIds = descendantIds

	return false, diags
}

// updateOrganizationalUnit invokes the API to update the organizational unit and from the response
// populates the computed attributes of the organizational unit.
func (r *clumioOrganizationalUnitResource) updateOrganizationalUnit(
	ctx context.Context, plan *clumioOrganizationalUnitResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Convert the schema to a Clumio API request to update the organizational unit.
	createReq := &models.PatchOrganizationalUnitV2Request{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
	}

	// Call the Clumio API to update the organizational unit.
	res, apiErr := r.sdkOrgUnits.PatchOrganizationalUnit(plan.Id.ValueString(), nil, createReq)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.Id.ValueString())
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

	// Convert the Clumio API response back to a schema and populate all computed fields of the
	// plan. ID however is not updated given that it is the field used to denote which resource to
	// update in the backend.
	var parentIdString types.String
	var childrenCount types.Int64
	var userCount types.Int64
	var configuredDatasourceTypes []*string
	var userWithRoleSlice []userWithRole
	var userSlice []*string
	var descendantIdSlice []*string

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	// API call can result in either 200 or 201 status code. Relevant data is returned inside a
	// field name mapped to the status code. Data is extracted from the correct field according to
	// the status code. Else return an empty response error.
	if res.StatusCode == http.StatusOK && res.Http200 != nil {
		parentIdString = types.StringPointerValue(res.Http200.ParentId)
		childrenCount = types.Int64PointerValue(res.Http200.ChildrenCount)
		userCount = types.Int64PointerValue(res.Http200.UserCount)
		configuredDatasourceTypes = res.Http200.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http200.Users)
		descendantIdSlice = res.Http200.DescendantIds
	} else if res.StatusCode == http.StatusAccepted && res.Http202 != nil {
		parentIdString = types.StringPointerValue(res.Http202.ParentId)
		childrenCount = types.Int64PointerValue(res.Http202.ChildrenCount)
		userCount = types.Int64PointerValue(res.Http202.UserCount)
		configuredDatasourceTypes = res.Http202.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http202.Users)
		descendantIdSlice = res.Http202.DescendantIds
	} else {
		summary := "Empty response returned."
		detail := "PatchOrganizationalUnit returned empty response returned which is not expected."
		diags.AddError(summary, detail)
		return diags
	}
	plan.ParentId = parentIdString
	plan.ChildrenCount = childrenCount
	plan.UserCount = userCount

	configuredDataTypes, conversionDiags := types.ListValueFrom(
		ctx, types.StringType, configuredDatasourceTypes)
	diags.Append(conversionDiags...)
	plan.ConfiguredDatasourceTypes = configuredDataTypes

	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	diags.Append(conversionDiags...)
	plan.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	diags.Append(conversionDiags...)
	plan.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, descendantIdSlice)
	diags.Append(conversionDiags...)
	plan.DescendantIds = descendantIds
	return diags
}

// deleteOrganizationalUnit invokes the API to delete the organizational unit.
func (r *clumioOrganizationalUnitResource) deleteOrganizationalUnit(
	ctx context.Context, state *clumioOrganizationalUnitResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to delete the organizational unit.
	res, apiErr := r.sdkOrgUnits.DeleteOrganizationalUnit(state.Id.ValueString(), nil)
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.Id.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}
	err := common.PollTask(ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf(
			"Unable to poll %s (ID: %v) for deletion", r.name, state.Id.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
	}
	return diags
}
