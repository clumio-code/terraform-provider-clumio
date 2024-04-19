// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_protection_group Terraform datasource.

package clumio_protection_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioProtectionGroupDataSourceModel is the datasource model for the clumio_protection_group
// Terraform datasource. It
// represents the schema of the datasource and the data it holds. This schema is used by customers
// to configure the datasource and by the Clumio provider to read and write the datasource.
type clumioProtectionGroupDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Schema defines the structure and constraints of the clumio_protection_group Terraform datasource.
// Schema is a method on the clumioProtectionGroupDataSource struct. It sets the schema for the
// clumio_protection_group Terraform datasource. The schema defines various attributes such as the
// protection_group name and id where 'id' is computed, meaning is determined by Clumio at runtime,
// whereas the 'name' attribute us used to determine the Clumio protection group to retrieve.
func (r *clumioProtectionGroupDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the protection group.",
				Computed:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the protection group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Description: "clumio_protection_group data source is used to retrieve details of a" +
			" protection group for use in other resources.",
	}
}
