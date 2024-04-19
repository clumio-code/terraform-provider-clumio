// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_organizational_unit Terraform datasource.

package clumio_organizational_unit

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioOrganizationalUnitDataSourceModel is the datasource model for the clumio_organizational_unit
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioOrganizationalUnitDataSourceModel struct {
	Name                types.String `tfsdk:"name"`
	OrganizationalUnits types.Set    `tfsdk:"organizational_units"`
}

// Schema defines the structure and constraints of the clumio_organizational_unit Terraform
// datasource. Schema is a method on the clumioOrganizationalUnitDataSource struct. It sets the
// schema for the clumio_organizational_unit Terraform datasource. The schema defines various
// attributes such as the name and organizational_units where 'organizational_units' is computed,
// meaning is determined by Clumio at runtime, whereas the 'name' attribute us used to determine
// the Clumio organizational units to retrieve.
func (r *clumioOrganizationalUnitDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaName: schema.StringAttribute{
				Description: "The name of the organizational unit.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaOrganizationalUnits: schema.SetNestedAttribute{
				Description: "OrganizationalUnits which match the given name.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the organizational unit.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "The name of the organizational unit.",
							Computed:    true,
						},
						schemaDescription: schema.StringAttribute{
							Description: "Brief description to denote details of the organizational" +
								" unit.",
							Computed: true,
						},
						schemaParentId: schema.StringAttribute{
							Description: "The identifier of the parent organizational unit under " +
								"which the organizational unit was created.",
							Computed: true,
						},
						schemaDescendantIds: schema.SetAttribute{
							Description: "List of all recursive descendent organizational units.",
							ElementType: types.StringType,
							Computed:    true,
						},
						schemaUsersWithRole: schema.SetNestedAttribute{
							Description: "List of user ids, with role assigned to this " +
								"organizational unit.",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									schemaUserId: schema.StringAttribute{
										Description: "Identifier of the user assigned to the " +
											"organizational unit.",
										Computed: true,
									},
									schemaAssignedRole: schema.StringAttribute{
										Description: "Identifier of the role associated with the " +
											"user assigned to the organizational unit.",
										Computed: true,
									},
								},
							},
						},
					},
				},
				Computed: true,
			},
		},
		Description: "clumio_organizational_unit data source is used to retrieve details of an" +
			" organizational unit for use in other resources.",
	}
}
