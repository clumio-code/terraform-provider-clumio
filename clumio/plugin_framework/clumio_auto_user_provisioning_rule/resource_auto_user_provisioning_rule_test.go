// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_auto_user_provisioning_rule Terraform resource.
// Please view the README.md file for more information on how to run these tests.

//go:build sso

package clumio_auto_user_provisioning_rule_test

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	autoUserProvisioningRules "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_rules"
	"github.com/clumio-code/clumio-go-sdk/models"
	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_auto_user_provisioning_rule resource. It tests the following scenarios:
//   - Creates an auto user provisioning rule and verifies that the plan was applied properly.
//   - Updates the auto user provisioning rule and verifies that the resource will be updated.
func TestAccClumioAutoUserProvisioningRule(t *testing.T) {
	autoUserProvisioningRuleName := "acceptance-test-auto-user-provisioning-rule"
	superAdminRole := "00000000-0000-0000-0000-000000000000"
	ouAdminRole := "10000000-0000-0000-0000-000000000000"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningRule(
					autoUserProvisioningRuleName, superAdminRole),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
							plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
							plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
						"name", regexp.MustCompile(autoUserProvisioningRuleName)),
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
						"role_id", regexp.MustCompile(superAdminRole)),
				),
				SkipFunc: clumioPf.SkipIfSSONotConfigured,
			},
			{
				Config: getTestAccResourceClumioAutoUserProvisioningRule(
					autoUserProvisioningRuleName, ouAdminRole),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
						"name",
						regexp.MustCompile(autoUserProvisioningRuleName)),
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
						"role_id",
						regexp.MustCompile(ouAdminRole)),
				),
				SkipFunc: clumioPf.SkipIfSSONotConfigured,
			},
		},
	})
}

// Tests that an external deletion of a clumio_auto_user_provisioning_rule resource leads to the
// resource needing to be re-created during the next plan. NOTE the Check function below as it is
// utilized to delete the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioAutoUserProvisioningRuleRecreate(t *testing.T) {

	autoUserProvisioningRuleName := "acceptance-test-auto-user-provisioning-rule"
	superAdminRole := "00000000-0000-0000-0000-000000000000"

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningRule(
					autoUserProvisioningRuleName, superAdminRole),
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
							"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
							plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
						"name",
						regexp.MustCompile(autoUserProvisioningRuleName)),
					// Delete the resource using the Clumio API after the plan is applied.
					deleteAutoUserProvisioningRule(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule"),
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

// Test imports an auto user provisioning rule by ID and ensures that the import is successful.
func TestAccResourceClumioAutoUserProvisioningRuleImport(t *testing.T) {

	autoUserProvisioningRuleName := "acceptance-test-auto-user-provisioning-rule"
	superAdminRole := "00000000-0000-0000-0000-000000000000"

	// Create the auto user provisioning rule to import using the Clumio API.
	clumioPf.UtilTestAccPreCheckClumio(t)
	id, err := createAutoUserProvisioningRoleUsingSDK()
	if err != nil {
		t.Errorf("Error creating auto user provisioning rule using API: %v", err.Error())
	}

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningRule(
					autoUserProvisioningRuleName, superAdminRole),
				ImportState:   true,
				ResourceName:  "clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported " +
							"auto user provisioning rule")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported auto user provisioning rule has different ID."+
								" Expected: %v, Actual: %v",
							id, instStates[0].ID)
						return errors.New(errMsg)
					}
					return nil
				},
				ImportStatePersist: true,
				Destroy:            true,
				SkipFunc:           clumioPf.SkipIfSSONotConfigured,
			},
		},
	})
}

// createAutoUserProvisioningRoleUsingSDK creates an auto user provisioning role using the Clumio API.
func createAutoUserProvisioningRoleUsingSDK() (string, error) {

	// Return if SSO is not configured.
	if strings.ToLower(os.Getenv(common.ClumioTestIsSSOConfigured)) != "true" {
		return "", nil
	}
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
	aupRules := autoUserProvisioningRules.NewAutoUserProvisioningRulesV1(config)
	condition := "{\"user.groups\":{\"$in\":[\"Group1\",\"Group2\"]}}"
	name := "acceptance-test-auto-user-provisioning-rule"
	orgUnitId := "00000000-0000-0000-0000-000000000000"
	roleId := "00000000-0000-0000-0000-000000000000"
	provision := &models.RuleProvision{
		OrganizationalUnitIds: []*string{
			&orgUnitId,
		},
		RoleId: &roleId,
	}
	res, apiErr := aupRules.CreateAutoUserProvisioningRule(
		&models.CreateAutoUserProvisioningRuleV1Request{
			Condition: &condition,
			Name:      &name,
			Provision: provision,
		})
	if apiErr != nil {
		return "", apiErr
	}
	return *res.RuleId, nil
}

// deleteAutoUserProvisioningRule returns a function that deletes an auto user provisioning rule
// using the Clumio API with information from the Terraform state. It is used to intentionally cause
// a difference between the Terraform state and the actual state of the resource in the backend.
func deleteAutoUserProvisioningRule(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		// Create a Clumio API client and delete the auto user provisioning rule.
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
		aupRules := autoUserProvisioningRules.NewAutoUserProvisioningRulesV1(config)
		_, apiErr := aupRules.DeleteAutoUserProvisioningRule(rs.Primary.ID)
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

// getTestAccResourceClumioAutoUserProvisioningRule returns the Terraform configuration for a basic
// clumio_auto_user_provisioning_rule resource.
func getTestAccResourceClumioAutoUserProvisioningRule(
	autoUserProvisioningRuleName string, roleId string) string {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioAutoUserProvisioningRule, baseUrl,
		autoUserProvisioningRuleName, roleId)
}

// testAccResourceClumioAutoUserProvisioningRule is the Terraform configuration for a basic
// clumio_auto_user_provisioning_rule resource.
const testAccResourceClumioAutoUserProvisioningRule = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_auto_user_provisioning_rule" "test_auto_user_provisioning_rule" {
  name = "%s"
  condition = "{\"user.groups\":{\"$in\":[\"Group1\",\"Group2\"]}}"
  role_id = "%s"
  organizational_unit_ids = ["00000000-0000-0000-0000-000000000000"]
}

`
