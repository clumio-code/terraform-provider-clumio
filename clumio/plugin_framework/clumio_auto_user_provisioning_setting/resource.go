// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the auto user provisioning setting SDK APIs to perform CRUD
// operations and set the attributes from the response of the API in the resource model.

package clumio_auto_user_provisioning_setting

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// createAutoUserProvisioningSetting invokes the API to create the auto user provisioning setting
// and from the response populates the computed attributes of the auto user provisioning setting.
func (r *autoUserProvisioningSettingResource) createAutoUserProvisioningSetting(
	_ context.Context, plan *autoUserProvisioningSettingResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: plan.IsEnabled.ValueBoolPointer(),
	}

	// Call the Clumio API to enable or disable the auto user provisioning setting. NOTE that this
	// setting is a singleton state for the entire organization and as such, creation of this
	// resource results in "updating" the current state rather than creating a new instance of a
	// setting.
	_, apiErr := r.sdkAUPSettings.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s ", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	// Since the API doesn't return an id, we are setting a uuid as the resource id.
	plan.ID = types.StringValue(uuid.New().String())

	return diags
}

// createAutoUserProvisioningSetting invokes the API to read the auto user provisioning setting and
// from the response populates the attributes of the auto user provisioning setting.
func (r *autoUserProvisioningSettingResource) readAutoUserProvisioningSetting(
	_ context.Context, state *autoUserProvisioningSettingResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to read the auto user provisioning setting. The setting is an Org-wide
	// value and as such, the Org ID associated with the API credentials will be utilized.
	res, apiErr := r.sdkAUPSettings.ReadAutoUserProvisioningSetting()
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s ", r.name)
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

	// Convert the Clumio API response back to a schema and update the state. IsEnabled is the only
	// setting that needs to be refreshed.
	state.IsEnabled = types.BoolPointerValue(res.IsEnabled)

	return diags
}

// updateAutoUserProvisioningSetting invokes the API to update the auto user provisioning setting.
func (r *autoUserProvisioningSettingResource) updateAutoUserProvisioningSetting(
	_ context.Context, plan *autoUserProvisioningSettingResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	isEnabled := plan.IsEnabled.ValueBool()
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}

	// Call the Clumio API to enable or disable auto user provisioning setting.
	_, apiErr := r.sdkAUPSettings.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s ", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}

	return diags
}

// deleteAutoUserProvisioningSetting invokes the API to disable the auto user provisioning setting.
func (r *autoUserProvisioningSettingResource) deleteAutoUserProvisioningSetting(
	_ context.Context, state *autoUserProvisioningSettingResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	isEnabled := false
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}

	// Call the Clumio API to disable auto user provisioning setting.
	_, apiErr := r.sdkAUPSettings.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to delete %s ", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}

	return diags
}
