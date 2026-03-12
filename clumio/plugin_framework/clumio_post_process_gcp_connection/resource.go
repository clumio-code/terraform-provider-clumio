// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_post_process_gcp_connection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var jsonMarshal = json.Marshal

// createUpdatePostProcessGcpConnection invokes the API to create/update the connection and from the response populates the
// computed attributes of the connection.
func (r *clumioPostProcessGCPConnectionResource) createUpdatePostProcessGcpConnection(_ context.Context, model *clumioPostProcessGCPConnectionResourceModel, requestType string) diag.Diagnostics {
	var diags diag.Diagnostics

	schemaPropertiesElements := model.Properties.Elements()
	propertiesMap := make(map[string]*string)
	for key, val := range schemaPropertiesElements {
		valStr := val.(types.String).ValueString()
		propertiesMap[key] = &valStr
	}

	templateConfig, err := GetTemplateConfiguration(model)
	if err != nil {
		summary := "Error in invoking Post-process Clumio GCP Connection."
		detail := "Unable to create template configurations from versions: " + err.Error()
		diags.AddError(summary, detail)
		return diags
	}

	configBytes, err := jsonMarshal(templateConfig)
	if err != nil {
		summary := "Unable to marshal template configuration"
		detail := err.Error()
		diags.AddError(summary, detail)
		return diags
	}
	configuration := string(configBytes)

	postprocessRequest := &models.PostProcessGcpConnectionV1Request{
		Configuration:       &configuration,
		ProjectId:           model.ProjectID.ValueStringPointer(),
		ProjectName:         model.ProjectName.ValueStringPointer(),
		ProjectNumber:       model.ProjectNumber.ValueStringPointer(),
		RequestType:         &requestType,
		ResourceProperties:  propertiesMap,
		ServiceAccountEmail: model.ServiceAccountEmail.ValueStringPointer(),
		Token:               model.Token.ValueStringPointer(),
		WifPoolId:           model.WifPoolId.ValueStringPointer(),
		WifProviderId:       model.WifProviderId.ValueStringPointer(),
	}

	_, apiErr := r.sdkConnections.PostProcessGcpConnection(postprocessRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to process %s Post-process Clumio GCP Connection. %s (project id: %v)",
			requestType, r.name, model.ProjectID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	// ID needs to be a value which is used by our backend to uniquely identify connection
	model.ID = types.StringPointerValue(model.Token.ValueStringPointer())
	return diags
}

// deletePostProcessGcpConnection invokes the API to delete the connection
func (r *clumioPostProcessGCPConnectionResource) deletePostProcessGcpConnection(_ context.Context, model *clumioPostProcessGCPConnectionResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	postprocessRequest := &models.PostProcessGcpConnectionV1Request{
		ProjectId:     model.ProjectID.ValueStringPointer(),
		ProjectName:   model.ProjectName.ValueStringPointer(),
		ProjectNumber: model.ProjectNumber.ValueStringPointer(),
		RequestType:   types.StringValue(deleteRequestType).ValueStringPointer(),
		Token:         model.Token.ValueStringPointer(),
	}

	_, apiErr := r.sdkConnections.PostProcessGcpConnection(postprocessRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to process %s Post-process Clumio GCP Connection. %s (project id: %v)",
			deleteRequestType, r.name, model.ProjectID.ValueString())
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}

	return diags
}
