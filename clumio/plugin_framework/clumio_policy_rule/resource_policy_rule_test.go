// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_policy_rule Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_policy_rule_test

import (
	"context"
	"errors"
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

// Basic test of the clumio_policy_rule resource. It tests the following scenarios:
//   - Creates two policy rules, one with before_rule_id set and verifies that the plan was applied properly.
//   - Updates the policy_rules and verifies that the resource will be updated.
func TestAccResourceClumioPolicyRule(t *testing.T) {

	// Define the policy and policy rule names
	policyName := "test_policy"
	policyTwoName := "test_policy_2"
	policyRuleName := "acceptance-test-policy-rule"
	policyRuleTwoName := "acceptance-test-policy-rule-2"

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyRule(policyName, policyRuleName, policyRuleTwoName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy_rule.test_policy_rule", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(
							"clumio_policy_rule.test_policy_rule_2", plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy_rule.test_policy_rule", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(
							"clumio_policy_rule.test_policy_rule_2", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy_rule.test_policy_rule", "name",
						regexp.MustCompile(policyRuleName)),
					resource.TestMatchResourceAttr(
						"clumio_policy_rule.test_policy_rule_2", "name",
						regexp.MustCompile(policyRuleTwoName)),
				),
			},
			{
				Config: getTestAccResourceClumioPolicyRule(policyTwoName, policyRuleName, policyRuleTwoName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_policy_rule.test_policy_rule", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(
							"clumio_policy_rule.test_policy_rule_2", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_policy_rule.test_policy_rule", "name",
						regexp.MustCompile(policyRuleName)),
					resource.TestMatchResourceAttr(
						"clumio_policy_rule.test_policy_rule_2", "name",
						regexp.MustCompile(policyRuleTwoName)),
				),
			},
		},
	})
}

// Test imports a policy rule by ID and ensures that the import is successful.
func TestAccResourceClumioPolicyRuleImport(t *testing.T) {

	// Return if it is not an acceptance test
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			resource.EnvTfAcc))
		return
	}

	// Retrieve the environment variables required for the test.
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	// Create the policy rule to import using the Clumio API.
	clumiopf.UtilTestAccPreCheckClumio(t)
	policy_id, id, err := createPolicyRuleUsingSDK()
	if err != nil {
		t.Errorf("Error creating policy rule using API: %v", err.Error())
	}

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceClumioPolicyRuleImport, baseUrl, policy_id),
				ImportState:   true,
				ResourceName:  "clumio_policy_rule.test_policy_rule",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported policy rule")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported policy rule has different ID. Expected: %v, Actual: %v",
							id, instStates[0].ID)
						return errors.New(errMsg)
					}
					return nil
				},
				ImportStatePersist: true,
				Destroy:            true,
			},
		},
		CheckDestroy: common.DeletePolicy(policy_id, false),
	})
}

// createPolicyRuleUsingSDK creates a policy and policy rule using the Clumio API. This
// is required simulate importing an existing policy_rule. Since the policy_rule requires
// the policy_id, first the policy is created and then using the policy_id, the
// policy_rule is created.
func createPolicyRuleUsingSDK() (string, string, error) {

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
	pd := sdkclients.NewPolicyDefinitionClient(config)
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
	policyRes, apiErr := pd.CreatePolicyDefinition(&models.CreatePolicyDefinitionV1Request{
		Name:       &name,
		Operations: operations,
		Timezone:   &timezone,
	})
	if apiErr != nil {
		return "", "", apiErr
	}

	policyRules := sdkclients.NewPolicyRuleClient(config)
	condition := "{\"entity_type\":{\"$in\":[\"aws_ebs_volume\",\"aws_ec2_instance\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"Foo\", \"value\":\"Bar\"}}}"
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: policyRes.Id,
		},
	}
	beforeRuleId := ""
	priority := &models.RulePriority{
		BeforeRuleId: &beforeRuleId,
	}
	res, apiErr := policyRules.CreatePolicyRule(&models.CreatePolicyRuleV1Request{
		Action:    action,
		Condition: &condition,
		Name:      &name,
		Priority:  priority,
	})
	if apiErr != nil {
		return "", "", apiErr
	}

	taskClient := sdkclients.NewTaskClient(config)
	// As creating a policy rule is an asynchronous operation, the task ID
	// returned by the API is used to poll for the completion of the task.
	err := common.PollTask(
		context.Background(), taskClient, *res.TaskId, 300*time.Second, 5*time.Second)
	if err != nil {
		return "", "", err
	}
	return *policyRes.Id, *res.Rule.Id, nil
}

// getTestAccResourceClumioPolicyRule returns the Terraform configuration for a basic
// clumio_policy_rule resource.
func getTestAccResourceClumioPolicyRule(policyName string,
	policyRuleName string, policyRuleTwoName string) string {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioPolicyRule, baseUrl, policyName, policyName,
		policyRuleName, policyName, policyRuleTwoName, policyName)
}

// testAccResourceClumioPolicyRule is the Terraform configuration for a basic clumio_policy_rule
// resource.
const testAccResourceClumioPolicyRule = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy" "%s" {
 name = "%s"
 activation_status = "activated"
 operations {
	action_setting = "window"
	type = "aws_ebs_volume_backup"
	backup_window_tz {
		start_time = "08:00"
		end_time = "20:00"
	}
	slas {
		retention_duration {
			unit = "days"
			value = 1
		}
		rpo_frequency {
			unit = "days"
			value = 1
		}
	}
 }
}

resource "clumio_policy_rule" "test_policy_rule" {
  name = "%s"
  policy_id = clumio_policy.%s.id
  before_rule_id = ""
  condition = "{\"entity_type\":{\"$in\":[\"aws_ebs_volume\",\"aws_ec2_instance\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"Foo\", \"value\":\"Bar\"}}}"
}

resource "clumio_policy_rule" "test_policy_rule_2" {
  name = "%s"
  policy_id = clumio_policy.%s.id
  before_rule_id = clumio_policy_rule.test_policy_rule.id
  condition = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_tag\":{\"$eq\":{\"key\":\"Foo\", \"value\":\"Bar\"}}}"
}
`

// testAccResourceClumioPolicyRuleImport is the Terraform configuration which is used to simulate
// the importing an existing policy rule.
const testAccResourceClumioPolicyRuleImport = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy_rule" "test_policy_rule" {
  name = "acceptance-test-import"
  policy_id = "%s"
  before_rule_id = ""
  condition = "{\"entity_type\":{\"$in\":[\"aws_ebs_volume\",\"aws_ec2_instance\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"Foo\", \"value\":\"Bar\"}}}"
}
`
