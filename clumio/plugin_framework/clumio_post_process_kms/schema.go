// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_post_process_kms Terraform resource.

package clumio_post_process_kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioPostProcessKmsResourceModel is the resource model for the clumio_post_process_kms Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type clumioPostProcessKmsResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Token                 types.String `tfsdk:"token"`
	AccountId             types.String `tfsdk:"account_id"`
	Region                types.String `tfsdk:"region"`
	RoleId                types.String `tfsdk:"role_id"`
	RoleArn               types.String `tfsdk:"role_arn"`
	RoleExternalId        types.String `tfsdk:"role_external_id"`
	CreatedMultiRegionCMK types.Bool   `tfsdk:"created_multi_region_cmk"`
	MultiRegionCMKKeyId   types.String `tfsdk:"multi_region_cmk_key_id"`
	TemplateVersion       types.Int64  `tfsdk:"template_version"`
}

// Schema defines the structure and constraints of the clumio_post_process_kms Terraform resource.
// Schema is a method on the clumioPostProcessKmsResource struct. It sets the schema for the
// clumio_post_process_kms Terraform resource.
func (r *clumioPostProcessKmsResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Post-Process Clumio KMS Resource used to post-process KMS in Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "The unique identifier of the post process kms.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaToken: schema.StringAttribute{
				Description: "Distinct 36-character token used to identify resources set up by " +
					"the Clumio BYOK template installation on the account being connected.",
				Required: true,
			},
			schemaAccountId: schema.StringAttribute{
				Description: "Identifier of the AWS account linked with Clumio.",
				Required:    true,
			},
			schemaRegion: schema.StringAttribute{
				Description: "Region of the AWS account linked with Clumio.",
				Required:    true,
			},
			schemaRoleId: schema.StringAttribute{
				Description: "Identifier of the IAM role to manage the customer-managed key.",
				Required:    true,
			},
			schemaRoleArn: schema.StringAttribute{
				Description: "The ARN of the IAM role to manage the customer-managed key.",
				Required:    true,
			},
			schemaRoleExternalId: schema.StringAttribute{
				Description: "Unique identifier Clumio uses to access the service role within " +
					"your account.",
				Required: true,
			},
			schemaCreatedMultiRegionCMK: schema.BoolAttribute{
				Description: "Indicates if a new customer-managed key was created.",
				Optional:    true,
			},
			schemaMultiRegionCMKKeyID: schema.StringAttribute{
				Description: "Identifier of the multi region customer-managed key.",
				Optional:    true,
			},
			schemaTemplateVersion: schema.Int64Attribute{
				Description: "Version of the BYOK template which was created.",
				Optional:    true,
			},
		},
	}
}
