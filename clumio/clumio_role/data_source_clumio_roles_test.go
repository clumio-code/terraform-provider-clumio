// Copyright 2021. Clumio, Inc.

// Acceptance test for Data Source clumio_role.
package clumio_role_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceClumioRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceClumioRoles, os.Getenv(common.ClumioApiBaseUrl)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.clumio_role.test", "id", "00000000-0000-0000-0000-000000000000"),
				),
			},
		},
	})
}

const testAccDataSourceClumioRoles = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_role" "test" {
	name = "Super Admin"
}
`
