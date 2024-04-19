// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_user Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_user_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_user datasource. It tests that the user matching
// the name provided in the config is fetched and its Id is set in the state.
func TestAccDataSourceClumioUser(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioUser(baseUrl, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_user.ds_user",
						"users.#", "1"),
				),
			},
		},
	})
}

// Test of the clumio_user datasource with name specified as empty string.
func TestAccDataSourceClumioEmptyUserName(t *testing.T) {
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

// Test of the clumio_user datasource with name and role_id not set in the config.
func TestAccDataSourceEmptyUser(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceEmptyClumioUser, baseUrl),
				ExpectError: regexp.MustCompile(".*Missing Attribute Configuration.*"),
			},
		},
	})
}

// getTestDataSourceClumioUser returns the Terraform configuration for a basic clumio_user
// datasource.
func getTestDataSourceClumioUser(baseUrl string, invalidName bool) string {

	name := "ds_test_user"

	dsName := name
	if invalidName {
		dsName = "some-name"
	}

	return fmt.Sprintf(testAccDataSourceClumioUser, baseUrl, name, dsName)
}

// testAccDataSourceClumioUser is the Terraform configuration for a basic
// clumio_user datasource.
const testAccDataSourceClumioUser = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_user" "ds_test_user"{
  full_name = "%s"
  email = "test@clumio.com"
  access_control_configuration = [
	{
		role_id = "00000000-0000-0000-0000-000000000000"
		organizational_unit_ids = ["00000000-0000-0000-0000-000000000000"]
	},
  ]
}

data "clumio_user" "ds_user" {
  depends_on = [ clumio_user.ds_test_user]
  name = "%s"
  role_id ="00000000-0000-0000-0000-000000000000" 
}
`

// testAccDataSourceEmptyName is the Terraform configuration for a
// clumio_user datasource with name set to empty string.
const testAccDataSourceEmptyName = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_user" "ds_user" {
  name=""
}
`

// testAccDataSourceEmptyClumioUser is the Terraform configuration for a clumio_user datasource with
// both name and role_id not set.
const testAccDataSourceEmptyClumioUser = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_user" "ds_user" {
}
`
