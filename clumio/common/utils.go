// Copyright 2021. Clumio, Inc.

// Contains the util functions used by the providers and resources

package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/clumio-code/clumio-go-sdk/config"
	tasks "github.com/clumio-code/clumio-go-sdk/controllers/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetStringValue returns the string value of the key if present.
func GetStringValue(d *schema.ResourceData, key string) string {
	value := ""
	if d.Get(key) != nil {
		value = fmt.Sprintf("%v", d.Get(key))
	}
	return value
}

// GetBoolValueWithDefault returns the bool value of the key if present,
// otherwise a default value.
func GetBoolValueWithDefault(d *schema.ResourceData, key string, defaultVal bool) bool {
	val := d.Get(key)
	if val == nil {
		return defaultVal
	}
	return val.(bool)
}

// GetIntValueWithDefault returns the int64 value of the key if present,
// otherwise a default value.
func GetIntValueWithDefault(d *schema.ResourceData, key string, defaultVal int64) int64 {
	val := d.Get(key)
	if val == nil {
		return defaultVal
	}
	if ival, ok := val.(int64); ok {
		return ival
	}
	return defaultVal
}

// GetStringSlice returns the string slice of the key if present.
func GetStringSlice(d *schema.ResourceData, key string) []*string {
	var value []*string
	if _, ok := d.GetOk(key); ok {
		value = make([]*string, 0)
		valSlice := d.Get(key).([]interface{})
		for _, val := range valSlice {
			strVal := val.(string)
			value = append(value, &strVal)
		}
	}
	return value
}

// GetStringSliceFromSet returns the string slice from the string set
// of the key if present.
func GetStringSliceFromSet(d *schema.ResourceData, key string) []*string {
	var value []*string
	if _, ok := d.GetOk(key); ok {
		value = make([]*string, 0)
		valSlice := d.Get(key).(*schema.Set).List()
		for _, val := range valSlice {
			strVal := val.(string)
			value = append(value, &strVal)
		}
	}
	return value
}

// Utility function to determine if string in slice exists.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//Utility function to return a string value from a map if the key exists
func GetStringValueFromMap(keyVals map[string]interface{}, key string) *string {
	if v, ok := keyVals[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

// Utility function to determine if it is unit or acceptance test
func IsAcceptanceTest() bool {
	return os.Getenv("TF_ACC") == "true" || os.Getenv("TF_ACC") == "True" ||
		os.Getenv("TF_ACC") == "1"
}

// RequireOneOf verifies that at least one environment variable is non-empty or returns an error.
//
// If at lease one environment variable is non-empty, returns the first name and value.
func RequireOneOf(names []string, usageMessage string) (string, string, error) {
	for _, variable := range names {
		value := os.Getenv(variable)

		if value != "" {
			return variable, value, nil
		}
	}

	return "", "", fmt.Errorf("at least one environment variable of %v must be set. Usage: %s", names, usageMessage)
}

// function to convert string in snake case to camel case
func SnakeCaseToCamelCase(key string) string {
	newKey := key
	if strings.Contains(key, "_") {
		parts := strings.Split(key, "_")
		newKey = parts[0]
		for _, part := range parts[1:] {
			newKey = newKey + strings.Title(part)
		}
	}
	return newKey
}

func SliceDifference(slice1 []interface{}, slice2 []interface{}) []interface{} {
	var diff []interface{}

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

// GetStringSliceFromInterfaceSlice returns the string slice from interface slice.
func GetStringSliceFromInterfaceSlice(input []interface{}) []*string {
	strSlice := make([]*string, 0)
	for _, val := range input {
		strVal := val.(string)
		strSlice = append(strSlice, &strVal)
	}
	return strSlice
}

// SchemaSetHashString returns a hash for TypeSet's with string elements.
func SchemaSetHashString(i interface{}) int {
	return schema.HashString(i.(string))
}

type sourceConfigInfo struct {
	sourceKey string
	isConfig  bool
}

var (
	// protectInfoMap is the mapping of the the datasource to the resource parameter and
	// if a config section is required, then isConfig will be true.
	protectInfoMap = map[string]sourceConfigInfo{
		"ebs": {
			sourceKey: "protect_ebs_version",
			isConfig:  false,
		},
		"rds": {
			sourceKey: "protect_rds_version",
			isConfig:  false,
		},
		"ec2_mssql": {
			sourceKey: "protect_ec2_mssql_version",
			isConfig:  false,
		},
		"warm_tier": {
			sourceKey: "protect_warm_tier_version",
			isConfig:  true,
		},
		"s3": {
			sourceKey: "protect_s3_version",
			isConfig:  false,
		},
		"dynamodb": {
			sourceKey: "protect_dynamodb_version",
			isConfig:  false,
		},
	}
	// warmtierInfoMap is the mapping of the the warm tier datasource to the resource
	// parameter and if a config section is required, then isConfig will be true.
	warmtierInfoMap = map[string]sourceConfigInfo{
		"dynamodb": {
			sourceKey: "protect_warm_tier_dynamodb_version",
			isConfig:  false,
		},
	}
)

// GetTemplateConfiguration returns the template configuration.
func GetTemplateConfiguration(d *schema.ResourceData, isCamelCase bool, isConsolidated bool) (
	map[string]interface{}, error) {
	templateConfigs := make(map[string]interface{})
	configMap, err := getConfigMapForKey(d, "config_version", false)
	if err != nil {
		return nil, err
	}
	if configMap == nil {
		return templateConfigs, nil
	}
	templateConfigs["config"] = configMap
	if !isConsolidated {
		discoverMap, err := getConfigMapForKey(d, "discover_version", true)
		if err != nil {
			return nil, err
		}
		if discoverMap == nil {
			return templateConfigs, nil
		}
		templateConfigs["discover"] = discoverMap
	}

	protectMap, err := getConfigMapForKey(d, "protect_config_version", true)
	if err != nil {
		return nil, err
	}
	if protectMap == nil {
		return templateConfigs, nil
	}
	err = populateConfigMap(d, protectInfoMap, protectMap, isCamelCase)
	if err != nil {
		return nil, err
	}
	warmTierKey := "warm_tier"
	if isCamelCase {
		warmTierKey = SnakeCaseToCamelCase(warmTierKey)
	}
	if protectWarmtierMap, ok := protectMap[warmTierKey]; ok {
		err = populateConfigMap(
			d, warmtierInfoMap, protectWarmtierMap.(map[string]interface{}), isCamelCase)
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
func populateConfigMap(d *schema.ResourceData, configInfoMap map[string]sourceConfigInfo,
	configMap map[string]interface{}, isCamelCase bool) error {
	for source, sourceInfo := range configInfoMap {
		configMapKey := source
		if isCamelCase {
			configMapKey = SnakeCaseToCamelCase(source)
		}
		protectSourceMap, err := getConfigMapForKey(
			d, sourceInfo.sourceKey, sourceInfo.isConfig)
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
func getConfigMapForKey(
	d *schema.ResourceData, key string, isConfig bool) (map[string]interface{}, error) {
	var mapToReturn map[string]interface{}
	if val, ok := d.GetOk(key); ok {
		keyMap := make(map[string]interface{})
		if keyVersion, ok := val.(string); ok {
			majorVersion, minorVersion, err := parseVersion(keyVersion)
			if err != nil {
				return nil, err
			}
			keyMap["enabled"] = true
			keyMap["version"] = majorVersion
			keyMap["minorVersion"] = minorVersion
		}
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

// PollTask polls created tasks to ensure that the resource
// was created successfully.
func PollTask(ctx context.Context, apiClient *ApiClient,
	taskId string, timeoutInSec int64, intervalInSec int64) error {
	t := tasks.NewTasksV1(apiClient.ClumioConfig)
	interval := time.Duration(intervalInSec) * time.Second
	ticker := time.NewTicker(interval)
	timeout := time.After(time.Duration(timeoutInSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resp, err := t.ReadTask(taskId)
			if err != nil {
				return err
			} else if *resp.Status == TaskSuccess {
				return nil
			} else if *resp.Status == TaskAborted {
				return errors.New("task aborted")
			} else if *resp.Status == TaskFailed {
				return errors.New("task failed")
			}
		case <-timeout:
			return errors.New("polling task timeout")
		}
	}
}

// GetClumioConfigForAPI returns the clumio config to be used for invoking APIs.
func GetClumioConfigForAPI(
	client *ApiClient, d *schema.ResourceData) config.Config {
	clumioConfig := client.ClumioConfig
	orgUnitId := GetStringValue(d, OrganizationalUnitId)
	if orgUnitId != "" {
		clumioConfig = config.Config{
			Token:                     client.ClumioConfig.Token,
			BaseUrl:                   client.ClumioConfig.BaseUrl,
			OrganizationalUnitContext: orgUnitId,
		}
	}
	return clumioConfig
}
