// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_auto_user_provisioning_setting_test Terraform
// resource. Please view the README.md file for more information on how to run these tests.

//go:build sso

package clumio_auto_user_provisioning_setting_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	sdkAUPSettings "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_settings"
	"github.com/clumio-code/clumio-go-sdk/models"
	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_auto_user_provisioning_setting resource. It tests the following
// scenarios:
//   - Sets auto user provisioning setting to disabled and verifies that the plan was applied
//     properly.
//   - Updates the auto user provisioning setting to enabled and verifies that the resource will be
//     updated.
func TestAccClumioAutoUserProvisioningSetting(t *testing.T) {
	enabled := "true"
	disabled := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningSetting(disabled),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
							plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
							plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
						"is_enabled",
						regexp.MustCompile(disabled)),
				),
				SkipFunc: clumioPf.SkipIfSSONotConfigured,
			},
			{
				Config: getTestAccResourceClumioAutoUserProvisioningSetting(enabled),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
						"is_enabled",
						regexp.MustCompile(enabled)),
				),
				SkipFunc: clumioPf.SkipIfSSONotConfigured,
			},
		},
	})
}

// Tests that an external disabling of auto user provisioning setting leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to disable
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioAutoUserProvisioningSettingResync(t *testing.T) {

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningSetting("true"),
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
							"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_setting.test_auto_user_provisioning_setting",
						"is_enabled",
						regexp.MustCompile("true")),
					// Delete the resource using the Clumio API after the plan is applied.
					disableClumioAutoUserProvisioningSetting(),
				),
				// This attribute is used to denote that the test expects that after the plan is
				// applied and a refresh is run, a non-empty plan is expected due to differences
				// from the state. Without this attribute set, the test would fail as it is unaware
				// that the resource was deleted externally.
				ExpectNonEmptyPlan: true,
				SkipFunc:           clumioPf.SkipIfSSONotConfigured,
			},
		},
	})
}

// disableClumioAutoUserProvisioningSetting returns a function that disables the auto user
// provisioning setting using the Clumio API. It is used to intentionally cause a difference between
// the Terraform state and the actual state of the resource in the backend.
func disableClumioAutoUserProvisioningSetting() resource.TestCheckFunc {

	return func(s *terraform.State) error {

		// Return if SSO is not configured.
		if strings.ToLower(os.Getenv(common.ClumioTestIsSSOConfigured)) != "true" {
			return nil
		}

		// Create a Clumio API client and disable the auto user provisioning setting.
		clumioApiToken := os.Getenv(common.ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
		clumioConfig := clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		aups := sdkAUPSettings.NewAutoUserProvisioningSettingsV1(clumioConfig)
		isEnabled := false
		_, apiErr := aups.UpdateAutoUserProvisioningSetting(
			&models.UpdateAutoUserProvisioningSettingV1Request{
				IsEnabled: &isEnabled,
			})
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

// getTestAccResourceClumioAutoUserProvisioningSetting returns the Terraform configuration for a
// clumio_auto_user_provisioning_setting resource with description attribute not set.
func getTestAccResourceClumioAutoUserProvisioningSetting(isEnabled string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioAutoUserProvisioningSetting, baseUrl, isEnabled)
}

// testAccResourceClumioAutoUserProvisioningSetting is the Terraform configuration for a
// clumio_auto_user_provisioning_setting resource with description attribute not set.
const testAccResourceClumioAutoUserProvisioningSetting = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_auto_user_provisioning_setting" "test_auto_user_provisioning_setting" {
  is_enabled = %s
}

`
