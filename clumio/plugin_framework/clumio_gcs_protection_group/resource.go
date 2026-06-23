// Copyright 2026. Clumio, Inc.

// This file holds the logic to invoke the GCP protection group SDK APIs to perform CRUD operations
// and set the attributes from the response of the API in the resource model.

package clumio_gcs_protection_group

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createProtectionGroup invokes the API to create the GCP protection group and from the response
// populates the computed attributes of the protection group.
func (r *clumioGCSProtectionGroupResource) createProtectionGroup(
	ctx context.Context, plan *clumioGCSProtectionGroupResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	bucketRule, d := mapSchemaBucketRuleToClumioBucketRule(ctx, plan.BucketRule)
	diags.Append(d...)
	includePrefixes, d := setToStringSlice(ctx, plan.IncludePrefixes)
	diags.Append(d...)
	excludePrefixes, d := setToStringSlice(ctx, plan.ExcludePrefixes)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	// Call the Clumio API to create the protection group.
	response, apiErr := r.sdkProtectionGroups.CreateGcpProtectionGroup(
		models.CreateGcpProtectionGroupV1Request{
			BucketRule:        bucketRule,
			ExcludePrefixes:   excludePrefixes,
			IncludePrefixes:   includePrefixes,
			LatestVersionOnly: plan.LatestVersionOnly.ValueBoolPointer(),
			Name:              plan.Name.ValueStringPointer(),
		})
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if response == nil || response.Id == nil {
		diags.AddError(common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return diags
	}
	plan.ID = types.StringPointerValue(response.Id)

	// Read back the created protection group to populate computed attributes. The user-managed
	// inputs (bucket_rule, prefixes, latest_version_only) are kept from the plan.
	_, readDiags := r.populateFromServer(ctx, plan, false)
	diags.Append(readDiags...)
	return diags
}

// readProtectionGroup invokes the API to read the GCP protection group and populates the model. It
// returns "true" when the protection group has been removed externally.
func (r *clumioGCSProtectionGroupResource) readProtectionGroup(
	ctx context.Context, state *clumioGCSProtectionGroupResourceModel) (bool, diag.Diagnostics) {

	// On read, the bucket_rule is mapped from the server so external drift is detected.
	return r.populateFromServer(ctx, state, true)
}

// updateProtectionGroup invokes the API to update the GCP protection group and populates the
// computed attributes of the protection group.
func (r *clumioGCSProtectionGroupResource) updateProtectionGroup(
	ctx context.Context, plan *clumioGCSProtectionGroupResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	bucketRule, d := mapSchemaBucketRuleToClumioBucketRule(ctx, plan.BucketRule)
	diags.Append(d...)
	includePrefixes, d := setToStringSlice(ctx, plan.IncludePrefixes)
	diags.Append(d...)
	excludePrefixes, d := setToStringSlice(ctx, plan.ExcludePrefixes)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	updateReq := models.UpdateGcpProtectionGroupV1Request{
		BucketRule:        bucketRule,
		ExcludePrefixes:   excludePrefixes,
		IncludePrefixes:   includePrefixes,
		LatestVersionOnly: plan.LatestVersionOnly.ValueBoolPointer(),
		Name:              plan.Name.ValueStringPointer(),
	}
	// bucket_rule and clear_bucket_rule are mutually exclusive. When the rule is removed from the
	// configuration, instruct the API to clear the existing rule.
	if bucketRule == nil {
		clearBucketRule := true
		updateReq.ClearBucketRule = &clearBucketRule
	}

	_, apiErr := r.sdkProtectionGroups.UpdateGcpProtectionGroup(plan.ID.ValueString(), updateReq)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	// Read back the updated protection group to populate computed attributes.
	_, readDiags := r.populateFromServer(ctx, plan, false)
	diags.Append(readDiags...)
	return diags
}

// deleteProtectionGroup invokes the API to delete the GCP protection group.
func (r *clumioGCSProtectionGroupResource) deleteProtectionGroup(
	_ context.Context, state *clumioGCSProtectionGroupResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	_, apiErr := r.sdkProtectionGroups.DeleteGcpProtectionGroup(state.ID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}
	return diags
}

// populateFromServer reads the protection group and maps server-owned fields into the model. When
// mapBucketRule is true (the Read path) the bucket_rule is also mapped from the server. It returns
// "true" when the protection group no longer exists.
func (r *clumioGCSProtectionGroupResource) populateFromServer(
	ctx context.Context, model *clumioGCSProtectionGroupResourceModel,
	mapBucketRule bool) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	readResponse, apiErr := r.sdkProtectionGroups.ReadGcpProtectionGroup(model.ID.ValueString(), nil)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("%s (ID: %v) not found. Removing from state",
				r.name, model.ID.ValueString()))
			return true, diags
		}
		summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, model.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return false, diags
	}
	if readResponse == nil {
		diags.AddError(common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return false, diags
	}
	if readResponse.IsDeleted != nil && *readResponse.IsDeleted {
		tflog.Warn(ctx, fmt.Sprintf("%s (ID: %v) is deleted. Removing from state",
			r.name, model.ID.ValueString()))
		return true, diags
	}

	diags.Append(populateComputedFromResponse(model, readResponse)...)
	if mapBucketRule {
		bucketRule, d := mapClumioBucketRuleToSchemaBucketRule(ctx, readResponse.BucketRule)
		diags.Append(d...)
		model.BucketRule = bucketRule
	}
	return false, diags
}

// populateComputedFromResponse maps the server-owned (computed) fields from a read response into the
// model. It does not touch the user-managed input fields (bucket_rule, prefixes,
// latest_version_only).
func populateComputedFromResponse(model *clumioGCSProtectionGroupResourceModel,
	resp *models.ReadGCPProtectionGroupResponse) diag.Diagnostics {

	var diags diag.Diagnostics

	model.ID = types.StringPointerValue(resp.Id)
	model.Name = types.StringPointerValue(resp.Name)
	model.OrganizationalUnitID = types.StringPointerValue(resp.OrganizationalUnitId)
	model.ProtectionStatus = types.StringPointerValue(resp.ProtectionStatus)
	model.BucketCount = types.Int64PointerValue(resp.BucketCount)
	model.BucketRuleMatchedCount = types.Int64PointerValue(resp.BucketRuleMatchedBucketCount)
	model.ManualAddedBucketCount = types.Int64PointerValue(resp.ManualAddedBucketCount)
	model.Location = types.StringPointerValue(resp.Location)
	model.LocationType = types.StringPointerValue(resp.LocationType)
	model.LastBackupTimestamp = types.StringPointerValue(resp.LastBackupTimestamp)
	model.TotalBackedUpObjectCount = types.Int64PointerValue(resp.TotalBackedUpObjectCount)
	model.TotalBackedUpSizeBytes = types.Int64PointerValue(resp.TotalBackedUpSizeBytes)
	model.CreatedTimestamp = types.StringPointerValue(resp.CreatedTimestamp)
	model.ModifiedTimestamp = types.StringPointerValue(resp.ModifiedTimestamp)
	model.Version = types.Int64PointerValue(resp.Version)

	bucketUuids, d := stringSliceToSet(resp.BucketUuids)
	diags.Append(d...)
	model.BucketUuids = bucketUuids

	matchedUuids, d := stringSliceToSet(resp.BucketRuleMatchedBucketUuids)
	diags.Append(d...)
	model.BucketRuleMatchedUuids = matchedUuids

	protectionInfo, d := mapProtectionInfo(resp.ProtectionInfo)
	diags.Append(d...)
	model.ProtectionInfo = protectionInfo

	labels, d := mapLabels(resp.Labels)
	diags.Append(d...)
	model.Labels = labels

	backupStatusStats, d := mapBackupStatusStats(resp.BackupStatusStats)
	diags.Append(d...)
	model.BackupStatusStats = backupStatusStats

	return diags
}
