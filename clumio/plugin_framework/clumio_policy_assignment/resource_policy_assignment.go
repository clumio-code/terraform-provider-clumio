// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_policy_assignment Terraform resource.
// This resource is used to assign a policy to an entity (for example, a protection_group).

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"net/http"

	sdkPolicyAssignments "github.com/clumio-code/clumio-go-sdk/controllers/policy_assignments"
	sdkPolicyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	sdkProtectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &clumioPolicyAssignmentResource{}
	_ resource.ResourceWithConfigure = &clumioPolicyAssignmentResource{}
)

// clumioPolicyAssignmentResource is the struct backing the clumio_policy_assignment Terraform resource.
// It holds the Clumio API client and any other required state needed to do policy assignment.
type clumioPolicyAssignmentResource struct {
	name   string
	client *common.ApiClient
}

// NewPolicyAssignmentResource creates a new instance of clumioPolicyAssignmentResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewPolicyAssignmentResource() resource.Resource {
	return &clumioPolicyAssignmentResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioPolicyAssignmentResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_policy_assignment"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioPolicyAssignmentResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioPolicyAssignmentResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyAssignmentResourceModel
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

	// Initialize the SDK clients. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyDefinitions := sdkPolicyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	protectionGroups := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	policyAssignments := sdkPolicyAssignments.NewPolicyAssignmentsV1(r.client.ClumioConfig)

	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := policyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read policy with id: %v ", policyId)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	correctPolicyType := false
	for _, operation := range policy.Operations {
		if *operation.ClumioType == protectionGroupBackup {
			correctPolicyType = true
		}
	}
	if !correctPolicyType {
		summary := "Invalid Policy operation."
		detail := fmt.Sprintf(
			"Policy id %s does not contain support protection_group_backup operation", policyId)
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Convert the schema to a Clumio API request to set policy assignments.
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(plan, false)

	// Call the Clumio API to set the policy assignments
	res, apiErr := policyAssignments.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to assign policy %v to entity %v ", policyId,
			*assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// As setting policy assignments is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll task after assigning policy %v to entity %v",
			policyId, *assignment.Entity.Id)
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	readResponse, apiErr := protectionGroups.ReadProtectionGroup(*assignment.Entity.Id)
	if apiErr != nil {
		summary := fmt.Sprintf(
			"Unable to read Protection Group %v.", *assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if readResponse == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}
	if readResponse.ProtectionInfo == nil ||
		*readResponse.ProtectionInfo.PolicyId != policyId {
		summary := "Protection group policy mismatch"
		detail := fmt.Sprintf("Protection group with id: %s does not have policy %s applied",
			*assignment.Entity.Id, policyId)
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Populate all computed fields of the plan including the ID given that the resource is getting created.
	entityType := plan.EntityType.ValueString()
	plan.ID = types.StringValue(
		fmt.Sprintf("%s_%s_%s", *assignment.PolicyId, *assignment.Entity.Id, entityType))
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioPolicyAssignmentResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyAssignmentResourceModel
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

	// Initialize the SDK clients. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyDefinitions := sdkPolicyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	protectionGroups := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)

	// Call the Clumio API to read the policy definition.
	policyId := state.PolicyID.ValueString()
	_, apiErr := policyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		if apiErr.ResponseCode == http.StatusNotFound {
			msgStr := fmt.Sprintf(
				"Clumio Policy with ID %s not found. Removing from state.",
				policyId)
			tflog.Warn(ctx, msgStr)
			resp.State.RemoveResource(ctx)
		} else {
			summary := fmt.Sprintf("Unable to read policy %v.", policyId)
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
		return
	}

	entityType := state.EntityType.ValueString()
	switch entityType {
	case entityTypeProtectionGroup:
		// Call the Clumio API to read the protection group.
		entityId := state.EntityID.ValueString()
		readResponse, apiErr := protectionGroups.ReadProtectionGroup(entityId)
		if apiErr != nil {
			if apiErr.ResponseCode == http.StatusNotFound {
				msgStr := fmt.Sprintf(
					"Clumio Protection Group with ID %s not found. Removing from state.",
					entityId)
				tflog.Warn(ctx, msgStr)
				resp.State.RemoveResource(ctx)
			} else {
				summary := fmt.Sprintf(
					"Unable to read Protection Group %v.", entityId)
				detail := common.ParseMessageFromApiError(apiErr)
				resp.Diagnostics.AddError(summary, detail)
			}
			return
		}
		if readResponse == nil {
			resp.Diagnostics.AddError(
				common.NilErrorMessageSummary, common.NilErrorMessageDetail)
			return
		}
		if readResponse.ProtectionInfo == nil ||
			*readResponse.ProtectionInfo.PolicyId != policyId {
			msgStr := fmt.Sprintf(
				"Protection group with id: %s does not have policy %s applied. Removing from state.",
				entityId, policyId)
			tflog.Warn(ctx, msgStr)
			resp.State.RemoveResource(ctx)
			return
		}
		state.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)
	default:
		summary := "Invalid entityType"
		detail := fmt.Sprintf("The entity type %v is not supported.", entityType)
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state.
func (r *clumioPolicyAssignmentResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan policyAssignmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If OrganizationalUnitID is not empty, then set it as the OrganizationalUnitContext.
	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}

	// Initialize the SDK clients. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyDefinitions := sdkPolicyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	protectionGroups := sdkProtectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	policyAssignments := sdkPolicyAssignments.NewPolicyAssignmentsV1(r.client.ClumioConfig)

	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := policyDefinitions.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read the policy with id : %v", policyId)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	correctPolicyType := false
	for _, operation := range policy.Operations {
		if *operation.ClumioType == protectionGroupBackup {
			correctPolicyType = true
		}
	}
	if !correctPolicyType {
		errMsg := fmt.Sprintf(
			"Policy id %s does not contain support protection_group_backup operation",
			policyId)
		resp.Diagnostics.AddError("Invalid Policy operation.", errMsg)
		return
	}

	// Call the Clumio API to update the policy assignments.
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(plan, false)
	res, apiErr := policyAssignments.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to assign policy %v to entity %v", policyId,
			*assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// As setting policy assignments is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		summary := fmt.Sprintf("Unable to poll task after assigning policy %v to entity %v",
			policyId, *assignment.Entity.Id)
		detail := err.Error()
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	// Call the Clumio API to read the protection group and verify that the policy is assigned
	// to the protection group.
	readResponse, apiErr := protectionGroups.ReadProtectionGroup(*assignment.Entity.Id)
	if apiErr != nil {
		summary := fmt.Sprintf(
			"Unable to read Protection Group %v.", *assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if readResponse.ProtectionInfo == nil ||
		*readResponse.ProtectionInfo.PolicyId != policyId {
		errMsg := fmt.Sprintf(
			"Protection group with id: %s does not have policy %s applied",
			*assignment.Entity.Id, policyId)
		resp.Diagnostics.AddError(errMsg, errMsg)
		return
	}
	plan.OrganizationalUnitID = types.StringPointerValue(readResponse.OrganizationalUnitId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioPolicyAssignmentResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state policyAssignmentResourceModel
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

	// Initialize the SDK client. SDK client initialization is being done after the
	// OrganizationalUnitContext is set in the ClumioConfig so that the API will get executed in the
	// context of the OrganizationalUnit.
	policyAssignments := sdkPolicyAssignments.NewPolicyAssignmentsV1(r.client.ClumioConfig)

	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(state, true)
	// Call the Clumio API to remove the policy assignment.
	_, apiErr := policyAssignments.SetPolicyAssignments(paRequest)
	if apiErr != nil {
		assignment := paRequest.Items[0]
		summary := fmt.Sprintf(
			"Unable to unassign policy from entity %v.", *assignment.Entity.Id)
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
}
