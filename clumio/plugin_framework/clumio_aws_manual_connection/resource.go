// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio AWS Manual Connection SDK APIs to perform CRUD
// operations and set the attributes from the response of the API in the resource model.

package clumio_aws_manual_connection

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// createAWSManualConnection invokes the API to create the manual connection and from the response
// populates the computed attributes of the connection.
func (r *clumioAWSManualConnectionResource) createAWSManualConnection(
	ctx context.Context, plan *clumioAWSManualConnectionResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the common util to deploy the manually configured resources for the connection.
	diags = r.clumioSetManualResourcesCommon(ctx, *plan)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	accountId := plan.AccountId.ValueString()
	awsRegion := plan.AwsRegion.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf("%v_%v", accountId, awsRegion))

	return diags
}

// updateAWSManualConnection invokes the API to update the manual connection and from the response
// populates the computed attributes of the connection.
func (r *clumioAWSManualConnectionResource) updateAWSManualConnection(
	ctx context.Context, plan *clumioAWSManualConnectionResourceModel,
	state *clumioAWSManualConnectionResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Block update if downgrading of assets is attempted.
	if isAssetConfigDowngraded(plan, state) {
		summary := fmt.Sprintf("Unable to update %s ", r.name)
		detail := "Downgrading assets is not allowed."
		diags.AddError(summary, detail)
	}

	// Call the Clumio API to update the manual connection.
	diags = r.clumioSetManualResourcesCommon(ctx, *plan)
	diags.Append(diags...)
	return diags
}
