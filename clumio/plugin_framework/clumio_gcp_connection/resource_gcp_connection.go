// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

// This file holds the resource implementation for the clumio_gcp_connection Terraform resource.
// This resource is used to connect GCP projects to Clumio.

package clumio_gcp_connection

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clumioGCPConnectionResource{}
	_ resource.ResourceWithConfigure = &clumioGCPConnectionResource{}
)

type clumioGCPConnectionResource struct {
	name           string
	client         *common.ApiClient
	sdkConnections sdkclients.GcpConnectionClient
}

// NewClumioGCPConnectionResource creates a new instance of clumioGCPConnectionResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioGCPConnectionResource() resource.Resource {
	return &clumioGCPConnectionResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioGCPConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_gcp_connection"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioGCPConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkConnections = sdkclients.NewGcpConnectionClient(r.client.ClumioConfig)
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioGCPConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioGCPConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the GCP connection.
	remove, diags := r.readGcpConnection(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioGCPConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the current Terraform plan.
	var plan clumioGCPConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call Clumio API to create GCP connection
	diags = r.createGcpConnection(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and sets the updated Terraform state.
func (r *clumioGCPConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the current Terraform plan.
	var plan clumioGCPConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call Clumio API to update GCP connection
	diags = r.updateGcpConnection(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API
func (r *clumioGCPConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the schema from the current Terraform plan.
	var state clumioGCPConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call Clumio API to delete GCP connection
	diags = r.deleteGcpConnection(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
