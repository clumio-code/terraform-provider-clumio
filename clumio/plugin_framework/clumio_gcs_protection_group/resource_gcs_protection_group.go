// Copyright 2026. Clumio, Inc.

// This file holds the resource implementation for the clumio_gcs_protection_group Terraform
// resource. This resource is used to manage GCS protection groups within Clumio.

package clumio_gcs_protection_group

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clumioGCSProtectionGroupResource{}
	_ resource.ResourceWithConfigure = &clumioGCSProtectionGroupResource{}
)

// clumioGCSProtectionGroupResource is the struct backing the clumio_gcs_protection_group Terraform
// resource. It holds the Clumio API client and any other required state needed to manage GCS
// protection groups within Clumio.
type clumioGCSProtectionGroupResource struct {
	name                string
	client              *common.ApiClient
	sdkProtectionGroups sdkclients.GcpProtectionGroupClient
}

// NewClumioGCSProtectionGroupResource creates a new instance of clumioGCSProtectionGroupResource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioGCSProtectionGroupResource() resource.Resource {
	return &clumioGCSProtectionGroupResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioGCSProtectionGroupResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_gcs_protection_group"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioGCSProtectionGroupResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkProtectionGroups = sdkclients.NewGcpProtectionGroupClient(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioGCSProtectionGroupResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioGCSProtectionGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createProtectionGroup(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioGCSProtectionGroupResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioGCSProtectionGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readProtectionGroup(ctx, &state)
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
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioGCSProtectionGroupResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioGCSProtectionGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateProtectionGroup(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioGCSProtectionGroupResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioGCSProtectionGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteProtectionGroup(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
