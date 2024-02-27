// Copyright 2024. Clumio, Inc.

package common

import (
	"fmt"
	"reflect"
	"testing"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Test all common utils
func TestUtils(t *testing.T) {
	t.Run("ParseMessageFromApiError - Parses and returns stringified response from api error", func(t *testing.T) {
		mockResponse := "{\"errors\":[{\"error_code\":111,\"error_message\":\"The request is invalid.\"}]}"
		mockByteArray := []byte(fmt.Sprintf("%v", mockResponse))
		mockApiError := apiutils.NewAPIError("test-reason", 500, mockByteArray)

		res := ParseMessageFromApiError(mockApiError)
		testResult := reflect.DeepEqual(res, mockResponse)
		if !testResult {
			t.Fatalf(TestResultsNotMatchingError, res, mockResponse)
		}
	})

	t.Run("ParseMessageFromApiError - Returns custom error message for auth error", func(t *testing.T) {
		mockResponse := "{\"errors\":[{\"error_code\":111,\"error_message\":\"The request is invalid.\"}]}"
		mockByteArray := []byte(fmt.Sprintf("%v", mockResponse))
		mockApiError := apiutils.NewAPIError("test-reason", 401, mockByteArray)

		res := ParseMessageFromApiError(mockApiError)
		testResult := reflect.DeepEqual(res, AuthError)
		if !testResult {
			t.Fatalf(TestResultsNotMatchingError, res, AuthError)
		}
	})

	t.Run("GetFieldNameFromNestedBlockPath - Returns the correct field name from path", func(t *testing.T) {
		mockPath := "operations[Value({\"action_setting\":\"window\",\"advanced_settings\":<null>,\"backup_aws_region\":<null>,\"backup_window_tz\":[{\"end_time\":\"07:00\",\"start_time\":\"05:00\"}],\"slas\":[{\"retention_duration\":<null>,\"rpo_frequency\":[{\"offsets\":<null>,\"unit\":\"days\",\"value\":1}]}],\"type\":\"aws_ebs_volume_backup\"})].slas[Value({\"retention_duration\":<null>,\"rpo_frequency\":[{\"offsets\":<null>,\"unit\":\"days\",\"value\":1}]})].retention_duration"
		expectedFieldName := "retention_duration"

		res := GetFieldNameFromNestedBlockPath(validator.SetRequest{
			Path: path.Root(mockPath),
		})
		testResult := reflect.DeepEqual(res, expectedFieldName)
		if !testResult {
			t.Fatalf(TestResultsNotMatchingError, expectedFieldName, res)
		}
	})
}
