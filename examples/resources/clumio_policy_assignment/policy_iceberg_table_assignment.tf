data "clumio_iceberg_tables" "ds_iceberg_tables" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
  name              = "Iceberg table name"
  catalog_type      = "aws_iceberg_glue_table"
}

resource "clumio_policy_assignment" "example" {
  entity_id   = data.clumio_iceberg_tables.ds_iceberg_tables.iceberg_tables[0].id
  entity_type = "aws_iceberg_glue_table"
  policy_id   = "policy_id"
}
