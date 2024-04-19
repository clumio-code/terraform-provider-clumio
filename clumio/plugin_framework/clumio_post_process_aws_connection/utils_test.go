// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_post_process_aws_connection

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following cases:
//   - The template configuration genrated contains the correct information.
//   - Getting template version with invalid config version returns an error.
//   - Getting template version with empty config version returns an empty config.
//   - Getting template version with invalid protect config version returns an error.
//   - Getting template version with empty protect config version returns an empty protect config.
//   - Getting template version with invalid asset version returns an error.
//   - Getting template version with invalid dynamoDB warm-tier version returns an error.
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

	// Test that the template configuration generated contains the correct information as
	// per the values in the model.
	t.Run("Basic success scenario that all configuration generated", func(t *testing.T) {
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
	})

	t.Run("Returns an error with invalid config version", func(t *testing.T) {

		prmWithInvalidVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue(invalidVersion),
		}

		config, err := GetTemplateConfiguration(prmWithInvalidVersion, true)
		assert.NotNil(t, err)
		assert.Nil(t, config)
	})

	t.Run("Returns an empty configuration with empty config version", func(t *testing.T) {

		prmWithEmptyVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue(emptyVersion),
		}

		config, err := GetTemplateConfiguration(prmWithEmptyVersion, true)
		assert.Nil(t, err)
		assert.Empty(t, config)
	})

	t.Run("Returns an error with invalid protect config version", func(t *testing.T) {

		prmWithInvalidVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion:        basetypes.NewStringValue(version),
			ProtectConfigVersion: basetypes.NewStringValue(invalidVersion),
		}

		config, err := GetTemplateConfiguration(prmWithInvalidVersion, true)
		assert.NotNil(t, err)
		assert.Nil(t, config)
	})

	t.Run("Returns an template configuration with empty protect config version", func(t *testing.T) {

		prmWithEmptyVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion:        basetypes.NewStringValue(version),
			ProtectConfigVersion: basetypes.NewStringValue(emptyVersion),
		}

		config, err := GetTemplateConfiguration(prmWithEmptyVersion, true)
		assert.Nil(t, err)
		assert.NotNil(t, config)

		configMap := config["config"].(map[string]interface{})
		assert.Equal(t, true, configMap["enabled"].(bool))
		assert.Equal(t, "1", configMap["version"].(string))
		assert.Equal(t, "0", configMap["minorVersion"])
	})

	t.Run("Returns an error with invalid asset version", func(t *testing.T) {

		prmWithInvalidVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion:        basetypes.NewStringValue(version),
			ProtectConfigVersion: basetypes.NewStringValue(version),
			ProtectEBSVersion:    basetypes.NewStringValue(invalidVersion),
		}

		config, err := GetTemplateConfiguration(prmWithInvalidVersion, true)
		assert.NotNil(t, err)
		assert.Nil(t, config)
	})

	t.Run("Returns an error with invalid warm tier dynamoDB version", func(t *testing.T) {

		prmWithInvalidVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion:                  basetypes.NewStringValue(version),
			ProtectConfigVersion:           basetypes.NewStringValue(version),
			ProtectEBSVersion:              basetypes.NewStringValue(version),
			ProtectWarmTierVersion:         basetypes.NewStringValue(version),
			ProtectWarmTierDynamoDBVersion: basetypes.NewStringValue(invalidVersion),
		}

		config, err := GetTemplateConfiguration(prmWithInvalidVersion, true)
		t.Log(err)
		assert.NotNil(t, err)
		assert.Nil(t, config)
	})

}

// Unit test for the following cases:
//   - Parse version with one character.
//   - Parse version with decimal point.
//   - Parse invalid version string returns an error.
func TestParseVersion(t *testing.T) {

	t.Run("Parse version with one character", func(t *testing.T) {
		version := "1"
		majorVersion, minorVersion, err := parseVersion(version)
		assert.Nil(t, err)
		assert.Equal(t, "1", majorVersion)
		assert.Equal(t, "", minorVersion)
	})

	t.Run("Parse version with decimal point", func(t *testing.T) {
		version := "1.2"
		majorVersion, minorVersion, err := parseVersion(version)
		assert.Nil(t, err)
		assert.Equal(t, "1", majorVersion)
		assert.Equal(t, "2", minorVersion)
	})

	t.Run("Parse invalid version string", func(t *testing.T) {
		version := "1.2.3"
		majorVersion, minorVersion, err := parseVersion(version)
		assert.NotNil(t, err)
		assert.Equal(t, "", majorVersion)
		assert.Equal(t, "", minorVersion)
	})
}
