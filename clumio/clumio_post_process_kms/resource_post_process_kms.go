// Copyright 2022. Clumio, Inc.

// clumio_post_process_kms definition and CRUD implementation.
package clumio_post_process_kms

import (
	"context"
	"fmt"

	kms "github.com/clumio-code/clumio-go-sdk/controllers/post_process_kms"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	minTemplateVersion = 10
)

// ClumioPostProcessKMS does the post-processing for Clumio AWS Connection.
func ClumioPostProcessKMS() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Post-Process Clumio KMS Resource used to post-process KMS in Clumio.",

		CreateContext: clumioPostProcessKMSCreate,
		ReadContext:   clumioPostProcessKMSRead,
		UpdateContext: clumioPostProcessKMSUpdate,
		DeleteContext: clumioPostProcessKMSDelete,

		Schema: map[string]*schema.Schema{
			schemaToken: {
				Type:        schema.TypeString,
				Description: "The AWS integration ID token.",
				Required:    true,
			},
			schemaAccountId: {
				Type:        schema.TypeString,
				Description: "The AWS Customer Account ID.",
				Required:    true,
			},
			schemaRegion: {
				Type:        schema.TypeString,
				Description: "The AWS Region.",
				Required:    true,
			},
			schemaRoleId: {
				Type:        schema.TypeString,
				Description: "The ID of the IAM role to manage the CMK.",
				Required:    true,
			},
			schemaRoleArn: {
				Type:        schema.TypeString,
				Description: "The ARN of the IAM role to manage the CMK.",
				Required:    true,
			},
			schemaRoleExternalId: {
				Type:        schema.TypeString,
				Description: "The external ID to use when assuming the IAM role.",
				Required:    true,
			},
			schemaCreatedMultiRegionCMK: {
				Type:        schema.TypeBool,
				Description: "Whether a new CMK was created.",
				Required:    true,
			},
			schemaMultiRegionCMKKeyID: {
				Type:        schema.TypeString,
				Description: "Multi Region CMK Key ID.",
				Required:    true,
			},
			schemaTemplateVersion: {
				Type:        schema.TypeInt,
				Description: "Template version",
				Optional:    true,
			},
		},
	}
}

// clumioPostProcessKMSCreate handles the Create action for the PostProcessKMS resource.
func clumioPostProcessKMSCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioPostProcessKMSCommon(ctx, d, meta, "Create")
}

// clumioPostProcessKMSRead handles the Read action for the PostProcessKMS resource.
func clumioPostProcessKMSRead(
	_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

// clumioPostProcessKMSUpdate handles the Update action for the PostProcessKMS resource.
func clumioPostProcessKMSUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioPostProcessKMSCommon(ctx, d, meta, "Update")
}

// clumioPostProcessKMSDelete handles the Delete action for the PostProcessKMS resource.
func clumioPostProcessKMSDelete(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioPostProcessKMSCommon(ctx, d, meta, "Delete")
}

// clumioPostProcessKMSCommon contains the common logic for all CRUD operations
// of PostProcessKMS resource.
func clumioPostProcessKMSCommon(_ context.Context, d *schema.ResourceData,
	meta interface{}, eventType string) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	postProcessAwsKMS := kms.NewPostProcessKmsV1(
		client.ClumioConfig)
	accountId := common.GetStringValue(d, schemaAccountId)
	awsRegion := common.GetStringValue(d, schemaRegion)
	multiRegionCMKKeyID := common.GetStringValue(d, schemaMultiRegionCMKKeyID)
	token := common.GetStringValue(d, schemaToken)
	roleId := common.GetStringValue(d, schemaRoleId)
	roleArn := common.GetStringValue(d, schemaRoleArn)
	roleExternalId := common.GetStringValue(d, schemaRoleExternalId)
	createdMrCmk := common.GetBoolValueWithDefault(d, schemaCreatedMultiRegionCMK, true)
	version := uint64(common.GetIntValueWithDefault(d, schemaTemplateVersion, minTemplateVersion))

	_, apiErr := postProcessAwsKMS.PostProcessKms(
		&models.PostProcessKmsV1Request{
			AccountNativeId:       &accountId,
			AwsRegion:             &awsRegion,
			RequestType:           &eventType,
			Token:                 &token,
			MultiRegionCmkKeyId:   &multiRegionCMKKeyID,
			RoleId:                &roleId,
			RoleArn:               &roleArn,
			RoleExternalId:        &roleExternalId,
			CreatedMultiRegionCmk: &createdMrCmk,
			Version:               &version,
		})
	if apiErr != nil {
		return diag.Errorf(
			"Error in invoking Post-process KMS. Error: %v",
			string(apiErr.Response))
	}
	if eventType == "Create" {
		d.SetId(fmt.Sprintf("%v/%v/%v", accountId, awsRegion, token))
	}
	return nil
}
