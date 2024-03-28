// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_aws_connection Terraform resource.
// This resource is used to connect AWS accounts to Clumio.

package clumio_aws_connection

import (
	"context"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the following Resource interfaces.
var (
	_ resource.Resource                = &clumioAWSConnectionResource{}
	_ resource.ResourceWithConfigure   = &clumioAWSConnectionResource{}
	_ resource.ResourceWithImportState = &clumioAWSConnectionResource{}
)

// clumioAWSConnectionResource is the struct backing the clumio_aws_connection Terraform resource.
// It holds the Clumio API client and any other required state needed to connect AWS accounts to
// Clumio.
type clumioAWSConnectionResource struct {
	name            string
	client          *common.ApiClient
	sdkConnections  sdkclients.AWSConnectionClient
	sdkEnvironments sdkclients.AWSEnvironmentClient
	sdkOrgUnits     sdkclients.OrganizationalUnitClient
	sdkTasks        sdkclients.TaskClient
	pollTimeout     time.Duration
	pollInterval    time.Duration
}

// NewClumioAWSConnectionResource creates a new instance of clumioAWSConnectionResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioAWSConnectionResource() resource.Resource {
	return &clumioAWSConnectionResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioAWSConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_aws_connection"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioAWSConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkConnections = sdkclients.NewAWSConnectionClient(r.client.ClumioConfig)
	r.sdkEnvironments = sdkclients.NewAWSEnvironmentClient(r.client.ClumioConfig)
	r.sdkOrgUnits = sdkclients.NewOrganizationalUnitClient(r.client.ClumioConfig)
	r.sdkTasks = sdkclients.NewTaskClient(r.client.ClumioConfig)
	r.pollTimeout = 3600 * time.Second
	r.pollInterval = 5 * time.Second
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioAWSConnectionResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to create the AWS connection.
	diags = r.createAWSConnection(ctx, &plan)
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
func (r *clumioAWSConnectionResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the AWS connection.
	remove, diags := r.readAWSConnection(ctx, &state)
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

// Update updates the resource via the Clumio API and updates the Terraform state. NOTE that the
// update for OU is a separate API call than the update for the AWS connection. Due to this it is
// possible for one portion of an update to go through but not the other. However, the update is
// idemptent so if a portion of the update fails, the next apply will attempt to update the failed
// portion again.
func (r *clumioAWSConnectionResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSConnectionResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.updateAWSConnection(ctx, &plan, &state)
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

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioAWSConnectionResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to delete the AWS connection.
	diags = r.deleteAWSConnection(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioAWSConnectionResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
