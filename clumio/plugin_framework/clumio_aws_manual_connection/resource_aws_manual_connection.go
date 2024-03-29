// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_aws_manual_connection Terraform
// resource. This resource is used in conjunction with clumio_aws_connection and provides the
// externally provisioned AWS resources needed for a Clumio connection to function. This resource
// takes the place of the clumio_post_process_aws_connection resource that is typically called as
// part of the clumio-aws-template module.

package clumio_aws_manual_connection

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// clumioAWSConnectionResource is the struct backing the clumio_aws_connection Terraform resource.
// It holds the Clumio API client and any other required state needed to manage AWS manual
// connections within Clumio.
type clumioAWSManualConnectionResource struct {
	name           string
	client         *common.ApiClient
	sdkConnections sdkclients.AWSConnectionClient
}

// NewClumioAWSManualConnectionResource is a helper function to simplify the provider implementation.
func NewClumioAWSManualConnectionResource() resource.Resource {
	return &clumioAWSManualConnectionResource{}
}

// Metadata returns the resource type name.
func (r *clumioAWSManualConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_aws_manual_connection"
	resp.TypeName = r.name
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioAWSManualConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkConnections = sdkclients.NewAWSConnectionClient(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *clumioAWSManualConnectionResource) Create(
	ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioAWSManualConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	diags = r.createAWSManualConnection(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = res.State.Set(ctx, &plan)
	res.Diagnostics.Append(diags...)
}

// Update updates the resource via the Clumio API and updates the Terraform state. Update only
// happens when there is a change in state of the AWS manual connection.
func (r *clumioAWSManualConnectionResource) Update(
	ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan clumioAWSManualConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Retrieve the schema from the current Terraform state.
	var state clumioAWSManualConnectionResourceModel
	diags = req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	diags = r.updateAWSManualConnection(ctx, &plan, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = res.State.Set(ctx, &plan)
	res.Diagnostics.Append(diags...)
}

// Read does not have an implementation as there is no API to read for clumio_aws_manual_connection.
func (*clumioAWSManualConnectionResource) Read(
	context.Context, resource.ReadRequest, *resource.ReadResponse) {
		// No implementation needed.
}

// Delete does not have an implementation as there is no API to delete for
// clumio_aws_manual_connection.
func (*clumioAWSManualConnectionResource) Delete(
	context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
		// No implementation needed.
}

// clumioSetManualResourcesCommon contains the logic for updating resources of a manual connection
// using Clumio API.
func (r *clumioAWSManualConnectionResource) clumioSetManualResourcesCommon(
	_ context.Context, state clumioAWSManualConnectionResourceModel) diag.Diagnostics {
	accountId := state.AccountId.ValueString()
	awsRegion := state.AwsRegion.ValueString()
	connectionId := accountId + "_" + awsRegion

	// Determine which asset types are enabled fo the connection
	assetsEnabled := []*string{}
	if state.AssetsEnabled.EBS.ValueBool() {
		enabled := EBS
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.S3.ValueBool() {
		enabled := S3
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.RDS.ValueBool() {
		enabled := RDS
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.DynamoDB.ValueBool() {
		enabled := DynamoDB
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.EC2MSSQL.ValueBool() {
		enabled := EC2MSSQL
		assetsEnabled = append(assetsEnabled, &enabled)
	}

	// Convert the schema into a Clumio API request, containing the enabled asset types and stack ARNs
	// to the manually configured resources
	req := models.UpdateAwsConnectionV1Request{
		AssetTypesEnabled: assetsEnabled,
		Resources: &models.Resources{
			ClumioIamRoleArn:     state.Resources.ClumioIAMRoleArn.ValueStringPointer(),
			ClumioEventPubArn:    state.Resources.ClumioEventPubArn.ValueStringPointer(),
			ClumioSupportRoleArn: state.Resources.ClumioSupportRoleArn.ValueStringPointer(),
			EventRules: &models.EventRules{
				CloudtrailRuleArn: state.Resources.EventRules.CloudtrailRuleArn.ValueStringPointer(),
				CloudwatchRuleArn: state.Resources.EventRules.CloudwatchRuleArn.ValueStringPointer(),
			},
			ServiceRoles: &models.ServiceRoles{
				S3: &models.S3ServiceRoles{
					ContinuousBackupsRoleArn: state.Resources.ServiceRoles.S3.ContinuousBackupsRoleArn.ValueStringPointer(),
				},
				Mssql: &models.MssqlServiceRoles{
					Ec2SsmInstanceProfileArn: state.Resources.ServiceRoles.Mssql.Ec2SsmInstanceProfileArn.ValueStringPointer(),
					SsmNotificationRoleArn:   state.Resources.ServiceRoles.Mssql.SsmNotificationRoleArn.ValueStringPointer(),
				},
			},
		},
	}

	// Call the Clumio API to update the AWS manual connection.
	_, apiErr := r.sdkConnections.UpdateAwsConnection(connectionId, req)
	if apiErr != nil {
		diagnostics := diag.Diagnostics{}
		summary := fmt.Sprintf("Unable to update resources for %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diagnostics.AddError(summary, detail)
		return diagnostics
	}
	return nil
}
