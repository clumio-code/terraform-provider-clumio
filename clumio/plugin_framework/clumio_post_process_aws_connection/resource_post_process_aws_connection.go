// Copyright 2023. Clumio, Inc.

// clumio_post_process_aws_connection definition and CRUD implementation.
package clumio_post_process_aws_connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"strings"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/post_process_aws_connection"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sourceConfigInfo struct {
	sourceKey string
	isConfig  bool
}

var (
	// protectInfoMap is the mapping of the the datasource to the resource parameter and
	// if a config section is required, then isConfig will be true.
	protectInfoMap = map[string]sourceConfigInfo{
		"ebs": {
			sourceKey: "ProtectEBSVersion",
			isConfig:  false,
		},
		"rds": {
			sourceKey: "ProtectRDSVersion",
			isConfig:  false,
		},
		"ec2_mssql": {
			sourceKey: "ProtectEC2MssqlVersion",
			isConfig:  false,
		},
		"warm_tier": {
			sourceKey: "ProtectWarmTierVersion",
			isConfig:  true,
		},
		"s3": {
			sourceKey: "ProtectS3Version",
			isConfig:  false,
		},
		"dynamodb": {
			sourceKey: "ProtectDynamoDBVersion",
			isConfig:  false,
		},
	}
	// warmtierInfoMap is the mapping of the the warm tier datasource to the resource
	// parameter and if a config section is required, then isConfig will be true.
	warmtierInfoMap = map[string]sourceConfigInfo{
		"dynamodb": {
			sourceKey: "ProtectWarmTierDynamoDBVersion",
			isConfig:  false,
		},
	}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &postProcessAWSConnectionResource{}
	_ resource.ResourceWithConfigure = &postProcessAWSConnectionResource{}
)

// NewPostProcessAWSConnectionResource is a helper function to simplify the provider implementation.
func NewPostProcessAWSConnectionResource() resource.Resource {
	return &postProcessAWSConnectionResource{}
}

// postProcessAWSConnectionResource is the resource implementation.
type postProcessAWSConnectionResource struct {
	client *common.ApiClient
}

// postProcessAWSConnectionResourceModel model
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
}

// Metadata returns the resource type name.
func (r *postProcessAWSConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_process_aws_connection"
}

// Schema defines the schema for the resource.
func (r *postProcessAWSConnectionResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Post-Process Clumio AWS Connection Resource used to" +
			" post-process AWS connection to Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
			},
			schemaToken: schema.StringAttribute{
				Description: "The AWS integration ID token.",
				Required:    true,
			},
			schemaRoleExternalId: schema.StringAttribute{
				Description: "A key that must be used by Clumio to assume the service role" +
					" in your account. This should be a secure string, like a password," +
					" but it does not need to be remembered (random characters are best).",
				Required: true,
			},
			schemaAccountId: schema.StringAttribute{
				Description: "The AWS Customer Account ID.",
				Required:    true,
			},
			schemaRegion: schema.StringAttribute{
				Description: "The AWS Region.",
				Required:    true,
			},
			schemaRoleArn: schema.StringAttribute{
				Description: "Clumio IAM Role Arn.",
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *postProcessAWSConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create handles the Create action for the PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postProcessAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = r.clumioPostProcessAWSConnectionCommon(ctx, plan, "Create")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountId := plan.AccountID.ValueString()
	awsRegion := plan.Region.ValueString()
	token := plan.Token.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf("%v/%v/%v", accountId, awsRegion, token))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read handles the Read action for the PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update handles the Update action for the PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state postProcessAWSConnectionResourceModel
	var plan postProcessAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	diags = r.clumioPostProcessAWSConnectionCommon(ctx, plan, "Update")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete handles the Create action for the PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postProcessAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = r.clumioPostProcessAWSConnectionCommon(ctx, state, "Delete")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// clumioPostProcessAWSConnectionCommon contains the common logic for all CRUD operations
// of PostProcessAWSConnection resource.
func (r *postProcessAWSConnectionResource) clumioPostProcessAWSConnectionCommon(
	_ context.Context, state postProcessAWSConnectionResourceModel, eventType string) diag.Diagnostics {
	postProcessAwsConnection := aws_connections.NewPostProcessAwsConnectionV1(
		r.client.ClumioConfig)
	accountId := state.AccountID.ValueString()
	awsRegion := state.Region.ValueString()
	roleArn := state.RoleArn.ValueString()
	token := state.Token.ValueString()
	roleExternalId := state.RoleExternalID.ValueString()
	clumioEventPubId := state.ClumioEventPubID.ValueString()
	schemaPropertiesElements := state.Properties.Elements()
	propertiesMap := make(map[string]*string)
	for key, val := range schemaPropertiesElements {
		valStr := val.String()
		propertiesMap[key] = &valStr
	}

	templateConfig, err := GetTemplateConfiguration(state, true, true)
	if err != nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Error forming template configuration.",
			fmt.Sprintf(errorFmt, err))
		return diagnostics
	}
	templateConfig["insights"] = templateConfig["discover"]
	delete(templateConfig, "discover")
	configBytes, err := json.Marshal(templateConfig)
	if err != nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Error marshalling template configuration.",
			fmt.Sprintf(errorFmt, err))
		return diagnostics
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
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Error in invoking Post-process Clumio AWS Connection.",
			fmt.Sprintf(errorFmt, err))
		return diagnostics
	}
	return nil
}

// GetTemplateConfiguration returns the template configuration.
func GetTemplateConfiguration(
	model postProcessAWSConnectionResourceModel, isCamelCase bool, isConsolidated bool) (
	map[string]interface{}, error) {
	templateConfigs := make(map[string]interface{})
	configMap, err := getConfigMapForKey(model.ConfigVersion.ValueString(), false)
	if err != nil {
		return nil, err
	}
	if configMap == nil {
		return templateConfigs, nil
	}
	templateConfigs["config"] = configMap
	if !isConsolidated {
		discoverMap, err := getConfigMapForKey(model.DiscoverVersion.ValueString(), true)
		if err != nil {
			return nil, err
		}
		if discoverMap == nil {
			return templateConfigs, nil
		}
		templateConfigs["discover"] = discoverMap
	}

	protectMap, err := getConfigMapForKey(model.ProtectConfigVersion.ValueString(), true)
	if err != nil {
		return nil, err
	}
	if protectMap == nil {
		return templateConfigs, nil
	}
	err = populateConfigMap(model, protectInfoMap, protectMap, isCamelCase)
	if err != nil {
		return nil, err
	}
	warmTierKey := "warm_tier"
	if isCamelCase {
		warmTierKey = common.SnakeCaseToCamelCase(warmTierKey)
	}
	if protectWarmtierMap, ok := protectMap[warmTierKey]; ok {
		err = populateConfigMap(
			model, warmtierInfoMap, protectWarmtierMap.(map[string]interface{}), isCamelCase)
		if err != nil {
			return nil, err
		}
	}
	if isConsolidated {
		templateConfigs["consolidated"] = protectMap
	} else {
		templateConfigs["protect"] = protectMap
	}
	return templateConfigs, nil
}

// populateConfigMap returns protect configuration information for the configs
// in the configInfoMap.
func populateConfigMap(model postProcessAWSConnectionResourceModel,
	configInfoMap map[string]sourceConfigInfo, configMap map[string]interface{},
	isCamelCase bool) error {

	for source, sourceInfo := range configInfoMap {
		configMapKey := source
		if isCamelCase {
			configMapKey = common.SnakeCaseToCamelCase(source)
		}
		reflectVal := reflect.ValueOf(&model)
		fieldVal := reflect.Indirect(reflectVal).FieldByName(sourceInfo.sourceKey)
		fieldStringVal := fieldVal.Interface().(types.String)
		protectSourceMap, err := getConfigMapForKey(
			fieldStringVal.ValueString(), sourceInfo.isConfig)
		if err != nil {
			return err
		}
		if protectSourceMap != nil {
			configMap[configMapKey] = protectSourceMap
		}
	}
	return nil
}

// getConfigMapForKey returns a config map for the key if it exists in ResourceData.
func getConfigMapForKey(val string, isConfig bool) (map[string]interface{}, error) {
	var mapToReturn map[string]interface{}
	if val != "" {
		keyMap := make(map[string]interface{})
		majorVersion, minorVersion, err := parseVersion(val)
		if err != nil {
			return nil, err
		}
		keyMap["enabled"] = true
		keyMap["version"] = majorVersion
		keyMap["minorVersion"] = minorVersion
		mapToReturn = keyMap
		// If isConfig is true it wraps the keyMap with another map with "config" as the key.
		if isConfig {
			configMap := make(map[string]interface{})
			configMap["config"] = keyMap
			mapToReturn = configMap
		}
	}
	return mapToReturn, nil
}

// parseVersion parses the version and minorVersion given the version string.
func parseVersion(version string) (string, string, error) {
	splits := strings.Split(version, ".")
	switch len(splits) {
	case 1:
		return version, "", nil
	case 2:
		return splits[0], splits[1], nil
	default:
		return "", "", errors.New(fmt.Sprintf("Invalid version: %v", version))
	}
}
