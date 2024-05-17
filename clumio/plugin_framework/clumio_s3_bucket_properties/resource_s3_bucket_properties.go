// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_s3_bucket Terraform resource. This
// resource is used to manage clumio S3 bucket properties.

package clumio_s3_bucket_properties

import (
	"context"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clumioS3BucketPropertiesResource{}
	_ resource.ResourceWithConfigure = &clumioS3BucketPropertiesResource{}
)

// clumioS3BucketPropertiesResource is the struct backing the clumio_s3_bucket_properties Terraform
// resource. It holds the Clumio API client and any other required state needed to manage
// S3 bucket properties within Clumio.

type clumioS3BucketPropertiesResource struct {
	name              string
	client            *common.ApiClient
	sdkS3BucketClient sdkclients.S3BucketClient
	pollTimeout       time.Duration
	pollInterval      time.Duration
}

// NewClumioS3BucketPropertiesResource creates a new instance of clumioS3BucketPropertiesResource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioS3BucketPropertiesResource() resource.Resource {
	return &clumioS3BucketPropertiesResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioS3BucketPropertiesResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_s3_bucket_properties"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioS3BucketPropertiesResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkS3BucketClient = sdkclients.NewS3BucketClient(r.client.ClumioConfig)
	r.pollInterval = 2 * time.Second
	r.pollTimeout = 60 * time.Second
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioS3BucketPropertiesResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioS3BucketPropertiesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createOrUpdateS3BucketProperties(ctx, &plan, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioS3BucketPropertiesResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioS3BucketPropertiesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readS3BucketProperties(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioS3BucketPropertiesResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioS3BucketPropertiesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createOrUpdateS3BucketProperties(ctx, &plan, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes it from the Terraform state.
func (r *clumioS3BucketPropertiesResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioS3BucketPropertiesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteS3BucketProperties(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
