// Copyright 2024. Clumio, Inc.

//go:build unit

package common

import (
	"context"
	"fmt"
	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	"github.com/clumio-code/clumio-go-sdk/models"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
	"time"
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

// Unit test for the utility function PollTask
func TestPollTask(t *testing.T) {

	mockTaskClient := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	taskId := "12345"

	t.Run("Success scenario", func(t *testing.T) {
		status := TaskSuccess
		readTaskResponse := models.ReadTaskResponse{
			Status: &status,
		}
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(&readTaskResponse, nil)
		err := PollTask(ctx, mockTaskClient, taskId, 5*time.Second, 1)
		assert.Nil(t, err)
	})

	t.Run("Read task returns aborted status", func(t *testing.T) {
		status := TaskAborted
		readTaskResponse := models.ReadTaskResponse{
			Status: &status,
		}
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(&readTaskResponse, nil)
		err := PollTask(ctx, mockTaskClient, taskId, 5*time.Second, 1)
		assert.NotNil(t, err)
	})

	t.Run("Read task returns failed status", func(t *testing.T) {
		status := TaskFailed
		readTaskResponse := models.ReadTaskResponse{
			Status: &status,
		}
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(&readTaskResponse, nil)
		err := PollTask(ctx, mockTaskClient, taskId, 5*time.Second, 1)
		assert.NotNil(t, err)
	})

	t.Run("Read task returns error", func(t *testing.T) {
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(nil,
			&apiutils.APIError{
				ResponseCode: http.StatusInternalServerError,
				Reason:       "Test",
				Response:     []byte("Test Error"),
			})
		err := PollTask(ctx, mockTaskClient, taskId, 5*time.Second, 1)
		assert.NotNil(t, err)
	})
}

// Unit test for the utility function PollForProtectionGroup.
// Tests the following scenarios:
//   - Success scenario for protection group polling.
//   - Read protection group returns HTTP 404 leading to polling timeout.
//   - Read protection group returns an error.
func TestPollForProtectionGroup(t *testing.T) {

	pgClient := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pgId := "12345"

	// Success scenario for protection group polling.
	t.Run("Success scenario", func(t *testing.T) {
		readPGResponse := models.ReadProtectionGroupResponse{
			Id: &pgId,
		}
		pgClient.EXPECT().ReadProtectionGroup(pgId).Times(1).Return(&readPGResponse, nil)
		res, err := PollForProtectionGroup(ctx, pgId, pgClient, 5*time.Second, 1)
		assert.Nil(t, err)
		assert.Equal(t, pgId, *res.Id)
	})

	// Read protection group returns HTTP 404 leading to polling timeout.
	t.Run("Polling timeout", func(t *testing.T) {
		notFoundError := apiutils.NewAPIError("Not found", http.StatusNotFound, nil)
		pgClient.EXPECT().ReadProtectionGroup(pgId).Times(1).Return(nil, notFoundError)
		res, err := PollForProtectionGroup(ctx, pgId, pgClient, 1, 1)
		assert.NotNil(t, err)
		assert.Equal(t, "polling timed out", err.Error())
		assert.Nil(t, res)
	})

	// Read protection group returns an error.
	t.Run("Read protection group returns an error", func(t *testing.T) {
		apiError := apiutils.NewAPIError("Test Error", http.StatusInternalServerError, nil)
		pgClient.EXPECT().ReadProtectionGroup(pgId).Times(1).Return(nil, apiError)
		res, err := PollForProtectionGroup(ctx, pgId, pgClient, 5*time.Second, 1)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

// Unit test for the utility function SliceDifferenceString.
func TestSliceDifferenceString(t *testing.T) {

	slice1 := []string{"test1", "test2", "test3"}
	slice2 := []string{"test2", "test3", "test4"}

	diff := SliceDifferenceString(slice1, slice2)
	assert.Equal(t, 1, len(diff))
	assert.Equal(t, "test1", diff[0])

	diff = SliceDifferenceString(slice2, slice1)
	assert.Equal(t, 1, len(diff))
	assert.Equal(t, "test4", diff[0])
}
