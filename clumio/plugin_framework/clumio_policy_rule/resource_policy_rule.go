// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_policy_rule Terraform resource.
// This resource facilitates the creation of rules that specify which assets are to be protected
// by particular policies.

package clumio_policy_rule

import (
	"context"
	"fmt"
	"net/http"

	sdkPolicyRules "github.com/clumio-code/clumio-go-sdk/controllers/policy_rules"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &policyRuleResource{}
	_ resource.ResourceWithConfigure   = &policyRuleResource{}
	_ resource.ResourceWithImportState = &policyRuleResource{}
)

// policyRuleResource is the struct backing the clumio_policy_rule Terraform resource. It holds the
// Clumio API client and any other required state needed to create a Clumio policy rule.
type policyRuleResource struct {
	name   string
	client *common.ApiClient
}

// NewPolicyRuleResource creates a new instance of policyRuleResource. Its attributes are
// initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewPolicyRuleResource() resource.Resource {
	return &policyRuleResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *policyRuleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_policy_rule"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *policyRuleResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *policyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *policyRuleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
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

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyRules := sdkPolicyRules.NewPolicyRulesV1(r.client.ClumioConfig)

	// Convert the schema to a Clumio API request to create a policy rule.
	priority := &models.RulePriority{
		BeforeRuleId: plan.BeforeRuleID.ValueStringPointer(),
	}
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: plan.PolicyID.ValueStringPointer(),
		},
	}
	prRequest := &models.CreatePolicyRuleV1Request{
		Action:    action,
		Condition: plan.Condition.ValueStringPointer(),
		Name:      plan.Name.ValueStringPointer(),
		Priority:  priority,
	}

	// Call the Clumio API to create the policy rule.
	res, apiErr := policyRules.CreatePolicyRule(prRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to create %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return

	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll %s (Name: %v) for creation", r.name, plan.Name.ValueString())
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.ID = types.StringPointerValue(res.Rule.Id)
	plan.OrganizationalUnitID = types.StringPointerValue(res.Rule.OrganizationalUnitId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *policyRuleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the OrganizationalUnitID is specified, then execute the API in that
	// OrganizationalUnit context.
	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyRules := sdkPolicyRules.NewPolicyRulesV1(r.client.ClumioConfig)

	priority := &models.RulePriority{
		BeforeRuleId: plan.BeforeRuleID.ValueStringPointer(),
	}
	policyId := plan.PolicyID.ValueString()
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: &policyId,
		},
	}
	prRequest := &models.UpdatePolicyRuleV1Request{
		Action:    action,
		Condition: plan.Condition.ValueStringPointer(),
		Name:      plan.Name.ValueStringPointer(),
		Priority:  priority,
	}

	// Call the Clumio API to update the policy rule.
	res, apiErr := policyRules.UpdatePolicyRule(plan.ID.ValueString(), prRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, plan.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// As the update of a policy rule is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll %s (ID: %v) for update", r.name, plan.ID.ValueString())
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the
	// plan. ID however is not updated given that it is the field used to denote which resource to
	// update in the backend.
	plan.OrganizationalUnitID = types.StringPointerValue(res.Rule.OrganizationalUnitId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *policyRuleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the OrganizationalUnitID is specified, then execute the API in that
	// OrganizationalUnit context.
	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyRules := sdkPolicyRules.NewPolicyRulesV1(r.client.ClumioConfig)

	// Call the Clumio API to read the policy rule.
	res, apiErr := policyRules.ReadPolicyRule(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf("%s (ID: %v) not found. Removing from state", r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			resp.State.RemoveResource(ctx)
		} else {
			summary := fmt.Sprintf("Unable to read %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and update the state. In addition to
	// computed fields, all fields are populated from the API response in case any values have been
	// changed externally. ID is not updated however given that it is the field used to query the
	// resource from the backend.
	state.Name = types.StringPointerValue(res.Name)
	state.Condition = types.StringPointerValue(res.Condition)
	if res.Priority != nil {
		state.BeforeRuleID = types.StringPointerValue(res.Priority.BeforeRuleId)
	}
	state.PolicyID = types.StringPointerValue(res.Action.AssignPolicy.PolicyId)
	state.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes it from the Terraform state.
func (r *policyRuleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the OrganizationalUnitID is specified, then execute the API in that
	// OrganizationalUnit context.
	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyRules := sdkPolicyRules.NewPolicyRulesV1(r.client.ClumioConfig)

	// Call the Clumio API to delete the policy rule.
	res, apiErr := policyRules.DeletePolicyRule(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}
	// As the delete of a policy rule is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll %s (ID: %v) for deletion", r.name, state.ID.ValueString())
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}
}
