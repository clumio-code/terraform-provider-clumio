// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio AWS Connection SDK API to perform read operation
// and set the attributes from the response of the API in the data source model.

package clumio_aws_connection

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readAWSConnection invokes the API to read the awsConnectionClient and from the response
// populates the attributes of the aws connection.
func (r *clumioAWSConnectionDataSource) readAWSConnection(
	_ context.Context, model *clumioAWSConnectionDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Prepare the query filter.
	accountNativeId := model.AccountNativeID.ValueString()
	region := model.AWSRegion.ValueString()
	filter := fmt.Sprintf(`{"account_native_id": {"$in":["%s"]}, "aws_region": {"$in":["%s"]}}`,
		accountNativeId, region)

	// Call the Clumio API to list the aws connections.
	res, apiErr := r.awsConnectionClient.ListAwsConnections(nil, nil, &filter)
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
	// Convert the Clumio API response for the AWS connections into the datasource schema model.
	if res.CurrentCount == nil || *res.CurrentCount == 0 {
		summary := "AWS connection not found"
		detail := fmt.Sprintf(
			"Expected aws connection with the specified Acount Native ID %s and AWS region %s"+
				" is not found", accountNativeId, region)
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response for the policies into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		model.Id = basetypes.NewStringValue(*res.Embedded.Items[0].Id)
	}
	return diags
}
