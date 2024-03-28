// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the auto user provisioning rule SDK APIs to perform CRUD
// operations and set the attributes from the response of the API in the resource model.

package clumio_auto_user_provisioning_rule

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

// createAutoUserProvisioningRule invokes the API to create the auto user provisioning rule and from
// the response populates the computed attributes of the auto user provisioning rule.
func (r *autoUserProvisioningRuleResource) createAutoUserProvisioningRule(
	ctx context.Context, plan *autoUserProvisioningRuleResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Convert the schema to a Clumio API request to create an auto user provisioning rule.
	ouIds := make([]*string, 0)
	conversionDiags := plan.OrganizationalUnitIDs.ElementsAs(ctx, &ouIds, false)
	diags.Append(conversionDiags...)
	provision := &models.RuleProvision{
		RoleId:                plan.RoleID.ValueStringPointer(),
		OrganizationalUnitIds: ouIds,
	}
	auprRequest := &models.CreateAutoUserProvisioningRuleV1Request{
		Name:      plan.Name.ValueStringPointer(),
		Condition: plan.Condition.ValueStringPointer(),
		Provision: provision,
	}

	// Call the Clumio API to create the auto user provisioning rule.
	res, apiErr := r.sdkAUPRules.CreateAutoUserProvisioningRule(auprRequest)
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

	// Set the ID in the schema.
	plan.ID = types.StringPointerValue(res.RuleId)

	return diags
}

// readAutoUserProvisioningRule invokes the API to read the auto user provisioning rule and from
// the response populates the attributes of the auto user provisioning rule.
func (r *autoUserProvisioningRuleResource) readAutoUserProvisioningRule(
	ctx context.Context, state *autoUserProvisioningRuleResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	// Call the Clumio API to read the auto user provisioning rule.
	res, apiErr := r.sdkAUPRules.ReadAutoUserProvisioningRule(state.ID.ValueString())
	if apiErr != nil {
		remove := false
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf(
				"%s (ID: %v) not found. Removing from state", r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			remove = true
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.ID.ValueString())
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
	state.Name = types.StringPointerValue(res.Name)
	state.Condition = types.StringPointerValue(res.Condition)
	state.RoleID = types.StringPointerValue(res.Provision.RoleId)
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, res.Provision.OrganizationalUnitIds)
	diags.Append(conversionDiags...)
	state.OrganizationalUnitIDs = ouIds

	return false, diags
}

// updateAutoUserProvisioningRule invokes the API to update the auto user provisioning rule and from
// the response populates the computed attributes of the auto user provisioning rule.
func (r *autoUserProvisioningRuleResource) updateAutoUserProvisioningRule(
	ctx context.Context, plan *autoUserProvisioningRuleResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	ouIds := make([]*string, 0)
	conversionDiags := plan.OrganizationalUnitIDs.ElementsAs(ctx, &ouIds, false)
	diags.Append(conversionDiags...)
	provision := &models.RuleProvision{
		RoleId:                plan.RoleID.ValueStringPointer(),
		OrganizationalUnitIds: ouIds,
	}
	auprRequest := &models.UpdateAutoUserProvisioningRuleV1Request{
		Name:      plan.Name.ValueStringPointer(),
		Condition: plan.Condition.ValueStringPointer(),
		Provision: provision,
	}

	// Call the Clumio API to update the auto user provisioning rule.
	res, apiErr := r.sdkAUPRules.UpdateAutoUserProvisioningRule(plan.ID.ValueString(), auprRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
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

// deleteAutoUserProvisioningRule invokes the API to delete the auto user provisioning rule.
func (r *autoUserProvisioningRuleResource) deleteAutoUserProvisioningRule(
	_ context.Context, state *autoUserProvisioningRuleResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to delete the auto user provisioning rule.
	_, apiErr := r.sdkAUPRules.DeleteAutoUserProvisioningRule(state.ID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}

	return diags
}
