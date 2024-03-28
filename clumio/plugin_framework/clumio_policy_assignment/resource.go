// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the policy assignment SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createPolicyAssignment invokes the API to create the policy assignment and from the response
// populates the computed attributes of the policy assignment.
func (r *clumioPolicyAssignmentResource) createPolicyAssignment(
	ctx context.Context, plan *policyAssignmentResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups
	sdkPolicyAssignments := r.sdkPolicyAssignments
	sdkPolicyDefinitions := r.sdkPolicyDefinitions

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK clients are temporarily re-initialized in the context of
	// the desired OU so that API calls are made on behalf of the OU.
	if plan.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, plan.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
		sdkPolicyAssignments = sdkclients.NewPolicyAssignmentClient(config)
		sdkPolicyDefinitions = sdkclients.NewPolicyDefinitionClient(config)
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
		if *operation.ClumioType == protectionGroupBackup {
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
	readResponse, apiErr := sdkProtectionGroups.ReadProtectionGroup(*assignment.Entity.Id)
	if apiErr != nil {
		summary := fmt.Sprintf(readProtectionGroupErrFmt, *assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}
	if readResponse.ProtectionInfo == nil ||
		*readResponse.ProtectionInfo.PolicyId != policyId {
		summary := "Protection group policy mismatch"
		detail := fmt.Sprintf("Protection group with id: %s does not have policy %s applied",
			*assignment.Entity.Id, policyId)
		diags.AddError(summary, detail)
		return diags
	}

	// Populate all computed fields of the plan including the ID given that the resource is getting
	// created.
	entityType := plan.EntityType.ValueString()
	plan.ID = types.StringValue(
		fmt.Sprintf("%s_%s_%s", *assignment.PolicyId, *assignment.Entity.Id, entityType))
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)

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

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK clients are temporarily re-initialized in the context of
	// the desired OU so that API calls are made on behalf of the OU.
	if state.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, state.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
		sdkPolicyDefinitions = sdkclients.NewPolicyDefinitionClient(config)
	}

	// Call the Clumio API to read the policy definition.
	policyId := state.PolicyID.ValueString()
	_, apiErr := sdkPolicyDefinitions.ReadPolicyDefinition(policyId, nil)
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

	entityType := state.EntityType.ValueString()
	switch entityType {
	case entityTypeProtectionGroup:
		// Call the Clumio API to read the protection group. Barring any errors, if the protection
		// group is not found or if the protection group no longer has the desired policy attached,
		// the function returns "true" to indicate to the caller that the expected resource no
		// longer exists.
		entityId := state.EntityID.ValueString()
		readResponse, apiErr := sdkProtectionGroups.ReadProtectionGroup(entityId)
		if apiErr != nil {
			remove := false
			if apiErr.ResponseCode == http.StatusNotFound {
				msgStr := fmt.Sprintf(
					"Clumio Protection Group with ID %s not found. Removing from state.",
					entityId)
				tflog.Warn(ctx, msgStr)
				remove = true
			} else {
				summary := fmt.Sprintf(readProtectionGroupErrFmt, entityId)
				detail := common.ParseMessageFromApiError(apiErr)
				diags.AddError(summary, detail)
			}
			return remove, diags
		}
		if readResponse == nil {
			summary := common.NilErrorMessageSummary
			detail := common.NilErrorMessageDetail
			diags.AddError(summary, detail)
			return false, diags
		}
		if readResponse.ProtectionInfo == nil ||
			*readResponse.ProtectionInfo.PolicyId != policyId {
			msgStr := fmt.Sprintf("Protection group with id: %s does not have policy %s applied."+
				" Removing from state.", entityId, policyId)
			tflog.Warn(ctx, msgStr)
			return true, diags
		}
		state.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
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
	sdkProtectionGroups := r.sdkProtectionGroups
	sdkPolicyDefinitions := r.sdkPolicyDefinitions

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK clients are temporarily re-initialized in the context of
	// the desired OU so that API calls are made on behalf of the OU.
	if plan.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, plan.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
		sdkPolicyAssignments = sdkclients.NewPolicyAssignmentClient(config)
		sdkPolicyDefinitions = sdkclients.NewPolicyDefinitionClient(config)
	}

	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := sdkPolicyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read the policy with id : %v", policyId)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	correctPolicyType := false
	for _, operation := range policy.Operations {
		if *operation.ClumioType == protectionGroupBackup {
			correctPolicyType = true
		}
	}

	if !correctPolicyType {
		summary := "Invalid Policy operation."
		detail := fmt.Sprintf(
			"Policy id %s does not contain support protection_group_backup operation",
			policyId)
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

	// Call the Clumio API to read the protection group and verify that the policy is assigned
	// to the protection group.
	readResponse, apiErr := sdkProtectionGroups.ReadProtectionGroup(*assignment.Entity.Id)
	if apiErr != nil {
		summary := fmt.Sprintf(readProtectionGroupErrFmt, *assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	if readResponse.ProtectionInfo == nil ||
		*readResponse.ProtectionInfo.PolicyId != policyId {
		errMsg := fmt.Sprintf(
			"Protection group with id: %s does not have policy %s applied",
			*assignment.Entity.Id, policyId)
		diags.AddError(errMsg, errMsg)
		return diags
	}

	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	return diags
}

// deletePolicyAssignment invokes the API to delete the policy assignment.
func (r *clumioPolicyAssignmentResource) deletePolicyAssignment(
	ctx context.Context, state *policyAssignmentResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkPolicyAssignments := r.sdkPolicyAssignments

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK client is temporarily re-initialized in the context of the
	// desired OU so that API calls are made on behalf of the OU.
	if state.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, state.OrganizationalUnitID.ValueString())
		sdkPolicyAssignments = sdkclients.NewPolicyAssignmentClient(config)
	}

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
