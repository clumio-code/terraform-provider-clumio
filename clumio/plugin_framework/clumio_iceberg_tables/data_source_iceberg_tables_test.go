// Copyright 2026. Clumio, Inc.

// This file holds acceptance tests for the clumio_iceberg_tables Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_iceberg_tables_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_iceberg_tables"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test of the clumio_iceberg_tables datasource queried by the table name.
func TestAccDataSourceClumioIcebergTables(t *testing.T) {

	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	tableName := os.Getenv(clumio_iceberg_tables.IcebergTableName)
	if tableName == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", clumio_iceberg_tables.IcebergTableName))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceClumioIcebergTables, baseUrl, accountNativeId,
					testAwsRegion, tableName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.clumio_iceberg_tables.ds_iceberg_tables",
						"iceberg_tables.0.name",
						regexp.MustCompile(tableName)),
					resource.TestMatchResourceAttr(
						"data.clumio_iceberg_tables.ds_iceberg_tables",
						"iceberg_tables.0.id",
						regexp.MustCompile(".+")),
				),
			},
		},
	})
}

// Test of the clumio_iceberg_tables datasource without any of the table filters.
func TestAccDataSourceClumioIcebergTablesNoFilter(t *testing.T) {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceIcebergTablesNoFilter, baseUrl),
				ExpectError: regexp.MustCompile(".*At least one of these attributes must be configured.*"),
			},
		},
	})
}

// Test of the clumio_iceberg_tables datasource with an unsupported catalog type.
func TestAccDataSourceClumioIcebergTablesInvalidCatalogType(t *testing.T) {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceIcebergTablesInvalidCatalogType, baseUrl),
				ExpectError: regexp.MustCompile(".*value must be one of.*"),
			},
		},
	})
}

// testAccDataSourceClumioIcebergTables is the Terraform configuration for a basic
// clumio_iceberg_tables data source.
const testAccDataSourceClumioIcebergTables = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_iceberg_tables" "ds_iceberg_tables" {
  account_native_id="%s"
  aws_region="%s"
  name="%s"
}
`

// testAccDataSourceIcebergTablesNoFilter is the Terraform configuration for a
// clumio_iceberg_tables datasource without a name, catalog or namespace.
const testAccDataSourceIcebergTablesNoFilter = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_iceberg_tables" "ds_iceberg_tables" {
  account_native_id="1234567890"
  aws_region="us-west-2"
}
`

// testAccDataSourceIcebergTablesInvalidCatalogType is the Terraform configuration for a
// clumio_iceberg_tables datasource with an unsupported catalog type.
const testAccDataSourceIcebergTablesInvalidCatalogType = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_iceberg_tables" "ds_iceberg_tables" {
  account_native_id="1234567890"
  aws_region="us-west-2"
  name="test-table"
  catalog_type="aws_iceberg_hive_table"
}
`
