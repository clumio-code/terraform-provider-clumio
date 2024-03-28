// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_post_process_aws_connection

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test to test that the template configuration generated contains the correct information as
// per the values in the model.
func TestGetTemplateConfiguration(t *testing.T) {

	prm := postProcessAWSConnectionResourceModel{
		ID:                             basetypes.NewStringValue(id),
		AccountID:                      basetypes.NewStringValue(accountId),
		Token:                          basetypes.NewStringValue(token),
		RoleExternalID:                 basetypes.NewStringValue(externalId),
		Region:                         basetypes.NewStringValue(region),
		ClumioEventPubID:               basetypes.NewStringValue(eventPubId),
		RoleArn:                        basetypes.NewStringValue(roleArn),
		ConfigVersion:                  basetypes.NewStringValue(version),
		DiscoverVersion:                basetypes.NewStringValue(version),
		ProtectConfigVersion:           basetypes.NewStringValue(version),
		ProtectEBSVersion:              basetypes.NewStringValue(version),
		ProtectS3Version:               basetypes.NewStringValue(version),
		ProtectDynamoDBVersion:         basetypes.NewStringValue(version),
		ProtectEC2MssqlVersion:         basetypes.NewStringValue(version),
		ProtectWarmTierVersion:         basetypes.NewStringValue(version),
		ProtectWarmTierDynamoDBVersion: basetypes.NewStringValue(version),
		IntermediateRoleArn:            basetypes.NewStringValue(intermediateRoleArn),
	}

	config, err := GetTemplateConfiguration(prm, true)
	assert.Nil(t, err)
	assert.NotNil(t, config)

	// Check if the versions are populated as expected. Checking for one of the data sources.
	consolidatedMap := config["consolidated"].(map[string]interface{})
	s3Map := consolidatedMap["s3"].(map[string]interface{})
	assert.Equal(t, s3Map["version"].(string), "1")
	assert.Equal(t, s3Map["minorVersion"].(string), "0")

	// Validate that if the version for a data source is not specified, it is not part of the
	// config map.
	_, ok := consolidatedMap["rds"]
	assert.False(t, ok)
}
