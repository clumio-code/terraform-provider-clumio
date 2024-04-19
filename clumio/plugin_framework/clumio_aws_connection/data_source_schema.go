// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_aws_connection Terraform datasource.

package clumio_aws_connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioAWSConnectionDataSourceModel is the datasource model for the clumio_aws_connection
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioAWSConnectionDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	AccountNativeID types.String `tfsdk:"account_native_id"`
	AWSRegion       types.String `tfsdk:"aws_region"`
}

// Schema defines the structure and constraints of the clumio_aws_connection Terraform datasource.
// Schema is a method on the clumioAWSConnectionDataSource struct. It sets the schema for the
// clumio_aws_connection Terraform datasource. The schema defines various attributes such as the
// aws_connection account_id, aws_region and id where 'id' is computed, meaning is determined by
// Clumio at runtime, whereas the 'account_id' and 'aws_region' attribute us used to determine the
// Clumio aws connection to retrieve.
func (r *clumioAWSConnectionDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the aws connection.",
				Computed:    true,
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "Identifier of the AWS account linked with Clumio.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "Region of the AWS account linked with Clumio.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Description: "clumio_aws_connection data source is used to retrieve details of the" +
			" aws connection for use in other resources.",
	}
}
