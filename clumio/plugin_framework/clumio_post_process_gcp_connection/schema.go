// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_post_process_gcp_connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioPostProcessGCPConnectionResourceModel is the resource model for the clumio_post_process_gcp_connection Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type clumioPostProcessGCPConnectionResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ProjectID           types.String `tfsdk:"project_id"`
	ProjectName         types.String `tfsdk:"project_name"`
	ProjectNumber       types.String `tfsdk:"project_number"`
	Token               types.String `tfsdk:"token"`
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	WifPoolId           types.String `tfsdk:"wif_pool_id"`
	WifProviderId       types.String `tfsdk:"wif_provider_id"`
	ConfigVersion       types.String `tfsdk:"config_version"`
	ProtectGcsVersion   types.String `tfsdk:"protect_gcs_version"`
	Properties          types.Map    `tfsdk:"properties"`
}

// Schema defines the structure and constraints of the clumio_post_process_gcp_connection Terraform
// resource. Schema is a method on the clumioPostProcessGCPConnectionResource struct. It sets the schema
// for the clumio_post_process_gcp_connection Terraform resource.
func (r *clumioPostProcessGCPConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for post processing connection between GCP projects and Clumio.",
		MarkdownDescription: "> ⚠️ **Beta Resource**\n>\n> This resource handles post-processing for connections between GCP projects and Clumio.\n> It is currently in **beta** and available only to select customers.\n> Behavior, schema, and APIs may change in future releases.\n>",
		Attributes: map[string]schema.Attribute{
			schemaID: schema.StringAttribute{
				Description: "Unique identifier of the connection",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaProjectID: schema.StringAttribute{
				Description: "The user-assigned ID of the GCP project associated with the connection.",
				Required:    true,
			},
			schemaProjectName: schema.StringAttribute{
				Description: "The user-friendly name of the GCP project associated with the connection.",
				Required:    true,
			},
			schemaProjectNumber: schema.StringAttribute{
				Description: "The GCP-assigned numeric INT64 project number associated with the connection.",
				Required:    true,
			},
			schemaToken: schema.StringAttribute{
				Description: "The 36-character Clumio GCP integration token used to identify the installation " +
					"of the Clumio GCP integration resources in the project.",
				Required: true,
			},
			schemaServiceAccountEmail: schema.StringAttribute{
				Description: "The email address of the GCP service account created for this connection.",
				Required:    true,
			},
			schemaWifPoolId: schema.StringAttribute{
				Description: "The Workload Identity Federation Pool ID created for this connection.",
				Required:    true,
			},
			schemaWifProviderId: schema.StringAttribute{
				Description: "The Workload Identity Federation Provider ID created for this connection.",
				Required:    true,
			},
			schemaConfigVersion: schema.StringAttribute{
				Description: "Clumio Config version. May be a single number or major.minor (e.g., 1, 1.0, 2.5, 10.11).",
				Required:    true,
			},
			schemaProtectGcsVersion: schema.StringAttribute{
				Description: "Clumio Config version for GCS. May be a single number or major.minor (e.g., 1, 1.0, 2.5, 10.11).",
				Optional:    true,
			},
			schemaProperties: schema.MapAttribute{
				Description: "A map to pass in additional information to be consumed " +
					"by Clumio Post Processing",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}
