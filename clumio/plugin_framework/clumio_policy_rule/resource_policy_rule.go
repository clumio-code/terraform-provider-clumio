// Copyright 2023. Clumio, Inc.

// clumio_policy_rule definition and CRUD implementation.
package clumio_policy_rule

import (
	"context"
	"fmt"

	policyRules "github.com/clumio-code/clumio-go-sdk/controllers/policy_rules"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &policyRuleResource{}
	_ resource.ResourceWithConfigure   = &policyRuleResource{}
	_ resource.ResourceWithImportState = &policyRuleResource{}
)

type policyRuleResource struct {
	client *common.ApiClient
}

// NewPolicyRuleResource is a helper function to simplify the provider implementation.
func NewPolicyRuleResource() resource.Resource {
	return &policyRuleResource{}
}

type policyRuleResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Condition            types.String `tfsdk:"condition"`
	BeforeRuleID         types.String `tfsdk:"before_rule_id"`
	PolicyID             types.String `tfsdk:"policy_id"`
	OrganizationalUnitID types.String `tfsdk:"organizational_unit_id"`
}

// Schema defines the schema for the data source.
func (r *policyRuleResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Rule Resource used to determine how" +
			" a policy should be assigned to assets.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Policy Rule Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the policy rule.",
				Required:    true,
			},
			schemaCondition: schema.StringAttribute{
				Description: "The condition of the policy rule. Possible conditions include: " +
					"1) `entity_type` is required and supports `$eq` and `$in` filters. " +
					"2) `aws_account_native_id` and `aws_region` are optional and both support " +
					"`$eq` and `$in` filters. " +
					"3) `aws_tag` is optional and supports `$eq`, `$in`, `$all`, and `$contains` " +
					"filters.",
				Required: true,
			},
			schemaBeforeRuleId: schema.StringAttribute{
				Description: "The policy rule ID before which this policy rule should be " +
					"inserted. An empty value will set the rule to have lowest priority. " +
					"NOTE: If in the Global Organizational Unit, rules can also be prioritized " +
					"against two virtual rules maintained by the system: `asset-level-rule` and " +
					"`child-ou-rule`. `asset-level-rule` corresponds to the priority of Direct " +
					"Assignments (when a policy is applied directly to an asset) whereas " +
					"`child-ou-rule` corresponds to the priority of rules created by child " +
					"organizational units.",
				Required: true,
			},
			schemaPolicyId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the policy. ",
				Required:    true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the organizational unit" +
					" to use as the context for assigning the policy.",
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Metadata returns the data source type name.
func (r *policyRuleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_rule"
}

// Configure adds the provider configured client to the data source.
func (r *policyRuleResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

func (r *policyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *policyRuleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	pr := policyRules.NewPolicyRulesV1(r.client.ClumioConfig)
	condition := plan.Condition.ValueString()
	name := plan.Name.ValueString()
	beforeRuleId := plan.BeforeRuleID.ValueString()
	priority := &models.RulePriority{
		BeforeRuleId: &beforeRuleId,
	}
	policyId := plan.PolicyID.ValueString()
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: &policyId,
		},
	}
	prRequest := &models.CreatePolicyRuleV1Request{
		Action:    action,
		Condition: &condition,
		Name:      &name,
		Priority:  priority,
	}
	res, apiErr := pr.CreatePolicyRule(prRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error starting task to create policy rule %v.", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return

	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating policy rule %v.", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}

	plan.ID = types.StringValue(*res.Rule.Id)
	plan.OrganizationalUnitID = types.StringValue(*res.Rule.OrganizationalUnitId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *policyRuleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	pr := policyRules.NewPolicyRulesV1(r.client.ClumioConfig)
	condition := plan.Condition.ValueString()
	name := plan.Name.ValueString()
	beforeRuleId := plan.BeforeRuleID.ValueString()
	priority := &models.RulePriority{
		BeforeRuleId: &beforeRuleId,
	}
	policyId := plan.PolicyID.ValueString()
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: &policyId,
		},
	}
	prRequest := &models.UpdatePolicyRuleV1Request{
		Action:    action,
		Condition: &condition,
		Name:      &name,
		Priority:  priority,
	}
	res, apiErr := pr.UpdatePolicyRule(plan.ID.ValueString(), prRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error starting task to update policy rule %v.", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating policy rule %v.", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	plan.OrganizationalUnitID = types.StringValue(*res.Rule.OrganizationalUnitId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *policyRuleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	pr := policyRules.NewPolicyRulesV1(r.client.ClumioConfig)

	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	res, apiErr := pr.ReadPolicyRule(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error retrieving policy rule %v.", state.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	state.Name = types.StringValue(*res.Name)
	state.Condition = types.StringValue(*res.Condition)
	if res.Priority != nil && res.Priority.BeforeRuleId != nil {
		state.BeforeRuleID = types.StringValue(*res.Priority.BeforeRuleId)
	}
	state.PolicyID = types.StringValue(*res.Action.AssignPolicy.PolicyId)
	state.OrganizationalUnitID = types.StringValue(*res.OrganizationalUnitId)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *policyRuleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	pr := policyRules.NewPolicyRulesV1(r.client.ClumioConfig)
	res, apiErr := pr.DeletePolicyRule(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error starting task to delete policy rule %v.",
				state.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting policy rule %v.", state.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
}

func (r *policyRuleResource) clearOUContext() {
	r.client.ClumioConfig.OrganizationalUnitContext = ""
}
