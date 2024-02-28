// Copyright 2023. Clumio, Inc.
//
// This file holds the resource implementation for the clumio_wallet Terraform resource.
// This resource is used to setup BYOK with Clumio.

package clumio_wallet

import (
	"context"
	"fmt"
	"net/http"

	sdkWallets "github.com/clumio-code/clumio-go-sdk/controllers/wallets"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioWalletResource{}
	_ resource.ResourceWithConfigure   = &clumioWalletResource{}
	_ resource.ResourceWithImportState = &clumioWalletResource{}
)

// clumioWalletResource is the struct backing the clumio_wallet Terraform resource. It holds the
// Clumio API client and any other required state needed to create a wallet.
type clumioWalletResource struct {
	name       string
	client     *common.ApiClient
	sdkWallets sdkWallets.WalletsV1Client
}

// NewClumioWalletResource creates a new instance of clumioWalletResource. Its attributes are
// initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewClumioWalletResource() resource.Resource {
	return &clumioWalletResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioWalletResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_wallet"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioWalletResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkWallets = sdkWallets.NewWalletsV1(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioWalletResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioWalletResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to create the wallet.
	res, apiErr := r.sdkWallets.CreateWallet(&models.CreateWalletV1Request{
		AccountNativeId: plan.AccountNativeId.ValueStringPointer(),
	})
	if apiErr != nil {
		summary := "Unable to create Clumio wallet"
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.Id = types.StringPointerValue(res.Id)
	plan.State = types.StringPointerValue(res.State)
	plan.Token = types.StringPointerValue(res.Token)
	plan.ClumioAccountId = types.StringPointerValue(res.ClumioAwsAccountId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioWalletResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the Terraform state.
	var state clumioWalletResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the wallet.
	res, apiErr := r.sdkWallets.ReadWallet(state.Id.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Wallet with ID %s not found. Removing from state.",
				state.Id.ValueString())
			tflog.Warn(ctx, msgStr)
			resp.State.RemoveResource(ctx)
		} else {
			summary := "Unable to read the Clumio wallet."
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and update the state. In addition to
	// computed fields, all fields are populated from the API response in case any values have been
	// changed externally. ID is not updated however given that it is the field used to query the
	// resource from the backend.
	state.State = types.StringPointerValue(res.State)
	state.AccountNativeId = types.StringPointerValue(res.AccountNativeId)
	state.Token = types.StringPointerValue(res.Token)
	state.ClumioAccountId = types.StringPointerValue(res.ClumioAwsAccountId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioWalletResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	summary := "Update not expected"
	detail := "None of the schema attributes allow updates."
	resp.Diagnostics.AddError(summary, detail)
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioWalletResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the Terraform state.
	var state clumioWalletResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to delete the wallet.
	_, apiErr := r.sdkWallets.DeleteWallet(state.Id.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := "Unable to delete the Clumio wallet"
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioWalletResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
