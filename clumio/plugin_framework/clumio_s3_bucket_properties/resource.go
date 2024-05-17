// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the S3 bucket SDK APIs to perform CRUD operations and set the
// attributes from the response of the API in the resource model.

package clumio_s3_bucket_properties

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createOrUpdateS3BucketProperties invokes the API to set the S3 bucket properties.
func (r *clumioS3BucketPropertiesResource) createOrUpdateS3BucketProperties(
	ctx context.Context, plan *clumioS3BucketPropertiesResourceModel, setId bool) diag.Diagnostics {

	var diags diag.Diagnostics

	// Convert the schema to a Clumio API request to update the S3 bucket properties.
	apiReq := &models.SetBucketPropertiesV1Request{
		EventBridgeEnabled:              plan.EventBridgeEnabled.ValueBoolPointer(),
		EventBridgeNotificationDisabled: plan.EvendBridgeNotificationDisabled.ValueBoolPointer(),
	}
	bucketId := plan.BucketID.ValueString()
	_, apiErr := r.sdkS3BucketClient.SetBucketProperties(bucketId, apiReq)
	if apiErr != nil {
		summary := fmt.Sprintf(setErrorFmt, r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	err := r.pollForS3Bucket(ctx, bucketId, plan, r.pollInterval, r.pollTimeout)
	if err != nil {
		summary := fmt.Sprintf(readS3BucketErrorFmt, bucketId)
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	if setId {
		// Since the API doesn't return an id, we are setting the bucket_id as the resource id.
		plan.ID = types.StringValue(bucketId)
	}

	return diags
}

// createOrUpdateS3BucketProperties invokes the API to read the S3 bucket properties.
func (r *clumioS3BucketPropertiesResource) readS3BucketProperties(
	ctx context.Context, state *clumioS3BucketPropertiesResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	bucketId := state.BucketID.ValueString()
	readResponse, apiErr := r.sdkS3BucketClient.ReadAwsS3Bucket(bucketId)
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf(
				"%s (ID: %v) not found. Removing from state", r.name, state.ID)
			tflog.Warn(ctx, summary)
			remove = true
		} else {
			summary := fmt.Sprintf(readS3BucketErrorFmt, bucketId)
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return remove, diags
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return false, diags
	}

	state.EventBridgeEnabled = basetypes.NewBoolPointerValue(readResponse.EventBridgeEnabled)
	return false, diags
}

// deleteS3BucketProperties invokes the API to set the S3 bucket property event_bridge_enabled to
// false.
func (r *clumioS3BucketPropertiesResource) deleteS3BucketProperties(
	_ context.Context, state *clumioS3BucketPropertiesResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// If event bridge is already disabled, there is no need to call the set bucket properties
	// with false.
	if state.EventBridgeEnabled.ValueBool() {
		enabled := false
		// Convert the schema to a Clumio API request to update the S3 bucket properties.
		apiReq := &models.SetBucketPropertiesV1Request{
			EventBridgeEnabled:              &enabled,
			EventBridgeNotificationDisabled: state.EvendBridgeNotificationDisabled.ValueBoolPointer(),
		}
		bucketId := state.BucketID.ValueString()
		_, apiErr := r.sdkS3BucketClient.SetBucketProperties(bucketId, apiReq)
		if apiErr != nil {
			if apiErr.ResponseCode == http.StatusNotFound {
				return diags
			}
			summary := fmt.Sprintf(setErrorFmt, r.name)
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return diags
		}
	}
	return diags
}
