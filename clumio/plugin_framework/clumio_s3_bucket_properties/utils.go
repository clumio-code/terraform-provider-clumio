// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_s3_bucket_properties Terraform
// resource.

package clumio_s3_bucket_properties

import (
	"context"
	"errors"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
)

// pollForS3Bucket queries the S3 bucket till the EventBridgeEnabled value is the same in the API
// response and the one in the given model.
func (r *clumioS3BucketPropertiesResource) pollForS3Bucket(ctx context.Context,
	bucketId string, model *clumioS3BucketPropertiesResourceModel, interval time.Duration,
	timeout time.Duration) error {

	ticker := time.NewTicker(interval)
	tickerTimeout := time.After(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("context canceled or timed out")
		case <-ticker.C:
			readResponse, err := r.sdkS3BucketClient.ReadAwsS3Bucket(bucketId)
			if err != nil {
				return errors.New(common.ParseMessageFromApiError(err))
			}
			if readResponse == nil ||
				*readResponse.EventBridgeEnabled != model.EventBridgeEnabled.ValueBool() {
				continue
			}
			return nil
		case <-tickerTimeout:
			return errors.New("polling timed out")
		}
	}
}
