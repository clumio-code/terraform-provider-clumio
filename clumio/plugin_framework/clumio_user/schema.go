// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_user Terraform resource.

package clumio_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// roleForOrganizationalUnitModel is the model mapping to an assigned role within an organizational
// unit.
type roleForOrganizationalUnitModel struct {
	RoleId                types.String `tfsdk:"role_id"`
	OrganizationalUnitIds types.Set    `tfsdk:"organizational_unit_ids"`
}

// clumioUserResourceModel is the resource model for the clumio_user Terraform resource. It
// represents the schema of the resource and the data it holds. This schema is used by customers to
// configure the resource and by the Clumio provider to read and write the resource.
type clumioUserResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Email                      types.String `tfsdk:"email"`
	FullName                   types.String `tfsdk:"full_name"`
	AssignedRole               types.String `tfsdk:"assigned_role"`
	OrganizationalUnitIds      types.Set    `tfsdk:"organizational_unit_ids"`
	AccessControlConfiguration types.Set    `tfsdk:"access_control_configuration"`
	Inviter                    types.String `tfsdk:"inviter"`
	IsConfirmed                types.Bool   `tfsdk:"is_confirmed"`
	IsEnabled                  types.Bool   `tfsdk:"is_enabled"`
	LastActivityTimestamp      types.String `tfsdk:"last_activity_timestamp"`
	OrganizationalUnitCount    types.Int64  `tfsdk:"organizational_unit_count"`
}

// Schema defines the structure and constraints of the clumio_user Terraform resource. Schema is a
// method on the clumioUserResource struct. It sets the schema for the clumio_user Terraform
// resource, which is used to create and manage users within Clumio. The schema defines various
// attributes such as the user ID, email, name, etc. Some of these attributes are computed, meaning
// they are determined by Clumio at runtime, while others are required or optional inputs from the
// user.
func (r *clumioUserResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Clumio User Resource to create and manage users in Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the user.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaEmail: schema.StringAttribute{
				Description: "The email address of the user to be added to Clumio.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaFullName: schema.StringAttribute{
				Description: "The full name of the user to be added to Clumio. For example, enter " +
					"the user's first name and last name. The name appears in the User Management " +
					"screen and in the body of the  email invitation.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAssignedRole: schema.StringAttribute{
				Description: "Identifier of the role to assign to the user.",
				Optional:    true,
				Computed:    true,
				DeprecationMessage: "Configure access_control_configuration instead. This attribute will" +
					" be removed in the next major version of the provider.",
			},
			schemaOrganizationalUnitIds: schema.SetAttribute{
				Description: "Identifiers of the organizational units  to be assigned to the user. The" +
					" Global Organizational Unit ID is \"00000000-0000-0000-0000-000000000000\"",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				DeprecationMessage: "Configure access_control_configuration instead. This attribute will" +
					" be removed in the next major version of the provider.",
			},
			schemaAccessControlConfiguration: schema.SetNestedAttribute{
				Description: "Identifiers of the organizational units, along with the identifier of the" +
					" role, to be assigned to the user.",
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaRoleId: schema.StringAttribute{
							Description: "Identifier of the role to assign to the user.",
							Optional:    true,
							Computed:    true,
						},
						schemaOrganizationalUnitIds: schema.SetAttribute{
							Description: "Identifiers of the organizational units to be assigned to the user." +
								"The Global Organizational Unit ID is \"00000000-0000-0000-0000-000000000000\"",
							Optional:    true,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			schemaInviter: schema.StringAttribute{
				Description: "Identifier of the user who sent the email invitation.",
				Computed:    true,
			},
			schemaIsConfirmed: schema.BoolAttribute{
				Description: "Determines whether the user has activated their Clumio account. If true," +
					" the user has activated the account.",
				Computed: true,
			},
			schemaIsEnabled: schema.BoolAttribute{
				Description: "Determines whether the user is enabled (in Activated or Invited status) in" +
					" Clumio. If true, the user is in Activated or Invited status in Clumio. Users in" +
					"Activated status can log in to Clumio. Users in Invited status have been invited to log" +
					"in to Clumio via an email invitation and the invitation is pending acceptance from the" +
					" user. If false, the user has been manually suspended and cannot log in to Clumio until" +
					"another Clumio user reactivates the account.",
				Computed: true,
			},
			schemaLastActivityTimestamp: schema.StringAttribute{
				Description: "The timestamp of when when the user was last active in the Clumio system." +
					" Represented in RFC-3339 format.",
				Computed: true,
			},
			schemaOrganizationalUnitCount: schema.Int64Attribute{
				Description: "The number of organizational units accessible to the user.",
				Computed:    true,
			},
		},
	}
}
