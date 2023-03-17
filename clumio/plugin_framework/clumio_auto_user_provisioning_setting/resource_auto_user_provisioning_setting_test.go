// Copyright 2023. Clumio, Inc.

// Acceptance test for clumio_auto_user_provisioning_setting resource.
package clumio_auto_user_provisioning_setting_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClumioAutoUserProvisioningSetting(t *testing.T) {
	enabled := "true"
	disabled := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningSetting(disabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting", "is_enabled",
						regexp.MustCompile(disabled)),
				),
			},
			{
				Config: getTestAccResourceClumioAutoUserProvisioningSetting(enabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting", "is_enabled",
						regexp.MustCompile(enabled)),
				),
			},
		},
	})
}

func getTestAccResourceClumioAutoUserProvisioningSetting(isEnabled string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioAutoUserProvisioningSetting, baseUrl, isEnabled)
}

const testAccResourceClumioAutoUserProvisioningSetting = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_auto_user_provisioning_setting" "test_auto_user_provisioning_setting" {
  is_enabled = %s
}

`
