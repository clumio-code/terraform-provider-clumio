// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_protection_group Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_protection_group_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_protection_group datasource. It tests that the protection_group matching
// the name provided in the config is fetched and its Id is set in the state.
func TestAccDataSourceClumioProtectionGroup(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioProtectionGroup(baseUrl, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.clumio_protection_group.ds_pg",
						"id"),
				),
			},
			{
				Config:      getTestDataSourceClumioProtectionGroup(baseUrl, true, false),
				ExpectError: regexp.MustCompile(".*Protection group not found.*"),
			},
		},
	})
}

// Test of the clumio_protection_group datasource with name specified as empty string.
func TestAccDataSourceClumioEmptyProtectionGroup(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceEmptyClumioProtectionGroup, baseUrl),
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
		},
	})
}

// getTestDataSourceClumioProtectionGroup returns the Terraform configuration for a basic
// clumio_protection_group datasource.
func getTestDataSourceClumioProtectionGroup(
	baseUrl string, invalidName bool, emptyName bool) string {

	name := "ds_test_pg"

	if invalidName {
		name = "some-name"
	} else if emptyName {
		name = ""
	}

	return fmt.Sprintf(testAccDataSourceClumioProtectionGroup, baseUrl, name)
}

// testAccDataSourceClumioProtectionGroup is the Terraform configuration for a basic
// clumio_protection_group datasource.
const testAccDataSourceClumioProtectionGroup = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "ds_test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "ds_test_pg"
  description = "Acceptance test protection group for protection group data source."
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}

data "clumio_protection_group" "ds_pg" {
	depends_on = [ clumio_protection_group.ds_test_pg]
	name = "%s"
}
`

// testAccDataSourceEmptyClumioProtectionGroup is the Terraform configuration for a
// clumio_protection_group datasource with name set to empty string.
const testAccDataSourceEmptyClumioProtectionGroup = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_protection_group" "ds_pg" {
	depends_on = [ clumio_protection_group.ds_test_pg]
	name = ""
}
`
