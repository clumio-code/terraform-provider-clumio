// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_policy_rule Terraform resource.
// This resource facilitates the creation of rules that specify which assets are to be protected
// by particular policies.

package clumio_policy_rule

import (
	"context"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &policyRuleResource{}
	_ resource.ResourceWithConfigure   = &policyRuleResource{}
	_ resource.ResourceWithImportState = &policyRuleResource{}
)

// policyRuleResource is the struct backing the clumio_policy_rule Terraform resource. It holds the
// Clumio API client and any other required state needed to create a Clumio policy rule.
type policyRuleResource struct {
	name           string
	client         *common.ApiClient
	sdkPolicyRules sdkclients.PolicyRuleClient
	sdkTasks       sdkclients.TaskClient
	pollInterval   time.Duration
	pollTimeout    time.Duration
}

// NewPolicyRuleResource creates a new instance of policyRuleResource. Its attributes are
// initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewPolicyRuleResource() resource.Resource {
	return &policyRuleResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *policyRuleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_policy_rule"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *policyRuleResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkPolicyRules = sdkclients.NewPolicyRuleClient(r.client.ClumioConfig)
	r.sdkTasks = sdkclients.NewTaskClient(r.client.ClumioConfig)
	r.pollTimeout = 3600 * time.Second
	r.pollInterval = 5 * time.Second
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *policyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *policyRuleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createPolicyRule(ctx, &plan)
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
func (r *policyRuleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readPolicyRule(ctx, &state)
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *policyRuleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updatePolicyRule(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes it from the Terraform state.
func (r *policyRuleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deletePolicyRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
