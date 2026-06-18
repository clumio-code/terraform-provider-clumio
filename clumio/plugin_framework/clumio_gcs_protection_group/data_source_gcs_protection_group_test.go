// Copyright 2026. Clumio, Inc.

// This file holds acceptance tests for the clumio_gcs_protection_group Terraform data source.
// Please view the README.md file for more information on how to run these tests.

//go:build gcp

package clumio_gcs_protection_group_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dataSourceName = "data.clumio_gcs_protection_group.ds_gcp_pg"

// TestAccDataSourceClumioGCSProtectionGroup creates a protection group and then looks it up by name
// through the data source, verifying the data source resolves the same id.
func TestAccDataSourceClumioGCSProtectionGroup(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf_acc_gcp_pg_ds"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { gcpAccPreCheck(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccDataSourceGCSProtectionGroup(baseUrl, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

// TestAccDataSourceClumioGCSProtectionGroupEmptyName verifies that an empty name is rejected by
// schema validation.
func TestAccDataSourceClumioGCSProtectionGroupEmptyName(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { gcpAccPreCheck(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceEmptyName, baseUrl),
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value.*"),
			},
		},
	})
}

// getTestAccDataSourceGCSProtectionGroup returns a configuration that creates a protection group and
// reads it back through the data source by name.
func getTestAccDataSourceGCSProtectionGroup(baseUrl, name string) string {
	return fmt.Sprintf(`
provider clumio {
  clumio_api_base_url = "%s"
}

resource "clumio_gcs_protection_group" "test_gcp_pg" {
  name = "%s"
}

data "clumio_gcs_protection_group" "ds_gcp_pg" {
  name = clumio_gcs_protection_group.test_gcp_pg.name
}
`, baseUrl, name)
}

// testAccDataSourceEmptyName is the Terraform configuration for a clumio_gcs_protection_group data
// source with name set to an empty string.
const testAccDataSourceEmptyName = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_gcs_protection_group" "ds_gcp_pg" {
  name = ""
}
`
