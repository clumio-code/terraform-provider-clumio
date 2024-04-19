// Copyright 2024. Clumio, Inc.

// This file holds the datasource implementation for the clumio_s3_bucket Terraform datasource. This
// datasource is used to retrieve the Clumio s3 buckets based on the specified attributes.

package clumio_s3_bucket

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioS3BucketDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioS3BucketDataSource{}
)

// clumioS3BucketDataSource is the struct backing the clumio_protection-group Terraform datasource.
// It holds the Clumio API client and any other required state needed to manage s3BucketClient
// within Clumio.
type clumioS3BucketDataSource struct {
	name           string
	client         *common.ApiClient
	s3BucketClient sdkclients.S3BucketClient
}

// NewClumioS3BucketDataSource creates a new instance of clumioS3BucketDataSource. Its attributes
// are initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewClumioS3BucketDataSource() datasource.DataSource {
	return &clumioS3BucketDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioS3BucketDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_s3_bucket"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioS3BucketDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.s3BucketClient = sdkclients.NewS3BucketClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioS3BucketDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioS3BucketDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readS3Bucket(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
