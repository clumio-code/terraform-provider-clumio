// Copyright 2025. Clumio, Inc.

// This file holds the resource implementation for the clumio_general_settings Terraform resource.
// This resource manages general organizational settings like auto-logout duration, password expiration, and IP allowlists.

package clumio_general_settings

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &clumioGeneralSettings{}
	_ resource.ResourceWithConfigure = &clumioGeneralSettings{}
)

// clumioGeneralSettings is the struct backing the clumio_general_settings Terraform resource.
// It holds the Clumio API client.
type clumioGeneralSettings struct {
	name               string
	client             *common.ApiClient
	sdkGeneralSettings sdkclients.GeneralSettingsClient
}

// NewGeneralSettingsResource creates a new instance of clumioGeneralSettings. Its attributes are
// initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewGeneralSettingsResource() resource.Resource {
	return &clumioGeneralSettings{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioGeneralSettings) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_general_settings"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioGeneralSettings) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkGeneralSettings = sdkclients.NewGeneralSettingsClient(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioGeneralSettings) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan generalSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource is a no-op. We just call the update API to set the values.
	diags = r.updateGeneralSettings(ctx, &plan)
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
func (r *clumioGeneralSettings) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state generalSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readGeneralSettings(ctx, &state)
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

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioGeneralSettings) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan generalSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateGeneralSettings(ctx, &plan)
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

// Delete resets the resource via the Clumio API and removes the Terraform state.
func (r *clumioGeneralSettings) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Call the Clumio API to reset the general settings to defaults.
	resetReq := generalSettingsResourceModel{
		AutoLogoutDuration:         types.Int64Value(defaultAutoLogoutDuration),
		IpAllowlist:                []types.String{types.StringPointerValue(&defaultIPAllow)},
		PasswordExpirationDuration: types.Int64Value(defaultPasswordExpiration),
	}
	diags := r.updateGeneralSettings(ctx, &resetReq)
	resp.Diagnostics.Append(diags...)
}
