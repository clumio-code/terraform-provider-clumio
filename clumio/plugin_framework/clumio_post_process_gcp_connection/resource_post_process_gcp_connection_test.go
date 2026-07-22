// Copyright (c) 2026 Clumio, a Commvault Company All Rights Reserved

// This file holds acceptance tests for the clumio_post_process_gcp_connection Terraform resource.
//
// These tests cover the beta GCP post-process resource and are gated behind the dedicated "gcp"
// build tag (run via `make testacc_gcp_connection`) so they stay out of the shared `basic`/`post_process` CI
// lanes until the GCP APIs are part of a published clumio-go-sdk release. They require a
// beta-enabled Clumio backend and the CLUMIO_TEST_GCP_PROJECT_ID/CLUMIO_TEST_GCP_PROJECT_ID2
// environment variables in addition to the Clumio API credentials.
//
// The post-process resource has no Read implementation, so a refresh is a no-op and the state
// cannot drift; the plan checks therefore expect an empty plan after refresh. The WIF, service
// account and project-number/name fields are not produced by any Clumio resource (in real usage
// they are outputs of the deployed GCP Terraform template), so they are supplied as fixed literals
// in the test configuration while the token and project_id are wired from the connection.

//go:build gcp_connection

package clumio_post_process_gcp_connection_test

import (
	"fmt"
	"os"
	"testing"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Basic test of the clumio_post_process_gcp_connection resource. It tests the following scenarios:
//   - Create scenario for post process GCP connection and verifies that the plan was applied
//     properly.
//   - Updates the config for post process GCP connection and verifies that the resource will be
//     updated.
func TestAccResourcePostProcessGcpConnection(t *testing.T) {
	projectId := os.Getenv(common.ClumioTestGcpProjectId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestGcpConnectionPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourcePostProcessGcpConnection(baseUrl, projectId, "1.1"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_gcp_connection.test",
							plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_gcp_connection.test",
							plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_post_process_gcp_connection.test", "config_version", "1.1"),
				),
			},
			{
				Config: getTestAccResourcePostProcessGcpConnection(baseUrl, projectId, "2.0"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_gcp_connection.test",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"clumio_post_process_gcp_connection.test", "config_version", "2.0"),
				),
			},
		},
	})
}

// getTestAccResourcePostProcessGcpConnection returns the Terraform configuration for the
// clumio_post_process_gcp_connection resource.
func getTestAccResourcePostProcessGcpConnection(
	baseUrl string, projectId string, configVersion string) string {
	return fmt.Sprintf(
		testAccResourcePostProcessGcpConnection, baseUrl, projectId, configVersion)
}

// testAccResourcePostProcessGcpConnection is the Terraform configuration for the
// clumio_post_process_gcp_connection resource. The token and project_id are sourced from a real
// clumio_gcp_connection, while the remaining required fields use fixed literals as they are not
// produced by any Clumio resource.
const testAccResourcePostProcessGcpConnection = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_gcp_connection" "test_conn" {
  project_id  = "%s"
  description = "test_post_process"
  regions     = ["us-west1"]
}

resource "clumio_post_process_gcp_connection" "test" {
  project_id            = clumio_gcp_connection.test_conn.project_id
  project_name          = "test-project"
  project_number        = "123456789012"
  token                 = clumio_gcp_connection.test_conn.token
  service_account_email = "clumio-test@test-project.iam.gserviceaccount.com"
  wif_pool_id           = "clumio-test-pool"
  wif_provider_id       = "clumio-test-provider"
  config_version        = "%s"
  protect_gcs_version   = "1.0"
  regions               = ["us-west1"]
  region_configuration = [
    {
      region                       = "us-west1"
      inventory_bridge_bucket_name = "clumio-inventory-bridge-us-west1-test-project"
    }
  ]
  properties = {
    key1 = "val1"
    key2 = "val2"
  }
}
`
