// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_aws_manual_connection_resources
// Terraform resource.
// This datasource is used to generate resources for deploying manual connections.
package clumio_aws_manual_connection_resources

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// clumioAwsManualConnectionResourcesDatasource is the struct backing the
// clumio_aws_manual_connection_resources Terraform resource. It holds the Clumio API client and any
//
//	other required state needed to connect AWS accounts to Clumio.
type clumioAwsManualConnectionResourcesDatasource struct {
	name         string
	client       *common.ApiClient
	awsTemplates sdkclients.AWSTemplatesClient
}

// NewAwsManualConnectionResourcesDataSource creates a new instance of
// clumioAwsManualConnectionResourcesDatasource. Its attributes are initialized later by Terraform
// via Metadata and Configure once the Provider is initialized.
func NewAwsManualConnectionResourcesDataSource() datasource.DataSource {
	return &clumioAwsManualConnectionResourcesDatasource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioAwsManualConnectionResourcesDatasource) Metadata(
	_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_aws_manual_connection_resources"
	res.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioAwsManualConnectionResourcesDatasource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.awsTemplates = sdkclients.NewAWSTemplatesClient(r.client.ClumioConfig)
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioAwsManualConnectionResourcesDatasource) Read(
	ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	// Retrieve the schema from the current Terraform state.
	var state clumioAwsManualConnectionResourcesModel
	diags := req.Config.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	diags = r.readAWSManualConnectionResources(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}
