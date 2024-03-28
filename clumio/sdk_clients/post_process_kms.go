// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsConnectionsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkPostProcessKms "github.com/clumio-code/clumio-go-sdk/controllers/post_process_kms"
)

type PostProcessKMSClient interface {
	sdkPostProcessKms.PostProcessKmsV1Client
}

func NewPostProcessKMSClient(config config.Config) PostProcessKMSClient {
	return sdkPostProcessKms.NewPostProcessKmsV1(config)
}
