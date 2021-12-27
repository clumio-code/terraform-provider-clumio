// Copyright 2021. Clumio, Inc.

// Acceptance test for resource_aws_connection.
package clumio

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioAwsConnection(t *testing.T) {
	accountNativeId := os.Getenv(clumioTestAwsAccountId)
	baseUrl := os.Getenv(clumioApiBaseUrl)
	testAwsRegion := os.Getenv(awsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckClumio(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioCallbackAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "account_native_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
			{
				Config: getTestAccResourceClumioCallbackAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "account_native_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
		},
	})
}

func getTestAccResourceClumioCallbackAwsConnection(
	baseUrl string, accountId string, awsRegion string, description string) string {
	return fmt.Sprintf(testAccResourceClumioAwsConnection, baseUrl, accountId,
		awsRegion, description)
}

const testAccResourceClumioAwsConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
  account_native_id = "%s"
  aws_region = "%s"
  description = "%s"
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  services_enabled = ["discover", "protect"]
}
`
