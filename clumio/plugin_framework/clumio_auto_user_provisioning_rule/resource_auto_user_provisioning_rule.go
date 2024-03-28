// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_auto_user_provisioning_rules Terraform
// resource. This resource is used to create auto user provisioning rules to determine the roles
// and organizational-units to be assigned to the users.

package clumio_auto_user_provisioning_rule

import (
	"context"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &autoUserProvisioningRuleResource{}
	_ resource.ResourceWithConfigure   = &autoUserProvisioningRuleResource{}
	_ resource.ResourceWithImportState = &autoUserProvisioningRuleResource{}
)

// autoUserProvisioningRuleResource is the struct backing the clumio_auto_user_provisioning_rules
// Terraform resource. It holds the Clumio API client and any other required state needed to set up
// the auto user provisioning rules to determine the roles and organizational-units to be assigned
// to the user.
type autoUserProvisioningRuleResource struct {
	name        string
	client      *common.ApiClient
	sdkAUPRules sdkclients.AutoUserProvisioningRuleClient
}

// NewAutoUserProvisioningRuleResource creates a new instance of autoUserProvisioningRuleResource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewAutoUserProvisioningRuleResource() resource.Resource {
	return &autoUserProvisioningRuleResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *autoUserProvisioningRuleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_auto_user_provisioning_rule"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *autoUserProvisioningRuleResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkAUPRules = sdkclients.NewAutoUserProvisioningRuleClient(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *autoUserProvisioningRuleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the current Terraform plan.
	var plan autoUserProvisioningRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createAutoUserProvisioningRule(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *autoUserProvisioningRuleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state autoUserProvisioningRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readAutoUserProvisioningRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *autoUserProvisioningRuleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the current Terraform plan.
	var plan autoUserProvisioningRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateAutoUserProvisioningRule(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *autoUserProvisioningRuleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state autoUserProvisioningRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteAutoUserProvisioningRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *autoUserProvisioningRuleResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
