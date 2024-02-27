// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_organizational_unit Terraform resource.
// This resource is used to manage organizational units within Clumio.

package clumio_organizational_unit

import (
	"context"
	"fmt"
	"net/http"

	sdkOrgUnits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioOrganizationalUnitResource{}
	_ resource.ResourceWithConfigure   = &clumioOrganizationalUnitResource{}
	_ resource.ResourceWithImportState = &clumioOrganizationalUnitResource{}
)

// clumioOrganizationalUnitResource is the struct backing the clumio_organizational_unit Terraform resource.
// It holds the Clumio API client and any other required state needed to
// manage organizational units within Clumio.
type clumioOrganizationalUnitResource struct {
	name        string
	client      *common.ApiClient
	sdkOrgUnits sdkOrgUnits.OrganizationalUnitsV2Client
}

// NewClumioOrganizationalUnitResource creates a new instance of clumioOrganizationalUnitResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioOrganizationalUnitResource() resource.Resource {
	return &clumioOrganizationalUnitResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioOrganizationalUnitResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_organizational_unit"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioOrganizationalUnitResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkOrgUnits = sdkOrgUnits.NewOrganizationalUnitsV2(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioOrganizationalUnitResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioOrganizationalUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	plan.Id = id
	plan.ParentId = parentIdString
	plan.ChildrenCount = childrenCount
	plan.UserCount = userCount

	configuredDataTypes, conversionDiags := types.ListValueFrom(ctx, types.StringType, configuredDatasourceTypes)
	resp.Diagnostics.Append(conversionDiags...)
	plan.ConfiguredDatasourceTypes = configuredDataTypes

	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, descendantIdSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.DescendantIds = descendantIds

	// Set the schema into the Terraform state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioOrganizationalUnitResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioOrganizationalUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the orgainizational unit.
	res, apiErr := r.sdkOrgUnits.ReadOrganizationalUnit(state.Id.ValueString(), nil)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state", r.name, state.Id.ValueString())
			tflog.Warn(ctx, summary)
			resp.State.RemoveResource(ctx)
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.Id.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
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

	configuredDataTypes, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.ConfiguredDatasourceTypes)
	resp.Diagnostics.Append(conversionDiags...)
	state.ConfiguredDatasourceTypes = configuredDataTypes

	userSlice, userWithRoleSlice := getUsersFromHTTPRes(res.Users)
	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	resp.Diagnostics.Append(conversionDiags...)
	state.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	resp.Diagnostics.Append(conversionDiags...)
	state.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.DescendantIds)
	resp.Diagnostics.Append(conversionDiags...)
	state.DescendantIds = descendantIds

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioOrganizationalUnitResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioOrganizationalUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the schema from the current Terraform state.
	var state clumioOrganizationalUnitResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	var parentIdString types.String
	var childrenCount types.Int64
	var userCount types.Int64
	var configuredDatasourceTypes []*string
	var userWithRoleSlice []userWithRole
	var userSlice []*string
	var descendantIdSlice []*string

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	// API call can result in either 200 or 201 status code. Relevant data is returned inside a field
	// name mapped to the status code. Data is extracted from the correct field according to the
	// status code. Else return an empty response error.
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
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	plan.ParentId = parentIdString
	plan.ChildrenCount = childrenCount
	plan.UserCount = userCount

	configuredDataTypes, conversionDiags := types.ListValueFrom(ctx, types.StringType, configuredDatasourceTypes)
	resp.Diagnostics.Append(conversionDiags...)
	plan.ConfiguredDatasourceTypes = configuredDataTypes

	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, descendantIdSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.DescendantIds = descendantIds

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioOrganizationalUnitResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioOrganizationalUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to delete the organizational unit.
	res, apiErr := r.sdkOrgUnits.DeleteOrganizationalUnit(state.Id.ValueString(), nil)
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.Id.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, pollTimeoutInSec, pollIntervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll %s (ID: %v) for deletion", r.name, state.Id.ValueString())
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioOrganizationalUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
