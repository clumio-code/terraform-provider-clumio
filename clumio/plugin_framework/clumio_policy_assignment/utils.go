// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_policy_assignment Terraform resource.

package clumio_policy_assignment

import "github.com/clumio-code/clumio-go-sdk/models"

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

// clearOUContext resets the OrganizationalUnitContext in the client.
func (r *clumioPolicyAssignmentResource) clearOUContext() {
	r.client.ClumioConfig.OrganizationalUnitContext = ""
}
