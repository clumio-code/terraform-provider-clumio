// Copyright 2023. Clumio, Inc.

package clumio_policy_assignment

const (
	// Constants used by the resource model for the clumio_policy_assignment Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaId         = "id"
	schemaEntityId   = "entity_id"
	schemaEntityType = "entity_type"
	schemaPolicyId   = "policy_id"

	entityTypeProtectionGroup  = "protection_group"
	entityTypeAWSDynamoDBTable = "aws_dynamodb_table"
	protectionGroupBackup      = "protection_group_backup"
	dynamodbTableBackup        = "aws_dynamodb_table_backup"

	//Common error messages used by the resource.
	readProtectionGroupErrFmt = "Unable to read Protection Group %v."
	readDynamoDBTableErrFmt   = "Unable to read DynamoDB table %v."
)

var (
	actionAssign   = "assign"
	actionUnassign = "unassign"
	policyIdEmpty  = ""
)
