// Copyright 2025. Clumio, Inc.

// This file holds the logic to invoke the report configuration SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_report_configuration

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createReportConfiguration invokes the API to create the report configuration and from the response
// populates the computed attributes of the report configuration.
func (r *clumioReportConfigurationResource) createReportConfiguration(
	_ context.Context, plan *reportConfigurationResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkReportConfigurations := r.sdkReportConfigurations

	// Convert the schema to a Clumio API request to set report configurations.
	request := &models.CreateComplianceReportConfigurationV1Request{
		Description:  plan.Description.ValueStringPointer(),
		Name:         plan.Name.ValueStringPointer(),
		Notification: mapSchemaNotificationToClumioNotification(plan.Notification),
		Parameter:    mapSchemaParameterToClumioParameter(plan.Parameter),
		Schedule:     mapSchemaScheduleToClumioSchedule(plan.Schedule),
	}

	// Call the Clumio API to set the report configurations
	res, apiErr := sdkReportConfigurations.CreateComplianceReportConfiguration(request)
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

	// Populate all computed fields of the plan including the ID given that the resource is getting
	// created.
	plan.ID = types.StringPointerValue(res.Id)
	plan.CreatedAt = types.StringPointerValue(res.Created)

	return diags
}

// readReportConfiguration invokes the APIs to read the report configuration and from the response
// populates the attributes of the report configuration. If the report configuration has been
// removed externally, the function returns "true" to indicate to the caller that the resource no
// longer exists.
func (r *clumioReportConfigurationResource) readReportConfiguration(
	ctx context.Context, state *reportConfigurationResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics
	sdkReportConfigurations := r.sdkReportConfigurations

	// Call the Clumio API to read the configuration.
	res, apiErr := sdkReportConfigurations.ReadComplianceReportConfiguration(state.ID.ValueString())
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
	description := types.StringPointerValue(res.Description)
	if !state.Description.IsNull() || description.ValueString() != "" {
		state.Description = description
	}
	state.CreatedAt = types.StringPointerValue(res.Created)
	state.Notification = mapClumioNotificationToSchemaNotification(res.Notification)
	state.Parameter = mapClumioParameterToSchemaParameter(res.Parameter)
	state.Schedule = mapClumioScheduleToSchemaSchedule(res.Schedule)

	return false, diags
}

// updateReportConfiguration invokes the API to update the report configuration and from the
// response populates the computed attributes of the report configuration. After update is done, it
// also verifies that the policy has been applied on the entity.
func (r *clumioReportConfigurationResource) updateReportConfiguration(
	_ context.Context, plan *reportConfigurationResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkReportConfigurations := r.sdkReportConfigurations

	// Convert the schema to a Clumio API request to set report configurations.
	request := &models.UpdateComplianceReportConfigurationV1Request{
		Description:  plan.Description.ValueStringPointer(),
		Name:         plan.Name.ValueStringPointer(),
		Notification: mapSchemaNotificationToClumioNotification(plan.Notification),
		Parameter:    mapSchemaParameterToClumioParameter(plan.Parameter),
		Schedule:     mapSchemaScheduleToClumioSchedule(plan.Schedule),
	}

	// Call the Clumio API to update the report configurations.
	res, apiErr := sdkReportConfigurations.UpdateComplianceReportConfiguration(
		plan.ID.ValueString(), request)
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

// deleteReportConfiguration invokes the API to delete the report configuration.
func (r *clumioReportConfigurationResource) deleteReportConfiguration(
	_ context.Context, state *reportConfigurationResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	sdkReportConfigurations := r.sdkReportConfigurations

	// Call the Clumio API to remove the report configuration.
	_, apiErr := sdkReportConfigurations.DeleteComplianceReportConfiguration(state.ID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}

	return diags
}
