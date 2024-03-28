// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the policy definition SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createPolicy invokes the API to create the policy definition and from the response populates the
// computed attributes of the policy.
func (r *policyResource) createPolicy(
	ctx context.Context, plan *policyResourceModel) diag.Diagnostics {

	// Convert the schema to a Clumio API request to create a policy.
	policyOperations, diags := mapSchemaOperationsToClumioOperations(ctx,
		plan.Operations)
	if diags.HasError() {
		return diags
	}
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
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.ID = types.StringPointerValue(res.Id)
	apiErr, diags = readPolicyAndUpdateModel(ctx, plan, r.sdkPolicyDefinitions)
	if diags.HasError() {
		return diags
	}
	if apiErr != nil {
		summary := fmt.Sprintf(errorPolicyReadMsg, r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	return diags
}

// readPolicy invokes the API to read the policy definition and from the response populates the
// attributes of the policy. If the policy has been removed externally, the function returns "true"
// to indicate to the caller that the resource no longer exists.
func (r *policyResource) readPolicy(
	ctx context.Context, state *policyResourceModel) (bool, diag.Diagnostics) {

	// Call the Clumio API to read the policy and convert the Clumio API response back to a schema
	// and update the state. In addition to computed fields, all fields are populated from the API
	// response in case any values have been changed externally. ID is not updated however given
	// that it is the field used to query the resource from the backend.
	apiErr, diags := readPolicyAndUpdateModel(ctx, state, r.sdkPolicyDefinitions)
	if diags.HasError() {
		return false, diags
	}
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Policy with ID %s not found. Removing from state.",
				state.ID.ValueString())
			tflog.Warn(ctx, msgStr)
			return true, nil
		} else {
			summary := fmt.Sprintf(errorPolicyReadMsg, r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return false, diags
		}
	}
	return false, diags
}

// updatePolicy invokes the API to update the policy definition and from the response populates the
// computed attributes of the policy.
func (r *policyResource) updatePolicy(
	ctx context.Context, plan *policyResourceModel) diag.Diagnostics {

	policyOperations, diags := mapSchemaOperationsToClumioOperations(ctx,
		plan.Operations)
	if diags.HasError() {
		return diags
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
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Since updating a policy is an asynchronous operation, poll till the update is completed.
	err := common.PollTask(ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	// As the policy is updated asynchronously, we need to read the policy after the update is
	// complete to get the updated policy attributes.
	apiErr, diags = readPolicyAndUpdateModel(ctx, plan, r.sdkPolicyDefinitions)
	if diags.HasError() {
		return diags
	}
	if apiErr != nil {
		summary := fmt.Sprintf(errorPolicyReadMsg, r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	return diags
}

// deletePolicy invokes the API to delete the policy definition
func (r *policyResource) deletePolicy(
	ctx context.Context, state *policyResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	// Call the Clumio API to delete the policy.
	res, apiErr := r.sdkPolicyDefinitions.DeletePolicyDefinition(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return diags
		}
		return nil
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Since deleting a policy is an asynchronous operation, poll till the deletion is completed.
	err := common.PollTask(ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}
	return diags
}
