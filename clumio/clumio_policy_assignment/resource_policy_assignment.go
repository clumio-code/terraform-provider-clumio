// Copyright 2022. Clumio, Inc.

// clumio_policy_assignment definition and CRUD implementation.

package clumio_policy_assignment

import (
	"context"
	"fmt"
	policyAssignments "github.com/clumio-code/clumio-go-sdk/controllers/policy_assignments"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioPolicyAssignment returns the resource for Clumio Policy Assignment.
func ClumioPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Assignment Resource used to assign (or unassign)" +
			" policies.\n\n NOTE: Currently policy assignment is supported only for" +
			" entity type \"protection_group\".",

		CreateContext: clumioPolicyAssignmentCreateUpdate,
		UpdateContext: clumioPolicyAssignmentCreateUpdate,
		ReadContext:   clumioPolicyAssignmentRead,
		DeleteContext: clumioPolicyAssignmentDelete,

		Schema: map[string]*schema.Schema{
			schemaEntityId: {
				Description: "The entity id.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			schemaEntityType: {
				Description: "The entity type. The supported entity type is" +
					"\"protection_group\".",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					_, ok := validEntityTypeMap[i.(string)]
					if !ok {
						validEntityTypes := make([]string, 0)
						for key := range validEntityTypeMap {
							validEntityTypes = append(validEntityTypes, key)
						}
						return diag.Errorf("Valid EntityTypes are %v", validEntityTypes)
					}
					return nil
				},
			},
			schemaPolicyId: {
				Description: "The Clumio-assigned ID of the policy. ",
				Type:        schema.TypeString,
				Required:    true,
			},
			schemaOrganizationalUnitId: {
				Type: schema.TypeString,
				Description: "The Clumio-assigned ID of the organizational unit" +
					" to use as the context for assigning the policy.",
				Optional: true,
			},
		},
	}
}

// clumioPolicyAssignmentCreateUpdate handles the Create/Update action for the Clumio Policy Assignment Resource.
func clumioPolicyAssignmentCreateUpdate(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	pa := policyAssignments.NewPolicyAssignmentsV1(clumioConfig)
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(d, false)
	_, apiErr := pa.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		return diag.Errorf("Error assigning policy %v to entity %v. Error: %v",
			*assignment.PolicyId, *assignment.Entity.Id, string(apiErr.Response))
	}
	id := fmt.Sprintf("%s_%s", *assignment.PolicyId, *assignment.Entity.Id)
	d.SetId(id)
	return nil
}

// clumioPolicyAssignmentDelete handles the Delete action for the Clumio Policy Assignment Resource.
func clumioPolicyAssignmentDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	pa := policyAssignments.NewPolicyAssignmentsV1(clumioConfig)
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(d, true)
	_, apiErr := pa.SetPolicyAssignments(paRequest)
	if apiErr != nil {
		assignment := paRequest.Items[0]
		return diag.Errorf("Error unassigning policy from entity %v. Error: %v",
			*assignment.Entity.Id, string(apiErr.Response))
	}
	return nil
}

// clumioPolicyAssignmentRead handles the Read action for the Clumio Policy Assignment Resource.
func clumioPolicyAssignmentRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// mapSchemaPolicyAssignmentToClumioPolicyAssignment maps the schema policy assignment
// to the Clumio API request policy assignment.
func mapSchemaPolicyAssignmentToClumioPolicyAssignment(
	d *schema.ResourceData, unassign bool) *models.SetPolicyAssignmentsV1Request {
	entityId := common.GetStringValue(d, schemaEntityId)
	entityType := common.GetStringValue(d, schemaEntityType)
	entity := &models.AssignmentEntity{
		Id:         &entityId,
		ClumioType: &entityType,
	}

	policyId := common.GetStringValue(d, schemaPolicyId)
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
