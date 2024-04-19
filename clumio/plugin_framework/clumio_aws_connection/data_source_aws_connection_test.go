// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_aws_connection Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_aws_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_aws_connection datasource. It tests that the aws_connection matching
// the name provided in the config is fetched and its Id is set in the state.
func TestAccDataSourceClumioAWSConnection(t *testing.T) {

	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioAWSConnection(
					baseUrl, accountNativeId, testAwsRegion, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.clumio_aws_connection.ds_conn",
						"id"),
				),
			},
			{
				Config: getTestDataSourceClumioAWSConnection(
					baseUrl, accountNativeId, testAwsRegion, true),
				ExpectError: regexp.MustCompile(".*AWS connection not found.*"),
			},
		},
	})
}

// getTestDataSourceClumioAWSConnection returns the Terraform configuration for a basic
// clumio_aws_connection datasource.
func getTestDataSourceClumioAWSConnection(
	baseUrl string, accountId string, awsRegion string, isInvalid bool) string {

	dsAccountId := accountId
	if isInvalid {
		dsAccountId = "12345678901"
	}
	return fmt.Sprintf(
		testAccDataSourceClumioAWSConnection, baseUrl, accountId, awsRegion, dsAccountId, awsRegion)
}

// testAccDataSourceClumioAWSConnection is the Terraform configuration for a basic
// clumio_aws_connection datasource.
const testAccDataSourceClumioAWSConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "ds_test_conn" {
  account_native_id = "%s"
  aws_region = "%s"
  description = "some description"
}

data "clumio_aws_connection" "ds_conn" {
	depends_on = [ clumio_aws_connection.ds_test_conn]
	account_native_id = "%s"
    aws_region = "%s"
}
`
