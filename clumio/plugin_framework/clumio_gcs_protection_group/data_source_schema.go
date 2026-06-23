// Copyright 2026. Clumio, Inc.

// This file holds the type definition and Schema data source function used by the data source model
// for the clumio_gcs_protection_group Terraform data source.

package clumio_gcs_protection_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioGCSProtectionGroupDataSourceModel is the data source model for the
// clumio_gcs_protection_group Terraform data source. The name attribute is used to query the
// protection group, whereas id is computed.
type clumioGCSProtectionGroupDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Schema defines the structure and constraints of the clumio_gcs_protection_group Terraform data
// source. It retrieves a GCP protection group by its name and exposes its id for use in other
// resources.
func (r *clumioGCSProtectionGroupDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the GCP protection group.",
				Computed:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the GCP protection group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Description: "clumio_gcs_protection_group data source is used to retrieve details of a" +
			" GCP protection group for use in other resources.",
	}
}
