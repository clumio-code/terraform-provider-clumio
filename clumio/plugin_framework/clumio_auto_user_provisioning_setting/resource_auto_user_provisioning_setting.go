// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_auto_user_provisioning_setting
// Terraform resource. This resource is used to enable or disable auto user provisioning setting.

package clumio_auto_user_provisioning_setting

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	sdkAUPSettings sdkclients.AutoUserProvisioningSettingClient
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
	r.sdkAUPSettings = sdkclients.NewAutoUserProvisioningSettingClient(r.client.ClumioConfig)
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

	diags = r.createAutoUserProvisioningSetting(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
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

	diags = r.readAutoUserProvisioningSetting(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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

	diags = r.updateAutoUserProvisioningSetting(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
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

	diags = r.deleteAutoUserProvisioningSetting(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
