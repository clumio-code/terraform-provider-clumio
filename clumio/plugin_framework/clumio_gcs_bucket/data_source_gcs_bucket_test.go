// Copyright 2026. Clumio, Inc.

// This file holds acceptance tests for the clumio_gcs_bucket Terraform data source. Please
// view the README.md file for more information on how to run these tests. These tests run in the
// gcp lane (make testacc_gcp), which assumes a GCP connection already exists in
// the target Clumio environment.

//go:build gcp

package clumio_gcs_bucket_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// gcpAccPreCheck validates the Clumio credentials required by the GCP acceptance lane.
func gcpAccPreCheck(t *testing.T) {
	clumioPf.UtilTestFailIfEmpty(t, common.ClumioApiToken, common.ClumioApiToken+" cannot be empty.")
	clumioPf.UtilTestFailIfEmpty(
		t, common.ClumioApiBaseUrl, common.ClumioApiBaseUrl+" cannot be empty.")
}

// TestAccDataSourceClumioGCSBucketEmptyName tests the clumio_gcs_bucket data source with an
// empty bucket_names set. It should return a validation error.
func TestAccDataSourceClumioGCSBucketEmptyName(t *testing.T) {
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

// TestAccDataSourceClumioGCSBucketNotFound tests that querying for a bucket name that matches no
// buckets returns a "not found" error.
func TestAccDataSourceClumioGCSBucketNotFound(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { gcpAccPreCheck(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccDataSourceNotFound, baseUrl),
				ExpectError: regexp.MustCompile(".*GCS bucket not found.*"),
			},
		},
	})
}

// testAccDataSourceEmptyName is the Terraform configuration for a clumio_gcs_bucket data source
// with bucket_names set to an empty list.
const testAccDataSourceEmptyName = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_gcs_bucket" "ds_gcs_bucket" {
  bucket_names=[]
}
`

// testAccDataSourceNotFound is the Terraform configuration for a clumio_gcs_bucket data source
// querying bucket names that are not expected to match any bucket.
const testAccDataSourceNotFound = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_gcs_bucket" "ds_gcs_bucket" {
  bucket_names = ["tf-acc-nonexistent-bucket-name-zzz"]
}
`
