// Copyright 2021. Clumio, Inc.

// clumio_post_process_aws_connection definition and CRUD implementation.
package clumio_post_process_aws_connection

import (
	"context"
	"encoding/json"
	"fmt"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/post_process_aws_connection"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioPostProcessAWSConnection does the post-processing for Clumio AWS Connection.
func ClumioPostProcessAWSConnection() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Post-Process Clumio AWS Connection Resource used to post-process AWS connection to Clumio.",

		CreateContext: clumioPostProcessAWSConnectionCreate,
		ReadContext:   clumioPostProcessAWSConnectionRead,
		UpdateContext: clumioPostProcessAWSConnectionUpdate,
		DeleteContext: clumioPostProcessAWSConnectionDelete,

		Schema: map[string]*schema.Schema{
			schemaToken: {
				Type:        schema.TypeString,
				Description: "The AWS integration ID token.",
				Required:    true,
			},
			schemaRoleExternalId: {
				Type: schema.TypeString,
				Description: "A key that must be used by Clumio to assume the service role" +
					" in your account. This should be a secure string, like a password," +
					" but it does not need to be remembered (random characters are best).",
				Required: true,
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
			schemaRoleArn: {
				Type:        schema.TypeString,
				Description: "Clumio IAM Role Arn.",
				Required:    true,
			},
			schemaConfigVersion: {
				Type:        schema.TypeString,
				Description: "Clumio Config version.",
				Required:    true,
			},
			schemaDiscoverVersion: {
				Type:        schema.TypeString,
				Description: "Clumio Discover version.",
				Required:    true,
			},
			schemaProtectConfigVersion: {
				Type:        schema.TypeString,
				Description: "Clumio Protect Config version.",
				Optional:    true,
			},
			schemaProtectEbsVersion: {
				Type:        schema.TypeString,
				Description: "Clumio EBS Protect version.",
				Optional:    true,
			},
			schemaProtectRdsVersion: {
				Type:        schema.TypeString,
				Description: "Clumio RDS Protect version.",
				Optional:    true,
			},
			schemaProtectEc2MssqlVersion: {
				Type:        schema.TypeString,
				Description: "Clumio EC2 MSSQL Protect version.",
				Optional:    true,
			},
			schemaProtectS3Version: {
				Type:        schema.TypeString,
				Description: "Clumio S3 Protect version.",
				Optional:    true,
			},
			schemaProtectDynamodbVersion: {
				Type:        schema.TypeString,
				Description: "Clumio DynamoDB Protect version.",
				Optional:    true,
			},
			schemaProtectWarmTierVersion: {
				Type:        schema.TypeString,
				Description: "Clumio Warm Tier Protect version.",
				Optional:    true,
			},
			schemaProtectWarmTierDynamodbVersion: {
				Type:        schema.TypeString,
				Description: "Clumio DynamoDB Warm Tier Protect version.",
				Optional:    true,
			},
			schemaClumioEventPubId: {
				Type:        schema.TypeString,
				Description: "Clumio Event Pub SNS topic ID.",
				Required:    true,
			},
			schemaProperties: {
				Type: schema.TypeMap,
				Description: "A map to pass in additional information to be consumed " +
					"by Clumio Post Processing",
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// clumioPostProcessAWSConnectionCreate handles the Create action for the
// PostProcessAWSConnection resource.
func clumioPostProcessAWSConnectionCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioPostProcessAWSConnectionCommon(ctx, d, meta, "Create")
}

// clumioPostProcessAWSConnectionRead handles the Create action for the
// PostProcessAWSConnection resource.
func clumioPostProcessAWSConnectionRead(
	_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

// clumioPostProcessAWSConnectionUpdate handles the Create action for the
// PostProcessAWSConnection resource.
func clumioPostProcessAWSConnectionUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioPostProcessAWSConnectionCommon(ctx, d, meta, "Update")
}

// clumioPostProcessAWSConnectionDelete handles the Create action for the
// PostProcessAWSConnection resource.
func clumioPostProcessAWSConnectionDelete(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioPostProcessAWSConnectionCommon(ctx, d, meta, "Delete")
}

// clumioPostProcessAWSConnectionCommon contains the common logic for all CRUD operations
// of PostProcessAWSConnection resource.
func clumioPostProcessAWSConnectionCommon(_ context.Context, d *schema.ResourceData,
	meta interface{}, eventType string) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	postProcessAwsConnection := aws_connections.NewPostProcessAwsConnectionV1(
		client.ClumioConfig)
	accountId := common.GetStringValue(d, schemaAccountId)
	awsRegion := common.GetStringValue(d, schemaRegion)
	roleArn := common.GetStringValue(d, schemaRoleArn)
	token := common.GetStringValue(d, schemaToken)
	roleExternalId := common.GetStringValue(d, schemaRoleExternalId)
	clumioEventPubId := common.GetStringValue(d, schemaClumioEventPubId)
	schemaPropertiesIface, ok := d.GetOk(schemaProperties)
	var propertiesMap map[string]*string
	if ok {
		schemaPropertiesMap := schemaPropertiesIface.(map[string]interface{})
		propertiesMap = make(map[string]*string)
		for key, val := range schemaPropertiesMap {
			valStr := val.(string)
			propertiesMap[key] = &valStr
		}
	}

	templateConfig, err := common.GetTemplateConfiguration(d, true)
	if err != nil {
		return diag.Errorf("Error forming template configuration. Error: %v", err)
	}
	templateConfig["insights"] = templateConfig["discover"]
	delete(templateConfig, "discover")
	configBytes, err := json.Marshal(templateConfig)
	if err != nil {
		return diag.Errorf("Error in marshalling template configuraton. Error: %v", err)
	}
	configuration := string(configBytes)
	_, apiErr := postProcessAwsConnection.PostProcessAwsConnection(
		&models.PostProcessAwsConnectionV1Request{
			AccountNativeId:  &accountId,
			AwsRegion:        &awsRegion,
			Configuration:    &configuration,
			RequestType:      &eventType,
			RoleArn:          &roleArn,
			RoleExternalId:   &roleExternalId,
			Token:            &token,
			ClumioEventPubId: &clumioEventPubId,
			Properties:       propertiesMap,
		})
	if apiErr != nil {
		return diag.Errorf(
			"Error in invoking Post-process Clumio AWS Connection. Error: %v",
			string(apiErr.Response))
	}
	if eventType == "Create" {
		d.SetId(fmt.Sprintf("%v/%v/%v", accountId, awsRegion, token))
	}
	return nil
}
