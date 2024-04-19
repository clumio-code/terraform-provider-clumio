// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_policy_rule Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_policy_rule_test

import (
	"fmt"
	"os"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_policy_rule datasource. It tests that the policy_rules matching the name
// and/or policy_id provided in the config are fetched and set in state.
func TestAccDataSourceClumioPolicyRule(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test where only policy_id are specified in the config.
			{
				Config: getTestDataSourceClumioPolicyRule(baseUrl, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy_rule.policy_rules",
						"policy_rules.#", "2"),
				),
			},
			// Test where only name are specified in the config.
			{
				Config: getTestDataSourceClumioPolicyRule(baseUrl, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy_rule.policy_rules",
						"policy_rules.#", "1"),
				),
			},
			// Test where both name and policy_id are specified in the config.
			{
				Config: getTestDataSourceClumioPolicyRule(baseUrl, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy_rule.policy_rules",
						"policy_rules.#", "1"),
				),
			},
		},
	})
}

// getTestDataSourceClumioPolicyRule returns the Terraform configuration for a basic
// clumio_policy_rule datasource.
func getTestDataSourceClumioPolicyRule(baseUrl string, includeName, includePolicyId bool) string {

	var policyId, datasourceName string

	if includeName {
		datasourceName = `name = "ds-acceptance-test-policy-rule1"`
	}
	if includePolicyId {
		policyId = `policy_id = clumio_policy.policy-rule-ds-test.id`
	}
	return fmt.Sprintf(testAccDataSourceClumioPolicyRule, baseUrl, datasourceName, policyId)
}

// testAccDataSourceClumioPolicyRule is the Terraform configuration for a basic clumio_policy_rule
// data source.
const testAccDataSourceClumioPolicyRule = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy" "policy-rule-ds-test" {
	name = "policy-rule-ds-test"
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

resource "clumio_policy_rule" "ds_test_policy_rule1" {
  name = "ds-acceptance-test-policy-rule1"
  policy_id = clumio_policy.policy-rule-ds-test.id
  before_rule_id = clumio_policy_rule.ds_test_policy_rule2.id
  condition = "{\"entity_type\":{\"$in\":[\"aws_ebs_volume\",\"aws_ec2_instance\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"Foo\", \"value\":\"Bar\"}}}"
}

resource "clumio_policy_rule" "ds_test_policy_rule2" {
  name = "ds-acceptance-test-policy-rule2"
  policy_id = clumio_policy.policy-rule-ds-test.id
  before_rule_id = ""
  condition = "{\"entity_type\":{\"$in\":[\"aws_ebs_volume\",\"aws_ec2_instance\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"Foo\", \"value\":\"Bar\"}}}"
}

data "clumio_policy_rule" "policy_rules" {
	depends_on = [ clumio_policy_rule.ds_test_policy_rule1, 
                   clumio_policy_rule.ds_test_policy_rule2]
	%s
	%s
}
`
