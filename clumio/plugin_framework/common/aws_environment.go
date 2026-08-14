// Copyright 2026. Clumio, Inc.

// This file holds a shared helper to look up the Clumio AWS environment associated with an AWS
// account and region. It is used by both the AWS connection resource and the AWS environment data
// source.

package common

import (
	"fmt"

	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/clumio-code/clumio-go-sdk/models"
)

// LookupAWSEnvironment retrieves the single Clumio AWS environment associated with the given AWS
// account and region. It returns an error if the environment cannot be uniquely resolved. NOTE: an
// AWS connection only gets associated with an environment in the backend once it becomes connected,
// so lookups for a non-connected account may fail.
func LookupAWSEnvironment(client sdkclients.AWSEnvironmentClient,
	accountNativeId string, awsRegion string) (*models.AWSEnvironment, error) {
	filterStr := fmt.Sprintf(
		"{\"account_native_id\":{\"$eq\":%v}, \"aws_region\":{\"$eq\":%v}}",
		JSONEscapeFilterValue(accountNativeId), JSONEscapeFilterValue(awsRegion))

	// Request up to two environments so the more-than-one-match check below can actually detect
	// ambiguity; an account and region combination should resolve to exactly one environment.
	limit := int64(2)
	envs, apiErr := client.ListAwsEnvironments(&limit, nil, &filterStr, nil, nil)
	if apiErr != nil {
		return nil, fmt.Errorf(
			"unable to retrieve environment corresponding to %v, %v (%v)",
			accountNativeId, awsRegion, ParseMessageFromApiError(apiErr))
	}
	if envs == nil {
		return nil, fmt.Errorf(
			"unable to retrieve environment corresponding to %v, %v, but received nil response",
			accountNativeId, awsRegion)
	}
	if envs.Embedded == nil || len(envs.Embedded.Items) == 0 {
		return nil, fmt.Errorf(
			"unable to retrieve environment corresponding to %v, %v, but no API error was returned",
			accountNativeId, awsRegion)
	}
	if len(envs.Embedded.Items) > 1 {
		count := int64(len(envs.Embedded.Items))
		if envs.CurrentCount != nil {
			count = *envs.CurrentCount
		}
		return nil, fmt.Errorf(
			"expected only one environment corresponding to %v, %v, but found %v",
			accountNativeId, awsRegion, count)
	}
	if envs.Embedded.Items[0].Id == nil || *envs.Embedded.Items[0].Id == "" {
		return nil, fmt.Errorf(
			"environment corresponding to %v, %v has no ID",
			accountNativeId, awsRegion)
	}

	return envs.Embedded.Items[0], nil
}
