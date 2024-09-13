// Copyright 2023. Clumio, Inc.

package clumio_aws_connection

const (
	// Constants used by the resource model for the clumio_aws_connection Terraform resource. These
	// values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId                 = "id"
	schemaAccountNativeId    = "account_native_id"
	schemaAwsRegion          = "aws_region"
	schemaDescription        = "description"
	schemaConnectionStatus   = "connection_status"
	schemaToken              = "token"
	schemaNamespace          = "namespace"
	schemaClumioAwsAccountId = "clumio_aws_account_id"
	schemaClumioAwsRegion    = "clumio_aws_region"
	schemaExternalId         = "role_external_id"
	schemaDataPlaneAccountId = "data_plane_account_id"

	awsEnvironment            = "aws_environment"
	statusConnected           = "connected"
	externalIDFmt             = "ExternalID_%s"
	defaultDataPlaneAccountId = "*"
	defaultOrgUnitId          = "00000000-0000-0000-0000-000000000000"
)
