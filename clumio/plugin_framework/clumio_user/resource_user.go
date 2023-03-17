// Copyright 2023. Clumio, Inc.
//
// clumio_user definition and CRUD implementation.

package clumio_user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/clumio-code/clumio-go-sdk/controllers/users"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioUserResource{}
	_ resource.ResourceWithConfigure   = &clumioUserResource{}
	_ resource.ResourceWithImportState = &clumioUserResource{}
)

// NewClumioUserResource is a helper function to simplify the provider implementation.
func NewClumioUserResource() resource.Resource {
	return &clumioUserResource{}
}

// clumioUserResource is the resource implementation.
type clumioUserResource struct {
	client *common.ApiClient
}

// clumioUserResource model
type clumioUserResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Email                   types.String `tfsdk:"email"`
	FullName                types.String `tfsdk:"full_name"`
	AssignedRole            types.String `tfsdk:"assigned_role"`
	OrganizationalUnitIds   types.Set    `tfsdk:"organizational_unit_ids"`
	Inviter                 types.String `tfsdk:"inviter"`
	IsConfirmed             types.Bool   `tfsdk:"is_confirmed"`
	IsEnabled               types.Bool   `tfsdk:"is_enabled"`
	LastActivityTimestamp   types.String `tfsdk:"last_activity_timestamp"`
	OrganizationalUnitCount types.Int64  `tfsdk:"organizational_unit_count"`
}

// Metadata returns the resource type name.
func (r *clumioUserResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *clumioUserResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Clumio User Resource to create and manage users in Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "User Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaEmail: schema.StringAttribute{
				Description: "The email address of the user to be added to Clumio.",
				Required:    true,
			},
			schemaFullName: schema.StringAttribute{
				Description: "The full name of the user to be added to Clumio." +
					" For example, enter the user's first name and last name. The name" +
					" appears in the User Management screen and in the body of the" +
					" email invitation.",
				Required: true,
			},
			schemaAssignedRole: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the role to assign to the user.",
				Optional:    true,
				Computed:    true,
			},
			schemaOrganizationalUnitIds: schema.SetAttribute{
				Description: "The Clumio-assigned IDs of the organizational units" +
					" to be assigned to the user. The Global Organizational Unit ID is " +
					"\"00000000-0000-0000-0000-000000000000\"",
				Required:    true,
				ElementType: types.StringType,
			},
			schemaInviter: schema.StringAttribute{
				Description: "The ID number of the user who sent the email invitation.",
				Computed:    true,
			},
			schemaIsConfirmed: schema.BoolAttribute{
				Description: "Determines whether the user has activated their Clumio" +
					" account. If true, the user has activated the account.",
				Computed: true,
			},
			schemaIsEnabled: schema.BoolAttribute{
				Description: "Determines whether the user is enabled (in Activated or" +
					" Invited status) in Clumio. If true, the user is in Activated or" +
					" Invited status in Clumio. Users in Activated status can log in to" +
					" Clumio. Users in Invited status have been invited to log in to" +
					" Clumio via an email invitation and the invitation is pending" +
					" acceptance from the user. If false, the user has been manually" +
					" suspended and cannot log in to Clumio until another Clumio user" +
					" reactivates the account.",
				Computed: true,
			},
			schemaLastActivityTimestamp: schema.StringAttribute{
				Description: "The timestamp of when when the user was last active in" +
					" the Clumio system. Represented in RFC-3339 format.",
				Computed: true,
			},
			schemaOrganizationalUnitCount: schema.Int64Attribute{
				Description: "The number of organizational units accessible to the user.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *clumioUserResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clumioUserResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan clumioUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	usersAPI := users.NewUsersV1(r.client.ClumioConfig)
	assignedRole := plan.AssignedRole.ValueString()
	email := plan.Email.ValueString()
	fullName := plan.FullName.ValueString()
	organizationalUnitElements := plan.OrganizationalUnitIds.Elements()
	organizationalUnitIds := make([]*string, len(organizationalUnitElements))
	for ind, element := range organizationalUnitElements {
		valString := element.String()
		organizationalUnitIds[ind] = &valString
	}
	res, apiErr := usersAPI.CreateUser(&models.CreateUserV1Request{
		AssignedRole:          &assignedRole,
		Email:                 &email,
		FullName:              &fullName,
		OrganizationalUnitIds: organizationalUnitIds,
	})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error creating Clumio User. Error: %v",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	plan.Id = types.StringValue(*res.Id)
	plan.Inviter = types.StringValue(*res.Inviter)
	plan.IsConfirmed = types.BoolValue(*res.IsConfirmed)
	plan.IsEnabled = types.BoolValue(*res.IsEnabled)
	plan.LastActivityTimestamp = types.StringValue(*res.LastActivityTimestamp)
	plan.OrganizationalUnitCount = types.Int64Value(*res.OrganizationalUnitCount)
	plan.Email = types.StringValue(*res.Email)
	plan.FullName = types.StringValue(*res.FullName)
	if res.AssignedRole != nil {
		plan.AssignedRole = types.StringValue(*res.AssignedRole)
	}
	orgUnitIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.AssignedOrganizationalUnitIds)
	resp.Diagnostics.Append(conversionDiags...)
	plan.OrganizationalUnitIds = orgUnitIds

	// Set the state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioUserResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state clumioUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	usersAPI := users.NewUsersV1(r.client.ClumioConfig)
	userId, perr := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if perr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(invalidUserFmt, state.Id.ValueString()),
			"Invalid user id")
	}

	res, apiErr := usersAPI.ReadUser(userId)
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			state.Id = types.StringValue("")
		}
		resp.Diagnostics.AddError(
			"Error retrieving Clumio User.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	state.Id = types.StringValue(*res.Id)
	state.Inviter = types.StringValue(*res.Inviter)
	state.IsConfirmed = types.BoolValue(*res.IsConfirmed)
	state.IsEnabled = types.BoolValue(*res.IsEnabled)
	state.LastActivityTimestamp = types.StringValue(*res.LastActivityTimestamp)
	state.OrganizationalUnitCount = types.Int64Value(*res.OrganizationalUnitCount)
	state.Email = types.StringValue(*res.Email)
	state.FullName = types.StringValue(*res.FullName)
	if res.AssignedRole != nil {
		state.AssignedRole = types.StringValue(*res.AssignedRole)
	}
	orgUnitIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.AssignedOrganizationalUnitIds)
	resp.Diagnostics.Append(conversionDiags...)
	state.OrganizationalUnitIds = orgUnitIds

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioUserResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan clumioUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state clumioUserResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Email != state.Email {
		resp.Diagnostics.AddError(
			fmt.Sprintf("email is not allowed to be changed"),
			"Error: email is not allowed to be changed")
		return
	}

	usersAPI := users.NewUsersV1(r.client.ClumioConfig)
	updateRequest := &models.UpdateUserV1Request{}
	if !plan.AssignedRole.IsUnknown() &&
		state.AssignedRole != plan.AssignedRole {
		assignedRole := plan.AssignedRole.ValueString()
		updateRequest.AssignedRole = &assignedRole
	}
	if !plan.FullName.IsUnknown() &&
		state.FullName != plan.FullName {
		fullName := plan.FullName.ValueString()
		updateRequest.FullName = &fullName
	}
	if !plan.OrganizationalUnitIds.IsUnknown() {
		added := common.SliceDifferenceAttrValue(plan.OrganizationalUnitIds.Elements(), state.OrganizationalUnitIds.Elements())
		deleted := common.SliceDifferenceAttrValue(state.OrganizationalUnitIds.Elements(), plan.OrganizationalUnitIds.Elements())
		addedStrings := common.GetStringSliceFromAttrValueSlice(added)
		deletedStrings := common.GetStringSliceFromAttrValueSlice(deleted)
		updateRequest.OrganizationalUnitAssignmentUpdates =
			&models.EntityGroupAssignmentUpdatesV1{
				Add:    addedStrings,
				Remove: deletedStrings,
			}
	}
	userId, perr := strconv.ParseInt(plan.Id.ValueString(), 10, 64)
	if perr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(invalidUserFmt, plan.Id.ValueString()),
			"Invalid user id")
	}

	res, apiErr := usersAPI.UpdateUser(userId, updateRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating Clumio User id: %v.", plan.Id.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	plan.Inviter = types.StringValue(*res.Inviter)
	plan.IsConfirmed = types.BoolValue(*res.IsConfirmed)
	plan.IsEnabled = types.BoolValue(*res.IsEnabled)
	plan.LastActivityTimestamp = types.StringValue(*res.LastActivityTimestamp)
	plan.OrganizationalUnitCount = types.Int64Value(*res.OrganizationalUnitCount)
	plan.Email = types.StringValue(*res.Email)
	plan.FullName = types.StringValue(*res.FullName)
	if res.AssignedRole != nil {
		plan.AssignedRole = types.StringValue(*res.AssignedRole)
	}
	orgUnitIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.AssignedOrganizationalUnitIds)
	resp.Diagnostics.Append(conversionDiags...)
	plan.OrganizationalUnitIds = orgUnitIds

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clumioUserResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state clumioUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	usersAPI := users.NewUsersV1(r.client.ClumioConfig)
	userId, perr := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if perr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(invalidUserFmt, state.Id.ValueString()),
			"Invalid user id")
	}

	_, apiErr := usersAPI.DeleteUser(userId)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting Clumio User %v.", userId),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
	}
}

func (r *clumioUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
