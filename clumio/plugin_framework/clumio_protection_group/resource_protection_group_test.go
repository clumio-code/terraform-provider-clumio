// Copyright 2023. Clumio, Inc.

// Acceptance test for clumio_s3_protection_group resource.
package clumio_protection_group_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioProtectionGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroup(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1"))),
			},
			{
				Config: getTestAccResourceClumioProtectionGroup(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1_updated"))),
			},
		},
	})
}

func getTestAccResourceClumioProtectionGroup(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	description := "test_pg_1"
	if update {
		description = "test_pg_1_updated"
	}
	val := fmt.Sprintf(testAccResourceClumioProtectionGroup, baseUrl, description)
	return val
}

const testAccResourceClumioProtectionGroup = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "%s"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}
`

func TestAccResourceClumioProtectionGroupWithOU(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroupWithOU(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1"))),
			},
			{
				Config: getTestAccResourceClumioProtectionGroupWithOU(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1_updated"))),
			},
		},
	})
}

func getTestAccResourceClumioProtectionGroupWithOU(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	description := "test_pg_1"
	if update {
		description = "test_pg_1_updated"
	}
	val := fmt.Sprintf(testAccResourceClumioProtectionGroupWithOU, baseUrl, description)
	return val
}

const testAccResourceClumioProtectionGroupWithOU = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou2" {
  name = "test_ou2"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "%s"
  organizational_unit_id = clumio_organizational_unit.test_ou2.id
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}
`
