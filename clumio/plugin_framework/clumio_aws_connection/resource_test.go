// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_aws_connection

import (
	"context"
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
	accountId          = "test-aws-account"
	region             = "test-region"
	clumioAccountId    = "test-clumio-account-id"
	description        = "test-description"
	resourceName       = "test_aws_connection"
	id                 = "mock-connection-id"
	connStatus         = "test-status"
	token              = "test-token"
	externalId         = "test-external-id"
	namespace          = "test-namespace"
	dataplaneAccountId = "test-dataplane-account-id"
	envId              = "test-env-id"
	taskId             = "test-task-id"
	status             = common.TaskSuccess

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Create AWS connection success scenario.
//   - SDK API for read OU returns error.
//   - SDK API for create AWS connection returns error.
//   - SDK API for create AWS connection returns nil response.
func TestCreateAWSConnection(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockAWSConnectionClient(t)
	mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
	mockOrgUnitsCient := sdkclients.NewMockOrganizationalUnitClient(t)
	ctx := context.Background()
	cr := clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections:  mockAwsConnClient,
		sdkEnvironments: mockAwsEnvClient,
		sdkOrgUnits:     mockOrgUnitsCient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the Clumio AWS connection resource model to be used as input to createAWSConnection()
	crm := clumioAWSConnectionResourceModel{
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Tests the success scenario for clumio aws connection create. It should not return Diagnostics.
	t.Run("Basic success scenario for create aws connection", func(t *testing.T) {

		createResponse := &models.CreateAWSConnectionResponse{
			AccountNativeId:    &accountId,
			AwsRegion:          &region,
			ClumioAwsAccountId: &clumioAccountId,
			ClumioAwsRegion:    &region,
			ConnectionStatus:   &connStatus,
			DataPlaneAccountId: &dataplaneAccountId,
			Description:        &description,
			ExternalId:         &externalId,
			Id:                 &id,
			Namespace:          &namespace,
			Token:              &token,
		}

		// Setup Expectations
		mockAwsConnClient.EXPECT().CreateAwsConnection(mock.Anything).Times(1).Return(
			createResponse, nil)

		diags := cr.createAWSConnection(ctx, &crm)
		assert.Nil(t, diags)
		assert.Equal(t, *createResponse.Id, crm.ID.ValueString())
		assert.Equal(t, *createResponse.ClumioAwsAccountId, crm.ClumioAWSAccountID.ValueString())
		assert.Equal(t, *createResponse.ClumioAwsRegion, crm.ClumioAWSRegion.ValueString())
		assert.Equal(t, *createResponse.ConnectionStatus, crm.ConnectionStatus.ValueString())
		assert.Equal(t, *createResponse.DataPlaneAccountId, crm.DataPlaneAccountID.ValueString())
		assert.Equal(t, *createResponse.ExternalId, crm.ExternalID.ValueString())
		assert.Equal(t, *createResponse.Token, crm.Token.ValueString())

	})

	// Tests that Diagnostics is returned in case the create aws connection API call returns an
	// error.
	t.Run("CreateAwsConnection returns an error", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().CreateAwsConnection(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := cr.createAWSConnection(ctx, &crm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create aws connection API call returns an
	// empty response.
	t.Run("CreateAwsConnection returns an empty response", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().CreateAwsConnection(mock.Anything).Times(1).Return(
			nil, nil)

		diags := cr.createAWSConnection(ctx, &crm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read AWS connection success scenario.
//   - SDK API for read AWS connection returns not found error.
//   - SDK API for read Clumio AWS connection returns error.
//   - SDK API for create Clumio AWS connection returns nil response.
func TestReadAWSConnection(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockAWSConnectionClient(t)
	ctx := context.Background()
	cr := clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockAwsConnClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the Clumio AWS connection resource model to be used as input to readAWSConnection()
	crm := clumioAWSConnectionResourceModel{
		ID:              basetypes.NewStringValue(id),
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Tests the success scenario for AWS connection read. It should not return Diagnostics.
	t.Run("success scenario for read aws connection", func(t *testing.T) {
		readResponse := &models.ReadAWSConnectionResponse{
			AccountNativeId:    &accountId,
			AwsRegion:          &region,
			ClumioAwsAccountId: &clumioAccountId,
			ClumioAwsRegion:    &region,
			ConnectionStatus:   &connStatus,
			DataPlaneAccountId: &dataplaneAccountId,
			Description:        &description,
			ExternalId:         &externalId,
			Id:                 &id,
			Namespace:          &namespace,
			Token:              &token,
		}
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsConnection(id, mock.Anything).Times(1).
			Return(readResponse, nil)

		remove, diags := cr.readAWSConnection(ctx, &crm)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that in case the AWS connection is not found, it returns true to indicate that the AWS
	// connection should be removed from the state.
	t.Run("read aws connection returns not found error", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsConnection(id, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := cr.readAWSConnection(ctx, &crm)
		assert.True(t, remove)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read AWS connection API call returns an error.
	t.Run("read aws connection returns error", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsConnection(id, mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := cr.readAWSConnection(ctx, &crm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read AWS connection API call returns an empty
	// response.
	t.Run("read aws connection returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsConnection(id, mock.Anything).Times(1).
			Return(nil, nil)

		remove, diags := cr.readAWSConnection(ctx, &crm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update AWS connection success scenario.
//   - Update AWS connection success scenario where only OU changes.
//   - SDK API to update OU for AWS connection returns an error.
//   - SDK API for read Clumio AWS connection returns an error.
//   - SDK API for create Clumio AWS connection returns nil response.
func TestUpdateAWSConnection(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockAWSConnectionClient(t)
	mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
	mockOrgUnitsCient := sdkclients.NewMockOrganizationalUnitClient(t)
	mockTaskClient := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	cr := clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections:  mockAwsConnClient,
		sdkEnvironments: mockAwsEnvClient,
		sdkOrgUnits:     mockOrgUnitsCient,
		sdkTasks:        mockTaskClient,
		pollTimeout:     5 * time.Second,
		pollInterval:    1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	descUpdated := "test-description-updated"

	// Populate the Clumio AWS connection resource model to be used as plan in updateAWSConnection()
	plan := clumioAWSConnectionResourceModel{
		ID:              basetypes.NewStringValue(id),
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Populate the Clumio AWS connection resource model to be used as state in updateAWSConnection()
	state := clumioAWSConnectionResourceModel{
		ID:              basetypes.NewStringValue(id),
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Tests the success scenario for AWS connection update. It should not return Diagnostics.
	t.Run("success scenario for update aws connection", func(t *testing.T) {

		plan.Description = basetypes.NewStringValue(descUpdated)
		updateResponse := &models.UpdateAWSConnectionResponse{
			AccountNativeId:    &accountId,
			AwsRegion:          &region,
			ClumioAwsAccountId: &clumioAccountId,
			ClumioAwsRegion:    &region,
			ConnectionStatus:   &connStatus,
			DataPlaneAccountId: &dataplaneAccountId,
			Description:        &description,
			ExternalId:         &externalId,
			Id:                 &id,
			Namespace:          &namespace,
			Token:              &token,
		}

		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(updateResponse, nil)

		diags := cr.updateAWSConnection(ctx, &plan, &state)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update AWS connection API call returns an
	// error.
	t.Run("update aws connection returns an error", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, apiError)

		plan.Description = basetypes.NewStringValue(description)
		diags := cr.updateAWSConnection(ctx, &plan, &state)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update AWS connection API call returns an
	// empty response.
	t.Run("read aws connection returns an empty response", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, nil)

		plan.Description = basetypes.NewStringValue(description)
		diags := cr.updateAWSConnection(ctx, &plan, &state)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete AWS connection success scenario.
//   - Delete AWS connection should not return error if AWS connection is not found.
//   - SDK API for delete AWS connection returns an error.
func TestDeleteAWSConnection(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockAWSConnectionClient(t)
	ctx := context.Background()
	cr := clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockAwsConnClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the Clumio AWS connection resource model to be used as input to deleteAWSConnection()
	crm := &clumioAWSConnectionResourceModel{
		ID:              basetypes.NewStringValue(id),
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Tests the success scenario for AWS connection deletion. It should not return diag.Diagnostics.
	t.Run("Success scenario for aws connection deletion", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().DeleteAwsConnection(id).Times(1).Return(nil, nil)

		diags := cr.deleteAWSConnection(ctx, crm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the AWS connection does not exist.
	t.Run("AWS connection not found should not return error", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().DeleteAwsConnection(id).Times(1).Return(nil, apiNotFoundError)

		diags := cr.deleteAWSConnection(ctx, crm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete AWS connection API call returns an error.
	t.Run("deleteAWSConnection returns an error", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().DeleteAwsConnection(id).Times(1).Return(nil, apiError)

		diags := cr.deleteAWSConnection(ctx, crm)
		assert.NotNil(t, diags)
	})

}
