// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_organizational_unit Terraform resource.

package clumio_organizational_unit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// userWithRole is a struct to format the results of getUsersFromHTTPRes util
type userWithRole struct {
	UserId       types.String `tfsdk:"user_id"`
	AssignedRole types.String `tfsdk:"assigned_role"`
}

// clumioOrganizationalUnitResourceModel is the resource model for the clumio_organizational_unit
// Terraform resource. It represents the schema of the resource and the data it holds. This schema
// is used by customers to configure the resource and by the Clumio provider to read and write the
// resource.
type clumioOrganizationalUnitResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	ParentId                  types.String `tfsdk:"parent_id"`
	ChildrenCount             types.Int64  `tfsdk:"children_count"`
	ConfiguredDatasourceTypes types.List   `tfsdk:"configured_datasource_types"`
	DescendantIds             types.List   `tfsdk:"descendant_ids"`
	UserCount                 types.Int64  `tfsdk:"user_count"`
	Users                     types.List   `tfsdk:"users"`
	UsersWithRole             types.List   `tfsdk:"users_with_role"`
}

// Schema defines the structure and constraints of the clumio_organizational_unit Terraform resource.
// Schema is a method on the clumioOrganizationalUnitResource struct. It sets the schema for the
// clumio_organizational_unit Terraform resource, which is used to manage Organizational Units
// within Clumio. The schema defines various attributes such as the organizational unit ID,
// name, description, etc. Some of these attributes are computed, meaning they are determined by
// Clumio at runtime, while others are required or optional inputs from the user.
func (r *clumioOrganizationalUnitResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for creating and managing Organizational Unit in Clumio. NOTE: If " +
			"this is the first time creating an Organizational Unit, the AWS \"data group\" that " +
			"denotes the top-most level that an Organizational Unit can see must be manually " +
			"selected once from the Clumio portal under \"Settings -> Access Management -> " +
			"Organizational Units\".",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio organizational unit.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "Unique name assigned to the organizational unit. " +
					"Root organizational unit is named as 'Global organizational unit'.",
				Required: true,
			},
			schemaDescription: schema.StringAttribute{
				Description: "Brief description to denote details of the organizational unit.",
				Optional:    true,
			},
			schemaParentId: schema.StringAttribute{
				Description: "The identifier of the parent organizational unit under which the new " +
					"organizational unit is to be created. If not provided, the resource will be " +
					"created under the default organizational unit associated with the credentials " +
					"used to create the organizational unit. Root organizational unit ID is " +
					"'00000000-0000-0000-0000-000000000000'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaChildrenCount: schema.Int64Attribute{
				Description: "Number of immediate children of the organizational unit.",
				Computed:    true,
			},
			schemaConfiguredDatasourceTypes: schema.ListAttribute{
				Description: "Datasource types configured in this organizational unit." +
					" Possible values include aws, microsoft365, vmware, or mssql.",
				ElementType: types.StringType,
				Computed:    true,
			},
			schemaDescendantIds: schema.ListAttribute{
				Description: "List of all recursive descendent organizational units.",
				ElementType: types.StringType,
				Computed:    true,
			},
			schemaUserCount: schema.Int64Attribute{
				Description: "Number of users to whom this organizational unit or any" +
					" of its descendants have been assigned.",
				Computed: true,
			},
			schemaUsers: schema.ListAttribute{
				Description: "List of user ids to assign this organizational unit.",
				ElementType: types.StringType,
				Computed:    true,
				DeprecationMessage: "This attribute will be removed in the next major version of " +
					"the provider.",
			},
			schemaUsersWithRole: schema.ListNestedAttribute{
				Description: "List of user ids, with role, to assign this organizational unit.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaUserId: schema.StringAttribute{
							Description: "Identifier of the user to assign to the organizational " +
								"unit.",
							Computed: true,
						},
						schemaAssignedRole: schema.StringAttribute{
							Description: "Identifier of the role to be associated with the user " +
								"when assigned to the organizational unit.",
							Computed: true,
						},
					},
				},
			},
		},
	}
}
