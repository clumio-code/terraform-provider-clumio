// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_s3_bucket_properties

import (
	"context"
	"testing"
	"time"

	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the utility function pollForS3Bucket.
// Tests the following scenarios:
//   - Success scenario for S3 bucket polling.
//   - Success scenario for S3 bucket polling with first response not containing expected result.
//   - Read S3 bucket returns an error.
func TestPollForS3Bucket(t *testing.T) {
	s3Client := sdkclients.NewMockS3BucketClient(t)
	ctx := context.Background()
	bucketId := "test-bucket-id"
	eventBridgeEnabled := true

	res := &clumioS3BucketPropertiesResource{
		sdkS3BucketClient: s3Client,
	}

	resModel := &clumioS3BucketPropertiesResourceModel{
		BucketID:           basetypes.NewStringValue(bucketId),
		EventBridgeEnabled: basetypes.NewBoolValue(eventBridgeEnabled),
	}

	apiError := apiutils.NewAPIError("test", 500, []byte("test"))
	// Success scenario for S3 bucket polling.
	t.Run("Success scenario", func(t *testing.T) {
		readResponse := models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &eventBridgeEnabled,
		}
		s3Client.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(&readResponse, nil)
		err := res.pollForS3Bucket(ctx, bucketId, resModel, 1, 5*time.Second)
		assert.Nil(t, err)
		assert.Equal(t, bucketId, *readResponse.Id)
	})

	// Success scenario for S3 bucket polling with the first API call not returning expected result.
	t.Run("Success scenario with second API call", func(t *testing.T) {
		enabled := false
		firstReadResponse := models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &enabled,
		}
		readResponse := models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &eventBridgeEnabled,
		}
		s3Client.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(&firstReadResponse, nil)
		s3Client.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(&readResponse, nil)
		err := res.pollForS3Bucket(ctx, bucketId, resModel, 1, 5*time.Second)
		assert.Nil(t, err)
		assert.Equal(t, bucketId, *readResponse.Id)
	})

	// Failure scenario where the Read S3 bucket API returns an error.
	t.Run("Read S3 bucket returns an error", func(t *testing.T) {
		s3Client.EXPECT().ReadAwsS3Bucket(bucketId).Times(1).Return(nil, apiError)
		err := res.pollForS3Bucket(ctx, bucketId, resModel, 1, 5*time.Second)
		assert.NotNil(t, err)
	})

	// Read S3 bucket with canceled context returns an error.
	t.Run("Context canceled", func(t *testing.T) {
		doneCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(-1*time.Hour))
		cancelFunc()
		assert.NotNil(t, doneCtx.Done())
		err := res.pollForS3Bucket(doneCtx, bucketId, resModel, 1, 1)
		assert.NotNil(t, err)
		assert.Equal(t, "context canceled or timed out", err.Error())
	})

}

// Test for timeout during polling of S3 bucket.
func TestPollS3BucketPollingTimeout(t *testing.T) {
	s3Client := sdkclients.NewMockS3BucketClient(t)
	ctx := context.Background()
	bucketId := "test-bucket-id"
	eventBridgeEnabled := true

	res := &clumioS3BucketPropertiesResource{
		sdkS3BucketClient: s3Client,
	}

	resModel := &clumioS3BucketPropertiesResourceModel{
		BucketID:           basetypes.NewStringValue(bucketId),
		EventBridgeEnabled: basetypes.NewBoolValue(eventBridgeEnabled),
	}

	// Read S3 bucket  first API call not returning expected result leading to polling timeout.
	t.Run("Polling timeout", func(t *testing.T) {
		enabled := false
		firstReadResponse := models.ReadBucketResponse{
			Id:                 &bucketId,
			EventBridgeEnabled: &enabled,
		}
		s3Client.EXPECT().ReadAwsS3Bucket(bucketId).Return(&firstReadResponse, nil)
		err := res.pollForS3Bucket(
			ctx, bucketId, resModel, 1, 100)
		assert.NotNil(t, err)
		assert.Equal(t, "polling timed out", err.Error())
	})
}
