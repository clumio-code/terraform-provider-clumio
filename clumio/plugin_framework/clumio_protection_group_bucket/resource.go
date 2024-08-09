// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the protection group SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_protection_group_bucket

import (
	"context"
	"fmt"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

// createProtectionGroupBucket invokes the API to assign the bucket to the protection group and from
// the response populates the computed attributes of the protection group bucket.
func (r *clumioProtectionGroupBucketResource) createProtectionGroupBucket(
	_ context.Context, plan *clumioProtectionGroupBucketResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups

	// Call the Clumio API to assign the bucket to the protection group.
	response, apiErr := sdkProtectionGroups.AddBucketProtectionGroup(
		plan.ProtectionGroupID.ValueString(),
		models.AddBucketProtectionGroupV1Request{
			BucketId: plan.BucketID.ValueStringPointer(),
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

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.ID = types.StringPointerValue(response.Id)
	return diags
}

// readProtectionGroupBucket invokes the API to list the protection group assets based on the
// protection group ID and bucket ID. If the bucket has been removed externally from the protection
// group, the function returns "true" to indicate to the caller that the resource no longer exists.
func (r *clumioProtectionGroupBucketResource) readProtectionGroupBucket(
	ctx context.Context, state *clumioProtectionGroupBucketResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	readResponse, apiErr := r.sdkS3Assets.ReadProtectionGroupS3Asset(
		state.ID.ValueString(), &common.DefaultLookBackDays)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state",
				r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			return true, diags
		} else {
			summary := fmt.Sprintf(
				"Unable to read bucket with ID: %v", state.BucketID.ValueString())
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

	// If the protection group bucket was deleted externally, issue a warning and return "true" to
	// signal to the caller that the resource has been removed.
	if readResponse.IsDeleted != nil && *readResponse.IsDeleted {
		msgStr := fmt.Sprintf(
			"Bucket with ID %s is not part of Protection Group with ID %s. Removing from state.",
			state.BucketID.ValueString(), state.ProtectionGroupID.ValueString())
		tflog.Warn(ctx, msgStr)
		return true, diags
	}
	state.BucketID = basetypes.NewStringPointerValue(readResponse.BucketId)
	state.ProtectionGroupID = basetypes.NewStringPointerValue(readResponse.GroupId)

	return false, diags
}

// deleteProtectionGroupBucket invokes the API to remove the bucket from the protection group.
func (r *clumioProtectionGroupBucketResource) deleteProtectionGroupBucket(
	_ context.Context, state *clumioProtectionGroupBucketResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkProtectionGroups := r.sdkProtectionGroups

	// Call the Clumio API to delete the protection group
	_, apiErr := sdkProtectionGroups.DeleteBucketProtectionGroup(
		state.ProtectionGroupID.ValueString(), state.BucketID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound &&
		apiErr.ResponseCode != http.StatusConflict {

		summary := fmt.Sprintf(
			"Unable to remove bucket with ID: %s from Protection Group with ID: %s",
			state.BucketID.ValueString(), state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}
	return diags
}
