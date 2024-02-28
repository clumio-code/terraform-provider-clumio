// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_auto_user_provisioning_rules Terraform
// resource. This resource is used to create auto user provisioning rules to determine the roles
// and organizational-units to be assigned to the users.

package clumio_auto_user_provisioning_rule

import (
	"context"
	"fmt"
	"net/http"

	sdkAUPRules "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_rules"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &autoUserProvisioningRuleResource{}
	_ resource.ResourceWithConfigure   = &autoUserProvisioningRuleResource{}
	_ resource.ResourceWithImportState = &autoUserProvisioningRuleResource{}
)

// autoUserProvisioningRuleResource is the struct backing the clumio_auto_user_provisioning_rules
// Terraform resource. It holds the Clumio API client and any other required state needed to set up
// the auto user provisioning rules to determine the roles and organizational-units to be assigned
// to the user.
type autoUserProvisioningRuleResource struct {
	name        string
	client      *common.ApiClient
	sdkAUPRules sdkAUPRules.AutoUserProvisioningRulesV1Client
}

// NewAutoUserProvisioningRuleResource creates a new instance of autoUserProvisioningRuleResource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewAutoUserProvisioningRuleResource() resource.Resource {
	return &autoUserProvisioningRuleResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *autoUserProvisioningRuleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_auto_user_provisioning_rule"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *autoUserProvisioningRuleResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkAUPRules = sdkAUPRules.NewAutoUserProvisioningRulesV1(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *autoUserProvisioningRuleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the current Terraform plan.
	var plan autoUserProvisioningRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the schema to a Clumio API request to create an auto user provisioning rule.
	ouIds := make([]*string, 0)
	conversionDiags := plan.OrganizationalUnitIDs.ElementsAs(ctx, &ouIds, false)
	resp.Diagnostics.Append(conversionDiags...)
	provision := &models.RuleProvision{
		RoleId:                plan.RoleID.ValueStringPointer(),
		OrganizationalUnitIds: ouIds,
	}
	auprRequest := &models.CreateAutoUserProvisioningRuleV1Request{
		Name:      plan.Name.ValueStringPointer(),
		Condition: plan.Condition.ValueStringPointer(),
		Provision: provision,
	}

	// Call the Clumio API to create the auto user provisioning rule.
	res, apiErr := r.sdkAUPRules.CreateAutoUserProvisioningRule(auprRequest)
	if apiErr != nil {
		summary := fmt.Sprintf(
			"Error creating auto user provisioning rule %v.", plan.Name.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Set the ID in the schema.
	plan.ID = types.StringPointerValue(res.RuleId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *autoUserProvisioningRuleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state autoUserProvisioningRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the auto user provisioning rule.
	res, apiErr := r.sdkAUPRules.ReadAutoUserProvisioningRule(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			summary := fmt.Sprintf(
				"%s (ID: %v) not found. Removing from state", r.name, state.ID.ValueString())
			tflog.Warn(ctx, summary)
			resp.State.RemoveResource(ctx)
		} else {
			summary := fmt.Sprintf(
				"Unable to retrieve auto user provisioning rule %v.", state.Name.ValueString())
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
	state.RoleID = types.StringPointerValue(res.Provision.RoleId)
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, res.Provision.OrganizationalUnitIds)
	resp.Diagnostics.Append(conversionDiags...)
	state.OrganizationalUnitIDs = ouIds

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *autoUserProvisioningRuleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the current Terraform plan.
	var plan autoUserProvisioningRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ouIds := make([]*string, 0)
	conversionDiags := plan.OrganizationalUnitIDs.ElementsAs(ctx, &ouIds, false)
	resp.Diagnostics.Append(conversionDiags...)
	provision := &models.RuleProvision{
		RoleId:                plan.RoleID.ValueStringPointer(),
		OrganizationalUnitIds: ouIds,
	}
	auprRequest := &models.UpdateAutoUserProvisioningRuleV1Request{
		Name:      plan.Name.ValueStringPointer(),
		Condition: plan.Condition.ValueStringPointer(),
		Provision: provision,
	}

	// Call the Clumio API to update the auto user provisioning rule.
	res, apiErr := r.sdkAUPRules.UpdateAutoUserProvisioningRule(plan.ID.ValueString(), auprRequest)
	if apiErr != nil {
		summary := fmt.Sprintf(
			"Error updating auto user provisioning rule %v.", plan.Name.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *autoUserProvisioningRuleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state autoUserProvisioningRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to delete the auto user provisioning rule.
	_, apiErr := r.sdkAUPRules.DeleteAutoUserProvisioningRule(state.ID.ValueString())
	if apiErr != nil {
		summary := fmt.Sprintf(
			"Error deleting auto user provisioning rule %v.", state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
	}
}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *autoUserProvisioningRuleResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
