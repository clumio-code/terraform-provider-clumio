// Copyright 2026. Clumio, Inc.

// This file holds the logic to invoke the Clumio Iceberg tables SDK API to perform read operation
// and set the attributes from the response of the API in the data source model.

package clumio_iceberg_tables

import (
	"context"
	"fmt"
	"strings"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readIcebergTables invokes the API to read the Iceberg tables and from the response populates the
// attributes of the Iceberg tables.
func (r *clumioIcebergTablesDataSource) readIcebergTables(
	ctx context.Context, model *clumioIcebergTablesDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	accountNativeId := model.AccountNativeID.ValueString()
	region := model.Region.ValueString()

	// The Iceberg tables API does not accept an AWS region filter, so the account and region are
	// first resolved into the Clumio environment which the API does accept.
	env, err := common.LookupAWSEnvironment(r.awsEnvironmentClient, accountNativeId, region)
	if err != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
		diags.AddError(summary, err.Error())
		return diags
	}

	// Prepare the query filter.
	filters := []string{
		fmt.Sprintf(`"environment_id": {"$eq":%s}`, common.JSONEscapeFilterValue(*env.Id)),
		fmt.Sprintf(
			`"account_native_id": {"$eq":%s}`, common.JSONEscapeFilterValue(accountNativeId)),
	}
	if name := model.Name.ValueString(); name != "" {
		filters = append(filters,
			fmt.Sprintf(`"table_name": {"$eq":%s}`, common.JSONEscapeFilterValue(name)))
	}
	if catalog := model.Catalog.ValueString(); catalog != "" {
		filters = append(filters,
			fmt.Sprintf(`"catalog": {"$eq":%s}`, common.JSONEscapeFilterValue(catalog)))
	}
	if namespace := model.Namespace.ValueString(); namespace != "" {
		filters = append(filters,
			fmt.Sprintf(`"namespace": {"$eq":%s}`, common.JSONEscapeFilterValue(namespace)))
	}
	if catalogType := model.CatalogType.ValueString(); catalogType != "" {
		filters = append(filters,
			fmt.Sprintf(`"catalog_type": {"$in":[%s]}`,
				common.JSONEscapeFilterValue(catalogType)))
	}
	filter := fmt.Sprintf("{%s}", strings.Join(filters, ","))

	// Call the Clumio API to list the Iceberg tables.
	limit := listLimit
	res, apiErr := r.icebergTableClient.ListAwsIcebergTables(&limit, nil, &filter, nil, nil)
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

	// Convert the Clumio API response for the Iceberg tables into the datasource schema model.
	if res.Embedded != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateIcebergTablesInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "Iceberg table not found."
		detail := "No Iceberg table found with the given query attributes."
		diags.AddError(summary, detail)
	}
	return diags
}
