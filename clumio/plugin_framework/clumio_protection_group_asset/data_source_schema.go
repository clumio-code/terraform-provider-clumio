// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_protection_group_asset Terraform datasource.

package clumio_protection_group_asset

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioProtectionGroupAssetDataSourceModel is the datasource model for the clumio_protection_group_asset
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioProtectionGroupAssetDataSourceModel struct {
	BucketID          types.String `tfsdk:"bucket_id"`
	ProtectionGroupID types.String `tfsdk:"protection_group_id"`
	Id                types.String `tfsdk:"id"`
}

// Schema defines the structure and constraints of the clumio_protection_group_asset Terraform
// datasource. Schema is a method on the clumioProtectionGroupAssetDataSource struct. It sets the
// schema for the clumio_protection_group_asset Terraform datasource. The schema defines various
// attributes such as the bucket_id, protection_group_id and id where 'id' is computed, meaning is
// determined by Clumio at runtime, whereas the 'bucket_id' and 'protection_group_id' attributes
// are used to determine the Clumio protection group asset to retrieve.
func (r *clumioProtectionGroupAssetDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio Protection Group asset.",
				Computed:    true,
			},
			schemaBucketId: schema.StringAttribute{
				Description: "Clumio assigned unique identifier of the AWS S3 bucket.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaProtectionGroupId: schema.StringAttribute{
				Description: "Unique identifier of the Protection Group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Description: "clumio_protection_group_asset data source is used to retrieve the id of the" +
			" protection group asset",
	}
}
