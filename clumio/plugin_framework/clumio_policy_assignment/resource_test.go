// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	resourceName  = "test_policy_assignment"
	entityId      = "test-pg-id"
	policyId      = "test-policy-id"
	otherPolicyId = "other-policy-id"
	policyType    = protectionGroupBackup
	ou            = "test-ou"
	testError     = "Test Error"
	taskId        = "test-task-id"
)

const (
	// The following constants are used as a test name in different tests.
	readPolicyError                  = "Read policy definition returns an error"
	readPolicyNotFoundError          = "Read policy definition returns not found error"
	setPolicyAssignmentPollingError  = "Polling for set policy assignment task returns an error"
	readProtectionGroupError         = "Read protection group returns an error"
	readProtectionGroupEmptyResponse = "Read protection group returns an empty response"
	readProtectionGroupNotFoundError = "Read protection group returns not found error"
)

// Unit test for the following cases:
//   - Create policy assignment success scenario.
//   - SDK API for read policy definition returns an error.
//   - SDK API for read policy definition returns a policy without required operation type.
//   - SDK API for set policy assignments returns an error.
//   - SDK API for set policy assignments returns an empty response.
//   - SDK API for read task returns an error.
func TestCreatePolicyAssignment(t *testing.T) {

	ctx := context.Background()
	mockPolicyDefinitions := sdkclients.NewMockPolicyDefinitionClient(t)
	mockPolicyAssignments := sdkclients.NewMockPolicyAssignmentClient(t)
	mockProtectionGroups := sdkclients.NewMockProtectionGroupClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	par := &clumioPolicyAssignmentResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicyDefinitions,
		sdkProtectionGroups:  mockProtectionGroups,
		sdkPolicyAssignments: mockPolicyAssignments,
		sdkTasks:             mockTasks,
		pollTimeout:          5 * time.Second,
		pollInterval:         1,
	}

	model := &policyAssignmentResourceModel{
		EntityID:   basetypes.NewStringValue(entityId),
		EntityType: basetypes.NewStringValue(entityTypeProtectionGroup),
		PolicyID:   basetypes.NewStringValue(policyId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for create policy assignment. It should not return Diagnostics.
	t.Run("Basic success scenario for create policy assignment", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup Expectations.
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

		diags := par.createPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests the success scenario for create policy assignment for DynamoDB table. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for create policy assignment for DynamoDB table",
		func(t *testing.T) {

			model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
			policyType := dynamodbTableBackup
			pdResp := &models.ReadPolicyResponse{
				Id: &policyId,
				Operations: []*models.PolicyOperation{
					{
						ClumioType: &policyType,
					},
				},
				OrganizationalUnitId: &ou,
			}

			paResp := &models.SetAssignmentsResponse{
				TaskId: &taskId,
			}

			taskStatus := common.TaskSuccess
			readTaskResponse := &models.ReadTaskResponse{
				Status: &taskStatus,
			}

			// Setup Expectations.
			mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
				Return(pdResp, nil)
			mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
				paResp, nil)
			mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

			diags := par.createPolicyAssignment(ctx, model)
			assert.Nil(t, diags)
			model.EntityType = basetypes.NewStringValue(entityTypeProtectionGroup)
		})

	// Tests that Diagnostics is returned in case the read policy definition API call returns an
	// error.
	t.Run(readPolicyError, func(t *testing.T) {

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := par.createPolicyAssignment(context.Background(), model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read policy definition API call returns a
	// policy without the required operation type for policy assignment.
	t.Run("Read policy definition returns policy with unsupported type", func(t *testing.T) {

		opType := "some-type"
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &opType,
				},
			},
		}

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)

		diags := par.createPolicyAssignment(context.Background(), model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the set policy assignments API call returns an
	// error.
	t.Run("Set policy assignments returns an error", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
		}

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := par.createPolicyAssignment(context.Background(), model)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the set policy assignments API call returns an
	// empty response.
	t.Run("Set policy assignments returns an empty response", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
		}

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			nil, nil)

		diags := par.createPolicyAssignment(context.Background(), model)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read task API call returns an error.
	t.Run(setPolicyAssignmentPollingError, func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
		}

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		diags := par.createPolicyAssignment(context.Background(), model)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read policy assignment success scenario.
//   - Read policy assignment with invalid entity type returns an error.
//   - SDK API for read policy definition returns an error.
//   - SDK API for read policy definition returns not found error.
//   - SDK API for read protection group returns an error.
//   - SDK API for read protection group returns an empty response.
//   - SDK API for read protection group returns not found error.
//   - protection group assigned policy does not match the given policy.
func TestReadPolicyAssignment(t *testing.T) {

	ctx := context.Background()
	mockPolicyDefinitions := sdkclients.NewMockPolicyDefinitionClient(t)
	mockPolicyAssignments := sdkclients.NewMockPolicyAssignmentClient(t)
	mockProtectionGroups := sdkclients.NewMockProtectionGroupClient(t)
	mockDynamoDBTables := sdkclients.NewMockDynamoDBTableClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	par := &clumioPolicyAssignmentResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicyDefinitions,
		sdkProtectionGroups:  mockProtectionGroups,
		sdkPolicyAssignments: mockPolicyAssignments,
		sdkDynamoDBTables:    mockDynamoDBTables,
		sdkTasks:             mockTasks,
	}

	id := fmt.Sprintf("%s_%s_%s", policyId, entityId, entityTypeProtectionGroup)
	model := &policyAssignmentResourceModel{
		ID:         basetypes.NewStringValue(id),
		EntityID:   basetypes.NewStringValue(entityId),
		EntityType: basetypes.NewStringValue(entityTypeProtectionGroup),
		PolicyID:   basetypes.NewStringValue(policyId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Tests the success scenario for read policy assignment. It should not return Diagnostics.
	t.Run("Basic success scenario for read policy assignment", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		readPgResp := &models.ReadProtectionGroupResponse{
			Id: &entityId,
			ProtectionInfo: &models.ProtectionInfoWithRule{
				PolicyId: &policyId,
			},
			OrganizationalUnitId: &ou,
		}

		// Setup Expectations.
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockProtectionGroups.EXPECT().ReadProtectionGroup(entityId, mock.Anything).Times(1).Return(
			readPgResp, nil)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests the success scenario for read policy assignment for DynamoDB Table. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for read policy assignment for DynamoDB table",
		func(t *testing.T) {

			model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
			policyType := dynamodbTableBackup
			pdResp := &models.ReadPolicyResponse{
				Id: &policyId,
				Operations: []*models.PolicyOperation{
					{
						ClumioType: &policyType,
					},
				},
				OrganizationalUnitId: &ou,
			}
			readTableResp := &models.ReadDynamoDBTableResponse{
				Id: &entityId,
				ProtectionInfo: &models.ProtectionInfoWithRule{
					PolicyId: &policyId,
				},
				OrganizationalUnitId: &ou,
			}

			// Setup Expectations.
			mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
				Return(pdResp, nil)
			mockDynamoDBTables.EXPECT().ReadAwsDynamodbTable(
				entityId, mock.Anything, mock.Anything).Times(1).Return(readTableResp, nil)

			remove, diags := par.readPolicyAssignment(ctx, model)
			assert.Nil(t, diags)
			assert.False(t, remove)
			model.EntityType = basetypes.NewStringValue(entityTypeProtectionGroup)
		})

	// Tests that Diagnostics is returned in case the read policy assignment with invalid entity
	// type.
	t.Run("Read policy assignment with invalid entity type", func(t *testing.T) {

		invalidType := "invalid-type"
		modelWithInvalidType := &policyAssignmentResourceModel{
			ID:         basetypes.NewStringValue(id),
			EntityID:   basetypes.NewStringValue(entityId),
			EntityType: basetypes.NewStringValue(invalidType),
			PolicyID:   basetypes.NewStringValue(policyId),
		}

		// Setup Expectations
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)

		remove, diags := par.readPolicyAssignment(ctx, modelWithInvalidType)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read policy assignment with invalid entity
	// type.
	t.Run("Read policy with invalid policy operation type", func(t *testing.T) {

		dynamodbPolicyType := dynamodbTableBackup
		// Setup Expectations
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &dynamodbPolicyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read policy definition API call returns an
	// error.
	t.Run(readPolicyError, func(t *testing.T) {

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that read policy assignment returns true to indicate that the resource should be
	// removed when read policy definition API call returns not found error.
	t.Run(readPolicyNotFoundError, func(t *testing.T) {

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// error.
	t.Run(readProtectionGroupError, func(t *testing.T) {

		// Setup Expectations
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockProtectionGroups.EXPECT().ReadProtectionGroup(entityId, mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// empty response.
	t.Run(readProtectionGroupEmptyResponse, func(t *testing.T) {

		// Setup Expectations
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockProtectionGroups.EXPECT().ReadProtectionGroup(entityId, mock.Anything).Times(1).
			Return(nil, nil)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that read policy assignment returns true to indicate that the resource should be
	// removed when read protection group API call returns not found error.
	t.Run(readProtectionGroupNotFoundError, func(t *testing.T) {

		// Setup Expectations
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockProtectionGroups.EXPECT().ReadProtectionGroup(entityId, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that read policy assignment returns true to indicate that the resource should be
	// removed when read protection group API call returns a protection group which has a policy
	// that is different to the one in the state.
	t.Run("Protection group assigned policy does not match given policy", func(t *testing.T) {

		readPgResp := &models.ReadProtectionGroupResponse{
			Id: &entityId,
			ProtectionInfo: &models.ProtectionInfoWithRule{
				PolicyId: &otherPolicyId,
			},
			OrganizationalUnitId: &ou,
		}
		// Setup Expectations
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockProtectionGroups.EXPECT().ReadProtectionGroup(entityId, mock.Anything).Times(1).
			Return(readPgResp, nil)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read DynamoDB table API call returns an
	// error.
	t.Run("Read DynamoDB table returns an error", func(t *testing.T) {

		model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
		// Setup Expectations
		policyType := dynamodbTableBackup
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockDynamoDBTables.EXPECT().ReadAwsDynamodbTable(entityId, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read DynamoDB table API call returns an empty
	// response.
	t.Run("Read DynamoDB table returns an empty response", func(t *testing.T) {

		model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
		// Setup Expectations
		policyType := dynamodbTableBackup
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockDynamoDBTables.EXPECT().ReadAwsDynamodbTable(entityId, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that read policy assignment returns true to indicate that the resource should be
	// removed when read protection group API call returns not found error.
	t.Run("Read DynamoDB table returns not found error", func(t *testing.T) {

		model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
		// Setup Expectations
		policyType := dynamodbTableBackup
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockDynamoDBTables.EXPECT().ReadAwsDynamodbTable(entityId, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiNotFoundError)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that read policy assignment returns true to indicate that the resource should be
	// removed when read DynamoDB table API call returns a protection group which has a policy
	// that is different to the one in the state.
	t.Run("DynamoDB table assigned policy does not match given policy", func(t *testing.T) {

		model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
		readTableResp := &models.ReadDynamoDBTableResponse{
			Id: &entityId,
			ProtectionInfo: &models.ProtectionInfoWithRule{
				PolicyId: &otherPolicyId,
			},
			OrganizationalUnitId: &ou,
		}
		// Setup Expectations
		policyType := dynamodbTableBackup
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockDynamoDBTables.EXPECT().ReadAwsDynamodbTable(entityId, mock.Anything, mock.Anything).
			Times(1).Return(readTableResp, nil)

		remove, diags := par.readPolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})
}

// Unit test for the following cases:
//   - Update policy assignment success scenario.
//   - SDK API for read policy definition returns an error.
//   - SDK API for read policy definition returns a policy without required operation type.
//   - SDK API for set policy assignments returns an error.
//   - SDK API for set policy assignments returns an empty response.
//   - SDK API for read task returns an error.
func TestUpdatePolicyAssignment(t *testing.T) {

	ctx := context.Background()
	mockPolicyDefinitions := sdkclients.NewMockPolicyDefinitionClient(t)
	mockPolicyAssignments := sdkclients.NewMockPolicyAssignmentClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	par := &clumioPolicyAssignmentResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicyDefinitions,
		sdkPolicyAssignments: mockPolicyAssignments,
		sdkTasks:             mockTasks,
		pollTimeout:          5 * time.Second,
		pollInterval:         1,
	}

	model := &policyAssignmentResourceModel{
		EntityID:   basetypes.NewStringValue(entityId),
		EntityType: basetypes.NewStringValue(entityTypeProtectionGroup),
		PolicyID:   basetypes.NewStringValue(policyId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for update policy assignment. It should not return Diagnostics.
	t.Run("Basic success scenario for update policy assignment", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup Expectations.
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

		diags := par.updatePolicyAssignment(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests the success scenario for update policy assignment. It should not return Diagnostics.
	t.Run("Basic success scenario for update policy assignment", func(t *testing.T) {

		model.EntityType = basetypes.NewStringValue(entityTypeAWSDynamoDBTable)
		policyType := dynamodbTableBackup
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
			OrganizationalUnitId: &ou,
		}

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup Expectations.
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

		diags := par.updatePolicyAssignment(ctx, model)
		assert.Nil(t, diags)
		model.EntityType = basetypes.NewStringValue(entityTypeProtectionGroup)
	})

	// Tests that Diagnostics is returned in case the read policy definition API call returns an
	// error.
	t.Run(readPolicyError, func(t *testing.T) {

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := par.updatePolicyAssignment(context.Background(), model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read policy definition API call returns a
	// policy without the required operation type for policy assignment.
	t.Run("Read policy definition returns policy with unsupported type", func(t *testing.T) {

		opType := "some-type"
		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &opType,
				},
			},
		}

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)

		diags := par.updatePolicyAssignment(context.Background(), model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the set policy assignments API call returns an
	// error.
	t.Run("Set policy assignments returns an error", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
		}

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := par.updatePolicyAssignment(context.Background(), model)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the set policy assignments API call returns an
	// empty response.
	t.Run("Set policy assignments returns an empty response", func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
		}

		// Setup Expectations
		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			nil, nil)

		diags := par.updatePolicyAssignment(context.Background(), model)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read task API call returns an error.
	t.Run(setPolicyAssignmentPollingError, func(t *testing.T) {

		pdResp := &models.ReadPolicyResponse{
			Id: &policyId,
			Operations: []*models.PolicyOperation{
				{
					ClumioType: &policyType,
				},
			},
		}

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		mockPolicyDefinitions.EXPECT().ReadPolicyDefinition(policyId, mock.Anything).Times(1).
			Return(pdResp, nil)
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		diags := par.updatePolicyAssignment(context.Background(), model)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete policy assignment success scenario.
//   - SDK API for set policy assignments returns an error.
//   - SDK API for set policy assignments returns not found error.
//   - SDK API for read task returns not found error.
func TestDeletePolicyAssignment(t *testing.T) {

	ctx := context.Background()
	mockPolicyAssignments := sdkclients.NewMockPolicyAssignmentClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	par := &clumioPolicyAssignmentResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyAssignments: mockPolicyAssignments,
		sdkTasks:             mockTasks,
		pollTimeout:          5 * time.Second,
		pollInterval:         1,
	}

	id := fmt.Sprintf("%s_%s_%s", policyId, entityId, entityTypeProtectionGroup)
	model := &policyAssignmentResourceModel{
		ID:         basetypes.NewStringValue(id),
		EntityID:   basetypes.NewStringValue(entityId),
		EntityType: basetypes.NewStringValue(entityTypeProtectionGroup),
		PolicyID:   basetypes.NewStringValue(policyId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for delete policy assignment. It should not return Diagnostics.
	t.Run("Basic success scenario for delete policy assignment", func(t *testing.T) {

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup Expectations
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

		diags := par.deletePolicyAssignment(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when set policy assignments API call returns an error.
	t.Run("set policy assignments returns an error", func(t *testing.T) {
		// Setup Expectations
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := par.deletePolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when set policy assignments API call returns an empty
	// response.
	t.Run("set policy assignments returns an empty response", func(t *testing.T) {
		// Setup Expectations
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := par.deletePolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read task API call returns an error.
	t.Run(setPolicyAssignmentPollingError, func(t *testing.T) {

		paResp := &models.SetAssignmentsResponse{
			TaskId: &taskId,
		}

		// Setup Expectations
		mockPolicyAssignments.EXPECT().SetPolicyAssignments(mock.Anything).Times(1).Return(
			paResp, nil)
		mockTasks.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		diags := par.deletePolicyAssignment(ctx, model)
		assert.NotNil(t, diags)
	})
}
