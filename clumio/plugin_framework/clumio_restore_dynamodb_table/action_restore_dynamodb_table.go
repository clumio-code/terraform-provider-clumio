// Copyright 2026. Clumio, Inc.

// This file holds the action implementation for the clumio_restore_dynamodb_table Terraform
// action. This action is used to initiate a restore of a DynamoDB table from a backup.

package clumio_restore_dynamodb_table

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ action.Action               = &clumioRestoreDynamoDBTableAction{}
	_ action.ActionWithConfigure  = &clumioRestoreDynamoDBTableAction{}
	_ action.ActionWithModifyPlan = &clumioRestoreDynamoDBTableAction{}
)

// clumioRestoreDynamoDBTableAction is the struct backing the clumio_restore_dynamodb_table
// Terraform action. It holds the Clumio API client and any other required state needed to restore
// a DynamoDB table within Clumio.
type clumioRestoreDynamoDBTableAction struct {
	name          string
	client        *common.ApiClient
	restoreClient sdkclients.RestoredDynamoDBTableClient
}

// NewClumioRestoreDynamoDBTableAction creates a new instance of
// clumioRestoreDynamoDBTableAction. Its attributes are initialized later by Terraform via Metadata
// and Configure once the Provider is initialized.
func NewClumioRestoreDynamoDBTableAction() action.Action {
	return &clumioRestoreDynamoDBTableAction{}
}

// Metadata returns the name of the action type. This is used by Terraform configurations to
// instantiate the action.
func (r *clumioRestoreDynamoDBTableAction) Metadata(
	_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	r.name = req.ProviderTypeName + "_restore_dynamodb_table"
	resp.TypeName = r.name
}

// Configure sets up the action with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioRestoreDynamoDBTableAction) Configure(
	_ context.Context, req action.ConfigureRequest, _ *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.restoreClient = sdkclients.NewRestoredDynamoDBTableClient(r.client.ClumioConfig)
}

// ModifyPlan validates the action configuration during the plan phase so that invalid restore
// source combinations are caught before the apply phase. Attributes whose values are not yet
// known during the plan phase are skipped and validated again during the invoke phase.
func (r *clumioRestoreDynamoDBTableAction) ModifyPlan(
	ctx context.Context, req action.ModifyPlanRequest, resp *action.ModifyPlanResponse) {

	var model restoreDynamoDBTableActionModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateRestoreSource(ctx, &model)...)
}

// Invoke initiates the restore of the DynamoDB table via the Clumio API and reports the identifier
// of the Clumio task tracking the restore as a progress event. It does not wait for the restore to
// complete.
func (r *clumioRestoreDynamoDBTableAction) Invoke(
	ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {

	// Retrieve the schema from the current Terraform configuration.
	var model restoreDynamoDBTableActionModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	taskId, diags := r.restoreDynamoDBTable(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The restore is intentionally fire-and-forget: report the task ID so that the restore can be
	// tracked via the Clumio API or UI and return without polling the task. The message is kept
	// short so the task ID stays easy to spot in the Terraform output.
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Restore task ID: %s", taskId),
	})
}
