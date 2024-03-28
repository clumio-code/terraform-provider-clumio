// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the user SDK APIs to perform CRUD operations and set the
// attributes from the response of the API in the resource model.

package clumio_user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createUser invokes the API to create the user and from the response populates the computed
// attributes of the user.
func (r *clumioUserResource) createUser(
	ctx context.Context, plan *clumioUserResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	if (!plan.AssignedRole.IsUnknown() || !plan.OrganizationalUnitIds.IsUnknown()) &&
		!plan.AccessControlConfiguration.IsUnknown() {
		summary := fmt.Sprintf(createErrorFmt, r.name)
		detail := "Both access_control_configuration and assigned_role/organizational_unit_ids" +
			" cannot be configured. Please configure access_control_configuration only."
		diags.AddError(summary, detail)
		return diags
	}

	// Convert access_control_configuration field from schema to API request format.
	accessControlConfiguration := make([]*models.RoleForOrganizationalUnits, 0)
	if !plan.AccessControlConfiguration.IsUnknown() {
		for _, element := range plan.AccessControlConfiguration.Elements() {
			roleForOU := roleForOrganizationalUnitModel{}
			element.(types.Object).As(ctx, &roleForOU, basetypes.ObjectAsOptions{})
			ouIds := make([]*string, 0)
			if !roleForOU.OrganizationalUnitIds.IsNull() {
				conversionDiags := roleForOU.OrganizationalUnitIds.ElementsAs(ctx, &ouIds, false)
				diags.Append(conversionDiags...)
			}
			accessControlConfiguration = append(accessControlConfiguration,
				&models.RoleForOrganizationalUnits{
					RoleId:                roleForOU.RoleId.ValueStringPointer(),
					OrganizationalUnitIds: ouIds,
				})
		}
	} else if !plan.OrganizationalUnitIds.IsUnknown() && !plan.AssignedRole.IsUnknown() {
		ouIds := make([]*string, 0)
		conversionDiags := plan.OrganizationalUnitIds.ElementsAs(ctx, &ouIds, false)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}

		accessControlConfiguration = append(accessControlConfiguration,
			&models.RoleForOrganizationalUnits{
				RoleId:                plan.AssignedRole.ValueStringPointer(),
				OrganizationalUnitIds: ouIds,
			})
	} else {
		summary := fmt.Sprintf(createErrorFmt, r.name)
		detail := "One of access_control_configuration or assigned_role/organizational_unit_ids" +
			" must be configured."
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the schema to a Clumio API request to create a Clumio user.
	apiReq := &models.CreateUserV2Request{
		Email:                      plan.Email.ValueStringPointer(),
		FullName:                   plan.FullName.ValueStringPointer(),
		AccessControlConfiguration: accessControlConfiguration,
	}

	// Call the Clumio API to create the user
	res, apiErr := r.sdkUsers.CreateUser(apiReq)
	if apiErr != nil {
		summary := fmt.Sprintf(createErrorFmt, r.name)
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
	plan.Id = types.StringPointerValue(res.Id)
	plan.Inviter = types.StringPointerValue(res.Inviter)
	plan.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
	plan.IsEnabled = types.BoolPointerValue(res.IsEnabled)
	plan.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
	plan.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)
	accessControlCfg := getAccessControlCfgFromHTTPRes(
		ctx, res.AccessControlConfiguration, &diags)
	plan.AccessControlConfiguration = accessControlCfg
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, res.AccessControlConfiguration[0].OrganizationalUnitIds)
	diags.Append(conversionDiags...)
	plan.OrganizationalUnitIds = ouIds
	plan.AssignedRole = basetypes.NewStringPointerValue(
		res.AccessControlConfiguration[0].RoleId)

	return diags
}

// readUser invokes the API to read the user and from the response populates the attributes of the
// user. If the user has been removed externally, the function returns "true" to indicate to the
// caller that the resource no longer exists.
func (r *clumioUserResource) readUser(ctx context.Context, state *clumioUserResourceModel) (
	bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	userId, perr := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if perr != nil {
		summary := invalidUserMsg
		detail := fmt.Sprintf(invalidUserFmt, state.Id.ValueString())
		diags.AddError(summary, detail)
	}

	// Call the Clumio API to read the user.
	res, apiErr := r.sdkUsers.ReadUser(userId)
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf(
				"%s (ID: %v) not found. Removing from state", r.name, state.Id.ValueString())
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
	state.Inviter = types.StringPointerValue(res.Inviter)
	state.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
	state.IsEnabled = types.BoolPointerValue(res.IsEnabled)
	state.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
	state.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)
	state.Email = types.StringPointerValue(res.Email)
	state.FullName = types.StringPointerValue(res.FullName)

	accessControlCfg := getAccessControlCfgFromHTTPRes(
		ctx, res.AccessControlConfiguration, &diags)
	state.AccessControlConfiguration = accessControlCfg
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, res.AccessControlConfiguration[0].OrganizationalUnitIds)
	diags.Append(conversionDiags...)
	state.OrganizationalUnitIds = ouIds
	state.AssignedRole = basetypes.NewStringPointerValue(res.AccessControlConfiguration[0].RoleId)

	return false, diags
}

// updateUser invokes the API to update the user and from the response populates the computed
// attributes of the user.
func (r *clumioUserResource) updateUser(ctx context.Context, plan *clumioUserResourceModel,
	state *clumioUserResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	if (!plan.AssignedRole.IsUnknown() || !plan.OrganizationalUnitIds.IsUnknown()) &&
		!plan.AccessControlConfiguration.IsUnknown() {
		summary := fmt.Sprintf("Unable to update %s", r.name)
		detail := "Both access_control_configuration and assigned_role/organizational_unit_ids" +
			" cannot be configured. Please configure access_control_configuration only."
		diags.AddError(summary, detail)
		return diags
	}

	updateRequest := &models.UpdateUserV2Request{}
	if !plan.FullName.IsUnknown() &&
		state.FullName != plan.FullName {
		updateRequest.FullName = plan.FullName.ValueStringPointer()
	}
	if !plan.AccessControlConfiguration.IsUnknown() {
		add, remove := getAccessControlCfgUpdates(ctx, state.AccessControlConfiguration.Elements(),
			plan.AccessControlConfiguration.Elements())

		updateRequest.AccessControlConfigurationUpdates = &models.EntityGroupAssignmentUpdates{
			Add:    add,
			Remove: remove,
		}
	} else if !plan.OrganizationalUnitIds.IsUnknown() && !plan.AssignedRole.IsUnknown() {

		// If AssignedRole and OrganizationalUnitIds are set in the plan, then first create
		// an accessControlConfiguration and use it to call getAccessControlCfgUpdates() to
		// calculate the EntityGroupAssignmentUpdates Add and Remove values. This is being done so
		// that we could reuse the getAccessControlCfgUpdates() even in the case
		// AccessControlConfiguration is not set in the plan.
		accessControlConfiguration := make([]*roleForOrganizationalUnitModel, 0)
		accessControlConfiguration = append(accessControlConfiguration,
			&roleForOrganizationalUnitModel{
				RoleId:                plan.AssignedRole,
				OrganizationalUnitIds: plan.OrganizationalUnitIds,
			})
		accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				schemaRoleId: types.StringType,
				schemaOrganizationalUnitIds: types.SetType{
					ElemType: types.StringType,
				},
			},
		}, accessControlConfiguration)
		diags.Append(conversionDiags...)
		add, remove := getAccessControlCfgUpdates(ctx, state.AccessControlConfiguration.Elements(),
			accessControlList.Elements())

		updateRequest.AccessControlConfigurationUpdates = &models.EntityGroupAssignmentUpdates{
			Add:    add,
			Remove: remove,
		}
	}
	userId, perr := strconv.ParseInt(plan.Id.ValueString(), 10, 64)
	if perr != nil {
		summary := invalidUserMsg
		detail := fmt.Sprintf(invalidUserFmt, plan.Id.ValueString())
		diags.AddError(summary, detail)
	}

	// Call the Clumio API to update the user.
	res, apiErr := r.sdkUsers.UpdateUser(userId, updateRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, state.Id.ValueString())
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
	plan.Inviter = types.StringPointerValue(res.Inviter)
	plan.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
	plan.IsEnabled = types.BoolPointerValue(res.IsEnabled)
	plan.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
	plan.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)

	accessControlCfg := getAccessControlCfgFromHTTPRes(
		ctx, res.AccessControlConfiguration, &diags)
	plan.AccessControlConfiguration = accessControlCfg
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, res.AccessControlConfiguration[0].OrganizationalUnitIds)
	diags.Append(conversionDiags...)
	plan.OrganizationalUnitIds = ouIds
	plan.AssignedRole = basetypes.NewStringPointerValue(
		res.AccessControlConfiguration[0].RoleId)

	return diags
}

// deleteUser invokes the API to delete the user.
func (r *clumioUserResource) deleteUser(
	_ context.Context, state *clumioUserResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	userId, perr := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if perr != nil {
		summary := invalidUserMsg
		detail := fmt.Sprintf(invalidUserFmt, state.Id.ValueString())
		diags.AddError(summary, detail)
	}

	// Call the Clumio API to delete the user.
	_, apiErr := r.sdkUsers.DeleteUser(userId)
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.Id.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}

	return diags
}
