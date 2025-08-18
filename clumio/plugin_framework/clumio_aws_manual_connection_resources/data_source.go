// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the AWS manual connection resources API and set the
// attributes from the response of the API in the data source model.

package clumio_aws_manual_connection_resources

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// readAWSManualConnectionResources invokes the API to read the AWS manual connection resources and
// from the response populates the schema attributes.
func (r *clumioAwsManualConnectionResourcesDatasource) readAWSManualConnectionResources(
	_ context.Context, state *clumioAwsManualConnectionResourcesModel) diag.Diagnostics {

	var diags diag.Diagnostics

	assetsEnabled := []*string{}
	if state.AssetsEnabled.EBS.ValueBool() {
		enabled := EBS
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.S3.ValueBool() {
		enabled := S3
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.RDS.ValueBool() {
		enabled := RDS
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.DynamoDB.ValueBool() {
		enabled := DynamoDB
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.EC2MSSQL.ValueBool() {
		enabled := EC2MSSQL
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	showManualResources := true
	returnGroupToken := false

	// Call the Clumio API to read the resources for the provided configuration.
	apiRes, apiErr := r.awsTemplates.CreateConnectionTemplate(
		&returnGroupToken,
		&models.CreateConnectionTemplateV1Request{
			ShowManualResources: &showManualResources,
			AssetTypesEnabled:   assetsEnabled,
			AwsAccountId:        state.AccountId.ValueStringPointer(),
			AwsRegion:           state.AwsRegion.ValueStringPointer(),
		})

	if apiErr != nil {
		summary := "Failed to get resources from API"
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	if apiRes == nil || apiRes.Resources == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the resources obtained from the Clumio API response into stringified format and update
	// the state with it.
	stringifiedResources := stringifyResources(apiRes.Resources)
	state.Resources = types.StringPointerValue(stringifiedResources)

	return diags
}
