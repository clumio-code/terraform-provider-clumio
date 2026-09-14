// Copyright 2026. Clumio, Inc.

// This file holds the datasource implementation for the clumio_iceberg_tables Terraform datasource.
// This datasource is used to retrieve the Clumio Iceberg tables based on the specified attributes.

package clumio_iceberg_tables

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioIcebergTablesDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioIcebergTablesDataSource{}
)

// clumioIcebergTablesDataSource is the struct backing the clumio_iceberg_tables Terraform
// datasource. It holds the Clumio API client and any other required state needed to read the
// Iceberg tables within Clumio.
type clumioIcebergTablesDataSource struct {
	name                 string
	client               *common.ApiClient
	icebergTableClient   sdkclients.IcebergTableClient
	awsEnvironmentClient sdkclients.AWSEnvironmentClient
}

// NewClumioIcebergTablesDataSource creates a new instance of clumioIcebergTablesDataSource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioIcebergTablesDataSource() datasource.DataSource {
	return &clumioIcebergTablesDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioIcebergTablesDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_iceberg_tables"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioIcebergTablesDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.icebergTableClient = sdkclients.NewIcebergTableClient(r.client.ClumioConfig)
	r.awsEnvironmentClient = sdkclients.NewAWSEnvironmentClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioIcebergTablesDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioIcebergTablesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readIcebergTables(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
