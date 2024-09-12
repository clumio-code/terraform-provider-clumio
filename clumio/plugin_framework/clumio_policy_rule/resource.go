// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the policy rules SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_policy_rule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createPolicyRule invokes the API to create the policy rule and from the response populates the
// computed attributes of the policy rule.
func (r *policyRuleResource) createPolicyRule(
	ctx context.Context, plan *policyRuleResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyRules := r.sdkPolicyRules

	// Convert the schema to a Clumio API request to create a policy rule.
	priority := &models.RulePriority{
		BeforeRuleId: plan.BeforeRuleID.ValueStringPointer(),
	}
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: plan.PolicyID.ValueStringPointer(),
		},
	}
	prRequest := &models.CreatePolicyRuleV1Request{
		Action:    action,
		Condition: plan.Condition.ValueStringPointer(),
		Name:      plan.Name.ValueStringPointer(),
		Priority:  priority,
	}

	// Call the Clumio API to create the policy rule.
	res, apiErr := sdkPolicyRules.CreatePolicyRule(prRequest)
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
	err := common.PollTask(
		ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf(
			"Unable to poll %s (Name: %v) for creation", r.name, plan.Name.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.ID = types.StringPointerValue(res.Rule.Id)

	return diags
}

// readPolicyRule invokes the API to read the policy rule and from the response populates the
// attributes of the policy rule. If the policy rule has been removed externally, the function
// returns "true" to indicate to the caller that the resource no longer exists.
func (r *policyRuleResource) readPolicyRule(ctx context.Context, state *policyRuleResourceModel) (
	bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	sdkPolicyRules := r.sdkPolicyRules

	// Call the Clumio API to read the policy rule.
	res, apiErr := sdkPolicyRules.ReadPolicyRule(state.ID.ValueString())
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state", r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			remove = true
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return remove, diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return false, diags
	}

	// Convert the Clumio API response back to a schema and update the state. In addition to
	// computed fields, all fields are populated from the API response in case any values have been
	// changed externally. ID is not updated however given that it is the field used to query the
	// resource from the backend.
	state.Name = types.StringPointerValue(res.Name)
	state.Condition = types.StringPointerValue(res.Condition)
	if res.Priority != nil {
		state.BeforeRuleID = types.StringPointerValue(res.Priority.BeforeRuleId)
	}
	state.PolicyID = types.StringPointerValue(res.Action.AssignPolicy.PolicyId)

	return false, diags
}

// updateProtectionGroup invokes the API to update the policy rule and from the response populates
// the computed attributes of the policy rule.
func (r *policyRuleResource) updatePolicyRule(
	ctx context.Context, plan *policyRuleResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyRules := r.sdkPolicyRules

	priority := &models.RulePriority{
		BeforeRuleId: plan.BeforeRuleID.ValueStringPointer(),
	}
	policyId := plan.PolicyID.ValueString()
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: &policyId,
		},
	}
	prRequest := &models.UpdatePolicyRuleV1Request{
		Action:    action,
		Condition: plan.Condition.ValueStringPointer(),
		Name:      plan.Name.ValueStringPointer(),
		Priority:  priority,
	}

	// Call the Clumio API to update the policy rule.
	res, apiErr := sdkPolicyRules.UpdatePolicyRule(plan.ID.ValueString(), prRequest)
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

	// As the update of a policy rule is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(
		ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf(
			"Unable to poll %s (ID: %v) for update", r.name, plan.ID.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	return diags
}

// deletePolicyRule invokes the API to delete the policy rule.
func (r *policyRuleResource) deletePolicyRule(
	ctx context.Context, state *policyRuleResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyRules := r.sdkPolicyRules

	// Call the Clumio API to delete the policy rule.
	res, apiErr := sdkPolicyRules.DeletePolicyRule(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}
	// As the delete of a policy rule is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(
		ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf(
			"Unable to poll %s (ID: %v) for deletion", r.name, state.ID.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
	}
	return diags
}
