// Copyright 2021. Clumio, Inc.

package common

const (
	ClumioApiToken                  = "CLUMIO_API_TOKEN"
	ClumioApiBaseUrl                = "CLUMIO_API_BASE_URL"
	ClumioOrganizationalUnitContext = "CLUMIO_ORGANIZATIONAL_UNIT_CONTEXT"
	AwsRegion                       = "AWS_REGION"
	ClumioTestAwsAccountId          = "CLUMIO_TEST_AWS_ACCOUNT_ID"
	ClumioTestAwsAccountId2         = "CLUMIO_TEST_AWS_ACCOUNT_ID2"
	ClumioTestIsSSOConfigured       = "CLUMIO_TEST_IS_SSO_CONFIGURED"

	TaskSuccess    = "completed"
	TaskAborted    = "aborted"
	TaskFailed     = "failed"
	TaskInProgress = "in_progress"

	// AWS Manual Connection Resources
	ClumioIAMRoleArn         = "clumio_iam_role_arn"
	ClumioEventPubArn        = "clumio_event_pub_arn"
	ClumioSupportRoleArn     = "clumio_support_role_arn"
	CloudwatchRuleArn        = "cloudwatch_rule_arn"
	CloudtrailRuleArn        = "cloudtrail_rule_arn"
	ContinuousBackupsRoleArn = "continuous_backups_role_arn"
	SsmNotificationRoleArn   = "ssm_notification_role_arn"
	Ec2SsmInstanceProfileArn = "ec2_ssm_instance_profile_arn"

	// Default Error Message
	NilErrorMessageSummary = "Unexpected API response"
	NilErrorMessageDetail  = "An empty response was returned by the API"

	// Testing error format
	TestResultsNotMatchingError = "Results don't match.\nExpected: %v\nActual: %v"

	// Auth Error
	AuthError = "Unauthorized access. Please ensure that your credentials are valid and/or your account is active."
)
