// Copyright 2023. Clumio, Inc.
//
// aws_manual_connection_resources datasource definition implementation.
package clumio_aws_manual_connection_resources

import (
	"context"
	"encoding/json"

	aws_templates "github.com/clumio-code/clumio-go-sdk/controllers/aws_templates"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EBS      = "EBS"
	S3       = "S3"
	DynamoDB = "DynamoDB"
	RDS      = "RDS"
	EC2MSSQL = "EC2MSSQL"
)

// AwsManualConnectionResources model
type AwsManualConnectionResourcesModel struct {
	ID            types.String                 `tfsdk:"id"`
	AccountId     types.String                 `tfsdk:"account_native_id"`
	AwsRegion     types.String                 `tfsdk:"aws_region"`
	AssetsEnabled *AssetTypesEnabledModel      `tfsdk:"asset_types_enabled"`
	Resources     types.String                 `tfsdk:"resources"`
}

// AssetTypesEnabled model
type AssetTypesEnabledModel struct {
	EBS      types.Bool `tfsdk:"ebs"`
	RDS      types.Bool `tfsdk:"rds"`
	DynamoDB types.Bool `tfsdk:"ddb"`
	S3       types.Bool `tfsdk:"s3"`
	EC2MSSQL types.Bool `tfsdk:"mssql"`
}

// awsManualConnectionResourcesDatasource is the datasource implementation.
type awsManualConnectionResourcesDatasource struct {
	client *common.ApiClient
}

// NewAwsManualConnectionResource is a helper function to simplify the provider implementation.
func NewAwsManualConnectionResourcesDataSource() datasource.DataSource {
	return &awsManualConnectionResourcesDatasource{}
}

// Metadata returns the datasource type name.
func (r *awsManualConnectionResourcesDatasource) Metadata(
	_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_aws_manual_connection_resources"
}

// Configure adds the provider configured client to the data source.
func (r *awsManualConnectionResourcesDatasource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
}

// Schema defines the schema for the resource.
func (*awsManualConnectionResourcesDatasource) Schema(_ context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Clumio AWS Manual Connection Resources Datasource to get resources for manual connections",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Combination of provided Account Native ID and Aws Region",
				Computed:    true,
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "AWS Account ID to be connected to Clumio",
				Required: true,
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "AWS Region to be connected to Clumio",
				Required:    true,
			},
			schemaAssetTypesEnabled: schema.ObjectAttribute{
				Description: "Assets to be connection to Clumio",
				Required:    true,
				AttributeTypes: map[string]attr.Type{
					schemaIsEbsEnabled: types.BoolType,
					schemaIsDynamoDBEnabled: types.BoolType,
					schemaIsRDSEnabled: types.BoolType,
					schemaIsS3Enabled: types.BoolType,
					schemaIsMssqlEnabled: types.BoolType,
				},
			},
			schemaResources: schema.StringAttribute{
				Description: "Generated manual resources for provided config",
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *awsManualConnectionResourcesDatasource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	// Get current state
	var state AwsManualConnectionResourcesModel
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

	awsAccountId := state.AccountId.ValueString()
	awsRegion := state.AwsRegion.ValueString()
	showManualResources := true

	awsTemplates := aws_templates.NewAwsTemplatesV1(r.client.ClumioConfig)
	apiRes, apiErr := awsTemplates.CreateConnectionTemplate(&models.CreateConnectionTemplateV1Request{
		ShowManualResources: &showManualResources,
		AssetTypesEnabled: assetsEnabled,
		AwsAccountId: &awsAccountId,
		AwsRegion: &awsRegion,
	})

	if apiRes.Resources == nil {
		res.Diagnostics.AddError("Failed to get resources from API", string(apiErr.Response))
	}
	if res.Diagnostics.HasError() {
		return
	}

	// Set refreshed state.
	stringifiedResources := stringifyResources(apiRes.Resources)
	state.Resources = types.StringValue(*stringifiedResources)
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}
}

func stringifyResources(resources *models.CategorisedResources) *string {
	bytes, err := json.Marshal(resources)
	if err != nil {
			return nil
	}
  stringifiedResources := string(bytes)
	return &stringifiedResources
}
