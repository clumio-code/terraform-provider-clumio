// Copyright 2023. Clumio, Inc.

package clumio_aws_manual_connection_resources

const (
	// Constants used by the resource model for the clumio_aws_connection Terraform resource. These
	// values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId = "id"
	schemaAccountNativeId = "account_native_id"
	schemaAwsRegion = "aws_region"
	schemaAssetTypesEnabled = "asset_types_enabled"
	schemaResources = "resources"

	schemaIsEbsEnabled = "ebs"
	schemaIsRDSEnabled = "rds" 
	schemaIsDynamoDBEnabled = "ddb"
	schemaIsS3Enabled = "s3"
	schemaIsMssqlEnabled = "mssql"

	EBS      = "EBS"
	S3       = "S3"
	DynamoDB = "DynamoDB"
	RDS      = "RDS"
	EC2MSSQL = "EC2MSSQL"
)
