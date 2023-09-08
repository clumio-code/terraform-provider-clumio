// Copyright 2023. Clumio, Inc.
//
// Acceptance test for data_source_role.

package clumio_role_test

import (
	"fmt"
	"os"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceClumioRoles(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioCallbackClumioRole(baseUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.clumio_role.test_role", "id", "30000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr(
						"data.clumio_role.test_role", "name", "Backup Admin"),
				),
			},
		},
	})
}

func getTestDataSourceClumioCallbackClumioRole(baseUrl string) string {
	return fmt.Sprintf(testAccDataSourceClumioRoles, baseUrl)
}

const testAccDataSourceClumioRoles = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_role" "test_role" {
	name = "Backup Admin"
}
`
