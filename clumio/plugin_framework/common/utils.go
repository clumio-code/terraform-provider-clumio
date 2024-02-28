// Copyright 2021. Clumio, Inc.

// Contains the util functions used by the providers and resources

package common

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	"github.com/clumio-code/clumio-go-sdk/controllers/tasks"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

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
			resp, apiErr := t.ReadTask(taskId)
			if apiErr != nil {
				return errors.New(ParseMessageFromApiError(apiErr))
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

// SliceDifferenceAttrValue returns the slice difference in attribute value slices.
func SliceDifferenceAttrValue(slice1 []attr.Value, slice2 []attr.Value) []attr.Value {
	var diff []attr.Value

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

// GetStringSliceFromAttrValueSlice returns the string slice from attribute value slice.
func GetStringSliceFromAttrValueSlice(input []attr.Value) []*string {
	strSlice := make([]*string, 0)
	for _, val := range input {
		strVal := val.String()
		strSlice = append(strSlice, &strVal)
	}
	return strSlice
}

// SliceDifferenceString returns the slice difference in string slices.
func SliceDifferenceString(slice1 []string, slice2 []string) []string {
	var diff []string

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

// GetStringPtrSliceFromStringSlice returns the string pointer slice from string slice.
func GetStringPtrSliceFromStringSlice(input []string) []*string {
	strSlice := make([]*string, 0)
	for _, val := range input {
		strVal := val
		strSlice = append(strSlice, &strVal)
	}
	return strSlice
}

// GetStringPtr returns the ptr of the string if not null, otherwise return nil.
func GetStringPtr(v basetypes.StringValue) *string {
	if v.IsNull() {
		return nil
	} else {
		s := v.ValueString()
		return &s
	}
}

// Parses the api error and returns the response in stringified format
func ParseMessageFromApiError(apiError *apiutils.APIError) string {
	// Handle auth errors separately
	if apiError.ResponseCode == http.StatusUnauthorized || apiError.ResponseCode == http.StatusForbidden {
		return AuthError
	}
	return string(apiError.Response)
}

// Parses through the path of a nested block and returns the lowest level field name
func GetFieldNameFromNestedBlockPath(req validator.SetRequest) string {
	// A nested block path looks something like this -
	// operations[Value({"backup_window_tz":[{"end_time":"07:00","start_time":"05:00"}],
	// "slas":[{"retention_duration":<null>,"rpo_frequency":[{"offsets":<null>,"unit":"days","value":1}]}],"type":"aws_ebs_volume_backup"})]
	// .slas[Value({"retention_duration":<null>,"rpo_frequency":[{"offsets":<null>,"unit":"days","value":1}]})].retention_duration
	// "retention_duration" is extracted from it as the lowest level field name
	return req.Path.String()[strings.LastIndex(req.Path.String(), ".")+1:]
}
