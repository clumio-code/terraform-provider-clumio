// Copyright 2021. Clumio, Inc.

package common

const (
	ClumioApiToken                  = "CLUMIO_API_TOKEN"
	ClumioApiBaseUrl                = "CLUMIO_API_BASE_URL"
	ClumioOrganizationalUnitContext = "CLUMIO_ORGANIZATIONAL_UNIT_CONTEXT"
	AwsAccessKeyId                  = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey              = "AWS_SECRET_ACCESS_KEY"
	AwsRegion                       = "AWS_REGION"
	ClumioTestAwsAccountId          = "CLUMIO_TEST_AWS_ACCOUNT_ID"

	TaskSuccess = "completed"
	TaskAborted = "aborted"
	TaskFailed  = "failed"

	ProtectService       = "protect"
	OrganizationalUnitId = "organizational_unit_id"

	// Set attribute error format
	SchemaAttributeSetError = "Error setting %s schema attribute. Error: %v"

	// Empty parameter error format
	SchemaEmptyParameterError = "Empty %s is invalid if specified."

	// AWS Manual Connection Resources
	ClumioIAMRoleArn = "clumio_iam_role_arn"
	ClumioEventPubArn = "clumio_event_pub_arn"
	ClumioSupportRoleArn = "clumio_support_role_arn"
	CloudwatchRuleArn = "cloudwatch_rule_arn"
	CloudtrailRuleArn = "cloudtrail_rule_arn"
	ContinuousBackupsRoleArn = "continuous_backups_role_arn"
	SsmNotificationRoleArn = "ssm_notification_role_arn"
	Ec2SsmInstanceProfileArn = "ec2_ssm_instance_profile_arn"
)
