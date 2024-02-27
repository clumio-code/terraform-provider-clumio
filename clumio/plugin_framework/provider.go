// Copyright 2023. Clumio, Inc.

// This file contains the functions related to provider definition and initialization utilizing the
// plugin framework.

package clumio_pf

import (
	"context"
	"os"
	"strings"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_auto_user_provisioning_rule"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_auto_user_provisioning_setting"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_manual_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_manual_connection_resources"
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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the following Provider interface.
var (
	_ provider.Provider = &clumioProvider{}
)

// clumioProvider is the struct backing the Clumio Provider for Terraform.
type clumioProvider struct{}

// New creates a new instance of clumioProvider.
func New() provider.Provider {
	return &clumioProvider{}
}

// Metadata returns the provider type name.
func (p *clumioProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clumio"
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

	// If given, provider attributes must be "known" values. In other words, their values cannot be
	// computed from other values in the configuration.
	if config.ClumioApiToken.IsUnknown() {
		attribute := path.Root("clumioApiToken")
		summary := "Unknown Clumio API Token"
		detail := "Value must not be computed from other values in the configuration."
		resp.Diagnostics.AddAttributeError(attribute, summary, detail)
	}
	if config.ClumioApiBaseUrl.IsUnknown() {
		attribute := path.Root("clumioApiBaseUrl")
		summary := "Unknown Clumio API Base URL"
		detail := "Value must not be computed from other values in the configuration."
		resp.Diagnostics.AddAttributeError(attribute, summary, detail)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Provider attributes can be set statically in the configuration or using environment. If
	// statically set, the value is available in the configuration. If set using environment, the
	// value is available in documented environment variables. Statically set values take precedence
	// over environment variables.
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

	// Ensure that all required values are set. If not, return an error.
	if clumioApiToken == "" {
		attribute := path.Root("clumioApiToken")
		summary := "Missing Clumio API Token"
		detail := "Value is required."
		resp.Diagnostics.AddAttributeError(attribute, summary, detail)
	}
	if clumioApiBaseUrl == "" {
		attribute := path.Root("clumioApiBaseUrl")
		summary := "Missing Clumio API Token"
		detail := "Value is required."
		resp.Diagnostics.AddAttributeError(attribute, summary, detail)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure that the base URL does not end with a slash.
	clumioApiBaseUrl = strings.TrimRight(clumioApiBaseUrl, "/")

	// Create the Clumio API client and make it available to instances of DataSource and Resource
	// types in their Configure methods.
	tflog.Debug(ctx, "Creating Clumio client")
	client := &common.ApiClient{
		ClumioConfig: clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				userAgentHeader:               userAgentHeaderValue,
				clumioTfProviderVersionHeader: clumioTfProviderVersionHeaderValue,
			},
		},
	}
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured Clumio client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider. Any new data source should be
// added here.
func (p *clumioProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		clumio_role.NewClumioRoleDataSource,
		clumio_aws_manual_connection_resources.NewAwsManualConnectionResourcesDataSource,
	}
}

// Resources defines the resources implemented in the provider. Any new resource should be added
// here.
func (p *clumioProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		clumio_aws_connection.NewClumioAWSConnectionResource,
		clumio_post_process_aws_connection.NewPostProcessAWSConnectionResource,
		clumio_policy.NewPolicyResource,
		clumio_policy_assignment.NewPolicyAssignmentResource,
		clumio_policy_rule.NewPolicyRuleResource,
		clumio_protection_group.NewClumioProtectionGroupResource,
		clumio_user.NewClumioUserResource,
		clumio_organizational_unit.NewClumioOrganizationalUnitResource,
		clumio_wallet.NewClumioWalletResource,
		clumio_post_process_kms.NewClumioPostProcessKmsResource,
		clumio_auto_user_provisioning_rule.NewAutoUserProvisioningRuleResource,
		clumio_auto_user_provisioning_setting.NewAutoUserProvisioningSettingResource,
		clumio_aws_manual_connection.NewClumioAWSManualConnectionResource,
	}
}
