// Copyright 2026. Clumio, Inc.

// This file holds the logic to invoke the Clumio DynamoDB table backups SDK API to perform read
// operation and set the attributes from the response of the API in the data source model.

package clumio_dynamodb_backups

import (
	"context"
	"fmt"
	"strings"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readDynamoDBBackups invokes the API to list the backups of the DynamoDB table given in the model
// and from the response populates the attributes of the backups, sorted newest first.
func (r *clumioDynamoDBBackupsDataSource) readDynamoDBBackups(
	ctx context.Context, model *clumioDynamoDBBackupsDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Prepare the query filter. table_id is a required attribute so it is always present.
	filters := []string{
		fmt.Sprintf(`"table_id": {"$eq":%s}`,
			common.JSONEscapeFilterValue(model.TableID.ValueString())),
	}

	tsFilters := make([]string, 0)
	if before := model.BeforeTimestamp.ValueString(); before != "" {
		tsFilters = append(tsFilters,
			fmt.Sprintf(`"$lte":%s`, common.JSONEscapeFilterValue(before)))
	}
	if after := model.AfterTimestamp.ValueString(); after != "" {
		tsFilters = append(tsFilters,
			fmt.Sprintf(`"$gt":%s`, common.JSONEscapeFilterValue(after)))
	}
	if len(tsFilters) > 0 {
		filters = append(filters,
			fmt.Sprintf(`"start_timestamp": {%s}`, strings.Join(tsFilters, ",")))
	}

	// The backups endpoint only supports the "$all" condition on the "type" field; "$in" and "$eq"
	// are rejected with a 400.
	if clumioType := model.ClumioType.ValueString(); clumioType != "" {
		filters = append(filters, fmt.Sprintf(`"type": {"$all":[%s]}`,
			common.JSONEscapeFilterValue(clumioType)))
	}
	filter := fmt.Sprintf("{%s}", strings.Join(filters, ","))

	// Call the Clumio API to list the DynamoDB table backups, sorted such that the most recent
	// backup is the first item of the response. Pagination is not followed: if the backup count
	// ever exceeds the limit, the oldest backups are dropped, while the descending sort keeps the
	// primary use case of backups[0] being the latest backup correct.
	limit := int64(10000)
	sort := "-" + schemaStartTimestamp
	res, apiErr := r.backupClient.ListBackupAwsDynamodbTables(&limit, nil, &sort, &filter)
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

	// Convert the Clumio API response for the DynamoDB table backups into the datasource schema
	// model.
	if res.Embedded != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateDynamoDBBackupsInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "DynamoDB table backup not found."
		detail := "No DynamoDB table backup found with the given query attributes."
		diags.AddError(summary, detail)
	}
	return diags
}
