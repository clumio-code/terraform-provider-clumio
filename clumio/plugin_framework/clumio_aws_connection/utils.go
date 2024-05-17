// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_aws_connection Terraform resource.

package clumio_aws_connection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// updateOrgUnitForConnection updates the Organizational Unit (OU) for the AWS connection if the new
// OU provided is either the parent of the current OU or one of its descendants.
func updateOrgUnitForConnection(
	ctx context.Context, r *clumioAWSConnectionResource, plan *clumioAWSConnectionResourceModel,
	state *clumioAWSConnectionResourceModel) error {

	// Retrieve the environment associated with the connection. The environment ID is required to
	// initiate the update of the Organizational Unit (OU) for the AWS connection. Once retrieved,
	// the environment ID is used to create an EntityModel that represents the connection. This
	// EntityModel is then used to update the OU for the connection later.
	environment, err := getEnvironmentForConnection(ctx, r, state)
	if err != nil {
		return err
	}
	entityId := *environment.Id
	entityType := awsEnvironment
	entityModels := []*models.EntityModel{
		{
			PrimaryEntity: &models.OrganizationalUnitPrimaryEntity{
				Id:         &entityId,
				ClumioType: &entityType,
			},
		},
	}

	// If the current OU is null and the new OU is the default OU or vice versa, there is no need
	// to update the OU.
	orgUnitId := plan.OrganizationalUnitID.ValueString()
	currentOrgUnitId := state.OrganizationalUnitID.ValueString()
	if orgUnitId == "" && currentOrgUnitId == defaultOrgUnitId {
		return nil
	} else if currentOrgUnitId == "" && orgUnitId == defaultOrgUnitId {
		return nil
	}

	var updateEntities *models.UpdateEntities

	if currentOrgUnitId == "" {
		// Current OU is null. So the connection needs to be added to the new OU.
		updateEntities = &models.UpdateEntities{Add: entityModels}
	} else {
		// Retrieve the current Organizational Unit (OU) associated with the connection. Once retrieved,
		// the OU's parent ID is compared with the new OU provided in the plan. If the new OU is the
		// parent of the current OU, the AWS connection is removed from the current OU. Else, the AWS
		// connection is added to the new OU.
		currentOrgUnit, err := getOrgUnitForConnection(ctx, r, state)
		if err != nil {
			return err
		}
		if currentOrgUnit.ParentId != nil && *currentOrgUnit.ParentId == orgUnitId || orgUnitId == "" {
			// As the new OU is the parent of the current OU, the AWS connection should be removed from
			// the current OU. Thus, the "orgUnitId" is set to the current OU's ID in order to remove
			// the connection from it.
			orgUnitId = *currentOrgUnit.Id
			updateEntities = &models.UpdateEntities{Remove: entityModels}
		} else {
			updateEntities = &models.UpdateEntities{Add: entityModels}
		}
	}

	// Call the Clumio API to update the Organizational Unit (OU) for the AWS connection.
	ouUpdateRequest := &models.PatchOrganizationalUnitV2Request{
		Entities: updateEntities,
	}
	res, apiErr := r.sdkOrgUnits.PatchOrganizationalUnit(orgUnitId, nil, ouUpdateRequest)
	if apiErr != nil {
		return fmt.Errorf(
			"unable to update the Organizational Unit for the connection (%v)",
			common.ParseMessageFromApiError(apiErr))
	}
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"unable to update the Organizational Unit for the connection (HTTP status code: %v)",
			res.StatusCode)
	}
	if res.Http202 == nil {
		return fmt.Errorf("unable to update the Organizational Unit for the connection (no task ID)")
	}

	// As the modification of the OU for a connection is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err = common.PollTask(ctx, r.sdkTasks, *res.Http202.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		return fmt.Errorf("unable to update the Organizational Unit for the connection (%v)", err)
	}

	return nil
}

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
	envs, apiErr := r.sdkEnvironments.ListAwsEnvironments(&limit, nil, &filterStr, nil)
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

// getOrgUnitForConnection returns the Organizational Unit (OU) associated with the given AWS
// connection.
func getOrgUnitForConnection(_ context.Context, r *clumioAWSConnectionResource,
	model *clumioAWSConnectionResourceModel) (*models.ReadOrganizationalUnitResponse, error) {

	// Retrieve the current Organizational Unit ID associated with the AWS connection.
	orgUnitId := model.OrganizationalUnitID.ValueString()

	// Call the Clumio API to retrieve the associated Organizational Unit.
	orgUnit, apiErr := r.sdkOrgUnits.ReadOrganizationalUnit(orgUnitId, nil)
	if apiErr != nil {
		return nil, fmt.Errorf(
			"unable to retrieve Organizational Unit %v (%v)",
			orgUnitId, common.ParseMessageFromApiError(apiErr))
	}
	return orgUnit, nil
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
