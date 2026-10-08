// Copyright 2023. Clumio, Inc.

package clumio_policy_assignment

const (
	// Constants used by the resource model for the clumio_policy_assignment Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId         = "id"
	schemaEntityId   = "entity_id"
	schemaEntityType = "entity_type"
	schemaPolicyId   = "policy_id"

	entityTypeProtectionGroup    = "protection_group"
	entityTypeGcpProtectionGroup = "gcp_protection_group"
	entityTypeAWSDynamoDBTable   = "aws_dynamodb_table"
	entityTypeIcebergGlueTable   = "aws_iceberg_glue_table"
	entityTypeIcebergS3Table     = "aws_iceberg_s3_table"
	protectionGroupBackup        = "protection_group_backup"
	gcpProtectionGroupBackup     = "gcp_protection_group_backup"
	dynamodbTableBackup          = "aws_dynamodb_table_backup"
)

var (
	actionAssign   = "assign"
	actionUnassign = "unassign"
	policyIdEmpty  = ""
)
