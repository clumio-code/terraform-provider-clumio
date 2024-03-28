// Copyright 2023. Clumio, Inc.
//
// This file holds the resource implementation for the clumio_post_process_kms Terraform resource.
// This resource is used to send the necessary information required by Clumio to post-process BYOK
// after the necessary resources have been created. This resource should only be invoked as part of
// the byok-template module.

package clumio_post_process_kms

import (
	"context"
	"fmt"

	sdkPostProcessKms "github.com/clumio-code/clumio-go-sdk/controllers/post_process_kms"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clumioPostProcessKmsResource{}
	_ resource.ResourceWithConfigure = &clumioPostProcessKmsResource{}
)

// clumioPostProcessKmsResource is the struct backing the clumio_post_process_kms Terraform resource.
// It holds the Clumio API client and any other required state needed do post process kms after the
// necessary resources have been created. This resource should only be invoked as part of the
// byok-template module.
type clumioPostProcessKmsResource struct {
	client            *common.ApiClient
	sdkPostProcessKMS sdkPostProcessKms.PostProcessKmsV1Client
}

// NewClumioPostProcessKmsResource creates a new instance of clumioPostProcessKmsResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioPostProcessKmsResource() resource.Resource {
	return &clumioPostProcessKmsResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioPostProcessKmsResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_process_kms"
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioPostProcessKmsResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkPostProcessKMS = sdkPostProcessKms.NewPostProcessKmsV1(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioPostProcessKmsResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioPostProcessKmsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.clumioPostProcessKmsCommon(ctx, plan, "Create")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	accountId := plan.AccountId.ValueString()
	awsRegion := plan.Region.ValueString()
	token := plan.Token.ValueString()
	plan.Id = types.StringValue(fmt.Sprintf("%v/%v/%v", accountId, awsRegion, token))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read does not have an implementation as there is no API to read for post process kms.
func (r *clumioPostProcessKmsResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No implementation needed.
}

// Update updates the resource via the Clumio API and removes the Terraform state.
func (r *clumioPostProcessKmsResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioPostProcessKmsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.clumioPostProcessKmsCommon(ctx, plan, "Update")
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
func (r *clumioPostProcessKmsResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the Terraform state.
	var state clumioPostProcessKmsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.clumioPostProcessKmsCommon(ctx, state, "Delete")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
