// Copyright 2023. Clumio, Inc.

// Package clumio_policy contains the Policy resource definition and CRUD implementation.
package clumio_policy

import (
	"context"
	"fmt"
	"strings"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	policyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
)

type policyResource struct {
	client *common.ApiClient
}

// NewPolicyResource is a helper function to simplify the provider implementation.
func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

type replicaModel struct {
	AlternativeReplica types.String `tfsdk:"alternative_replica"`
	PreferredReplica   types.String `tfsdk:"preferred_replica"`
}

type backupTierModel struct {
	BackupTier types.String `tfsdk:"backup_tier"`
}

type pitrConfigModel struct {
	Apply types.String `tfsdk:"apply"`
}

type advancedSettingsModel struct {
	EC2MssqlDatabaseBackup []*replicaModel    `tfsdk:"ec2_mssql_database_backup"`
	EC2MssqlLogBackup      []*replicaModel    `tfsdk:"ec2_mssql_log_backup"`
	MssqlDatabaseBackup    []*replicaModel    `tfsdk:"mssql_database_backup"`
	MssqlLogBackup         []*replicaModel    `tfsdk:"mssql_log_backup"`
	ProtectionGroupBackup  []*backupTierModel `tfsdk:"protection_group_backup"`
	EBSVolumeBackup        []*backupTierModel `tfsdk:"aws_ebs_volume_backup"`
	EC2InstanceBackup      []*backupTierModel `tfsdk:"aws_ec2_instance_backup"`
	RDSPitrConfigSync      []*pitrConfigModel `tfsdk:"aws_rds_config_sync"`
	RDSLogicalBackup       []*backupTierModel `tfsdk:"aws_rds_resource_granular_backup"`
}

type policyOperationModel struct {
	ActionSetting    types.String             `tfsdk:"action_setting"`
	OperationType    types.String             `tfsdk:"type"`
	BackupWindowTz   []*backupWindowModel     `tfsdk:"backup_window_tz"`
	Slas             []*slaModel              `tfsdk:"slas"`
	AdvancedSettings []*advancedSettingsModel `tfsdk:"advanced_settings"`
	BackupAwsRegion  types.String             `tfsdk:"backup_aws_region"`
}

type unitValueModel struct {
	Unit  types.String `tfsdk:"unit"`
	Value types.Int64  `tfsdk:"value"`
}

type rpoModel struct {
	Unit    types.String `tfsdk:"unit"`
	Value   types.Int64  `tfsdk:"value"`
	Offsets types.List   `tfsdk:"offsets"`
}

type slaModel struct {
	RetentionDuration []*unitValueModel `tfsdk:"retention_duration"`
	RPOFrequency      []*rpoModel       `tfsdk:"rpo_frequency"`
}

type backupWindowModel struct {
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

type policyResourceModel struct {
	ID                   types.String            `tfsdk:"id"`
	LockStatus           types.String            `tfsdk:"lock_status"`
	Name                 types.String            `tfsdk:"name"`
	Timezone             types.String            `tfsdk:"timezone"`
	ActivationStatus     types.String            `tfsdk:"activation_status"`
	OrganizationalUnitId types.String            `tfsdk:"organizational_unit_id"`
	Operations           []*policyOperationModel `tfsdk:"operations"`
}

// Metadata returns the data source type name.
func (r *policyResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

// Schema defines the schema for the data source.
func (r *policyResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	unitAttribute := schema.StringAttribute{
		Required: true,
		Description: "The measurement unit of the SLA parameter. Values include" +
			" hours, days, months, and years.",
	}

	valueAttribute := schema.Int64Attribute{
		Required:    true,
		Description: "The measurement value of the SLA parameter.",
	}

	unitValueSchemaAttributes := map[string]schema.Attribute{
		schemaUnit:  unitAttribute,
		schemaValue: valueAttribute,
	}

	rpoValueSchemaAttributes := map[string]schema.Attribute{
		schemaUnit:  unitAttribute,
		schemaValue: valueAttribute,
		schemaOffsets: schema.ListAttribute{
			Optional:    true,
			Description: "The offset values of the SLA parameter.",
			ElementType: types.Int64Type,
		},
	}

	databaseBackupSchemaAttributes := map[string]schema.Attribute{
		schemaAlternativeReplica: schema.StringAttribute{
			Optional:    true,
			Description: fmt.Sprintf(alternativeReplicaDescFmt, "database"),
		},
		schemaPreferredReplica: schema.StringAttribute{
			Optional:    true,
			Description: fmt.Sprintf(preferredReplicaDescFmt, "database"),
		},
	}

	logBackupSchemaAttributes := map[string]schema.Attribute{
		schemaAlternativeReplica: schema.StringAttribute{
			Optional:    true,
			Description: fmt.Sprintf(alternativeReplicaDescFmt, "log"),
		},
		schemaPreferredReplica: schema.StringAttribute{
			Optional:    true,
			Description: fmt.Sprintf(preferredReplicaDescFmt, "log"),
		},
	}

	advancedSettingsSchemaBlocks := map[string]schema.Block{
		schemaEc2MssqlDatabaseBackup: schema.SetNestedBlock{
			Description: mssqlDatabaseBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: databaseBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaEc2MssqlLogBackup: schema.SetNestedBlock{
			Description: mssqlLogBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: logBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaMssqlDatabaseBackup: schema.SetNestedBlock{
			Description: mssqlDatabaseBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: databaseBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaMssqlLogBackup: schema.SetNestedBlock{
			Description: mssqlLogBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: logBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaProtectionGroupBackup: schema.SetNestedBlock{
			Description: "Additional policy configuration settings for the" +
				" protection_group_backup operation. If this operation is not of" +
				" type protection_group_backup, then this field is omitted from" +
				" the response.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaBackupTier: schema.StringAttribute{
						Optional: true,
						Description: "Backup tier to store the backup in. Valid values are:" +
							" cold, frozen",
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaEBSVolumeBackup: schema.SetNestedBlock{
			Description: ebsBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaBackupTier: schema.StringAttribute{
						Optional:    true,
						Description: secureVaultLiteDesc,
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaEC2InstanceBackup: schema.SetNestedBlock{
			Description: ec2BackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaBackupTier: schema.StringAttribute{
						Optional:    true,
						Description: secureVaultLiteDesc,
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaRDSPitrConfigSync: schema.SetNestedBlock{
			Description: rdsPitrConfigSyncDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaApply: schema.StringAttribute{
						Optional:    true,
						Description: pitrConfigDesc,
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaRdsLogicalBackup: schema.SetNestedBlock{
			Description: rdsLogicalBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaBackupTier: schema.StringAttribute{
						Optional:    true,
						Description: rdsLogicalBackupAdvancedSettingDesc,
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
	}

	backupWindowSchemaAttributes := map[string]schema.Attribute{
		schemaStartTime: schema.StringAttribute{
			Description: "The time when the backup window opens." +
				" Specify the start time in the format `hh:mm`," +
				" where `hh` represents the hour of the day and" +
				" `mm` represents the minute of the day based on" +
				" the 24 hour clock.",
			Optional: true,
		},
		schemaEndTime: schema.StringAttribute{
			Description: "The time when the backup window closes." +
				" Specify the end time in the format `hh:mm`," +
				" where `hh` represents the hour of the day and" +
				" `mm` represents the minute of the day based on" +
				" the 24 hour clock. Leave empty if you do not want" +
				" to specify an end time. If the backup window closes" +
				" while a backup is in progress, the entire backup process" +
				" is aborted. The next backup will be performed when the " +
				" backup window re-opens.",
			Optional: true,
			// Use computed property to accept both empty string and null value.
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	slaSchemaBlocks := map[string]schema.Block{
		schemaRetentionDuration: schema.SetNestedBlock{
			Description: "The retention time for this SLA. " +
				"For example, to retain the backup for 1 month," +
				" set unit=months and value=1.",
			NestedObject: schema.NestedBlockObject{
				Attributes: unitValueSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.IsRequired(),
				setvalidator.SizeAtMost(1),
			},
		},
		schemaRpoFrequency: schema.SetNestedBlock{
			Description: "The minimum frequency between " +
				"backups for this SLA. Also known as the " +
				"recovery point objective (RPO) interval. For" +
				" example, to configure the minimum frequency" +
				" between backups to be every 2 days, set " +
				"unit=days and value=2. To configure the SLA " +
				"for on-demand backups, set unit=on_demand " +
				"and leave the value field empty. Also you can " +
				"specify a day of week for Weekly SLA. For example, " +
				"set offsets=[1] will trigger backup on every " +
				"Monday.",
			NestedObject: schema.NestedBlockObject{
				Attributes: rpoValueSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.IsRequired(),
				setvalidator.SizeAtMost(1),
			},
		},
	}

	operationSchemaAttributes := map[string]schema.Attribute{
		schemaActionSetting: schema.StringAttribute{
			Description: "Determines whether the policy should take action" +
				" now or during the specified backup window. Valid values:" +
				"immediate: to start backup process immediately" +
				"window: to start backup in the specified window",
			Required: true,
		},
		schemaOperationType: schema.StringAttribute{
			Description: "The type of operation to be performed. Depending on the type " +
				"selected, `advanced_settings` may also be required. See the API " +
				"Documentation for \"List policies\" for more information about the " +
				"supported types.",
			Required: true,
		},
		schemaBackupAwsRegion: schema.StringAttribute{
			Description: "The region in which this backup is stored. This might be used " +
				"for cross-region backup. Possible values are AWS region string, for " +
				"example: `us-east-1`, `us-west-2`, .... If no value is provided, it " +
				"defaults to in-region (the asset's source region).",
			Optional: true,
		},
	}

	operationSchemaBlocks := map[string]schema.Block{
		schemaBackupWindowTz: schema.SetNestedBlock{
			Description: "The start and end times for the customized" +
				" backup window that reflects the user-defined timezone.",
			NestedObject: schema.NestedBlockObject{
				Attributes: backupWindowSchemaAttributes,
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaAdvancedSettings: schema.SetNestedBlock{
			Description: "Additional operation-specific policy settings.",
			NestedObject: schema.NestedBlockObject{
				Blocks: advancedSettingsSchemaBlocks,
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		schemaSlas: schema.SetNestedBlock{
			Description: "The service level agreement (SLA) for the policy." +
				" A policy can include one or more SLAs. For example, " +
				"a policy can retain daily backups for a month each, " +
				"and monthly backups for a year each.",
			NestedObject: schema.NestedBlockObject{
				Blocks: slaSchemaBlocks,
			},
			Validators: []validator.Set{
				setvalidator.IsRequired(),
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Resource used to schedule backups on" +
			" Clumio supported data sources.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Policy Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaLockStatus: schema.StringAttribute{
				Description: "Policy Lock Status.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The name of the policy.",
				Required:    true,
			},
			schemaTimezone: schema.StringAttribute{
				Description: "The time zone for the policy, in IANA format. For example: " +
					"`America/Los_Angeles`, `America/New_York`, `Etc/UTC`, etc. " +
					"For more information, see the Time Zone Database " +
					"(https://www.iana.org/time-zones) on the IANA website.",
				Optional: true,
				// Use computed property to accept both empty string and null value.
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaActivationStatus: schema.StringAttribute{
				Description: "The status of the policy. Valid values are:" +
					"activated: Backups will take place regularly according to the policy SLA." +
					"deactivated: Backups will not begin until the policy is reactivated." +
					" The assets associated with the policy will have their compliance" +
					" status set to deactivated.",
				Optional: true,
				// Use computed property to accept both empty string and null value.
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the organizational unit" +
					" associated with the policy.",
				Optional: true,
				// Use computed property to accept both empty string and null value.
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			schemaOperations: schema.SetNestedBlock{
				Description: "Each data source to be protected should have details provided in " +
					"the list of operations. These details include information such as how often " +
					"to protect the data source, whether a backup window is desired, which type " +
					"of protection to perform, etc.",
				NestedObject: schema.NestedBlockObject{
					Attributes: operationSchemaAttributes,
					Blocks:     operationSchemaBlocks,
				},
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *policyResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create handles the Create action for the Clumio Policy Resource.
func (r *policyResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	pd := policyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	activationStatus := plan.ActivationStatus.ValueString()
	name := plan.Name.ValueString()
	timezone := plan.Timezone.ValueString()
	policyOperations, diags := mapSchemaOperationsToClumioOperations(ctx,
		plan.Operations)
	resp.Diagnostics.Append(diags...)
	orgUnitId := plan.OrganizationalUnitId.ValueString()
	pdRequest := &models.CreatePolicyDefinitionV1Request{
		ActivationStatus:     &activationStatus,
		Name:                 &name,
		Timezone:             &timezone,
		Operations:           policyOperations,
		OrganizationalUnitId: &orgUnitId,
	}
	res, apiErr := pd.CreatePolicyDefinition(pdRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error creating policy definition.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	plan.ID = types.StringValue(*res.Id)
	apiErr, diags = readPolicyAndUpdateModel(ctx, &plan, pd)
	resp.Diagnostics.Append(diags...)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			errorPolicyReadMsg,
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *policyResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	pd := policyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	// Get current state
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiErr, diags := readPolicyAndUpdateModel(ctx, &state, pd)
	resp.Diagnostics.Append(diags...)
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			state.ID = types.StringValue("")
		} else {
			resp.Diagnostics.AddError(
				errorPolicyReadMsg,
				fmt.Sprintf(errorFmt, string(apiErr.Response)))
			return
		}
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *policyResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := policyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	activationStatus := plan.ActivationStatus.ValueString()
	name := plan.Name.ValueString()
	timezone := plan.Timezone.ValueString()
	policyOperations, policyDiag := mapSchemaOperationsToClumioOperations(ctx,
		plan.Operations)
	if policyDiag != nil {
		resp.Diagnostics.Append(policyDiag...)
	}
	orgUnitId := plan.OrganizationalUnitId.ValueString()
	pdRequest := &models.UpdatePolicyDefinitionV1Request{
		ActivationStatus:     &activationStatus,
		Name:                 &name,
		Timezone:             &timezone,
		Operations:           policyOperations,
		OrganizationalUnitId: &orgUnitId,
	}
	res, apiErr := pd.UpdatePolicyDefinition(plan.ID.ValueString(), nil, pdRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error updating Policy Definition %v.", plan.ID.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error updating Policy Definition %v.", plan.ID.ValueString()),
			fmt.Sprintf(errorFmt, err.Error()))
		return
	}
	apiErr, diags = readPolicyAndUpdateModel(ctx, &plan, pd)
	resp.Diagnostics.Append(diags...)
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			plan.ID = types.StringValue("")
		} else {
			resp.Diagnostics.AddError(
				errorPolicyReadMsg,
				fmt.Sprintf(errorFmt, string(apiErr.Response)))
			return
		}
	}
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *policyResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd := policyDefinitions.NewPolicyDefinitionsV1(r.client.ClumioConfig)
	res, apiErr := pd.DeletePolicyDefinition(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting Policy Definition %v.", state.ID.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	err := common.PollTask(ctx, r.client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting Policy Definition %v.", state.ID.ValueString()),
			fmt.Sprintf(errorFmt, err.Error()))
		return
	}
}

func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func readPolicyAndUpdateModel(ctx context.Context,
	state *policyResourceModel, pd policyDefinitions.PolicyDefinitionsV1Client) (
	*apiutils.APIError, diag.Diagnostics) {
	res, apiErr := pd.ReadPolicyDefinition(state.ID.ValueString(), nil)
	if apiErr != nil {
		tflog.Error(ctx, fmt.Sprintf("Error retrieving policy with ID: %s. Error: %v", state.ID.ValueString(), apiErr))
		return apiErr, nil
	}
	state.LockStatus = types.StringValue(*res.LockStatus)
	state.Name = types.StringValue(*res.Name)
	state.Timezone = types.StringValue(*res.Timezone)
	if res.ActivationStatus != nil {
		state.ActivationStatus = types.StringValue(*res.ActivationStatus)
	}
	if res.OrganizationalUnitId != nil {
		state.OrganizationalUnitId = types.StringValue(*res.OrganizationalUnitId)
	}
	stateOp, diags := mapClumioOperationsToSchemaOperations(ctx, res.Operations)
	state.Operations = stateOp
	return nil, diags
}

// mapSchemaOperationsToClumioOperations maps the schema operations to the Clumio API
// request operations.
func mapSchemaOperationsToClumioOperations(ctx context.Context,
	schemaOperations []*policyOperationModel) ([]*models.PolicyOperationInput,
	diag.Diagnostics) {
	var diags diag.Diagnostics
	policyOperations := make([]*models.PolicyOperationInput, 0)
	for _, operation := range schemaOperations {
		actionSetting := operation.ActionSetting.ValueString()
		operationType := operation.OperationType.ValueString()
		backupAwsRegionPtr := common.GetStringPtr(operation.BackupAwsRegion)

		var backupWindowTz *models.BackupWindow
		if operation.BackupWindowTz != nil {
			startTime := operation.BackupWindowTz[0].StartTime.ValueString()
			endTime := operation.BackupWindowTz[0].EndTime.ValueString()
			backupWindowTz = &models.BackupWindow{
				EndTime:   &endTime,
				StartTime: &startTime,
			}
		}

		advancedSettings := getOperationAdvancedSettings(operation)

		var backupSLAs []*models.BackupSLA
		if operation.Slas != nil {
			backupSLAs = make([]*models.BackupSLA, 0)

			for _, operationSla := range operation.Slas {
				backupSLA := &models.BackupSLA{}
				if operationSla.RetentionDuration != nil {
					unit := operationSla.RetentionDuration[0].Unit.ValueString()
					value := operationSla.RetentionDuration[0].Value.ValueInt64()
					backupSLA.RetentionDuration = &models.RetentionBackupSLAParam{
						Unit:  &unit,
						Value: &value,
					}
				}
				if operationSla.RPOFrequency != nil {
					var offsets []*int64
					unit := operationSla.RPOFrequency[0].Unit.ValueString()
					value := operationSla.RPOFrequency[0].Value.ValueInt64()
					diags = operationSla.RPOFrequency[0].Offsets.ElementsAs(ctx, &offsets, true)
					backupSLA.RpoFrequency = &models.RPOBackupSLAParam{
						Unit:    &unit,
						Value:   &value,
						Offsets: offsets,
					}
				}
				backupSLAs = append(backupSLAs, backupSLA)
			}
		}

		policyOperation := &models.PolicyOperationInput{
			ActionSetting:    &actionSetting,
			BackupWindowTz:   backupWindowTz,
			Slas:             backupSLAs,
			ClumioType:       &operationType,
			AdvancedSettings: advancedSettings,
			BackupAwsRegion:  backupAwsRegionPtr,
		}
		policyOperations = append(policyOperations, policyOperation)
	}
	return policyOperations, diags
}

// mapClumioOperationsToSchemaOperations maps the Operations from the API response to
// the schema operations.
func mapClumioOperationsToSchemaOperations(ctx context.Context,
	operations []*models.PolicyOperation) ([]*policyOperationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	schemaOperations := make([]*policyOperationModel, 0)
	for _, operation := range operations {
		schemaOperation := &policyOperationModel{}
		schemaOperation.ActionSetting = types.StringValue(*operation.ActionSetting)
		schemaOperation.OperationType = types.StringValue(*operation.ClumioType)

		if operation.BackupAwsRegion != nil {
			schemaOperation.BackupAwsRegion = types.StringValue(*operation.BackupAwsRegion)
		}

		if operation.BackupWindowTz != nil {
			window := &backupWindowModel{}
			window.StartTime = types.StringValue(*operation.BackupWindowTz.StartTime)
			window.EndTime = types.StringValue(*operation.BackupWindowTz.EndTime)
			schemaOperation.BackupWindowTz = []*backupWindowModel{window}
		}

		if operation.Slas != nil {
			backupSlas := make([]*slaModel, 0)
			for _, sla := range operation.Slas {
				backupSla := &slaModel{}
				if sla.RetentionDuration != nil {
					backupSla.RetentionDuration = []*unitValueModel{
						{
							Unit:  types.StringValue(*sla.RetentionDuration.Unit),
							Value: types.Int64Value(*sla.RetentionDuration.Value),
						},
					}
				}
				if sla.RpoFrequency != nil {
					offsets, rpoDiags := types.ListValueFrom(ctx,
						types.Int64Type, sla.RpoFrequency.Offsets)
					diags = rpoDiags
					backupSla.RPOFrequency = []*rpoModel{
						{
							Unit:    types.StringValue(*sla.RpoFrequency.Unit),
							Value:   types.Int64Value(*sla.RpoFrequency.Value),
							Offsets: offsets,
						},
					}
				}
				backupSlas = append(backupSlas, backupSla)
			}
			schemaOperation.Slas = backupSlas
		}
		if operation.AdvancedSettings != nil {
			advSettings := &advancedSettingsModel{}
			if operation.AdvancedSettings.Ec2MssqlDatabaseBackup != nil {
				advSettings.EC2MssqlDatabaseBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringValue(
							*operation.AdvancedSettings.Ec2MssqlDatabaseBackup.AlternativeReplica),
						PreferredReplica: types.StringValue(
							*operation.AdvancedSettings.Ec2MssqlDatabaseBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.Ec2MssqlLogBackup != nil {
				advSettings.EC2MssqlLogBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringValue(
							*operation.AdvancedSettings.Ec2MssqlLogBackup.AlternativeReplica),
						PreferredReplica: types.StringValue(
							*operation.AdvancedSettings.Ec2MssqlLogBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.MssqlDatabaseBackup != nil {
				advSettings.MssqlDatabaseBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringValue(
							*operation.AdvancedSettings.MssqlDatabaseBackup.AlternativeReplica),
						PreferredReplica: types.StringValue(
							*operation.AdvancedSettings.MssqlDatabaseBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.MssqlLogBackup != nil {
				advSettings.MssqlLogBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringValue(
							*operation.AdvancedSettings.MssqlLogBackup.AlternativeReplica),
						PreferredReplica: types.StringValue(
							*operation.AdvancedSettings.MssqlLogBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.ProtectionGroupBackup != nil {
				advSettings.ProtectionGroupBackup = []*backupTierModel{
					{
						BackupTier: types.StringValue(
							*operation.AdvancedSettings.ProtectionGroupBackup.BackupTier),
					},
				}
			}
			if operation.AdvancedSettings.AwsEbsVolumeBackup != nil {
				advSettings.EBSVolumeBackup = []*backupTierModel{
					{
						BackupTier: types.StringValue(
							*operation.AdvancedSettings.AwsEbsVolumeBackup.BackupTier),
					},
				}
			}
			if operation.AdvancedSettings.AwsEc2InstanceBackup != nil {
				advSettings.EC2InstanceBackup = []*backupTierModel{
					{
						BackupTier: types.StringValue(
							*operation.AdvancedSettings.AwsEc2InstanceBackup.BackupTier),
					},
				}
			}
			if operation.AdvancedSettings.AwsRdsConfigSync != nil {
				advSettings.RDSPitrConfigSync = []*pitrConfigModel{
					{
						Apply: types.StringValue(
							*operation.AdvancedSettings.AwsRdsConfigSync.Apply),
					},
				}
			}
			if operation.AdvancedSettings.AwsRdsResourceGranularBackup != nil {
				advSettings.RDSLogicalBackup = []*backupTierModel{
					{
						BackupTier: types.StringValue(
							*operation.AdvancedSettings.AwsRdsResourceGranularBackup.BackupTier),
					},
				}
			}
			schemaOperation.AdvancedSettings = []*advancedSettingsModel{advSettings}
		}
		schemaOperations = append(schemaOperations, schemaOperation)
	}

	return schemaOperations, diags
}

// getOperationAdvancedSettings returns the models.PolicyAdvancedSettings after parsing
// the advanced_settings from the schema.
func getOperationAdvancedSettings(
	operation *policyOperationModel) *models.PolicyAdvancedSettings {
	var advancedSettings *models.PolicyAdvancedSettings
	if operation.AdvancedSettings != nil {
		advancedSettings = &models.PolicyAdvancedSettings{}
		if operation.AdvancedSettings[0].EBSVolumeBackup != nil {
			backupTier :=
				operation.AdvancedSettings[0].EBSVolumeBackup[0].BackupTier.ValueString()
			advancedSettings.AwsEbsVolumeBackup = &models.EBSBackupAdvancedSetting{
				BackupTier: &backupTier,
			}
		}
		if operation.AdvancedSettings[0].EC2InstanceBackup != nil {
			backupTier :=
				operation.AdvancedSettings[0].EC2InstanceBackup[0].BackupTier.ValueString()
			advancedSettings.AwsEc2InstanceBackup = &models.EC2BackupAdvancedSetting{
				BackupTier: &backupTier,
			}
		}
		if operation.AdvancedSettings[0].ProtectionGroupBackup != nil {
			backupTier :=
				operation.AdvancedSettings[0].ProtectionGroupBackup[0].BackupTier.ValueString()
			advancedSettings.ProtectionGroupBackup =
				&models.ProtectionGroupBackupAdvancedSetting{
					BackupTier: &backupTier,
				}
		}
		if operation.AdvancedSettings[0].EC2MssqlDatabaseBackup != nil {
			alternativeReplica := operation.AdvancedSettings[0].EC2MssqlDatabaseBackup[0].
				AlternativeReplica.ValueString()
			preferredReplica := operation.AdvancedSettings[0].EC2MssqlDatabaseBackup[0].
				PreferredReplica.ValueString()
			advancedSettings.Ec2MssqlDatabaseBackup =
				&models.EC2MSSQLDatabaseBackupAdvancedSetting{
					AlternativeReplica: &alternativeReplica,
					PreferredReplica:   &preferredReplica,
				}
		}
		if operation.AdvancedSettings[0].EC2MssqlLogBackup != nil {
			alternativeReplica := operation.AdvancedSettings[0].EC2MssqlLogBackup[0].
				AlternativeReplica.ValueString()
			preferredReplica := operation.AdvancedSettings[0].EC2MssqlLogBackup[0].
				PreferredReplica.ValueString()
			advancedSettings.Ec2MssqlLogBackup =
				&models.EC2MSSQLLogBackupAdvancedSetting{
					AlternativeReplica: &alternativeReplica,
					PreferredReplica:   &preferredReplica,
				}
		}
		if operation.AdvancedSettings[0].MssqlDatabaseBackup != nil {
			alternativeReplica := operation.AdvancedSettings[0].MssqlDatabaseBackup[0].
				AlternativeReplica.ValueString()
			preferredReplica := operation.AdvancedSettings[0].MssqlDatabaseBackup[0].
				PreferredReplica.ValueString()
			advancedSettings.MssqlDatabaseBackup =
				&models.MSSQLDatabaseBackupAdvancedSetting{
					AlternativeReplica: &alternativeReplica,
					PreferredReplica:   &preferredReplica,
				}
		}
		if operation.AdvancedSettings[0].MssqlLogBackup != nil {
			alternativeReplica := operation.AdvancedSettings[0].MssqlLogBackup[0].
				AlternativeReplica.ValueString()
			preferredReplica := operation.AdvancedSettings[0].MssqlLogBackup[0].
				PreferredReplica.ValueString()
			advancedSettings.MssqlLogBackup =
				&models.MSSQLLogBackupAdvancedSetting{
					AlternativeReplica: &alternativeReplica,
					PreferredReplica:   &preferredReplica,
				}
		}
		if operation.AdvancedSettings[0].RDSPitrConfigSync != nil {
			apply := operation.AdvancedSettings[0].RDSPitrConfigSync[0].Apply.ValueString()
			advancedSettings.AwsRdsConfigSync =
				&models.RDSConfigSyncAdvancedSetting{
					Apply: &apply,
				}
		}
		if operation.AdvancedSettings[0].RDSLogicalBackup != nil {
			backupTier :=
				operation.AdvancedSettings[0].RDSLogicalBackup[0].BackupTier.ValueString()
			advancedSettings.AwsRdsResourceGranularBackup =
				&models.RDSLogicalBackupAdvancedSetting{
					BackupTier: &backupTier,
				}
		}
	}
	return advancedSettings
}
