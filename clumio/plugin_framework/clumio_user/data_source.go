// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio User SDK API to perform read operation
// and set the attributes from the response of the API in the data source model.

package clumio_user

import (
	"context"
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readUser invokes the API to read the userClient and from the response populates the attributes of
// the user.
func (r *clumioUserDataSource) readUser(
	ctx context.Context, model *clumioUserDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	filters := make([]string, 0)

	// Prepare the query nameFilter.
	name := model.Name.ValueString()
	if name != "" {
		nameFilter := fmt.Sprintf(`"name": {"$contains":"%s"}`, name)
		filters = append(filters, nameFilter)
	}

	roleId := model.RoleId.ValueString()
	if roleId != "" {
		roleFilter := fmt.Sprintf(`"role_id": {"$eq":"%s"}`, roleId)
		filters = append(filters, roleFilter)
	}
	filter := fmt.Sprintf("{%s}", strings.Join(filters, ","))

	// Call the Clumio API to list the users.
	limit := int64(10000)
	res, apiErr := r.userClient.ListUsers(&limit, nil, &filter)
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

	// Convert the Clumio API response for the users into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		populateDiag := populateUsersInDataSourceModel(ctx, model, res.Embedded.Items)
		diags.Append(populateDiag...)
	} else {
		summary := "User not found."
		detail := fmt.Sprintf(
			"No user found with the given name %s and role_id %s", name, roleId)
		diags.AddError(summary, detail)
	}
	return diags
}
