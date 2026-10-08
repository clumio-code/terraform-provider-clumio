// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_policy_assignment Terraform resource.

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// mapSchemaPolicyAssignmentToClumioPolicyAssignment maps the schema policy assignment
// to the Clumio API request policy assignment.
func mapSchemaPolicyAssignmentToClumioPolicyAssignment(
	model policyAssignmentResourceModel,
	unassign bool) *models.SetPolicyAssignmentsV1Request {

	entityId := model.EntityID.ValueString()
	entityType := model.EntityType.ValueString()
	entity := &models.AssignmentEntity{
		Id:         &entityId,
		ClumioType: &entityType,
	}

	policyId := model.PolicyID.ValueString()
	action := actionAssign
	if unassign {
		policyId = policyIdEmpty
		action = actionUnassign
	}

	assignmentInput := &models.AssignmentInputModel{
		Action:   &action,
		Entity:   entity,
		PolicyId: &policyId,
	}
	return &models.SetPolicyAssignmentsV1Request{
		Items: []*models.AssignmentInputModel{
			assignmentInput,
		},
	}
}

// validateAssignment returns true when the entity is gone or no longer has the policy applied,
// so that the caller removes the assignment from state.
func validateAssignment[T any](ctx context.Context, entityName, entityId, policyId string,
	resp *T, apiErr *apiutils.APIError, assignedPolicyId func(*T) *string) (
	bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf(
				"%s with ID %s not found. Removing from state.", entityName, entityId))
			return true, diags
		}
		diags.AddError(fmt.Sprintf("Unable to read %s %v.", entityName, entityId),
			common.ParseMessageFromApiError(apiErr))
		return false, diags
	}
	if resp == nil {
		diags.AddError(common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return false, diags
	}
	if id := assignedPolicyId(resp); id == nil || *id != policyId {
		tflog.Warn(ctx, fmt.Sprintf("%s with ID %s does not have policy %s applied."+
			" Removing from state.", entityName, entityId, policyId))
		return true, diags
	}
	return false, diags
}

// protectionPolicyId returns the ID of the policy that protects the entity, or nil.
func protectionPolicyId(info *models.ProtectionInfoWithRule) *string {
	if info == nil {
		return nil
	}
	return info.PolicyId
}
