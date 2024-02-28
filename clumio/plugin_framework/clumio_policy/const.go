// Copyright 2023. Clumio, Inc.

package clumio_policy

const (
	// Constants used by the resource model for the clumio_policy Terraform resource. These values
	// should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaActivationStatus       = "activation_status"
	schemaName                   = "name"
	schemaTimezone               = "timezone"
	schemaOperations             = "operations"
	schemaOrganizationalUnitId   = "organizational_unit_id"
	schemaActionSetting          = "action_setting"
	schemaOperationType          = "type"
	schemaBackupWindowTz         = "backup_window_tz"
	schemaSlas                   = "slas"
	schemaStartTime              = "start_time"
	schemaEndTime                = "end_time"
	schemaRetentionDuration      = "retention_duration"
	schemaRpoFrequency           = "rpo_frequency"
	schemaUnit                   = "unit"
	schemaValue                  = "value"
	schemaOffsets                = "offsets"
	schemaId                     = "id"
	schemaLockStatus             = "lock_status"
	schemaAdvancedSettings       = "advanced_settings"
	schemaAlternativeReplica     = "alternative_replica"
	schemaPreferredReplica       = "preferred_replica"
	schemaBackupAwsRegion        = "backup_aws_region"
	schemaEc2MssqlDatabaseBackup = "ec2_mssql_database_backup"
	schemaEc2MssqlLogBackup      = "ec2_mssql_log_backup"
	schemaMssqlDatabaseBackup    = "mssql_database_backup"
	schemaMssqlLogBackup         = "mssql_log_backup"
	schemaProtectionGroupBackup  = "protection_group_backup"
	schemaBackupTier             = "backup_tier"
	schemaEBSVolumeBackup        = "aws_ebs_volume_backup"
	schemaEC2InstanceBackup      = "aws_ec2_instance_backup"
	schemaRDSPitrConfigSync      = "aws_rds_config_sync"
	schemaApply                  = "apply"
	schemaRdsLogicalBackup       = "aws_rds_resource_granular_backup"

	alternativeReplicaDescFmt = "The alternative replica for MSSQL %s backups. This" +
		" setting only applies to Availability Group databases. Possible" +
		" values include \"primary\", \"sync_secondary\", and \"stop\"." +
		" If \"stop\" is provided, then backups will not attempt to switch" +
		" to a different replica when the preferred replica is unavailable." +
		" Otherwise, recurring backups will attempt to use either" +
		" the primary replica or the secondary replica accordingly."

	preferredReplicaDescFmt = "The primary preferred replica for MSSQL %s backups." +
		" This setting only applies to Availability Group databases." +
		" Possible values include \"primary\" and \"sync_secondary\"." +
		" Recurring backup will first attempt to use either the primary" +
		" replica or the secondary replica accordingly."

	mssqlDatabaseBackupDesc = "Additional policy configuration settings for the" +
		" mssql_database_backup operation. If this operation is not of" +
		" type mssql_database_backup, then this field is omitted from the" +
		" response."

	mssqlLogBackupDesc = "Additional policy configuration settings for the" +
		" mssql_log_backup operation. If this operation is not of" +
		" type mssql_log_backup, then this field is omitted from the" +
		" response."

	ebsBackupDesc = "Optional configuration settings for the aws_ebs_volume_backup operation."

	ec2BackupDesc = "Optional configuration settings for the aws_ec2_instance_backup operation."

	secureVaultLiteDesc = "Backup tier to store the SecureVault Lite backup in." +
		" Valid values are: `standard` and `lite`. If not provided, the default is `standard`."

	rdsPitrConfigSyncDesc = "Optional configuration settings for the aws_rds_config_sync operation."

	pitrConfigDesc = "Additional policy configuration for syncing the configuration of Pitr in aws." +
		" Possible values include \"immediate\" and \"maintenance_window\"." +
		" If \"immediate\" is provided, then configuration sync will be kicked in immediately." +
		" Otherwise configuration sync will be executed in a specific time user has provided."

	rdsLogicalBackupDesc = "Optional configuration settings for the aws_rds_resource_granular_backup operation."

	rdsLogicalBackupAdvancedSettingDesc = "Backup tier to store the RDS backup in. Valid values" +
		" are: `standard` (for Granular Record Retrieval) and `frozen` (for SecureVault Archive)." +
		" To update existing policies with RDS Granular Record Retrieval, the default is `standard`" +
		" if backup_tier is not provided. To update existing policies that do not have RDS Granular" +
		" Record Retrieval, or to create new policies, the only supported option is `frozen`."

	errorPolicyReadMsg = "Unable to read %s (ID: %v)"

	timeoutInSec  = 3600
	intervalInSec = 5

	// Constants for activation status allowed values
	activationStatusActivated  = "activated"
	activationStatusDectivated = "deactivated"
)
