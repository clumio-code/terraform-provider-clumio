// Copyright 2021. Clumio, Inc.

// Contains the util functions used by the providers and resources

package clumio

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getStringValue returns the string value of the key if present.
func getStringValue(d *schema.ResourceData, key string) string {
	value := ""
	if d.Get(key) != nil {
		value = fmt.Sprintf("%v", d.Get(key))
	}
	return value
}

// getStringSlice returns the string slice of the key if present.
func getStringSlice(d *schema.ResourceData, key string) []*string {
	var value []*string
	if d.Get(key) != nil {
		value := make([]*string, 0)
		valSlice :=  d.Get(key).([]interface{})
		for _, val := range valSlice{
			strVal := val.(string)
			value = append(value, &strVal)
		}
	}
	return value
}


//Utility function to return a string value from a map if the key exists
func getStringValueFromMap(keyVals map[string]interface{}, key string) *string{
	if v, ok := keyVals[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

// Utility function to determine if it is unit or acceptance test
func isAcceptanceTest() bool{
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
func snakeCaseToCamelCase(key string)string{
	newKey := key
	if strings.Contains(key, "_"){
		parts := strings.Split(key, "_")
		newKey = parts[0]
		for _, part := range parts[1:]{
			newKey = newKey + strings.Title(part)
		}
	}
	return newKey
}