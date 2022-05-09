// Copyright 2021. Clumio, Inc.

// clumio_aws_connection definition and CRUD implementation.

package clumio_aws_connection

import (
	"context"
	"strings"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/aws_connections"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioAWSConnection returns the resource for Clumio AWS Connection.
func ClumioAWSConnection() *schema.Resource {
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
			schemaId: {
				Description: "Clumio AWS Connection Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaAccountNativeId: {
				Description: "AWS Account Id to connect to Clumio.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			schemaAwsRegion: {
				Description: "AWS Region of account.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			schemaDescription: {
				Description: "Clumio AWS Connection Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			schemaOrganizationalUnitId: {
				Description: "Clumio Organizational Unit Id.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			schemaProtectAssetTypesEnabled: {
				Description: "The asset types enabled for protect. This is only" +
					" populated if protect is enabled. Valid values are any of" +
					" [EBS, RDS, DynamoDB, EC2MSSQL, S3].",
				Type: schema.TypeSet,
				Set:  common.SchemaSetHashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Deprecated: "This is no longer required as the asset types to be" +
					"enabled are based on the variables passed to the " +
					"clumio_terraform_aws_template module.",
			},
			schemaServicesEnabled: {
				Description: "The services to be enabled for this configuration." +
					" Valid values are [discover], [discover, protect]. This is only set" +
					" when the registration is created, the enabled services are" +
					" obtained directly from the installed template after that.",
				Type: schema.TypeSet,
				Set:  common.SchemaSetHashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Deprecated: "This is no longer required as by default discover and" +
					" protect are enabled.",
			},
			schemaConnectionStatus: {
				Description: "The status of the connection. Possible values include " +
					"connecting, connected and unlinked.",
				Type:     schema.TypeString,
				Computed: true,
			},
			schemaToken: {
				Description: "The 36-character Clumio AWS integration ID token used to" +
					" identify the installation of the Terraform template on the account.",
				Type:     schema.TypeString,
				Computed: true,
			},
			schemaNamespace: {
				Description: "K8S Namespace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaClumioAwsAccountId: {
				Description: "Clumio AWS AccountId.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaClumioAwsRegion: {
				Description: "Clumio AWS Region.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func clumioAWSConnectionCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.ClumioConfig)
	accountNativeId := common.GetStringValue(d, schemaAccountNativeId)
	awsRegion := common.GetStringValue(d, schemaAwsRegion)
	description := common.GetStringValue(d, schemaDescription)
	organizationalUnitId := common.GetStringValue(d, schemaOrganizationalUnitId)
	res, apiErr := awsConnection.CreateAwsConnection(&models.CreateAwsConnectionV1Request{
		AccountNativeId:      &accountNativeId,
		AwsRegion:            &awsRegion,
		Description:          &description,
		OrganizationalUnitId: &organizationalUnitId,
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
	client := meta.(*common.ApiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.ClumioConfig)
	res, apiErr := awsConnection.ReadAwsConnection(d.Id())
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			d.SetId("")
			return nil
		}
		return diag.Errorf(
			"Error retrieving Clumio AWS Connection. Error: %v", string(apiErr.Response))

	}
	err := d.Set(schemaConnectionStatus, *res.ConnectionStatus)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaConnectionStatus, err)
	}
	err = d.Set(schemaToken, *res.Token)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaToken, err)
	}
	if res.Namespace != nil {
		err = d.Set(schemaNamespace, *res.Namespace)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaNamespace, err)
		}
	}
	err = d.Set(schemaClumioAwsAccountId, *res.ClumioAwsAccountId)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaClumioAwsAccountId, err)
	}
	err = d.Set(schemaClumioAwsRegion, *res.ClumioAwsRegion)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaClumioAwsRegion, err)
	}
	err = d.Set(schemaAccountNativeId, *res.AccountNativeId)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaAccountNativeId, err)
	}
	err = d.Set(schemaAwsRegion, *res.AwsRegion)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaAwsRegion, err)
	}
	if res.Description != nil {
		err = d.Set(schemaDescription, *res.Description)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaDescription, err)
		}
	}
	if res.OrganizationalUnitId != nil {
		err = d.Set(schemaOrganizationalUnitId, *res.OrganizationalUnitId)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaOrganizationalUnitId, err)
		}
	}
	return nil
}

func clumioAWSConnectionUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange(schemaDescription) {
		return nil
	}
	client := meta.(*common.ApiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.ClumioConfig)
	description := common.GetStringValue(d, schemaDescription)
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
	client := meta.(*common.ApiClient)
	awsConnection := aws_connections.NewAwsConnectionsV1(client.ClumioConfig)
	_, apiErr := awsConnection.DeleteAwsConnection(d.Id())
	if apiErr != nil {
		return diag.Errorf(
			"Error deleting Clumio AWS Connection %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return nil
}
