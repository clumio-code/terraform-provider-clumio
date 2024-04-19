// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the client initialization functions.

//go:build unit

package sdkclients

import (
	"testing"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/stretchr/testify/assert"
)

// Unit tests for all the client initialization functions.
func TestClientInitialization(t *testing.T) {

	config := sdkconfig.Config{}

	auprClient := NewAutoUserProvisioningRuleClient(config)
	assert.NotNil(t, auprClient)

	aupsClient := NewAutoUserProvisioningSettingClient(config)
	assert.NotNil(t, aupsClient)

	awsConnClient := NewAWSConnectionClient(config)
	assert.NotNil(t, awsConnClient)

	awsEnvClient := NewAWSEnvironmentClient(config)
	assert.NotNil(t, awsEnvClient)

	ouClient := NewOrganizationalUnitClient(config)
	assert.NotNil(t, ouClient)

	paClient := NewPolicyAssignmentClient(config)
	assert.NotNil(t, paClient)

	pdClient := NewPolicyDefinitionClient(config)
	assert.NotNil(t, pdClient)

	prClient := NewPolicyRuleClient(config)
	assert.NotNil(t, prClient)

	postProcessConnClient := NewPostProcessAWSConnectionClient(config)
	assert.NotNil(t, postProcessConnClient)

	postProcessKmsClient := NewPostProcessKMSClient(config)
	assert.NotNil(t, postProcessKmsClient)

	pgClient := NewProtectionGroupClient(config)
	assert.NotNil(t, pgClient)

	roleClient := NewRoleClient(config)
	assert.NotNil(t, roleClient)

	taskClient := NewTaskClient(config)
	assert.NotNil(t, taskClient)

	userClient := NewUserClient(config)
	assert.NotNil(t, userClient)

	walletClient := NewWalletClient(config)
	assert.NotNil(t, walletClient)

	templatesClient := NewAWSTemplatesClient(config)
	assert.NotNil(t, templatesClient)

	s3BucketClient := NewS3BucketClient(config)
	assert.NotNil(t, s3BucketClient)
}
