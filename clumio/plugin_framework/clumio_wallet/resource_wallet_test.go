// Copyright 2023. Clumio, Inc.
//
// Acceptance test for resource_wallet.

package clumio_wallet_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWallet(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceWallet(baseUrl, accountNativeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_wallet.test_wallet", "account_native_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
		},
	})
}

func getTestAccResourceWallet(baseUrl string, accountId string) string {
	return fmt.Sprintf(testAccResourcePostWallet, baseUrl, accountId)
}

const testAccResourcePostWallet = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_wallet" "test_wallet" {
  account_native_id = "%s"
}
`
