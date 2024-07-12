// Copyright 2024. Clumio, Inc.

// This file holds the datasource implementation for the clumio_protection_group_asset Terraform
// datasource. This datasource is used to retrieve the Clumio protection group asset based on the
// specified attributes.

package clumio_protection_group_asset

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioProtectionGroupAssetDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioProtectionGroupAssetDataSource{}
)

// clumioProtectionGroupAssetDataSource is the struct backing the clumio_protection_group_asset
// Terraform datasource. It holds the Clumio API client and any other required state needed to
// manage s3AssetsClient within Clumio.
type clumioProtectionGroupAssetDataSource struct {
	name           string
	client         *common.ApiClient
	s3AssetsClient sdkclients.ProtectionGroupS3AssetsClient
}

// NewClumioProtectionGroupAssetDataSource creates a new instance of clumioProtectionGroupAssetDataSource. Its attributes
// are initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewClumioProtectionGroupAssetDataSource() datasource.DataSource {
	return &clumioProtectionGroupAssetDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioProtectionGroupAssetDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_protection_group_asset"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioProtectionGroupAssetDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.s3AssetsClient = sdkclients.NewProtectionGroupS3AssetsClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioProtectionGroupAssetDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioProtectionGroupAssetDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readProtectionGroupAsset(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
