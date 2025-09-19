// Copyright 2025. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_general_settings Terraform resource.

package clumio_general_settings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// generalSettingsResourceModel is the resource model for the clumio_general_settings
// Terraform resource. It represents the schema of the resource and the data it holds. This schema
// is used by customers to configure the resource and by the Clumio provider to read and write the
// resource.
type generalSettingsResourceModel struct {
	AutoLogoutDuration         types.Int64    `tfsdk:"auto_logout_duration"`
	IpAllowlist                []types.String `tfsdk:"ip_allowlist"`
	PasswordExpirationDuration types.Int64    `tfsdk:"password_expiration_duration"`
}

// Schema defines the structure and constraints of the clumio_general_settings Terraform
// resource. Schema is a method on the clumioGeneralSettings struct. It sets the schema
// for the clumio_general_settings Terraform resource, which is used to configure a report.
func (r *clumioGeneralSettings) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio general settings Resource used to manage general organizational " +
			"settings. These settings are persistent and cannot be deleted, only updated or reset.",
		Attributes: map[string]schema.Attribute{
			schemaAutoLogoutDuration: schema.Int64Attribute{
				Description: "The length of time before a user is logged out of the Clumio system" +
					" due to inactivity. Measured in seconds. The valid range is between 600 " +
					"seconds (10 minutes) and 3600 seconds (60 minutes). If not configured, the " +
					"value defaults to 900 seconds (15 minutes).",
				Optional:   true,
				Computed:   true,
				Validators: []validator.Int64{int64validator.Between(600, 3600)},
				Default:    int64default.StaticInt64(defaultAutoLogoutDuration),
			},
			schemaIPAllowlist: schema.SetAttribute{
				Description: "The designated range of IP addresses that are allowed to access the" +
					" Clumio REST API. API requests that originate from outside this list will be" +
					" blocked. The IP address of the server from which this request is being made" +
					" must be in this list; otherwise, the request will fail. Set the parameter " +
					"to individual IP addresses and/or a range of IP addresses in CIDR notation. " +
					"For example, [`193.168.1.0/24`, `193.172.1.1`]. If not configured, the value" +
					" defaults to [`0.0.0.0/0`] meaning all addresses will be allowed.",
				Optional:    true,
				ElementType: types.StringType,
			},
			schemaPasswordExpiration: schema.Int64Attribute{
				Description: "The length of time a user password is valid before it must be " +
					"changed. Measured in seconds. The valid range is between 2592000 seconds " +
					"(30 days) and 15552000 seconds (180 days). If not configured, the value " +
					"defaults to 7776000 seconds (90 days).",
				Optional:   true,
				Computed:   true,
				Validators: []validator.Int64{int64validator.Between(2592000, 15552000)},
				Default:    int64default.StaticInt64(defaultPasswordExpiration),
			},
		},
	}
}
