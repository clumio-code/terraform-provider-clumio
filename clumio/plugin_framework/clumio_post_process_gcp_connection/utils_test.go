// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

//go:build unit

package clumio_post_process_gcp_connection

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following cases:
//   - Parse version with one character.
//   - Parse version with decimal point.
//   - Parse invalid version string returns an error.
//   - Parse invalid empty string
//   - Parse invalid non-numeric major
//   - Parse invalid non-numeric minor
//   - Parse invalid arbitrary string
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
		majorVersion, minorVersion, err := parseVersion("1.2.3")
		assert.NotNil(t, err)
		assert.Equal(t, "", majorVersion)
		assert.Equal(t, "", minorVersion)

		majorVersion, minorVersion, err = parseVersion("Invalid_version")
		assert.NotNil(t, err)
		assert.Equal(t, "", majorVersion)
		assert.Equal(t, "", minorVersion)
	})

	t.Run("Parse invalid empty string", func(t *testing.T) {
		major, minor, err := parseVersion("")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})

	t.Run("Parse invalid non-numeric major", func(t *testing.T) {
		major, minor, err := parseVersion("a")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})

	t.Run("Parse invalid non-numeric minor", func(t *testing.T) {
		major, minor, err := parseVersion("1.a")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})

	t.Run("Parse invalid arbitrary string", func(t *testing.T) {
		major, minor, err := parseVersion("Invalid_version")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})

	t.Run("Parse version with consecutive dots", func(t *testing.T) {
		major, minor, err := parseVersion("1..2")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})

	t.Run("Parse version with trailing dot", func(t *testing.T) {
		major, minor, err := parseVersion("1.")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})

	t.Run("Parse version with leading dot", func(t *testing.T) {
		major, minor, err := parseVersion(".1")
		assert.Error(t, err)
		assert.Equal(t, "", major)
		assert.Equal(t, "", minor)
	})
}

// Unit test for the following cases:
//   - The template configuration generated contains the correct information.
//   - Getting template version with invalid config version returns an error.
//   - Getting template version with invalid gcs config version returns an error.
//   - Getting template version with empty config version returns as disabled.
//   - Getting template version with empty gcs config version returns as disabled.
func TestGetTemplateConfiguration(t *testing.T) {
	t.Run("template configuration generated contains the correct information", func(t *testing.T) {
		result, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion:     basetypes.NewStringValue("1.2"),
			ProtectGcsVersion: basetypes.NewStringValue("2.2"),
		})

		assert.Nil(t, err)
		config := result["config"].(map[string]any)
		assert.Equal(t, config["enabled"].(bool), true)
		assert.Equal(t, config["version"].(string), "1")
		assert.Equal(t, config["minor_version"].(string), "2")

		gcs := result["gcs"].(map[string]any)

		assert.Equal(t, gcs["enabled"].(bool), true)
		assert.Equal(t, gcs["version"].(string), "2")
		assert.Equal(t, gcs["minor_version"].(string), "2")
	})

	t.Run("template configuration handles null value for ProtectGcsVersion", func(t *testing.T) {
		_, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue("1.2"),
		})
		assert.Nil(t, err)

		_, err = GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion:     basetypes.NewStringValue("1.2"),
			ProtectGcsVersion: types.String{},
		})
		assert.Nil(t, err)
	})

	t.Run("template version with invalid config version returns an error", func(t *testing.T) {
		result, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue("1.2.3"),
		})

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("template version with invalid config semantic versioning returns an error", func(t *testing.T) {
		result, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue("invalid_versioning"),
		})

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("template version with invalid gcs version returns an error", func(t *testing.T) {
		result, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion:     basetypes.NewStringValue("1.2"),
			ProtectGcsVersion: basetypes.NewStringValue("2.2.2"),
		})

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("template version with empty config version returns as disabled", func(t *testing.T) {
		result, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue(""),
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)
		config := result["config"].(map[string]any)
		assert.Equal(t, config["enabled"].(bool), false)
	})

	t.Run("template version with empty gcs version returns as disabled", func(t *testing.T) {
		result, err := GetTemplateConfiguration(&clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion:     basetypes.NewStringValue("1.2"),
			ProtectGcsVersion: basetypes.NewStringValue(""),
		})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		config := result["gcs"].(map[string]any)
		assert.Equal(t, config["enabled"].(bool), false)
	})
}
