// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_organizational_unit Terraform datasource. Please
// view the README.md file for more information on how to run these tests.

//go:build basic

package clumio_organizational_unit_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_organizational_unit datasource. It tests that the organizational units
// matching the name provided in the config is fetched and set in the state.
func TestAccDataSourceClumioOU(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioOU(baseUrl, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_organizational_unit.ds_organizational_unit",
						"organizational_units.#", "1"),
				),
			},
		},
	})
}

// Test of the clumio_organizational_unit datasource with name specified as empty string.
func TestAccDataSourceClumioEmptyOUName(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceEmptyName, baseUrl),
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
		},
	})
}

// getTestDataSourceClumioOU returns the Terraform configuration for a basic clumio_organizational_unit
// datasource.
func getTestDataSourceClumioOU(baseUrl string, invalidName bool) string {

	name := "ds_test_organizational_unit"

	dsName := name
	if invalidName {
		dsName = "some-name"
	}

	return fmt.Sprintf(testAccDataSourceClumioOU, baseUrl, name, dsName)
}

// testAccDataSourceClumioOU is the Terraform configuration for a basic
// clumio_organizational_unit datasource.
const testAccDataSourceClumioOU = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "ds_test_organizational_unit"{
  name = "%s"
  description = "ds_test_organizational_unit"
}

data "clumio_organizational_unit" "ds_organizational_unit" {
  depends_on = [ clumio_organizational_unit.ds_test_organizational_unit]
  name = "%s"
}
`

// testAccDataSourceEmptyName is the Terraform configuration for a
// clumio_organizational_unit datasource with name set to empty string.
const testAccDataSourceEmptyName = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_organizational_unit" "ds_organizational_unit" {
  name=""
}
`
