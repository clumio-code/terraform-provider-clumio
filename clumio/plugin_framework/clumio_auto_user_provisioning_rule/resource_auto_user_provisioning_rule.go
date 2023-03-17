// Copyright 2023. Clumio, Inc.

// clumio_auto_user_provisioning_rule definition and CRUD implementation.
package clumio_auto_user_provisioning_rule

import (
	"context"
	"fmt"

	aupRules "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_rules"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &autoUserProvisioningRuleResource{}
	_ resource.ResourceWithConfigure = &autoUserProvisioningRuleResource{}
)

// autoUserProvisioningRuleResource is the resource implementation.
type autoUserProvisioningRuleResource struct {
	client *common.ApiClient
}

// NewAutoUserProvisioningRuleResource is a helper function to simplify the provider implementation.
func NewAutoUserProvisioningRuleResource() resource.Resource {
	return &autoUserProvisioningRuleResource{}
}

// autoUserProvisioningRuleResource model
type autoUserProvisioningRuleResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Condition             types.String `tfsdk:"condition"`
	RoleID                types.String `tfsdk:"role_id"`
	OrganizationalUnitIDs types.Set    `tfsdk:"organizational_unit_ids"`
}

// Metadata returns the resource type name.
func (r *autoUserProvisioningRuleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_user_provisioning_rule"
}

// Schema defines the schema for the data source.
func (r *autoUserProvisioningRuleResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Auto User Provisioning Rule Resource used to determine " +
			"the Role and Organizational Units to be assigned to a user based on their groups.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Auto User Provisioning Rule Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the auto user provisioning rule.",
				Required:    true,
			},
			schemaCondition: schema.StringAttribute{
				Description: "The condition of the auto user provisioning rule. Possible conditions include:\n" +
					"\t1) `This group` - User must belong to the specified group\n" +
					"\t2) `ANY of these groups` - User must belong to at least one of the specified groups\n" +
					"\t3) `ALL of these groups` - User must belong to all the specified groups\n" +
					"\t4) `Group CONTAINS this keyword` - User's group must contain the specified keyword\n" +
					"\t5) `Group CONTAINS ANY of these keywords` - User's group must contain at least one of the specified keywords\n" +
					"\t6) `Group CONTAINS ALL of these keywords` - User's group must contain all the specified keywords\n",
				Required: true,
			},
			schemaRoleId: schema.StringAttribute{
				Description: "The role ID of the role to be assigned to the user.",
				Required:    true,
			},
			schemaOrganizationalUnitIds: schema.SetAttribute{
				Description: "The Clumio-assigned IDs of the organizational units " +
					"to be assigned to the user.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *autoUserProvisioningRuleResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *autoUserProvisioningRuleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan autoUserProvisioningRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aupr := aupRules.NewAutoUserProvisioningRulesV1(r.client.ClumioConfig)
	name := plan.Name.ValueString()
	condition := plan.Condition.ValueString()
	roleId := plan.RoleID.ValueString()
	ouIds := make([]*string, 0)
	conversionDiags := plan.OrganizationalUnitIDs.ElementsAs(ctx, &ouIds, false)
	resp.Diagnostics.Append(conversionDiags...)
	provision := &models.RuleProvision{
		RoleId:                &roleId,
		OrganizationalUnitIds: ouIds,
	}
	auprRequest := &models.CreateAutoUserProvisioningRuleV1Request{
		Name:      &name,
		Condition: &condition,
		Provision: provision,
	}
	res, apiErr := aupr.CreateAutoUserProvisioningRule(auprRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating auto user provisioning rule %v.", name),
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
		return
	}

	plan.ID = types.StringValue(*res.RuleId)

	// Set the state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *autoUserProvisioningRuleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state autoUserProvisioningRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aupr := aupRules.NewAutoUserProvisioningRulesV1(r.client.ClumioConfig)
	res, apiErr := aupr.ReadAutoUserProvisioningRule(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error retrieving auto user provisioning rule %v.", state.Name.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
		return
	}

	state.Name = types.StringValue(*res.Name)
	state.Condition = types.StringValue(*res.Condition)
	state.RoleID = types.StringValue(*res.Provision.RoleId)
	ouIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.Provision.OrganizationalUnitIds)
	resp.Diagnostics.Append(conversionDiags...)
	state.OrganizationalUnitIDs = ouIds

	// Set refreshed state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *autoUserProvisioningRuleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan autoUserProvisioningRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aupr := aupRules.NewAutoUserProvisioningRulesV1(r.client.ClumioConfig)
	name := plan.Name.ValueString()
	condition := plan.Condition.ValueString()
	roleId := plan.RoleID.ValueString()
	ouIds := make([]*string, 0)
	conversionDiags := plan.OrganizationalUnitIDs.ElementsAs(ctx, &ouIds, false)
	resp.Diagnostics.Append(conversionDiags...)
	provision := &models.RuleProvision{
		RoleId:                &roleId,
		OrganizationalUnitIds: ouIds,
	}
	auprRequest := &models.UpdateAutoUserProvisioningRuleV1Request{
		Name:      &name,
		Condition: &condition,
		Provision: provision,
	}

	res, apiErr := aupr.UpdateAutoUserProvisioningRule(plan.ID.ValueString(), auprRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating auto user provisioning rule %v.", name),
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
		return
	}

	plan.Name = types.StringValue(*res.Name)
	plan.Condition = types.StringValue(*res.Condition)
	plan.RoleID = types.StringValue(*res.Provision.RoleId)
	orgUnitIds, conversionDiags := types.SetValueFrom(ctx, types.StringType, res.Provision.OrganizationalUnitIds)
	resp.Diagnostics.Append(conversionDiags...)
	plan.OrganizationalUnitIDs = orgUnitIds

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *autoUserProvisioningRuleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state autoUserProvisioningRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aupr := aupRules.NewAutoUserProvisioningRulesV1(r.client.ClumioConfig)
	_, apiErr := aupr.DeleteAutoUserProvisioningRule(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting auto user provisioning rule %v.", state.ID.String()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
	}
}
