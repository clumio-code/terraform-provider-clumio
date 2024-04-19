// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_user Terraform datasource.

package clumio_user

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioUserDataSourceModel is the datasource model for the clumio_user
// Terraform datasource. It
// represents the schema of the datasource and the data it holds. This schema is used by customers
// to configure the datasource and by the Clumio provider to read and write the datasource.
type clumioUserDataSourceModel struct {
	Name   types.String `tfsdk:"name"`
	RoleId types.String `tfsdk:"role_id"`
	Users  types.Set    `tfsdk:"users"`
}

// Schema defines the structure and constraints of the clumio_user Terraform datasource.
// Schema is a method on the clumioUserDataSource struct. It sets the schema for the
// clumio_user Terraform datasource. The schema defines various attributes such as the
// user name and id where 'id' is computed, meaning is determined by Clumio at runtime,
// whereas the 'name' attribute us used to determine the Clumio user to retrieve.
func (r *clumioUserDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaName: schema.StringAttribute{
				Description: "The name of the user.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaRoleId: schema.StringAttribute{
				Description: "Unique identifier of the role assigned to the user.",
				Optional:    true,
			},
			schemaUsers: schema.SetNestedAttribute{
				Description: "Users that match the given name and/or role_id.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the user.",
							Computed:    true,
						},
						schemaFullName: schema.StringAttribute{
							Description: "The name of the user.",
							Computed:    true,
						},
						schemaAccessControlConfiguration: schema.SetNestedAttribute{
							Description: "Identifiers of the organizational units, along with the " +
								"identifier of the role assigned to the user.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									schemaRoleId: schema.StringAttribute{
										Description: "Identifier of the role assigned to the user.",
										Computed:    true,
									},
									schemaOrganizationalUnitIds: schema.SetAttribute{
										Description: "Identifiers of the organizational units " +
											"assigned to the user.",
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
				Computed: true,
			},
		},
		Description: "clumio_user data source is used to retrieve details of a" +
			" user for use in other resources.",
	}
}

// ConfigValidators to check if at least one of name or role_id is specified.
func (r *clumioUserDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot(schemaName),
			path.MatchRoot(schemaRoleId),
		),
	}
}
