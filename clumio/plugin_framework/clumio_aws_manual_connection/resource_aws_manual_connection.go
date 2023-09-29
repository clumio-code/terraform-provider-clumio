package clumio_aws_manual_connection

import (
	"context"
	"fmt"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/aws_connections"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EBS = "EBS"
	S3 = "S3"
	DynamoDB = "DynamoDB"
	RDS = "RDS"
	EC2MSSQL = "EC2MSSQL" 
)

// AwsManualConnection model
type AwsManualConnectionModel struct {
	ID            types.String        `tfsdk:"id"`
	AccountId     types.String        `tfsdk:"account_id"`
	AwsRegion     types.String        `tfsdk:"aws_region"`
	AssetsEnabled *AssetsEnabledModel `tfsdk:"assets_enabled"`
	Resources     *ResourcesModel     `tfsdk:"resources"`
}

type AssetsEnabledModel struct {
	EBS      types.Bool `tfsdk:"ebs"`
	RDS      types.Bool `tfsdk:"rds"`
	DynamoDB types.Bool `tfsdk:"ddb"`
	S3       types.Bool `tfsdk:"s3"`
	EC2MSSQL types.Bool `tfsdk:"mssql"`
}

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

type EventRules struct {
	// Event rule for tracking resource changes in selected assets
	CloudtrailRuleArn types.String `tfsdk:"cloudtrail_rule_arn"`
	// Event rule for tracking tag and resource changes in selected assets
	CloudwatchRuleArn types.String `tfsdk:"cloudwatch_rule_arn"`
}

type ServiceRoles struct {
	// Service roles required for mssql
	Mssql *MssqlServiceRoles `tfsdk:"mssql"`
	// Service roles required for s3
	S3 *S3ServiceRoles `tfsdk:"s3"`
}

type MssqlServiceRoles struct {
	// Role assumable by ssm service
	SsmNotificationRoleArn types.String `tfsdk:"ssm_notification_role_arn"`
	// Instance created for ec2 instance profile role
	Ec2SsmInstanceProfileArn types.String `tfsdk:"ec2_ssm_instance_profile_arn"`
}

type S3ServiceRoles struct {
	// Role assumed for continuous backups
	ContinuousBackupsRoleArn types.String `tfsdk:"continuous_backups_role_arn"`
}

// awsManualConnectionResource is the resource implementation.
type awsManualConnectionResource struct {
	client *common.ApiClient
}

// NewAwsManualConnectionResource is a helper function to simplify the provider implementation.
func NewAwsManualConnectionResource() resource.Resource {
	return &awsManualConnectionResource{}
}

// Schema defines the schema for the resource.
func (r *awsManualConnectionResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Clumio AWS Manual Connection Resource used to setup manual resources for connections.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Clumio AWS Connection Id",
				Computed:    true,
			},
			schemaAccountId: schema.StringAttribute{
				Description: "AWS Account Id of the connection",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "AWS Region of the connection",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAssetsEnabled: schema.ObjectAttribute{
				Description: "Assets enabled for the connection",
				Required:    true,
				AttributeTypes: map[string]attr.Type{
					schemaIsEbsEnabled: types.BoolType,
					schemaIsDynamoDBEnabled: types.BoolType,
					schemaIsRDSEnabled: types.BoolType,
					schemaIsS3Enabled: types.BoolType,
					schemaIsMssqlEnabled: types.BoolType,
				},
			},
			schemaResources: schema.ObjectAttribute{
				Description: "Manual resources for the connection",
				Required: true,
				AttributeTypes: map[string]attr.Type{
					schemaClumioIAMRoleArn: types.StringType,
					schemaClumioEventPubArn: types.StringType,
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
									schemaSsmNotificationRoleArn: types.StringType,
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

// Create implements resource.Resource.
func (r *awsManualConnectionResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan AwsManualConnectionModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	diags = r.clumioSetManualResourcesCommon(ctx, plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	accountId := plan.AccountId.ValueString()
	awsRegion := plan.AwsRegion.ValueString()
	
	plan.ID = types.StringValue(fmt.Sprintf("%v_%v", accountId, awsRegion))

	diags = res.State.Set(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

}

// Delete implements resource.Resource.
func (*awsManualConnectionResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

// Read implements resource.Resource.
func (*awsManualConnectionResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
}

// Update implements resource.Resource.
func (r *awsManualConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan AwsManualConnectionModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	var state AwsManualConnectionModel
	diags = req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	changeInState := isAssetConfigChanged(&plan, &state) || isResourceConfigChanged(&plan, &state)
	if changeInState {
		diags = r.clumioSetManualResourcesCommon(ctx, plan)
		res.Diagnostics.Append(diags...)
		if res.Diagnostics.HasError() {
			return
		}
		plan.ID = types.StringValue(state.ID.ValueString())
		diags = res.State.Set(ctx, &plan)

		res.Diagnostics.Append(diags...)
		if res.Diagnostics.HasError() {
			return
		}
	}
}

// Metadata returns the resource type name.
func (r *awsManualConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_manual_connection"
}

// Configure adds the provider configured client to the data source.
func (r *awsManualConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// clumioSetManualResourcesCommon contains the logic for updating resources of a manual connection.
func (r *awsManualConnectionResource) clumioSetManualResourcesCommon(
	ctx context.Context, state AwsManualConnectionModel) diag.Diagnostics {
	awsConnection := aws_connections.NewAwsConnectionsV1(r.client.ClumioConfig)
	accountId := state.AccountId.ValueString()
	awsRegion := state.AwsRegion.ValueString()

	assetsEnabled := []*string{}
	if state.AssetsEnabled.EBS.ValueBool() {
		enabled := EBS
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.S3.ValueBool() {
		enabled := S3
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.RDS.ValueBool() {
		enabled := RDS
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.DynamoDB.ValueBool() {
		enabled := DynamoDB
		assetsEnabled = append(assetsEnabled, &enabled)
	}
	if state.AssetsEnabled.EC2MSSQL.ValueBool() {
		enabled := EC2MSSQL
		assetsEnabled = append(assetsEnabled, &enabled)
	}

	connectionId := accountId + "_" + awsRegion

	clumioIamRoleArn := state.Resources.ClumioIAMRoleArn.ValueString()
	clumioEventPubArn := state.Resources.ClumioEventPubArn.ValueString()
	clumioSupportRoleArn := state.Resources.ClumioSupportRoleArn.ValueString()
	cloudtrailRuleArn := state.Resources.EventRules.CloudtrailRuleArn.ValueString()
	cloudwatchRuleArn := state.Resources.EventRules.CloudwatchRuleArn.ValueString()
	continuousBackupsRoleArn := state.Resources.ServiceRoles.S3.ContinuousBackupsRoleArn.ValueString()
	ec2SsmInstanceProfileArn := state.Resources.ServiceRoles.Mssql.Ec2SsmInstanceProfileArn.ValueString()
	ssmNotificationRoleArn := state.Resources.ServiceRoles.Mssql.SsmNotificationRoleArn.ValueString()

	_, apiErr := awsConnection.UpdateAwsConnection(
		connectionId,
		models.UpdateAwsConnectionV1Request{
			AssetTypesEnabled: assetsEnabled,
			Resources: &models.Resources{
				ClumioIamRoleArn: &clumioIamRoleArn,
				ClumioEventPubArn: &clumioEventPubArn,
				ClumioSupportRoleArn: &clumioSupportRoleArn,
				EventRules: &models.EventRules{
					CloudtrailRuleArn: &cloudtrailRuleArn,
					CloudwatchRuleArn: &cloudwatchRuleArn,
				},
				ServiceRoles: &models.ServiceRoles{
					S3: &models.S3ServiceRoles{
						ContinuousBackupsRoleArn: &continuousBackupsRoleArn,
					},
					Mssql: &models.MssqlServiceRoles{
						Ec2SsmInstanceProfileArn: &ec2SsmInstanceProfileArn,
						SsmNotificationRoleArn: &ssmNotificationRoleArn,
					},
				},
			},
		},
	)
	if apiErr != nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Error in updating resources of Clumio AWS Manual Connection.", string(apiErr.Response))
		return diagnostics
	}
	return nil
}

// Check if any new assets were added (downgrades aren't supported)
func isAssetConfigChanged(plan *AwsManualConnectionModel, state *AwsManualConnectionModel) bool {
	// If EBS was enabled now
	if plan.AssetsEnabled.EBS.ValueBool() && !state.AssetsEnabled.EBS.ValueBool() {
		return true
	}
	// If S3 was enabled now
	if plan.AssetsEnabled.S3.ValueBool() && !state.AssetsEnabled.S3.ValueBool() {
		return true
	}
	// If RDS was enabled now
	if plan.AssetsEnabled.RDS.ValueBool() && !state.AssetsEnabled.RDS.ValueBool() {
		return true
	}
	// If DynamoDB was enabled now
	if plan.AssetsEnabled.DynamoDB.ValueBool() && !state.AssetsEnabled.DynamoDB.ValueBool() {
		return true
	}
	// If EC2MSSQL was enabled now
	if plan.AssetsEnabled.EC2MSSQL.ValueBool() && !state.AssetsEnabled.EC2MSSQL.ValueBool() {
		return true
	}
	return false
}

func isResourceConfigChanged(plan *AwsManualConnectionModel, state *AwsManualConnectionModel) bool {
	// Change in General Roles
	if plan.Resources.ClumioIAMRoleArn.ValueString() != state.Resources.ClumioIAMRoleArn.ValueString() {
		return true
	}
	if plan.Resources.ClumioSupportRoleArn.ValueString() != state.Resources.ClumioSupportRoleArn.ValueString() {
		return true
	}
	if plan.Resources.ClumioSupportRoleArn.ValueString() != state.Resources.ClumioSupportRoleArn.ValueString() {
		return true
	}
	// Change in Event Rules
	if plan.Resources.EventRules.CloudtrailRuleArn.ValueString() != 
		state.Resources.EventRules.CloudtrailRuleArn.ValueString() {
		return true
	}
	if plan.Resources.EventRules.CloudwatchRuleArn.ValueString() != 
		state.Resources.EventRules.CloudwatchRuleArn.ValueString() {
		return true
	}
	// Change in Service Roles
	if plan.Resources.ServiceRoles.S3.ContinuousBackupsRoleArn.ValueString() != 
		state.Resources.ServiceRoles.S3.ContinuousBackupsRoleArn.ValueString() {
		return true
	}
	if plan.Resources.ServiceRoles.Mssql.SsmNotificationRoleArn.ValueString() != 
		state.Resources.ServiceRoles.Mssql.SsmNotificationRoleArn.ValueString() {
		return true
	}
	if plan.Resources.ServiceRoles.Mssql.Ec2SsmInstanceProfileArn.ValueString() != 
		state.Resources.ServiceRoles.Mssql.Ec2SsmInstanceProfileArn.ValueString() {
		return true
	}
	return false
}
