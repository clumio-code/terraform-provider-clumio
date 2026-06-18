// Copyright 2026. Clumio, Inc.

// This file holds acceptance tests for the clumio_gcs_protection_group Terraform resource. Please
// view the README.md file for more information on how to run these tests.

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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const resourceName = "clumio_gcs_protection_group.test_gcp_pg"

// gcpAccPreCheck validates the Clumio credentials required by the GCP acceptance lane.
func gcpAccPreCheck(t *testing.T) {
	clumioPf.UtilTestFailIfEmpty(t, common.ClumioApiToken, common.ClumioApiToken+" cannot be empty.")
	clumioPf.UtilTestFailIfEmpty(
		t, common.ClumioApiBaseUrl, common.ClumioApiBaseUrl+" cannot be empty.")
}

// TestAccResourceClumioGCSProtectionGroup exercises the full lifecycle of the resource:
//   - Create with a bucket_rule and include_prefixes, then verify a refresh produces an empty plan.
//   - Update the name and verify the resource is updated in place.
//   - Remove the bucket_rule and include_prefixes and verify they are cleared.
func TestAccResourceClumioGCSProtectionGroup(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf_acc_gcp_pg_1"
	updatedName := "tf_acc_gcp_pg_2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			gcpAccPreCheck(t)
			clumioPf.UtilTestFailIfEmpty(t, common.ClumioTestGcpProjectId,
				common.ClumioTestGcpProjectId+" cannot be empty")
		},
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with optional attributes set.
			{
				Config: getTestAccResourceGCSProtectionGroup(baseUrl, name, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName, "bucket_rule.gcp_project_id.eq",
						os.Getenv(common.ClumioTestGcpProjectId)),
					resource.TestCheckResourceAttr(resourceName, "include_prefixes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
			// Update the name in place.
			{
				Config: getTestAccResourceGCSProtectionGroup(baseUrl, updatedName, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			// Remove the optional attributes and verify they are cleared.
			{
				Config: getTestAccResourceGCSProtectionGroup(baseUrl, updatedName, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "include_prefixes.#", "0"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// TestAccResourceClumioGCSProtectionGroupNoOptional tests creation without any optional attributes
// and verifies that a subsequent refresh produces an empty plan.
func TestAccResourceClumioGCSProtectionGroupNoOptional(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf_acc_gcp_pg_no_optional"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { gcpAccPreCheck(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceGCSProtectionGroup(baseUrl, name, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "latest_version_only", "true"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAccResourceClumioGCSProtectionGroupEmptyName verifies that an empty name is rejected by schema
// validation.
func TestAccResourceClumioGCSProtectionGroupEmptyName(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { gcpAccPreCheck(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getTestAccResourceGCSProtectionGroup(baseUrl, "", false, false),
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value.*"),
			},
		},
	})
}

// TestAccResourceClumioGCSProtectionGroupLabelRule exercises the gcp_label bucket_rule operators in
// their distinct forms and verifies each round-trips with an empty refresh plan:
//   - eq: a single-label map operator.
//   - all: a multi-label map operator.
//   - in: repeated {key, value} blocks carrying the same key with multiple values.
func TestAccResourceClumioGCSProtectionGroupLabelRule(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf_acc_gcp_pg_label"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { gcpAccPreCheck(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// eq: single-label map operator.
			{
				Config: getTestAccResourceGCSProtectionGroupLabelRule(baseUrl, name,
					`      eq = { dongjoo = "test" }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "bucket_rule.gcp_label.eq.dongjoo", "test"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// all: multi-label map operator.
			{
				Config: getTestAccResourceGCSProtectionGroupLabelRule(baseUrl, name,
					`      all = { dongjoo = "test", team = "data" }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "bucket_rule.gcp_label.all.dongjoo", "test"),
					resource.TestCheckResourceAttr(
						resourceName, "bucket_rule.gcp_label.all.team", "data"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(
						resourceName, plancheck.ResourceActionUpdate)},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// in: repeated blocks, same key with multiple values (a map could not express this).
			{
				Config: getTestAccResourceGCSProtectionGroupLabelRule(baseUrl, name,
					`      in {
        key   = "env"
        value = "prod"
      }
      in {
        key   = "env"
        value = "staging"
      }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "bucket_rule.gcp_label.in.#", "2"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// getTestAccResourceGCSProtectionGroupLabelRule returns a clumio_gcs_protection_group configuration
// whose bucket_rule contains a gcp_label block with the supplied body (the operator lines).
func getTestAccResourceGCSProtectionGroupLabelRule(baseUrl, name, labelBody string) string {
	return fmt.Sprintf(`
provider clumio {
  clumio_api_base_url = "%s"
}

resource "clumio_gcs_protection_group" "test_gcp_pg" {
  name = "%s"
  bucket_rule {
    gcp_label {
%s
    }
  }
}
`, baseUrl, name, labelBody)
}

// getTestAccResourceGCSProtectionGroup returns the Terraform configuration for a
// clumio_gcs_protection_group resource. The bucketRule and prefixes flags toggle the optional
// attributes.
func getTestAccResourceGCSProtectionGroup(
	baseUrl, name string, bucketRule bool, prefixes bool) string {

	bucketRuleStr := ""
	if bucketRule {
		projectId := os.Getenv(common.ClumioTestGcpProjectId)
		bucketRuleStr = fmt.Sprintf(`
  bucket_rule {
    gcp_project_id {
      eq = "%s"
    }
  }`, projectId)
	}
	prefixesStr := ""
	if prefixes {
		prefixesStr = `include_prefixes = ["data/"]`
	}
	return fmt.Sprintf(`
provider clumio {
  clumio_api_base_url = "%s"
}

resource "clumio_gcs_protection_group" "test_gcp_pg" {
  name = "%s"
  %s
  %s
}
`, baseUrl, name, bucketRuleStr, prefixesStr)
}
