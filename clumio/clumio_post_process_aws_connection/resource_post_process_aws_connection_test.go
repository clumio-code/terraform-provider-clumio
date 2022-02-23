// Copyright 2021. Clumio, Inc.

// Acceptance test for resource_post_rprocess_aws_connection.
package clumio_post_process_aws_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePostProcessAwsConnection(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourcePostProcessAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_aws_connection.test", "account_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
			{
				Config: getTestAccResourcePostProcessAwsConnectionUpdate(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_post_process_aws_connection.test", "account_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
		},
	})
}

func getTestAccResourcePostProcessAwsConnection(
	baseUrl string, accountId string, awsRegion string, description string) string {
	return fmt.Sprintf(testAccResourcePostProcessAwsConnection, baseUrl, accountId,
		awsRegion, description)
}

const testAccResourcePostProcessAwsConnection = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
  account_native_id = %s
  aws_region = "%s"
  description = "%s"
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  services_enabled = ["discover", "protect"]
}

resource "clumio_post_process_aws_connection" "test" {
  token = clumio_aws_connection.test_conn.token
  role_external_id = "role_external_id_${clumio_aws_connection.test_conn.token}"
  role_arn = "testRoleArn"
  account_id = clumio_aws_connection.test_conn.account_native_id
  region = clumio_aws_connection.test_conn.aws_region
  clumio_event_pub_id = "test_event_pub_id"
  config_version = "1.1"
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

func getTestAccResourcePostProcessAwsConnectionUpdate(
	baseUrl string, accountId string, awsRegion string, description string) string {
	return fmt.Sprintf(testAccResourcePostProcessAwsConnectionUpdate, baseUrl, accountId,
		awsRegion, description)
}

const testAccResourcePostProcessAwsConnectionUpdate = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
  account_native_id = %s
  aws_region = "%s"
  description = "%s"
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  services_enabled = ["discover", "protect"]
}

resource "clumio_post_process_aws_connection" "test" {
  token = clumio_aws_connection.test_conn.token
  role_external_id = "role_external_id_${clumio_aws_connection.test_conn.token}"
  role_arn = "testRoleArn"
  account_id = clumio_aws_connection.test_conn.account_native_id
  region = clumio_aws_connection.test_conn.aws_region
  clumio_event_pub_id = "test_event_pub_id"
  config_version = "1.2"
  discover_version = "4.1"
  protect_config_version = "19.2"
  protect_ebs_version = "20.1"
  protect_rds_version = "18.1"
  protect_ec2_mssql_version = "2.1"
  protect_warm_tier_version = "2.1"
  protect_warm_tier_dynamodb_version = "2.1"
  protect_dynamodb_version = "1.1"
  protect_s3_version = "2.1"
}
`
