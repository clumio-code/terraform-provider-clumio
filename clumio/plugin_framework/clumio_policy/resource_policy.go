// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_policy Terraform resource.
// This resource is used to create a policy for scheduling backups on Clumio supported data sources.

package clumio_policy

import (
	"context"
	"fmt"
	sdkPolicyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
)

// policyResource is the struct backing the clumio_policy Terraform resource. It holds the Clumio
// API client and any other required state needed to create a Clumio Policy.
type policyResource struct {
	name                 string
	client               *common.ApiClient
	sdkPolicyDefinitions sdkPolicyDefinitions.PolicyDefinitionsV1Client
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
	r.sdkPolicyDefinitions = sdkPolicyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
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

	// Convert the schema to a Clumio API request to create a policy.
	policyOperations, diags := mapSchemaOperationsToClumioOperations(ctx,
		plan.Operations)
	resp.Diagnostics.Append(diags...)
	pdRequest := &models.CreatePolicyDefinitionV1Request{
		ActivationStatus:     plan.ActivationStatus.ValueStringPointer(),
		Name:                 plan.Name.ValueStringPointer(),
		Timezone:             plan.Timezone.ValueStringPointer(),
		Operations:           policyOperations,
		OrganizationalUnitId: plan.OrganizationalUnitId.ValueStringPointer(),
	}

	// Call the Clumio API to create the policy.
	res, apiErr := r.sdkPolicyDefinitions.CreatePolicyDefinition(pdRequest)
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
	plan.ID = types.StringPointerValue(res.Id)
	apiErr, diags = readPolicyAndUpdateModel(ctx, &plan, r.sdkPolicyDefinitions)
	resp.Diagnostics.Append(diags...)
	if apiErr != nil {
		summary := fmt.Sprintf(errorPolicyReadMsg, r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
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

	// Call the Clumio API to read the policy and convert the Clumio API response back to a schema
	// and update the state. In addition to computed fields, all fields are populated from the API
	// response in case any values have been changed externally. ID is not updated however given
	// that it is the field used to query the resource from the backend.
	apiErr, diags := readPolicyAndUpdateModel(ctx, &state, r.sdkPolicyDefinitions)
	resp.Diagnostics.Append(diags...)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Policy with ID %s not found. Removing from state.",
				state.ID.ValueString())
			tflog.Warn(ctx, msgStr)
			resp.State.RemoveResource(ctx)
		} else {
			summary := fmt.Sprintf(errorPolicyReadMsg, r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
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

	policyOperations, policyDiag := mapSchemaOperationsToClumioOperations(ctx,
		plan.Operations)
	if policyDiag != nil {
		resp.Diagnostics.Append(policyDiag...)
	}
	pdRequest := &models.UpdatePolicyDefinitionV1Request{
		ActivationStatus:     plan.ActivationStatus.ValueStringPointer(),
		Name:                 plan.Name.ValueStringPointer(),
		Timezone:             plan.Timezone.ValueStringPointer(),
		Operations:           policyOperations,
		OrganizationalUnitId: plan.OrganizationalUnitId.ValueStringPointer(),
	}

	// Call the Clumio API to update the policy.
	res, apiErr := r.sdkPolicyDefinitions.UpdatePolicyDefinition(
		plan.ID.ValueString(), nil, pdRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Since updating a policy is an asynchronous operation, poll till the update is completed.
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// As the policy is updated asynchronously, we need to read the policy after the update is
	// complete to get the updated policy attributes.
	apiErr, diags = readPolicyAndUpdateModel(ctx, &plan, r.sdkPolicyDefinitions)
	resp.Diagnostics.Append(diags...)
	if apiErr != nil {
		summary := fmt.Sprintf(errorPolicyReadMsg, r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
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
	res, apiErr := r.sdkPolicyDefinitions.DeletePolicyDefinition(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
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

	// Since deleting a policy is an asynchronous operation, poll till the deletion is completed.
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
