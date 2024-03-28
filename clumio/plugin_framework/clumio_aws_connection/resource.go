// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio AWS Connection SDK APIs to perform CRUD operations
// and set the attributes from the response of the API in the resource model.

package clumio_aws_connection

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

// createAWSConnection invokes the API to create the connection and from the response populates the
// computed attributes of the connection.
func (r *clumioAWSConnectionResource) createAWSConnection(
	ctx context.Context, plan *clumioAWSConnectionResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	if !plan.OrganizationalUnitID.IsNull() {
		_, err := getOrgUnitForConnection(ctx, r, plan)
		if err != nil {
			summary := fmt.Sprintf("invalid %s", schemaOrganizationalUnitId)
			detail := err.Error()
			diags.AddError(summary, detail)
			return diags
		}
	}
	// Convert the schema to a Clumio API request to create an AWS connection.
	createReq := &models.CreateAwsConnectionV1Request{
		AccountNativeId:      plan.AccountNativeID.ValueStringPointer(),
		AwsRegion:            plan.AWSRegion.ValueStringPointer(),
		Description:          plan.Description.ValueStringPointer(),
		OrganizationalUnitId: plan.OrganizationalUnitID.ValueStringPointer(),
	}

	// Call the Clumio API to create the AWS connection.
	res, apiErr := r.sdkConnections.CreateAwsConnection(createReq)
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
	plan.ID = types.StringPointerValue(res.Id)
	plan.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	plan.ConnectionStatus = types.StringPointerValue(res.ConnectionStatus)
	plan.Token = types.StringPointerValue(res.Token)
	plan.Namespace = types.StringPointerValue(res.Namespace)
	plan.ClumioAWSAccountID = types.StringPointerValue(res.ClumioAwsAccountId)
	plan.ClumioAWSRegion = types.StringPointerValue(res.ClumioAwsRegion)
	setExternalId(plan, res.ExternalId, res.Token)
	setDataPlaneAccountId(plan, res.DataPlaneAccountId)

	return diags
}

// readAWSConnection invokes the API to read the connection and from the response populates the
// attributes of the connection. If the connection has been removed externally, the function returns
// "true" to indicate to the caller that the resource no longer exists.
func (r *clumioAWSConnectionResource) readAWSConnection(
	ctx context.Context, state *clumioAWSConnectionResourceModel) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	// Call the Clumio API to read the AWS connection.
	returnExternalId := "true"
	res, apiErr := r.sdkConnections.ReadAwsConnection(state.ID.ValueString(), &returnExternalId)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state", r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			return true, diags
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return false, diags
		}
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
	state.AccountNativeID = types.StringPointerValue(res.AccountNativeId)
	state.AWSRegion = types.StringPointerValue(res.AwsRegion)

	// Since the Description field is optional, it should only be populated if it initially
	// contained a non-null value or if there is a specific value that needs to be assigned.
	description := types.StringPointerValue(res.Description)
	if !state.Description.IsNull() || description.ValueString() != "" {
		state.Description = description
	}
	state.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	state.ConnectionStatus = types.StringPointerValue(res.ConnectionStatus)
	state.Token = types.StringPointerValue(res.Token)
	state.Namespace = types.StringPointerValue(res.Namespace)
	state.ClumioAWSAccountID = types.StringPointerValue(res.ClumioAwsAccountId)
	state.ClumioAWSRegion = types.StringPointerValue(res.ClumioAwsRegion)
	setExternalId(state, res.ExternalId, res.Token)
	setDataPlaneAccountId(state, res.DataPlaneAccountId)

	return false, diags
}

// updateAWSConnection invokes the API to update the connection and from the response populates the
// computed attributes of the connection.
func (r *clumioAWSConnectionResource) updateAWSConnection(
	ctx context.Context, plan *clumioAWSConnectionResourceModel,
	state *clumioAWSConnectionResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Update the Organizational Unit (OU) associated with the AWS connection if it has been
	// explicitly set and is different from the current OU.
	if !plan.OrganizationalUnitID.IsUnknown() &&
		plan.OrganizationalUnitID != state.OrganizationalUnitID {

		err := updateOrgUnitForConnection(ctx, r, plan, state)
		if err != nil {
			summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, state.ID.ValueString())
			detail := err.Error()
			diags.AddError(summary, detail)
			return diags
		}
		// "description" is the only field within the schema that can cause an update to the
		// resource. As such, if this has not changed, no need to go further and return early.
		if plan.Description == state.Description {
			return diags
		}
	}

	// Call the Clumio API to update the AWS connection. The "Description" parameter, while optional
	// in the REST API, is deliberately provided to ensure the update process is executed, even in
	// the absence of a specified description. This approach is necessary as the API's response
	// varies when a description is omitted.
	description := plan.Description.ValueString()
	updateReq := models.UpdateAwsConnectionV1Request{
		Description: &description,
	}
	res, apiErr := r.sdkConnections.UpdateAwsConnection(plan.ID.ValueString(), updateReq)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, state.ID.ValueString())
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

	// Convert the Clumio API response back to a schema and populate all computed fields of the
	// plan. ID however is not updated given that it is the field used to denote which resource to
	// update in the backend. Additionally the external ID is currently not returned during an
	// update call and thus is not updated below. This is okay however as the external ID is not
	// expected to change once a connection is created.
	plan.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	plan.ConnectionStatus = types.StringPointerValue(res.ConnectionStatus)
	plan.Token = types.StringPointerValue(res.Token)
	plan.Namespace = types.StringPointerValue(res.Namespace)
	plan.ClumioAWSAccountID = types.StringPointerValue(res.ClumioAwsAccountId)
	plan.ClumioAWSRegion = types.StringPointerValue(res.ClumioAwsRegion)
	setDataPlaneAccountId(plan, res.DataPlaneAccountId)

	return diags
}

// deleteAWSConnection invokes the API to delete the connection.
func (r *clumioAWSConnectionResource) deleteAWSConnection(
	_ context.Context, state *clumioAWSConnectionResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Call the Clumio API to delete the AWS connection.
	_, apiErr := r.sdkConnections.DeleteAwsConnection(state.ID.ValueString())
	if apiErr != nil && apiErr.ResponseCode != http.StatusNotFound {
		summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}
	return diags
}
