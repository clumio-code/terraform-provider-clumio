// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_aws_connection

import (
	"context"
	"fmt"
	"testing"
	"time"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Success scenario for getting the AWS environment for the connection.
//   - SDK API for listing AWS environments returns an error.
//   - SDK API for listing AWS environments returns no AWS environment.
//   - SDK API for listing AWS environments returns more than one AWS environment.
func TestGetEnvironmentForConnection(t *testing.T) {

	mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
	ctx := context.Background()
	cr := &clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkEnvironments: mockAwsEnvClient,
	}

	// Populate the Clumio AWS connection resource model to be used as input for
	// getEnvironmentForConnection().
	state := &clumioAWSConnectionResourceModel{
		ID:              basetypes.NewStringValue(id),
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Tests that the get AWS environment is successful.
	t.Run("success scenario for get AWS environment", func(t *testing.T) {

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{
					{
						Id: &envId,
					},
				},
			},
		}

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(listEnvsResponse, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.Nil(t, err)
		assert.Equal(t, envId, *env.Id)
	})

	// Tests that the get AWS environment fails due to ListAwsEnvironments API returning an error.
	t.Run("list aws environments returns an error", func(t *testing.T) {

		apiError := &apiutils.APIError{
			ResponseCode: 500,
			Reason:       "test",
			Response:     []byte(testError),
		}

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

	// Tests that the get AWS environment fails due to ListAwsEnvironments API returning no
	// environment.
	t.Run("list aws environments returns no environment", func(t *testing.T) {

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(&models.ListAWSEnvironmentsResponse{}, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

	t.Run("falls back to global OU context when current context cannot see environment", func(t *testing.T) {
		previousFactory := newAWSEnvironmentClient
		previousClient := cr.client
		defer func() {
			newAWSEnvironmentClient = previousFactory
			cr.client = previousClient
		}()

		mockGlobalAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		newAWSEnvironmentClient = func(config sdkconfig.Config) sdkclients.AWSEnvironmentClient {
			assert.Equal(t, "", config.OrganizationalUnitContext)
			return mockGlobalAwsEnvClient
		}

		cr.client = &common.ApiClient{
			ClumioConfig: sdkconfig.Config{
				OrganizationalUnitContext: ouId,
			},
		}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(&models.ListAWSEnvironmentsResponse{}, nil)
		mockGlobalAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(&models.ListAWSEnvironmentsResponse{
				Embedded: &models.AWSEnvironmentListEmbedded{
					Items: []*models.AWSEnvironment{{Id: &envId}},
				},
			}, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.Nil(t, err)
		assert.Equal(t, envId, *env.Id)
	})

	// Tests that the get AWS environment fails due to ListAwsEnvironments API returning more than
	// one environment.
	t.Run("list aws environments returns more than one environment", func(t *testing.T) {

		currentCount := int64(2)
		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{
					{
						Id: &envId,
					},
					{
						Id: &envId,
					},
				},
			},
			CurrentCount: &currentCount,
		}

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(listEnvsResponse, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

	t.Run("list aws environments returns environment without id", func(t *testing.T) {
		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{}},
			},
		}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(listEnvsResponse, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

}

func TestGetDesiredOrganizationalUnitID(t *testing.T) {
	t.Run("returns global OU when client is nil", func(t *testing.T) {
		assert.Equal(t, defaultOrgUnitId, getDesiredOrganizationalUnitID(nil))
	})

	t.Run("returns global OU when provider context is empty", func(t *testing.T) {
		client := &common.ApiClient{ClumioConfig: sdkconfig.Config{}}
		assert.Equal(t, defaultOrgUnitId, getDesiredOrganizationalUnitID(client))
	})

	t.Run("returns provider context when set", func(t *testing.T) {
		client := &common.ApiClient{ClumioConfig: sdkconfig.Config{
			OrganizationalUnitContext: ouId,
		}}
		assert.Equal(t, ouId, getDesiredOrganizationalUnitID(client))
	})
}

func TestSetOrganizationalUnitID(t *testing.T) {
	t.Run("uses response OU when present", func(t *testing.T) {
		state := &clumioAWSConnectionResourceModel{}
		setOrganizationalUnitID(state, &ouId)
		assert.Equal(t, ouId, state.OrganizationalUnitID.ValueString())
	})

	t.Run("normalizes empty OU to global", func(t *testing.T) {
		state := &clumioAWSConnectionResourceModel{}
		empty := ""
		setOrganizationalUnitID(state, &empty)
		assert.Equal(t, defaultOrgUnitId, state.OrganizationalUnitID.ValueString())
	})

	t.Run("normalizes nil OU to global", func(t *testing.T) {
		state := &clumioAWSConnectionResourceModel{}
		setOrganizationalUnitID(state, nil)
		assert.Equal(t, defaultOrgUnitId, state.OrganizationalUnitID.ValueString())
	})
}

func TestGetOrgUnitForConnection(t *testing.T) {
	mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)

	t.Run("success scenario for get OU", func(t *testing.T) {
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)

		orgUnit, err := getOrgUnitForConnection(mockOrgUnitsClient, ouId)
		assert.Nil(t, err)
		assert.Equal(t, ouId, *orgUnit.Id)
	})

	t.Run("read OU returns error", func(t *testing.T) {
		apiError := &apiutils.APIError{
			ResponseCode: 500,
			Reason:       "test",
			Response:     []byte(testError),
		}

		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			nil, apiError)

		orgUnit, err := getOrgUnitForConnection(mockOrgUnitsClient, parentOuId)
		assert.NotNil(t, err)
		assert.Nil(t, orgUnit)
	})
}

func TestUpdateOrgUnitForConnection(t *testing.T) {
	ctx := context.Background()

	t.Run("moves connection from global OU to child OU by add", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(defaultOrgUnitId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(parentOuId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Add, 1)
				assert.Nil(t, body.Entities.Remove)
				assert.Equal(t, envId, *body.Entities.Add[0].PrimaryEntity.Id)
				assert.Equal(t, awsEnvironment, *body.Entities.Add[0].PrimaryEntity.ClumioType)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("moves connection from non-global OU to another OU by add", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(parentOuId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Add, 1)
				assert.Nil(t, body.Entities.Remove)
				assert.Equal(t, envId, *body.Entities.Add[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("uses global OU client for OU operations when provider context is set", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockGlobalOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)
		mockGlobalTaskClient := sdkclients.NewMockTaskClient(t)
		previousFactory := newOrganizationalUnitClient
		previousTaskFactory := newTaskClient
		defer func() {
			newOrganizationalUnitClient = previousFactory
			newTaskClient = previousTaskFactory
		}()

		newOrganizationalUnitClient = func(config sdkconfig.Config) sdkclients.OrganizationalUnitClient {
			assert.Equal(t, "", config.OrganizationalUnitContext)
			return mockGlobalOrgUnitsClient
		}
		newTaskClient = func(config sdkconfig.Config) sdkclients.TaskClient {
			assert.Equal(t, "", config.OrganizationalUnitContext)
			return mockGlobalTaskClient
		}

		cr := &clumioAWSConnectionResource{
			client: &common.ApiClient{ClumioConfig: sdkconfig.Config{
				OrganizationalUnitContext: ouId,
			}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockGlobalOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)
		mockGlobalOrgUnitsClient.EXPECT().PatchOrganizationalUnit(parentOuId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Add, 1)
			}).Return(patchResponse, nil).Once()
		mockGlobalTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("moves connection to parent OU by remove", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId, ParentId: &parentOuId}, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(ouId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Remove, 1)
				assert.Nil(t, body.Entities.Add)
				assert.Equal(t, envId, *body.Entities.Remove[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("moves connection to global OU by remove from current OU", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(defaultOrgUnitId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(ouId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Remove, 1)
				assert.Nil(t, body.Entities.Add)
				assert.Equal(t, envId, *body.Entities.Remove[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})
}

// Unit test for setExternalId that checks and sets the ExternalID.
func TestSetExternalId(t *testing.T) {

	// Populate the Clumio AWS connection resource model to be used as input for
	// setExternalId().
	state := &clumioAWSConnectionResourceModel{
		ExternalID: basetypes.NewStringValue(externalId),
	}

	setExternalId(state, nil, &token)
	assert.Equal(t, fmt.Sprintf(externalIDFmt, token), state.ExternalID.ValueString())
}

// Unit test for setDataPlaneAccountId that checks and sets the DataPlaneAccountID.
func TestSetDataPlaneAccountId(t *testing.T) {

	// Populate the Clumio AWS connection resource model to be used as input for
	// setExternalId().
	state := &clumioAWSConnectionResourceModel{
		DataPlaneAccountID: basetypes.NewStringValue(dataplaneAccountId),
	}

	setDataPlaneAccountId(state, nil)
	assert.Equal(t, defaultDataPlaneAccountId, state.DataPlaneAccountID.ValueString())
}
