// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_aws_connection Terraform resource.
// This resource is used to connect AWS accounts to Clumio.

package clumio_aws_connection

import (
	"context"
	"fmt"
	"net/http"

	sdkConnections "github.com/clumio-code/clumio-go-sdk/controllers/aws_connections"
	sdkEnvironments "github.com/clumio-code/clumio-go-sdk/controllers/aws_environments"
	sdkOrgUnits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the following Resource interfaces.
var (
	_ resource.Resource                = &clumioAWSConnectionResource{}
	_ resource.ResourceWithConfigure   = &clumioAWSConnectionResource{}
	_ resource.ResourceWithImportState = &clumioAWSConnectionResource{}
)

// clumioAWSConnectionResource is the struct backing the clumio_aws_connection Terraform resource.
// It holds the Clumio API client and any other required state needed to connect AWS accounts to
// Clumio.
type clumioAWSConnectionResource struct {
	name            string
	client          *common.ApiClient
	sdkConnections  sdkConnections.AwsConnectionsV1Client
	sdkEnvironments sdkEnvironments.AwsEnvironmentsV1Client
	sdkOrgUnits     sdkOrgUnits.OrganizationalUnitsV1Client
}

// NewClumioAWSConnectionResource creates a new instance of clumioAWSConnectionResource. Its
// attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewClumioAWSConnectionResource() resource.Resource {
	return &clumioAWSConnectionResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *clumioAWSConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	r.name = req.ProviderTypeName + "_aws_connection"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioAWSConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkConnections = sdkConnections.NewAwsConnectionsV1(r.client.ClumioConfig)
	r.sdkEnvironments = sdkEnvironments.NewAwsEnvironmentsV1(r.client.ClumioConfig)
	r.sdkOrgUnits = sdkOrgUnits.NewOrganizationalUnitsV1(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioAWSConnectionResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.OrganizationalUnitID.IsNull() {
		_, err := getOrgUnitForConnection(ctx, r, &plan)
		if err != nil {
			summary := fmt.Sprintf("invalid %s", schemaOrganizationalUnitId)
			detail := err.Error()
			resp.Diagnostics.AddError(summary, detail)
			return
		}
	}
	// Convert the schema to a Clumio API request to create an AWS connection.
	createReq := &models.CreateAwsConnectionV1Request{
		AccountNativeId:      plan.AccountNativeID.ValueStringPointer(),
		AwsRegion:            plan.AWSRegion.ValueStringPointer(),
		Description:          plan.Description.ValueStringPointer(),
		OrganizationalUnitId: plan.OrganizationalUnitID.ValueStringPointer(),
	}

	// Call the Clumio API to create the AWS connection.
	res, apiErr := r.sdkConnections.CreateAwsConnection(createReq)
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

	// Convert the Clumio API response back to a schema and populate all computed fields of the plan
	// including the ID given that the resource is getting created.
	plan.ID = types.StringPointerValue(res.Id)
	plan.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	plan.ConnectionStatus = types.StringPointerValue(res.ConnectionStatus)
	plan.Token = types.StringPointerValue(res.Token)
	plan.Namespace = types.StringPointerValue(res.Namespace)
	plan.ClumioAWSAccountID = types.StringPointerValue(res.ClumioAwsAccountId)
	plan.ClumioAWSRegion = types.StringPointerValue(res.ClumioAwsRegion)
	setExternalId(&plan, res.ExternalId, res.Token)
	setDataPlaneAccountId(&plan, res.DataPlaneAccountId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves the resource from the Clumio API and sets the Terraform state.
func (r *clumioAWSConnectionResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to read the AWS connection.
	returnExternalId := "true"
	res, apiErr := r.sdkConnections.ReadAwsConnection(state.ID.ValueString(), &returnExternalId)
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
	state.AccountNativeID = types.StringPointerValue(res.AccountNativeId)
	state.AWSRegion = types.StringPointerValue(res.AwsRegion)

	// Since the Description field is optional, it should only be populated if it initially
	// contained a non-null value or if there is a specific value that needs to be assigned.
	description := types.StringPointerValue(res.Description)
	if !state.Description.IsNull() || description.ValueString() != "" {
		state.Description = description
	}
	state.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	state.ConnectionStatus = types.StringPointerValue(res.ConnectionStatus)
	state.Token = types.StringPointerValue(res.Token)
	state.Namespace = types.StringPointerValue(res.Namespace)
	state.ClumioAWSAccountID = types.StringPointerValue(res.ClumioAwsAccountId)
	state.ClumioAWSRegion = types.StringPointerValue(res.ClumioAwsRegion)
	setExternalId(&state, res.ExternalId, res.Token)
	setDataPlaneAccountId(&state, res.DataPlaneAccountId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource via the Clumio API and updates the Terraform state. NOTE that the
// update for OU is a separate API call than the update for the AWS connection. Due to this it is
// possible for one portion of an update to go through but not the other. However, the update is
// idemptent so if a portion of the update fails, the next apply will attempt to update the failed
// portion again.
func (r *clumioAWSConnectionResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSConnectionResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Organizational Unit (OU) associated with the AWS connection if it has been
	// explicitly set and is different from the current OU.
	if !plan.OrganizationalUnitID.IsUnknown() &&
		plan.OrganizationalUnitID != state.OrganizationalUnitID {

		err := updateOrgUnitForConnection(ctx, r, &plan, &state)
		if err != nil {
			summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, state.ID.ValueString())
			detail := err.Error()
			resp.Diagnostics.AddError(summary, detail)
			return
		}
		// "description" is the only field within the schema that can cause an update to the
		// resource. As such, if this has not changed, update the Terraform state with the change
		// to the Organization Unit and return.
		if plan.Description == state.Description {
			diags = resp.State.Set(ctx, plan)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// Call the Clumio API to update the AWS connection. The "Description" parameter, while optional
	// in the REST API, is deliberately provided to ensure the update process is executed, even in
	// the absence of a specified description. This approach is necessary as the API's response
	// varies when a description is omitted.
	description := plan.Description.ValueString()
	updateReq := models.UpdateAwsConnectionV1Request{
		Description: &description,
	}
	res, apiErr := r.sdkConnections.UpdateAwsConnection(plan.ID.ValueString(), updateReq)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to update %s (ID: %v)", r.name, state.ID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		resp.Diagnostics.AddError(summary, detail)
		return
	}
	if res == nil {
		resp.Diagnostics.AddError(
			common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return
	}

	// Convert the Clumio API response back to a schema and populate all computed fields of the
	// plan. ID however is not updated given that it is the field used to denote which resource to
	// update in the backend. Additionally the external ID is currently not returned during an
	// update call and thus is not updated below. This is okay however as the external ID is not
	// expected to change once a connection is created.
	plan.OrganizationalUnitID = types.StringPointerValue(res.OrganizationalUnitId)
	plan.ConnectionStatus = types.StringPointerValue(res.ConnectionStatus)
	plan.Token = types.StringPointerValue(res.Token)
	plan.Namespace = types.StringPointerValue(res.Namespace)
	plan.ClumioAWSAccountID = types.StringPointerValue(res.ClumioAwsAccountId)
	plan.ClumioAWSRegion = types.StringPointerValue(res.ClumioAwsRegion)
	setDataPlaneAccountId(&plan, res.DataPlaneAccountId)

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource via the Clumio API and removes the Terraform state.
func (r *clumioAWSConnectionResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Clumio API to delete the AWS connection.
	_, apiErr := r.sdkConnections.DeleteAwsConnection(state.ID.ValueString())
	if apiErr != nil {
		if apiErr.ResponseCode != http.StatusNotFound {
			summary := fmt.Sprintf("Unable to delete %s (ID: %v)", r.name, state.ID.ValueString())
			detail := common.ParseMessageFromApiError(apiErr)
			resp.Diagnostics.AddError(summary, detail)
		}
	}

}

// ImportState retrieves the resource via the Clumio API and sets the Terraform state. The import
// is done by the ID of the resource.
func (r *clumioAWSConnectionResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
