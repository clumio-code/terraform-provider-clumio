// Copyright 2024. Clumio, Inc.

// This file hold various utility functions and variables used by the clumio_post_process_aws_connection
// Terraform resource.

package clumio_post_process_aws_connection

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		"iceberg_on_glue": {
			sourceKey: "ProtectIcebergOnGlueVersion",
			isConfig:  false,
		},
		"iceberg_on_s3_table": {
			sourceKey: "ProtectIcebergOnS3TablesVersion",
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
	model postProcessAWSConnectionResourceModel, isCamelCase bool) (
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
	templateConfigs["consolidated"] = protectMap
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

// PollForConnectionIngestionAndTargetStatus polls till the connection ingestion and target setup
// status fields become either completed or failed.
func pollForConnectionIngestionAndTargetStatus(
	ctx context.Context, sdkAWSConnection sdkclients.AWSConnectionClient,
	model postProcessAWSConnectionResourceModel, timeout time.Duration,
	interval time.Duration) (bool, error) {

	ticker := time.NewTicker(interval)
	tickerTimeout := time.After(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return false, errors.New("context canceled or timed out")
		case <-ticker.C:
			connectionId := fmt.Sprintf("%s_%s",
				model.AccountID.ValueString(), model.Region.ValueString())
			returnExternalId := "false"
			// Call the Clumio API to read the AWS connection.
			res, apiErr := sdkAWSConnection.ReadAwsConnection(connectionId, &returnExternalId)
			if apiErr != nil {
				return false, errors.New(common.ParseMessageFromApiError(apiErr))
			}
			shouldReturn, targetSetupError, err := performValidation(res, model)
			if shouldReturn {
				return targetSetupError, err
			}
		case <-tickerTimeout:
			return false, errors.New("polling timed out")
		}
	}
}

// performValidation checks the status from the response to determine whether ReadAwsConnection
// needs to be performed again.
func performValidation(res *models.ReadAWSConnectionResponse,
	model postProcessAWSConnectionResourceModel) (bool, bool, error) {

	if res == nil {
		return false, false, nil
	}
	ingestionComplete, targetSetupComplete := true, true
	ingestionErr, targetSetupErr := false, false
	ingestionComplete, ingestionErr = isIngestionComplete(
		model.WaitForIngestion.ValueBool(), *res.IngestionStatus)
	targetSetupComplete, targetSetupErr = isTargetSetupComplete(
		model.WaitForDataPlaneResources.ValueBool(), *res.TargetSetupStatus)
	if ingestionErr && targetSetupErr {
		return true, true, errors.New("Ingestion task failed for the connection as well as" +
			" one or more of the data plane resources setup tasks failed.")
	} else if ingestionErr {
		return true, false, errors.New("Ingestion task failed for the connection.")
	} else if targetSetupErr {
		return true, true, errors.New(
			"One or more of the data plane resources setup tasks failed.")
	}
	if ingestionComplete && targetSetupComplete {
		return true, false, nil
	}
	return false, false, nil
}

// isIngestionComplete returns true if either the ingestionStatus is either completed or failed or if
// waitForIngesion is false. Returns false when ingestionStatus is in_progress.
func isIngestionComplete(waitForIngestion bool, ingestionStatus string) (bool, bool) {

	if waitForIngestion {
		if ingestionStatus == inProgress {
			return false, false
		} else if ingestionStatus == failed {
			return true, true
		}
	}
	return true, false
}

// isTargetSetupComplete returns true if either the targetSetupStatus is either completed or failed
// or if waitForTargetSetup is false. Returns false when targetSetupStatus is in_progress.
func isTargetSetupComplete(waitForTargetSetup bool, targetSetupStatus string) (bool, bool) {

	if waitForTargetSetup {
		if targetSetupStatus == inProgress {
			return false, false
		} else if targetSetupStatus == failed {
			return true, true
		}
	}
	return true, false
}
