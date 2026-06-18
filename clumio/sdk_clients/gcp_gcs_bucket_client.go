// Copyright (c) 2026 Clumio, a Commvault Company All Rights Reserved

// Contains the wrapper interface for Clumio GO SDK GcpGcsBucketsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	gcpgcsbuckets "github.com/clumio-code/clumio-go-sdk/controllers/gcp_gcs_buckets"
)

type GcpGcsBucketClient interface {
	gcpgcsbuckets.GcpGcsBucketsV1Client
}

func NewGcpGcsBucketClient(config config.Config) GcpGcsBucketClient {
	return gcpgcsbuckets.NewGcpGcsBucketsV1(config)
}
