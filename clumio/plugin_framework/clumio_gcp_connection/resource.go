// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_gcp_connection

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

// readGcpConnection invokes the API to read the connection and from the response populates the
// attributes of the connection. If the connection has been removed externally, the function returns
// "true" to indicate to the caller that the resource no longer exists.
func (r *clumioGCPConnectionResource) readGcpConnection(ctx context.Context, state *clumioGCPConnectionResourceModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	res, apiErr := r.sdkConnections.ReadGcpConnection(state.ProjectID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state", r.name, state.ProjectID.ValueString())
			tflog.Warn(ctx, summary)
			return true, diags
		}

		summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.ProjectID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return false, diags
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the state
	state.Token = types.StringPointerValue(res.Token)
	state.ClumioControlPlaneId = types.StringPointerValue(res.ControlPlaneId)
	state.ClumioControlPlaneRole = types.StringPointerValue(res.ControlPlaneRole)

	regionsValue, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.Regions)
	diags.Append(conversionDiags...)
	state.Regions = regionsValue

	state.DeploymentType = types.StringPointerValue(res.DeploymentType)

	// Description and ProjectID are not computed values
	return false, diags
}

// createGcpConnection invokes the API to create the connection and from the response populates the
// computed attributes of the connection.
func (r *clumioGCPConnectionResource) createGcpConnection(ctx context.Context, plan *clumioGCPConnectionResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Convert regions from the plan to a slice of string pointers for the API request.
	var regions []*string
	if !plan.Regions.IsNull() && !plan.Regions.IsUnknown() {
		conversionDiags := plan.Regions.ElementsAs(ctx, &regions, false)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
	}

	// Convert schema to CreateGcpConnectionV1Request model
	req := &models.CreateGcpConnectionV1Request{
		DeploymentType: plan.DeploymentType.ValueStringPointer(),
		Description:    plan.Description.ValueStringPointer(),
		ProjectId:      plan.ProjectID.ValueStringPointer(),
		Regions:        regions,
	}

	// Call Clumio API to create a connection
	res, apiErr := r.sdkConnections.CreateGcpConnection(req)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s (project id: %v)", r.name, plan.ProjectID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// ID needs to be a value which is used by our backend to uniquely identify connection
	plan.ID = types.StringPointerValue(res.Token)
	plan.ClumioControlPlaneId = types.StringPointerValue(res.ControlPlaneId)
	plan.ClumioControlPlaneRole = types.StringPointerValue(res.ControlPlaneRole)
	plan.Token = types.StringPointerValue(res.Token)

	regionsValue, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.Regions)
	diags.Append(conversionDiags...)
	plan.Regions = regionsValue

	plan.DeploymentType = types.StringPointerValue(res.DeploymentType)

	return diags
}

// updateGcpConnection invokes the API to update the connection and from the response populates the
// computed attributes of the connection.
func (r *clumioGCPConnectionResource) updateGcpConnection(ctx context.Context, plan *clumioGCPConnectionResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Convert regions from the plan to a slice of string pointers for the API request.
	var regions []*string
	if !plan.Regions.IsNull() && !plan.Regions.IsUnknown() {
		conversionDiags := plan.Regions.ElementsAs(ctx, &regions, false)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
	}

	// Convert schema to UpdateGcpConnectionV1Request model
	req := &models.UpdateGcpConnectionV1Request{
		DeploymentType: plan.DeploymentType.ValueStringPointer(),
		Description:    plan.Description.ValueStringPointer(),
		Regions:        regions,
	}

	// Call Clumio API to update a connection
	_, apiErr := r.sdkConnections.UpdateGcpConnection(plan.ProjectID.ValueString(), req)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (project id: %v)", r.name, plan.ProjectID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	return diags
}

// deleteGcpConnection invokes the API to delete the connection
func (r *clumioGCPConnectionResource) deleteGcpConnection(ctx context.Context, state *clumioGCPConnectionResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Call Clumio API to delete a connection
	_, apiErr := r.sdkConnections.DeleteGcpConnection(state.ProjectID.ValueString())
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to delete %s (project id: %v)", r.name, state.ProjectID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	return diags
}
