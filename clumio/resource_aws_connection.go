// Copyright 2021. Clumio, Inc.

// clumio_aws_connection definition and CRUD implementation.
package clumio

import (
	"context"
	"strings"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/aws_connections"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// clumioAWSConnection returns the resource for Clumio AWS Connection.

func clumioAWSConnection() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio AWS Connection Resource used to connect AWS accounts to Clumio.",

		CreateContext: clumioAWSConnectionCreate,
		ReadContext:   clumioAWSConnectionRead,
		UpdateContext: clumioAWSConnectionUpdate,
		DeleteContext: clumioAWSConnectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Clumio AWS Connection Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"account_native_id": {
				Description: "AWS Account Id to connect to Clumio.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"aws_region": {
				Description: "AWS Region of account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Clumio AWS Connection Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"organizational_unit_id": {
				Description: "Clumio Organizational Unit Id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"protect_asset_types_enabled": {
				Description: "The asset types enabled for protect. This is only" +
					" populated if protect is enabled. Valid values are any of" +
					" [EBS, RDS, DynamoDB, EC2MSSQL, S3].",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"services_enabled": {
				Description: "The services to be enabled for this configuration." +
					" Valid values are [discover], [discover, protect]. This is only set" +
					" when the registration is created, the enabled services are" +
					" obtained directly from the installed template after that.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"connection_status": {
				Description: "The status of the connection. Possible values include " +
					"connecting, connected and unlinked.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": {
				Description: "The 36-character Clumio AWS integration ID token.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"namespace": {
				Description: "K8S Namespace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"clumio_aws_account_id": {
				Description: "Clumio AWS AccountId.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"clumio_aws_region": {
				Description: "Clumio AWS Region.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func clumioAWSConnectionCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.clumioConfig)
	accountNativeId := getStringValue(d, "account_native_id")
	awsRegion := getStringValue(d, "aws_region")
	description := getStringValue(d, "description")
	organizationalUnitId := getStringValue(d, "organizational_unit_id")
	res, apiErr := awsConnection.CreateAwsConnection(&models.CreateAwsConnectionV1Request{
		AccountNativeId:          &accountNativeId,
		AwsRegion:                &awsRegion,
		Description:              &description,
		OrganizationalUnitId:     &organizationalUnitId,
		ProtectAssetTypesEnabled: getStringSlice(d, "protect_asset_types_enabled"),
		ServicesEnabled:          getStringSlice(d, "services_enabled"),
	})
	if apiErr != nil {
		return diag.Errorf(
			"Error creating Clumio AWS Connection. Error: %v", string(apiErr.Response))
	}
	d.SetId(*res.Id)
	return clumioAWSConnectionRead(ctx, d, meta)
}

func clumioAWSConnectionRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.clumioConfig)
	res, apiErr := awsConnection.ReadAwsConnection(d.Id())
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			d.SetId("")
			return nil
		}
		return diag.Errorf(
			"Error creating Clumio AWS Connection. Error: %v", string(apiErr.Response))

	}
	err := d.Set("token", *res.Token)
	if err != nil {
		return diag.Errorf(
			"Error setting token schema attribute. Error: %v", err)
	}
	err = d.Set("namespace", res.Namespace)
	if err != nil {
		return diag.Errorf(
			"Error setting namespace schema attribute. Error: %v", err)
	}
	err = d.Set("clumio_aws_account_id", res.ClumioAwsAccountId)
	if err != nil {
		return diag.Errorf("Error setting clumio_aws_account_id schema attribute."+
			" Error: %v", err)
	}
	err = d.Set("clumio_aws_region", res.ClumioAwsRegion)
	if err != nil {
		return diag.Errorf("Error setting clumio_aws_region schema attribute."+
			" Error: %v", err)
	}
	return nil
}

func clumioAWSConnectionUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("description") {
		return nil
	}
	client := meta.(*apiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.clumioConfig)
	description := getStringValue(d, "description")
	_, apiErr := awsConnection.UpdateAwsConnection(d.Id(),
		models.UpdateAwsConnectionV1Request{
			Description: &description,
		})
	if apiErr != nil {
		return diag.Errorf(
			"Error updating description of Clumio AWS Connection %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return clumioAWSConnectionRead(ctx, d, meta)
}

func clumioAWSConnectionDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.clumioConfig)
	_, apiErr := awsConnection.DeleteAwsConnection(d.Id())
	if apiErr != nil {
		return diag.Errorf(
			"Error deleting Clumio AWS Connection %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return nil
}
