// Copyright 2026. Clumio, Inc.

// This file holds the logic to invoke the GCP protection group SDK API to perform read operations
// and set the attributes from the response of the API in the data source model.

package clumio_gcs_protection_group

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// readProtectionGroup invokes the API to list GCP protection groups matching the given name and
// from the response populates the id of the data source model.
func (r *clumioGCSProtectionGroupDataSource) readProtectionGroup(
	_ context.Context, model *clumioGCSProtectionGroupDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Prepare the query filter to look up the protection group by its exact name. The name is
	// JSON-encoded so that any quotes/backslashes are escaped and the filter stays valid JSON.
	name := model.Name.ValueString()
	nameJSON, err := json.Marshal(name)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to read %s", r.name),
			fmt.Sprintf("Failed to encode the name filter: %s", err.Error()))
		return diags
	}
	filter := fmt.Sprintf(`{"name": {"$eq":%s}}`, string(nameJSON))

	// Call the Clumio API to list the GCP protection groups.
	res, apiErr := r.sdkProtectionGroups.ListGcpProtectionGroups(nil, nil, &filter, nil, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		diags.AddError(common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return diags
	}
	if res.CurrentCount == nil || *res.CurrentCount == 0 ||
		res.Embedded == nil || len(res.Embedded.Items) == 0 {
		summary := "GCP protection group not found"
		detail := fmt.Sprintf("No GCP protection group found with name %q.", name)
		diags.AddError(summary, detail)
		return diags
	}

	model.Id = basetypes.NewStringPointerValue(res.Embedded.Items[0].Id)
	return diags
}
