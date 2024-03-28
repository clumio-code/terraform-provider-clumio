// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_post_process_aws_connection Terraform resource.
// Please view the README.md file for more information on how to run these tests.

//go:build post_process

package clumio_post_process_aws_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Basic test of the clumio_post_process_aws_connection resource. It tests the following scenarios:
//   - Create scenario for post process aws connection and verifies that the plan was
//     applied properly.
//   - Updates the config for post process aws connection and verifies that the resource will
//     be updated.
func TestAccResourcePostProcessAwsConnection(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId2)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourcePostProcessAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description", false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_aws_connection.test",
							plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_aws_connection.test",
							plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_aws_connection.test", "config_version",
						regexp.MustCompile("1.1")),
				),
			},
			{
				Config: getTestAccResourcePostProcessAwsConnection(baseUrl, accountNativeId,
					testAwsRegion, "test_description_updated", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_post_process_aws_connection.test",
							plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_aws_connection.test", "config_version",
						regexp.MustCompile("2.0")),
				),
			},
		},
	})
}

// testAccResourcePostProcessAwsConnection returns the Terraform configuration for the
// clumio_post_process_aws_connection resource.
func getTestAccResourcePostProcessAwsConnection(
	baseUrl string, accountId string, awsRegion string, description string, update bool) string {
	configVersion := "1.1"
	if update {
		configVersion = "2.0"
	}
	return fmt.Sprintf(testAccResourcePostProcessAwsConnection, baseUrl, accountId, awsRegion,
		description, configVersion)
}

// testAccResourcePostProcessAwsConnection is the Terraform configuration for the
// clumio_post_process_aws_connection resource.
const testAccResourcePostProcessAwsConnection = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
  account_native_id = "%s"
  aws_region = "%s"
  description = "%s"
}

resource "clumio_post_process_aws_connection" "test" {
  token = clumio_aws_connection.test_conn.token
  role_external_id = clumio_aws_connection.test_conn.role_external_id
  role_arn = "arn:aws:iam::${clumio_aws_connection.test_conn.account_native_id}:role/testRoleArn"
  account_id = clumio_aws_connection.test_conn.account_native_id
  region = clumio_aws_connection.test_conn.aws_region
  clumio_event_pub_id = "arn:aws:iam::${clumio_aws_connection.test_conn.account_native_id}:role/ev"
  config_version = "%s"
  discover_version = "4.1"
  protect_config_version = "19.2"
  protect_ebs_version = "20.1"
  protect_rds_version = "18.1"
  protect_ec2_mssql_version = "2.1"
  protect_warm_tier_version = "2.1"
  protect_warm_tier_dynamodb_version = "2.1"
  protect_dynamodb_version = "1.1"
  protect_s3_version = "2.1"
  properties = {
	key1 = "val1"
	key2 = "val2"
  }
}
`
