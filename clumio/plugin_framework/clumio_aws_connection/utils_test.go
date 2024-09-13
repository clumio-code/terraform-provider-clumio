// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_aws_connection

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"testing"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
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
