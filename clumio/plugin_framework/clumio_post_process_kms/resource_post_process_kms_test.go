// Copyright 2023. Clumio, Inc.
//
// Acceptance test for resource_post_rprocess_aws_connection.

package clumio_post_process_kms_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePostProcessKMS(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourcePostProcessKMS(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_kms.test", "account_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
			{
				Config: getTestAccResourcePostProcessAwsConnectionUpdate(
					baseUrl, accountNativeId, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_kms.test", "account_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
		},
	})
}

func getTestAccResourcePostProcessKMS(
	baseUrl string, accountId string, awsRegion string, description string) string {
	return fmt.Sprintf(testAccResourcePostProcessAwsConnection, baseUrl, accountId,
		awsRegion)
}

const testAccResourcePostProcessAwsConnection = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_wallet" "test_wallet" {
  account_native_id = "%s"
}

resource "clumio_post_process_kms" "test" {
  token = clumio_wallet.test_wallet.token
  account_id = clumio_wallet.test_wallet.account_native_id
  region = "%s"
  multi_region_cmk_key_id = "test_multi_region_cmk_key_id"
  role_external_id = "test_role_external_id"
  role_arn = "test_role_arn"
  role_id = "test_role_id"
}
`

func getTestAccResourcePostProcessAwsConnectionUpdate(
	baseUrl string, accountId string, awsRegion string) string {
	return fmt.Sprintf(testAccResourcePostProcessAwsConnectionUpdate, baseUrl, accountId,
		awsRegion)
}

const testAccResourcePostProcessAwsConnectionUpdate = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_wallet" "test_wallet {
  account_native_id = %s
}

resource "clumio_post_process_aws_connection" "test" {
  token = clumio_wallet.test_wallet.token
  account_native_id = clumio_wallet.test_wallet.account_native_id
  region = "%s"
  multi_region_cmk_key_id = "test_multi_region_cmk_key_id"
  role_external_id = "test_role_external_id"
  role_arn = "test_role_arn"
  role_id = "test_role_id"
}
`
