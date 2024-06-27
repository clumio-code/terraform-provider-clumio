// Copyright 2024. Clumio, Inc.

// This file holds the resource implementation for the clumio_protection_group_bucket Terraform
// resource. This resource is used to manage S3 bucket assignment to a Protection Group.

package clumio_protection_group_bucket

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioProtectionGroupBucketResource{}
	_ resource.ResourceWithConfigure   = &clumioProtectionGroupBucketResource{}
	_ resource.ResourceWithImportState = &clumioProtectionGroupBucketResource{}
)

// clumioProtectionGroupBucketResource is the struct backing the clumio_protection_group_bucket
// Terraform resource. It holds the Clumio API client and any other required state needed to manage
// protection group bucket assignment within Clumio.
type clumioProtectionGroupBucketResource struct {
	name                string
	client              *common.ApiClient
	sdkProtectionGroups sdkclients.ProtectionGroupClient
	sdkS3Assets         sdkclients.ProtectionGroupS3AssetsClient
}

// NewClumioProtectionGroupBucketResource creates a new instance of
// clumioProtectionGroupBucketResource. Its attributes are initialized later by Terraform via
// Metadata and Configure once the Provider is initialized.
func NewClumioProtectionGroupBucketResource() resource.Resource {
	return &clumioProtectionGroupBucketResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioProtectionGroupBucketResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_protection_group_bucket"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioProtectionGroupBucketResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkProtectionGroups = sdkclients.NewProtectionGroupClient(r.client.ClumioConfig)
	r.sdkS3Assets = sdkclients.NewProtectionGroupS3AssetsClient(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioProtectionGroupBucketResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioProtectionGroupBucketResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createProtectionGroupBucket(ctx, &plan)
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
func (r *clumioProtectionGroupBucketResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioProtectionGroupBucketResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readProtectionGroupBucket(ctx, &state)
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
func (r *clumioProtectionGroupBucketResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Not implemented as none of the schema attribute supports update.
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioProtectionGroupBucketResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioProtectionGroupBucketResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteProtectionGroupBucket(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioProtectionGroupBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
