// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_post_process_gcp_connection

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &clumioPostProcessGCPConnectionResource{}
var _ resource.ResourceWithConfigure = &clumioPostProcessGCPConnectionResource{}

type clumioPostProcessGCPConnectionResource struct {
	name           string
	client         *common.ApiClient
	sdkConnections sdkclients.GcpConnectionClient
}

// NewClumioPostProcessGCPConnectionResource creates a new instance of clumioPostProcessGCPConnectionResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioPostProcessGCPConnectionResource() resource.Resource {
	return &clumioPostProcessGCPConnectionResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioPostProcessGCPConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_post_process_gcp_connection"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioPostProcessGCPConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkConnections = sdkclients.NewGcpConnectionClient(r.client.ClumioConfig)
}

// Read does not have an implementation as there is no API to read for post process gcp connection.
func (r *clumioPostProcessGCPConnectionResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Not supported
}

// Create creates a resource via Clumio API and sets initial Terraform state
func (r *clumioPostProcessGCPConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the schema from the current Terraform plan
	var plan clumioPostProcessGCPConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to process the GCP connection.
	diags = r.createUpdatePostProcessGcpConnection(ctx, &plan, createRequestType)
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

// Update updates a resource via Clumio API and sets the updated Terraform state.
func (r *clumioPostProcessGCPConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve schema from current terraform state
	var plan clumioPostProcessGCPConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to process the GCP connection.
	diags = r.createUpdatePostProcessGcpConnection(ctx, &plan, updateRequestType)
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

// Delete deletes the resource via API
func (r *clumioPostProcessGCPConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve schema from current terraform state
	var state clumioPostProcessGCPConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call Clumio API to delete GCP connection
	diags = r.deletePostProcessGcpConnection(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
