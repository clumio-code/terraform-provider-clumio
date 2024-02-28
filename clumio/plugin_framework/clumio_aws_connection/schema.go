// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_aws_connection Terraform resource.

package clumio_aws_connection

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

// clumioAWSConnectionResourceModel is the resource model for the clumio_aws_connection Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type clumioAWSConnectionResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	AccountNativeID      types.String `tfsdk:"account_native_id"`
	AWSRegion            types.String `tfsdk:"aws_region"`
	Description          types.String `tfsdk:"description"`
	OrganizationalUnitID types.String `tfsdk:"organizational_unit_id"`
	ConnectionStatus     types.String `tfsdk:"connection_status"`
	Token                types.String `tfsdk:"token"`
	Namespace            types.String `tfsdk:"namespace"`
	ClumioAWSAccountID   types.String `tfsdk:"clumio_aws_account_id"`
	ClumioAWSRegion      types.String `tfsdk:"clumio_aws_region"`
	ExternalID           types.String `tfsdk:"role_external_id"`
	DataPlaneAccountID   types.String `tfsdk:"data_plane_account_id"`
}

// Schema defines the structure and constraints of the clumio_aws_connection Terraform resource.
// Schema is a method on the clumioAWSConnectionResource struct. It sets the schema for the
// clumio_aws_connection Terraform resource, which is used to connect AWS accounts to Clumio. The
// schema defines various attributes such as the connection ID, AWS account ID, AWS region,
// description, etc. Some of these attributes are computed, meaning they are determined by Clumio at
// runtime, while others are required or optional inputs from the user.
func (r *clumioAWSConnectionResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for establishing a connection between AWS accounts and Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio AWS connection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "Identifier of the AWS account to be linked with Clumio.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "Region of the AWS account to be linked with Clumio.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaDescription: schema.StringAttribute{
				Description: "Brief description to denote details of the connection.",
				Optional:    true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "Identifier of the Clumio organizational unit associated with the " +
					"connection. If not provided, the connection will be associated with the " +
					"default organizational unit associated with the credentials used to create " +
					"the connection.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaConnectionStatus: schema.StringAttribute{
				Description: "Current state of the connection (e.g, `connecting`, `connected`, " +
					"`unlinked`, etc.)",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaToken: schema.StringAttribute{
				Description: "Distinct 36-character token used to identify resources set up by " +
					"the Clumio AWS template installation on the account being connected.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaNamespace: schema.StringAttribute{
				Description:        "K8S Namespace.",
				Computed:           true,
				DeprecationMessage: "This attribute will be removed in the next major version of the provider.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaClumioAwsAccountId: schema.StringAttribute{
				Description: "Identifier of the AWS account associated with Clumio. This " +
					"identifier is provided so that access to the service role for Clumio can be " +
					"restricted to just this account.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaClumioAwsRegion: schema.StringAttribute{
				Description: "Region of the AWS account associated with Clumio.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaExternalId: schema.StringAttribute{
				Description: "Unique identifier Clumio uses to access the service role within " +
					"your account.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaDataPlaneAccountId: schema.StringAttribute{
				Description: "Identifier of the AWS account data plane within Clumio.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
