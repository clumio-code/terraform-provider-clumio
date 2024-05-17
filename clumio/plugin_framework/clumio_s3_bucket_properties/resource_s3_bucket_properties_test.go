// Copyright 2024. Clumio, Inc.
//
// This files holds acceptance tests for the clumio_s3_bucket_properties Terraform resource. Please
// view the README.md file for more information on how to run these tests.

//go:build bucket

package clumio_s3_bucket_properties_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const BucketID = "BUCKET_ID"

// Basic test of the clumio_s3_bucket_properties resource. It tests the following scenarios:
//   - Creates S3 bucket properties and verifies that the plan was applied properly.
//   - Updates the S3 bucket properties and verifies that the resource will be updated.
func TestAccResourceClumioS3BucketProperties(t *testing.T) {
	// Retrieve the environment variables required for the test.
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	bucketId := os.Getenv(BucketID)
	if bucketId == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", BucketID))
		return
	}
	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioS3BucketProperties(baseUrl, bucketId, false),
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
							"clumio_s3_bucket_properties.test_s3_bucket_properties",
							plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_s3_bucket_properties.test_s3_bucket_properties",
						"event_bridge_enabled", regexp.MustCompile("true")),
					resource.TestMatchResourceAttr(
						"clumio_s3_bucket_properties.test_s3_bucket_properties",
						"event_bridge_notification_disabled", regexp.MustCompile("true")),
				),
			},
			{
				Config: getTestAccResourceClumioS3BucketProperties(baseUrl, bucketId, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_s3_bucket_properties.test_s3_bucket_properties",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_s3_bucket_properties.test_s3_bucket_properties",
						"event_bridge_enabled", regexp.MustCompile("false")),
					resource.TestMatchResourceAttr(
						"clumio_s3_bucket_properties.test_s3_bucket_properties",
						"event_bridge_notification_disabled", regexp.MustCompile("true")),
				),
			},
		},
	})
}

// getTestAccResourceClumioS3BucketProperties returns the Terraform configuration for a basic
// clumio_s3_bucket_properties resource.
func getTestAccResourceClumioS3BucketProperties(baseUrl, bucketId string, update bool) string {
	eventBridgeEnabled := "true"
	if update {
		eventBridgeEnabled = "false"
	}
	return fmt.Sprintf(
		testAccResourceClumioS3BucketProperties, baseUrl, bucketId, eventBridgeEnabled, "true")
}

// testAccResourceClumioS3BucketProperties is the Terraform configuration for a basic
// clumio_s3_bucket_properties resource.
const testAccResourceClumioS3BucketProperties = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_s3_bucket_properties" "test_s3_bucket_properties" {
  bucket_id = "%s"
  event_bridge_enabled = %s
  event_bridge_notification_disabled = %s
}
`
