// Copyright 2025. Clumio, Inc.

// This file holds the logic to invoke the general settings SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_general_settings

import (
	"context"
	"fmt"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// readGeneralSettings invokes the APIs to read the general settings and from the response
// populates the attributes of the general settings.
func (r *clumioGeneralSettings) readGeneralSettings(
	_ context.Context, state *generalSettingsResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkGeneralSettings := r.sdkGeneralSettings

	// Call the Clumio API to read the configuration.
	res, apiErr := sdkGeneralSettings.ReadGeneralSettings()
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
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

	// Convert the Clumio API response back to a schema and update the state.
	state.AutoLogoutDuration = types.Int64PointerValue(res.AutoLogoutDuration)
	state.PasswordExpirationDuration = types.Int64PointerValue(res.PasswordExpirationDuration)

	ipAllowlist := make([]types.String, 0)
	for _, ip := range res.IpAllowlist {
		ipAllowlist = append(ipAllowlist, types.StringPointerValue(ip))
	}
	state.IpAllowlist = ipAllowlist

	return diags
}

// updateGeneralSettings invokes the API to update the general settings and from the
// response populates the computed attributes of the general settings.
func (r *clumioGeneralSettings) updateGeneralSettings(
	_ context.Context, plan *generalSettingsResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkGeneralSettings := r.sdkGeneralSettings

	ipAllowlist := make([]*string, 0)
	if plan.IpAllowlist != nil {
		for _, ip := range plan.IpAllowlist {
			ipStr := ip.ValueString()
			ipAllowlist = append(ipAllowlist, &ipStr)
		}
	}

	// Convert the schema to a Clumio API request to set general settings.
	request := &models.UpdateGeneralSettingsV2Request{
		AutoLogoutDuration:         plan.AutoLogoutDuration.ValueInt64Pointer(),
		IpAllowlist:                ipAllowlist,
		PasswordExpirationDuration: plan.PasswordExpirationDuration.ValueInt64Pointer(),
	}

	// Call the Clumio API to update the general settings.
	res, apiErr := sdkGeneralSettings.UpdateGeneralSettings(request)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s", r.name)
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

	return diags
}
