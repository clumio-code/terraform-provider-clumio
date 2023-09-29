// Copyright 2023. Clumio, Inc.
//
// This file contains the functions related to provider definition and initialization utilizing plugin framework.

package clumio_pf

import (
	"context"
	"fmt"
	"os"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_auto_user_provisioning_rule"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_auto_user_provisioning_setting"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_manual_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_organizational_unit"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_policy"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_policy_assignment"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_policy_rule"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_post_process_aws_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_post_process_kms"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_protection_group"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_role"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_user"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_wallet"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &clumioProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &clumioProvider{}
}

// clumioProvider is the provider implementation.
type clumioProvider struct{}

// clumioProviderModel maps provider schema data to a Go type.
type clumioProviderModel struct {
	ClumioApiToken                  types.String `tfsdk:"clumio_api_token"`
	ClumioApiBaseUrl                types.String `tfsdk:"clumio_api_base_url"`
	ClumioOrganizationalUnitContext types.String `tfsdk:"clumio_organizational_unit_context"`
}

// Metadata returns the provider type name.
func (p *clumioProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clumio"
}

// Schema defines the provider-level schema for configuration data.
func (p *clumioProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"clumio_api_token": schema.StringAttribute{
				MarkdownDescription: "The API token required to invoke Clumio APIs.",
				Optional:            true,
			},
			"clumio_api_base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for Clumio APIs. The following are the valid " +
					"values for clumio_api_base_url. Use the appropriate value depending" +
					" on the region for which your credentials were created:\n\n\t\t" +
					"us-west: https://us-west-2.api.clumio.com\n\n\t\t" +
					"us-east: https://us-east-1.api.clumio.com\n\n\t\t" +
					"canada:  https://ca-central-1.ca.api.clumio.com",
				Optional: true,
			},
			"clumio_organizational_unit_context": schema.StringAttribute{
				MarkdownDescription: "Organizational Unit context in which to create the" +
					" clumio resources. If not set, the resources will be created in" +
					" the context of the Global Organizational Unit. The value should" +
					" be the id of the Organizational Unit and not the name.",
				Optional: true,
			},
		},
	}
}

// Configure prepares a Clumio API client for data sources and resources.
func (p *clumioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Clumio client")

	// Retrieve provider data from configuration
	var config clumioProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ClumioApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("clumioApiToken"),
			"Unknown Clumio API Token",
			fmt.Sprintf(errorFmt, token, common.ClumioApiToken),
		)
	}

	if config.ClumioApiBaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("clumioApiBaseUrl"),
			"Unknown Clumio API Base URL",
			fmt.Sprintf(errorFmt, baseUrl, common.ClumioApiBaseUrl),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	clumioApiToken := os.Getenv(common.ClumioApiToken)
	clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
	clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)

	if !config.ClumioApiToken.IsNull() {
		clumioApiToken = config.ClumioApiToken.ValueString()
	}

	if !config.ClumioApiBaseUrl.IsNull() {
		clumioApiBaseUrl = config.ClumioApiBaseUrl.ValueString()
	}

	if !config.ClumioOrganizationalUnitContext.IsNull() {
		clumioOrganizationalUnitContext = config.ClumioOrganizationalUnitContext.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if clumioApiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("clumioApiToken"),
			"Missing Clumio API Token",
			fmt.Sprintf(errorFmt, token, common.ClumioApiToken),
		)
	}

	if clumioApiBaseUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("clumioApiBaseUrl"),
			"Missing Clumio API Username",
			fmt.Sprintf(errorFmt, baseUrl, common.ClumioApiBaseUrl),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Clumio client")

	client := &common.ApiClient{
		ClumioConfig: clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				userAgentHeader:            userAgentHeaderValue,
				clumioTfProviderVersionKey: clumioTfProviderVersionValue,
			},
		},
	}
	// Make the Clumio client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Clumio client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *clumioProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		clumio_role.NewClumioRoleDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *clumioProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		clumio_aws_connection.NewClumioAWSConnectionResource,
		clumio_post_process_aws_connection.NewPostProcessAWSConnectionResource,
		clumio_policy.NewPolicyResource,
		clumio_policy_assignment.NewPolicyAssignmentResource,
		clumio_policy_rule.NewPolicyRuleResource,
		clumio_protection_group.NewProtectionGroupResource,
		clumio_user.NewClumioUserResource,
		clumio_organizational_unit.NewClumioOrganizationalUnitResource,
		clumio_wallet.NewClumioWalletResource,
		clumio_post_process_kms.NewClumioPostProcessKmsResource,
		clumio_auto_user_provisioning_rule.NewAutoUserProvisioningRuleResource,
		clumio_auto_user_provisioning_setting.NewAutoUserProvisioningSettingResource,
		clumio_aws_manual_connection.NewAwsManualConnectionResource,
	}
}
