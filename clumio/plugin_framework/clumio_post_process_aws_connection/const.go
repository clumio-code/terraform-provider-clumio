// Copyright 2023. Clumio, Inc.

package clumio_post_process_aws_connection

const (
	// Constants used by the resource model for the clumio_post_process_aws_connection Terraform
	// resource. These values should match the schema tfsdk tags on the resource model struct in
	// schema.go.
	schemaId                              = "id"
	schemaToken                           = "token"
	schemaRoleExternalId                  = "role_external_id"
	schemaAccountId                       = "account_id"
	schemaRegion                          = "region"
	schemaRoleArn                         = "role_arn"
	schemaConfigVersion                   = "config_version"
	schemaDiscoverVersion                 = "discover_version"
	schemaProtectConfigVersion            = "protect_config_version"
	schemaProtectEbsVersion               = "protect_ebs_version"
	schemaProtectRdsVersion               = "protect_rds_version"
	schemaProtectS3Version                = "protect_s3_version"
	schemaProtectDynamodbVersion          = "protect_dynamodb_version"
	schemaProtectWarmTierVersion          = "protect_warm_tier_version"
	schemaProtectWarmTierDynamodbVersion  = "protect_warm_tier_dynamodb_version"
	schemaProtectEc2MssqlVersion          = "protect_ec2_mssql_version"
	schemaProtectIcebergOnGlueVersion     = "protect_iceberg_on_glue_version"
	schemaProtectIcebergOnS3TablesVersion = "protect_iceberg_on_s3_tables_version"
	schemaClumioEventPubId                = "clumio_event_pub_id"
	schemaProperties                      = "properties"
	schemaIntermediateRoleArn             = "intermediate_role_arn"
	schemaWaitForIngestion                = "wait_for_ingestion"
	schemaWaitForDataPlaneResources       = "wait_for_data_plane_resources"

	eventTypeCreate = "Create"
	eventTypeUpdate = "Update"

	inProgress = "in_progress"
	completed  = "completed"
	failed     = "failed"

	//Connected connection status
	connected = "connected"
)
