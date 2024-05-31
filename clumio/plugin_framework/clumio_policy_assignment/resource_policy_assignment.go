// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_policy_assignment Terraform resource.
// This resource is used to assign a policy to an entity (for example, a protection_group).

package clumio_policy_assignment

import (
	"context"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"time"
)

var (
	_ resource.Resource              = &clumioPolicyAssignmentResource{}
	_ resource.ResourceWithConfigure = &clumioPolicyAssignmentResource{}
)

// clumioPolicyAssignmentResource is the struct backing the clumio_policy_assignment Terraform resource.
// It holds the Clumio API client and any other required state needed to do policy assignment.
type clumioPolicyAssignmentResource struct {
	name                 string
	client               *common.ApiClient
	sdkPolicyDefinitions sdkclients.PolicyDefinitionClient
	sdkProtectionGroups  sdkclients.ProtectionGroupClient
	sdkPolicyAssignments sdkclients.PolicyAssignmentClient
	sdkDynamoDBTables    sdkclients.DynamoDBTableClient
	sdkTasks             sdkclients.TaskClient
	pollTimeout          time.Duration
	pollInterval         time.Duration
}

// NewPolicyAssignmentResource creates a new instance of clumioPolicyAssignmentResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewPolicyAssignmentResource() resource.Resource {
	return &clumioPolicyAssignmentResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioPolicyAssignmentResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_policy_assignment"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioPolicyAssignmentResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkPolicyDefinitions = sdkclients.NewPolicyDefinitionClient(r.client.ClumioConfig)
	r.sdkProtectionGroups = sdkclients.NewProtectionGroupClient(r.client.ClumioConfig)
	r.sdkPolicyAssignments = sdkclients.NewPolicyAssignmentClient(r.client.ClumioConfig)
	r.sdkTasks = sdkclients.NewTaskClient(r.client.ClumioConfig)
	r.pollTimeout = 300 * time.Second
	r.pollInterval = 5 * time.Second
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioPolicyAssignmentResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyAssignmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createPolicyAssignment(ctx, &plan)
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
func (r *clumioPolicyAssignmentResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyAssignmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readPolicyAssignment(ctx, &state)
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}
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
func (r *clumioPolicyAssignmentResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyAssignmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updatePolicyAssignment(ctx, &plan)
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

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioPolicyAssignmentResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyAssignmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deletePolicyAssignment(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
