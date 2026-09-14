// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in data_source.go

//go:build unit

package clumio_iceberg_tables

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
//   - Read Iceberg tables success scenario.
//   - The environment lookup returns an error.
//   - SDK API for read Iceberg tables returns an error.
//   - SDK API for read Iceberg tables returns an empty response.
//   - SDK API for read Iceberg tables returns empty items in the response.
func TestDatasourceReadIcebergTables(t *testing.T) {

	ctx := context.Background()
	icebergClient := sdkclients.NewMockIcebergTableClient(t)
	envClient := sdkclients.NewMockAWSEnvironmentClient(t)

	resourceName := "test_iceberg_tables"
	accountNativeId := "test-account-native-id"
	region := "test-region"
	environmentId := "test-environment-id"
	id := "test-iceberg-table-id"
	name := "test-iceberg-table"
	catalog := "test-catalog"
	namespace := "test-namespace"
	catalogType := catalogTypeGlue
	testError := "Test Error"

	ds := clumioIcebergTablesDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		icebergTableClient:   icebergClient,
		awsEnvironmentClient: envClient,
	}

	dsm := &clumioIcebergTablesDataSourceModel{
		AccountNativeID: basetypes.NewStringValue(accountNativeId),
		Region:          basetypes.NewStringValue(region),
		Name:            basetypes.NewStringValue(name),
		Catalog:         basetypes.NewStringValue(catalog),
		Namespace:       basetypes.NewStringValue(namespace),
		CatalogType:     basetypes.NewStringValue(catalogType),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	envResponse := &models.ListAWSEnvironmentsResponse{
		Embedded: &models.AWSEnvironmentListEmbedded{
			Items: []*models.AWSEnvironment{
				{
					Id: &environmentId,
				},
			},
		},
	}

	// Tests the success scenario for Iceberg tables read. It should not return Diagnostics.
	t.Run("Basic success scenario for read Iceberg tables", func(t *testing.T) {

		readResponse := &models.ListIcebergTablesResponse{
			Embedded: &models.IcebergTableListEmbedded{
				Items: []*models.IcebergTable{
					{
						Id:          &id,
						Name:        &name,
						Catalog:     &catalog,
						CatalogType: &catalogType,
						Namespace:   &namespace,
					},
				},
			},
		}

		// Setup expectations.
		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(envResponse, nil)
		icebergClient.EXPECT().ListAwsIcebergTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).RunAndReturn(
			func(_ *int64, _ *string, filter *string, _ *string, _ *int64) (
				*models.ListIcebergTablesResponse, *apiutils.APIError) {
				// The API rejects an aws_region filter, so the region must be translated into
				// the environment ID.
				assert.NotContains(t, *filter, schemaRegion)
				assert.Contains(t, *filter, environmentId)
				return readResponse, nil
			}).Times(1)

		diags := ds.readIcebergTables(ctx, dsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned if the environment cannot be resolved.
	t.Run("environment lookup returns an error", func(t *testing.T) {

		// Setup expectations.
		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		diags := ds.readIcebergTables(ctx, dsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list Iceberg tables API call returns an
	// error.
	t.Run("list Iceberg tables returns an error", func(t *testing.T) {

		// Setup expectations.
		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(envResponse, nil)
		icebergClient.EXPECT().ListAwsIcebergTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		diags := ds.readIcebergTables(ctx, dsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list Iceberg tables API call returns an
	// empty response.
	t.Run("list Iceberg tables returns an empty response", func(t *testing.T) {

		// Setup expectations.
		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(envResponse, nil)
		icebergClient.EXPECT().ListAwsIcebergTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, nil)

		diags := ds.readIcebergTables(ctx, dsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list Iceberg tables API call returns
	// empty items in the response.
	t.Run("list Iceberg tables returns empty items in response", func(t *testing.T) {

		readResponse := &models.ListIcebergTablesResponse{
			Embedded: &models.IcebergTableListEmbedded{
				Items: []*models.IcebergTable{},
			},
		}

		// Setup expectations.
		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(envResponse, nil)
		icebergClient.EXPECT().ListAwsIcebergTables(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(readResponse, nil)

		diags := ds.readIcebergTables(ctx, dsm)
		assert.NotNil(t, diags)
	})
}
