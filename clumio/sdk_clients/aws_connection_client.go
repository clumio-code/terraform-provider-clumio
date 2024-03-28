// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsConnectionsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/aws_connections"
)

type AWSConnectionClient interface {
	aws_connections.AwsConnectionsV1Client
}

func NewAWSConnectionClient(config config.Config) AWSConnectionClient {
	return aws_connections.NewAwsConnectionsV1(config)
}
