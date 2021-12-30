// Copyright 2021. Clumio, Inc.

// Acceptance test for clumio_organizational_unit resource.
package clumio

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

)

func TestAccResourceClumioOrganizationalUnit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckClumio(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioOrganizationalUnit(false),
			},
			{
				Config: getTestAccResourceClumioOrganizationalUnit(true),
			},
		},
	})
}

func getTestAccResourceClumioOrganizationalUnit(update bool) string {
	baseUrl := os.Getenv(clumioApiBaseUrl)
	name := "acceptance-test-ou"
	user := "clumio_user.test_user1.id"
	if update{
		name = "acceptance-test-ou-updated"
		user = "clumio_user.test_user2.id"
	}
	return fmt.Sprintf(testAccResourceClumioOrganizationalUnit, baseUrl, name, user)
}

const testAccResourceClumioOrganizationalUnit = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_user" "test_user1" {
  full_name = "acceptance-test-user"
  email = "test1@clumio.com"
}

resource "clumio_user" "test_user2" {
  full_name = "acceptance-test-user"
  email = "test2@clumio.com"
}

resource "clumio_organizational_unit" "test_ou" {
  name = "%s"
  users = [%s]
}
`
