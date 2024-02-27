// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_role Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_role_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_role datasource. It tests that the role matching the name provided in
// the config is fetched and set in state.
func TestAccDataSourceClumioRoles(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioCallbackClumioRole(baseUrl, "Backup Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.clumio_role.test_role", "id", "30000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr(
						"data.clumio_role.test_role", "name", "Backup Admin"),
				),
			},
			{
				Config:      getTestDataSourceClumioCallbackClumioRole(baseUrl, "abcd"),
				ExpectError: regexp.MustCompile(".*Role not found.*"),
			},
		},
	})
}

// getTestDataSourceClumioCallbackClumioRole returns the Terraform configuration for a basic
// clumio_role datasource.
func getTestDataSourceClumioCallbackClumioRole(baseUrl string, roleName string) string {
	return fmt.Sprintf(testAccDataSourceClumioRoles, baseUrl, roleName)
}

// testAccDataSourceClumioRoles is the Terraform configuration for a basic clumio_role dataource.
const testAccDataSourceClumioRoles = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_role" "test_role" {
	name = "%s"
}
`
