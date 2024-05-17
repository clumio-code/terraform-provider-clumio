// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_policy_assignment Terraform resource.

package clumio_policy_assignment

import (
	"context"

	validators "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// policyAssignmentResourceModel is the resource model for the clumio_policy_assignment Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type policyAssignmentResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	EntityID             types.String `tfsdk:"entity_id"`
	EntityType           types.String `tfsdk:"entity_type"`
	PolicyID             types.String `tfsdk:"policy_id"`
	OrganizationalUnitID types.String `tfsdk:"organizational_unit_id"`
}

// Schema defines the structure and constraints of the clumio_policy_assignment Terraform resource.
// Schema is a method on the clumioPolicyAssignmentResource struct. It sets the schema for the
// clumio_policy_assignment Terraform resource, which is used to assign a policy to an entity.
func (r *clumioPolicyAssignmentResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Assignment Resource used to assign (or unassign)" +
			" policies.\n\n NOTE: Currently policy assignment is supported only for" +
			" entity type \"protection_group\".",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the policy assignment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaEntityId: schema.StringAttribute{
				Description:   "Identifier of the resource to which the policy will be assigned.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			schemaEntityType: schema.StringAttribute{
				Description: "Type of resource to which the policy will be assigned. " +
					"Only `protection_group` is currently supported.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					validators.OneOf(entityTypeProtectionGroup),
				},
			},
			schemaPolicyId: schema.StringAttribute{
				Description: "Identifier of the Clumio policy to be assigned.",
				Required:    true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "Identifier of the Clumio organizational unit associated with the " +
					"resource for which the policy will be assigned. If not provided, the resource " +
					"will be assumed to be in the default organizational unit associated with the " +
					"credentials used to create the assignment.",
				Optional: true,
				Computed: true,
				DeprecationMessage: "Use the provider schema attribute " +
					"clumio_organizational_unit_context to create the resource in the context of " +
					"an Organizational Unit.",
				Validators: []validator.String{
					validators.LengthAtLeast(1),
				},
			},
		},
	}
}
