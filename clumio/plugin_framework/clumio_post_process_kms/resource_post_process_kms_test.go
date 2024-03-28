// Copyright 2023. Clumio, Inc.
//
// This files holds acceptance tests for the clumio_post_process_kms Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build post_process

package clumio_post_process_kms_test

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

// Basic test of the clumio_post_process_kms resource. It tests the following scenarios:
//   - Create scenario for post process kms and verifies that the plan was applied properly.
//   - Updates the config for post process kms and verifies that the resource will be updated.
func TestAccResourcePostProcessKMS(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourcePostProcessKMS(
					baseUrl, accountNativeId, testAwsRegion, false),
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
							"clumio_post_process_kms.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_kms.test", "template_version",
						regexp.MustCompile("1")),
				),
				SkipFunc: shouldSkip,
			},
			{
				Config: getTestAccResourcePostProcessKMS(
					baseUrl, accountNativeId, testAwsRegion, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_kms.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_kms.test", "template_version",
						regexp.MustCompile("2")),
				),
				SkipFunc: shouldSkip,
			},
		},
	})
}

// Function to determine if the acceptance test step should be skipped.
func shouldSkip() (bool, error) {
	// Due to a limitation in the Wallets API, the cleanup of the wallet after the
	// test is complete fails if the wallet is configured with a key. So we skip
	// this acceptance test till it is fixed.
	return true, nil
}

// getTestAccResourcePostProcessKMS returns the Terraform configuration for the
// clumio_post_process_kms resource.
func getTestAccResourcePostProcessKMS(
	baseUrl string, accountId string, awsRegion string, update bool) string {

	templateVersion := "1"
	if update {
		templateVersion = "2"
	}
	return fmt.Sprintf(
		testAccResourcePostProcessKms, baseUrl, accountId, awsRegion, templateVersion)
}

// testAccResourcePostProcessKms is the Terraform configuration for the clumio_post_process_kms
// resource.
const testAccResourcePostProcessKms = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_wallet" "test_wallet" {
  account_native_id = "%s"
}

resource "clumio_post_process_kms" "test" {
  token = clumio_wallet.test_wallet.token
  account_id = clumio_wallet.test_wallet.account_native_id
  region = "%s"
  multi_region_cmk_key_id = "test_multi_region_cmk_key_id"
  role_external_id = "test_role_external_id"
  role_arn = "arn:aws:iam::${clumio_wallet.test_wallet.account_native_id}:role/testArn"
  role_id = "test_role_id"
  template_version = "%s"
}
`
