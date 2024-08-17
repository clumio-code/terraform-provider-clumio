// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio DynamoDB tables SDK API to perform read operation
// and set the attributes from the response of the API in the data source model.

package clumio_dynamodb_tables

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"strings"
)

// readDynamoDBTables invokes the API to read the dynamoDBTableClient and from the response
// populates the attributes of the DynamoDB tables.
func (r *clumioDynamoDBTablesDataSource) readDynamoDBTables(
	ctx context.Context, model *clumioDynamoDBTablesDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Prepare the query filter.
	nameFilter := ""
	accountNativeIdFilter := ""
	regionFilter := ""
	tableNativeIdFilter := ""
	filters := make([]string, 0)

	tableNativeId := model.TableNativeID.ValueString()
	if tableNativeId != "" {
		tableNativeIdFilter = fmt.Sprintf(`"table_native_id": {"$eq":"%s"}`, tableNativeId)
		filters = append(filters, tableNativeIdFilter)
	}
	name := model.Name.ValueString()
	if tableNativeId == "" && name != "" {
		nameFilter = fmt.Sprintf(`"name": {"$contains":"%s"}`, name)
		filters = append(filters, nameFilter)
	}

	accountNativeId := model.AccountNativeID.ValueString()
	if accountNativeId != "" {
		accountNativeIdFilter = fmt.Sprintf(`"account_native_id": {"$eq":"%s"}`, accountNativeId)
		filters = append(filters, accountNativeIdFilter)
	}

	region := model.Region.ValueString()
	if region != "" {
		regionFilter = fmt.Sprintf(`"aws_region": {"$eq":"%s"}`, region)
		filters = append(filters, regionFilter)
	}

	filter := fmt.Sprintf("{%s}", strings.Join(filters, ","))
	// Call the Clumio API to list the DynamoDB tables.
	limit := int64(10000)
	res, apiErr := r.dynamoDBTableClient.ListAwsDynamodbTables(&limit, nil, &filter, nil, nil)
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

	// Convert the Clumio API response for the DynamoDB tables into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateDynamoDBTablesInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "DynamoDB table not found."
		detail := "No DynamoDB table found with the given query attributes."
		diags.AddError(summary, detail)
	}
	return diags
}
