// Copyright 2023. Clumio, Inc.

package clumio_aws_connection

const (
	schemaId                   = "id"
	schemaAccountNativeId      = "account_native_id"
	schemaAwsRegion            = "aws_region"
	schemaDescription          = "description"
	schemaOrganizationalUnitId = "organizational_unit_id"
	schemaConnectionStatus     = "connection_status"
	schemaToken                = "token"
	schemaNamespace            = "namespace"
	schemaClumioAwsAccountId   = "clumio_aws_account_id"
	schemaClumioAwsRegion      = "clumio_aws_region"
	schemaExternalId           = "role_external_id"
	schemaDataPlaneAccountId   = "data_plane_account_id"

	awsEnvironment            = "aws_environment"
	statusConnected           = "connected"
	errorFmt                  = "Error: %v"
	externalIDFmt             = "ExternalID_%s"
	defaultDataPlaneAccountId = "*"

	http202           = 202
	pollTimeoutInSec  = 3600
	pollIntervalInSec = 5
)
