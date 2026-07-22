// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in data_source.go

//go:build unit

package clumio_gcs_bucket

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Read GCS bucket success scenario.
//   - SDK API for list GCS buckets returns an error.
//   - SDK API for list GCS buckets returns a nil response.
//   - SDK API for list GCS buckets returns an empty items list.
func TestDatasourceReadGCSBucket(t *testing.T) {

	ctx := context.Background()
	mockClient := sdkclients.NewMockGcpGcsBucketClient(t)
	resourceName := "test_gcs_bucket"
	bucketName := "test-gcs-bucket"
	id := "test-gcs-bucket-id"
	projectId := "test-project-id"
	projectUuid := "test-project-uuid"
	location := "us-central1"
	locationType := "Region"
	locationUuid := "test-location-uuid"
	ou := "test-ou"
	testError := "Test Error"
	objectCount := int64(42)
	sizeBytes := int64(1024)
	pgCount := int64(2)
	isDeleted := false
	isVersioningEnabled := true
	createdTimestamp := "test-created-timestamp"
	lastBackupTimestamp := "test-last-backup-timestamp"

	rds := clumioGCSBucketDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		gcsBucketClient: mockClient,
	}

	bucketNames, conversionDiags := types.SetValueFrom(ctx, types.StringType, []string{bucketName})
	assert.Nil(t, conversionDiags)
	rdsm := &clumioGCSBucketDataSourceModel{
		BucketNames: bucketNames,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for GCS bucket read. It should not return Diagnostics and should
	// populate the gcs_buckets attribute.
	t.Run("Basic success scenario for read GCS bucket", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListGCSBucketsResponse{
			Embedded: &models.GCSBucketListEmbedded{
				Items: []*models.GCSBucket{
					{
						Id:                   &id,
						BucketName:           &bucketName,
						ProjectId:            &projectId,
						ProjectUuid:          &projectUuid,
						Location:             &location,
						LocationType:         &locationType,
						LocationUuid:         &locationUuid,
						OrganizationalUnitId: &ou,
						ObjectCount:          &objectCount,
						SizeBytes:            &sizeBytes,
						ProtectionGroupCount: &pgCount,
						IsDeleted:            &isDeleted,
						IsVersioningEnabled:  &isVersioningEnabled,
						CreatedTimestamp:     &createdTimestamp,
						LastBackupTimestamp:  &lastBackupTimestamp,
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		mockClient.EXPECT().ListGcpGcsBuckets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readGCSBucket(ctx, rdsm)
		assert.Nil(t, diags)
		assert.False(t, rdsm.GCSBuckets.IsNull())
		assert.Len(t, rdsm.GCSBuckets.Elements(), 1)
	})

	// Tests that Diagnostics is returned in case the list GCS buckets API call returns an error.
	t.Run("list GCS buckets returns an error", func(t *testing.T) {

		// Setup expectations.
		mockClient.EXPECT().ListGcpGcsBuckets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rds.readGCSBucket(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list GCS buckets API call returns a nil
	// response.
	t.Run("list GCS buckets returns a nil response", func(t *testing.T) {

		// Setup expectations.
		mockClient.EXPECT().ListGcpGcsBuckets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rds.readGCSBucket(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list GCS buckets API call returns an empty
	// items list in the response.
	t.Run("list GCS buckets returns an empty items list", func(t *testing.T) {

		readResponse := &models.ListGCSBucketsResponse{
			Embedded: &models.GCSBucketListEmbedded{
				Items: []*models.GCSBucket{},
			},
		}

		// Setup expectations.
		mockClient.EXPECT().ListGcpGcsBuckets(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readGCSBucket(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
