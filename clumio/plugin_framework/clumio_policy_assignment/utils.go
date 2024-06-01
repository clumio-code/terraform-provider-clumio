// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_policy_assignment Terraform resource.

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

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

// readAndValidateDynamoDBTable reads the Protection Group and validates that the given policy is
// assigned to the Protection Group.
func (r *clumioPolicyAssignmentResource) readAndValidateProtectionGroup(ctx context.Context,
	sdkProtectionGroups sdkclients.ProtectionGroupClient, state *policyAssignmentResourceModel,
	policyId string) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics
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
	return false, diags
}

// readAndValidateDynamoDBTable reads the DynamoDB table and validates that the given policy is
// assigned to the DynamoDB table.
func (r *clumioPolicyAssignmentResource) readAndValidateDynamoDBTable(ctx context.Context,
	sdkDynamoDBTables sdkclients.DynamoDBTableClient, state *policyAssignmentResourceModel,
	policyId string) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	// Call the Clumio API to read the DynamoDB table. Barring any errors, if the DynamoDB
	// table is not found or if the DynamoDB table no longer has the desired policy attached,
	// the function returns "true" to indicate to the caller that the expected resource no
	// longer exists.
	entityId := state.EntityID.ValueString()
	readResponse, apiErr := sdkDynamoDBTables.ReadAwsDynamodbTable(entityId, nil)
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"DynamoDB table with ID %s not found. Removing from state.",
				entityId)
			tflog.Warn(ctx, msgStr)
			remove = true
		} else {
			summary := fmt.Sprintf(readDynamoDBTableErrFmt, entityId)
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
		msgStr := fmt.Sprintf("DynamoDB table with id: %s does not have policy %s applied."+
			" Removing from state.", entityId, policyId)
		tflog.Warn(ctx, msgStr)
		return true, diags
	}
	return false, diags
}
