// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio Organizational Unit SDK API to perform read
// operation and set the attributes from the response of the API in the data source model.

package clumio_organizational_unit

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readOrganizationalUnit invokes the API to read the organizationalUnitClient and from the response
// populates the attributes of the organizational unit.
func (r *clumioOrganizationalUnitDataSource) readOrganizationalUnit(
	ctx context.Context, model *clumioOrganizationalUnitDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics

	// Prepare the query nameFilter.
	name := model.Name.ValueString()
	nameFilter := fmt.Sprintf(`{"name": {"$contains":"%s"}}`, name)

	// Call the Clumio API to list the organizational units.
	limit := int64(10000)
	res, apiErr := r.organizationalUnitClient.ListOrganizationalUnits(&limit, nil, &nameFilter)
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

	// Convert the Clumio API response for the organizational units into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateOrganizationalUnitsInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "Organizational unit not found."
		detail := fmt.Sprintf(
			"No organizational unit found with the given name %s", name)
		diags.AddError(summary, detail)
	}
	return diags
}
