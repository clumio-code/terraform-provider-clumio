// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK PostProcessAwsConnectionV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkPostProcessConn "github.com/clumio-code/clumio-go-sdk/controllers/post_process_aws_connection"
)

type PostProcessAWSConnectionClient interface {
	sdkPostProcessConn.PostProcessAwsConnectionV1Client
}

func NewPostProcessAWSConnectionClient(config config.Config) PostProcessAWSConnectionClient {
	return sdkPostProcessConn.NewPostProcessAwsConnectionV1(config)
}
