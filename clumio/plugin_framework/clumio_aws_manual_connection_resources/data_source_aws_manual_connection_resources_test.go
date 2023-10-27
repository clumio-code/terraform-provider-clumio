// Copyright 2023. Clumio, Inc.
//
// Acceptance test for aws_manual_connection_resources
package clumio_aws_manual_connection_resources_test

import (
	"fmt"
	"os"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAwsManualConnectionResources(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAccountId := os.Getenv(common.ClumioTestAwsAccountId)
	testAwsRegion := os.Getenv(common.AwsRegion)
	testAssetTypes := map[string]bool{
		"EBS": true,
		"S3": true,
		"RDS": true,
		"DynamoDB": true,
		"EC2MSSQL": true,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { 
			clumio_pf.UtilTestAccPreCheckClumio(t) 
		},
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestClumioAwsManualConnectionResources(
					baseUrl, testAccountId, testAwsRegion, testAssetTypes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(
						"clumio_aws_manual_connection_resources.test_update_resources",
						"resources",
						func(value string) error {
							if len(value) == 0 {
								return fmt.Errorf("Resources field does not exist")
							}
							return nil
						},
					),
				),
			},
		},
	})
}

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

const testResourceClumioAwsManualConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

data "clumio_aws_manual_connection_resources" "test_get_resources" {
	account_native_id = "%s"
	aws_region = "%s"
	asset_types_enabled = {
		ebs = %t
		rds = %t
		ddb = %t
		s3 = %t
		mssql = %t
	}
}
`
