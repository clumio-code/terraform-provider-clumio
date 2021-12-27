// Copyright 2021. Clumio, Inc.

// clumio_post_process_aws_connection definition and CRUD implementation.
package clumio

import (
	"context"
	"encoding/json"
	"fmt"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/post_process_aws_connection"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// clumioPostProcessAWSConnection does the post-processing for Clumio AWS Connection.
func clumioPostProcessAWSConnection() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Post process Clumio AWS Connection Resource used to post-process AWS connection to Clumio.",

		CreateContext: clumioPostProcessAWSConnectionCreate,
		ReadContext:   clumioPostProcessAWSConnectionRead,
		UpdateContext: clumioPostProcessAWSConnectionUpdate,
		DeleteContext: clumioPostProcessAWSConnectionDelete,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Description: "The AWS integration ID token.",
				Required:    true,
			},
			"role_external_id": {
				Type: schema.TypeString,
				Description: "A key that must be used by Clumio to assume the service role" +
					" in your account. This should be a secure string, like a password," +
					" but it does not need to be remembered (random characters are best).",
				Required: true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The AWS Customer Account ID.",
				Required:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The AWS Region.",
				Required:    true,
			},
			"role_arn": {
				Type:        schema.TypeString,
				Description: "Clumio IAM Role Arn.",
				Required:    true,
			},
			"config_version": {
				Type:        schema.TypeString,
				Description: "Clumio Config version.",
				Required:    true,
			},
			"discover_version": {
				Type:        schema.TypeString,
				Description: "Clumio Discover version.",
				Required:    true,
			},
			"protect_config_version": {
				Type:        schema.TypeString,
				Description: "Clumio Protect Config version.",
				Optional:    true,
			},
			"protect_ebs_version": {
				Type:        schema.TypeString,
				Description: "Clumio EBS Protect version.",
				Optional:    true,
			},
			"protect_rds_version": {
				Type:        schema.TypeString,
				Description: "Clumio RDS Protect version.",
				Optional:    true,
			},
			"protect_ec2_mssql_version": {
				Type:        schema.TypeString,
				Description: "Clumio EC2 MSSQL Protect version.",
				Optional:    true,
			},
			"protect_s3_version": {
				Type:        schema.TypeString,
				Description: "Clumio S3 Protect version.",
				Optional:    true,
			},
			"protect_dynamodb_version": {
				Type:        schema.TypeString,
				Description: "Clumio DynamoDB Protect version.",
				Optional:    true,
			},
			"protect_warm_tier_version": {
				Type:        schema.TypeString,
				Description: "Clumio Warmtier Protect version.",
				Optional:    true,
			},
			"protect_warm_tier_dynamodb_version": {
				Type:        schema.TypeString,
				Description: "Clumio DynamoDB Warmtier Protect version.",
				Optional:    true,
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
	client := meta.(*apiClient)
	postProcessAwsConnection := aws_connections.NewPostProcessAwsConnectionV1(
		client.clumioConfig)
	accountId := getStringValue(d, "account_id")
	awsRegion := getStringValue(d, "region")
	roleArn := getStringValue(d, "role_arn")
	token := getStringValue(d, "token")
	roleExternalId := getStringValue(d, "role_external_id")

	templateConfig, err := getTemplateConfiguration(d, true)
	if err != nil{
		return diag.Errorf("Error forming template configuration. Error: %v", err)
	}
	templateConfig["insights"] = templateConfig["discover"]
	delete(templateConfig, "discover")
	configBytes, err := json.Marshal(templateConfig)
	if err != nil{
		return diag.Errorf("Error in marshalling template configuraton. Error: %v", err)
	}
	configuration := string(configBytes)
	_, apiErr := postProcessAwsConnection.PostProcessAwsConnection(
		&models.PostProcessAwsConnectionV1Request{
			AccountNativeId: &accountId,
			AwsRegion:       &awsRegion,
			Configuration:   &configuration,
			RequestType:     &eventType,
			RoleArn:         &roleArn,
			RoleExternalId:  &roleExternalId,
			Token:           &token,
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