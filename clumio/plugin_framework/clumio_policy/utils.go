// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_policy Terraform resource.

package clumio_policy

import (
	"context"
	"fmt"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkPolicyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// readPolicyAndUpdateModel calls the Clumio API to read the policy and convert the Clumio API
// response back to a schema and update the state. In addition to computed fields, all fields are
// populated from the API response in case any values have been changed externally. ID is not
// updated however given that it is the field used to query the resource from the backend.
func readPolicyAndUpdateModel(ctx context.Context,
	state *policyResourceModel, pd sdkPolicyDefinitions.PolicyDefinitionsV1Client) (
	*apiutils.APIError, diag.Diagnostics) {

	// Call the Clumio API to read the policy.
	res, apiErr := pd.ReadPolicyDefinition(state.ID.ValueString(), nil)
	if apiErr != nil {
		errMsg := fmt.Sprintf(
			"Error retrieving policy with ID: %s. Error: %v", state.ID.ValueString(), apiErr)
		tflog.Error(ctx, errMsg)
		return apiErr, nil
	}
	state.LockStatus = types.StringPointerValue(res.LockStatus)
	state.Name = types.StringPointerValue(res.Name)
	state.Timezone = types.StringPointerValue(res.Timezone)
	state.ActivationStatus = types.StringPointerValue(res.ActivationStatus)
	state.OrganizationalUnitId = types.StringPointerValue(res.OrganizationalUnitId)
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
		backupAwsRegionPtr := common.GetStringPtr(operation.BackupAwsRegion)

		var backupWindowTz *models.BackupWindow
		if operation.BackupWindowTz != nil {
			backupWindowTz = &models.BackupWindow{
				EndTime:   operation.BackupWindowTz[0].EndTime.ValueStringPointer(),
				StartTime: operation.BackupWindowTz[0].StartTime.ValueStringPointer(),
			}
		}

		advancedSettings := getOperationAdvancedSettings(operation)

		var backupSLAs []*models.BackupSLA
		if operation.Slas != nil {
			backupSLAs = make([]*models.BackupSLA, 0)

			for _, operationSla := range operation.Slas {
				backupSLA := &models.BackupSLA{}
				if operationSla.RetentionDuration != nil {
					backupSLA.RetentionDuration = &models.RetentionBackupSLAParam{
						Unit:  operationSla.RetentionDuration[0].Unit.ValueStringPointer(),
						Value: operationSla.RetentionDuration[0].Value.ValueInt64Pointer(),
					}
				}
				if operationSla.RPOFrequency != nil {
					var offsets []*int64
					diags = operationSla.RPOFrequency[0].Offsets.ElementsAs(ctx, &offsets, true)
					backupSLA.RpoFrequency = &models.RPOBackupSLAParam{
						Unit:    operationSla.RPOFrequency[0].Unit.ValueStringPointer(),
						Value:   operationSla.RPOFrequency[0].Value.ValueInt64Pointer(),
						Offsets: offsets,
					}
				}
				backupSLAs = append(backupSLAs, backupSLA)
			}
		}

		policyOperation := &models.PolicyOperationInput{
			ActionSetting:    operation.ActionSetting.ValueStringPointer(),
			BackupWindowTz:   backupWindowTz,
			Slas:             backupSLAs,
			ClumioType:       operation.OperationType.ValueStringPointer(),
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
		schemaOperation.ActionSetting = types.StringPointerValue(operation.ActionSetting)
		schemaOperation.OperationType = types.StringPointerValue(operation.ClumioType)

		if operation.BackupAwsRegion != nil {
			schemaOperation.BackupAwsRegion = types.StringPointerValue(operation.BackupAwsRegion)
		}

		if operation.BackupWindowTz != nil {
			window := &backupWindowModel{}
			window.StartTime = types.StringPointerValue(operation.BackupWindowTz.StartTime)
			window.EndTime = types.StringPointerValue(operation.BackupWindowTz.EndTime)
			schemaOperation.BackupWindowTz = []*backupWindowModel{window}
		}

		if operation.Slas != nil {
			backupSlas := make([]*slaModel, 0)
			for _, sla := range operation.Slas {
				backupSla := &slaModel{}
				if sla.RetentionDuration != nil {
					backupSla.RetentionDuration = []*unitValueModel{
						{
							Unit:  types.StringPointerValue(sla.RetentionDuration.Unit),
							Value: types.Int64PointerValue(sla.RetentionDuration.Value),
						},
					}
				}
				if sla.RpoFrequency != nil {
					offsets, rpoDiags := types.ListValueFrom(ctx,
						types.Int64Type, sla.RpoFrequency.Offsets)
					diags = rpoDiags
					backupSla.RPOFrequency = []*rpoModel{
						{
							Unit:    types.StringPointerValue(sla.RpoFrequency.Unit),
							Value:   types.Int64PointerValue(sla.RpoFrequency.Value),
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
						AlternativeReplica: types.StringPointerValue(
							operation.AdvancedSettings.Ec2MssqlDatabaseBackup.AlternativeReplica),
						PreferredReplica: types.StringPointerValue(
							operation.AdvancedSettings.Ec2MssqlDatabaseBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.Ec2MssqlLogBackup != nil {
				advSettings.EC2MssqlLogBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringPointerValue(
							operation.AdvancedSettings.Ec2MssqlLogBackup.AlternativeReplica),
						PreferredReplica: types.StringPointerValue(
							operation.AdvancedSettings.Ec2MssqlLogBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.MssqlDatabaseBackup != nil {
				advSettings.MssqlDatabaseBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringPointerValue(
							operation.AdvancedSettings.MssqlDatabaseBackup.AlternativeReplica),
						PreferredReplica: types.StringPointerValue(
							operation.AdvancedSettings.MssqlDatabaseBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.MssqlLogBackup != nil {
				advSettings.MssqlLogBackup = []*replicaModel{
					{
						AlternativeReplica: types.StringPointerValue(
							operation.AdvancedSettings.MssqlLogBackup.AlternativeReplica),
						PreferredReplica: types.StringPointerValue(
							operation.AdvancedSettings.MssqlLogBackup.PreferredReplica),
					},
				}
			}
			if operation.AdvancedSettings.ProtectionGroupBackup != nil {
				advSettings.ProtectionGroupBackup = []*backupTierModel{
					{
						BackupTier: types.StringPointerValue(
							operation.AdvancedSettings.ProtectionGroupBackup.BackupTier),
					},
				}
			}
			if operation.AdvancedSettings.AwsEbsVolumeBackup != nil {
				advSettings.EBSVolumeBackup = []*backupTierModel{
					{
						BackupTier: types.StringPointerValue(
							operation.AdvancedSettings.AwsEbsVolumeBackup.BackupTier),
					},
				}
			}
			if operation.AdvancedSettings.AwsEc2InstanceBackup != nil {
				advSettings.EC2InstanceBackup = []*backupTierModel{
					{
						BackupTier: types.StringPointerValue(
							operation.AdvancedSettings.AwsEc2InstanceBackup.BackupTier),
					},
				}
			}
			if operation.AdvancedSettings.AwsRdsConfigSync != nil {
				advSettings.RDSPitrConfigSync = []*pitrConfigModel{
					{
						Apply: types.StringPointerValue(
							operation.AdvancedSettings.AwsRdsConfigSync.Apply),
					},
				}
			}
			if operation.AdvancedSettings.AwsRdsResourceGranularBackup != nil {
				advSettings.RDSLogicalBackup = []*backupTierModel{
					{
						BackupTier: types.StringPointerValue(
							operation.AdvancedSettings.AwsRdsResourceGranularBackup.BackupTier),
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
			advancedSettings.AwsEbsVolumeBackup = &models.EBSBackupAdvancedSetting{
				BackupTier: operation.AdvancedSettings[0].EBSVolumeBackup[0].BackupTier.
					ValueStringPointer(),
			}
		}
		if operation.AdvancedSettings[0].EC2InstanceBackup != nil {
			advancedSettings.AwsEc2InstanceBackup = &models.EC2BackupAdvancedSetting{
				BackupTier: operation.AdvancedSettings[0].EC2InstanceBackup[0].
					BackupTier.ValueStringPointer(),
			}
		}
		if operation.AdvancedSettings[0].ProtectionGroupBackup != nil {
			advancedSettings.ProtectionGroupBackup =
				&models.ProtectionGroupBackupAdvancedSetting{
					BackupTier: operation.AdvancedSettings[0].ProtectionGroupBackup[0].
						BackupTier.ValueStringPointer(),
				}
		}
		if operation.AdvancedSettings[0].EC2MssqlDatabaseBackup != nil {
			advancedSettings.Ec2MssqlDatabaseBackup =
				&models.EC2MSSQLDatabaseBackupAdvancedSetting{
					AlternativeReplica: operation.AdvancedSettings[0].EC2MssqlDatabaseBackup[0].
						AlternativeReplica.ValueStringPointer(),
					PreferredReplica: operation.AdvancedSettings[0].EC2MssqlDatabaseBackup[0].
						PreferredReplica.ValueStringPointer(),
				}
		}
		if operation.AdvancedSettings[0].EC2MssqlLogBackup != nil {
			advancedSettings.Ec2MssqlLogBackup =
				&models.EC2MSSQLLogBackupAdvancedSetting{
					AlternativeReplica: operation.AdvancedSettings[0].EC2MssqlLogBackup[0].
						AlternativeReplica.ValueStringPointer(),
					PreferredReplica: operation.AdvancedSettings[0].EC2MssqlLogBackup[0].
						PreferredReplica.ValueStringPointer(),
				}
		}
		if operation.AdvancedSettings[0].MssqlDatabaseBackup != nil {
			advancedSettings.MssqlDatabaseBackup =
				&models.MSSQLDatabaseBackupAdvancedSetting{
					AlternativeReplica: operation.AdvancedSettings[0].MssqlDatabaseBackup[0].
						AlternativeReplica.ValueStringPointer(),
					PreferredReplica: operation.AdvancedSettings[0].MssqlDatabaseBackup[0].
						PreferredReplica.ValueStringPointer(),
				}
		}
		if operation.AdvancedSettings[0].MssqlLogBackup != nil {
			advancedSettings.MssqlLogBackup =
				&models.MSSQLLogBackupAdvancedSetting{
					AlternativeReplica: operation.AdvancedSettings[0].MssqlLogBackup[0].
						AlternativeReplica.ValueStringPointer(),
					PreferredReplica: operation.AdvancedSettings[0].MssqlLogBackup[0].
						PreferredReplica.ValueStringPointer(),
				}
		}
		if operation.AdvancedSettings[0].RDSPitrConfigSync != nil {
			advancedSettings.AwsRdsConfigSync =
				&models.RDSConfigSyncAdvancedSetting{
					Apply: operation.AdvancedSettings[0].RDSPitrConfigSync[0].Apply.
						ValueStringPointer(),
				}
		}
		if operation.AdvancedSettings[0].RDSLogicalBackup != nil {
			advancedSettings.AwsRdsResourceGranularBackup =
				&models.RDSLogicalBackupAdvancedSetting{
					BackupTier: operation.AdvancedSettings[0].RDSLogicalBackup[0].
						BackupTier.ValueStringPointer(),
				}
		}
	}
	return advancedSettings
}
