// Copyright 2026. Clumio, Inc.

// This file holds acceptance tests for the clumio_aws_environment Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_aws_environment_test

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

// TestAccDataSourceClumioAWSEnvironment is a basic test of the clumio_aws_environment datasource.
// It verifies that the environment matching the account and region in the config is fetched and its
// id is set, and that a non-existent account returns an error.
//
// NOTE: an AWS environment only exists once the corresponding AWS connection is connected. The test
// is gated on TABLE_NATIVE_ID being set, which signals that a connected test account (which has
// DynamoDB data) is available.
func TestAccDataSourceClumioAWSEnvironment(t *testing.T) {

	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	if os.Getenv(clumio_dynamodb_tables.TableNativeId) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance test skipped unless env '%s' set", clumio_dynamodb_tables.TableNativeId))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceClumioAWSEnvironment, baseUrl, accountNativeId, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.clumio_aws_environment.ds_env", "id"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccDataSourceClumioAWSEnvironment, baseUrl, "12345678901", testAwsRegion),
				ExpectError: regexp.MustCompile(
					".*unable to retrieve environment corresponding to.*"),
			},
		},
	})
}

// testAccDataSourceClumioAWSEnvironment is the Terraform configuration for a basic
// clumio_aws_environment datasource.
const testAccDataSourceClumioAWSEnvironment = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_aws_environment" "ds_env" {
  account_native_id = "%s"
  aws_region        = "%s"
}
`
