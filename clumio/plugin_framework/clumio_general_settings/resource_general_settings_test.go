// Copyright 2025. Clumio, Inc.

// This files holds acceptance tests for the clumio_general_settings Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build general_settings

package clumio_general_settings_test

import (
	"fmt"
	"os"
	"testing"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Basic test of the clumio_general_settings resource. It tests the following scenario:
//   - Creates a general settings and verifies that the plan was applied properly.
//   - Updates the general settings with new values and verifies that the resource will
//     be updated.
func TestAccResourceClumioGeneralSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioGeneralSettings(1200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(resourceAddress,
							plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							resourceAddress,
							plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: getTestAccResourceClumioGeneralSettings(2400),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(resourceAddress,
							plancheck.ResourceActionUpdate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(resourceAddress,
							plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// getTestAccResourceClumioGeneralSettings returns the Terraform configuration for a basic
// clumio_general_settings resource.
func getTestAccResourceClumioGeneralSettings(autoLogoutDuration int64) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioGeneralSettings, baseUrl, autoLogoutDuration)
}

const resourceAddress = "clumio_general_settings.test_general_settings"

// testAccResourceClumioGeneralSettings is the Terraform configuration for a basic
// clumio_general_settings resource.
const testAccResourceClumioGeneralSettings = `
provider clumio{
  clumio_api_base_url = "%s"
}

data "http" "myip" {
  url = "https://ipv4.icanhazip.com"
}

resource "clumio_general_settings" "test_general_settings" {
  auto_logout_duration         = %d
  password_expiration_duration = 7776000
  ip_allowlist                 = ["${chomp(data.http.myip.response_body)}/32", "192.168.1.2"]
}
`
