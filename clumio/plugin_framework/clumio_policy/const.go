// Copyright 2023. Clumio, Inc.

package clumio_policy

const (
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
	schemaEc2MssqlDatabaseBackup = "ec2_mssql_database_backup"
	schemaEc2MssqlLogBackup      = "ec2_mssql_log_backup"
	schemaMssqlDatabaseBackup    = "mssql_database_backup"
	schemaMssqlLogBackup         = "mssql_log_backup"
	schemaProtectionGroupBackup  = "protection_group_backup"
	schemaBackupTier             = "backup_tier"
	schemaEBSVolumeBackup        = "aws_ebs_volume_backup"
	schemaEC2InstanceBackup      = "aws_ec2_instance_backup"

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

	eC2BackupDesc = "Optional configuration settings for the aws_ec2_instance_backup operation."

	secureVaultLiteDescFmt = "Backup tier to store the SecureVault Lite backup in." +
		" Valid values are: `standard` and `lite`. If not provided, the default is `standard`."

	errorFmt           = "Error: %v"
	errorPolicyReadMsg = "Error retrieving Clumio Policy."

	timeoutInSec  = 3600
	intervalInSec = 5
)