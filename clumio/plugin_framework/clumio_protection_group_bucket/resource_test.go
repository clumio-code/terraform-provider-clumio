// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_protection_group_bucket

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

var (
	name         = "test-protection-group-bucket"
	resourceName = "test_protection_group_bucket"
	id           = "mock-pg-bucket-id"
	testError    = "Test Error"
	bucketId     = "test-bucket-id"
	pgId         = "test-pg-id"
)

// Unit test for the following cases:
//   - Create protection group bucket success scenario.
//   - SDK API for add bucket to protection group returns an error.
//   - SDK API for add bucket to protection group returns an empty response.
func TestCreateProtectionGroupBucket(t *testing.T) {

	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupBucketResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the protection group resource model to be used as input to createProtectionGroupBucket()
	pgrm := clumioProtectionGroupBucketResourceModel{
		BucketID:          basetypes.NewStringValue(bucketId),
		ProtectionGroupID: basetypes.NewStringValue(pgId),
	}

	// Create the response of the SDK CreateProtectionGroupDefinition() API.
	createResponse := &models.AddBucketToProtectionGroupResponse{
		Id:       &id,
		BucketId: &bucketId,
		GroupId:  &pgId,
	}

	// Tests the success scenario for protection group create. It should not return Diagnostics.
	t.Run("Basic success scenario for create protection group", func(t *testing.T) {

		// Setup Expectations
		mockProtectionGroup.EXPECT().AddBucketProtectionGroup(pgId, mock.Anything).Times(1).
			Return(createResponse, nil)

		diags := pr.createProtectionGroupBucket(ctx, &pgrm)
		assert.Nil(t, diags)
		assert.Equal(t, pgrm.ID.ValueString(), *createResponse.Id)
	})

	// Tests that Diagnostics is returned in case the add bucket to protection group API call
	// returns an error.
	t.Run("CreateProtectionGroupBucket returns error", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().AddBucketProtectionGroup(pgId, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.createProtectionGroupBucket(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the add bucket to protection group API call
	// returns an empty response.
	t.Run("CreateProtectionGroupBucket returns an empty response", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().AddBucketProtectionGroup(pgId, mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.createProtectionGroupBucket(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read protection group bucket success scenario.
//   - SDK API for read protection group s3 assets returns not found error.
//   - SDK API for read protection group s3 assets returns no S3 asset.
//   - SDK API for read protection group s3 assets returns an error.
//   - SDK API for read protection group s3 assets returns an empty response.
func TestReadProtectionGroupBucket(t *testing.T) {

	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	mockSDKProtectionGroupS3Assets := sdkclients.NewMockProtectionGroupS3AssetsClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupBucketResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
		sdkS3Assets:         mockSDKProtectionGroupS3Assets,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	count := int64(1)

	// Populate the protection group resource model to be used as input to readProtectionGroupBucket()
	pgrm := clumioProtectionGroupBucketResourceModel{
		ID:                basetypes.NewStringValue(id),
		BucketID:          basetypes.NewStringValue(bucketId),
		ProtectionGroupID: basetypes.NewStringValue(pgId),
	}
	// Tests the success scenario for protection group bucket read. It should not return Diagnostics.
	t.Run("Basic success scenario for read protection group", func(t *testing.T) {

		// Create the response of the SDK ReadProtectionGroupDefinition() API.
		listResponse := &models.ListProtectionGroupS3AssetsResponse{
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

		// Setup Expectations
		mockSDKProtectionGroupS3Assets.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything).Times(1).Return(listResponse, nil)

		remove, diags := pr.readProtectionGroupBucket(ctx, &pgrm)
		assert.False(t, remove)
		assert.Nil(t, diags)
		assert.Equal(t, pgrm.ID.ValueString(), *listResponse.Embedded.Items[0].Id)
		assert.Equal(t, pgrm.BucketID.ValueString(), *listResponse.Embedded.Items[0].BucketId)
		assert.Equal(t, pgrm.ID.ValueString(), *listResponse.Embedded.Items[0].Id)
	})

	// Tests that Diagnostics is returned in case the ListProtectionGroupS3Assets API call returns
	// HTTP 404 error.
	t.Run("ReadProtectionGroup returns http 404 error", func(t *testing.T) {
		// Setup Expectations
		mockSDKProtectionGroupS3Assets.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := pr.readProtectionGroupBucket(context.Background(), &pgrm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the ListProtectionGroupS3Assets API call returns
	// no S3 asset.
	t.Run("ReadProtectionGroup returns no S3 asset", func(t *testing.T) {

		// Create the response of the SDK ReadProtectionGroupDefinition() API.
		count = 0
		listResponse := &models.ListProtectionGroupS3AssetsResponse{
			TotalCount: &count,
		}
		// Setup Expectations
		mockSDKProtectionGroupS3Assets.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listResponse, nil)

		remove, diags := pr.readProtectionGroupBucket(context.Background(), &pgrm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// error.
	t.Run("ReadProtectionGroup returns an error", func(t *testing.T) {
		// Setup Expectations
		mockSDKProtectionGroupS3Assets.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		remove, diags := pr.readProtectionGroupBucket(context.Background(), &pgrm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// empty response.
	t.Run("ReadProtectionGroup returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockSDKProtectionGroupS3Assets.EXPECT().ListProtectionGroupS3Assets(
			mock.Anything, mock.Anything, mock.Anything).Times(1).Return(nil, nil)

		remove, diags := pr.readProtectionGroupBucket(context.Background(), &pgrm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Delete protection group success scenario.
//   - Delete protection group should not return an error if protection group is not found.
//   - SDK API for delete protection group returns an error.
func TestDeleteProtectionGroupBucket(t *testing.T) {

	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupBucketResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the protection group resource model to be used as input to readProtectionGroupBucket()
	pgrm := &clumioProtectionGroupBucketResourceModel{
		ID:                basetypes.NewStringValue(id),
		BucketID:          basetypes.NewStringValue(bucketId),
		ProtectionGroupID: basetypes.NewStringValue(pgId),
	}

	// Tests the success scenario for protection group deletion. It should not return
	// diag.Diagnostics.
	t.Run("Success scenario for protection group deletion", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().DeleteBucketProtectionGroup(pgId, bucketId).Times(1).
			Return(&models.DeleteBucketFromProtectionGroupResponse{}, nil)

		diags := pr.deleteProtectionGroupBucket(ctx, pgrm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the protection group does not exist.
	t.Run("Policy not found should not return error", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().DeleteBucketProtectionGroup(pgId, bucketId).Times(1).
			Return(nil, apiNotFoundError)

		diags := pr.deleteProtectionGroupBucket(ctx, pgrm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete protection group API call returns error.
	t.Run("deleteProtectionGroupBucket returns an error", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().DeleteBucketProtectionGroup(pgId, bucketId).Times(1).
			Return(nil, apiError)

		diags := pr.deleteProtectionGroupBucket(ctx, pgrm)
		assert.NotNil(t, diags)
	})

}
