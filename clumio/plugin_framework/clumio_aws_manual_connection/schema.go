// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_aws_manual_connection Terraform resource.

package clumio_aws_manual_connection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioAWSManualConnectionResourceModel is the resource model for the clumio_aws_manual_connection
// Terraform resource. It represents the schema of the resource and the data it holds. This schema
// is used by ustomers to configure the resource and by the Clumio provider to read and write the
// resource.
type clumioAWSManualConnectionResourceModel struct {
	ID            types.String        `tfsdk:"id"`
	AccountId     types.String        `tfsdk:"account_id"`
	AwsRegion     types.String        `tfsdk:"aws_region"`
	AssetsEnabled *AssetsEnabledModel `tfsdk:"assets_enabled"`
	Resources     *ResourcesModel     `tfsdk:"resources"`
}

// AssetsEnabledModel maps to the 'assets_enabled' field in clumioAWSManualConnectionResourceModel
// and is used to denote which asset types are enabled for the manual connection.
type AssetsEnabledModel struct {
	EBS      types.Bool `tfsdk:"ebs"`
	RDS      types.Bool `tfsdk:"rds"`
	DynamoDB types.Bool `tfsdk:"ddb"`
	S3       types.Bool `tfsdk:"s3"`
	EC2MSSQL types.Bool `tfsdk:"mssql"`
}

// ResourcesModel maps to the 'resources' field in clumioAWSManualConnectionResourceModel and is
// used to denote the stack ARNs to the configured manual resources for the connection.
type ResourcesModel struct {
	// IAM role with permissions to enable Clumio to backup and restore your assets
	ClumioIAMRoleArn types.String `tfsdk:"clumio_iam_role_arn"`
	// IAM role with permissions used by Clumio to create AWS support cases
	ClumioSupportRoleArn types.String `tfsdk:"clumio_support_role_arn"`
	// SNS topic to publish messages to Clumio services
	ClumioEventPubArn types.String `tfsdk:"clumio_event_pub_arn"`
	// Event rules for tracking changes in assets
	EventRules *EventRules `tfsdk:"event_rules"`
	// Asset-specific service roles
	ServiceRoles *ServiceRoles `tfsdk:"service_roles"`
}

// EventRules maps to 'event_rules' field in ResourcesModel and contains stack ARNs to the event
// rules used for tracking changes in assets.
type EventRules struct {
	// Event rule for tracking resource changes in selected assets
	CloudtrailRuleArn types.String `tfsdk:"cloudtrail_rule_arn"`
	// Event rule for tracking tag and resource changes in selected assets
	CloudwatchRuleArn types.String `tfsdk:"cloudwatch_rule_arn"`
}

// ServiceRoles maps to 'service_roles' field in ResourcesModel and contains stack ARNs to the
// asset-specific service roles for the connection.
type ServiceRoles struct {
	// Service roles required for mssql
	Mssql *MssqlServiceRoles `tfsdk:"mssql"`
	// Service roles required for s3
	S3 *S3ServiceRoles `tfsdk:"s3"`
}

// MssqlServiceRoles maps to 'mssql' field in ServiceRoles and contains stack ARNs to Mssql specific
// service roles
type MssqlServiceRoles struct {
	// Role assumable by ssm service
	SsmNotificationRoleArn types.String `tfsdk:"ssm_notification_role_arn"`
	// Instance created for ec2 instance profile role
	Ec2SsmInstanceProfileArn types.String `tfsdk:"ec2_ssm_instance_profile_arn"`
}

// S3ServiceRoles maps to 's3' field in ServiceRoles and contains stack ARNs to S3 specific
// service roles
type S3ServiceRoles struct {
	// Role assumed for continuous backups
	ContinuousBackupsRoleArn types.String `tfsdk:"continuous_backups_role_arn"`
}

// Schema defines the structure and constraints of the clumio_aws_manual_connection Terraform
// resource. Schema is a method on the clumioAWSManualConnectionResource struct. It sets the schema
// for the clumio_aws_manual_connection Terraform resource, which is used deploy resources for
// manual connections in Clumio. The schema defines various attributes such as the connection ID,
// AWS account ID, AWS region, assets enabled, etc. Some of these attributes are computed, meaning they
// are determined by Clumio at runtime, while others are required inputs from the user.
func (r *clumioAWSManualConnectionResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Clumio AWS Manual Connection Resource used to setup manual resources for" +
			" connections.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio AWS manual connection.",
				Computed:    true,
			},
			schemaAccountId: schema.StringAttribute{
				Description: "Identifier of the AWS account to be linked with Clumio.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "Region of the AWS account to be linked with Clumio.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAssetsEnabled: schema.ObjectAttribute{
				Description: "Assets enabled for the connection. Note that `mssql` is only " +
					"available for legacy connections.",
				Required: true,
				AttributeTypes: map[string]attr.Type{
					schemaIsEbsEnabled:      types.BoolType,
					schemaIsDynamoDBEnabled: types.BoolType,
					schemaIsRDSEnabled:      types.BoolType,
					schemaIsS3Enabled:       types.BoolType,
					schemaIsMssqlEnabled:    types.BoolType,
				},
			},
			schemaResources: schema.ObjectAttribute{
				Description: "An object containing the ARNs of the resources created for the manual AWS" +
					" connection. Please refer to this guide for instructions on how to create them. - " +
					"https://help.clumio.com/docs/manual-setup-for-aws-account-integration. If any" +
					" of the ARNs are not applicable to the manual connection, provide an empty" +
					" string \"\".",
				Required: true,
				AttributeTypes: map[string]attr.Type{
					schemaClumioIAMRoleArn:     types.StringType,
					schemaClumioEventPubArn:    types.StringType,
					schemaClumioSupportRoleArn: types.StringType,
					schemaEventRules: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							schemaCloudtrailRuleArn: types.StringType,
							schemaCloudwatchRuleArn: types.StringType,
						},
					},
					schemaServiceRoles: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							schemaS3: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									schemaContinuousBackupsRoleArn: types.StringType,
								},
							},
							schemaMssql: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									schemaSsmNotificationRoleArn:   types.StringType,
									schemaEc2SsmInstanceProfileArn: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}
