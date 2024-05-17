// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Post Process AWS Connection SDK APIs to perform CRUD
// operations and set the attributes from the response of the API in the resource model.

package clumio_post_process_aws_connection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioPostProcessAWSConnectionCommon contains the common logic for create, update and delete
// operations of PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) clumioPostProcessAWSConnectionCommon(ctx context.Context,
	model postProcessAWSConnectionResourceModel, eventType string) diag.Diagnostics {

	var diags diag.Diagnostics

	schemaPropertiesElements := model.Properties.Elements()
	propertiesMap := make(map[string]*string)
	for key, val := range schemaPropertiesElements {
		valStr := val.(types.String).ValueString()
		propertiesMap[key] = &valStr
	}

	// Using the schema properties in the model, create the template configuration required for
	// post processing the aws connection.
	templateConfig, err := GetTemplateConfiguration(model, true)
	if err != nil {
		summary := "Unable to form template configuration"
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}
	templateConfig["insights"] = templateConfig["discover"]
	delete(templateConfig, "discover")
	configBytes, err := json.Marshal(templateConfig)
	if err != nil {
		summary := "Unable to marshal template configuration"
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}
	configuration := string(configBytes)

	if eventType == eventTypeCreate {
		readExtId := "false"
		connId := fmt.Sprintf("%s_%s", model.AccountID.ValueString(), model.Region.ValueString())
		connRes, apiErr := r.sdkAWSConnection.ReadAwsConnection(connId, &readExtId)
		if apiErr != nil {
			summary := fmt.Sprintf("Unable to read Clumio AWS connection with id: %s", connId)
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return diags
		}
		if connRes != nil && *connRes.ConnectionStatus == connected {
			eventType = eventTypeUpdate
		}
	}
	// Call the Clumio API to post process aws connection.
	req := &models.PostProcessAwsConnectionV1Request{
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
	}
	_, apiErr := r.sdkPostProcessConn.PostProcessAwsConnection(req)
	if apiErr != nil {
		summary := "Error in invoking Post-process Clumio AWS Connection."
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if (eventType == eventTypeCreate || eventType == eventTypeUpdate) &&
		(model.WaitForIngestion.ValueBool() || model.WaitForDataPlaneResources.ValueBool()) {

		targetSetupErr, err := pollForConnectionIngestionAndTargetStatus(
			ctx, r.sdkAWSConnection, model, r.pollTimeout, r.pollInterval)
		if err != nil {
			if targetSetupErr {
				summary := "Error in polling for connection ingestion and/or data plane resources" +
					" setup status."
				Detail := err.Error()
				diags.AddError(summary, Detail)
			} else {
				summary := "Error in polling for ingestion status."
				Detail := err.Error()
				diags.AddWarning(summary, Detail)
			}
		}
	}
	return diags
}
