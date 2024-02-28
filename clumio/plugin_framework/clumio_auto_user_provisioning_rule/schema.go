// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_auto_user_provisioning_rule Terraform resource.

package clumio_auto_user_provisioning_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// autoUserProvisioningRuleResourceModel is the resource model for the
// clumio_auto_user_provisioning_rule Terraform resource. It represents the schema of the resource
// and the data it holds. This schema is used by customers to configure the resource and by the
// Clumio provider to read and write the resource.
type autoUserProvisioningRuleResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Condition             types.String `tfsdk:"condition"`
	RoleID                types.String `tfsdk:"role_id"`
	OrganizationalUnitIDs types.Set    `tfsdk:"organizational_unit_ids"`
}

// Schema defines the structure and constraints of the clumio_auto_user_provisioning_rule Terraform
// resource. Schema is a method on the autoUserProvisioningRuleResource struct. It sets the schema
// for the clumio_auto_user_provisioning_rule Terraform resource, which is used to create and manage
// wallets.
func (r *autoUserProvisioningRuleResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Auto User Provisioning Rule Resource used to determine " +
			"the Role and Organizational Units to be assigned to a user based on their groups.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the auto user provisioning rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the auto user provisioning rule.",
				Required:    true,
			},
			schemaCondition: schema.StringAttribute{
				Description: "The condition of the auto user provisioning rule. Possible conditions" +
					" include:\n" +
					"\t1) `This group` - User must belong to the specified group\n" +
					"\t2) `ANY of these groups` - User must belong to at least one of the specified" +
					" groups\n" +
					"\t3) `ALL of these groups` - User must belong to all the specified groups\n" +
					"\t4) `Group CONTAINS this keyword` - User's group must contain the specified" +
					" keyword\n" +
					"\t5) `Group CONTAINS ANY of these keywords` - User's group must contain at" +
					" least one of the specified keywords\n" +
					"\t6) `Group CONTAINS ALL of these keywords` - User's group must contain all" +
					" the specified keywords\n",
				Required: true,
			},
			schemaRoleId: schema.StringAttribute{
				Description: "Identifier of the Clumio role to be assigned to the user.",
				Required:    true,
			},
			schemaOrganizationalUnitIds: schema.SetAttribute{
				Description: "List of Clumio organizational unit identifiers to be assigned to the" +
					" user.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}
