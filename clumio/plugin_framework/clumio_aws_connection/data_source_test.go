// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_aws_connection

import (
	"context"
	"testing"

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
//   - Read aws connection success scenario.
//   - SDK API for read aws connection returns an error.
//   - SDK API for read aws connection returns an empty response.
//   - SDK API for read aws connection returns a response with empty items.
func TestDatasourceReadAWSConnection(t *testing.T) {

	ctx := context.Background()
	pgClient := sdkclients.NewMockAWSConnectionClient(t)
	resourceName := "test_aws_connection"
	id := "test-aws-connection-id"
	testError := "Test Error"
	region = "test-region"
	accountId = "test-account-id"

	rds := clumioAWSConnectionDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		awsConnectionClient: pgClient,
	}

	rdsm := &clumioAWSConnectionDataSourceModel{
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for aws connection read. It should not return Diagnostics.
	t.Run("Basic success scenario for read aws connection", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListAWSConnectionsResponse{
			Embedded: &models.AWSConnectionListEmbedded{
				Items: []*models.AWSConnection{
					{
						Id:              &id,
						AccountNativeId: &accountId,
						AwsRegion:       &region,
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		pgClient.EXPECT().ListAwsConnections(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readAWSConnection(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list aws connections API call returns an
	// error.
	t.Run("list aws connections returns an error", func(t *testing.T) {

		// Setup expectations.
		pgClient.EXPECT().ListAwsConnections(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rds.readAWSConnection(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list aws connections API call returns an
	// empty response.
	t.Run("list aws connections returns an empty response", func(t *testing.T) {

		// Setup expectations.
		pgClient.EXPECT().ListAwsConnections(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rds.readAWSConnection(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list aws connections API call returns a
	// response with empty items.
	t.Run("list aws connections returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListAWSConnectionsResponse{
			Embedded: &models.AWSConnectionListEmbedded{
				Items: []*models.AWSConnection{},
			},
		}

		// Setup expectations.
		pgClient.EXPECT().ListAwsConnections(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readAWSConnection(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
