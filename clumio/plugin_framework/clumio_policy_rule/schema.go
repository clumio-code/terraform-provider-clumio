// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_policy_rule Terraform resource.

package clumio_policy_rule

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

// policyRuleResourceModel is the resource model for the clumio_policy_rule Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type policyRuleResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Condition            types.String `tfsdk:"condition"`
	BeforeRuleID         types.String `tfsdk:"before_rule_id"`
	PolicyID             types.String `tfsdk:"policy_id"`
	OrganizationalUnitID types.String `tfsdk:"organizational_unit_id"`
}

// Schema defines the structure and constraints of the clumio_policy_rule Terraform resource.
// Schema is a method on the policyRuleResource struct. It sets the schema for the
// clumio_policy_rule Terraform resource.
func (r *policyRuleResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Rule Resource used to determine how" +
			" a policy should be assigned to assets.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the policy rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the policy rule.",
				Required:    true,
			},
			schemaCondition: schema.StringAttribute{
				Description: "The condition of the policy rule. Possible conditions include:\n\t" +
					"1) `entity_type` is required and supports `$eq` and `$in` filters. " +
					"`entity_type` must be one of `aws_rds_instance`, `aws_ebs_volume`, " +
					"`aws_ec2_instance`, `aws_dynamodb_table` or `aws_rds_cluster`.\n\t" +
					"2) `aws_account_native_id` and `aws_region` are optional and both support " +
					"`$eq` and `$in` filters.\n\t" +
					"3) `aws_tag` is optional and supports `$eq`, `$in`, `$all`, and `$contains` " +
					"filters.",
				Required: true,
			},
			schemaBeforeRuleId: schema.StringAttribute{
				Description: "The policy rule ID before which this policy rule should be " +
					"inserted. An empty value will set the rule to have lowest priority. " +
					"NOTE: If in the Global Organizational Unit, rules can also be prioritized " +
					"against two virtual rules maintained by the system: `asset-level-rule` and " +
					"`child-ou-rule`. `asset-level-rule` corresponds to the priority of Direct " +
					"Assignments (when a policy is applied directly to an asset) whereas " +
					"`child-ou-rule` corresponds to the priority of rules created by child " +
					"organizational units.",
				Required: true,
			},
			schemaPolicyId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the policy. ",
				Required:    true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the organizational unit" +
					" to use as the context for assigning the policy.",
				Optional: true,
				Computed: true,
				DeprecationMessage: "Use the provider schema attribute " +
					"clumio_organizational_unit_context to create the resource in the context of " +
					"an Organizational Unit.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}
