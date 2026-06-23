// Copyright (c) 2026 Clumio, a Commvault Company All Rights Reserved

// This file holds acceptance tests for the clumio_gcp_connection Terraform resource.
//
// These tests cover the beta GCP connection resource and are gated behind the dedicated "gcp" build
// tag (run via `make testacc_gcp_connection`) so they stay out of the shared `basic`/`post_process` CI lanes
// until the GCP APIs are part of a published clumio-go-sdk release. They require a beta-enabled
// Clumio backend and the CLUMIO_TEST_GCP_PROJECT_ID/CLUMIO_TEST_GCP_PROJECT_ID2 environment
// variables in addition to the Clumio API credentials.

//go:build gcp_connection

package clumio_gcp_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_gcp_connection resource. It tests the following scenarios:
//   - Creates a connection and verifies that the plan was applied properly.
//   - Updates the description and regions and verifies that the resource will be updated.
//   - Ensures that updates to the project ID requires that the resource is re-created as opposed to
//     just updated.
func TestAccResourceClumioGcpConnection(t *testing.T) {
	// Retrieve the environment variables required for the test.
	projectId := os.Getenv(common.ClumioTestGcpProjectId)
	projectId2 := os.Getenv(common.ClumioTestGcpProjectId2)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestGcpConnectionPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioGcpConnection(
					baseUrl, projectId, "test_description", `["us-west1"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "project_id", projectId),
				),
			},
			{
				Config: getTestAccResourceClumioGcpConnection(
					baseUrl, projectId, "test_description_updated", `["us-west1", "us-central1"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "project_id", projectId),
					resource.TestMatchResourceAttr(
						"clumio_gcp_connection.test_conn", "description",
						regexp.MustCompile("test_description_updated")),
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "regions.#", "2"),
				),
			},
			{
				Config: getTestAccResourceClumioGcpConnection(
					baseUrl, projectId2, "test_description_updated", `["us-west1", "us-central1"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "project_id", projectId2),
				),
			},
		},
	})
}

// Tests that an external deletion of a clumio_gcp_connection resource leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioGcpConnectionRecreate(t *testing.T) {
	// Retrieve the environment variables required for the test.
	projectId := os.Getenv(common.ClumioTestGcpProjectId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestGcpConnectionPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioGcpConnection(
					baseUrl, projectId, "test_description", `["us-west1"]`),
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
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "project_id", projectId),
					// Delete the resource using the Clumio API after the plan is applied.
					deleteGcpConnection("clumio_gcp_connection.test_conn"),
				),
				// This attribute is used to denote that the test expects that after the plan is
				// applied and a refresh is run, a non-empty plan is expected due to differences from
				// the state. Without this attribute set, the test would fail as it is unaware that
				// the resource was deleted externally.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Tests that changing the deployment_type attribute forces the resource to be re-created, as the
// attribute is configured with a RequiresReplace plan modifier.
func TestAccResourceClumioGcpConnectionDeploymentTypeReplace(t *testing.T) {
	// Retrieve the environment variables required for the test.
	projectId := os.Getenv(common.ClumioTestGcpProjectId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestGcpConnectionPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioGcpConnectionDeploymentType(
					baseUrl, projectId, "direct_terraform"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "deployment_type", "direct_terraform"),
				),
			},
			{
				Config: getTestAccResourceClumioGcpConnectionDeploymentType(
					baseUrl, projectId, "infrastructure_manager"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_gcp_connection.test_conn", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_gcp_connection.test_conn", "deployment_type",
						"infrastructure_manager"),
				),
			},
		},
	})
}

// deleteGcpConnection returns a function that deletes a GCP connection using the Clumio API with
// information from the Terraform state. It is used to intentionally cause a difference between the
// Terraform state and the actual state of the resource in the backend. NOTE the GCP connection is
// keyed by project_id (not the resource id, which holds the connection token), so the delete is
// performed using the project_id attribute from state.
func deleteGcpConnection(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Retrieve the resource by name from state.
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		projectId := rs.Primary.Attributes["project_id"]
		if projectId == "" {
			return fmt.Errorf("project_id is not set")
		}

		// Create a Clumio API client and delete the GCP connection.
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
		gcpConnection := sdkclients.NewGcpConnectionClient(config)
		_, apiErr := gcpConnection.DeleteGcpConnection(projectId)
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

// getTestAccResourceClumioGcpConnection returns the Terraform configuration for a basic
// clumio_gcp_connection resource.
func getTestAccResourceClumioGcpConnection(
	baseUrl string, projectId string, description string, regions string) string {
	return fmt.Sprintf(testAccResourceClumioGcpConnection, baseUrl, projectId, description, regions)
}

// getTestAccResourceClumioGcpConnectionDeploymentType returns the Terraform configuration for a
// clumio_gcp_connection resource with an explicit deployment_type.
func getTestAccResourceClumioGcpConnectionDeploymentType(
	baseUrl string, projectId string, deploymentType string) string {
	return fmt.Sprintf(
		testAccResourceClumioGcpConnectionDeploymentType, baseUrl, projectId, deploymentType)
}

// testAccResourceClumioGcpConnection is the Terraform configuration for a basic
// clumio_gcp_connection resource.
const testAccResourceClumioGcpConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_gcp_connection" "test_conn" {
  project_id  = "%s"
  description = "%s"
  regions     = %s
}
`

// testAccResourceClumioGcpConnectionDeploymentType is the Terraform configuration for a
// clumio_gcp_connection resource with an explicit deployment_type attribute.
const testAccResourceClumioGcpConnectionDeploymentType = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_gcp_connection" "test_conn" {
  project_id      = "%s"
  deployment_type = "%s"
  regions         = ["us-west1"]
}
`
