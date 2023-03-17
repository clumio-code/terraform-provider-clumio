// Copyright 2023. Clumio, Inc.
//
// clumio_role definition and CRUD implementation.

package clumio_role

import (
	"context"
	"fmt"

	"github.com/clumio-code/clumio-go-sdk/controllers/roles"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioRoleDataSource{}
)

// NewClumioRoleDataSource is a helper function to simplify the provider implementation.
func NewClumioRoleDataSource() datasource.DataSource {
	return &clumioRoleDataSource{}
}

// clumioRoleDataSource is the data source implementation.
type clumioRoleDataSource struct {
	client *common.ApiClient
}

type permissionModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// clumioRoleDataSourceModel model
type clumioRoleDataSourceModel struct {
	Id          types.String       `tfsdk:"id"`
	Name        types.String       `tfsdk:"name"`
	Description types.String       `tfsdk:"description"`
	UserCount   types.Int64        `tfsdk:"user_count"`
	Permissions []*permissionModel `tfsdk:"permissions"`
}

// Metadata returns the resource type name.
func (r *clumioRoleDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the resource.
func (r *clumioRoleDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the role.",
				Computed:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "Unique name assigned to the role.",
				Required:    true,
			},
			schemaDescription: schema.StringAttribute{
				Description: "A description of the role.",
				Computed:    true,
			},
			schemaUserCount: schema.Int64Attribute{
				Description: "Number of users to whom the role has been assigned.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			schemaPermissions: schema.ListNestedBlock{
				Description: "Permissions contained in the role.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						schemaDescription: schema.StringAttribute{
							Description: "Description of the permission.",
							Computed:    true,
						},
						schemaId: schema.StringAttribute{
							Description: "The Clumio-assigned ID of the permission.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "Name of the permission.",
							Computed:    true,
						},
					},
				},
			},
		},
		Description: "Clumio Roles Data Source used to list the Clumio Roles.",
	}
}

// Configure adds the provider configured client to the data source.
func (r *clumioRoleDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioRoleDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Get current state
	var state clumioRoleDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rolesApi := roles.NewRolesV1(r.client.ClumioConfig)
	res, apiErr := rolesApi.ListRoles()
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error listing Clumio roles.",
			fmt.Sprintf("Error: %v", string(apiErr.Response)))
		return
	}
	var expectedRole *models.RoleWithETag
	for _, roleItem := range res.Embedded.Items {
		if *roleItem.Name == state.Name.ValueString() {
			expectedRole = roleItem
			break
		}
	}

	state.Id = types.StringValue(*expectedRole.Id)
	state.Name = types.StringValue(*expectedRole.Name)
	state.Description = types.StringValue(*expectedRole.Description)
	state.UserCount = types.Int64Value(*expectedRole.UserCount)

	permissions := make([]*permissionModel, len(expectedRole.Permissions))
	for ind, permission := range expectedRole.Permissions {
		permissionModel := &permissionModel{}
		permissionModel.Description = types.StringValue(*permission.Description)
		permissionModel.Id = types.StringValue(*permission.Id)
		permissionModel.Name = types.StringValue(*permission.Name)
		permissions[ind] = permissionModel
	}
	state.Permissions = permissions

	// Set refreshed state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	return
}
