// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the policy assignment SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createPolicyAssignment invokes the API to create the policy assignment and from the response
// populates the computed attributes of the policy assignment.
func (r *clumioPolicyAssignmentResource) createPolicyAssignment(
	ctx context.Context, plan *policyAssignmentResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyAssignments := r.sdkPolicyAssignments
	sdkPolicyDefinitions := r.sdkPolicyDefinitions
	policyOperationType := protectionGroupBackup
	if plan.EntityType.ValueString() == entityTypeAWSDynamoDBTable {
		policyOperationType = dynamodbTableBackup
	}
	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := sdkPolicyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read policy with id: %v ", policyId)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	correctPolicyType := false
	for _, operation := range policy.Operations {
		if *operation.ClumioType == policyOperationType {
			correctPolicyType = true
		}
	}
	if !correctPolicyType {
		summary := "Invalid Policy operation."
		detail := fmt.Sprintf(
			"Policy id %s does not contain support protection_group_backup operation", policyId)
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the schema to a Clumio API request to set policy assignments.
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(*plan, false)

	// Call the Clumio API to set the policy assignments
	res, apiErr := sdkPolicyAssignments.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to assign policy %v to entity %v ", policyId,
			*assignment.Entity.Id)
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

	// As setting policy assignments is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll task after assigning policy %v to entity %v",
			policyId, *assignment.Entity.Id)
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	// Populate all computed fields of the plan including the ID given that the resource is getting
	// created.
	entityType := plan.EntityType.ValueString()
	plan.ID = types.StringValue(
		fmt.Sprintf("%s_%s_%s", *assignment.PolicyId, *assignment.Entity.Id, entityType))

	return diags
}

// readPolicyAssignment invokes the APIs to read the policy definition and entity corresponding to
// the policyId and entityId in the state and checks if the entity has the policy applied. If either
// the entity or the policy definition has been removed externally or if the policy applied on the
// entity is different to the policy in the state, then the function returns "true" to indicate to
// the caller that the resource no longer exists.
func (r *clumioPolicyAssignmentResource) readPolicyAssignment(
	ctx context.Context, state *policyAssignmentResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups
	sdkPolicyDefinitions := r.sdkPolicyDefinitions
	sdkDynamoDBTables := r.sdkDynamoDBTables

	// Call the Clumio API to read the policy definition.
	policyId := state.PolicyID.ValueString()
	policy, apiErr := sdkPolicyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Policy with ID %s not found. Removing from state.",
				policyId)
			tflog.Warn(ctx, msgStr)
			remove = true
		} else {
			summary := fmt.Sprintf("Unable to read policy %v.", policyId)
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return remove, diags
	}
	policyOperationType := protectionGroupBackup
	if state.EntityType.ValueString() == entityTypeAWSDynamoDBTable {
		policyOperationType = dynamodbTableBackup
	}
	correctPolicyType := false
	for _, operation := range policy.Operations {
		if *operation.ClumioType == policyOperationType {
			correctPolicyType = true
		}
	}
	if !correctPolicyType {
		msgStr := fmt.Sprintf(
			"Policy does not support required policy operation: %s", policyOperationType)
		tflog.Warn(ctx, msgStr)
		return true, diags
	}

	entityType := state.EntityType.ValueString()
	switch entityType {
	case entityTypeProtectionGroup:
		return r.readAndValidateProtectionGroup(ctx, sdkProtectionGroups, state, policyId)
	case entityTypeAWSDynamoDBTable:
		return r.readAndValidateDynamoDBTable(ctx, sdkDynamoDBTables, state, policyId)
	default:
		summary := "Invalid entityType"
		detail := fmt.Sprintf("The entity type %v is not supported.", entityType)
		diags.AddError(summary, detail)
	}
	return false, diags
}

// updatePolicyAssignment invokes the API to update the policy assignment and from the response
// populates the computed attributes of the policy assignment. After update is done, it also
// verifies that the policy has been applied on the entity.
func (r *clumioPolicyAssignmentResource) updatePolicyAssignment(
	ctx context.Context, plan *policyAssignmentResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyAssignments := r.sdkPolicyAssignments
	sdkPolicyDefinitions := r.sdkPolicyDefinitions

	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := sdkPolicyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read the policy with id : %v", policyId)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	policyOperationType := protectionGroupBackup
	if plan.EntityType.ValueString() == entityTypeAWSDynamoDBTable {
		policyOperationType = dynamodbTableBackup
	}
	correctPolicyType := false
	for _, operation := range policy.Operations {
		if *operation.ClumioType == policyOperationType {
			correctPolicyType = true
		}
	}

	if !correctPolicyType {
		summary := "Invalid Policy operation."
		detail := fmt.Sprintf(
			"Policy id %s does not contain support %s operation", policyId, policyOperationType)
		diags.AddError(summary, detail)
		return diags
	}

	// Call the Clumio API to update the policy assignments.
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(*plan, false)
	res, apiErr := sdkPolicyAssignments.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to assign policy %v to entity %v", policyId,
			*assignment.Entity.Id)
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

	// As setting policy assignments is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll task after assigning policy %v to entity %v",
			policyId, *assignment.Entity.Id)
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	return diags
}

// deletePolicyAssignment invokes the API to delete the policy assignment.
func (r *clumioPolicyAssignmentResource) deletePolicyAssignment(
	ctx context.Context, state *policyAssignmentResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyAssignments := r.sdkPolicyAssignments

	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(*state, true)
	// Call the Clumio API to remove the policy assignment.
	res, apiErr := sdkPolicyAssignments.SetPolicyAssignments(paRequest)
	if apiErr != nil {
		assignment := paRequest.Items[0]
		summary := fmt.Sprintf(
			"Unable to unassign policy from entity %v.", *assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// As setting policy assignments is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.sdkTasks, *res.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll task after unassigning policy %v to entity %v",
			state.PolicyID.ValueString(), state.EntityID.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
	}

	return diags
}
