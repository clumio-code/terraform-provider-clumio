// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy_assignment

import (
	"testing"

	"github.com/clumio-code/clumio-go-sdk/models"
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

// Unit test for the various types of operations are supported with policy assignment.
func TestIsOperationsSupported(t *testing.T) {

	operationDDB := dynamodbTableBackup
	operationS3 := protectionGroupBackup
	operationBacktrack := awsS3Backtrack
	operationContinuous := awsS3Continuous
	operationEBS := "aws_ebs_backup"
	modelOperations := []*models.PolicyOperation{
		{
			ClumioType: &operationDDB,
		},
	}

	// Tests that the DynamoDB table entity allows DynamoDB table backup type operation.
	t.Run("Allow DynamoDB policy assignment", func(t *testing.T) {
		assert.Nil(t, isOperationsSupported(entityTypeAWSDynamoDBTable, policyId, modelOperations))
	})

	// Tests that the Protection Group entity allows S3 backup type operation.
	t.Run("Allow S3 backup policy assignment", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationS3
		assert.Nil(t, isOperationsSupported(entityTypeProtectionGroup, policyId, modelOperations))
	})

	// Tests that the Protection Group entity allows S3 backtrack type operation.
	t.Run("Allow S3 backtrack policy assignment", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationBacktrack
		assert.Nil(t, isOperationsSupported(entityTypeProtectionGroup, policyId, modelOperations))
	})

	// Tests that the Protection Group entity allows S3 continuous backup type operation.
	t.Run("Allow S3 continuous policy assignment", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationContinuous
		assert.Nil(t, isOperationsSupported(entityTypeProtectionGroup, policyId, modelOperations))
	})

	// Tests that the Protection Group entity does not allow DyanmoDB table backup type operation.
	t.Run("Inhibit different type of assignment of PG", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationDDB
		diags := isOperationsSupported(entityTypeProtectionGroup, policyId, modelOperations)
		assert.NotNil(t, diags)
		assert.True(t, diags.HasError())
	})

	// Tests that the Protection Group entity does not allow other backup type operation.
	t.Run("Inhibit different type of assignment of PG", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationEBS
		diags := isOperationsSupported(entityTypeProtectionGroup, policyId, modelOperations)
		assert.NotNil(t, diags)
		assert.True(t, diags.HasError())
	})

	// Tests that the DynamoDB Table does not allow protection group backup type operation.
	t.Run("Inhibit different type of assignment of DDB", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationS3
		diags := isOperationsSupported(entityTypeAWSDynamoDBTable, policyId, modelOperations)
		assert.NotNil(t, diags)
		assert.True(t, diags.HasError())
	})

	// Tests that the DynamoDB Table does not allow other backup type operation.
	t.Run("Inhibit different type of assignment of DDB", func(t *testing.T) {
		modelOperations[0].ClumioType = &operationEBS
		diags := isOperationsSupported(entityTypeAWSDynamoDBTable, policyId, modelOperations)
		assert.NotNil(t, diags)
		assert.True(t, diags.HasError())
	})
}
