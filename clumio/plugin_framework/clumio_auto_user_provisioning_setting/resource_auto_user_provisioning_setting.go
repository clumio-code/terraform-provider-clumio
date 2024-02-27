// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_auto_user_provisioning_setting
// Terraform resource. This resource is used to enable or disable auto user provisioning setting.

package clumio_auto_user_provisioning_setting

import (
	"context"

	sdkAUPSettings "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_settings"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &autoUserProvisioningSettingResource{}
	_ resource.ResourceWithConfigure = &autoUserProvisioningSettingResource{}
)

// autoUserProvisioningSettingResource is the struct backing the clumio_auto_user_provisioning_setting
// Terraform resource. It holds the Clumio API client and any other required state needed to enable
// or disable auto user provisioning setting.
type autoUserProvisioningSettingResource struct {
	name           string
	client         *common.ApiClient
	sdkAUPSettings sdkAUPSettings.AutoUserProvisioningSettingsV1Client
}

// NewAutoUserProvisioningSettingResource creates a new instance of autoUserProvisioningSettingResource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewAutoUserProvisioningSettingResource() resource.Resource {

	return &autoUserProvisioningSettingResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *autoUserProvisioningSettingResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_auto_user_provisioning_setting"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *autoUserProvisioningSettingResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkAUPSettings = sdkAUPSettings.NewAutoUserProvisioningSettingsV1(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *autoUserProvisioningSettingResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan autoUserProvisioningSettingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: plan.IsEnabled.ValueBoolPointer(),
	}

	// Call the Clumio API to enable or disable the auto user provisioning setting. NOTE that this
	// setting is a singleton state for the entire organization and as such, creation of this
	// resource results in "updating" the current state rather than creating a new instance of a
	// setting.
	_, apiErr := r.sdkAUPSettings.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		summary := "Unable to set auto user provisioning setting."
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Since the API doesn't return an id, we are setting a uuid as the resource id.
	plan.ID = types.StringValue(uuid.New().String())

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *autoUserProvisioningSettingResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the Terraform state.
	var state autoUserProvisioningSettingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the auto user provisioning setting. The setting is an Org-wide
	// value and as such, the Org ID associated with the API credentials will be utilized.
	res, apiErr := r.sdkAUPSettings.ReadAutoUserProvisioningSetting()
	if apiErr != nil {
		summary := "Unable to retrieve auto user provisioning setting."
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and update the state. IsEnabled is the only
	// setting that needs to be refreshed.
	state.IsEnabled = types.BoolPointerValue(res.IsEnabled)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *autoUserProvisioningSettingResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan autoUserProvisioningSettingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	isEnabled := plan.IsEnabled.ValueBool()
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}

	// Call the Clumio API to enable or disable auto user provisioning setting.
	res, apiErr := r.sdkAUPSettings.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		summary := "Unable to update auto user provisioning setting."
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *autoUserProvisioningSettingResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the Terraform state.
	var state autoUserProvisioningSettingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	isEnabled := false
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}

	// Call the Clumio API to disable auto user provisioning setting.
	_, apiErr := r.sdkAUPSettings.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		summary := "Unable to delete auto user provisioning setting."
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
	}
}
