// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_protection_group_asset

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
//   - Read protection group asset success scenario.
//   - SDK API for reading protection group asset returns an error.
//   - SDK API for reading protection group asset returns an empty response.
func TestDatasourceReadProtectionGroupAsset(t *testing.T) {

	ctx := context.Background()
	s3AssetsClient := sdkclients.NewMockProtectionGroupS3AssetsClient(t)
	resourceName := "test_pg_asset"
	id := "test-pg-asset-id"
	pgId := "test-pg-id"
	bucketId := "test-bucket-id"
	testError := "Test Error"

	ds := clumioProtectionGroupAssetDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		s3AssetsClient: s3AssetsClient,
	}

	dsModel := &clumioProtectionGroupAssetDataSourceModel{
		ProtectionGroupID: basetypes.NewStringValue(pgId),
		BucketID:          basetypes.NewStringValue(bucketId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for protection group asset read. It should not return Diagnostics.
	t.Run("Basic success scenario for read protection group asset", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListProtectionGroupS3AssetsResponse{
			Embedded: &models.ProtectionGroupBucketListEmbedded{
				Items: []*models.ProtectionGroupBucket{
					{
						Id:       &id,
						BucketId: &bucketId,
						GroupId:  &pgId,
					},
				},
			},
			TotalCount: &count,
		}

		// Setup expectations.
		s3AssetsClient.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := ds.readProtectionGroupAsset(ctx, dsModel)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list protection group assets API call returns an
	// error.
	t.Run("list protection group assets returns an error", func(t *testing.T) {

		// Setup expectations.
		s3AssetsClient.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := ds.readProtectionGroupAsset(ctx, dsModel)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list protection group assets API call returns an
	// empty response.
	t.Run("list protection group assets returns an empty response", func(t *testing.T) {

		// Setup expectations.
		s3AssetsClient.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		diags := ds.readProtectionGroupAsset(ctx, dsModel)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list protection group assets API call returns
	// an empty response.
	t.Run("list protection group assets returns an empty items in response",
		func(t *testing.T) {

			readResponse := &models.ListProtectionGroupS3AssetsResponse{
				Embedded: &models.ProtectionGroupBucketListEmbedded{
					Items: []*models.ProtectionGroupBucket{},
				},
			}

			// Setup expectations.
			s3AssetsClient.EXPECT().ListProtectionGroupS3Assets(
				mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
				Return(readResponse, nil)

			diags := ds.readProtectionGroupAsset(ctx, dsModel)
			assert.NotNil(t, diags)
		})
}
