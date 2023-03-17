// Copyright 2023. Clumio, Inc.

// clumio_policy_assignment definition and CRUD implementation.

package clumio_policy_assignment

import (
	"context"
	"fmt"
	"strings"

	policyAssignments "github.com/clumio-code/clumio-go-sdk/controllers/policy_assignments"
	policyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	validators "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &clumioPolicyAssignmentResource{}
	_ resource.ResourceWithConfigure   = &clumioPolicyAssignmentResource{}
	_ resource.ResourceWithImportState = &clumioPolicyAssignmentResource{}
)

type clumioPolicyAssignmentResource struct {
	client *common.ApiClient
}

// NewPolicyAssignmentResource is a helper function to simplify the provider implementation.
func NewPolicyAssignmentResource() resource.Resource {
	return &clumioPolicyAssignmentResource{}
}

type policyAssignmentResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	EntityID             types.String `tfsdk:"entity_id"`
	EntityType           types.String `tfsdk:"entity_type"`
	PolicyID             types.String `tfsdk:"policy_id"`
	OrganizationalUnitID types.String `tfsdk:"organizational_unit_id"`
}

// Schema defines the schema for the data source.
func (r *clumioPolicyAssignmentResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Assignment Resource used to assign (or unassign)" +
			" policies.\n\n NOTE: Currently policy assignment is supported only for" +
			" entity type \"protection_group\".",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaEntityId: schema.StringAttribute{
				Description:   "The entity id.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			schemaEntityType: schema.StringAttribute{
				Description: "The entity type. The supported entity type is" +
					"\"protection_group\".",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					validators.OneOf(entityTypeProtectionGroup),
				},
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
func (r *clumioPolicyAssignmentResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_assignment"
}

// Configure adds the provider configured client to the data source.
func (r *clumioPolicyAssignmentResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

func (r *clumioPolicyAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clumioPolicyAssignmentResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	pa := policyAssignments.NewPolicyAssignmentsV1(r.client.ClumioConfig)
	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	pdv1 := policyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := pdv1.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading the policy with id : %v", policyId),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
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

	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(plan, false)
	res, apiErr := pa.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error assigning policy %v to entity %v", policyId,
				*assignment.Entity.Id),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error assigning policy %v to entity %v.", policyId,
				*assignment.Entity.Id),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	entityType := plan.EntityType.ValueString()
	plan.ID = types.StringValue(
		fmt.Sprintf("%s_%s_%s", *assignment.PolicyId, *assignment.Entity.Id, entityType))
	protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(*assignment.Entity.Id)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error reading Protection Group %v.", *assignment.Entity.Id),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	if *readResponse.ProtectionInfo.PolicyId != policyId {
		errMsg := fmt.Sprintf(
			"Protection group with id: %s does not have policy %s applied",
			*assignment.Entity.Id, policyId)
		resp.Diagnostics.AddError(errMsg, errMsg)
		return
	}
	plan.OrganizationalUnitID = types.StringValue(*readResponse.OrganizationalUnitId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioPolicyAssignmentResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	idSplits := strings.Split(state.ID.ValueString(), "_")
	if len(idSplits) < 3 {
		resp.Diagnostics.AddError("Invalid ID.",
			fmt.Sprintf("Invalid id %s for policy_assignment", state.ID.ValueString()))
	}
	policyId, entityId, entityType :=
		idSplits[0], idSplits[1], strings.Join(idSplits[2:], "_")
	switch entityType {
	case entityTypeProtectionGroup:
		protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
		readResponse, apiErr := protectionGroup.ReadProtectionGroup(entityId)
		if apiErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf(
					"Error reading Protection Group %v.", entityId),
				fmt.Sprintf(errorFmt, string(apiErr.Response)))
			return
		}
		if *readResponse.ProtectionInfo.PolicyId != policyId {
			errMsg := fmt.Sprintf("Protection group with id: %s does not have policy %s applied",
				entityId, policyId)
			resp.Diagnostics.AddError(errMsg, errMsg)
			return
		}
		state.PolicyID = types.StringValue(policyId)
		state.EntityID = types.StringValue(entityId)
		state.EntityType = types.StringValue(entityType)
		state.OrganizationalUnitID = types.StringValue(*readResponse.OrganizationalUnitId)
	default:
		errMsg := fmt.Sprintf("Invalid entityType: %v", entityType)
		resp.Diagnostics.AddError(errMsg, errMsg)
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioPolicyAssignmentResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	pa := policyAssignments.NewPolicyAssignmentsV1(r.client.ClumioConfig)
	// Validation to check if the policy id mentioned supports protection_group_backup operation.
	pdv1 := policyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	policyId := plan.PolicyID.ValueString()
	policy, apiErr := pdv1.ReadPolicyDefinition(policyId, nil)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading the policy with id : %v", policyId),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
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

	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(plan, false)
	res, apiErr := pa.SetPolicyAssignments(paRequest)
	assignment := paRequest.Items[0]
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error assigning policy %v to entity %v", policyId,
				*assignment.Entity.Id),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error assigning policy %v to entity %v.", policyId,
				*assignment.Entity.Id),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(*assignment.Entity.Id)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error reading Protection Group %v.", *assignment.Entity.Id),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	if *readResponse.ProtectionInfo.PolicyId != policyId {
		errMsg := fmt.Sprintf(
			"Protection group with id: %s does not have policy %s applied",
			*assignment.Entity.Id, policyId)
		resp.Diagnostics.AddError(errMsg, errMsg)
		return
	}
	plan.OrganizationalUnitID = types.StringValue(*readResponse.OrganizationalUnitId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clumioPolicyAssignmentResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	pa := policyAssignments.NewPolicyAssignmentsV1(r.client.ClumioConfig)
	paRequest := mapSchemaPolicyAssignmentToClumioPolicyAssignment(state, true)
	_, apiErr := pa.SetPolicyAssignments(paRequest)
	if apiErr != nil {
		assignment := paRequest.Items[0]
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error unassigning policy from entity %v.", *assignment.Entity.Id),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
}

// mapSchemaPolicyAssignmentToClumioPolicyAssignment maps the schema policy assignment
// to the Clumio API request policy assignment.
func mapSchemaPolicyAssignmentToClumioPolicyAssignment(
	model policyAssignmentResourceModel,
	unassign bool) *models.SetPolicyAssignmentsV1Request {
	entityId := model.EntityID.ValueString()
	entityType := model.EntityType.ValueString()
	entity := &models.AssignmentEntity{
		Id:         &entityId,
		ClumioType: &entityType,
	}

	policyId := model.PolicyID.ValueString()
	action := actionAssign
	if unassign {
		policyId = policyIdEmpty
		action = actionUnassign
	}

	assignmentInput := &models.AssignmentInputModel{
		Action:   &action,
		Entity:   entity,
		PolicyId: &policyId,
	}
	return &models.SetPolicyAssignmentsV1Request{
		Items: []*models.AssignmentInputModel{
			assignmentInput,
		},
	}
}

func (r *clumioPolicyAssignmentResource) clearOUContext() {
	r.client.ClumioConfig.OrganizationalUnitContext = ""
}
