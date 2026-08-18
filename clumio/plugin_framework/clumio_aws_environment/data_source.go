// Copyright 2026. Clumio, Inc.

// This file holds the logic to invoke the Clumio AWS environments SDK API to perform read
// operation and set the attributes from the response of the API in the data source model.

package clumio_aws_environment

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// readAWSEnvironment invokes the API to read the AWS environment corresponding to the AWS account
// and region given in the model and from the response populates the identifier of the environment.
func (r *clumioAWSEnvironmentDataSource) readAWSEnvironment(
	ctx context.Context, model *clumioAWSEnvironmentDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Look up the AWS environment for the given account and region. Note that an AWS environment is
	// only created once the corresponding AWS connection is connected, so the lookup may fail for a
	// non-connected account.
	env, err := common.LookupAWSEnvironment(
		r.envClient, model.AccountNativeID.ValueString(), model.AWSRegion.ValueString())
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to read %s", r.name), err.Error())
		return diags
	}
	model.ID = types.StringPointerValue(env.Id)

	return diags
}
