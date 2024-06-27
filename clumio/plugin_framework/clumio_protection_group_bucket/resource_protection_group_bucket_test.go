// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_protection_group_bucket Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_protection_group_bucket_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const BucketID = "BUCKET_ID"

// Basic test of the clumio_protection_group_bucket resource. It tests the following scenarios:
//   - Assigns a bucket to the protection group and verifies that the plan was applied properly.
func TestAccResourceClumioProtectionGroupBucket(t *testing.T) {

	bucketId := os.Getenv(BucketID)
	if bucketId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", BucketID))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroupBucket(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group_bucket.test_pg_bucket", "bucket_id",
						regexp.MustCompile(bucketId))),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_protection_group_bucket.test_pg_bucket", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// Tests that an external deletion of a clumio_protection_group_bucket resource leads to the
// resource needing to be re-created during the next plan. NOTE the Check function below as it is
// utilized to delete the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioProtectionGroupRecreate(t *testing.T) {

	bucketId := os.Getenv(BucketID)
	if bucketId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", BucketID))
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroupBucket(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group_bucket.test_pg_bucket", "bucket_id",
						regexp.MustCompile(bucketId)),
					DeleteProtectionGroupBucket("clumio_protection_group_bucket.test_pg_bucket"),
				),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_protection_group_bucket.test_pg_bucket",
							plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

// getTestAccResourceClumioProtectionGroupBucket returns the Terraform configuration for a basic
// clumio_protection_group_bucket resource.
func getTestAccResourceClumioProtectionGroupBucket() string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	bucketId := os.Getenv(BucketID)
	val := fmt.Sprintf(testAccResourceClumioProtectionGroupBucket, baseUrl, bucketId)
	return val
}

// DeleteProtectionGroupBucket deletes the bucket from the protection group using the Clumio API.
// It takes as argument, either the resource name or the actual id of the protection group bucket.
func DeleteProtectionGroupBucket(resourceName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.Attributes["bucket_id"] == "" {
			return fmt.Errorf("Bucket ID is not set")
		}
		bucketId := rs.Primary.Attributes["bucket_id"]
		pgId := rs.Primary.Attributes["protection_group_id"]

		clumioApiToken := os.Getenv(common.ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
		config := sdkconfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		pd := protectionGroups.NewProtectionGroupsV1(config)
		_, apiErr := pd.DeleteBucketProtectionGroup(pgId, bucketId)
		if apiErr != nil {
			return apiErr
		}
		time.Sleep(3 * time.Second)

		return nil
	}
}

// testAccResourceClumioProtectionGroupBucket is the Terraform configuration for a basic
// clumio_protection_group_bucket resource.
const testAccResourceClumioProtectionGroupBucket = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg_assignment"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_assignment"
  description = "test_pg_assignment"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}

resource "clumio_protection_group_bucket" "test_pg_bucket"{
  protection_group_id = clumio_protection_group.test_pg_assignment.id
  bucket_id = "%s"
}
`
