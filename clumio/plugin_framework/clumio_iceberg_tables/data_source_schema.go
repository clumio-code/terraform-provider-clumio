// Copyright 2026. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_iceberg_tables Terraform datasource.

package clumio_iceberg_tables

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	validators "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioIcebergTablesDataSourceModel is the datasource model for the clumio_iceberg_tables
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioIcebergTablesDataSourceModel struct {
	AccountNativeID types.String `tfsdk:"account_native_id"`
	Region          types.String `tfsdk:"aws_region"`
	Name            types.String `tfsdk:"name"`
	Catalog         types.String `tfsdk:"catalog"`
	CatalogType     types.String `tfsdk:"catalog_type"`
	Namespace       types.String `tfsdk:"namespace"`
	IcebergTables   types.List   `tfsdk:"iceberg_tables"`
}

// Schema defines the structure and constraints of the clumio_iceberg_tables Terraform datasource.
// Schema is a method on the clumioIcebergTablesDataSource struct. It sets the schema for the
// clumio_iceberg_tables Terraform datasource. The schema defines various attributes such as the
// account_native_id, aws_region, name, catalog, catalog_type, namespace and iceberg_tables where
// 'iceberg_tables' is computed, meaning is determined by Clumio at runtime, whereas the other
// attributes are used to determine the Clumio Iceberg tables to retrieve.
func (r *clumioIcebergTablesDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaAccountNativeId: schema.StringAttribute{
				Description: "The identifier of the AWS account under which the Iceberg tables" +
					" were created.",
				Required: true,
			},
			schemaRegion: schema.StringAttribute{
				Description: "The AWS region associated with the Iceberg tables.",
				Required:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The Iceberg table name to be queried.",
				Optional:    true,
			},
			schemaCatalog: schema.StringAttribute{
				Description: "The catalog of the Iceberg tables to be queried.",
				Optional:    true,
			},
			schemaCatalogType: schema.StringAttribute{
				Description: "The catalog type of the Iceberg tables to be queried. The" +
					" supported values are `aws_iceberg_glue_table` and `aws_iceberg_s3_table`.",
				Optional: true,
				Validators: []validator.String{
					validators.OneOf(catalogTypeGlue, catalogTypeS3Table),
				},
			},
			schemaNamespace: schema.StringAttribute{
				Description: "The namespace of the Iceberg tables to be queried.",
				Optional:    true,
			},
			schemaIcebergTables: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the Iceberg table.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "Name of the Iceberg table.",
							Computed:    true,
						},
						schemaCatalog: schema.StringAttribute{
							Description: "Catalog of the Iceberg table.",
							Computed:    true,
						},
						schemaCatalogType: schema.StringAttribute{
							Description: "Catalog type of the Iceberg table.",
							Computed:    true,
						},
						schemaNamespace: schema.StringAttribute{
							Description: "Namespace of the Iceberg table.",
							Computed:    true,
						},
					},
				},
				Computed:    true,
				Description: "List of Iceberg tables which matched the query criteria.",
			},
		},
		Description: "clumio_iceberg_tables data source is used to retrieve details of the" +
			" Iceberg tables for use in other resources.",
	}
}

// ConfigValidators checks that at least one of name, catalog or namespace is specified.
func (r *clumioIcebergTablesDataSource) ConfigValidators(
	_ context.Context) []datasource.ConfigValidator {

	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot(schemaName),
			path.MatchRoot(schemaCatalog),
			path.MatchRoot(schemaNamespace),
		),
	}
}
