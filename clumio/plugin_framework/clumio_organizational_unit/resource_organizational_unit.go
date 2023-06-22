// Copyright 2023. Clumio, Inc.
//
// clumio_organizational_unit definition and CRUD implementation.

package clumio_organizational_unit

import (
	"context"
	"fmt"
	"strings"

	orgUnits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioOrganizationalUnitResource{}
	_ resource.ResourceWithConfigure   = &clumioOrganizationalUnitResource{}
	_ resource.ResourceWithImportState = &clumioOrganizationalUnitResource{}
)

// NewClumioOrganizationalUnitResource is a helper function to simplify the provider implementation.
func NewClumioOrganizationalUnitResource() resource.Resource {
	return &clumioOrganizationalUnitResource{}
}

// clumioOrganizationalUnitResource is the resource implementation.
type clumioOrganizationalUnitResource struct {
	client *common.ApiClient
}

type userWithRole struct {
	UserId       types.String `tfsdk:"user_id"`
	AssignedRole types.String `tfsdk:"assigned_role"`
}

// clumioOrganizationalUnitResource model
type clumioOrganizationalUnitResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	ParentId                  types.String `tfsdk:"parent_id"`
	ChildrenCount             types.Int64  `tfsdk:"children_count"`
	ConfiguredDatasourceTypes types.List   `tfsdk:"configured_datasource_types"`
	DescendantIds             types.List   `tfsdk:"descendant_ids"`
	UserCount                 types.Int64  `tfsdk:"user_count"`
	Users                     types.List   `tfsdk:"users"`
	UsersWithRole             types.List   `tfsdk:"users_with_role"`
}

// Metadata returns the resource type name.
func (r *clumioOrganizationalUnitResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_unit"
}

// Schema defines the schema for the resource.
func (r *clumioOrganizationalUnitResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for creating and managing Organizational Unit in Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "OrganizationalUnit Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "Unique name assigned to the organizational unit.",
				Required:    true,
			},
			schemaDescription: schema.StringAttribute{
				Description: "A description of the organizational unit.",
				Optional:    true,
			},
			schemaParentId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the parent organizational unit" +
					" under which the new organizational unit is to be created.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			schemaChildrenCount: schema.Int64Attribute{
				Description: "Number of immediate children of the organizational unit.",
				Computed:    true,
			},
			schemaConfiguredDatasourceTypes: schema.ListAttribute{
				Description: "Datasource types configured in this organizational unit." +
					" Possible values include aws, microsoft365, vmware, or mssql.",
				ElementType: types.StringType,
				Computed:    true,
			},
			schemaDescendantIds: schema.ListAttribute{
				Description: "List of all recursive descendant organizational units" +
					" of this OU.",
				ElementType: types.StringType,
				Computed:    true,
			},
			schemaUserCount: schema.Int64Attribute{
				Description: "Number of users to whom this organizational unit or any" +
					" of its descendants have been assigned.",
				Computed: true,
			},
			schemaUsers: schema.ListAttribute{
				Description:        "List of user ids to assign this organizational unit.",
				ElementType:        types.StringType,
				Computed:           true,
				DeprecationMessage: "This attribute will be removed in the next major version of the provider.",
			},
			schemaUsersWithRole: schema.ListNestedAttribute{
				Description: "List of user ids, with role, to assign this organizational unit.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaUserId: schema.StringAttribute{
							Description: "The Clumio-assigned ID of the user.",
							Computed:    true,
						},
						schemaAssignedRole: schema.StringAttribute{
							Description: "The Clumio-assigned ID of the role assigned to the user.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *clumioOrganizationalUnitResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clumioOrganizationalUnitResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan clumioOrganizationalUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV2(r.client.ClumioConfig)
	name := plan.Name.ValueString()
	parentId := plan.ParentId.ValueString()
	request := &models.CreateOrganizationalUnitV2Request{
		Name:     &name,
		ParentId: &parentId,
	}
	description := plan.Description.ValueString()
	if !plan.Description.IsNull() {
		request.Description = &description
	}

	res, apiErr := orgUnitsAPI.CreateOrganizationalUnit(nil, request)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error creating Clumio organizational unit.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	var id types.String
	var nameString types.String
	var descriptionString types.String
	var parentIdString types.String
	var childrenCount types.Int64
	var userCount types.Int64
	var configuredDatasourceTypes []*string
	var userSlice []*string
	var userWithRoleSlice []userWithRole
	var descendantIdSlice []*string

	if res.StatusCode == http200 {
		id = types.StringValue(*res.Http200.Id)
		nameString = types.StringValue(*res.Http200.Name)
		if res.Http200.Description != nil {
			descriptionString = types.StringValue(*res.Http200.Description)
		}
		if res.Http200.ParentId != nil {
			parentIdString = types.StringValue(*res.Http200.ParentId)
		}
		childrenCount = types.Int64Value(*res.Http200.ChildrenCount)
		userCount = types.Int64Value(*res.Http200.UserCount)
		configuredDatasourceTypes = res.Http200.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http200.Users)
		descendantIdSlice = res.Http200.DescendantIds
	} else if res.StatusCode == http202 {
		id = types.StringValue(*res.Http202.Id)
		nameString = types.StringValue(*res.Http202.Name)
		if res.Http202.Description != nil {
			descriptionString = types.StringValue(*res.Http202.Description)
		}
		if res.Http202.ParentId != nil {
			parentIdString = types.StringValue(*res.Http202.ParentId)
		}
		childrenCount = types.Int64Value(*res.Http202.ChildrenCount)
		userCount = types.Int64Value(*res.Http202.UserCount)
		configuredDatasourceTypes = res.Http202.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http202.Users)
		descendantIdSlice = res.Http202.DescendantIds
	}
	plan.Id = id
	plan.Name = nameString
	plan.Description = descriptionString
	plan.ParentId = parentIdString
	plan.ChildrenCount = childrenCount
	plan.UserCount = userCount

	configuredDataTypes, conversionDiags := types.ListValueFrom(ctx, types.StringType, configuredDatasourceTypes)
	resp.Diagnostics.Append(conversionDiags...)
	plan.ConfiguredDatasourceTypes = configuredDataTypes

	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, descendantIdSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.DescendantIds = descendantIds

	// Set the state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioOrganizationalUnitResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state clumioOrganizationalUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV2(r.client.ClumioConfig)
	res, apiErr := orgUnitsAPI.ReadOrganizationalUnit(state.Id.ValueString(), nil)
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			state.Id = types.StringValue("")
		}
		resp.Diagnostics.AddError(
			"Error retrieving Clumio organizational unit.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	state.Id = types.StringValue(*res.Id)
	state.Name = types.StringValue(*res.Name)
	if res.Description != nil {
		state.Description = types.StringValue(*res.Description)
	}
	if res.ParentId != nil {
		state.ParentId = types.StringValue(*res.ParentId)
	}
	state.ChildrenCount = types.Int64Value(*res.ChildrenCount)
	state.UserCount = types.Int64Value(*res.UserCount)

	configuredDataTypes, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.ConfiguredDatasourceTypes)
	resp.Diagnostics.Append(conversionDiags...)
	state.ConfiguredDatasourceTypes = configuredDataTypes

	userSlice, userWithRoleSlice := getUsersFromHTTPRes(res.Users)
	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	resp.Diagnostics.Append(conversionDiags...)
	state.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	resp.Diagnostics.Append(conversionDiags...)
	state.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, res.DescendantIds)
	resp.Diagnostics.Append(conversionDiags...)
	state.DescendantIds = descendantIds

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioOrganizationalUnitResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan clumioOrganizationalUnitResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state clumioOrganizationalUnitResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV2(r.client.ClumioConfig)
	name := plan.Name.ValueString()
	request := &models.PatchOrganizationalUnitV2Request{
		Name: &name,
	}
	description := plan.Description.ValueString()
	if !plan.Description.IsNull() {
		request.Description = &description
	}
	res, apiErr := orgUnitsAPI.PatchOrganizationalUnit(plan.Id.ValueString(), nil, request)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating Clumio organizational unit id: %v.", plan.Id.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	var id types.String
	var nameString types.String
	var descriptionString types.String
	var parentIdString types.String
	var childrenCount types.Int64
	var userCount types.Int64
	var configuredDatasourceTypes []*string
	var userWithRoleSlice []userWithRole
	var userSlice []*string
	var descendantIdSlice []*string

	if res.StatusCode == http200 {
		id = types.StringValue(*res.Http200.Id)
		nameString = types.StringValue(*res.Http200.Name)
		if res.Http200.Description != nil {
			descriptionString = types.StringValue(*res.Http200.Description)
		}
		if res.Http200.ParentId != nil {
			parentIdString = types.StringValue(*res.Http200.ParentId)
		}
		childrenCount = types.Int64Value(*res.Http200.ChildrenCount)
		userCount = types.Int64Value(*res.Http200.UserCount)
		configuredDatasourceTypes = res.Http200.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http200.Users)
		descendantIdSlice = res.Http200.DescendantIds
	} else if res.StatusCode == http202 {
		id = types.StringValue(*res.Http202.Id)
		nameString = types.StringValue(*res.Http202.Name)
		if res.Http202.Description != nil {
			descriptionString = types.StringValue(*res.Http202.Description)
		}
		if res.Http202.ParentId != nil {
			parentIdString = types.StringValue(*res.Http202.ParentId)
		}
		childrenCount = types.Int64Value(*res.Http202.ChildrenCount)
		userCount = types.Int64Value(*res.Http202.UserCount)
		configuredDatasourceTypes = res.Http202.ConfiguredDatasourceTypes
		userSlice, userWithRoleSlice = getUsersFromHTTPRes(res.Http202.Users)
		descendantIdSlice = res.Http202.DescendantIds
	}
	plan.Id = id
	plan.Name = nameString
	plan.Description = descriptionString
	plan.ParentId = parentIdString
	plan.ChildrenCount = childrenCount
	plan.UserCount = userCount

	configuredDataTypes, conversionDiags := types.ListValueFrom(ctx, types.StringType, configuredDatasourceTypes)
	resp.Diagnostics.Append(conversionDiags...)
	plan.ConfiguredDatasourceTypes = configuredDataTypes

	users, conversionDiags := types.ListValueFrom(ctx, types.StringType, userSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.Users = users
	usersWithRole, conversionDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaUserId:       types.StringType,
			schemaAssignedRole: types.StringType,
		},
	}, userWithRoleSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.UsersWithRole = usersWithRole

	descendantIds, conversionDiags := types.ListValueFrom(ctx, types.StringType, descendantIdSlice)
	resp.Diagnostics.Append(conversionDiags...)
	plan.DescendantIds = descendantIds

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clumioOrganizationalUnitResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state clumioOrganizationalUnitResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV2(r.client.ClumioConfig)
	res, apiErr := orgUnitsAPI.DeleteOrganizationalUnit(state.Id.ValueString(), nil)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting Clumio organizational unit %v.", state.Id.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, pollTimeoutInSec, pollIntervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error while polling for Delete OU Task"),
			fmt.Sprintf(errorFmt, string(err.Error())))
	}
}

func (r *clumioOrganizationalUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getUsersFromHTTPRes(users []*models.UserWithRole) ([]*string, []userWithRole) {
	userSlice, userWithRoleSlice := make([]*string, len(users)), make([]userWithRole, len(users))
	for idx, user := range users {
		userSlice[idx] = user.UserId
		userWithRoleSlice[idx] = userWithRole{
			UserId:       types.StringValue(*user.UserId),
			AssignedRole: types.StringValue(*user.AssignedRole),
		}
	}

	return userSlice, userWithRoleSlice
}
