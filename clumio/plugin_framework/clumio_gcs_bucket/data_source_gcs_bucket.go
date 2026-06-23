// Copyright 2026. Clumio, Inc.

// This file holds the data source implementation for the clumio_gcs_bucket Terraform data
// source. This data source is used to retrieve the Clumio GCS buckets based on the specified
// attributes.

package clumio_gcs_bucket

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioGCSBucketDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioGCSBucketDataSource{}
)

// clumioGCSBucketDataSource is the struct backing the clumio_gcs_bucket Terraform data
// source. It holds the Clumio API client and any other required state needed to query GCS buckets
// within Clumio.
type clumioGCSBucketDataSource struct {
	name            string
	client          *common.ApiClient
	gcsBucketClient sdkclients.GcpGcsBucketClient
}

// NewClumioGCSBucketDataSource creates a new instance of clumioGCSBucketDataSource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioGCSBucketDataSource() datasource.DataSource {
	return &clumioGCSBucketDataSource{}
}

// Metadata returns the name of the data source type. This is used by Terraform configurations to
// instantiate the data source.
func (r *clumioGCSBucketDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_gcs_bucket"
	resp.TypeName = r.name
}

// Configure sets up the data source with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioGCSBucketDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.gcsBucketClient = sdkclients.NewGcpGcsBucketClient(r.client.ClumioConfig)
}

// Read retrieves the data source from the Clumio API and sets the Terraform state.
func (r *clumioGCSBucketDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform configuration.
	var state clumioGCSBucketDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readGCSBucket(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
