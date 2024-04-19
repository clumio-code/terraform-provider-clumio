// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsS3BucketsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkS3Buckets "github.com/clumio-code/clumio-go-sdk/controllers/aws_s3_buckets"
)

type S3BucketClient interface {
	sdkS3Buckets.AwsS3BucketsV1Client
}

func NewS3BucketClient(config config.Config) S3BucketClient {
	return sdkS3Buckets.NewAwsS3BucketsV1(config)
}
