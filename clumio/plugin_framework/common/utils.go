// Copyright 2021. Clumio, Inc.

// Contains the util functions used by the providers and resources

package common

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/clumio-code/clumio-go-sdk/controllers/tasks"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
