// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_s3_bucket

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
//   - Read S3 bucket success scenario.
//   - SDK API for read s3 bucket returns an error.
//   - SDK API for read s3 bucket returns an empty response.
func TestDatasourceReadS3Bucket(t *testing.T) {

	ctx := context.Background()
	ouClient := sdkclients.NewMockS3BucketClient(t)
	name := "test-s3-bucket"
	resourceName := "test_s3_bucket"
	id := "test-s3-bucket-id"
	ou := "test-ou"
	testError := "Test Error"
	eventBridgeEnabled := true
	accountNativeId := "test-account-native-id"
	region := "test-region"
	lastBackupTimestamp := "test-last-backup-timestamp"
	lastContinuousBackupTimestamp := "test-last-continuous-backup-timestamp"
	pgCount := int64(2)

	rds := clumioS3BucketDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		s3BucketClient: ouClient,
	}

	names, conversionDiags := types.SetValueFrom(ctx, types.StringType, []string{name})
	assert.Nil(t, conversionDiags)
	rdsm := &clumioS3BucketDataSourceModel{
		BucketNames: names,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for s3 bucket read. It should not return Diagnostics.
	t.Run("Basic success scenario for read s3 bucket", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListBucketsResponse{
			Embedded: &models.BucketListEmbedded{
				Items: []*models.Bucket{
					{
						Id:                            &id,
						Name:                          &ou,
						AccountNativeId:               &accountNativeId,
						AwsRegion:                     &region,
						ProtectionGroupCount:          &pgCount,
						EventBridgeEnabled:            &eventBridgeEnabled,
						LastBackupTimestamp:           &lastBackupTimestamp,
						LastContinuousBackupTimestamp: &lastContinuousBackupTimestamp,
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		ouClient.EXPECT().ListAwsS3Buckets(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readS3Bucket(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list s3 buckets API call returns an
	// error.
	t.Run("list s3 buckets returns an error", func(t *testing.T) {

		// Setup expectations.
		ouClient.EXPECT().ListAwsS3Buckets(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rds.readS3Bucket(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list s3 buckets API call returns an
	// empty response.
	t.Run("list s3 buckets returns an empty response", func(t *testing.T) {

		// Setup expectations.
		ouClient.EXPECT().ListAwsS3Buckets(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rds.readS3Bucket(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list s3 buckets API call returns an
	// empty response.
	t.Run("list s3 buckets returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListBucketsResponse{
			Embedded: &models.BucketListEmbedded{
				Items: []*models.Bucket{},
			},
		}

		// Setup expectations.
		ouClient.EXPECT().ListAwsS3Buckets(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readS3Bucket(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
