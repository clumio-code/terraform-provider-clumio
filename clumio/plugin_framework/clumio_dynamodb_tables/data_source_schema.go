// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_dynamo_db_tables Terraform datasource.

package clumio_dynamodb_tables

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioDynamoDBTablesDataSourceModel is the datasource model for the clumio_dynamo_db_tables
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioDynamoDBTablesDataSourceModel struct {
	AccountNativeID types.String `tfsdk:"account_native_id"`
	Region          types.String `tfsdk:"aws_region"`
	Name            types.String `tfsdk:"name"`
	TableNativeID   types.String `tfsdk:"table_native_id"`
	DynamoDBTables  types.List   `tfsdk:"dynamodb_tables"`
}

// Schema defines the structure and constraints of the clumio_dynamo_db_tables Terraform datasource.
// Schema is a method on the clumioDynamoDBTablesDataSource struct. It sets the schema for the
// clumio_dynamo_db_tables Terraform datasource. The schema defines various attributes such as the
// account_native_id, region, table_name, table_native_id and dynamodb_tables where
// 'dynamodb_tables' is computed, meaning is determined by Clumio at runtime, whereas the other
// attributes are used to determine the Clumio DynamoDB tables to retrieve.
func (r *clumioDynamoDBTablesDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaAccountNativeId: schema.StringAttribute{
				Description: "The identifier of the AWS account under which the DynamoDB " +
					"bucket was created.",
				Required: true,
			},
			schemaRegion: schema.StringAttribute{
				Description: "The AWS region associated with the DynamoDB tables.",
				Required:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The DynamoDB table name to be queried.",
				Optional:    true,
			},
			schemaTableNativeId: schema.StringAttribute{
				Description: "Native identifier of the DynamoDB table to be queried.",
				Optional:    true,
			},
			schemaDynamoDBTables: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the DynamoDB table.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "Name of the DynamoDB table.",
							Computed:    true,
						},
						schemaTableNativeId: schema.StringAttribute{
							Description: "Native identifier of the DynamoDB table.",
							Computed:    true,
						},
					},
				},
				Computed:    true,
				Description: "List of DynamoDB tables which matched the query criteria.",
			},
		},
		Description: "clumio_dynamo_db_tables data source is used to retrieve details of the" +
			" DynamoDB tables for use in other resources.",
	}
}

// ConfigValidators to check if at least one of name, operation_types or activation_status is specified.
func (r *clumioDynamoDBTablesDataSource) ConfigValidators(
	_ context.Context) []datasource.ConfigValidator {

	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot(schemaName),
			path.MatchRoot(schemaTableNativeId),
		),
	}
}
