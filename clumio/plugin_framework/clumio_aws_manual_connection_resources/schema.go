// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_aws_manual_connection_resources Terraform resource.

package clumio_aws_manual_connection_resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioAwsManualConnectionResourcesModel is the resource model for the
// clumio_aws_manual_connection_resources Terraform resource. It represents the schema of the
// resource and the data it holds. This schema is used by customers to configure the resource and by
// the Clumio provider to read and write the resource.
type clumioAwsManualConnectionResourcesModel struct {
	ID            types.String            `tfsdk:"id"`
	AccountId     types.String            `tfsdk:"account_native_id"`
	AwsRegion     types.String            `tfsdk:"aws_region"`
	AssetsEnabled *assetTypesEnabledModel `tfsdk:"asset_types_enabled"`
	Resources     types.String            `tfsdk:"resources"`
}

// assetTypesEnabled is a model used inside clumioAwsManualConnectionResourcesModel for determining
// the asset types enabled for the configuration
type assetTypesEnabledModel struct {
	EBS      types.Bool `tfsdk:"ebs"`
	RDS      types.Bool `tfsdk:"rds"`
	DynamoDB types.Bool `tfsdk:"ddb"`
	S3       types.Bool `tfsdk:"s3"`
	EC2MSSQL types.Bool `tfsdk:"mssql"`
}

// Schema defines the structure and constraints of the clumio_aws_manual_connection_resources
// Terraform datasource. Schema is a method on the clumioAwsManualConnectionResourcesDatasource
// struct. It sets the schema for the clumio_aws_manual_connection_resources Terraform datasource,
// which is used to generate resources required for deployment of manual connections. The schema
// defines various attributes such as the connection ID, AWS account ID, AWS region etc. The
// "resources" attribute is computed, meaning it is determine by Clumio at runtime, while others are
// required inputs from the user.
func (*clumioAwsManualConnectionResourcesDatasource) Schema(
	_ context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Clumio AWS Manual Connection Resources Datasource to get resources for manual connections",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Combination of provided Account Native ID and Aws Region",
				Computed:    true,
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "AWS Account ID to be connected to Clumio",
				Required:    true,
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "AWS Region to be connected to Clumio",
				Required:    true,
			},
			schemaAssetTypesEnabled: schema.ObjectAttribute{
				Description: "Assets to be connected to Clumio",
				Required:    true,
				AttributeTypes: map[string]attr.Type{
					schemaIsEbsEnabled:      types.BoolType,
					schemaIsDynamoDBEnabled: types.BoolType,
					schemaIsRDSEnabled:      types.BoolType,
					schemaIsS3Enabled:       types.BoolType,
					schemaIsMssqlEnabled:    types.BoolType,
				},
			},
			schemaResources: schema.StringAttribute{
				Description: "Generated manual resources for provided configuration",
				Computed:    true,
			},
		},
	}
}
