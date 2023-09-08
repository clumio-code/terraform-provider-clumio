// Copyright 2023. Clumio, Inc.
//
// Acceptance test for resource_user.

package clumio_user_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	assignedRoleBefore = "30000000-0000-0000-0000-000000000000"
	assignedRoleAfter  = "20000000-0000-0000-0000-000000000000"
)

func TestAccResourceClumioUserV1(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioUserV1(baseUrl, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_user.test_user", "assigned_role",
						regexp.MustCompile(assignedRoleBefore)),
				),
			},
			{
				Config: getTestAccResourceClumioUserV1(baseUrl, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_user.test_user", "assigned_role",
						regexp.MustCompile(assignedRoleAfter)),
				),
			},
		},
	})
}

func getTestAccResourceClumioUserV1(baseUrl string, update bool) string {
	orgUnitId := "clumio_organizational_unit.test_ou1.id"
	assignedRole := assignedRoleBefore
	if update {
		orgUnitId = "clumio_organizational_unit.test_ou2.id"
		assignedRole = assignedRoleAfter
	}
	return fmt.Sprintf(testAccResourceClumioUserV1, baseUrl, assignedRole, orgUnitId)
}

const testAccResourceClumioUserV1 = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou1" {
  name = "test_ou1"
  description = "test-ou-1"
}

resource "clumio_organizational_unit" "test_ou2" {
  name = "test_ou2"
  description = "test-ou-2"
}

resource "clumio_user" "test_user" {
  full_name = "acceptance-test-user"
  email = "test@clumio.com"
  assigned_role = "%s"
  organizational_unit_ids = [%s]
}
`

func TestAccResourceClumioUser(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioUser(baseUrl, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_user.test_user", "assigned_role",
						regexp.MustCompile(assignedRoleBefore)),
				),
			},
			{
				Config: getTestAccResourceClumioUser(baseUrl, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_user.test_user", "assigned_role",
						regexp.MustCompile(assignedRoleAfter)),
				),
			},
		},
	})
}

func getTestAccResourceClumioUser(baseUrl string, update bool) string {
	orgUnitId := "clumio_organizational_unit.test_ou1.id"
	assignedRole := assignedRoleBefore
	if update {
		orgUnitId = "clumio_organizational_unit.test_ou2.id"
		assignedRole = assignedRoleAfter
	}
	return fmt.Sprintf(testAccResourceClumioUser, baseUrl, assignedRole, orgUnitId)
}

const testAccResourceClumioUser = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou1" {
  name = "test_ou1"
  description = "test-ou-1"
}

resource "clumio_organizational_unit" "test_ou2" {
  name = "test_ou2"
  description = "test-ou-2"
}

resource "clumio_user" "test_user" {
  full_name = "acceptance-test-user"
  email = "test@clumio.com"
  access_control_configuration = [
	{
		role_id = "%s"
		organizational_unit_ids = [%s]
	},
  ]
}
`
