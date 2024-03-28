// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Post Process KMS SDK APIs to perform CRUD operations and
// set the attributes from the response of the API in the resource model.

package clumio_post_process_kms

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// clumioPostProcessKmsCommon contains the common logic for create, update and delete operations of
// PostProcessKms resource.
func (r *clumioPostProcessKmsResource) clumioPostProcessKmsCommon(
	_ context.Context, state clumioPostProcessKmsResourceModel, eventType string) diag.Diagnostics {

	var diags diag.Diagnostics

	templateVersion := uint64(*state.TemplateVersion.ValueInt64Pointer())

	// Call the Clumio API to post process kms.
	_, apiErr := r.sdkPostProcessKMS.PostProcessKms(
		&models.PostProcessKmsV1Request{
			AccountNativeId:       state.AccountId.ValueStringPointer(),
			AwsRegion:             state.Region.ValueStringPointer(),
			RequestType:           &eventType,
			Token:                 state.Token.ValueStringPointer(),
			MultiRegionCmkKeyId:   state.MultiRegionCMKKeyId.ValueStringPointer(),
			RoleId:                state.RoleId.ValueStringPointer(),
			RoleArn:               state.RoleArn.ValueStringPointer(),
			RoleExternalId:        state.RoleExternalId.ValueStringPointer(),
			CreatedMultiRegionCmk: state.CreatedMultiRegionCMK.ValueBoolPointer(),
			Version:               &templateVersion,
		})
	if apiErr != nil {
		summary := "Error in invoking Post-process Clumio KMS."
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
	}
	return diags
}
