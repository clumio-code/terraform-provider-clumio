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
	"github.com/clumio-code/clumio-go-sdk/models"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
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

// PollTask polls created tasks till it completes either with success, aborted, failed
// or it returns an error.
func PollTask(ctx context.Context, taskClient tasks.TasksV1Client,
	taskId string, timeout time.Duration, interval time.Duration) error {

	ticker := time.NewTicker(interval)
	tickerTimeout := time.After(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resp, apiErr := taskClient.ReadTask(taskId)
			if apiErr != nil {
				return errors.New(ParseMessageFromApiError(apiErr))
			} else if *resp.Status == TaskSuccess {
				return nil
			} else if *resp.Status == TaskAborted {
				return errors.New("task aborted")
			} else if *resp.Status == TaskFailed {
				return errors.New("task failed")
			}
		case <-tickerTimeout:
			return errors.New("polling task timeout")
		}
	}
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
	// "slas":[{"retention_duration":<null>,"rpo_frequency":[{"offsets":<null>,
	// 	"unit":"days","value":1}]}],"type":"aws_ebs_volume_backup"})]
	// .slas[Value({"retention_duration":<null>,"rpo_frequency":[{"offsets":<null>,
	// 	"unit":"days","value":1}]})].retention_duration
	// "retention_duration" is extracted from it as the lowest level field name
	return req.Path.String()[strings.LastIndex(req.Path.String(), ".")+1:]
}

// PollForProtectionGroup polls till the protection group becomes available after create or update
// protection group as they are asynchronous operations.
func PollForProtectionGroup(
	ctx context.Context, id string, protectionGroup sdkclients.ProtectionGroupClient,
	timeout time.Duration, interval time.Duration) (*models.ReadProtectionGroupResponse, error) {

	ticker := time.NewTicker(interval)
	tickerTimeout := time.After(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("context canceled or timed out")
		case <-ticker.C:
			readResponse, err := protectionGroup.ReadProtectionGroup(id)
			if err != nil {
				if err.ResponseCode != http.StatusNotFound {
					return nil, errors.New(ParseMessageFromApiError(err))
				}
				continue
			}
			return readResponse, nil
		case <-tickerTimeout:
			return nil, errors.New("polling timed out")
		}
	}
}

// GetSDKConfigForOU returns a copy of the given SDK config with the OrganizationalUnitContext set
// to the specified organizationalUnitId.
func GetSDKConfigForOU(clumioConfig sdkconfig.Config, organizationalUnitId string) sdkconfig.Config {

	return sdkconfig.Config{
		Token:                     clumioConfig.Token,
		BaseUrl:                   clumioConfig.BaseUrl,
		OrganizationalUnitContext: organizationalUnitId,
		CustomHeaders:             clumioConfig.CustomHeaders,
	}
}
