// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_s3_bucket_properties

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
	resourceName       = "test_s3_bucket_properties"
	id                 = "test-id"
	bucketId           = "test-bucket-id"
	eventBridgeEnabled = true
	testError          = "Test Error"
)

// Unit test for the following cases:
//   - Create or Update S3 bucket properties success scenario.
//   - Create or Update S3 bucket properties success scenario where the first check for
//     EventBridgeEnabled fails.
//   - SDK API for create or update S3 bucket properties returns error.
//   - SDK API for create or update S3 bucket properties returns nil response.
func TestCreateOrUpdateS3BucketProperties(t *testing.T) {

	mockS3BucketClient := sdkclients.NewMockS3BucketClient(t)
	ctx := context.Background()
	resource := clumioS3BucketPropertiesResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkS3BucketClient: mockS3BucketClient,
		pollTimeout:       5 * time.Second,
		pollInterval:      1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the Clumio S3 bucket properties resource model to be used as input to createOrUpdateS3BucketProperties()
	resourceModel := clumioS3BucketPropertiesResourceModel{
		BucketID:                        basetypes.NewStringValue(bucketId),
		EventBridgeEnabled:              basetypes.NewBoolValue(eventBridgeEnabled),
		EvendBridgeNotificationDisabled: basetypes.NewBoolValue(true),
	}

	// Tests the success scenario for clumio S3 bucket properties create or update. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for create or update S3 bucket properties",
		func(t *testing.T) {

			readResponse := &models.ReadBucketResponse{
				Id:                 &bucketId,
				EventBridgeEnabled: &eventBridgeEnabled,
			}

			// Setup Expectations
			mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
				Return(nil, nil)
			mockS3BucketClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(
				readResponse, nil)

			diags := resource.createOrUpdateS3BucketProperties(ctx, &resourceModel, true)
			assert.Nil(t, diags)
			assert.Equal(t, *readResponse.EventBridgeEnabled,
				resourceModel.EventBridgeEnabled.ValueBool())
		})

	// Tests the success scenario for clumio S3 bucket properties create. It should not return
	// Diagnostics. Tests the case where the first read S3 bucket returns the eventBridgeEnabled
	// value as false causing another call to read S3 bucket.
	t.Run("Success scenario for create S3 bucket properties", func(t *testing.T) {

		enabled := false
		firstReadResponse := &models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &enabled,
		}
		readResponse := &models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &eventBridgeEnabled,
		}

		// Setup Expectations
		mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
			Return(nil, nil)
		mockS3BucketClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(firstReadResponse, nil)
		mockS3BucketClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(readResponse, nil)

		diags := resource.createOrUpdateS3BucketProperties(ctx, &resourceModel, true)
		assert.Nil(t, diags)
		assert.Equal(t, *readResponse.EventBridgeEnabled,
			resourceModel.EventBridgeEnabled.ValueBool())
	})

	// Tests that Diagnostics is returned in case the set bucket properties call returns an error.
	t.Run("Set bucket properties returns an error", func(t *testing.T) {
		// Setup Expectations
		mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := resource.createOrUpdateS3BucketProperties(ctx, &resourceModel, true)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read S3 bucket properties API call returns an
	// error.
	t.Run("Read S3 bucket returns an error", func(t *testing.T) {
		// Setup Expectations
		mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
			Return(nil, nil)
		mockS3BucketClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(nil, apiError)

		diags := resource.createOrUpdateS3BucketProperties(ctx, &resourceModel, true)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read S3 bucket properties success scenario.
//   - SDK API for read S3 bucket returns not found error.
//   - SDK API for read Clumio S3 bucket returns error.
//   - SDK API for create Clumio S3 bucket returns nil response.
func TestReadS3BucketProperties(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockS3BucketClient(t)
	ctx := context.Background()
	cr := clumioS3BucketPropertiesResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkS3BucketClient: mockAwsConnClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the Clumio S3 bucket properties resource model to be used as input to readS3BucketProperties()
	crm := clumioS3BucketPropertiesResourceModel{
		ID:                 basetypes.NewStringValue(id),
		BucketID:           basetypes.NewStringValue(bucketId),
		EventBridgeEnabled: basetypes.NewBoolValue(eventBridgeEnabled),
	}

	// Tests the success scenario for S3 bucket properties read. It should not return Diagnostics.
	t.Run("success scenario for read S3 bucket properties", func(t *testing.T) {
		readResponse := &models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &eventBridgeEnabled,
		}
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).
			Return(readResponse, nil)

		remove, diags := cr.readS3BucketProperties(ctx, &crm)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that in case the S3 bucket properties is not found, it returns true to indicate that
	// the S3 bucket properties resource should be removed from the state.
	t.Run("read S3 bucket returns not found error", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(nil, apiNotFoundError)

		remove, diags := cr.readS3BucketProperties(ctx, &crm)
		assert.True(t, remove)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read S3 bucket properties API call returns an
	// error.
	t.Run("read S3 bucket returns error", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(nil, apiError)

		remove, diags := cr.readS3BucketProperties(ctx, &crm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read S3 bucket properties API call returns an empty
	// response.
	t.Run("read S3 bucket returns nil response", func(t *testing.T) {
		// Setup Expectations
		mockAwsConnClient.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(nil, nil)

		remove, diags := cr.readS3BucketProperties(ctx, &crm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Delete S3 bucket properties success scenario.
//   - Delete S3 bucket properties should not return error if S3 bucket properties is not found.
//   - SDK API for delete S3 bucket properties returns an error.
func TestDeleteS3BucketProperties(t *testing.T) {

	mockS3BucketClient := sdkclients.NewMockS3BucketClient(t)
	ctx := context.Background()
	resource := clumioS3BucketPropertiesResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkS3BucketClient: mockS3BucketClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the Clumio S3 bucket properties resource model to be used as input to createOrUpdateS3BucketProperties()
	resourceModel := clumioS3BucketPropertiesResourceModel{
		ID:                              basetypes.NewStringValue(id),
		BucketID:                        basetypes.NewStringValue(bucketId),
		EventBridgeEnabled:              basetypes.NewBoolValue(eventBridgeEnabled),
		EvendBridgeNotificationDisabled: basetypes.NewBoolValue(true),
	}

	// Tests the success scenario for clumio S3 bucket properties delete. It should not return
	// Diagnostics.
	t.Run("Basic success scenario for delete S3 bucket properties", func(t *testing.T) {

		// Setup Expectations
		mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
			Return(nil, nil)

		diags := resource.deleteS3BucketProperties(ctx, &resourceModel)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the set bucket properties call returns an error.
	t.Run("Set bucket properties returns an error", func(t *testing.T) {
		// Setup Expectations
		mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := resource.deleteS3BucketProperties(ctx, &resourceModel)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the set bucket properties call returns an error.
	t.Run("Set bucket properties returns not found error", func(t *testing.T) {
		// Setup Expectations
		mockS3BucketClient.EXPECT().SetBucketProperties(bucketId, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		diags := resource.deleteS3BucketProperties(ctx, &resourceModel)
		assert.Nil(t, diags)
	})

}
