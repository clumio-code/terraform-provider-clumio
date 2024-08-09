// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the protection group SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_protection_group

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createProtectionGroup invokes the API to create the protection group and from the response
// populates the computed attributes of the protection group.
func (r *clumioProtectionGroupResource) createProtectionGroup(
	ctx context.Context, plan *clumioProtectionGroupResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK client is temporarily re-initialized in the context of the
	// desired OU so that API calls are made on behalf of the OU.
	if plan.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, plan.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
	}

	// Call the Clumio API to create the protection group.
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(plan.ObjectFilter)
	response, apiErr := sdkProtectionGroups.CreateProtectionGroup(
		models.CreateProtectionGroupV1Request{
			BucketRule:   plan.BucketRule.ValueStringPointer(),
			Description:  plan.Description.ValueStringPointer(),
			Name:         plan.Name.ValueStringPointer(),
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if response == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Poll to read the protection group till it becomes available
	readResponse, err := common.PollForProtectionGroup(
		ctx, *response.Id, sdkProtectionGroups, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll %s (ID: %v) for creation",
			r.name, plan.Name.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.ID = types.StringPointerValue(response.Id)
	plan.Name = types.StringPointerValue(readResponse.Name)
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	plan.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(readResponse.ObjectFilter)
	plan.ProtectionStatus = types.StringPointerValue(readResponse.ProtectionStatus)
	plan.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	return diags
}

// readProtectionGroup invokes the API to read the protection group and from the response populates
// the attributes of the protection group. If the protection group has been removed externally, the
// function returns "true" to indicate to the caller that the resource no longer exists.
func (r *clumioProtectionGroupResource) readProtectionGroup(
	ctx context.Context, state *clumioProtectionGroupResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK client is temporarily re-initialized in the context of the
	// desired OU so that API calls are made on behalf of the OU.
	if state.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, state.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
	}

	// Call the Clumio API to read the protection group
	readResponse, apiErr := sdkProtectionGroups.ReadProtectionGroup(state.ID.ValueString(), nil)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state",
				r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			return true, diags
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return false, diags
		}
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return false, diags
	}

	// If the protection group was deleted externally, issue a warning and return "true" to signal
	// to the caller that the resource has been removed.
	if readResponse.IsDeleted != nil && *readResponse.IsDeleted {
		msgStr := fmt.Sprintf(
			"Clumio Protection Group with ID %s not found. Removing from state.",
			state.ID.ValueString())
		tflog.Warn(ctx, msgStr)
		return true, diags
	}

	// Convert the Clumio API response back to a schema and update the state. An optional field need
	// only be populated if it initially contained a non-null value or if there is a specific value
	// that should be assigned.
	description := types.StringPointerValue(readResponse.Description)
	bucketRule := types.StringPointerValue(readResponse.BucketRule)
	if !state.Description.IsNull() || description.ValueString() != "" {
		state.Description = description
	}
	if !state.BucketRule.IsNull() || bucketRule.ValueString() != "" {
		state.BucketRule = bucketRule
	}
	state.Name = types.StringPointerValue(readResponse.Name)
	state.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	state.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(readResponse.ObjectFilter)
	state.ProtectionStatus = types.StringPointerValue(readResponse.ProtectionStatus)
	state.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	return false, diags
}

// updateProtectionGroup invokes the API to update the protection group and from the response
// populates the computed attributes of the protection group.
func (r *clumioProtectionGroupResource) updateProtectionGroup(
	ctx context.Context, plan *clumioProtectionGroupResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups

	// Call the Clumio API to get the current version of the protection group.
	readResp, apiErr := sdkProtectionGroups.ReadProtectionGroup(plan.ID.ValueString(), nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	version := readResp.Version

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK client is temporarily re-initialized in the context of the
	// desired OU so that API calls are made on behalf of the OU.
	if plan.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, plan.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
	}

	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(plan.ObjectFilter)

	// Call the Clumio API to update the protection group.
	updateReq := &models.UpdateProtectionGroupV1Request{
		BucketRule:   plan.BucketRule.ValueStringPointer(),
		Description:  plan.Description.ValueStringPointer(),
		Name:         plan.Name.ValueStringPointer(),
		ObjectFilter: objectFilter,
	}
	response, apiErr := sdkProtectionGroups.UpdateProtectionGroup(plan.ID.ValueString(),
		updateReq)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if response == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Poll to read the protection group till it is updated
	readResponse, err := common.PollForProtectionGroupUpdate(
		ctx, *response.Id, version, updateReq, sdkProtectionGroups, r.pollTimeout, r.pollInterval)
	if err != nil {
		summary := fmt.Sprintf(
			"Unable to poll %s (ID: %v) for update", r.name, plan.ID.ValueString())
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}
	// Convert the Clumio API response back to a schema and populate all computed fields of the plan.
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	plan.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(readResponse.ObjectFilter)
	plan.ProtectionStatus = types.StringPointerValue(readResponse.ProtectionStatus)
	plan.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	return diags
}

// deleteProtectionGroup invokes the API to delete the protection group
func (r *clumioProtectionGroupResource) deleteProtectionGroup(
	_ context.Context, state *clumioProtectionGroupResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups

	// If the OrganizationalUnitID is specified, then execute the API in that Organizational Unit
	// (OU) context. To that end, the SDK client is temporarily re-initialized in the context of the
	// desired OU so that API calls are made on behalf of the OU.
	if state.OrganizationalUnitID.ValueString() != "" {
		config := common.GetSDKConfigForOU(
			r.client.ClumioConfig, state.OrganizationalUnitID.ValueString())
		sdkProtectionGroups = sdkclients.NewProtectionGroupClient(config)
	}

	// Call the Clumio API to delete the protection group
	_, apiErr := sdkProtectionGroups.DeleteProtectionGroup(state.ID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}
	return diags
}
