// Copyright 2026. Clumio, Inc.

// This file holds the data source implementation for the clumio_gcs_protection_group Terraform data
// source. This data source is used to retrieve a GCP protection group based on the specified
// attributes.

package clumio_gcs_protection_group

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioGCSProtectionGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioGCSProtectionGroupDataSource{}
)

// clumioGCSProtectionGroupDataSource is the struct backing the clumio_gcs_protection_group
// Terraform data source. It holds the Clumio API client and any other required state needed to
// query GCP protection groups within Clumio.
type clumioGCSProtectionGroupDataSource struct {
	name                string
	client              *common.ApiClient
	sdkProtectionGroups sdkclients.GcpProtectionGroupClient
}

// NewClumioGCSProtectionGroupDataSource creates a new instance of
// clumioGCSProtectionGroupDataSource. Its attributes are initialized later by Terraform via
// Metadata and Configure once the Provider is initialized.
func NewClumioGCSProtectionGroupDataSource() datasource.DataSource {
	return &clumioGCSProtectionGroupDataSource{}
}

// Metadata returns the name of the data source type. This is used by Terraform configurations to
// instantiate the data source.
func (r *clumioGCSProtectionGroupDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_gcs_protection_group"
	resp.TypeName = r.name
}

// Configure sets up the data source with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioGCSProtectionGroupDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkProtectionGroups = sdkclients.NewGcpProtectionGroupClient(r.client.ClumioConfig)
}

// Read retrieves the data source from the Clumio API and sets the Terraform state.
func (r *clumioGCSProtectionGroupDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform configuration.
	var state clumioGCSProtectionGroupDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readProtectionGroup(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
