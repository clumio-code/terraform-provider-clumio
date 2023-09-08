// Copyright 2023. Clumio, Inc.

// Acceptance test for clumio_auto_user_provisioning_rule resource.
package clumio_auto_user_provisioning_rule_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClumioAutoUserProvisioningRule(t *testing.T) {
	autoUserProvisioningRuleName := "acceptance-test-auto-user-provisioning-rule"
	superAdminRole := "00000000-0000-0000-0000-000000000000"
	ouAdminRole := "10000000-0000-0000-0000-000000000000"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAutoUserProvisioningRule(
					autoUserProvisioningRuleName, superAdminRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule", "name",
						regexp.MustCompile(autoUserProvisioningRuleName)),
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule", "role_id",
						regexp.MustCompile(superAdminRole)),
				),
			},
			{
				Config: getTestAccResourceClumioAutoUserProvisioningRuleUpdate(
					autoUserProvisioningRuleName, ouAdminRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule", "name",
						regexp.MustCompile(autoUserProvisioningRuleName)),
					resource.TestMatchResourceAttr(
						"clumio_auto_user_provisioning_rule.test_auto_user_provisioning_rule", "role_id",
						regexp.MustCompile(ouAdminRole)),
				),
			},
		},
	})
}

func getTestAccResourceClumioAutoUserProvisioningRule(autoUserProvisioningRuleName string,
	roleId string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioAutoUserProvisioningRule, baseUrl,
		autoUserProvisioningRuleName, roleId)
}

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

func getTestAccResourceClumioAutoUserProvisioningRuleUpdate(autoUserProvisioningRuleName string,
	roleId string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioAutoUserProvisioningRuleUpdate, baseUrl,
		autoUserProvisioningRuleName, roleId)
}

const testAccResourceClumioAutoUserProvisioningRuleUpdate = `
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
