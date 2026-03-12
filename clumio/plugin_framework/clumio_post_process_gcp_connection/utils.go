// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_post_process_gcp_connection

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const InvalidVersionError = "invalid version %v"

// configVersionMap is the mapping of resource model version key to the config value
var configVersionMap = map[string]func(
	*clumioPostProcessGCPConnectionResourceModel) types.String{
	"config": func(m *clumioPostProcessGCPConnectionResourceModel) types.String {
		return m.ConfigVersion
	},
	"gcs": func(m *clumioPostProcessGCPConnectionResourceModel) types.String {
		return m.ProtectGcsVersion
	},
}

// GetTemplateConfiguration generates the template configuration from the schema.
func GetTemplateConfiguration(
	model *clumioPostProcessGCPConnectionResourceModel) (
	map[string]any, error) {

	templateConfigs := make(map[string]any)
	for configKey, getVersion := range configVersionMap {
		ts := getVersion(model)
		configMap, err := getConfigMapForKey(ts.ValueString())

		if err != nil {
			return nil, err
		}
		templateConfigs[configKey] = configMap
	}

	return templateConfigs, nil
}

// getConfigMapForKey returns a config map for the key if it exists in ResourceData.
func getConfigMapForKey(val string) (map[string]any, error) {

	if val == "" {
		return map[string]any{
			"enabled":       false,
			"version":       "0",
			"minor_version": "0",
		}, nil
	}

	majorVersion, minorVersion, err := parseVersion(val)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"enabled":       true,
		"version":       majorVersion,
		"minor_version": minorVersion,
	}, nil
}

// parseVersion parses the version and minorVersion given the version string.
func parseVersion(version string) (major string, minor string, err error) {
	if version == "" {
		return "", "", fmt.Errorf(InvalidVersionError, version)
	}

	parts := strings.Split(version, ".")
	if len(parts) > 2 {
		return "", "", fmt.Errorf(InvalidVersionError, version)
	}

	for _, p := range parts {
		if p == "" {
			return "", "", fmt.Errorf(InvalidVersionError, version)
		}
		if _, convErr := strconv.Atoi(p); convErr != nil {
			return "", "", fmt.Errorf(InvalidVersionError, version)
		}
	}

	if len(parts) == 1 {
		return parts[0], "", nil
	}

	return parts[0], parts[1], nil
}
