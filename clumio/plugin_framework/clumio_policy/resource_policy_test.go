// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_policy Terraform resource. Please view the
// README.md file for more information on how to run these tests.

//go:build basic

package clumio_policy_test

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_policy resource. It tests the following scenarios:
//   - Creates a policy with fixed start time and no end time and verifies that the plan was applied
//     properly.
//   - Updates the policy with a new timezone and verifies that the resource will be updated.
//   - Updates the policy with a new start time and verifies that the resource will be updated.
//   - Updates the policy with both start and end time and verifies that the resource will be
//     updated.
//   - Updates the policy with a new start and end time and verifies that the resource will be
//     updated.
//   - Create a Policy for Secure Vault Lite and verifies that the plan was applied properly.
//   - Update the policy for Secure Vault Lite by adding a new SLA and verifies that the resource
//     will be updated.
//   - Creates a policy with hourly and minutely SLAs populated and verifies that the plan was
//     applied properly.
//   - Update the policy with different hourly and minutely SLAs and verifies that the resource will
//     be updated.
//   - Creates a policy with weekly SLA populated and verifies that the plan was applied properly.
//   - Update the policy with different weekly SLA and verifies that the resource will be updated.
//   - Creates a policy with backup region populated and verifies that the plan was applied properly.
//   - Update the policy with a different backup region and verifies that the resource will be
//     updated.
//   - Creates a policy without any operations and verifies that an error is returned.
func TestAccResourceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Basic create policy test.
				Config: getTestAccResourceClumioPolicyFixedStart(0),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.test_policy", plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.test_policy", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "timezone",
						regexp.MustCompile("UTC")),
				),
			},
			{
				// Update test for updating the timezone to US/Pacific
				Config: getTestAccResourceClumioPolicyFixedStart(1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.test_policy", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "timezone",
						regexp.MustCompile("US/Pacific")),
				),
			},
			{
				// Update test for updating the backup_window_tz start time to 05:00
				Config: getTestAccResourceClumioPolicyFixedStart(2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.test_policy", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "operations.1.backup_window_tz.0.start_time",
						regexp.MustCompile("05:00")),
				),
			},
			{
				// Update test for updating the backup_window_tz start time to 01:00 and
				// end time to 05:00
				Config: getTestAccResourceClumioPolicyWindow(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.test_policy", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "operations.1.backup_window_tz.0.start_time",
						regexp.MustCompile("01:00")),
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "operations.1.backup_window_tz.0.end_time",
						regexp.MustCompile("05:00")),
				),
			},
			{
				// Update test for updating the backup_window_tz start time to 03:00 and
				// end time to 07:00
				Config: getTestAccResourceClumioPolicyWindow(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.test_policy", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "operations.1.backup_window_tz.0.start_time",
						regexp.MustCompile("03:00")),
					resource.TestMatchResourceAttr(
						"clumio_policy.test_policy", "operations.1.backup_window_tz.0.end_time",
						regexp.MustCompile("07:00")),
				),
			},
			{
				// Create Policy test for Secure Vault Lite. Sets the advanced_setting
				// for aws_ebs_volume_backup -> backup_tier to lite.
				Config: getTestAccResourceClumioPolicySecureVaultLite(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.secure_vault_lite_success",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.secure_vault_lite_success",
						"operations.0.advanced_settings.0.aws_ebs_volume_backup.0.backup_tier",
						regexp.MustCompile("lite")),
					resource.TestMatchResourceAttr(
						"clumio_policy.secure_vault_lite_success", "operations.0.slas.#",
						regexp.MustCompile("1")),
				),
			},
			{
				// Update test for adding a new SLA the Secure Vault Lite policy.
				Config: getTestAccResourceClumioPolicySecureVaultLite(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.secure_vault_lite_success",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy.secure_vault_lite_success",
						"operations.0.advanced_settings.0.aws_ebs_volume_backup.0.backup_tier",
						regexp.MustCompile("lite")),
					resource.TestMatchResourceAttr(
						"clumio_policy.secure_vault_lite_success", "operations.0.slas.#",
						regexp.MustCompile("2")),
				),
			},
			{
				// Create policy test with hourly and minutely SLAs populated.
				Config: getTestAccResourceClumioPolicyHourlyMinutely(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.hourly_minutely_policy",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.0.slas.0.rpo_frequency.0.unit",
						regexp.MustCompile("hours")),
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.0.slas.0.rpo_frequency.0.value",
						regexp.MustCompile("4")),
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.1.slas.0.rpo_frequency.0.unit",
						regexp.MustCompile("minutes")),
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.1.slas.0.rpo_frequency.0.value",
						regexp.MustCompile("15")),
				),
			},
			{
				// Update policy test to update the hourly and minutely SLAs.
				Config: getTestAccResourceClumioPolicyHourlyMinutely(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.hourly_minutely_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.0.slas.0.rpo_frequency.0.unit",
						regexp.MustCompile("hours")),
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.0.slas.0.rpo_frequency.0.value",
						regexp.MustCompile("12")),
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.1.slas.0.rpo_frequency.0.unit",
						regexp.MustCompile("minutes")),
					resource.TestMatchResourceAttr("clumio_policy.hourly_minutely_policy",
						"operations.1.slas.0.rpo_frequency.0.value",
						regexp.MustCompile("30")),
				),
			},
			{
				// Create policy test with weekly SLA populated.
				Config: getTestAccResourceClumioPolicyWeekly(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.weekly_policy",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.weekly_policy",
						"operations.0.slas.0.rpo_frequency.0.unit",
						regexp.MustCompile("weeks")),
					resource.TestMatchResourceAttr("clumio_policy.weekly_policy",
						"operations.0.slas.0.rpo_frequency.0.offsets.0",
						regexp.MustCompile("1")),
				),
			},
			{
				// Update policy test for updating weekly SLA.
				Config: getTestAccResourceClumioPolicyWeekly(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.weekly_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.weekly_policy",
						"operations.0.slas.0.rpo_frequency.0.unit",
						regexp.MustCompile("weeks")),
					resource.TestMatchResourceAttr("clumio_policy.weekly_policy",
						"operations.0.slas.0.rpo_frequency.0.offsets.0",
						regexp.MustCompile("3")),
				),
			},
			{
				// Create policy test with Backup Region populated.
				Config: getTestAccResourceClumioPolicyBackupRegion(0),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.backup_region_policy",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.backup_region_policy",
						"operations.0.backup_aws_region", regexp.MustCompile("us-west-2")),
				),
			},
			{
				// Update policy test for updating Backup Region.
				Config: getTestAccResourceClumioPolicyBackupRegion(1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.backup_region_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("clumio_policy.backup_region_policy",
						"operations.0.backup_aws_region"),
				),
			},
			{
				// Test to check if an error is returned when backup_aws_region is specified as ""
				Config:      getTestAccResourceClumioPolicyBackupRegion(2),
				ExpectError: regexp.MustCompile(".*Error running apply.*"),
			},
		},
	})
}

// RDS compliance test of the clumio_policy resource. It tests the following scenarios:
//   - Creates a policy for RDS with the frozen backup tier and verifies that the plan was applied
//     properly.
//   - Updates the RDS compliance policy and verifies that the resource will be updated.
func TestRdsComplianceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyRDSCompliance(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.test_policy",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.test_policy",
						"operations.0.advanced_settings.0.aws_rds_resource_granular_backup.0.backup_tier",
						regexp.MustCompile("frozen")),
					resource.TestMatchResourceAttr("clumio_policy.test_policy",
						"operations.0.slas.0.retention_duration.0.value",
						regexp.MustCompile("31")),
				),
			},
			{
				Config: getTestAccResourceClumioPolicyRDSCompliance(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.test_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.test_policy",
						"operations.0.advanced_settings.0.aws_rds_resource_granular_backup.0.backup_tier",
						regexp.MustCompile("frozen")),
					resource.TestMatchResourceAttr("clumio_policy.test_policy",
						"operations.0.slas.0.retention_duration.0.value",
						regexp.MustCompile("28")),
				),
			},
		},
	})
}

// RDS PITR test of the clumio_policy resource. It tests the following scenarios:
//   - Creates a policy for RDS PITR and verifies that the plan was applied properly.
//   - Updates the policy into RDS Airgap and verifies that the resource will be updated.
//   - Updates the policy into RDS PITR with apply set to immediate for aws_rds_config_sync advanced
//     setting and verifies that the resource will be updated.
//   - Updates the policy into RDS PITR with apply set to maintenance_window for aws_rds_config_sync
//     advanced setting and verifies that the resource will be updated.
func TestRdsPitrClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// PITR only
				Config: getTestClumioPolicyRds(true, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.tf_rds_policy",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.tf_rds_policy",
						"operations.0.type", regexp.MustCompile("aws_rds_resource_aws_snapshot")),
					resource.TestMatchResourceAttr("clumio_policy.tf_rds_policy",
						"operations.0.slas.0.retention_duration.0.value", regexp.MustCompile("7")),
				),
			},
			{
				// Airgap only
				Config: getTestClumioPolicyRds(false, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.tf_rds_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.tf_rds_policy",
						"operations.0.type",
						regexp.MustCompile("aws_rds_resource_rolling_backup")),
					resource.TestMatchResourceAttr("clumio_policy.tf_rds_policy",
						"operations.0.slas.0.retention_duration.0.value",
						regexp.MustCompile("31")),
				),
			},
			{
				// RDS PITR immediate
				Config: getTestClumioPolicyRdsPitrAdv(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.tf_rds_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.tf_rds_policy",
						"operations.0.advanced_settings.0.aws_rds_config_sync.0.apply",
						regexp.MustCompile("immediate")),
				),
			},
			{
				// RDS PITR maintenance_window
				Config: getTestClumioPolicyRdsPitrAdv(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.tf_rds_policy",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.tf_rds_policy",
						"operations.0.advanced_settings.0.aws_rds_config_sync.0.apply",
						regexp.MustCompile("maintenance_window")),
				),
			},
		},
	})
}

// Tests that an external deletion of a clumio_policy resource leads to the resource needing to be
// re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioPolicyRecreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyWindow(false),
				Check: resource.ComposeTestCheckFunc(
					common.DeletePolicy("clumio_policy.test_policy", true),
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
							"clumio_policy.test_policy", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

// Test imports a Policy by ID and ensures that the import is successful.
func TestAccResourceClumioPolicyImport(t *testing.T) {

	// Return if it is not an acceptance test
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			resource.EnvTfAcc))
		return
	}

	clumiopf.UtilTestAccPreCheckClumio(t)
	id, err := createPolicyUsingSDK()
	if err != nil {
		t.Errorf("Error creating Policy using API: %v", err.Error())
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceClumioPolicyImport, os.Getenv(common.ClumioApiBaseUrl)),
				ImportState:   true,
				ResourceName:  "clumio_policy.test_policy",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported policy")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported policy has different ID. Expected: %v, Actual: %v",
							id, instStates[0].ID)
						return errors.New(errMsg)
					}
					return nil
				},
				ImportStatePersist: true,
				Destroy:            true,
			},
		},
	})
}

// Tests to check if creating a policy without operations returns error.
func TestAccResourceClumioPolicyErrorMissingOperations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceClumioPolicyEmptyOperations, os.Getenv(common.ClumioApiBaseUrl)),
				ExpectError: regexp.MustCompile("Invalid Block"),
			},
		},
	})
}

// Tests to check if creating a policy with empty organizational_unit_id returns error.
func TestAccResourceClumioPolicyErrorEmptyOU(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestClumioPolicyEmptyParams("organizational_unit_id"),
				ExpectError: regexp.MustCompile(
					"Attribute organizational_unit_id string length must be at least 1"),
			},
		},
	})
}

// Tests to check if creating a policy with empty timezone returns error.
func TestAccResourceClumioPolicyErrorEmptyTimezone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestClumioPolicyEmptyParams("timezone"),
				ExpectError: regexp.MustCompile(
					"Attribute timezone string length must be at least 1"),
			},
		},
	})
}

// Tests to check if creating a policy without operations returns error.
func TestAccResourceClumioPolicyErrorEmptyActivationStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getTestClumioPolicyEmptyParams("activation_status"),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

// Tests to check if creating a policy with child-level timezone works as expected.
func TestAccResourceClumioPolicyTimezone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyTimezone("US/Pacific", "US/Eastern"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction("clumio_policy.tf_child_timezone_policy",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.tf_child_timezone_policy",
						"operations.0.timezone",
						regexp.MustCompile("US/Pacific")),
					resource.TestMatchResourceAttr("clumio_policy.tf_child_timezone_policy",
						"operations.1.timezone",
						regexp.MustCompile("US/Eastern")),
				),
			},
			{
				Config: getTestAccResourceClumioPolicyTimezone("US/Eastern", "US/Pacific"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy.tf_child_timezone_policy", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("clumio_policy.tf_child_timezone_policy",
						"operations.0.timezone",
						regexp.MustCompile("US/Eastern")),
					resource.TestMatchResourceAttr("clumio_policy.tf_child_timezone_policy",
						"operations.1.timezone",
						regexp.MustCompile("US/Pacific")),
				),
			},
		},
	})
}

// getTestAccResourceClumioPolicyWindow returns the Terraform configuration for a clumio_policy resource
// containing a backup window.
func getTestAccResourceClumioPolicyWindow(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-policy-1234"
	timezone := "UTC"
	window := `
	backup_window_tz {
		start_time = "01:00"
		end_time = "05:00"
	}`
	if update {
		name = "acceptance-test-policy-4321"
		timezone = "US/Pacific"
		window = `
		backup_window_tz {
			start_time = "03:00"
			end_time = "07:00"
		}`
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, name, timezone, window)
}

// getTestAccResourceClumioPolicyFixedStart returns the Terraform configuration for a clumio_policy resource
// containing a fixed start time in the backup window.
func getTestAccResourceClumioPolicyFixedStart(update int) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-policy-1234"
	timezone := "UTC"
	window := `
	backup_window_tz {
		start_time = "01:00"
	}`
	if update == 1 {
		name = "acceptance-test-policy-4321"
		timezone = "US/Pacific"
		window = `
		backup_window_tz {
			start_time = "05:00"
		}`
	} else if update == 2 {
		window = `
		backup_window_tz {
			start_time = "05:00"
			end_time = ""
		}`
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, name, timezone, window)
}

// getTestAccResourceClumioPolicySecureVaultLite returns the Terraform configuration for a secure vault lite clumio_policy
// resource.
func getTestAccResourceClumioPolicySecureVaultLite(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "SecureVaultLite Test"
	sla := ``
	if update {
		sla = `
		slas {
			retention_duration {
				unit = "months"
				value = 3
			}
			rpo_frequency {
				unit = "months"
				value = 1
			}
		}`
	}
	return fmt.Sprintf(testAccResourceClumioPolicyVaultLite, baseUrl, name, sla)
}

// getTestAccResourceClumioPolicyHourlyMinutely returns the Terraform configuration for a clumio_policy resource
// containing hourly and minutely SLA.
func getTestAccResourceClumioPolicyHourlyMinutely(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Hourly & Minutely Policy Create"
	hourlySla := `
	slas {
		retention_duration {
			unit = "days"
			value = 15
		}
		rpo_frequency {
			unit = "hours"
			value = 4
		}
	}
	`
	minutelySla := `
	slas {
		retention_duration {
			unit = "days"
			value = 5
		}
		rpo_frequency {
			unit = "minutes"
			value = 15
		}
	}
	`
	if update {
		name = "Hourly & Minutely Policy Update"
		hourlySla = `
		slas {
			retention_duration {
				unit = "days"
				value = 15
			}
			rpo_frequency {
				unit = "hours"
				value = 12
			}
		}
		`
		minutelySla = `
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "minutes"
				value = 30
			}
		}
		`
	}
	return fmt.Sprintf(testAccResourceClumioPolicyHourlyMinutely, baseUrl, name, hourlySla, minutelySla)
}

// getTestAccResourceClumioPolicyWeekly returns the Terraform configuration for a clumio_policy resource
// containing a weekly SLA.
func getTestAccResourceClumioPolicyWeekly(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Weekly Policy Create"
	weeklySla := `
	slas {
		retention_duration {
			unit = "weeks"
			value = 4
		}
		rpo_frequency {
			unit = "weeks"
			value = 1
			offsets = [1]
		}
	}
	`
	if update {
		name = "Weekly Policy Update"
		weeklySla = `
		slas {
			retention_duration {
				unit = "weeks"
				value = 5
			}
			rpo_frequency {
				unit = "weeks"
				value = 1
				offsets = [3]
			}
		}
		`
	}
	return fmt.Sprintf(testAccResourceClumioPolicyWeekly, baseUrl, name, weeklySla)
}

// getTestAccResourceClumioPolicyBackupRegion returns the Terraform configuration for a clumio_policy resource
// containing a backup region.
func getTestAccResourceClumioPolicyBackupRegion(scenario int) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Backup Region Policy Create"
	timezone := "UTC"
	region := `
	backup_aws_region = "us-west-2"`
	if scenario == 1 {
		name = "Backup Region Policy Update"
		region = `` // valid as the region is optional
	} else if scenario == 2 {
		name = "Backup Region Policy Update 2"
		region = `
	backup_aws_region = ""` // invalid as empty region is not allowed as request.
	}
	return fmt.Sprintf(testAccResourceClumioPolicyBackupRegion, baseUrl, name, timezone, region)
}

// getTestAccResourceClumioPolicyRDSCompliance returns the Terraform configuration for a clumio_policy resource to
// support RDS backup for compliance.
func getTestAccResourceClumioPolicyRDSCompliance(update bool) string {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Rds Compliance Policy Create"
	slas := `
	slas {
		retention_duration {
			unit  = "days"
			value = 31
		}
		rpo_frequency {
			unit  = "days"
			value = 7
		}
	}
	advanced_settings {
		aws_rds_resource_granular_backup {
			backup_tier = "frozen"
		}
	}
	`
	if update {
		name = "Rds Compliance Policy Update"
		slas = `
		slas {
			retention_duration {
				unit  = "days"
				value = 28
			}
			rpo_frequency {
				unit  = "days"
				value = 7
			}
		}
		advanced_settings {
			aws_rds_resource_granular_backup {
				backup_tier = "frozen"
			}
		}
		`
	}
	return fmt.Sprintf(testAccResourceClumioRdsPolicy, baseUrl, name, slas)
}

// getTestClumioPolicyRds returns the Terraform configuration for a clumio_policy resource to
// support RDS backup.
func getTestClumioPolicyRds(pitr bool, airgap bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf-rds-policy"
	operations := ""
	// TODO: add advanced settings on it.
	pitrTemplate := `
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_aws_snapshot"
		slas {
			retention_duration {
				unit  = "days"
				value = 7
			}
			rpo_frequency {
				unit  = "days"
				value = 1
			}
		}
	}`
	airgapTemplate := `
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_rolling_backup"
		slas {
			retention_duration {
				unit  = "days"
				value = 31
			}
			rpo_frequency {
				unit  = "days"
				value = 7
			}
		}
	}`

	if pitr {
		operations += pitrTemplate
		name += "-pitr"
	}
	if airgap {
		operations += airgapTemplate
		name += "-airgap"
	}
	return fmt.Sprintf(testClumioPolicyRdsPolicyTemplate, baseUrl, name, operations)
}

// getTestClumioPolicyRdsPitrAdv returns the Terraform configuration for a clumio_policy resource
// to support RDS PITR backup with advanced settings.
func getTestClumioPolicyRdsPitrAdv(immediate bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf-rds-pitr-adv-policy"
	rdsPitrConfigAdv := "immediate"
	if !immediate {
		rdsPitrConfigAdv = "maintenance_window"
	}
	operations := fmt.Sprintf(`
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_aws_snapshot"
		slas {
			retention_duration {
				unit  = "days"
				value = 7
			}
			rpo_frequency {
				unit  = "days"
				value = 1
			}
		}
		advanced_settings {
			aws_rds_config_sync {
				apply = "%s"
			}
		}
	}`, rdsPitrConfigAdv)
	return fmt.Sprintf(testClumioPolicyRdsPolicyTemplate, baseUrl, name, operations)
}

// getTestAccResourceClumioPolicyTimezone returns the Terraform configuration for a clumio_policy resource
// containing a child-level timezone and a backup window.
func getTestAccResourceClumioPolicyTimezone(timezone1, timezone2 string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioPolicyTimezone, baseUrl, timezone1, timezone2)
}

// getTestClumioPolicyEmptyParams returns the Terraform configuration for a clumio_policy resource with
// organizational_unit_id set to empty string.
func getTestClumioPolicyEmptyParams(param string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf-acceptance-test-policy-empty-params"
	emptyParam := fmt.Sprintf("%s = \"\"", param)
	return fmt.Sprintf(testAccResourceClumioPolicyEmptyParams, baseUrl, name, emptyParam)
}

// createPolicyUsingSDK creates a policy using the Clumio API.
func createPolicyUsingSDK() (string, error) {
	clumioApiToken := os.Getenv(common.ClumioApiToken)
	clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
	clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
	client := &common.ApiClient{
		ClumioConfig: sdkconfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		},
	}
	pd := sdkclients.NewPolicyDefinitionClient(client.ClumioConfig)
	name := "acceptance-test-import"
	timezone := "UTC"
	actionSetting := "immediate"
	clumioType := "aws_ebs_volume_backup"
	unit := "days"
	retValue := int64(5)
	rpoValue := int64(1)
	slas := []*models.BackupSLA{
		{
			RetentionDuration: &models.RetentionBackupSLAParam{
				Unit:  &unit,
				Value: &retValue,
			},
			RpoFrequency: &models.RPOBackupSLAParam{
				Unit:  &unit,
				Value: &rpoValue,
			},
		},
	}
	operations := []*models.PolicyOperationInput{
		{
			ActionSetting: &actionSetting,
			Slas:          slas,
			ClumioType:    &clumioType,
		},
	}
	res, apiErr := pd.CreatePolicyDefinition(&models.CreatePolicyDefinitionV1Request{
		Name:       &name,
		Operations: operations,
		Timezone:   &timezone,
	})
	if apiErr != nil {
		return "", apiErr
	}
	return *res.Id, nil
}

// testAccResourceClumioPolicy is the Terraform configuration for a basic clumio_policy resource.
const testAccResourceClumioPolicy = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
	name = "%s"
	timezone = "%s"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
		%s
	}
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
`

// testAccResourceClumioPolicyTimezone is the Terraform configuration for a testing child-level timezone
// of clumio_policy resource.
const testAccResourceClumioPolicyTimezone = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "tf_child_timezone_policy" {
	name = "test_child_timezone_policy"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
		backup_window_tz {
			start_time = "05:00"
			end_time = ""
		}
		timezone = "%s"
	}
	operations {
		action_setting = "immediate"
		type = "aws_ec2_instance_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 10
			}
			rpo_frequency {
				unit = "days"
				value = 2
			}
		}
		backup_window_tz {
			start_time = "05:00"
			end_time = "10:00"
		}
		timezone = "%s"
	}
}
`

// testAccResourceClumioPolicyVaultLite is the Terraform configuration for a secure vault lite
// clumio_policy resource.
const testAccResourceClumioPolicyVaultLite = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "secure_vault_lite_success" {
  name = "%s"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 30
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
		%s
		advanced_settings {
			aws_ebs_volume_backup {
				backup_tier = "lite"
			}
		}
	}
}
`

// testAccResourceClumioPolicyHourlyMinutely is the Terraform configuration for a clumio_policy
// resource with hourly/minutely SLA configured.
const testAccResourceClumioPolicyHourlyMinutely = `
provider clumio{
	clumio_api_base_url = "%s"
}
resource "clumio_policy" "hourly_minutely_policy" {
	name = "%s"
	operations {
		action_setting = "immediate"
		type = "ec2_mssql_database_backup"
		%s
		advanced_settings {
			ec2_mssql_database_backup {
				alternative_replica = "sync_secondary"
				preferred_replica = "primary"
			}
		}
	}	
	operations {
		action_setting = "immediate"
		type = "ec2_mssql_log_backup"
		%s
		advanced_settings {
			ec2_mssql_log_backup {
				alternative_replica = "sync_secondary"
				preferred_replica = "primary"
			}
		}
	}
}
`

// testAccResourceClumioPolicyWeekly is the Terraform configuration for a clumio_policy resource
// with weekly SLA configured.
const testAccResourceClumioPolicyWeekly = `
provider clumio{
	clumio_api_base_url = "%s"
}
resource "clumio_policy" "weekly_policy" {
	name = "%s"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		%s
	}
}
`

// testClumioPolicyRdsPolicyTemplate is the Terraform configuration for a clumio_policy resource
// for RDS backup.
const testClumioPolicyRdsPolicyTemplate = `
provider clumio{
	clumio_api_base_url = "%s"
}
resource "clumio_policy" "tf_rds_policy" {
	name = "%s"
	%s
}
`

// testAccResourceClumioRdsPolicy is the Terraform configuration for a clumio_policy resource
// supporting RDS granular and snapshot backups.
const testAccResourceClumioRdsPolicy = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
	name = "%s"
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_granular_backup"
		%s
	}
}
`

// testAccResourceClumioPolicyImport is the Terraform configuration for importing a clumio_policy
// resource.
const testAccResourceClumioPolicyImport = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
	name = "acceptance-test-import"
	timezone = "UTC"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
	}
}
`

// testAccResourceClumioPolicyEmptyOperations is the Terraform configuration for a clumio_policy
// without any operations.
const testAccResourceClumioPolicyEmptyOperations = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
	name = "acceptance-test-import"
	timezone = "UTC"
}
`

// testAccResourceClumioPolicyBackupRegion is the Terraform configuration for a clumio_policy
// resource with backup_aws_region configured.
const testAccResourceClumioPolicyBackupRegion = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "backup_region_policy" {
	name = "%s"
	timezone = "%s"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
		%s
	}
}
`

// testAccResourceClumioPolicyEmptyParams is the Terraform configuration for a clumio_policy resource
// with organizational_unit_id or activation_status set to empty string.
const testAccResourceClumioPolicyEmptyParams = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "backup_region_policy" {
	name = "%s"
	%s
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
	}
}
`
