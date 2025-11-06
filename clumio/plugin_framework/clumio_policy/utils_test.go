// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_policy

import (
	"context"
	"testing"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

var (
	actionSetting                  = "immediate"
	operationType                  = "aws_ebs_volume_backup"
	operationType2                 = "aws_ec2_instance_backup"
	operationType3                 = "aws_rds_resource_aws_snapshot"
	operationType4                 = "aws_rds_resource_granular_backup"
	operationType5                 = "ec2_mssql_database_backup"
	operationType6                 = "ec2_mssql_log_backup"
	operationType7                 = "mssql_database_backup"
	operationType8                 = "mssql_log_backup"
	operationType9                 = "protection_group_backup"
	operationType10                = "aws_s3_continuous_backup"
	operationType11                = "aws_iceberg_table_backup"
	retUnit                        = "days"
	retValue                       = int64(5)
	retUnit2                       = "hours"
	retValue2                      = int64(6)
	rpoUnit                        = "days"
	rpoValue                       = int64(1)
	rpoUnit2                       = "hours"
	rpoValue2                      = int64(2)
	offset                         = int64(1)
	backupRegion                   = "us-east-1"
	backupTier                     = "mock-backup-tier"
	startTime                      = "start-time"
	endTime                        = "end-time"
	rdsConfigSyncApply             = "immediate"
	preferredReplica               = "mock-preferred-replica"
	alternativeReplica             = "mock-alternate-replica"
	disableEventbridgeNotification = true
)

// Unit test for the utility function to convert ClumioOperations to SchemaOperations. This tests
// an operation with multiple SLAs.
func TestMapClumioOperationsToSchemaOperations(t *testing.T) {

	ctx := context.Background()

	modelOperations := []*models.PolicyOperation{
		{
			ActionSetting:   &actionSetting,
			BackupAwsRegion: &backupRegion,
			BackupWindowTz: &models.BackupWindow{
				EndTime:   &endTime,
				StartTime: &startTime,
			},
			ClumioType: &operationType,
			Slas: []*models.BackupSLA{
				{
					RetentionDuration: &models.RetentionBackupSLAParam{
						Unit:  &retUnit,
						Value: &retValue,
					},
					RpoFrequency: &models.RPOBackupSLAParam{
						Unit:    &rpoUnit,
						Value:   &rpoValue,
						Offsets: []*int64{&offset},
					},
				},
				{
					RetentionDuration: &models.RetentionBackupSLAParam{
						Unit:  &retUnit2,
						Value: &retValue2,
					},
					RpoFrequency: &models.RPOBackupSLAParam{
						Unit:    &rpoUnit2,
						Value:   &rpoValue2,
						Offsets: []*int64{&offset},
					},
				},
			},
		},
	}

	schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
	assert.Nil(t, diags)

	// Ensure the operation attributes are correct.
	modelOp := modelOperations[0]
	schemaOp := schemaOperations[0]
	assert.Equal(t, *modelOp.ActionSetting, schemaOp.ActionSetting.ValueString())
	assert.Equal(t, *modelOp.BackupAwsRegion, schemaOp.BackupAwsRegion.ValueString())
	assert.Equal(t, *modelOp.BackupWindowTz.StartTime,
		schemaOp.BackupWindowTz[0].StartTime.ValueString())
	assert.Equal(t, *modelOp.BackupWindowTz.EndTime,
		schemaOp.BackupWindowTz[0].EndTime.ValueString())
	assert.Equal(t, *modelOp.ClumioType, schemaOp.OperationType.ValueString())

	// Ensure the first SLA's attributes are correct.
	modelSla := *modelOperations[0].Slas[0]
	schemaSla := schemaOperations[0].Slas[0]
	assert.Equal(t, *modelSla.RetentionDuration.Unit,
		schemaSla.RetentionDuration[0].Unit.ValueString())
	assert.Equal(t, *modelSla.RetentionDuration.Value,
		schemaSla.RetentionDuration[0].Value.ValueInt64())
	assert.Equal(t, *modelSla.RpoFrequency.Unit, schemaSla.RPOFrequency[0].Unit.ValueString())
	assert.Equal(t, *modelSla.RpoFrequency.Value, schemaSla.RPOFrequency[0].Value.ValueInt64())
	var offsets []*int64
	diags = schemaSla.RPOFrequency[0].Offsets.ElementsAs(ctx, &offsets, true)
	assert.Nil(t, diags)

	// Ensure the second SLA's attributes are correct.
	modelSla = *modelOperations[0].Slas[1]
	schemaSla = schemaOperations[0].Slas[1]
	assert.Equal(t, *modelSla.RetentionDuration.Unit,
		schemaSla.RetentionDuration[0].Unit.ValueString())
	assert.Equal(t, *modelSla.RetentionDuration.Value,
		schemaSla.RetentionDuration[0].Value.ValueInt64())
	assert.Equal(t, *modelSla.RpoFrequency.Unit,
		schemaSla.RPOFrequency[0].Unit.ValueString())
	assert.Equal(t, *modelSla.RpoFrequency.Value,
		schemaSla.RPOFrequency[0].Value.ValueInt64())
	var offsets2 []*int64
	diags = schemaSla.RPOFrequency[0].Offsets.ElementsAs(ctx, &offsets2, true)
	assert.Nil(t, diags)
	assert.Equal(t, *modelSla.RpoFrequency.Offsets[0], *offsets2[0])
}

// Unit tests for MapClumioOperationsToSchemaOperations containing different Advanced Settings
// corresponding to different operations.
func TestMapClumioOperationsToSchemaOperationsAdvSettings(t *testing.T) {
	ctx := context.Background()
	modelOperations := []*models.PolicyOperation{
		{
			ActionSetting:   &actionSetting,
			BackupAwsRegion: &backupRegion,
			BackupWindowTz: &models.BackupWindow{
				EndTime:   &endTime,
				StartTime: &startTime,
			},
			Slas: []*models.BackupSLA{
				{
					RetentionDuration: &models.RetentionBackupSLAParam{
						Unit:  &retUnit,
						Value: &retValue,
					},
					RpoFrequency: &models.RPOBackupSLAParam{
						Unit:  &rpoUnit,
						Value: &rpoValue,
					},
				},
			},
		},
	}

	t.Run("Test EBS Volume Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			AwsEbsVolumeBackup: &models.EBSBackupAdvancedSetting{
				BackupTier: &backupTier,
			},
		}
		modelOperations[0].ClumioType = &operationType
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.AwsEbsVolumeBackup.BackupTier,
			schemaOpAdvSettings.EBSVolumeBackup[0].BackupTier.ValueString())
	})
	t.Run("Test EC2 Instance Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			AwsEc2InstanceBackup: &models.EC2BackupAdvancedSetting{
				BackupTier: &backupTier,
			},
		}
		modelOperations[0].ClumioType = &operationType2
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.AwsEc2InstanceBackup.BackupTier,
			schemaOpAdvSettings.EC2InstanceBackup[0].BackupTier.ValueString())
	})
	t.Run("Test RDS Config Sync Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			AwsRdsConfigSync: &models.RDSConfigSyncAdvancedSetting{
				Apply: &rdsConfigSyncApply,
			},
		}
		modelOperations[0].ClumioType = &operationType3
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.AwsRdsConfigSync.Apply,
			schemaOpAdvSettings.RDSPitrConfigSync[0].Apply.ValueString())
	})
	t.Run("Test RDS Granular Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			AwsRdsResourceGranularBackup: &models.RDSLogicalBackupAdvancedSetting{
				BackupTier: &backupTier,
			},
		}
		modelOperations[0].ClumioType = &operationType4
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.AwsRdsResourceGranularBackup.BackupTier,
			schemaOpAdvSettings.RDSLogicalBackup[0].BackupTier.ValueString())
	})
	t.Run("Test EC2 MSSQL Database Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			Ec2MssqlDatabaseBackup: &models.EC2MSSQLDatabaseBackupAdvancedSetting{
				AlternativeReplica: &alternativeReplica,
				PreferredReplica:   &preferredReplica,
			},
		}
		modelOperations[0].ClumioType = &operationType5
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.Ec2MssqlDatabaseBackup.PreferredReplica,
			schemaOpAdvSettings.EC2MssqlDatabaseBackup[0].PreferredReplica.ValueString())
		assert.Equal(t, *modelOpAdvSettings.Ec2MssqlDatabaseBackup.AlternativeReplica,
			schemaOpAdvSettings.EC2MssqlDatabaseBackup[0].AlternativeReplica.ValueString())
	})
	t.Run("Test EC2 MSSQL Log Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			Ec2MssqlLogBackup: &models.EC2MSSQLLogBackupAdvancedSetting{
				AlternativeReplica: &alternativeReplica,
				PreferredReplica:   &preferredReplica,
			},
		}
		modelOperations[0].ClumioType = &operationType6
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.Ec2MssqlLogBackup.PreferredReplica,
			schemaOpAdvSettings.EC2MssqlLogBackup[0].PreferredReplica.ValueString())
		assert.Equal(t, *modelOpAdvSettings.Ec2MssqlLogBackup.AlternativeReplica,
			schemaOpAdvSettings.EC2MssqlLogBackup[0].AlternativeReplica.ValueString())
	})
	t.Run("Test Mssql Database Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			MssqlDatabaseBackup: &models.MSSQLDatabaseBackupAdvancedSetting{
				AlternativeReplica: &alternativeReplica,
				PreferredReplica:   &preferredReplica,
			},
		}
		modelOperations[0].ClumioType = &operationType7
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.MssqlDatabaseBackup.PreferredReplica,
			schemaOpAdvSettings.MssqlDatabaseBackup[0].PreferredReplica.ValueString())
		assert.Equal(t, *modelOpAdvSettings.MssqlDatabaseBackup.AlternativeReplica,
			schemaOpAdvSettings.MssqlDatabaseBackup[0].AlternativeReplica.ValueString())
	})
	t.Run("Test Mssql Log Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			MssqlLogBackup: &models.MSSQLLogBackupAdvancedSetting{
				AlternativeReplica: &alternativeReplica,
				PreferredReplica:   &preferredReplica,
			},
		}
		modelOperations[0].ClumioType = &operationType8
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.MssqlLogBackup.PreferredReplica,
			schemaOpAdvSettings.MssqlLogBackup[0].PreferredReplica.ValueString())
		assert.Equal(t, *modelOpAdvSettings.MssqlLogBackup.AlternativeReplica,
			schemaOpAdvSettings.MssqlLogBackup[0].AlternativeReplica.ValueString())
	})
	t.Run("Test Protection Group Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			ProtectionGroupBackup: &models.ProtectionGroupBackupAdvancedSetting{
				BackupTier: &backupTier,
			},
		}
		modelOperations[0].ClumioType = &operationType9
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.ProtectionGroupBackup.BackupTier,
			schemaOpAdvSettings.ProtectionGroupBackup[0].BackupTier.ValueString())
	})
	t.Run("Test S3 Continuous Backup Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			ProtectionGroupContinuousBackup: &models.ProtectionGroupContinuousBackupAdvancedSetting{
				DisableEventbridgeNotification: &disableEventbridgeNotification,
			},
		}
		modelOperations[0].ClumioType = &operationType10
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.ProtectionGroupContinuousBackup.DisableEventbridgeNotification,
			schemaOpAdvSettings.S3ContinuousBackup[0].DisableEventbridgeNotification.ValueBool())
	})

	t.Run("Test Iceberg Advanced Setting", func(t *testing.T) {
		modelOperations[0].AdvancedSettings = &models.PolicyAdvancedSettings{
			AwsIcebergTableBackup: &models.IcebergBackupAdvancedSetting{
				BackupTier: &backupTier,
			},
		}
		modelOperations[0].ClumioType = &operationType11
		schemaOperations, diags := mapClumioOperationsToSchemaOperations(ctx, modelOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, *modelOpAdvSettings.AwsIcebergTableBackup.BackupTier,
			schemaOpAdvSettings.IcebergTableBackup[0].BackupTier.ValueString())
	})
}

// Unit test for the utility function to convert SchemaOperations to ClumioOperations. This tests
// an operation with multiple SLAs.
func TestMapSchemaOperationsToClumioOperations(t *testing.T) {
	ctx := context.Background()
	offsets, diags := types.ListValueFrom(ctx,
		types.Int64Type, []*int64{&offset})
	assert.Nil(t, diags)
	schemaOperations := []*policyOperationModel{
		{
			ActionSetting: basetypes.NewStringValue(actionSetting),
			OperationType: basetypes.NewStringValue(operationType),
			BackupWindowTz: []*backupWindowModel{
				{
					StartTime: basetypes.NewStringValue(startTime),
					EndTime:   basetypes.NewStringValue(endTime),
				},
			},
			Slas: []*slaModel{
				{
					RetentionDuration: []*unitValueModel{
						{
							Unit:  basetypes.NewStringValue(retUnit),
							Value: basetypes.NewInt64Value(retValue),
						},
					},
					RPOFrequency: []*rpoModel{
						{
							Unit:    basetypes.NewStringValue(rpoUnit),
							Value:   basetypes.NewInt64Value(rpoValue),
							Offsets: offsets,
						},
					},
				},
				{
					RetentionDuration: []*unitValueModel{
						{
							Unit:  basetypes.NewStringValue(retUnit2),
							Value: basetypes.NewInt64Value(retValue2),
						},
					},
					RPOFrequency: []*rpoModel{
						{
							Unit:    basetypes.NewStringValue(rpoUnit2),
							Value:   basetypes.NewInt64Value(rpoValue2),
							Offsets: offsets,
						},
					},
				},
			},
			BackupAwsRegion: basetypes.NewStringValue(backupRegion),
		},
	}
	modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
	assert.Nil(t, diags)

	// Ensure the operation attributes are correct.
	modelOp := modelOperations[0]
	schemaOp := schemaOperations[0]
	assert.Equal(t, schemaOp.ActionSetting.ValueString(), *modelOp.ActionSetting)
	assert.Equal(t, schemaOp.OperationType.ValueString(), *modelOp.ClumioType)
	assert.Equal(t, schemaOp.BackupAwsRegion.ValueString(), *modelOp.BackupAwsRegion)
	assert.Equal(t, schemaOp.BackupWindowTz[0].StartTime.ValueString(),
		*modelOp.BackupWindowTz.StartTime)
	assert.Equal(t, schemaOp.BackupWindowTz[0].EndTime.ValueString(),
		*modelOp.BackupWindowTz.EndTime)

	// Ensure the first SLA's attributes are correct.
	modelSla := modelOp.Slas[0]
	schemaSla := schemaOp.Slas[0]
	assert.Equal(t, schemaSla.RetentionDuration[0].Unit.ValueString(),
		*modelSla.RetentionDuration.Unit)
	assert.Equal(t, schemaSla.RetentionDuration[0].Value.ValueInt64(),
		*modelSla.RetentionDuration.Value)
	assert.Equal(t, schemaSla.RPOFrequency[0].Unit.ValueString(), *modelSla.RpoFrequency.Unit)
	assert.Equal(t, schemaSla.RPOFrequency[0].Value.ValueInt64(), *modelSla.RpoFrequency.Value)
	assert.Equal(t, offset, *modelSla.RpoFrequency.Offsets[0])

	// Ensure the second SLA's attributes are correct.
	modelSla = modelOp.Slas[1]
	schemaSla = schemaOp.Slas[1]
	assert.Equal(t, schemaSla.RetentionDuration[0].Unit.ValueString(),
		*modelSla.RetentionDuration.Unit)
	assert.Equal(t, schemaSla.RetentionDuration[0].Value.ValueInt64(),
		*modelSla.RetentionDuration.Value)
	assert.Equal(t, schemaSla.RPOFrequency[0].Unit.ValueString(), *modelSla.RpoFrequency.Unit)
	assert.Equal(t, schemaSla.RPOFrequency[0].Value.ValueInt64(), *modelSla.RpoFrequency.Value)
	assert.Equal(t, offset, *modelSla.RpoFrequency.Offsets[0])
}

// Unit tests for MapClumioOperationsToSchemaOperations containing different Advanced Settings
// corresponding to different operations.
func TestMapSchemaOperationsToClumioOperationsAdvSettings(t *testing.T) {

	ctx := context.Background()
	offsets, diags := types.ListValueFrom(ctx,
		types.Int64Type, []*int64{&offset})
	assert.Nil(t, diags)
	schemaOperations := []*policyOperationModel{
		{
			ActionSetting: basetypes.NewStringValue(actionSetting),
			OperationType: basetypes.NewStringValue(operationType),
			BackupWindowTz: []*backupWindowModel{
				{
					StartTime: basetypes.NewStringValue(startTime),
					EndTime:   basetypes.NewStringValue(endTime),
				},
			},
			Slas: []*slaModel{
				{
					RetentionDuration: []*unitValueModel{
						{
							Unit:  basetypes.NewStringValue(retUnit),
							Value: basetypes.NewInt64Value(retValue),
						},
					},
					RPOFrequency: []*rpoModel{
						{
							Unit:    basetypes.NewStringValue(rpoUnit),
							Value:   basetypes.NewInt64Value(rpoValue),
							Offsets: offsets,
						},
					},
				},
				{
					RetentionDuration: []*unitValueModel{
						{
							Unit:  basetypes.NewStringValue(retUnit2),
							Value: basetypes.NewInt64Value(retValue2),
						},
					},
					RPOFrequency: []*rpoModel{
						{
							Unit:    basetypes.NewStringValue(rpoUnit2),
							Value:   basetypes.NewInt64Value(rpoValue2),
							Offsets: offsets,
						},
					},
				},
			},
			AdvancedSettings: []*advancedSettingsModel{
				{
					EBSVolumeBackup: []*backupTierModel{
						{
							BackupTier: basetypes.NewStringValue(backupTier),
						},
					},
				},
			},
			BackupAwsRegion: basetypes.NewStringValue(backupRegion),
		},
	}

	t.Run("Test EBS Volume Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				EBSVolumeBackup: []*backupTierModel{
					{
						BackupTier: basetypes.NewStringValue(backupTier),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.EBSVolumeBackup[0].BackupTier.ValueString(),
			*modelOpAdvSettings.AwsEbsVolumeBackup.BackupTier)
	})

	t.Run("Test EC2 Instance Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				EC2InstanceBackup: []*backupTierModel{
					{
						BackupTier: basetypes.NewStringValue(backupTier),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType2)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.EC2InstanceBackup[0].BackupTier.ValueString(),
			*modelOpAdvSettings.AwsEc2InstanceBackup.BackupTier)
	})

	t.Run("Test RDS Config Sync Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				RDSPitrConfigSync: []*pitrConfigModel{
					{
						Apply: basetypes.NewStringValue(rdsConfigSyncApply),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType3)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.RDSPitrConfigSync[0].Apply.ValueString(),
			*modelOpAdvSettings.AwsRdsConfigSync.Apply)
	})

	t.Run("Test RDS Granular Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				RDSLogicalBackup: []*backupTierModel{
					{
						BackupTier: basetypes.NewStringValue(backupTier),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType4)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.RDSLogicalBackup[0].BackupTier.ValueString(),
			*modelOpAdvSettings.AwsRdsResourceGranularBackup.BackupTier)
	})

	t.Run("Test EC2 MSSQL Database Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				EC2MssqlDatabaseBackup: []*replicaModel{
					{
						PreferredReplica:   basetypes.NewStringValue(preferredReplica),
						AlternativeReplica: basetypes.NewStringValue(alternativeReplica),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType5)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.EC2MssqlDatabaseBackup[0].PreferredReplica.ValueString(),
			*modelOpAdvSettings.Ec2MssqlDatabaseBackup.PreferredReplica)
		assert.Equal(t, schemaOpAdvSettings.EC2MssqlDatabaseBackup[0].AlternativeReplica.ValueString(),
			*modelOpAdvSettings.Ec2MssqlDatabaseBackup.AlternativeReplica)
	})

	t.Run("Test EC2 MSSQL Log Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				EC2MssqlLogBackup: []*replicaModel{
					{
						PreferredReplica:   basetypes.NewStringValue(preferredReplica),
						AlternativeReplica: basetypes.NewStringValue(alternativeReplica),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType6)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.EC2MssqlLogBackup[0].PreferredReplica.ValueString(),
			*modelOpAdvSettings.Ec2MssqlLogBackup.PreferredReplica)
		assert.Equal(t, schemaOpAdvSettings.EC2MssqlLogBackup[0].AlternativeReplica.ValueString(),
			*modelOpAdvSettings.Ec2MssqlLogBackup.AlternativeReplica)
	})

	t.Run("Test MSSQL Database Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				MssqlDatabaseBackup: []*replicaModel{
					{
						PreferredReplica:   basetypes.NewStringValue(preferredReplica),
						AlternativeReplica: basetypes.NewStringValue(alternativeReplica),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType7)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.MssqlDatabaseBackup[0].PreferredReplica.ValueString(),
			*modelOpAdvSettings.MssqlDatabaseBackup.PreferredReplica)
		assert.Equal(t, schemaOpAdvSettings.MssqlDatabaseBackup[0].AlternativeReplica.ValueString(),
			*modelOpAdvSettings.MssqlDatabaseBackup.AlternativeReplica)
	})

	t.Run("Test MSSQL Log Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				MssqlLogBackup: []*replicaModel{
					{
						PreferredReplica:   basetypes.NewStringValue(preferredReplica),
						AlternativeReplica: basetypes.NewStringValue(alternativeReplica),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType8)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.MssqlLogBackup[0].PreferredReplica.ValueString(),
			*modelOpAdvSettings.MssqlLogBackup.PreferredReplica)
		assert.Equal(t, schemaOpAdvSettings.MssqlLogBackup[0].AlternativeReplica.ValueString(),
			*modelOpAdvSettings.MssqlLogBackup.AlternativeReplica)
	})

	t.Run("Test Protection Group Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				ProtectionGroupBackup: []*backupTierModel{
					{
						BackupTier: basetypes.NewStringValue(backupTier),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType9)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.ProtectionGroupBackup[0].BackupTier.ValueString(),
			*modelOpAdvSettings.ProtectionGroupBackup.BackupTier)
	})

	t.Run("Test S3 Continuous Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				S3ContinuousBackup: []*ContinuousConfigModel{
					{
						DisableEventbridgeNotification: basetypes.NewBoolValue(
							disableEventbridgeNotification),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType10)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.S3ContinuousBackup[0].DisableEventbridgeNotification.ValueBool(),
			*modelOpAdvSettings.ProtectionGroupContinuousBackup.DisableEventbridgeNotification)
	})

	t.Run("Test Iceberg Backup Advanced Setting", func(t *testing.T) {
		schemaOperations[0].AdvancedSettings = []*advancedSettingsModel{
			{
				IcebergTableBackup: []*backupTierModel{
					{
						BackupTier: basetypes.NewStringValue(backupTier),
					},
				},
			},
		}
		schemaOperations[0].OperationType = basetypes.NewStringValue(operationType10)
		modelOperations, diags := mapSchemaOperationsToClumioOperations(ctx, schemaOperations)
		assert.Nil(t, diags)
		modelOpAdvSettings := modelOperations[0].AdvancedSettings
		schemaOpAdvSettings := schemaOperations[0].AdvancedSettings[0]
		assert.Equal(t, schemaOpAdvSettings.IcebergTableBackup[0].BackupTier.ValueString(),
			*modelOpAdvSettings.AwsIcebergTableBackup.BackupTier)
	})
}
