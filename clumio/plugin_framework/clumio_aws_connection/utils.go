// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_aws_connection Terraform resource.

package clumio_aws_connection

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// getEnvironmentForConnection returns the environment associated with the given AWS connection.
// NOTE: An AWS connection only gets associated with an environment in the backend once it becomes
// connected. As such, attempts to retrieve the environment for a non-connected AWS connection may
// fail.
func getEnvironmentForConnection(_ context.Context, r *clumioAWSConnectionResource,
	state *clumioAWSConnectionResourceModel) (*models.AWSEnvironment, error) {

	// Construct the filter string required to retrieve the environment using the state of the AWS
	// connection.
	accountNativeId := state.AccountNativeID.ValueString()
	awsRegion := state.AWSRegion.ValueString()
	filterStr := fmt.Sprintf(
		"{\"account_native_id\":{\"$eq\":\"%v\"}, \"aws_region\":{\"$eq\":\"%v\"}}",
		accountNativeId, awsRegion)

	// Call the Clumio API to retrieve the associated environment.
	limit := int64(1)
	envs, apiErr := r.sdkEnvironments.ListAwsEnvironments(&limit, nil, &filterStr, nil, nil)
	if apiErr != nil {
		return nil, fmt.Errorf(
			"unable to retrieve environment corresponding to %v, %v (%v)",
			accountNativeId, awsRegion, common.ParseMessageFromApiError(apiErr))
	}
	if envs.Embedded == nil || len(envs.Embedded.Items) == 0 {
		return nil, fmt.Errorf(
			"unable to retrieve environment corresponding to %v, %v, but no API error was returned",
			accountNativeId, awsRegion)
	}
	if len(envs.Embedded.Items) > 1 {
		return nil, fmt.Errorf(
			"expected only one environment corresponding to %v, %v, but found %v",
			accountNativeId, awsRegion, *envs.CurrentCount)
	}

	// Return the environment.
	return envs.Embedded.Items[0], nil
}

// setExternalId checks and sets the ExternalID in the given state.
func setExternalId(state *clumioAWSConnectionResourceModel, externalId *string, token *string) {
	if externalId != nil && *externalId != "" {
		state.ExternalID = types.StringPointerValue(externalId)
	} else {
		state.ExternalID = types.StringValue(fmt.Sprintf(externalIDFmt, *token))
	}
}

// setDataPlaneAccountId checks and sets the DataPlaneAccountID in the given state.
func setDataPlaneAccountId(state *clumioAWSConnectionResourceModel, dataPlaneAccountId *string) {
	if dataPlaneAccountId != nil && *dataPlaneAccountId != "" {
		state.DataPlaneAccountID = types.StringPointerValue(dataPlaneAccountId)
	} else {
		state.DataPlaneAccountID = types.StringValue(defaultDataPlaneAccountId)
	}
}
