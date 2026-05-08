// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_aws_connection Terraform resource.

package clumio_aws_connection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var newAWSEnvironmentClient = sdkclients.NewAWSEnvironmentClient
var newOrganizationalUnitClient = sdkclients.NewOrganizationalUnitClient
var newTaskClient = sdkclients.NewTaskClient

// getDesiredOrganizationalUnitID returns the OU implied by the configured provider context.
func getDesiredOrganizationalUnitID(client *common.ApiClient) string {
	if client == nil || client.ClumioConfig.OrganizationalUnitContext == "" {
		return defaultOrgUnitId
	}
	return client.ClumioConfig.OrganizationalUnitContext
}

// setOrganizationalUnitID normalizes the stored OU to a concrete known value.
func setOrganizationalUnitID(state *clumioAWSConnectionResourceModel, organizationalUnitID *string) {
	if organizationalUnitID != nil && *organizationalUnitID != "" {
		state.OrganizationalUnitID = types.StringPointerValue(organizationalUnitID)
	} else {
		state.OrganizationalUnitID = types.StringValue(defaultOrgUnitId)
	}
}

// updateOrgUnitForConnection moves the AWS connection between OUs by patching entity membership.
func updateOrgUnitForConnection(
	ctx context.Context, r *clumioAWSConnectionResource, plan *clumioAWSConnectionResourceModel,
	state *clumioAWSConnectionResourceModel) error {

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

	orgUnitClient := r.sdkOrgUnits
	taskClient := r.sdkTasks
	if r.client != nil && r.client.ClumioConfig.OrganizationalUnitContext != "" {
		globalConfig := common.GetSDKConfigForOU(r.client.ClumioConfig, "")
		orgUnitClient = newOrganizationalUnitClient(globalConfig)
		taskClient = newTaskClient(globalConfig)
	}

	currentOrgUnitID := state.OrganizationalUnitID.ValueString()
	orgUnitID := plan.OrganizationalUnitID.ValueString()
	var updateEntities *models.UpdateEntities
	if currentOrgUnitID == defaultOrgUnitId {
		updateEntities = &models.UpdateEntities{Add: entityModels}
	} else if orgUnitID == defaultOrgUnitId {
		orgUnitID = currentOrgUnitID
		updateEntities = &models.UpdateEntities{Remove: entityModels}
	} else {
		currentOrgUnit, err := getOrgUnitForConnection(orgUnitClient, currentOrgUnitID)
		if err != nil {
			return err
		}
		if currentOrgUnit.ParentId != nil && *currentOrgUnit.ParentId == orgUnitID {
			orgUnitID = currentOrgUnitID
			updateEntities = &models.UpdateEntities{Remove: entityModels}
		} else {
			updateEntities = &models.UpdateEntities{Add: entityModels}
		}
	}

	ouUpdateRequest := &models.PatchOrganizationalUnitV2Request{
		Entities: updateEntities,
	}
	res, apiErr := orgUnitClient.PatchOrganizationalUnit(orgUnitID, nil, ouUpdateRequest)
	if apiErr != nil {
		return fmt.Errorf(
			"Unable to update the Organizational Unit for the connection (%v)",
			common.ParseMessageFromApiError(apiErr))
	}
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"Unable to update the Organizational Unit for the connection (HTTP status code: %v)",
			res.StatusCode)
	}
	if res.Http202 == nil || res.Http202.TaskId == nil {
		return fmt.Errorf("Unable to update the Organizational Unit for the connection (no task ID)")
	}

	err = common.PollTask(ctx, taskClient, *res.Http202.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		return fmt.Errorf("Unable to update the Organizational Unit for the connection (%v)", err)
	}

	return nil
}

// getOrgUnitForConnection returns the organizational unit for the given ID.
func getOrgUnitForConnection(
	client sdkclients.OrganizationalUnitClient, organizationalUnitID string) (
	*models.ReadOrganizationalUnitResponse, error) {

	orgUnit, apiErr := client.ReadOrganizationalUnit(organizationalUnitID, nil)
	if apiErr != nil {
		return nil, fmt.Errorf(
			"Unable to retrieve Organizational Unit %v (%v)",
			organizationalUnitID, common.ParseMessageFromApiError(apiErr))
	}
	return orgUnit, nil
}

// getEnvironmentForConnection returns the environment associated with the given AWS connection.
// NOTE: An AWS connection only gets associated with an environment in the backend once it becomes
// connected. As such, attempts to retrieve the environment for a non-connected AWS connection may
// fail.
func getEnvironmentForConnection(_ context.Context, r *clumioAWSConnectionResource,
	state *clumioAWSConnectionResourceModel) (*models.AWSEnvironment, error) {

	accountNativeId := state.AccountNativeID.ValueString()
	awsRegion := state.AWSRegion.ValueString()
	environment, err := lookupEnvironmentForConnection(r.sdkEnvironments, accountNativeId, awsRegion)
	if err == nil {
		return environment, nil
	}

	// Environment visibility can lag or be scoped differently across OU contexts. Retry from the
	// Global OU context before failing the move.
	if r.client != nil && r.client.ClumioConfig.OrganizationalUnitContext != "" {
		globalConfig := common.GetSDKConfigForOU(r.client.ClumioConfig, "")
		globalEnvClient := newAWSEnvironmentClient(globalConfig)
		environment, globalErr := lookupEnvironmentForConnection(
			globalEnvClient, accountNativeId, awsRegion)
		if globalErr == nil {
			return environment, nil
		}
	}

	return nil, err
}

func lookupEnvironmentForConnection(client sdkclients.AWSEnvironmentClient,
	accountNativeId string, awsRegion string) (*models.AWSEnvironment, error) {
	filterStr := fmt.Sprintf(
		"{\"account_native_id\":{\"$eq\":\"%v\"}, \"aws_region\":{\"$eq\":\"%v\"}}",
		accountNativeId, awsRegion)

	limit := int64(1)
	envs, apiErr := client.ListAwsEnvironments(&limit, nil, &filterStr, nil, nil)
	if apiErr != nil {
		return nil, fmt.Errorf(
			"Unable to retrieve environment corresponding to %v, %v (%v)",
			accountNativeId, awsRegion, common.ParseMessageFromApiError(apiErr))
	}
	if envs == nil {
		return nil, fmt.Errorf(
			"Unable to retrieve environment corresponding to %v, %v, but received nil response",
			accountNativeId, awsRegion)
	}
	if envs.Embedded == nil || len(envs.Embedded.Items) == 0 {
		return nil, fmt.Errorf(
			"Unable to retrieve environment corresponding to %v, %v, but no API error was returned",
			accountNativeId, awsRegion)
	}
	if len(envs.Embedded.Items) > 1 {
		count := int64(len(envs.Embedded.Items))
		if envs.CurrentCount != nil {
			count = *envs.CurrentCount
		}
		return nil, fmt.Errorf(
			"Expected only one environment corresponding to %v, %v, but found %v",
			accountNativeId, awsRegion, count)
	}
	if envs.Embedded.Items[0].Id == nil || *envs.Embedded.Items[0].Id == "" {
		return nil, fmt.Errorf(
			"Environment corresponding to %v, %v has no ID",
			accountNativeId, awsRegion)
	}

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
