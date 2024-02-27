// Copyright 2023. Clumio, Inc.

// This file holds the datasource implementation for the clumio_role Terraform datasource. This
// datasource is used to manage roles within Clumio

package clumio_role

import (
	"context"
	"fmt"

	"github.com/clumio-code/clumio-go-sdk/controllers/roles"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioRoleDataSource{}
)

// clumioRoleDataSource is the struct backing the clumio_role Terraform datasource. It holds the 
// Clumio API client and any other required state needed to manage roles within Clumio.
type clumioRoleDataSource struct {
	name            string
	client          *common.ApiClient
	roles           roles.RolesV1Client
}

// NewClumioRoleDataSource creates a new instance of clumioRoleDataSource. Its attributes are 
// initialized later by Terraform via Metadata and Configure once the Provider is initialized.
func NewClumioRoleDataSource() datasource.DataSource {
	return &clumioRoleDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioRoleDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_role"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioRoleDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.roles = roles.NewRolesV1(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioRoleDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioRoleDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the list of roles.
	res, apiErr := r.roles.ListRoles()
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error listing Clumio roles.",
			fmt.Sprintf("Error: %v", common.ParseMessageFromApiError(apiErr)))
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Find the expected role from the list of fetched roles via its name.
	var expectedRole *models.RoleWithETag
	for _, roleItem := range res.Embedded.Items {
		if *roleItem.Name == state.Name.ValueString() {
			expectedRole = roleItem
			break
		}
	}
	// Throw error if role with provided name was not found.
	if expectedRole == nil {
			summary := "Role not found"
			detail := fmt.Sprintf("Couldn't find a role with the provided name %s", state.Name.ValueString())
			resp.Diagnostics.AddError(summary, detail)
			return
	}

	// Convert the Clumio API response for the expected role into a schema and update the state.
	state.Id = types.StringPointerValue(expectedRole.Id)
	state.Name = types.StringPointerValue(expectedRole.Name)
	state.Description = types.StringPointerValue(expectedRole.Description)
	state.UserCount = types.Int64PointerValue(expectedRole.UserCount)

	// Go through all permissions inside the expected role API response and map it to permissionModel.
	// Then update the state with it.
	permissions := make([]*permissionModel, len(expectedRole.Permissions))
	for ind, permission := range expectedRole.Permissions {
		permissionModel := &permissionModel{}
		permissionModel.Description = types.StringPointerValue(permission.Description)
		permissionModel.Id = types.StringPointerValue(permission.Id)
		permissionModel.Name = types.StringPointerValue(permission.Name)
		permissions[ind] = permissionModel
	}
	state.Permissions = permissions

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	return
}
