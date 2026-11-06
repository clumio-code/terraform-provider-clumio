// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_policy Terraform resource.
// This resource is used to create a policy for scheduling backups on Clumio supported data sources.

package clumio_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
	_ resource.ResourceWithModifyPlan  = &policyResource{}
)

// policyResource is the struct backing the clumio_policy Terraform resource. It holds the Clumio
// API client and any other required state needed to create a Clumio Policy.
type policyResource struct {
	name                 string
	client               *common.ApiClient
	sdkPolicyDefinitions sdkclients.PolicyDefinitionClient
	sdkTasks             sdkclients.TaskClient
	pollTimeout          time.Duration
	pollInterval         time.Duration
}

// NewPolicyResource creates a new instance of policyResource. Its attributes are initialized later
// by Terraform via Metadata and Configure once the Provider is initialized.
func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *policyResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_policy"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *policyResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkPolicyDefinitions = sdkclients.NewPolicyDefinitionClient(r.client.ClumioConfig)
	r.sdkTasks = sdkclients.NewTaskClient(r.client.ClumioConfig)
	r.pollTimeout = 3600 * time.Second
	r.pollInterval = 5 * time.Second
}

// ModifyPlan rejects an operation activation_status that the policy-level status overrides.
func (r *policyResource) ModifyPlan(
	ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	if req.Plan.Raw.IsNull() {
		return
	}
	var planStatus, stateStatus types.String
	var operations types.Set
	resp.Diagnostics.Append(
		req.Plan.GetAttribute(ctx, path.Root(schemaActivationStatus), &planStatus)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(schemaOperations), &operations)...)
	if !req.State.Raw.IsNull() {
		resp.Diagnostics.Append(
			req.State.GetAttribute(ctx, path.Root(schemaActivationStatus), &stateStatus)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var overridden, detail string
	switch {
	case planStatus.ValueString() == activationStatusDectivated:
		overridden = activationStatusActivated
		detail = "A deactivated policy deactivates every operation. Remove the operation" +
			" activation_status, or activate the policy."
	case planStatus.ValueString() == activationStatusActivated &&
		stateStatus.ValueString() == activationStatusDectivated:
		overridden = activationStatusDectivated
		detail = "Activating a deactivated policy activates every operation. Remove the" +
			" operation activation_status for this apply, then set it in a later apply."
	}

	typeCount := make(map[string]int)
	var statuses []types.String
	var opTypes []string
	for _, element := range operations.Elements() {
		operation, ok := element.(types.Object)
		if !ok {
			continue
		}
		opType, _ := operation.Attributes()[schemaOperationType].(types.String)
		status, _ := operation.Attributes()[schemaActivationStatus].(types.String)
		typeCount[opType.ValueString()]++
		opTypes = append(opTypes, opType.ValueString())
		statuses = append(statuses, status)
	}
	for i, status := range statuses {
		if status.IsNull() {
			continue
		}
		if typeCount[opTypes[i]] > 1 {
			resp.Diagnostics.AddAttributeError(path.Root(schemaOperations),
				"Operation activation_status on a duplicated operation type",
				fmt.Sprintf("Operation type %q occurs more than once. Remove activation_status"+
					" from these operations.", opTypes[i]))
			return
		}
		if overridden != "" && status.ValueString() == overridden {
			resp.Diagnostics.AddAttributeError(path.Root(schemaOperations),
				"Operation activation_status conflicts with the policy activation_status", detail)
			return
		}
	}
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *policyResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to create the policy.
	diags = r.createPolicy(ctx, &plan)
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
func (r *policyResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the policy.
	remove, diags := r.readPolicy(ctx, &state)
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
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *policyResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to update the policy.
	diags = r.updatePolicy(ctx, &plan)
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
func (r *policyResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to delete the policy.
	diags = r.deletePolicy(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
