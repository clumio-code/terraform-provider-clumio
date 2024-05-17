// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_post_process_aws_connection Terraform resource.

package clumio_post_process_aws_connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// postProcessAWSConnectionResourceModel is the resource model for the
// clumio_post_process_aws_connection Terraform resource. It represents the schema of the resource
// and the data it holds. This schema is used by customers to configure the resource and by the
// Clumio provider to read and write the resource.
type postProcessAWSConnectionResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	AccountID                      types.String `tfsdk:"account_id"`
	Token                          types.String `tfsdk:"token"`
	RoleExternalID                 types.String `tfsdk:"role_external_id"`
	Region                         types.String `tfsdk:"region"`
	ClumioEventPubID               types.String `tfsdk:"clumio_event_pub_id"`
	RoleArn                        types.String `tfsdk:"role_arn"`
	ConfigVersion                  types.String `tfsdk:"config_version"`
	DiscoverVersion                types.String `tfsdk:"discover_version"`
	ProtectConfigVersion           types.String `tfsdk:"protect_config_version"`
	ProtectEBSVersion              types.String `tfsdk:"protect_ebs_version"`
	ProtectRDSVersion              types.String `tfsdk:"protect_rds_version"`
	ProtectS3Version               types.String `tfsdk:"protect_s3_version"`
	ProtectDynamoDBVersion         types.String `tfsdk:"protect_dynamodb_version"`
	ProtectEC2MssqlVersion         types.String `tfsdk:"protect_ec2_mssql_version"`
	ProtectWarmTierVersion         types.String `tfsdk:"protect_warm_tier_version"`
	ProtectWarmTierDynamoDBVersion types.String `tfsdk:"protect_warm_tier_dynamodb_version"`
	Properties                     types.Map    `tfsdk:"properties"`
	IntermediateRoleArn            types.String `tfsdk:"intermediate_role_arn"`
	WaitForIngestion               types.Bool   `tfsdk:"wait_for_ingestion"`
	WaitForDataPlaneResources      types.Bool   `tfsdk:"wait_for_data_plane_resources"`
}

// Schema defines the structure and constraints of the clumio_post_process_aws_connection Terraform
// resource. Schema is a method on the postProcessAWSConnectionResource struct. It sets the schema
// for the clumio_post_process_aws_connection Terraform resource.
func (r *postProcessAWSConnectionResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Post-Process Clumio AWS Connection Resource used to" +
			" post-process AWS connection to Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "The unique identifier of the post process aws connection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaToken: schema.StringAttribute{
				Description: "Distinct 36-character token used to identify resources set up by " +
					"the Clumio AWS template installation on the account being connected.",
				Required: true,
			},
			schemaRoleExternalId: schema.StringAttribute{
				Description: "Unique identifier Clumio uses to access the service role within " +
					"your account.",
				Required: true,
			},
			schemaAccountId: schema.StringAttribute{
				Description: "Identifier of the AWS account to be linked with Clumio.",
				Required:    true,
			},
			schemaRegion: schema.StringAttribute{
				Description: "Region of the AWS account to be linked with Clumio.",
				Required:    true,
			},
			schemaRoleArn: schema.StringAttribute{
				Description: "ARN of the role which allows Clumio to access the linked account.",
				Required:    true,
			},
			schemaConfigVersion: schema.StringAttribute{
				Description: "Clumio Config version.",
				Required:    true,
			},
			schemaDiscoverVersion: schema.StringAttribute{
				Description: "Clumio Discover version.",
				Optional:    true,
			},
			schemaProtectConfigVersion: schema.StringAttribute{
				Description: "Clumio Protect Config version.",
				Optional:    true,
			},
			schemaProtectEbsVersion: schema.StringAttribute{
				Description: "Clumio EBS Protect version.",
				Optional:    true,
			},
			schemaProtectRdsVersion: schema.StringAttribute{
				Description: "Clumio RDS Protect version.",
				Optional:    true,
			},
			schemaProtectEc2MssqlVersion: schema.StringAttribute{
				Description: "Clumio EC2 MSSQL Protect version.",
				Optional:    true,
			},
			schemaProtectS3Version: schema.StringAttribute{
				Description: "Clumio S3 Protect version.",
				Optional:    true,
			},
			schemaProtectDynamodbVersion: schema.StringAttribute{
				Description: "Clumio DynamoDB Protect version.",
				Optional:    true,
			},
			schemaProtectWarmTierVersion: schema.StringAttribute{
				Description: "Clumio Warm Tier Protect version.",
				Optional:    true,
			},
			schemaProtectWarmTierDynamodbVersion: schema.StringAttribute{
				Description: "Clumio DynamoDB Warm Tier Protect version.",
				Optional:    true,
			},
			schemaClumioEventPubId: schema.StringAttribute{
				Description: "Clumio Event Pub SNS topic ID.",
				Required:    true,
			},
			schemaProperties: schema.MapAttribute{
				Description: "A map to pass in additional information to be consumed " +
					"by Clumio Post Processing",
				Optional:    true,
				ElementType: types.StringType,
			},
			schemaIntermediateRoleArn: schema.StringAttribute{
				Description: "Intermediate Role arn to be assumed before accessing" +
					" ClumioRole in customer account.",
				Optional: true,
			},
			schemaWaitForIngestion: schema.BoolAttribute{
				Description: "Wait for the AWS connection ingestion task to complete.",
				Optional:    true,
			},
			schemaWaitForDataPlaneResources: schema.BoolAttribute{
				Description: "Wait for the data plane resources to be created.",
				Optional:    true,
			},
		},
	}
}
