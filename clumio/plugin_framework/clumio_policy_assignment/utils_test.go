// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy_assignment

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following PolicyAssignment mapping cases:
//   - Mapping corresponding to assigning a policy to entity.
//   - Mapping corresponding to un-assigning the policy from the entity.
func TestMapSchemaPolicyAssignmentToClumioPolicyAssignment(t *testing.T) {

	model := policyAssignmentResourceModel{
		EntityID:   basetypes.NewStringValue(entityId),
		EntityType: basetypes.NewStringValue(entityTypeProtectionGroup),
		PolicyID:   basetypes.NewStringValue(policyId),
	}

	// Tests that the schema policy assignment gets converted to a SDK model policy assignment
	// for the entity assign scenario.
	t.Run("assign scenario mapping", func(t *testing.T) {
		clumioPA := mapSchemaPolicyAssignmentToClumioPolicyAssignment(model, false)
		assert.Equal(t, 1, len(clumioPA.Items))
		assert.Equal(t, *clumioPA.Items[0].Action, actionAssign)
		assert.Equal(t, *clumioPA.Items[0].Entity.Id, entityId)
		assert.Equal(t, *clumioPA.Items[0].PolicyId, policyId)
	})

	// Tests that the schema policy assignment gets converted to a SDK model policy assignment
	// for the entity unassign scenario.
	t.Run("unassign scenario mapping", func(t *testing.T) {
		clumioPA := mapSchemaPolicyAssignmentToClumioPolicyAssignment(model, true)
		assert.Equal(t, 1, len(clumioPA.Items))
		assert.Equal(t, *clumioPA.Items[0].Action, actionUnassign)
		assert.Equal(t, *clumioPA.Items[0].Entity.Id, entityId)
		assert.Equal(t, *clumioPA.Items[0].PolicyId, policyIdEmpty)
	})

}
