// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_dynamo_db_tables

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Read DynamoDB tables success scenario.
//   - SDK API for read DynamoDB tables returns an error.
//   - SDK API for read DynamoDB tables returns an empty response.
func TestDatasourceReadDynamoDBTables(t *testing.T) {

	ctx := context.Background()
	ouClient := sdkclients.NewMockDynamoDBTableClient(t)
	name := "test-s3-bucket"
	resourceName := "test_dynamodb_tables"
	id := "test-s3-bucket-id"
	ou := "test-ou"
	testError := "Test Error"
	accountNativeId := "test-account-native-id"
	tableNativeId := "test-table-native-id"
	region := "test-region"

	rds := clumioDynamoDBTablesDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		dynamoDBTableClient: ouClient,
	}

	rdsm := &clumioDynamoDBTablesDataSourceModel{
		AccountNativeID: basetypes.NewStringValue(accountNativeId),
		Region:          basetypes.NewStringValue(region),
		TableNativeID:   basetypes.NewStringValue(tableNativeId),
		Name:            basetypes.NewStringValue(name),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for DynamoDB tables read. It should not return Diagnostics.
	t.Run("Basic success scenario for read DynamoDB tables", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListDynamoDBTableResponse{
			Embedded: &models.DynamoDBTableListEmbedded{
				Items: []*models.DynamoDBTable{
					{
						Id:   &id,
						Name: &ou,
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		ouClient.EXPECT().ListAwsDynamodbTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Times(1).Return(readResponse, nil)

		diags := rds.readDynamoDBTables(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list DynamoDB tables API call returns an
	// error.
	t.Run("list DynamoDB tables returns an error", func(t *testing.T) {

		// Setup expectations.
		ouClient.EXPECT().ListAwsDynamodbTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Times(1).Return(nil, apiError)

		diags := rds.readDynamoDBTables(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list DynamoDB tables API call returns an
	// empty response.
	t.Run("list DynamoDB tables returns an empty response", func(t *testing.T) {

		// Setup expectations.
		ouClient.EXPECT().ListAwsDynamodbTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Times(1).Return(nil, nil)

		diags := rds.readDynamoDBTables(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list DynamoDB tables API call returns an
	// empty response.
	t.Run("list DynamoDB tables returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListDynamoDBTableResponse{
			Embedded: &models.DynamoDBTableListEmbedded{
				Items: []*models.DynamoDBTable{},
			},
		}

		// Setup expectations.
		ouClient.EXPECT().ListAwsDynamodbTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Times(1).Return(readResponse, nil)

		diags := rds.readDynamoDBTables(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
