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
)
