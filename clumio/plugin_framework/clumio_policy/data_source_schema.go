// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_policy Terraform datasource.

package clumio_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioPolicyDataSourceModel is the datasource model for the clumio_policy Terraform datasource. It
// represents the schema of the datasource and the data it holds. This schema is used by customers
// to configure the datasource and by the Clumio provider to read and write the datasource.
type clumioPolicyDataSourceModel struct {
	Name             types.String `tfsdk:"name"`
	ActivationStatus types.String `tfsdk:"activation_status"`
	OperationTypes   types.Set    `tfsdk:"operation_types"`
	Policies         types.Set    `tfsdk:"policies"`
}

// Schema defines the structure and constraints of the clumio_policy Terraform datasource. Schema is
// a method on the clumioPolicyDataSource struct. It sets the schema for the clumio_policy Terraform
// datasource. The schema defines various attributes such as the policy name, operation_types,
// activation_status, etc, some of which are computed, meaning they are determined by Clumio at
// runtime, whereas 'name', 'operation_types' and 'activation_status' attributes are used to
// determine the Clumio policies to retrieve.
func (r *clumioPolicyDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaName: schema.StringAttribute{
				Description: "The name of the policy to be included in the read policies query.",
				Optional:    true,
			},
			schemaOperationTypes: schema.SetAttribute{
				Description: "Operation types to be included in the read policies query.",
				Optional:    true,
				ElementType: types.StringType,
			},
			schemaActivationStatus: schema.StringAttribute{
				Description: "Activation status to be included in the query filter. Valid values" +
					" are activated/deactivated.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(activationStatusActivated, activationStatusDectivated),
				},
			},
			schemaPolicies: schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the policy.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "The name of the policy.",
							Computed:    true,
						},
						schemaOperationTypes: schema.SetAttribute{
							Description: "Operation types supported by the policy.",
							Computed:    true,
							ElementType: types.StringType,
						},
						schemaActivationStatus: schema.StringAttribute{
							Description: "Activation status of the policy.",
							Computed:    true,
						},
						schemaTimezone: schema.StringAttribute{
							Description: "The time zone for the policy, in IANA format.",
							Computed:    true,
						},
						schemaOrganizationalUnitId: schema.StringAttribute{
							Description: "Identifier of the Clumio organizational unit associated " +
								"with the policy.",
							Computed: true,
						},
					},
				},
				Computed:    true,
				Description: "List of policies which matched the query criteria.",
			},
		},
		Description: "clumio_policy data source is used to retrieve details of the policies for use" +
			" in other resources. At least one of 'name', 'activation_status' or 'operation_types'" +
			" must be specified in the config.",
	}
}

// ConfigValidators to check if at least one of name, operation_types or activation_status is specified.
func (r *clumioPolicyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot(schemaName),
			path.MatchRoot(schemaOperationTypes),
			path.MatchRoot(schemaActivationStatus),
		),
	}
}
