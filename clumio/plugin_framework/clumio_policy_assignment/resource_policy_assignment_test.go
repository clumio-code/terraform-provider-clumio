// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_policy_assignment Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_policy_assignment_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_policy_assignment resource. It tests the following scenario:
//   - Creates a policy assignment and verifies that the plan was applied properly.
func TestAccResourceClumioPolicyAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyAssignment("test-1"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy_assignment.test_policy_assignment",
							plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy_assignment.test_policy_assignment",
							plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// Tests that an external deletion of a clumio_policy_assignment resource leads to the resource needing
// to be re-created during the next plan. Following tests are included:
// - Deleting policy externally should recreate the resource.
// - Deleting protection externally group should recreate the resource.
// - Removing the policy from the protection group should recreate the resource.
func TestAccResourceClumioPolicyAssignmentRecreate(t *testing.T) {

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyAssignment("test-1"),
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
							"clumio_policy.test_policy", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(
							"clumio_policy_assignment.test_policy_assignment", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					// Delete the resource using the Clumio API after the plan is applied.
					common.DeletePolicy("clumio_policy.test_policy", true),
				),
				// This attribute is used to denote that the test expects that after the plan is
				// applied and a refresh is run, a non-empty plan is expected due to differences
				// from the state. Without this attribute set, the test would fail as it is unaware
				// that the resource was deleted externally.
				ExpectNonEmptyPlan: true,
			},
			{
				Config: getTestAccResourceClumioPolicyAssignment("test-2"),
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
							"clumio_protection_group.test_pg_policy_assignment",
							plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(
							"clumio_policy_assignment.test_policy_assignment",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					// Delete the resource using the Clumio API after the plan is applied.
					common.DeleteProtectionGroup(
						"clumio_protection_group.test_pg_policy_assignment", true),
				),
				// This attribute is used to denote that the test expects that after the plan is
				// applied and a refresh is run, a non-empty plan is expected due to differences
				// from the state. Without this attribute set, the test would fail as it is unaware
				// that the resource was deleted externally.
				ExpectNonEmptyPlan: true,
			},
			{
				Config: getTestAccResourceClumioPolicyAssignment("test-3"),
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
							"clumio_policy_assignment.test_policy_assignment", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					// Delete the resource using the Clumio API after the plan is applied.
					removePolicyFromProtectionGroup(
						"clumio_protection_group.test_pg_policy_assignment"),
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

// Tests that empty organizational_unit returns error
func TestAccResourceClumioPolicyAssignmentWithEmptyOU(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceClumioPolicyAssignmentEmptyOU,
					os.Getenv(common.ClumioApiBaseUrl)),
				ExpectError: regexp.MustCompile(
					"Attribute organizational_unit_id string length must be at least 1"),
			},
		},
	})
}

// getTestAccResourceClumioPolicyAssignment returns the Terraform configuration for a basic
// clumio_policy_assignment resource.
func getTestAccResourceClumioPolicyAssignment(suffix string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioPolicyAssignment, baseUrl, suffix)
}

// removePolicyFromProtectionGroup is a utility function to remove the policy from the protection
// group to simulate recreation of the policy_assignment_resource.
func removePolicyFromProtectionGroup(pgResName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		// retrieve the pg resource by name from state
		pgRes, ok := s.RootModule().Resources[pgResName]
		if !ok {
			return fmt.Errorf("Not found: %s", pgResName)
		}

		if pgRes.Primary.ID == "" {
			return fmt.Errorf("Protection Group ID is not set")
		}
		pgId := pgRes.Primary.ID

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
		pa := sdkclients.NewPolicyAssignmentClient(config)
		entityType := "protection_group"
		action := "unassign"
		entity := &models.AssignmentEntity{
			Id:         &pgId,
			ClumioType: &entityType,
		}
		assignmentInput := &models.AssignmentInputModel{
			Action: &action,
			Entity: entity,
		}
		res, apiErr := pa.SetPolicyAssignments(&models.SetPolicyAssignmentsV1Request{
			Items: []*models.AssignmentInputModel{
				assignmentInput,
			},
		})
		if apiErr != nil {
			return apiErr
		}
		if res == nil {
			return fmt.Errorf(common.NilErrorMessageDetail)
		}

		taskClient := sdkclients.NewTaskClient(config)

		// As creating a policy assignment is an asynchronous operation, the task ID
		// returned by the API is used to poll for the completion of the task.
		err := common.PollTask(
			context.Background(), taskClient, *res.TaskId, 300*time.Second, 5*time.Second)
		if err != nil {
			return err
		}
		return nil
	}
}

// testAccResourceClumioPolicyAssignment is the Terraform configuration for a basic
// clumio_policy_assignment resource.
const testAccResourceClumioPolicyAssignment = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg_policy_assignment"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_policy_assignment-%s"
  description = "test-description"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}

resource "clumio_policy" "test_policy" {
  name = "acceptance-test-policy-1234"
  operations {
	action_setting = "immediate"
	type = "protection_group_backup"
	slas {
		retention_duration {
			unit = "months"
			value = 3
		}
		rpo_frequency {
			unit = "days"
			value = 2
		}
	}
    advanced_settings {
		protection_group_backup {
			backup_tier = "cold"
		}
    }
  }
}

resource "clumio_policy_assignment" "test_policy_assignment" {
  entity_id = clumio_protection_group.test_pg_policy_assignment.id
  entity_type = "protection_group"
  policy_id = clumio_policy.test_policy.id
}
`

// testAccResourceClumioPolicyAssignmentEmptyOU is the Terraform configuration for creating a policy
// assignment with empty value for organizational_unit_id
const testAccResourceClumioPolicyAssignmentEmptyOU = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy_assignment" "test_policy_assignment" {
  entity_id = "some_pg_id"
  entity_type = "protection_group"
  policy_id = "some_policy_id"
  organizational_unit_id = ""
}
`
