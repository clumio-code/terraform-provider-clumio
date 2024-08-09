// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio S3 Bucket SDK API to perform read
// operation and set the attributes from the response of the API in the data source model.

package clumio_protection_group_asset

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"strings"
)

// readProtectionGroupAsset invokes the API to read the s3BucketClient and from the response
// populates the attributes of the s3 bucket.
func (r *clumioProtectionGroupAssetDataSource) readProtectionGroupAsset(
	ctx context.Context, model *clumioProtectionGroupAssetDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	filters := make([]string, 0)
	pgId := model.ProtectionGroupID.ValueString()
	pgFilter := fmt.Sprintf(`"protection_group_id": {"$eq":"%s"}`, pgId)
	filters = append(filters, pgFilter)
	bucketId := model.BucketID.ValueString()
	bucketFilter := fmt.Sprintf(`"bucket_id": {"$eq":"%s"}`, bucketId)
	filters = append(filters, bucketFilter)
	filter := fmt.Sprintf("{%s}", strings.Join(filters, ","))

	// Call the Clumio API to list the S3 Assets for the protection group.
	readResponse, apiErr := r.s3AssetsClient.ListProtectionGroupS3Assets(
		nil, nil, &filter, &common.DefaultLookBackDays)
	if apiErr != nil {
		summary := fmt.Sprintf(
			"Unable to read protection group asset with Protection Group ID: %v and "+
				"Bucket ID: %v", pgId, bucketId)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if readResponse == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	if readResponse.TotalCount == nil || *readResponse.TotalCount == 0 {
		summary := "Protection group asset not found."
		detail := fmt.Sprintf(
			"Expected one asset with Bucket ID %s and Protection Group ID %s.",
			bucketId, pgId)
		diags.AddError(summary, detail)
		return diags
	}

	model.Id = basetypes.NewStringPointerValue(readResponse.Embedded.Items[0].Id)

	return diags
}
