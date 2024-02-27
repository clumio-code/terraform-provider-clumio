// Copyright 2023. Clumio, Inc.

// This file holds the resource implementation for the clumio_post_process_aws_connection Terraform
// resource. This resource is used to send the necessary information required by Clumio to
// post-process an AWS connection after the necessary resources have been created. This resource
// should only be invoked as part of the aws-template module.

package clumio_post_process_aws_connection

import (
	"context"
	"encoding/json"
	"fmt"

	sdkPostProcessConn "github.com/clumio-code/clumio-go-sdk/controllers/post_process_aws_connection"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sourceConfigInfo struct {
	sourceKey string
	isConfig  bool
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &postProcessAWSConnectionResource{}
	_ resource.ResourceWithConfigure = &postProcessAWSConnectionResource{}
)

// postProcessAWSConnectionResource is the resource implementation.
type postProcessAWSConnectionResource struct {
	client              *common.ApiClient
	sdkPostPsConnection sdkPostProcessConn.PostProcessAwsConnectionV1Client
}

// NewPostProcessAWSConnectionResource creates a new instance of postProcessAWSConnectionResource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider is
// initialized.
func NewPostProcessAWSConnectionResource() resource.Resource {
	return &postProcessAWSConnectionResource{}
}

// Metadata returns the name of the resource type. This is used by Terraform configurations to
// instantiate the resource.
func (r *postProcessAWSConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_process_aws_connection"
}

// Configure sets up the resource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *postProcessAWSConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
	r.sdkPostPsConnection = sdkPostProcessConn.NewPostProcessAwsConnectionV1(r.client.ClumioConfig)
}

// Create creates the resource via the Clumio API and sets the initial Terraform state.
func (r *postProcessAWSConnectionResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan postProcessAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.clumioPostProcessAWSConnectionCommon(ctx, plan, "Create")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId := plan.AccountID.ValueString()
	awsRegion := plan.Region.ValueString()
	token := plan.Token.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf("%v/%v/%v", accountId, awsRegion, token))

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read does not have an implementation as there is no API to read for post process aws connection.
func (r *postProcessAWSConnectionResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource via the Clumio API and removes the Terraform state.
func (r *postProcessAWSConnectionResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// Retrieve the schema from the Terraform plan.
	var plan postProcessAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.clumioPostProcessAWSConnectionCommon(ctx, plan, "Update")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
func (r *postProcessAWSConnectionResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve the schema from the Terraform state.
	var state postProcessAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.clumioPostProcessAWSConnectionCommon(ctx, state, "Delete")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// clumioPostProcessAWSConnectionCommon contains the common logic for all CRUD operations
// of PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) clumioPostProcessAWSConnectionCommon(
	_ context.Context, model postProcessAWSConnectionResourceModel, eventType string) diag.Diagnostics {

	schemaPropertiesElements := model.Properties.Elements()
	propertiesMap := make(map[string]*string)
	for key, val := range schemaPropertiesElements {
		valStr := val.(types.String).ValueString()
		propertiesMap[key] = &valStr
	}

	// Using the schema properties in the model, create the template configuration required for
	// post processing the aws connection.
	templateConfig, err := GetTemplateConfiguration(model, true, true)
	if err != nil {
		diagnostics := diag.Diagnostics{}
		summary := "Unable to form template configuration"
		detail := err.Error()
		diagnostics.AddError(summary, detail)
		return diagnostics
	}
	templateConfig["insights"] = templateConfig["discover"]
	delete(templateConfig, "discover")
	configBytes, err := json.Marshal(templateConfig)
	if err != nil {
		diagnostics := diag.Diagnostics{}
		summary := "Unable to marshal template configuration"
		detail := err.Error()
		diagnostics.AddError(summary, detail)
		return diagnostics
	}
	configuration := string(configBytes)

	// Call the Clumio API to post process aws connection.
	_, apiErr := r.sdkPostPsConnection.PostProcessAwsConnection(
		&models.PostProcessAwsConnectionV1Request{
			AccountNativeId:     model.AccountID.ValueStringPointer(),
			AwsRegion:           model.Region.ValueStringPointer(),
			Configuration:       &configuration,
			RequestType:         &eventType,
			RoleArn:             model.RoleArn.ValueStringPointer(),
			RoleExternalId:      model.RoleExternalID.ValueStringPointer(),
			Token:               model.Token.ValueStringPointer(),
			ClumioEventPubId:    model.ClumioEventPubID.ValueStringPointer(),
			Properties:          propertiesMap,
			IntermediateRoleArn: model.IntermediateRoleArn.ValueStringPointer(),
		})
	if apiErr != nil {
		diagnostics := diag.Diagnostics{}
		summary := "Error in invoking Post-process Clumio AWS Connection."
		detail := common.ParseMessageFromApiError(apiErr)
		diagnostics.AddError(summary, detail)
		return diagnostics
	}
	return nil
}
