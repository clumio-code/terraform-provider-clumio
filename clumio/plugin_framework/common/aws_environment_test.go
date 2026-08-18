// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for LookupAWSEnvironment.

//go:build unit

package common

import (
	"testing"

	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Success: exactly one environment with an ID is returned.
//   - API error, nil response, empty items, more than one match, and missing ID all error.
func TestLookupAWSEnvironment(t *testing.T) {

	accountNativeId := "test-account-native-id"
	awsRegion := "test-region"
	envId := "test-env-id"

	// Success scenario: a single environment with an ID is returned without error.
	t.Run("returns the environment on a single match", func(t *testing.T) {

		client := sdkclients.NewMockAWSEnvironmentClient(t)
		client.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(&models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}, nil)

		env, err := LookupAWSEnvironment(client, accountNativeId, awsRegion)
		assert.NoError(t, err)
		assert.Equal(t, envId, *env.Id)
	})

	// Error scenarios: each should return an error and no environment.
	twoId := "id-2"
	tests := []struct {
		name string
		resp *models.ListAWSEnvironmentsResponse
		err  error
	}{
		{name: "api error", err: &apiutils.APIError{ResponseCode: 500, Reason: "test"}},
		{name: "nil response", resp: nil},
		{name: "empty items", resp: &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{Items: []*models.AWSEnvironment{}}}},
		{name: "more than one match", resp: &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}, {Id: &twoId}}}}},
		{name: "missing ID", resp: &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: nil}}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := sdkclients.NewMockAWSEnvironmentClient(t)
			var apiErr *apiutils.APIError
			if tt.err != nil {
				apiErr = tt.err.(*apiutils.APIError)
			}
			client.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
				mock.Anything, mock.Anything).Times(1).Return(tt.resp, apiErr)

			env, err := LookupAWSEnvironment(client, accountNativeId, awsRegion)
			assert.Error(t, err)
			assert.Nil(t, env)
		})
	}
}
