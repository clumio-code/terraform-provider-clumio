// Copyright 2023. Clumio, Inc.

package clumio_aws_manual_connection

const (
	// Constants used by the resource model for the clumio_aws_manual_connection Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId                       = "id"
	schemaAccountId                = "account_id"
	schemaAwsRegion                = "aws_region"
	schemaAssetsEnabled            = "assets_enabled"
	schemaResources                = "resources"
	schemaClumioIAMRoleArn         = "clumio_iam_role_arn"
	schemaClumioEventPubArn        = "clumio_event_pub_arn"
	schemaClumioSupportRoleArn     = "clumio_support_role_arn"
	schemaEventRules               = "event_rules"
	schemaCloudwatchRuleArn        = "cloudwatch_rule_arn"
	schemaCloudtrailRuleArn        = "cloudtrail_rule_arn"
	schemaServiceRoles             = "service_roles"
	schemaS3                       = "s3"
	schemaMssql                    = "mssql"
	schemaContinuousBackupsRoleArn = "continuous_backups_role_arn"
	schemaSsmNotificationRoleArn   = "ssm_notification_role_arn"
	schemaEc2SsmInstanceProfileArn = "ec2_ssm_instance_profile_arn"
)
