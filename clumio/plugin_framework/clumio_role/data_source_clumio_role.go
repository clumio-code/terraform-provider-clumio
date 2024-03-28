// Copyright 2023. Clumio, Inc.

// This file holds the datasource implementation for the clumio_role Terraform datasource. This
// datasource is used to manage roles within Clumio

package clumio_role

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioRoleDataSource{}
)

// clumioRoleDataSource is the struct backing the clumio_role Terraform datasource. It holds the
// Clumio API client and any other required state needed to manage roles within Clumio.
type clumioRoleDataSource struct {
	name   string
	client *common.ApiClient
	roles  sdkclients.RoleClient
}

// NewClumioRoleDataSource creates a new instance of clumioRoleDataSource. Its attributes are
// initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewClumioRoleDataSource() datasource.DataSource {
	return &clumioRoleDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioRoleDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_role"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioRoleDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.roles = sdkclients.NewRoleClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioRoleDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioRoleDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readRole(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
