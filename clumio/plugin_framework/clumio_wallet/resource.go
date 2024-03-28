// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio Wallet SDK APIs to perform CRUD operations
// and set the attributes from the response of the API in the resource model.

package clumio_wallet

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createWallet invokes the API to create the wallet and from the response populates the computed
// attributes of the wallet.
func (r *clumioWalletResource) createWallet(
	_ context.Context, plan *clumioWalletResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to create the wallet.
	res, apiErr := r.sdkWallets.CreateWallet(&models.CreateWalletV1Request{
		AccountNativeId: plan.AccountNativeId.ValueStringPointer(),
	})
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.Id = types.StringPointerValue(res.Id)
	plan.State = types.StringPointerValue(res.State)
	plan.Token = types.StringPointerValue(res.Token)
	plan.ClumioAccountId = types.StringPointerValue(res.ClumioAwsAccountId)
	return diags
}

// readWallet invokes the API to read the wallet and from the response populates the attributes of
// the wallet. If the wallet has been removed externally, the function returns "true" to indicate to
// the caller that the resource no longer exists.
func (r *clumioWalletResource) readWallet(
	ctx context.Context, state *clumioWalletResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	// Call the Clumio API to read the wallet.
	res, apiErr := r.sdkWallets.ReadWallet(state.Id.ValueString())
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Wallet with ID %s not found. Removing from state.",
				state.Id.ValueString())
			tflog.Warn(ctx, msgStr)
			remove = true
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.Id.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
		}
		return remove, diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return false, diags
	}

	// Convert the Clumio API response back to a schema and update the state. In addition to
	// computed fields, all fields are populated from the API response in case any values have been
	// changed externally. ID is not updated however given that it is the field used to query the
	// resource from the backend.
	state.State = types.StringPointerValue(res.State)
	state.AccountNativeId = types.StringPointerValue(res.AccountNativeId)
	state.Token = types.StringPointerValue(res.Token)
	state.ClumioAccountId = types.StringPointerValue(res.ClumioAwsAccountId)

	return false, diags
}

// deleteWallet invokes the API to delete the wallet.
func (r *clumioWalletResource) deleteWallet(
	_ context.Context, state *clumioWalletResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to delete the wallet.
	_, apiErr := r.sdkWallets.DeleteWallet(state.Id.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.Id.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}

	return diags
}
