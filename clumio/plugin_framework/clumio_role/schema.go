// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_role Terraform datasource.

package clumio_role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// permissionModel is a model used inside clumioRoleDataSourceModel to map to permissions available
// in each role
type permissionModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// clumioRoleDataSourceModel is the datasource model for the clumio_role Terraform datasource. It
// represents the schema of the datasource and the data it holds. This schema is used by customers
// to configure the datasource and by the Clumio provider to read and write the datasource.
type clumioRoleDataSourceModel struct {
	Id          types.String       `tfsdk:"id"`
	Name        types.String       `tfsdk:"name"`
	Description types.String       `tfsdk:"description"`
	UserCount   types.Int64        `tfsdk:"user_count"`
	Permissions []*permissionModel `tfsdk:"permissions"`
}

// Schema defines the structure and constraints of the clumio_role Terraform datasource. Schema is a
// method on the clumioRoleDataSource struct. It sets the schema for the clumio_role Terraform
// datasource, which fetches a role that can be assigned to a user within Clumio. The schema
// defines various attributes such as the role ID, name, description, user count, etc, some of which
// are computed, meaning they are determined by Clumio at runtime, whereas 'name' attribute is used
// determine the Clumio role to retrieve.
func (r *clumioRoleDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Clumio Role data source used to retrieve details of a role for use in other" +
			" resources.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the role.",
				Computed:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The unique name of the role from which to populate the data source.",
				Required:    true,
			},
			schemaDescription: schema.StringAttribute{
				Description: "Brief description to denote details of the role.",
				Computed:    true,
			},
			schemaUserCount: schema.Int64Attribute{
				Description: "Number of users to whom the role has been assigned.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			schemaPermissions: schema.ListNestedBlock{
				Description: "Permissions contained in the role.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						schemaDescription: schema.StringAttribute{
							Description: "Brief description to denote details of the permission.",
							Computed:    true,
						},
						schemaId: schema.StringAttribute{
							Description: "Unique identifier for the permission.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "Name of the permission.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
