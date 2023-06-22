// Copyright 2023. Clumio, Inc.
//
// Acceptance test for resource_organizational_unit.

package clumio_organizational_unit_test

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
	ouNameBefore   = "acceptance-test-ou"
	ouNameAfter    = "acceptance-test-ou-updated"
	descNameBefore = "test-ou-description"
	descNameAfter  = "test-ou-description-updated"
)

func TestAccResourceClumioOrganizationalUnit(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioOrganizationalUnit(baseUrl, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "name",
						regexp.MustCompile(ouNameBefore)),
				),
			},
			{
				Config: getTestAccResourceClumioOrganizationalUnit(baseUrl, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "name",
						regexp.MustCompile(ouNameAfter)),
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "description",
						regexp.MustCompile(descNameAfter)),
				),
			},
		},
	})
}

func getTestAccResourceClumioOrganizationalUnit(baseUrl string, update bool) string {
	content :=
		`name = "acceptance-test-ou"`
	if update {
		content =
			`name = "acceptance-test-ou-updated"
			 description = "test-ou-description-updated"
			`
	}
	return fmt.Sprintf(testAccResourceClumioOrganizationalUnit, baseUrl, content)
}

// description = "test-ou-description-updated"
const testAccResourceClumioOrganizationalUnit = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou" {
   %s
}
`
