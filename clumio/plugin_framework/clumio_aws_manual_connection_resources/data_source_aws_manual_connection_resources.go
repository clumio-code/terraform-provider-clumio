// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_aws_manual_connection_resources
// Terraform resource.
// This datasource is used to generate resources for deploying manual connections.
package clumio_aws_manual_connection_resources

import (
	"context"

	awsTemplates "github.com/clumio-code/clumio-go-sdk/controllers/aws_templates"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioAwsManualConnectionResourcesDatasource is the struct backing the
// clumio_aws_manual_connection_resources Terraform resource. It holds the Clumio API client and any
//  other required state needed to connect AWS accounts to Clumio.
type clumioAwsManualConnectionResourcesDatasource struct {
	name            string
	client          *common.ApiClient
	awsTemplates    awsTemplates.AwsTemplatesV1Client
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
	r.awsTemplates = awsTemplates.NewAwsTemplatesV1(r.client.ClumioConfig)
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

	// Call the Clumio API to read the resources for the provided configuration.
	apiRes, apiErr := r.awsTemplates.CreateConnectionTemplate(&models.CreateConnectionTemplateV1Request{
		ShowManualResources: &showManualResources,
		AssetTypesEnabled: assetsEnabled,
		AwsAccountId: state.AccountId.ValueStringPointer(),
		AwsRegion: state.AwsRegion.ValueStringPointer(),
	})
	if apiRes.Resources == nil {
		res.Diagnostics.AddError("Failed to get resources from API", common.ParseMessageFromApiError(
			apiErr))
	}
	if res.Diagnostics.HasError() {
		return
	}

	// Convert the resources obtained from the Clumio API response into stringified format and update
	// the state with it.
	stringifiedResources := stringifyResources(apiRes.Resources)
	state.Resources = types.StringPointerValue(stringifiedResources)

	// Set the schema into the Terraform state.
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}
}
