// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy

import (
	"context"
	"testing"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	name         = "test-policy"
	resourceName = "test_policy"
	id           = "mock-policy-id"
	ou           = "mock-ou"
	testError    = "Test Error"
)

// Unit test for the following cases:
//   - Create policy success scenario.
//   - SDK API for create policy returns error.
//   - SDK API for create policy returns nil response.
//   - SDK API for read policy returns error.
func TestCreatePolicy(t *testing.T) {

	mockPolicy := sdkclients.NewMockPolicyDefinitionClient(t)
	pr := policyResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicy,
	}

	timezone := "UTC"
	activationStatus := activationStatusActivated
	actionSetting := "immediate"
	operationType := "aws_ebs_volume_backup"
	retUnit := "days"
	retValue := int64(5)
	rpoUnit := "days"
	rpoValue := int64(1)

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the policy resource model to be used as input to createPolicy()
	prm := policyResourceModel{
		Name:             basetypes.NewStringValue(name),
		Timezone:         basetypes.NewStringValue(timezone),
		ActivationStatus: basetypes.NewStringValue(activationStatus),
		Operations: []*policyOperationModel{
			{
				ActionSetting: basetypes.NewStringValue(actionSetting),
				OperationType: basetypes.NewStringValue(operationType),
				Slas: []*slaModel{
					{
						RetentionDuration: []*unitValueModel{
							{
								Unit:  basetypes.NewStringValue(retUnit),
								Value: basetypes.NewInt64Value(retValue),
							},
						},
						RPOFrequency: []*rpoModel{
							{
								Unit:    basetypes.NewStringValue(rpoUnit),
								Value:   basetypes.NewInt64Value(rpoValue),
								Offsets: basetypes.NewListNull(types.Int64Type),
							},
						},
					},
				},
			},
		},
	}

	// Tests the success scenario for policy create. It should not return Diagnostics.
	t.Run("Basic success scenario for create policy", func(t *testing.T) {

		// Create the response of the SDK CreatePolicyDefinition() API.
		lockStatus := "unlocked"
		createResponse := &models.CreatePolicyResponse{
			ActivationStatus: &activationStatus,
			Id:               &id,
			LockStatus:       &lockStatus,
			Name:             &name,
			Operations: []*models.PolicyOperation{
				{
					ActionSetting: &actionSetting,
					Slas: []*models.BackupSLA{
						{
							RetentionDuration: &models.RetentionBackupSLAParam{
								Unit:  &retUnit,
								Value: &retValue,
							},
							RpoFrequency: &models.RPOBackupSLAParam{
								Unit:  &rpoUnit,
								Value: &rpoValue,
							},
						},
					},
				},
			},
			OrganizationalUnitId: &ou,
			Timezone:             &timezone,
		}

		// Create the response of the SDK ReadPolicyDefinition() API.
		readResponse := &models.ReadPolicyResponse{
			ActivationStatus: &activationStatus,
			Id:               &id,
			LockStatus:       &lockStatus,
			Name:             &name,
			Operations: []*models.PolicyOperation{
				{
					ActionSetting: &actionSetting,
					Slas: []*models.BackupSLA{
						{
							RetentionDuration: &models.RetentionBackupSLAParam{
								Unit:  &retUnit,
								Value: &retValue,
							},
							RpoFrequency: &models.RPOBackupSLAParam{
								Unit:  &rpoUnit,
								Value: &rpoValue,
							},
						},
					},
				},
			},
			OrganizationalUnitId: &ou,
			Timezone:             &timezone,
		}

		// Setup Expectations
		mockPolicy.EXPECT().CreatePolicyDefinition(mock.Anything).Times(1).
			Return(createResponse, nil)
		mockPolicy.EXPECT().ReadPolicyDefinition(mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := pr.createPolicy(context.Background(), &prm)
		assert.Nil(t, diags)
		assert.Equal(t, id, prm.ID.ValueString())
		assert.Equal(t, lockStatus, prm.LockStatus.ValueString())
		assert.Equal(t, ou, prm.OrganizationalUnitId.ValueString())
	})

	// Tests that Diagnostics is returned in case the create policy API call returns error.
	t.Run("CreatePolicyDefinition returns error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().CreatePolicyDefinition(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.createPolicy(context.Background(), &prm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create policy API call returns an empty
	// response.
	t.Run("CreatePolicyDefinition returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().CreatePolicyDefinition(mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.createPolicy(context.Background(), &prm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned if read policy fails which is done as part of the
	// createPolicy function.
	t.Run("ReadPolicyDefinition returns error", func(t *testing.T) {

		lockStatus := "unlocked"
		createResponse := &models.CreatePolicyResponse{
			ActivationStatus: &activationStatus,
			Id:               &id,
			LockStatus:       &lockStatus,
			Name:             &name,
			Operations: []*models.PolicyOperation{
				{
					ActionSetting: &actionSetting,
					Slas: []*models.BackupSLA{
						{
							RetentionDuration: &models.RetentionBackupSLAParam{
								Unit:  &retUnit,
								Value: &retValue,
							},
							RpoFrequency: &models.RPOBackupSLAParam{
								Unit:  &rpoUnit,
								Value: &rpoValue,
							},
						},
					},
				},
			},
			OrganizationalUnitId: &ou,
			Timezone:             &timezone,
		}

		// Setup Expectations
		mockPolicy.EXPECT().CreatePolicyDefinition(mock.Anything).Times(1).
			Return(createResponse, nil)
		mockPolicy.EXPECT().ReadPolicyDefinition(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.createPolicy(context.Background(), &prm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

}

// Unit test for the following cases:
//   - Read policy success scenario.
//   - SDK API for read policy returns not found error.
//   - SDK API for read policy returns error.
//   - SDK API for create policy returns nil response.
func TestReadPolicy(t *testing.T) {

	mockPolicy := sdkclients.NewMockPolicyDefinitionClient(t)
	pr := policyResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicy,
	}

	timezone := "UTC"
	activationStatus := activationStatusActivated
	actionSetting := "immediate"
	operationType := "aws_ebs_volume_backup"
	retUnit := "days"
	retValue := int64(5)
	rpoUnit := "days"
	rpoValue := int64(1)
	lockStatus := "unlocked"

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	prm := policyResourceModel{
		ID: basetypes.NewStringValue(id),
	}
	// Create the response of the SDK ReadPolicyDefinition() API.
	readResponse := &models.ReadPolicyResponse{
		ActivationStatus: &activationStatus,
		Id:               &id,
		LockStatus:       &lockStatus,
		Name:             &name,
		Operations: []*models.PolicyOperation{
			{
				ClumioType:    &operationType,
				ActionSetting: &actionSetting,
				Slas: []*models.BackupSLA{
					{
						RetentionDuration: &models.RetentionBackupSLAParam{
							Unit:  &retUnit,
							Value: &retValue,
						},
						RpoFrequency: &models.RPOBackupSLAParam{
							Unit:  &rpoUnit,
							Value: &rpoValue,
						},
					},
				},
			},
		},
		OrganizationalUnitId: &ou,
		Timezone:             &timezone,
	}

	// Tests the success scenario for policy read. It should not return Diagnostics.
	t.Run("success scenario for read policy", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().ReadPolicyDefinition(id, mock.Anything).Times(1).
			Return(readResponse, nil)

		remove, diags := pr.readPolicy(context.Background(), &prm)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that in case the policy is not found, it returns true to indicate that the policy
	// should be removed from the state.
	t.Run("read policy returns not found error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().ReadPolicyDefinition(id, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := pr.readPolicy(context.Background(), &prm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read policy API call returns error.
	t.Run("read policy returns error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().ReadPolicyDefinition(id, mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := pr.readPolicy(context.Background(), &prm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read policy API call returns an empty response.
	t.Run("read policy returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().ReadPolicyDefinition(id, mock.Anything).Times(1).
			Return(nil, nil)

		remove, diags := pr.readPolicy(context.Background(), &prm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

}

// Unit test for the following cases:
//   - Update policy success scenario.
//   - SDK API for update policy returns error.
//   - SDK API for update policy returns nil response.
//   - Polling of delete policy task returns error.
//   - SDK API for read policy returns error.
func TestUpdatePolicy(t *testing.T) {

	mockPolicy := sdkclients.NewMockPolicyDefinitionClient(t)
	mockTask := sdkclients.NewMockTaskClient(t)
	pr := policyResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicy,
		sdkTasks:             mockTask,
		pollTimeout:          5 * time.Second,
		pollInterval:         1,
	}

	timezone := "UTC"
	activationStatus := activationStatusActivated
	actionSetting := "immediate"
	operationType := "aws_ebs_volume_backup"
	retUnit := "days"
	retValue := int64(5)
	rpoUnit := "days"
	rpoValue := int64(1)

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the policy resource model to be used as input to updatePolicy()
	prm := policyResourceModel{
		ID:               basetypes.NewStringValue(id),
		Name:             basetypes.NewStringValue(name),
		Timezone:         basetypes.NewStringValue(timezone),
		ActivationStatus: basetypes.NewStringValue(activationStatus),
		Operations: []*policyOperationModel{
			{
				ActionSetting: basetypes.NewStringValue(actionSetting),
				OperationType: basetypes.NewStringValue(operationType),
				Slas: []*slaModel{
					{
						RetentionDuration: []*unitValueModel{
							{
								Unit:  basetypes.NewStringValue(retUnit),
								Value: basetypes.NewInt64Value(retValue),
							},
						},
						RPOFrequency: []*rpoModel{
							{
								Unit:    basetypes.NewStringValue(rpoUnit),
								Value:   basetypes.NewInt64Value(rpoValue),
								Offsets: basetypes.NewListNull(types.Int64Type),
							},
						},
					},
				},
			},
		},
	}

	// Create the response of the SDK CreatePolicyDefinition() API.
	lockStatus := "unlocked"
	taskId := "12345"
	updateResponse := &models.UpdatePolicyResponse{
		TaskId:           &taskId,
		ActivationStatus: &activationStatus,
		Id:               &id,
		LockStatus:       &lockStatus,
		Name:             &name,
		Operations: []*models.PolicyOperation{
			{
				ActionSetting: &actionSetting,
				Slas: []*models.BackupSLA{
					{
						RetentionDuration: &models.RetentionBackupSLAParam{
							Unit:  &retUnit,
							Value: &retValue,
						},
						RpoFrequency: &models.RPOBackupSLAParam{
							Unit:  &rpoUnit,
							Value: &rpoValue,
						},
					},
				},
			},
		},
		OrganizationalUnitId: &ou,
		Timezone:             &timezone,
	}

	taskStatus := common.TaskSuccess
	readTaskResponse := &models.ReadTaskResponse{
		Status: &taskStatus,
	}

	// Tests the success scenario for policy update. It should not return Diagnostics.
	t.Run("Basic success scenario for update policy", func(t *testing.T) {

		// Create the response of the SDK ReadPolicyDefinition() API.
		readResponse := &models.ReadPolicyResponse{
			ActivationStatus: &activationStatus,
			Id:               &id,
			LockStatus:       &lockStatus,
			Name:             &name,
			Operations: []*models.PolicyOperation{
				{
					ActionSetting: &actionSetting,
					Slas: []*models.BackupSLA{
						{
							RetentionDuration: &models.RetentionBackupSLAParam{
								Unit:  &retUnit,
								Value: &retValue,
							},
							RpoFrequency: &models.RPOBackupSLAParam{
								Unit:  &rpoUnit,
								Value: &rpoValue,
							},
						},
					},
				},
			},
			OrganizationalUnitId: &ou,
			Timezone:             &timezone,
		}

		// Setup Expectations
		mockPolicy.EXPECT().UpdatePolicyDefinition(id, mock.Anything, mock.Anything).Times(1).
			Return(updateResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)
		mockPolicy.EXPECT().ReadPolicyDefinition(mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := pr.updatePolicy(context.Background(), &prm)
		assert.Nil(t, diags)
		assert.Equal(t, id, prm.ID.ValueString())
		assert.Equal(t, lockStatus, prm.LockStatus.ValueString())
		assert.Equal(t, ou, prm.OrganizationalUnitId.ValueString())
	})

	// Tests that Diagnostics is returned if the update policy API call returns error.
	t.Run("UpdatePolicyDefinition returns error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().UpdatePolicyDefinition(id, mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.updatePolicy(context.Background(), &prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned if the update policy API call returns an empty response.
	t.Run("UpdatePolicyDefinition returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().UpdatePolicyDefinition(id, mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.updatePolicy(context.Background(), &prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when polling of the update policy task fails.
	t.Run("Task poll returns error", func(t *testing.T) {

		// Setup Expectations
		mockPolicy.EXPECT().UpdatePolicyDefinition(id, mock.Anything, mock.Anything).Times(1).
			Return(updateResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		diags := pr.updatePolicy(context.Background(), &prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned if read policy fails which is done as part of the
	// updatePolicy function.
	t.Run("ReadPolicyDefinition returns error", func(t *testing.T) {

		// Setup Expectations
		mockPolicy.EXPECT().UpdatePolicyDefinition(id, mock.Anything, mock.Anything).Times(1).
			Return(updateResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)
		mockPolicy.EXPECT().ReadPolicyDefinition(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.updatePolicy(context.Background(), &prm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

}

// Unit test for the following cases:
//   - Delete policy success scenario.
//   - Delete policy should not return error if policy is not found.
//   - SDK API for delete policy returns error.
//   - SDK API for delete policy returns nil response.
//   - Polling of delete policy task returns error.
func TestDeletePolicy(t *testing.T) {

	mockPolicy := sdkclients.NewMockPolicyDefinitionClient(t)
	mockTask := sdkclients.NewMockTaskClient(t)
	pr := policyResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyDefinitions: mockPolicy,
		sdkTasks:             mockTask,
		pollTimeout:          5 * time.Second,
		pollInterval:         1,
	}
	// Populate the policy resource model to be used as input to createPolicy()
	prm := &policyResourceModel{
		ID: basetypes.NewStringValue(id),
	}

	taskId := "12345"
	deleteResponse := &models.DeletePolicyResponse{
		TaskId: &taskId,
	}
	taskStatus := common.TaskSuccess
	readTaskResponse := &models.ReadTaskResponse{
		Status: &taskStatus,
	}
	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}
	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "Test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for policy deletion. It should not return diag.Diagnostics.
	t.Run("Success scenario for policy deletion", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().DeletePolicyDefinition(id).Times(1).Return(deleteResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

		diags := pr.deletePolicy(context.Background(), prm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the policy does not exist.
	t.Run("Policy not found should not return error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().DeletePolicyDefinition(id).Times(1).Return(nil, apiNotFoundError)

		diags := pr.deletePolicy(context.Background(), prm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete policy API call returns error.
	t.Run("deletePolicy returns error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().DeletePolicyDefinition(id).Times(1).Return(nil, apiError)

		diags := pr.deletePolicy(context.Background(), prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when policy deletion returns an empty response.
	t.Run("deletePolicy returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().DeletePolicyDefinition(id).Times(1).Return(nil, nil)

		diags := pr.deletePolicy(context.Background(), prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when polling of the delete policy task fails.
	t.Run("Task poll returns error", func(t *testing.T) {
		// Setup Expectations
		mockPolicy.EXPECT().DeletePolicyDefinition(id).Times(1).Return(deleteResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		diags := pr.deletePolicy(context.Background(), prm)
		assert.NotNil(t, diags)
	})
}
