// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_auto_user_provisioning_setting Terraform resource.

package clumio_auto_user_provisioning_setting

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// autoUserProvisioningSettingResourceModel is the resource model for the
// clumio_auto_user_provisioning_setting Terraform resource. It represents the schema of the
// resource and the data it holds. This schema is used by customers to configure the resource and
// by the Clumio provider to read and write the resource.
type autoUserProvisioningSettingResourceModel struct {
	ID        types.String `tfsdk:"id"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
}

// Schema defines the structure and constraints of the clumio_auto_user_provisioning_setting
// Terraform resource. Schema is a method on the autoUserProvisioningSettingResource struct. It sets
// the schema for the clumio_auto_user_provisioning_setting Terraform resource, which is used to
// enable/disable the auto user provisioning feature.
func (r *autoUserProvisioningSettingResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Auto User Provisioning Setting Resource used to enable or disable the" +
			" auto user provisioning feature.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the auto user provisioning setting.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaIsEnabled: schema.BoolAttribute{
				Description: "Whether auto user provisioning is enabled or not.",
				Required:    true,
			},
		},
	}
}
