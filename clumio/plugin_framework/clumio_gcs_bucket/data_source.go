// Copyright 2026. Clumio, Inc.

// This file holds the logic to invoke the Clumio GCS Bucket SDK API to perform read operations and
// set the attributes from the response of the API in the data source model.

package clumio_gcs_bucket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readGCSBucket invokes the API to list the GCS buckets matching the given bucket names and from the
// response populates the gcs_buckets attribute of the data source model.
func (r *clumioGCSBucketDataSource) readGCSBucket(
	ctx context.Context, model *clumioGCSBucketDataSourceModel) diag.Diagnostics {

	// Prepare the query name filter.
	var bucketNames []string
	diags := model.BucketNames.ElementsAs(ctx, &bucketNames, true)
	if diags.HasError() {
		return diags
	}
	// JSON-encode the names so the "$in" filter stays valid JSON regardless of their contents.
	namesJSON, err := json.Marshal(bucketNames)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to read %s", r.name),
			fmt.Sprintf("Failed to encode the name filter: %s", err.Error()))
		return diags
	}
	nameFilter := fmt.Sprintf(`{"name": {"$in":%s}}`, string(namesJSON))

	// Call the Clumio API to list the GCS buckets.
	limit := int64(10000)
	res, apiErr := r.gcsBucketClient.ListGcpGcsBuckets(&limit, nil, &nameFilter, nil, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		diags.AddError(common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return diags
	}

	// Convert the Clumio API response for the GCS buckets into the data source schema model.
	if res.Embedded != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateGCSBucketsInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "GCS bucket not found."
		detail := fmt.Sprintf("No GCS bucket found with the given bucket names %v", bucketNames)
		diags.AddError(summary, detail)
	}
	return diags
}
