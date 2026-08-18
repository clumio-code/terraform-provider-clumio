// Copyright 2026. Clumio, Inc.

// This file holds acceptance tests for the clumio_dynamodb_backups Terraform datasource. Please
// view the README.md file for more information on how to run these tests.

//go:build basic

package clumio_dynamodb_backups_test

import (
	"fmt"
	"os"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_dynamodb_tables"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDataSourceClumioDynamoDBBackups is a basic test of the clumio_dynamodb_backups datasource.
// It resolves a DynamoDB table by its native id via the clumio_dynamodb_tables datasource, then
// verifies that the backups datasource returns at least one backup (the latest at index 0).
//
// Requires TABLE_NATIVE_ID to point to a table that has at least one Clumio backup; the test is
// skipped otherwise.
func TestAccDataSourceClumioDynamoDBBackups(t *testing.T) {

	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	tableNativeId := os.Getenv(clumio_dynamodb_tables.TableNativeId)
	if tableNativeId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance test skipped unless env '%s' set", clumio_dynamodb_tables.TableNativeId))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceClumioDynamoDBBackups,
					baseUrl, accountNativeId, testAwsRegion, tableNativeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.clumio_dynamodb_backups.ds_backups", "backups.0.id"),
				),
			},
		},
	})
}

// testAccDataSourceClumioDynamoDBBackups is the Terraform configuration for a basic
// clumio_dynamodb_backups datasource. It chains through the clumio_dynamodb_tables datasource to
// resolve the Clumio-assigned table_id from the table native id.
const testAccDataSourceClumioDynamoDBBackups = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_dynamodb_tables" "ds_tables" {
  account_native_id = "%s"
  aws_region        = "%s"
  table_native_id   = "%s"
}

data "clumio_dynamodb_backups" "ds_backups" {
  table_id = data.clumio_dynamodb_tables.ds_tables.dynamodb_tables[0].id
}
`
