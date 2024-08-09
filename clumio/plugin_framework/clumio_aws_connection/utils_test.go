// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_aws_connection

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

// Unit test for the following cases:
//   - Update OU for the connection success scenario.
//   - Update OU to current parent OU success scenario.
//   - SDK API for listing AWS environments returns an error.
//   - SDK API for reading organizational unit returns an error.
//   - SDK API for patching organizational unit returns an error.
//   - SDK API for patching organizational unit returns HTTP bad request.
//   - SDK API for patching organizational unit returns an empty response.
//   - SDK API for polling task after OU update returns an error.
func TestUpdateOrgUnitForConnection(t *testing.T) {

	mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
	mockOrgUnitsCient := sdkclients.NewMockOrganizationalUnitClient(t)
	mockTaskClient := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	cr := &clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkEnvironments: mockAwsEnvClient,
		sdkOrgUnits:     mockOrgUnitsCient,
		sdkTasks:        mockTaskClient,
		pollTimeout:     5 * time.Second,
		pollInterval:    1,
	}

	// Populate the Clumio AWS connection resource model to be used as plan in
	// updateOrgUnitForConnection().
	plan := &clumioAWSConnectionResourceModel{
		ID:                   basetypes.NewStringValue(id),
		AccountNativeID:      basetypes.NewStringValue(accountId),
		AWSRegion:            basetypes.NewStringValue(region),
		Description:          basetypes.NewStringValue(description),
		OrganizationalUnitID: basetypes.NewStringValue(ou),
	}

	// Populate the Clumio AWS connection resource model to be used as state in
	// updateOrgUnitForConnection().
	state := &clumioAWSConnectionResourceModel{
		ID:                   basetypes.NewStringValue(id),
		AccountNativeID:      basetypes.NewStringValue(accountId),
		AWSRegion:            basetypes.NewStringValue(region),
		Description:          basetypes.NewStringValue(description),
		OrganizationalUnitID: basetypes.NewStringValue(ou),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests that the OU update is successful.
	t.Run("Success scenario for ou update", func(t *testing.T) {

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &ou,
			}, nil)
		mockOrgUnitsCient.EXPECT().PatchOrganizationalUnit(ou, mock.Anything, mock.Anything).
			Times(1).Return(&models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}, nil)
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(&models.ReadTaskResponse{
			Status: &status,
		}, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	// Tests that the OU update to current parent OU is successful.
	t.Run("Success scenario for ou update to current parent ou", func(t *testing.T) {

		updateOU := "updated-ou"
		planUpdateOU := &clumioAWSConnectionResourceModel{
			ID:                   basetypes.NewStringValue(id),
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			Description:          basetypes.NewStringValue(description),
			OrganizationalUnitID: basetypes.NewStringValue(updateOU),
		}

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id:       &ou,
				ParentId: &updateOU,
			}, nil)
		mockOrgUnitsCient.EXPECT().PatchOrganizationalUnit(ou, mock.Anything, mock.Anything).
			Times(1).Return(&models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}, nil)
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(&models.ReadTaskResponse{
			Status: &status,
		}, nil)

		err := updateOrgUnitForConnection(ctx, cr, planUpdateOU, state)
		assert.Nil(t, err)
	})

	// Tests that the OU update fails due to ListAwsEnvironments API returning an error.
	t.Run("list aws environments returns an error", func(t *testing.T) {

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.NotNil(t, err)
	})

	// Tests that the OU update fails due to ReadOrganizationalUnit API returning an error.
	t.Run("Read OU returns an error", func(t *testing.T) {

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			nil, apiError)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.NotNil(t, err)
	})

	// Tests that the OU update fails due to PatchOrganizationalUnit API returning an error.
	t.Run("Patch OU returns an error due to API Error", func(t *testing.T) {

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &ou,
			}, nil)
		mockOrgUnitsCient.EXPECT().PatchOrganizationalUnit(ou, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.NotNil(t, err)
	})

	// Tests that the OU update fails due to PatchOrganizationalUnit API returning an HTTP Bad
	// Request.
	t.Run("Patch OU returns an error due to HTTP bad request", func(t *testing.T) {

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &ou,
			}, nil)
		mockOrgUnitsCient.EXPECT().PatchOrganizationalUnit(ou, mock.Anything, mock.Anything).
			Times(1).Return(&models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 400,
		}, apiError)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.NotNil(t, err)
	})

	// Tests that the OU update fails due to PatchOrganizationalUnit API returning an empty
	// response.
	t.Run("Patch OU returns an error due to an empty response", func(t *testing.T) {

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &ou,
			}, nil)
		mockOrgUnitsCient.EXPECT().PatchOrganizationalUnit(ou, mock.Anything, mock.Anything).
			Times(1).Return(&models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202:    nil,
		}, apiError)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.NotNil(t, err)
	})

	// Tests that the OU update fails due to ReadTask API returning an error.
	t.Run("Read task returns an error", func(t *testing.T) {

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
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &ou,
			}, nil)
		mockOrgUnitsCient.EXPECT().PatchOrganizationalUnit(ou, mock.Anything, mock.Anything).
			Times(1).Return(&models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}, nil)
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.NotNil(t, err)
	})

}

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
		ID:                   basetypes.NewStringValue(id),
		AccountNativeID:      basetypes.NewStringValue(accountId),
		AWSRegion:            basetypes.NewStringValue(region),
		Description:          basetypes.NewStringValue(description),
		OrganizationalUnitID: basetypes.NewStringValue(ou),
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

}

// Unit test for the following cases:
//   - Get organizational unit for connection success scenario.
//   - SDK API for read OU returns error.
func TestGetOrgUnitForConnection(t *testing.T) {

	mockOrgUnitsCient := sdkclients.NewMockOrganizationalUnitClient(t)
	ctx := context.Background()
	cr := &clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkOrgUnits: mockOrgUnitsCient,
	}

	// Populate the Clumio AWS connection resource model to be used as input for
	// getOrgUnitForConnection()
	model := &clumioAWSConnectionResourceModel{
		ID:                   basetypes.NewStringValue(id),
		AccountNativeID:      basetypes.NewStringValue(accountId),
		AWSRegion:            basetypes.NewStringValue(region),
		Description:          basetypes.NewStringValue(description),
		OrganizationalUnitID: basetypes.NewStringValue(ou),
	}

	// Tests that the get OU for the AWS connection is successful.
	t.Run("success scenario for get OU for connection", func(t *testing.T) {

		readOUResponse := &models.ReadOrganizationalUnitResponse{
			Id: &id,
		}

		// Setup Expectations
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).
			Return(readOUResponse, nil)

		ouResp, err := getOrgUnitForConnection(ctx, cr, model)
		assert.Nil(t, err)
		assert.Equal(t, id, *ouResp.Id)
	})

	// Tests that the get OU for the AWS connection fails due to ReadOrganizationalUnit API
	// returning an error.
	t.Run("read OU returns an error", func(t *testing.T) {

		apiError := &apiutils.APIError{
			ResponseCode: 500,
			Reason:       "test",
			Response:     []byte(testError),
		}

		// Setup Expectations
		mockOrgUnitsCient.EXPECT().ReadOrganizationalUnit(ou, mock.Anything).Times(1).
			Return(nil, apiError)

		ouResp, err := getOrgUnitForConnection(ctx, cr, model)
		assert.NotNil(t, err)
		assert.Nil(t, ouResp)
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
