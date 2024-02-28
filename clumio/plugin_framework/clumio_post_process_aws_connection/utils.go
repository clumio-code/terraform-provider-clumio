// Copyright 2024. Clumio, Inc.

// This file hold various utility functions and variables used by the clumio_post_process_aws_connection
// Terraform resource.

package clumio_post_process_aws_connection

import (
	"errors"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
)

var (
	// protectInfoMap is the mapping of the datasource to the resource parameter and
	// if a config section is required, then isConfig will be true.
	protectInfoMap = map[string]sourceConfigInfo{
		"ebs": {
			sourceKey: "ProtectEBSVersion",
			isConfig:  false,
		},
		"rds": {
			sourceKey: "ProtectRDSVersion",
			isConfig:  false,
		},
		"ec2_mssql": {
			sourceKey: "ProtectEC2MssqlVersion",
			isConfig:  false,
		},
		"warm_tier": {
			sourceKey: "ProtectWarmTierVersion",
			isConfig:  true,
		},
		"s3": {
			sourceKey: "ProtectS3Version",
			isConfig:  false,
		},
		"dynamodb": {
			sourceKey: "ProtectDynamoDBVersion",
			isConfig:  false,
		},
	}
	// warmtierInfoMap is the mapping of the the warm tier datasource to the resource
	// parameter and if a config section is required, then isConfig will be true.
	warmtierInfoMap = map[string]sourceConfigInfo{
		"dynamodb": {
			sourceKey: "ProtectWarmTierDynamoDBVersion",
			isConfig:  false,
		},
	}
)

// GetTemplateConfiguration generates the template configuration from the schema.
func GetTemplateConfiguration(
	model postProcessAWSConnectionResourceModel, isCamelCase bool, isConsolidated bool) (
	map[string]interface{}, error) {

	templateConfigs := make(map[string]interface{})
	configMap, err := getConfigMapForKey(model.ConfigVersion.ValueString(), false)
	if err != nil {
		return nil, err
	}
	if configMap == nil {
		return templateConfigs, nil
	}
	templateConfigs["config"] = configMap
	if !isConsolidated {
		discoverMap, err := getConfigMapForKey(model.DiscoverVersion.ValueString(), true)
		if err != nil {
			return nil, err
		}
		if discoverMap == nil {
			return templateConfigs, nil
		}
		templateConfigs["discover"] = discoverMap
	}

	protectMap, err := getConfigMapForKey(model.ProtectConfigVersion.ValueString(), true)
	if err != nil {
		return nil, err
	}
	if protectMap == nil {
		return templateConfigs, nil
	}
	err = populateConfigMap(model, protectInfoMap, protectMap, isCamelCase)
	if err != nil {
		return nil, err
	}
	warmTierKey := "warm_tier"
	if isCamelCase {
		warmTierKey = common.SnakeCaseToCamelCase(warmTierKey)
	}
	if protectWarmtierMap, ok := protectMap[warmTierKey]; ok {
		err = populateConfigMap(
			model, warmtierInfoMap, protectWarmtierMap.(map[string]interface{}), isCamelCase)
		if err != nil {
			return nil, err
		}
	}
	if isConsolidated {
		templateConfigs["consolidated"] = protectMap
	} else {
		templateConfigs["protect"] = protectMap
	}
	return templateConfigs, nil
}

// populateConfigMap returns protect configuration information for the configs
// in the configInfoMap.
func populateConfigMap(model postProcessAWSConnectionResourceModel,
	configInfoMap map[string]sourceConfigInfo, configMap map[string]interface{},
	isCamelCase bool) error {

	for source, sourceInfo := range configInfoMap {
		configMapKey := source
		if isCamelCase {
			configMapKey = common.SnakeCaseToCamelCase(source)
		}
		reflectVal := reflect.ValueOf(&model)
		fieldVal := reflect.Indirect(reflectVal).FieldByName(sourceInfo.sourceKey)
		fieldStringVal := fieldVal.Interface().(types.String)
		protectSourceMap, err := getConfigMapForKey(
			fieldStringVal.ValueString(), sourceInfo.isConfig)
		if err != nil {
			return err
		}
		if protectSourceMap != nil {
			configMap[configMapKey] = protectSourceMap
		}
	}
	return nil
}

// getConfigMapForKey returns a config map for the key if it exists in ResourceData.
func getConfigMapForKey(val string, isConfig bool) (map[string]interface{}, error) {

	var mapToReturn map[string]interface{}
	if val != "" {
		keyMap := make(map[string]interface{})
		majorVersion, minorVersion, err := parseVersion(val)
		if err != nil {
			return nil, err
		}
		keyMap["enabled"] = true
		keyMap["version"] = majorVersion
		keyMap["minorVersion"] = minorVersion
		mapToReturn = keyMap
		// If isConfig is true it wraps the keyMap with another map with "config" as the key.
		if isConfig {
			configMap := make(map[string]interface{})
			configMap["config"] = keyMap
			mapToReturn = configMap
		}
	}
	return mapToReturn, nil
}

// parseVersion parses the version and minorVersion given the version string.
func parseVersion(version string) (string, string, error) {

	splits := strings.Split(version, ".")
	switch len(splits) {
	case 1:
		return version, "", nil
	case 2:
		return splits[0], splits[1], nil
	default:
		return "", "", errors.New(fmt.Sprintf("Invalid version: %v", version))
	}
}
