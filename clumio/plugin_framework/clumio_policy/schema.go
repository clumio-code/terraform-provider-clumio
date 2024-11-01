// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_policy Terraform resource.

package clumio_policy

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// policyResourceModel is the resource model for the clumio_policy Terraform resource. It represents
// the schema of the resource and the data it holds. This schema is used by customers to configure
// the resource and by the Clumio provider to read and write the resource.
type policyResourceModel struct {
	ID               types.String            `tfsdk:"id"`
	LockStatus       types.String            `tfsdk:"lock_status"`
	Name             types.String            `tfsdk:"name"`
	Timezone         types.String            `tfsdk:"timezone"`
	ActivationStatus types.String            `tfsdk:"activation_status"`
	Operations       []*policyOperationModel `tfsdk:"operations"`
}

// replicaModel maps to some of the attributes in the advancedSettingsModel which require a
// the preferred and alternative replica to be specified.
type replicaModel struct {
	AlternativeReplica types.String `tfsdk:"alternative_replica"`
	PreferredReplica   types.String `tfsdk:"preferred_replica"`
}

// backupTierModel maps to some of the attributes in advancedSettingsModel which require a
// backup tier to be specified.
type backupTierModel struct {
	BackupTier types.String `tfsdk:"backup_tier"`
}

// pitrConfigModel maps to the RDSPitrConfigSync attribute in advancedSettingsModel which
// determines when the configuration will be applied.
type pitrConfigModel struct {
	Apply types.String `tfsdk:"apply"`
}

// advancedSettingsModel maps to the AdvancedSettings attribute in policyOperationModel which
// contains additional operation-specific policy settings.
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

// policyOperationModel maps to the Operations attribute in policyResourceModel and contains
// information such as how often to protect the data source, whether a backup window is desired,
// which type of protection to perform, etc
type policyOperationModel struct {
	ActionSetting    types.String             `tfsdk:"action_setting"`
	OperationType    types.String             `tfsdk:"type"`
	BackupWindowTz   []*backupWindowModel     `tfsdk:"backup_window_tz"`
	Slas             []*slaModel              `tfsdk:"slas"`
	AdvancedSettings []*advancedSettingsModel `tfsdk:"advanced_settings"`
	BackupAwsRegion  types.String             `tfsdk:"backup_aws_region"`
	Timezone         types.String             `tfsdk:"timezone"`
}

// unitValueModel maps tho the RetentionDuration attribute in slaModel and and it provides the unit
// and value for the retention duration.
type unitValueModel struct {
	Unit  types.String `tfsdk:"unit"`
	Value types.Int64  `tfsdk:"value"`
}

// unitValueModel maps tho the RPOFrequency attribute in slaModel and and it provides the unit,
// value and offsets for the RPO Fequency.
type rpoModel struct {
	Unit    types.String `tfsdk:"unit"`
	Value   types.Int64  `tfsdk:"value"`
	Offsets types.List   `tfsdk:"offsets"`
}

// slaModel maps to the Slas attribute in policyOperationModel and it refers to the service level
// agreement (SLA) for the policy.
type slaModel struct {
	RetentionDuration []*unitValueModel `tfsdk:"retention_duration"`
	RPOFrequency      []*rpoModel       `tfsdk:"rpo_frequency"`
}

// backupWindowModel maps to the BackupWindowTz attribute policyOperationModel and it refers to
// time window during which the backups can run.
type backupWindowModel struct {
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
}

// Schema defines the structure and constraints of the clumio_policy Terraform resource. Schema is a
// method on the policyResource struct. It sets the schema for the clumio_policy Terraform resource,
// which is used to create a policy for scheduling backups on Clumio supported data sources.
func (r *policyResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	retentionUnitAttribute := schema.StringAttribute{
		Required: true,
		Description: "The measurement unit of the SLA parameter. Values include" +
			" days, weeks, months and years.",
	}

	rpoUnitAttribute := schema.StringAttribute{
		Required: true,
		Description: "The measurement unit of the SLA parameter. Values include" +
			" minutes, hours, days, weeks, months and years.",
	}

	valueAttribute := schema.Int64Attribute{
		Required:    true,
		Description: "The measurement value of the SLA parameter.",
	}

	retentionSchemaAttributes := map[string]schema.Attribute{
		schemaUnit:  retentionUnitAttribute,
		schemaValue: valueAttribute,
	}

	rpoSchemaAttributes := map[string]schema.Attribute{
		schemaUnit:  rpoUnitAttribute,
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
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
		schemaEc2MssqlLogBackup: schema.SetNestedBlock{
			Description: mssqlLogBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: logBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
		schemaMssqlDatabaseBackup: schema.SetNestedBlock{
			Description: mssqlDatabaseBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: databaseBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
		schemaMssqlLogBackup: schema.SetNestedBlock{
			Description: mssqlLogBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: logBackupSchemaAttributes,
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
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
							" `cold` and `frozen`.\n\t- `cold` = Clumio SecureVault Standard\n\t" +
							"- `frozen` = Clumio SecureVault Archive",
					},
				},
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
		schemaEBSVolumeBackup: schema.SetNestedBlock{
			Description: ebsBackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaBackupTier: schema.StringAttribute{
						Optional:    true,
						Description: ebsEc2BackupTierDesc,
					},
				},
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
		schemaEC2InstanceBackup: schema.SetNestedBlock{
			Description: ec2BackupDesc,
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					schemaBackupTier: schema.StringAttribute{
						Optional:    true,
						Description: ebsEc2BackupTierDesc,
					},
				},
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
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
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
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
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
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
				Attributes: retentionSchemaAttributes,
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.IsRequired()),
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
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
				Attributes: rpoSchemaAttributes,
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.IsRequired()),
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
	}

	operationSchemaAttributes := map[string]schema.Attribute{
		schemaActionSetting: schema.StringAttribute{
			Description: "Determines whether the policy should take action" +
				" now or during the specified backup window. Valid values are: `immediate` and `window`." +
				" `immediate` starts the backup process immediately while `window` starts the backup" +
				" in the specified window.",
			Required: true,
		},
		schemaOperationType: schema.StringAttribute{
			Description: "The type of operation to be performed. Depending on the type " +
				"selected, `advanced_settings` may also be required. See the [API " +
				"Documentation for List policies]" +
				"(https://help.clumio.com/reference/list-policy-definitions) for more information " +
				"about the supported types.",
			Required: true,
		},
		schemaBackupAwsRegion: schema.StringAttribute{
			Description: "The region in which this backup is stored. This might be used " +
				"for cross-region backup. Possible values are AWS region string, for " +
				"example: `us-east-1`, `us-west-2`, .... If no value is provided, it " +
				"defaults to in-region (the asset's source region).",
			Optional: true,
		},
		schemaTimezone: schema.StringAttribute{
			Description: "The time zone for the policy, in IANA format. For example: " +
				"`America/Los_Angeles`, `America/New_York`, `Etc/UTC`, etc. " +
				"For more information, see the Time Zone Database " +
				"(https://www.iana.org/time-zones) on the IANA website.",
			Optional: true,
			// Use computed property to accept null value.
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
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
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
			},
		},
		schemaAdvancedSettings: schema.SetNestedBlock{
			Description: "Additional operation-specific policy settings.",
			NestedObject: schema.NestedBlockObject{
				Blocks: advancedSettingsSchemaBlocks,
			},
			Validators: []validator.Set{
				common.WrapSetValidator(setvalidator.SizeAtMost(1)),
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
				common.WrapSetValidator(setvalidator.IsRequired()),
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Resource used to schedule backups on" +
			" Clumio supported data sources.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier of the policy.",
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
				Description: "The user-assigned name of the policy. " +
					"Note that having identical names for different policies is permissible.",
				Required: true,
			},
			schemaTimezone: schema.StringAttribute{
				Description: "The time zone for the policy, in IANA format. For example: " +
					"`America/Los_Angeles`, `America/New_York`, `Etc/UTC`, etc. " +
					"For more information, see the Time Zone Database " +
					"(https://www.iana.org/time-zones) on the IANA website.",
				Optional: true,
				DeprecationMessage: "Global timezone is deprecated. Instead, use the timezone " +
					"attribute within each policy operation.",
				// Use computed property to accept null value.
				Computed: true,
			},
			schemaActivationStatus: schema.StringAttribute{
				Description: "The status of the policy. Valid values are: `activated` and `deactivated`." +
					" `activated` backups will take place regularly according to the policy SLA." +
					" `deactivated` backups will not begin until the policy is reactivated." +
					" The assets associated with the policy will have their compliance" +
					" status set to deactivated.",
				Optional: true,
				// Use computed property to accept both empty string and null value.
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(activationStatusActivated, activationStatusDectivated),
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
					common.WrapSetValidator(setvalidator.IsRequired()),
				},
			},
		},
	}
}
