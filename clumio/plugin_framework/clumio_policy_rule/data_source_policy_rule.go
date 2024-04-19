// Copyright 2024. Clumio, Inc.

// This file holds the datasource implementation for the clumio_policy_rule Terraform datasource.
// This datasource is used to retrieve the policy rules within Clumio based on the specified
// attributes.

package clumio_policy_rule

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioPolicyRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioPolicyRuleDataSource{}
)

// clumioPolicyRuleDataSource is the struct backing the clumio_policy_rule Terraform datasource. It
// holds the Clumio API client and any other required state needed to manage sdkPolicyRules
// within Clumio.
type clumioPolicyRuleDataSource struct {
	name           string
	client         *common.ApiClient
	sdkPolicyRules sdkclients.PolicyRuleClient
}

// NewClumioPolicyRuleDataSource creates a new instance of clumioPolicyRuleDataSource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioPolicyRuleDataSource() datasource.DataSource {
	return &clumioPolicyRuleDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioPolicyRuleDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_policy_rule"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioPolicyRuleDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkPolicyRules = sdkclients.NewPolicyRuleClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioPolicyRuleDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioPolicyRuleDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readPolicyRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
