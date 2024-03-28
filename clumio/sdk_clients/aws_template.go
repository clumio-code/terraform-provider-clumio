// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsTemplatesV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	awstemplates "github.com/clumio-code/clumio-go-sdk/controllers/aws_templates"
)

type AWSTemplatesClient interface {
	awstemplates.AwsTemplatesV1Client
}

func NewAWSTemplatesClient(config config.Config) AWSTemplatesClient {
	return awstemplates.NewAwsTemplatesV1(config)
}
