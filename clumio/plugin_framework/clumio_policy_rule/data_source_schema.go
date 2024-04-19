// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_policy_rule Terraform datasource.

package clumio_policy_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioPolicyRuleDataSourceModel is the datasource model for the clumio_policy_rule Terraform
// datasource. It represents the schema of the datasource and the data it holds. This schema is used
// by customers to configure the datasource and by the Clumio provider to read and write the
// datasource.
type clumioPolicyRuleDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	PolicyId    types.String `tfsdk:"policy_id"`
	PolicyRules types.Set    `tfsdk:"policy_rules"`
}

// Schema defines the structure and constraints of the clumio_policy_rule Terraform datasource.
// Schema is a method on the clumioPolicyRuleDataSource struct. It sets the schema for the
// clumio_policy_rule Terraform datasource, which fetches the policy_rules. The schema
// defines various attributes such as the name, policy_id, etc, some of which are computed,
// meaning they are determined by Clumio at runtime, whereas 'name' and 'policy_id' attributes are
// used to determine the Clumio policy rules to retrieve.
func (r *clumioPolicyRuleDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaName: schema.StringAttribute{
				Description: "The name of the policy rule to filter in the list of policy rules" +
					" returned by the API.",
				Optional: true,
			},
			schemaPolicyId: schema.StringAttribute{
				Description: "Unique identifier of the policy to filter in the list of policy rules" +
					" returned by the API.",
				Optional: true,
			},
			schemaPolicyRules: schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the policy rule.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "The name of the policy rule.",
							Computed:    true,
						},
						schemaPolicyId: schema.StringAttribute{
							Description: "Unique identifier of the policy associated with the " +
								"policy rule.",
							Computed: true,
						},
						schemaBeforeRuleId: schema.StringAttribute{
							Description: "The policy rule ID before which this policy rule should be " +
								"executed.",
							Computed: true,
						},
						schemaCondition: schema.StringAttribute{
							Description: "The condition of the policy rule. Possible conditions " +
								"include: " +
								"1) `entity_type` is required and supports `$eq` and `$in` filters." +
								"2) `aws_account_native_id` and `aws_region` are optional and both" +
								" support `$eq` and `$in` filters. " +
								"3) `aws_tag` is optional and supports `$eq`, `$in`, `$all`, and " +
								"`$contains` filters.",
							Computed: true,
						},
					},
				},
				Computed:    true,
				Description: "List of policies which matched the query criteria.",
			},
		},
		Description: "clumio_policy_rule data source is used to retrieve details of the policy rules" +
			" for use in other resources.",
	}
}
