// Copyright 2024. Clumio, Inc.

// This file holds the datasource implementation for the clumio_organizational_unit Terraform
// datasource. This datasource is used to retrieve the Clumio organizational units based on the
// specified attributes.

package clumio_organizational_unit

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioOrganizationalUnitDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioOrganizationalUnitDataSource{}
)

// clumioOrganizationalUnitDataSource is the struct backing the clumio_protection-group Terraform
// datasource. It holds the Clumio API client and any other required state needed to manage
// organizationalUnitClient within Clumio.
type clumioOrganizationalUnitDataSource struct {
	name                     string
	client                   *common.ApiClient
	organizationalUnitClient sdkclients.OrganizationalUnitClient
}

// NewClumioOrganizationalUnitDataSource creates a new instance of clumioOrganizationalUnitDataSource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioOrganizationalUnitDataSource() datasource.DataSource {
	return &clumioOrganizationalUnitDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioOrganizationalUnitDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_organizational_unit"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioOrganizationalUnitDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.organizationalUnitClient = sdkclients.NewOrganizationalUnitClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioOrganizationalUnitDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioOrganizationalUnitDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readOrganizationalUnit(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
