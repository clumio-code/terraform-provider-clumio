// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_aws_manual_connection_resources
// Terraform datasource. Please view the README.md file for more information on how to run these
// tests.

//go:build manual_connection

package clumio_aws_manual_connection_resources_test

import (
	"fmt"
	"os"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_aws_manual_connection_resources datasource. It tests whether the
// resources field is returned and set in the state properly.
func TestAwsManualConnectionResources(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAccountId := os.Getenv(common.ClumioTestAwsAccountId2)
	testAwsRegion := os.Getenv(common.AwsRegion)
	testAssetTypes := map[string]bool{
		"EBS":      true,
		"S3":       true,
		"RDS":      true,
		"DynamoDB": true,
		"EC2MSSQL": true,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			clumioPf.UtilTestAccPreCheckClumio(t)
		},
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestClumioAwsManualConnectionResources(
					baseUrl, testAccountId, testAwsRegion, testAssetTypes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.clumio_aws_manual_connection_resources.test_get_resources",
						"resources"),
				),
			},
		},
	})
}

// getTestClumioAwsManualConnectionResources returns the Terraform configuration for a basic
// clumio_aws_manual_connection_resources resource.
func getTestClumioAwsManualConnectionResources(
	baseUrl string, accountId string, awsRegion string,
	testAssetTypes map[string]bool) string {
	return fmt.Sprintf(testResourceClumioAwsManualConnection,
		baseUrl,
		accountId,
		awsRegion,
		testAssetTypes["EBS"],
		testAssetTypes["RDS"],
		testAssetTypes["DynamoDB"],
		testAssetTypes["S3"],
		testAssetTypes["EC2MSSQL"],
	)
}

// testResourceClumioAwsManualConnection is the Terraform configuration for a basic
// clumio_aws_manual_connection_resources resource.
const testResourceClumioAwsManualConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn_manual" {
  account_native_id = "%s"
  aws_region = "%s"
  description = "test connection for manual resources"
}

data "clumio_aws_manual_connection_resources" "test_get_resources" {
	account_native_id = clumio_aws_connection.test_conn_manual.account_native_id
	aws_region = clumio_aws_connection.test_conn_manual.aws_region
	asset_types_enabled = {
		ebs = %t
		rds = %t
		ddb = %t
		s3 = %t
		mssql = %t
	}
}
`
