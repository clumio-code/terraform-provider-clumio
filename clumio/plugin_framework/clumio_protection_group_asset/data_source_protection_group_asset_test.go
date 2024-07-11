// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_protection_group_asset Terraform datasource. Please
// view the README.md file for more information on how to run these tests.

//go:build basic

package clumio_protection_group_asset_test

import (
	"fmt"
	"os"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	BucketID          = "BUCKET_ID"
	ProtectionGroupID = "PROTECTION_GROUP_ID"
)

// Test of the clumio_protection_group_asset datasource with name specified as empty string.
func TestAccDataSourceClumioEmptyOUName(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	bucketId := os.Getenv(BucketID)
	if bucketId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", BucketID))
		return
	}
	pgId := os.Getenv(ProtectionGroupID)
	if pgId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", ProtectionGroupID))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceProtectionGroupAsset, baseUrl, pgId, bucketId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.clumio_protection_group_asset.pg_asset", "id"),
				),
			},
		},
	})
}

// testAccDataSourceProtectionGroupAsset is the Terraform configuration for a
// clumio_protection_group_asset datasource.
const testAccDataSourceProtectionGroupAsset = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_protection_group_asset" "pg_asset" {
  protection_group_id = "%s"
  bucket_id = "%s"
}
`
