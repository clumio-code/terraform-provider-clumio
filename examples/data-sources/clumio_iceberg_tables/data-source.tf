# Example using the table name
data "clumio_iceberg_tables" "ds_iceberg_tables" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
  name              = "Iceberg table name"
}

# Example using the catalog, namespace and catalog type
data "clumio_iceberg_tables" "ds_iceberg_tables" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
  catalog           = "Iceberg catalog name"
  namespace         = "Iceberg namespace"
  catalog_type      = "aws_iceberg_glue_table"
}
