// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsEnvironmentsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	awsenvironments "github.com/clumio-code/clumio-go-sdk/controllers/aws_environments"
)

type AWSEnvironmentClient interface {
	awsenvironments.AwsEnvironmentsV1Client
}

func NewAWSEnvironmentClient(config config.Config) AWSEnvironmentClient {
	return awsenvironments.NewAwsEnvironmentsV1(config)
}
