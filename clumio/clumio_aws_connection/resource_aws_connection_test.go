// Copyright 2021. Clumio, Inc.

// Acceptance test for resource_aws_connection.
package clumio_aws_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioAwsConnection(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
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
}
`
