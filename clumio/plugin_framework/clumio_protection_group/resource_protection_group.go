// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_protection_group Terraform resource.
// This resource is used to manage protection groups within Clumio.

package clumio_protection_group

import (
	"context"
	"fmt"
	"net/http"

	sdkProtectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioProtectionGroupResource{}
	_ resource.ResourceWithConfigure   = &clumioProtectionGroupResource{}
	_ resource.ResourceWithImportState = &clumioProtectionGroupResource{}
)

// clumioProtectionGroupResource is the struct backing the clumio_protection_group Terraform
// resource. It holds the Clumio API client and any other required state needed to manage protection
// groups within Clumio.
type clumioProtectionGroupResource struct {
	name   string
	client *common.ApiClient
}

// NewClumioProtectionGroupResource creates a new instance of clumioProtectionGroupResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioProtectionGroupResource() resource.Resource {
	return &clumioProtectionGroupResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioProtectionGroupResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_protection_group"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioProtectionGroupResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioProtectionGroupResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioProtectionGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the OrganizationalUnitID is specified, then execute the API in that OrganizationalUnit
	// context.
	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	protectionGroup := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)

	// Call the Clumio API to create the protection group.
	name := plan.Name.ValueString()
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(plan.ObjectFilter)
	response, apiErr := protectionGroup.CreateProtectionGroup(
		models.CreateProtectionGroupV1Request{
			BucketRule:   plan.BucketRule.ValueStringPointer(),
			Description:  plan.Description.ValueStringPointer(),
			Name:         plan.Name.ValueStringPointer(),
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating Protection Group %v.", name),
			fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr)))
		return
	}

	// Poll to read the protection group till it becomes available
	err := pollForProtectionGroup(ctx, *response.Id, protectionGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading the created Protection Group: %v", name),
			fmt.Sprintf(errorFmt, err.Error()))
		return
	}

	// Read the protection group
	plan.ID = types.StringPointerValue(response.Id)
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(plan.ID.ValueString())
	if apiErr != nil {
		summary := fmt.Sprintf(errorProtectionGroupReadFmt, plan.Name.ValueString())
		detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if readResponse == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.Name = types.StringPointerValue(readResponse.Name)
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	plan.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(readResponse.ObjectFilter)
	plan.ProtectionStatus = types.StringPointerValue(readResponse.ProtectionStatus)
	plan.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioProtectionGroupResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioProtectionGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If an organizational unit id is provided, defer clearing the context
	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	protectionGroup := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)

	name := plan.Name.ValueString()
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(plan.ObjectFilter)

	// Call the Clumio API to update the protection group.
	response, apiErr := protectionGroup.UpdateProtectionGroup(plan.ID.ValueString(),
		&models.UpdateProtectionGroupV1Request{
			BucketRule:   plan.BucketRule.ValueStringPointer(),
			Description:  plan.Description.ValueStringPointer(),
			Name:         plan.Name.ValueStringPointer(),
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating Protection Group %v.", name),
			fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr)))
		return
	}
	if response == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Poll to read the protection group till it is updated
	err := pollForProtectionGroup(ctx, *response.Id, protectionGroup)
	if err != nil {
		summary := fmt.Sprintf("Error reading the updated Protection Group: %v", name)
		detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(plan.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(errorProtectionGroupReadFmt, plan.Name.ValueString()),
			fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr)))
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan.
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	plan.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(readResponse.ObjectFilter)
	plan.ProtectionStatus = types.StringPointerValue(readResponse.ProtectionStatus)
	plan.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioProtectionGroupResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioProtectionGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the OrganizationalUnitID is specified, then execute the API in that OrganizationalUnit
	// context.
	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	protectionGroup := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)

	// Call the Clumio API to read the protection group
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Protection Group with ID %s not found. Removing from state.",
				state.ID.ValueString())
			tflog.Warn(ctx, msgStr)
			resp.State.RemoveResource(ctx)
		} else {
			summary := fmt.Sprintf(errorProtectionGroupReadFmt, state.Name.ValueString())
			detail := fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr))
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
	if readResponse == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// If the protection group was deleted externally, throw an error and reset the state
	if readResponse.IsDeleted != nil && *readResponse.IsDeleted {
		msgStr := fmt.Sprintf(
			"Clumio Protection Group with ID %s not found. Removing from state.",
			state.ID.ValueString())
		tflog.Warn(ctx, msgStr)
		resp.State.RemoveResource(ctx)
		return
	}

	// Convert the Clumio API response back to a schema and update the state. An optional field need
	// only be populated if it initially contained a non-null value or if there is a specific value
	// that should be assigned.
	description := types.StringPointerValue(readResponse.Description)
	bucketRule := types.StringPointerValue(readResponse.BucketRule)
	if !state.Description.IsNull() || description.ValueString() != "" {
		state.Description = description
	}
	if !state.BucketRule.IsNull() || bucketRule.ValueString() != "" {
		state.BucketRule = bucketRule
	}
	state.Name = types.StringPointerValue(readResponse.Name)
	state.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	state.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(readResponse.ObjectFilter)
	state.ProtectionStatus = types.StringPointerValue(readResponse.ProtectionStatus)
	state.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioProtectionGroupResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioProtectionGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the OrganizationalUnitID is specified, then execute the API in that OrganizationalUnit
	// context.
	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	protectionGroup := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)

	// Call the Clumio API to delete the protection group
	_, apiErr := protectionGroup.DeleteProtectionGroup(state.ID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting Protection Group %v.", state.Name.ValueString()),
			fmt.Sprintf(errorFmt, common.ParseMessageFromApiError(apiErr)))
		return
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioProtectionGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
