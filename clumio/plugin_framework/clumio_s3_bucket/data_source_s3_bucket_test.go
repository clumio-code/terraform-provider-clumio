// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_s3_bucket Terraform datasource. Please
// view the README.md file for more information on how to run these tests.

//go:build basic

package clumio_s3_bucket_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test of the clumio_s3_bucket datasource with name specified as empty string.
func TestAccDataSourceClumioEmptyOUName(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceEmptyName, baseUrl),
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value.*"),
			},
		},
	})
}

// testAccDataSourceEmptyName is the Terraform configuration for a
// clumio_s3_bucket datasource with name set to empty string.
const testAccDataSourceEmptyName = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_s3_bucket" "ds_s3_bucket" {
  bucket_names=[]
}
`
