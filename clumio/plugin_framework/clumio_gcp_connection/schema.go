// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_gcp_connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioGCPConnectionResourceModel is the resource model for the clumio_gcp_connection Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type clumioGCPConnectionResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	ClumioControlPlaneId   types.String `tfsdk:"clumio_control_plane_id"`
	ClumioControlPlaneRole types.String `tfsdk:"clumio_control_plane_role"`
	ProjectID              types.String `tfsdk:"project_id"`
	Description            types.String `tfsdk:"description"`
	Regions                types.List   `tfsdk:"regions"`
	Token                  types.String `tfsdk:"token"`
}

// Schema defines the structure and constraints of the clumio_gcp_connection Terraform resource.
// Schema is a method on the clumioGCPConnectionResource struct. It sets the schema for the
// clumio_gcp_connection Terraform resource, which is used to connect GCP projects to Clumio.
// Some of these attributes are computed, meaning they are determined by Clumio at
// runtime, while others are required or optional inputs from the user.
func (r *clumioGCPConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for establishing a connection between GCP projects and Clumio.",
		MarkdownDescription: "> ⚠️ **Beta Resource**\n>\n> This resource establishes a connection between GCP projects and Clumio.\n> It is currently in **beta** and available only to select customers.\n> Behavior, schema, and APIs may change in future releases.\n>",
		Attributes: map[string]schema.Attribute{
			schemaID: schema.StringAttribute{
				Description: "Unique identifier of the connection",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaClumioControlPlaneId: schema.StringAttribute{
				Description: "Identifier for the Clumio Control Plan. This " +
					"identifier is provided so that access to the service role for Clumio can be " +
					"restricted to just this control plane.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaClumioControlPlaneRole: schema.StringAttribute{
				Description: "Identifier for the Clumio Control Role. This " +
					"identifier will be federated into GCP",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaToken: schema.StringAttribute{
				Description: "The 36-character Clumio GCP integration token used to identify the " +
					"installation of the Clumio GCP integration resources in the project.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaProjectId: schema.StringAttribute{
				Description: "The user-assigned ID of the GCP project associated with the connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					// If the value of this attribute changes, Terraform will destroy and recreate the resource.
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaDescription: schema.StringAttribute{
				Description: "The user defined description for the connection.",
				Optional:    true,
			},
			schemaRegions: schema.ListAttribute{
				Description: "The GCP regions to be used for inventory.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}
