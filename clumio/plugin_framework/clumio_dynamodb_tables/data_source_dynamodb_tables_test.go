// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_dynamo_db_tables Terraform datasource. Please
// view the README.md file for more information on how to run these tests.

//go:build basic

package clumio_dynamodb_tables_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_dynamodb_tables"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceClumioDynamoDBTables(t *testing.T) {

	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	tableId := os.Getenv(clumio_dynamodb_tables.TableNativeId)
	if tableId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", clumio_dynamodb_tables.TableNativeId))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceClumioDynamoDBTables, baseUrl, accountNativeId,
					testAwsRegion, tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.clumio_dynamodb_tables.ds_dynamodb_tables",
						"dynamodb_tables.0.table_native_id",
						regexp.MustCompile(tableId)),
				),
			},
		},
	})
}

// Test of the clumio_dynamo_db_tables datasource with name specified as empty string.
func TestAccDataSourceClumioEmptyName(t *testing.T) {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceEmptyName, baseUrl),
				ExpectError: regexp.MustCompile(
					".*At least one of these attributes must be configured: \\[name,table_native_id\\].*"),
			},
		},
	})
}

// testAccDataSourceClumioDynamoDBTables is the Terraform configuration for a basic
// clumio_dynamo_db_tables data source.
const testAccDataSourceClumioDynamoDBTables = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_dynamodb_tables" "ds_dynamodb_tables" {
  account_native_id="%s"
  aws_region="%s"
  table_native_id="%s"
}
`

// testAccDataSourceEmptyName is the Terraform configuration for a
// clumio_dynamo_db_tables datasource with name set to empty string.
const testAccDataSourceEmptyName = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_dynamodb_tables" "ds_dynamodb_tables" {
  account_native_id="1234567890"
  aws_region="us-west-2"
}
`
