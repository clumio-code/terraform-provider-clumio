// Copyright 2021. Clumio, Inc.

// Acceptance test for clumio_user resource.
package clumio

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

)

func TestAccResourceClumioUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckClumio(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioUser(false),
			},
			{
				Config: getTestAccResourceClumioUser(true),
			},
		},
	})
}

func getTestAccResourceClumioUser(update bool) string {
	baseUrl := os.Getenv(clumioApiBaseUrl)
	orgUnitId := "clumio_organizational_unit.test_ou1.id"
	if update{
		orgUnitId = "clumio_organizational_unit.test_ou2.id"
	}
	val := fmt.Sprintf(testAccResourceClumioUser, baseUrl, orgUnitId)
	return val
}

const testAccResourceClumioUser = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou1" {
  name = "test_ou1"
}

resource "clumio_organizational_unit" "test_ou2" {
  name = "test_ou2"
}

resource "clumio_user" "test_user" {
  full_name = "acceptance-test-user"
  email = "test@clumio.com"
  organizational_unit_ids = [%s]
}
`
