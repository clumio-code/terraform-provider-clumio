// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_user Terraform resource. This resource
// is used to manage users within Clumio.

package clumio_user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	sdkUsers "github.com/clumio-code/clumio-go-sdk/controllers/users"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioUserResource{}
	_ resource.ResourceWithConfigure   = &clumioUserResource{}
	_ resource.ResourceWithImportState = &clumioUserResource{}
)

// clumioUserResource is the struct backing the clumio_user Terraform resource. It holds the Clumio
// API client and any other required state needed to manage users within Clumio.
type clumioUserResource struct {
	name            string
	client          *common.ApiClient
	sdkUsersV1      sdkUsers.UsersV1Client
	sdkUsersV2      sdkUsers.UsersV2Client
}

// NewClumioUserResource creates a new instance of clumioUserResource. Its attributes are 
// initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewClumioUserResource() resource.Resource {
	return &clumioUserResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioUserResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_user"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioUserResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkUsersV1 = sdkUsers.NewUsersV1(r.client.ClumioConfig)
	r.sdkUsersV2 = sdkUsers.NewUsersV2(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioUserResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (!plan.AssignedRole.IsUnknown() || !plan.OrganizationalUnitIds.IsUnknown()) &&
		!plan.AccessControlConfiguration.IsUnknown() {
		summary := "Error creating Clumio user."
		detail := fmt.Sprintf(errorFmt,
			"Both access_control_configuration and assigned_role/organizational_unit_ids"+
				" cannot be configured. Please configure access_control_configuration only.")
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	if !plan.AssignedRole.IsUnknown() || !plan.OrganizationalUnitIds.IsUnknown() {
		organizationalUnitElements := plan.OrganizationalUnitIds.Elements()
		organizationalUnitIds := make([]*string, len(organizationalUnitElements))
		for ind, element := range organizationalUnitElements {
			valString := element.String()
			organizationalUnitIds[ind] = &valString
		}

		// For backwards compatibility purposes, we're created the user via both the old V1 and new V2
		// API

		// Call the Clumio API to create the user using the old v1 API.
		res, apiErr := r.sdkUsersV1.CreateUser(&models.CreateUserV1Request{
			AssignedRole:          plan.AssignedRole.ValueStringPointer(),
			Email:                 plan.Email.ValueStringPointer(),
			FullName:              plan.FullName.ValueStringPointer(),
			OrganizationalUnitIds: organizationalUnitIds,
		})
		if apiErr != nil {
			summary := "Error creating Clumio User."
			detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
			resp.Diagnostics.AddError(summary, detail)
			return
		}
		if res == nil {
			resp.Diagnostics.AddError(
				common.NilErrorMessageSummary, common.NilErrorMessageDetail)
			return
		}

		// Convert the Clumio API response back to a schema and populate all computed fields of the plan
		// including the ID given that the resource is getting created.
		plan.Id = types.StringPointerValue(res.Id)
		plan.Inviter = types.StringPointerValue(res.Inviter)
		plan.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
		plan.IsEnabled = types.BoolPointerValue(res.IsEnabled)
		plan.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
		plan.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)
		plan.AssignedRole = types.StringPointerValue(res.AssignedRole)
		orgUnitIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.AssignedOrganizationalUnitIds)
		resp.Diagnostics.Append(conversionDiags...)
		plan.OrganizationalUnitIds = orgUnitIds

		accessControl := []roleForOrganizationalUnitModel{
			{
				RoleId:                plan.AssignedRole,
				OrganizationalUnitIds: orgUnitIds,
			},
		}
		accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				schemaRoleId: types.StringType,
				schemaOrganizationalUnitIds: types.SetType{
					ElemType: types.StringType,
				},
			},
		}, accessControl)
		resp.Diagnostics.Append(conversionDiags...)
		plan.AccessControlConfiguration = accessControlList

		// Set the schema into the Terraform state.
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		return
	}

	accessControlConfiguration := make([]*models.RoleForOrganizationalUnits, 0)
	if !plan.AccessControlConfiguration.IsNull() {
		for _, element := range plan.AccessControlConfiguration.Elements() {
			roleForOU := roleForOrganizationalUnitModel{}
			element.(types.Object).As(ctx, &roleForOU, basetypes.ObjectAsOptions{})
			ouIds := make([]*string, 0)
			if !roleForOU.OrganizationalUnitIds.IsNull() {
				conversionDiags := roleForOU.OrganizationalUnitIds.ElementsAs(ctx, &ouIds, false)
				resp.Diagnostics.Append(conversionDiags...)
			}
			accessControlConfiguration = append(accessControlConfiguration, &models.RoleForOrganizationalUnits{
				RoleId:                roleForOU.RoleId.ValueStringPointer(),
				OrganizationalUnitIds: ouIds,
			})
		}
	}

	// Call the Clumio API to create the user using the v2 API.
	res, apiErr := r.sdkUsersV2.CreateUser(&models.CreateUserV2Request{
		AccessControlConfiguration: accessControlConfiguration,
		Email:                      plan.Email.ValueStringPointer(),
		FullName:                   plan.FullName.ValueStringPointer(),
	})
	if apiErr != nil {
		summary := "Error creating Clumio User."
		detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.Id = types.StringPointerValue(res.Id)
	plan.Inviter = types.StringPointerValue(res.Inviter)
	plan.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
	plan.IsEnabled = types.BoolPointerValue(res.IsEnabled)
	plan.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
	plan.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)
	accessControlCfg, assignedRole, ouIds := getAccessControlCfgFromHTTPRes(
		ctx, res.AccessControlConfiguration, &resp.Diagnostics)
	plan.AccessControlConfiguration = accessControlCfg
	plan.AssignedRole = assignedRole
	plan.OrganizationalUnitIds = ouIds

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioUserResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId, perr := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if perr != nil {
		summary := invalidUserMsg
		detail := fmt.Sprintf(invalidUserFmt, state.Id.ValueString())
		resp.Diagnostics.AddError(summary, detail)
	}

	// Call the Clumio API to read the user.
	res, apiErr := r.sdkUsersV2.ReadUser(userId)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio User with ID %s not found. Removing from state.",
				state.Id.ValueString())
			tflog.Warn(ctx, msgStr)
			resp.State.RemoveResource(ctx)
		} else {
			summary := "Error retrieving Clumio User."
			detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
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

	accessControlCfg, assignedRole, ouIds := getAccessControlCfgFromHTTPRes(
		ctx, res.AccessControlConfiguration, &resp.Diagnostics)
	state.AccessControlConfiguration = accessControlCfg
	state.AssignedRole = assignedRole
	state.OrganizationalUnitIds = ouIds

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioUserResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the schema from the current Terraform state.
	var state clumioUserResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (!plan.AssignedRole.IsUnknown() || !plan.OrganizationalUnitIds.IsUnknown()) &&
		!plan.AccessControlConfiguration.IsUnknown() {
		summary := "Error creating Clumio user."
		detail := fmt.Sprintf(errorFmt,
			"Both access_control_configuration and assigned_role/organizational_unit_ids"+
				" cannot be configured. Please configure access_control_configuration only.")
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	if !plan.AssignedRole.IsUnknown() || !plan.OrganizationalUnitIds.IsUnknown() {
		updateRequest := &models.UpdateUserV1Request{}

		if !plan.AssignedRole.IsUnknown() &&
			state.AssignedRole != plan.AssignedRole {
			updateRequest.AssignedRole = plan.AssignedRole.ValueStringPointer()
		}
		if !plan.FullName.IsUnknown() {
			updateRequest.FullName = plan.FullName.ValueStringPointer()
		}
		if !plan.OrganizationalUnitIds.IsUnknown() {
			added := common.SliceDifferenceAttrValue(
				plan.OrganizationalUnitIds.Elements(), state.OrganizationalUnitIds.Elements())
			deleted := common.SliceDifferenceAttrValue(
				state.OrganizationalUnitIds.Elements(), plan.OrganizationalUnitIds.Elements())
			addedStrings := common.GetStringSliceFromAttrValueSlice(added)
			deletedStrings := common.GetStringSliceFromAttrValueSlice(deleted)
			updateRequest.OrganizationalUnitAssignmentUpdates =
				&models.EntityGroupAssignmentUpdatesV1{
					Add:    addedStrings,
					Remove: deletedStrings,
				}
		}

		userId, perr := strconv.ParseInt(plan.Id.ValueString(), 10, 64)
		if perr != nil {
			summary := invalidUserMsg
			detail := fmt.Sprintf(invalidUserFmt, plan.Id.ValueString())
			resp.Diagnostics.AddError(summary, detail)
		}

		// For backwards compatibility purposes we're updating the user via both V1 and V2 API
		// Call the Clumio API to update the user using the old v1 API.
		res, apiErr := r.sdkUsersV1.UpdateUser(userId, updateRequest)
		if apiErr != nil {
			summary := fmt.Sprintf("Error updating Clumio User id: %v.", plan.Id.ValueString())
			detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
			resp.Diagnostics.AddError(summary, detail)
			return
		}
		if res == nil {
			resp.Diagnostics.AddError(
				common.NilErrorMessageSummary, common.NilErrorMessageDetail)
			return
		}

		// Convert the Clumio API response back to a schema and populate all computed fields of the
		// plan. ID however is not updated given that it is the field used to denote which resource to
		// update in the backend.
		plan.Inviter = types.StringPointerValue(res.Inviter)
		plan.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
		plan.IsEnabled = types.BoolPointerValue(res.IsEnabled)
		plan.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
		plan.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)
		plan.AssignedRole = types.StringPointerValue(res.AssignedRole)
		orgUnitIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.AssignedOrganizationalUnitIds)
		resp.Diagnostics.Append(conversionDiags...)
		plan.OrganizationalUnitIds = orgUnitIds
		accessControl := []roleForOrganizationalUnitModel{
			{
				RoleId:                types.StringPointerValue(res.AssignedRole),
				OrganizationalUnitIds: orgUnitIds,
			},
		}
		accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				schemaRoleId: types.StringType,
				schemaOrganizationalUnitIds: types.SetType{
					ElemType: types.StringType,
				},
			},
		}, accessControl)
		resp.Diagnostics.Append(conversionDiags...)
		plan.AccessControlConfiguration = accessControlList

		// Set the schema into the Terraform state.
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		return
	}

	updateRequest := &models.UpdateUserV2Request{}
	if !plan.FullName.IsUnknown() &&
		state.FullName != plan.FullName {
		updateRequest.FullName = plan.FullName.ValueStringPointer()
	}
	if !plan.AccessControlConfiguration.IsUnknown() {
		add, remove := getAccessControlCfgUpdates(
			ctx, state.AccessControlConfiguration.Elements(), plan.AccessControlConfiguration.Elements())
		updateRequest.AccessControlConfigurationUpdates = &models.EntityGroupAssignmentUpdates{
			Add:    add,
			Remove: remove,
		}
	}
	userId, perr := strconv.ParseInt(plan.Id.ValueString(), 10, 64)
	if perr != nil {
		summary := invalidUserMsg
		detail := fmt.Sprintf(invalidUserFmt, plan.Id.ValueString())
		resp.Diagnostics.AddError(summary, detail)
	}

	// For backwards compatibility purposes we're updating the user via both V1 and V2 API
	// Call the Clumio API to update the user using the V2 API.
	res, apiErr := r.sdkUsersV2.UpdateUser(userId, updateRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Error updating Clumio User id: %v.", plan.Id.ValueString())
		detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the
	// plan. ID however is not updated given that it is the field used to denote which resource to
	// update in the backend.
	plan.Inviter = types.StringPointerValue(res.Inviter)
	plan.IsConfirmed = types.BoolPointerValue(res.IsConfirmed)
	plan.IsEnabled = types.BoolPointerValue(res.IsEnabled)
	plan.LastActivityTimestamp = types.StringPointerValue(res.LastActivityTimestamp)
	plan.OrganizationalUnitCount = types.Int64PointerValue(res.OrganizationalUnitCount)

	accessControlCfg, assignedRole, ouIds := getAccessControlCfgFromHTTPRes(
		ctx, res.AccessControlConfiguration, &resp.Diagnostics)
	plan.AccessControlConfiguration = accessControlCfg
	plan.AssignedRole = assignedRole
	plan.OrganizationalUnitIds = ouIds

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes it from the Terraform state.
func (r *clumioUserResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId, perr := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if perr != nil {
		summary := invalidUserMsg
		detail := fmt.Sprintf(invalidUserFmt, state.Id.ValueString())
		resp.Diagnostics.AddError(summary, detail)
	}

	// Call the Clumio API to delete the user.
	_, apiErr := r.sdkUsersV2.DeleteUser(userId)
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Error deleting Clumio User %v.", userId)
			detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done using the ID of the resource.
func (r *clumioUserResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

