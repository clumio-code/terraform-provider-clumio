// Copyright 2021. Clumio, Inc.

// Acceptance test for clumio_rule resource.
package clumio_policy_rule_test

import (
	"fmt"
	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioPolicyRule(t *testing.T) {
	policyName := "test_policy"
	policyTwoName := "test_policy_2"
	policyRuleName := "acceptance-test-policy-rule"
	policyRuleTwoName := "acceptance-test-policy-rule-2"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyRule(policyName, policyRuleName, policyRuleTwoName),
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

func getTestAccResourceClumioPolicyRule(policyName string,
	policyRuleName string, policyRuleTwoName string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioPolicyRule, baseUrl, policyName, policyName,
		policyRuleName, policyName, policyRuleTwoName, policyName)
}

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
	backup_window {
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
  condition = "{\"entity_type\":{\"$in\":[\"aws_ebs_volume\",\"aws_ec2_instance\"]}}"
}

resource "clumio_policy_rule" "test_policy_rule_2" {
  name = "%s"
  policy_id = clumio_policy.%s.id
  before_rule_id = clumio_policy_rule.test_policy_rule.id
  condition = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}}"
}

`


