// Copyright 2025. Clumio, Inc.

// This files holds acceptance tests for the clumio_report_configuration Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_report_configuration_test

import (
	"fmt"
	"os"
	"testing"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_report_configuration resource. It tests the following scenario:
//   - Creates a report configuration and verifies that the plan was applied properly.
//   - Updates the report configuration with a new description and verifies that the resource will
//     be updated.
func TestAccResourceClumioReportConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioReportConfiguration("test-1"),
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
				Config: getTestAccResourceClumioReportConfiguration("test-2"),
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

// Tests that an external deletion of a clumio_report_configuration resource leads to the resource
// needing to be re-created during the next plan. NOTE the Check function below as it is utilized to
// delete the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioReportConfigurationRecreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioReportConfiguration("test"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(resourceAddress,
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					// Delete the resource using the Clumio API after the plan is applied.
					deleteReportConfiguration(resourceAddress),
				),
				// This attribute is used to denote that the test expects that after the plan is
				// applied and a refresh is run, a non-empty plan is expected due to differences
				// from the state. Without this attribute set, the test would fail as it is unaware
				// that the resource was deleted externally.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// getTestAccResourceClumioReportConfiguration returns the Terraform configuration for a basic
// clumio_report_configuration resource.
func getTestAccResourceClumioReportConfiguration(suffix string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioReportConfiguration, baseUrl, suffix)
}

// deleteReportConfiguration returns a function that deletes a report configuration using the Clumio
// API with information from the Terraform state. It is used to intentionally cause a difference
// between the Terraform state and the actual state of the resource in the backend.
func deleteReportConfiguration(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		// Create a Clumio API client and delete the wallet.
		clumioApiToken := os.Getenv(common.ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
		config := clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		client := sdkclients.NewReportConfigurationClient(config)
		_, apiErr := client.DeleteComplianceReportConfiguration(rs.Primary.ID)
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

const resourceAddress = "clumio_report_configuration.test_report_configuration"

// testAccResourceClumioReportConfiguration is the Terraform configuration for a basic
// clumio_report_configuration resource.
const testAccResourceClumioReportConfiguration = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_report_configuration" "test_report_configuration" {
  name = "test-report-configuration"
  description = "%s"
  notification {
	email_list = ["email1@clumio.com", "email2@clumio.com"]
  }
  parameter {
    controls {
	  asset_backup {
	    look_back_period {
		  unit = "days"
		  value = 7
		}
		minimum_retention_duration {
		  unit = "days"
		  value = 30
		}
		window_size {
		  unit = "days"
		  value = 7
		}
	  }
	  asset_protection {
	    should_ignore_deactivated_policy = true
	  }
	  policy {
	    minimum_retention_duration {
		  unit = "days"
		  value = 1
		}
		minimum_rpo_frequency {
		  unit = "days"
		  value = 1
		}
	  }
	}
	filters {
	  asset {
	    tag_op_mode = "equal"
		tags {
		  key = "environment"
		  value = "production"
		}
	  }
	}
  }
  schedule {
    day_of_week = "sunday"
	frequency = "weekly"
	start_time = "15:00"
	timezone = "America/New_York"
  }
}
`
