// Copyright 2026. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_aws_environment Terraform datasource.

package clumio_aws_environment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioAWSEnvironmentDataSourceModel is the datasource model for the clumio_aws_environment
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioAWSEnvironmentDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	AccountNativeID types.String `tfsdk:"account_native_id"`
	AWSRegion       types.String `tfsdk:"aws_region"`
}

// Schema defines the structure and constraints of the clumio_aws_environment Terraform
// datasource. The schema defines the account_native_id and aws_region attributes, which are used
// to determine the Clumio AWS environment to retrieve, and the computed attribute 'id' which holds
// the identifier of the environment.
func (r *clumioAWSEnvironmentDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the AWS environment.",
				Computed:    true,
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "Identifier of the AWS account associated with the environment.",
				Required:    true,
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "The AWS region associated with the environment.",
				Required:    true,
			},
		},
		Description: "clumio_aws_environment data source is used to retrieve the identifier of " +
			"the Clumio AWS environment associated with an AWS account and region for use in " +
			"other resources or actions.",
	}
}
