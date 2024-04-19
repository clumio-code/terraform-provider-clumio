// Copyright 2024. Clumio, Inc.

// This files holds acceptance tests for the clumio_policy Terraform datasource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_policy_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_policy datasource. It tests that the policies matching the name and/or
// opration_types and/or activation_status provided in the config are fetched and set in state.
func TestAccDataSourceClumioPolicy(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClumioPolicy(baseUrl, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy.policies",
						"policies.#", "3"),
					resource.TestMatchResourceAttr("data.clumio_policy.policies",
						"policies.0.timezone", regexp.MustCompile("UTC")),
					resource.TestMatchResourceAttr("data.clumio_policy.policies",
						"policies.1.timezone", regexp.MustCompile("UTC")),
					resource.TestMatchResourceAttr("data.clumio_policy.policies",
						"policies.2.timezone", regexp.MustCompile("UTC")),
				),
			},
			// Test where name and operation_types are specified in the config.
			{
				Config: getTestDataSourceClumioPolicy(baseUrl, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy.policies",
						"policies.#", "2"),
				),
			},
			// Test where name and activation_status are specified in the config.
			{
				Config: getTestDataSourceClumioPolicy(baseUrl, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy.policies",
						"policies.#", "2"),
				),
			},
			// Test where name, activation_status and operation_types are specified in the config.
			{
				Config: getTestDataSourceClumioPolicy(baseUrl, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.clumio_policy.policies",
						"policies.#", "1"),
				),
			},
		},
	})
}

// Test to validate that an error is returned if none of name, operation_types and activation_status
// are specified in the config.
func TestEmptyPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test error scenario where none of the optional schema attributes are specified in the
			// config.
			{
				Config:      testAccEmptyDataSourceClumioPolicy,
				ExpectError: regexp.MustCompile(".*Missing Attribute Configuration.*"),
			},
		},
	})
}

// getTestDataSourceClumioPolicy returns the Terraform configuration for a basic clumio_policy
// datasource.
func getTestDataSourceClumioPolicy(baseUrl string,
	includeName, includeOperationTypes, includeActivationStatus bool) string {

	name := "datasource-test-policy"
	name1 := fmt.Sprintf(`%s1`, name)
	name2 := fmt.Sprintf(`%s2`, name)
	name3 := fmt.Sprintf(`%s3`, name)

	var operationType, activationStatus, datasourceName string
	if includeName {
		datasourceName = fmt.Sprintf(`name = "%s"`, name)
	}
	if includeOperationTypes {
		operationType = fmt.Sprintf(`operation_types = ["%s"]`, "protection_group_backup")
	}
	if includeActivationStatus {
		activationStatus = fmt.Sprintf(`activation_status = "%s"`, "activated")
	}
	return fmt.Sprintf(testAccDataSourceClumioPolicy, baseUrl, name1, name2, name3, datasourceName,
		operationType, activationStatus)
}

// testAccDataSourceClumioPolicy is the Terraform configuration for a basic clumio_policy datasource.
const testAccDataSourceClumioPolicy = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy" "datasource-test-policy1" {
	name = "%s"
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

resource "clumio_policy" "datasource-test-policy2" {
	name = "%s"
	timezone = "UTC"
	activation_status = "deactivated"
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

resource "clumio_policy" "datasource-test-policy3" {
	name = "%s"
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

data "clumio_policy" "policies" {
	depends_on = [ clumio_policy.datasource-test-policy1, 
                   clumio_policy.datasource-test-policy2,
                   clumio_policy.datasource-test-policy3]
	%s
	%s
	%s
}
`

// testAccDataSourceClumioPolicy is the Terraform configuration for a clumio_policy datasource where
// all of the optional attributes are not specified.
const testAccEmptyDataSourceClumioPolicy = `
provider clumio{
}

data "clumio_policy" "policies" {}
`
