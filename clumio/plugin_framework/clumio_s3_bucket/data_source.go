// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio S3 Bucket SDK API to perform read
// operation and set the attributes from the response of the API in the data source model.

package clumio_s3_bucket

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"strings"
)

// readS3Bucket invokes the API to read the s3BucketClient and from the response
// populates the attributes of the s3 bucket.
func (r *clumioS3BucketDataSource) readS3Bucket(
	ctx context.Context, model *clumioS3BucketDataSourceModel) diag.Diagnostics {

	// Prepare the query nameFilter.
	var bucketNames []string
	diags := model.BucketNames.ElementsAs(ctx, &bucketNames, true)
	if diags.HasError() {
		return diags
	}
	nameFilter := fmt.Sprintf(`{"name": {"$in":["%s"]}}`, strings.Join(bucketNames, "\", \""))

	// Call the Clumio API to list the s3 buckets.
	limit := int64(10000)
	res, apiErr := r.s3BucketClient.ListAwsS3Buckets(&limit, nil, &nameFilter)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response for the s3 buckets into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateS3BucketsInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "S3 bucket not found."
		detail := fmt.Sprintf(
			"No S3 bucket found with the given bucket names %v", bucketNames)
		diags.AddError(summary, detail)
	}
	return diags
}
