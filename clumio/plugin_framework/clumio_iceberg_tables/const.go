// Copyright 2026. Clumio, Inc.

package clumio_iceberg_tables

const (
	// Constants used by the datasource model for the clumio_iceberg_tables Terraform datasource.
	// These values should match the schema tfsdk tags on the datasource model struct in
	// data_source_schema.go.
	schemaId              = "id"
	schemaName            = "name"
	schemaRegion          = "aws_region"
	schemaAccountNativeId = "account_native_id"
	schemaCatalog         = "catalog"
	schemaCatalogType     = "catalog_type"
	schemaNamespace       = "namespace"
	schemaIcebergTables   = "iceberg_tables"

	catalogTypeGlue    = "aws_iceberg_glue_table"
	catalogTypeS3Table = "aws_iceberg_s3_table"

	// The list Iceberg tables API caps the page size at 100.
	listLimit = int64(100)

	// IcebergTableName is the environment variable holding the name of the Iceberg table used by
	// the acceptance tests.
	IcebergTableName = "ICEBERG_TABLE_NAME"
)
