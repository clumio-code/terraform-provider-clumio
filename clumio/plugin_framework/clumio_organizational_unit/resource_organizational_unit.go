// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_organizational_unit Terraform resource.
// This resource is used to manage organizational units within Clumio.

package clumio_organizational_unit

import (
	"context"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioOrganizationalUnitResource{}
	_ resource.ResourceWithConfigure   = &clumioOrganizationalUnitResource{}
	_ resource.ResourceWithImportState = &clumioOrganizationalUnitResource{}
)

// clumioOrganizationalUnitResource is the struct backing the clumio_organizational_unit Terraform resource.
// It holds the Clumio API client and any other required state needed to
// manage organizational units within Clumio.
type clumioOrganizationalUnitResource struct {
	name         string
	client       *common.ApiClient
	sdkOrgUnits  sdkclients.OrganizationalUnitClient
	sdkTasks     sdkclients.TaskClient
	pollTimeout  time.Duration
	pollInterval time.Duration
}

// NewClumioOrganizationalUnitResource creates a new instance of clumioOrganizationalUnitResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioOrganizationalUnitResource() resource.Resource {
	return &clumioOrganizationalUnitResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioOrganizationalUnitResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_organizational_unit"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioOrganizationalUnitResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkOrgUnits = sdkclients.NewOrganizationalUnitClient(r.client.ClumioConfig)
	r.sdkTasks = sdkclients.NewTaskClient(r.client.ClumioConfig)
	r.pollTimeout = 3600 * time.Second
	r.pollInterval = 5 * time.Second
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioOrganizationalUnitResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioOrganizationalUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.createOrganizationalUnit(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioOrganizationalUnitResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioOrganizationalUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	remove, diags := r.readOrganizationalUnit(ctx, &state)
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
func (r *clumioOrganizationalUnitResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioOrganizationalUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateOrganizationalUnit(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioOrganizationalUnitResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioOrganizationalUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.deleteOrganizationalUnit(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioOrganizationalUnitResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
