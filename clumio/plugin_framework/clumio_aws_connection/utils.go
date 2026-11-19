// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_aws_connection Terraform resource.

package clumio_aws_connection

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var newAWSEnvironmentClient = sdkclients.NewAWSEnvironmentClient
var newOrganizationalUnitClient = sdkclients.NewOrganizationalUnitClient
var newTaskClient = sdkclients.NewTaskClient

// setOrganizationalUnitID normalizes the stored OU to a concrete known value.
func setOrganizationalUnitID(state *clumioAWSConnectionResourceModel, organizationalUnitID *string) {
	if organizationalUnitID != nil && *organizationalUnitID != "" {
		state.OrganizationalUnitID = types.StringPointerValue(organizationalUnitID)
	} else {
		state.OrganizationalUnitID = types.StringValue(defaultOrgUnitId)
	}
}

// updateOrgUnitForConnection moves the AWS connection between OUs by patching entity membership.
// The backend accepts an add only when the entity already sits in the target OU's parent, and a
// remove lifts the entity to the OU's parent; a repeated add or remove is a no-op. Each PATCH
// therefore moves the entity one level, so the move walks the entity up to the lowest common
// ancestor of the current and target OUs and then down to the target, one PATCH per level.
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

	// An OU-scoped API key can read only its own OU and the OUs below it, so the walk never reads
	// an OU above the common ancestor. Both walks finish before the first PATCH.
	targetID := plan.OrganizationalUnitID.ValueString()
	removeFrom, commonID, err := getOrgUnitsUpToCommonAncestor(
		orgUnitClient, state.OrganizationalUnitID.ValueString(), targetID)
	if err != nil {
		return err
	}
	addTo, err := getOrgUnitsBelowAncestor(orgUnitClient, targetID, commonID)
	if err != nil {
		return err
	}

	for _, orgUnitID := range removeFrom {
		err = patchOrgUnitEntities(ctx, r, orgUnitClient, taskClient, orgUnitID,
			&models.UpdateEntities{Remove: entityModels})
		if err != nil {
			return err
		}
	}
	for _, orgUnitID := range slices.Backward(addTo) {
		err = patchOrgUnitEntities(ctx, r, orgUnitClient, taskClient, orgUnitID,
			&models.UpdateEntities{Add: entityModels})
		if err != nil {
			return err
		}
	}
	return nil
}

// getOrgUnitsUpToCommonAncestor climbs from the current OU to the first OU that is the target or
// contains it. It returns the OUs climbed out of, bottom-up, and that common ancestor.
func getOrgUnitsUpToCommonAncestor(client sdkclients.OrganizationalUnitClient,
	currentID string, targetID string) ([]string, string, error) {

	var climbed []string
	for currentID != targetID && currentID != defaultOrgUnitId {
		orgUnit, err := getOrgUnitForConnection(client, currentID)
		if err != nil {
			return nil, "", err
		}
		if slices.ContainsFunc(orgUnit.DescendantIds, func(id *string) bool {
			return id != nil && *id == targetID
		}) {
			break
		}
		climbed = append(climbed, currentID)
		currentID = getParentOrgUnitID(orgUnit)
		if slices.Contains(climbed, currentID) {
			return nil, "", fmt.Errorf("organizational unit %v is its own ancestor", currentID)
		}
	}
	return climbed, currentID, nil
}

// getOrgUnitsBelowAncestor climbs from the target OU up to the given ancestor and returns the
// OUs passed, bottom-up, excluding the ancestor.
func getOrgUnitsBelowAncestor(client sdkclients.OrganizationalUnitClient,
	targetID string, ancestorID string) ([]string, error) {

	var climbed []string
	for orgUnitID := targetID; orgUnitID != ancestorID; {
		if orgUnitID == defaultOrgUnitId || slices.Contains(climbed, orgUnitID) {
			return nil, fmt.Errorf("organizational unit %v is not below organizational unit %v",
				targetID, ancestorID)
		}
		orgUnit, err := getOrgUnitForConnection(client, orgUnitID)
		if err != nil {
			return nil, err
		}
		climbed = append(climbed, orgUnitID)
		orgUnitID = getParentOrgUnitID(orgUnit)
	}
	return climbed, nil
}

// getParentOrgUnitID returns the parent of the given OU, treating a missing parent as Global.
func getParentOrgUnitID(orgUnit *models.ReadOrganizationalUnitResponse) string {
	if orgUnit.ParentId == nil || *orgUnit.ParentId == "" {
		return defaultOrgUnitId
	}
	return *orgUnit.ParentId
}

// patchOrgUnitEntities applies one entity add or remove to the given OU and waits for its task.
func patchOrgUnitEntities(ctx context.Context, r *clumioAWSConnectionResource,
	orgUnitClient sdkclients.OrganizationalUnitClient, taskClient sdkclients.TaskClient,
	organizationalUnitID string, updateEntities *models.UpdateEntities) error {

	action := fmt.Sprintf("add the connection to organizational unit %v", organizationalUnitID)
	if updateEntities.Remove != nil {
		action = fmt.Sprintf("remove the connection from organizational unit %v", organizationalUnitID)
	}
	ouUpdateRequest := &models.PatchOrganizationalUnitV2Request{
		Entities: updateEntities,
	}
	res, apiErr := orgUnitClient.PatchOrganizationalUnit(organizationalUnitID, nil, ouUpdateRequest)
	if apiErr != nil {
		return fmt.Errorf("unable to %v (%v)", action, common.ParseMessageFromApiError(apiErr))
	}
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unable to %v (HTTP status code: %v)", action, res.StatusCode)
	}
	if res.Http202 == nil || res.Http202.TaskId == nil {
		return fmt.Errorf("unable to %v (no task ID)", action)
	}

	err := common.PollTask(ctx, taskClient, *res.Http202.TaskId, r.pollTimeout, r.pollInterval)
	if err != nil {
		return fmt.Errorf("unable to %v (%w)", action, err)
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
			"unable to retrieve Organizational Unit %v (%v)",
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
	environment, err := common.LookupAWSEnvironment(r.sdkEnvironments, accountNativeId, awsRegion)
	if err == nil {
		return environment, nil
	}

	// Environment visibility can lag or be scoped differently across OU contexts. Retry from the
	// Global OU context before failing the move.
	if r.client != nil && r.client.ClumioConfig.OrganizationalUnitContext != "" {
		globalConfig := common.GetSDKConfigForOU(r.client.ClumioConfig, "")
		globalEnvClient := newAWSEnvironmentClient(globalConfig)
		environment, globalErr := common.LookupAWSEnvironment(
			globalEnvClient, accountNativeId, awsRegion)
		if globalErr == nil {
			return environment, nil
		}
	}

	return nil, err
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
