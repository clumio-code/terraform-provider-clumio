// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio Protection Group SDK API to perform read operation
// and set the attributes from the response of the API in the data source model.

package clumio_protection_group

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readProtectionGroup invokes the API to read the protectionGroupClient and from the response
// populates the attributes of the protection group.
func (r *clumioProtectionGroupDataSource) readProtectionGroup(
	ctx context.Context, model *clumioProtectionGroupDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Prepare the query filter.
	name := model.Name.ValueString()
	filter := fmt.Sprintf(`{"name": {"$eq":"%s"}}`, name)

	// Call the Clumio API to list the protection groups.
	res, apiErr := r.protectionGroupClient.ListProtectionGroups(
		nil, nil, &filter, &common.DefaultLookBackDays)
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
	// Convert the Clumio API response for the policies into the datasource schema model.
	if res.CurrentCount == nil || *res.CurrentCount == 0 {
		summary := "Protection group not found"
		detail := fmt.Sprintf(
			"Expected protection group with the specified name %s is not found", name)
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response for the policies into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		model.Id = basetypes.NewStringValue(*res.Embedded.Items[0].Id)
	}
	return diags
}
