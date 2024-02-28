// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_aws_manual_connection Terraform resource.
// Please view the README.md file for more information on how to run these tests.

//go:build manual_connection

package clumio_aws_manual_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_manual_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Basic test of the clumio_aws_manual_connection resource. It tests updating an undeployed
// connection with enabled asset types and stack ARN links to manually configured resources and
// ensuring that the connection is deployed and moves to "connected" status.
func TestClumioAwsManualConnection(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAccountId := os.Getenv(common.ClumioTestAwsAccountId)
	testAwsRegion := os.Getenv(common.AwsRegion)
	testAssetTypes := map[string]bool{
		"EBS":      true,
		"S3":       true,
		"RDS":      true,
		"DynamoDB": true,
		"EC2MSSQL": true,
	}
	testResources := &clumio_aws_manual_connection.ResourcesModel{
		ClumioIAMRoleArn:     basetypes.NewStringValue(os.Getenv(common.ClumioIAMRoleArn)),
		ClumioSupportRoleArn: basetypes.NewStringValue(os.Getenv(common.ClumioSupportRoleArn)),
		ClumioEventPubArn:    basetypes.NewStringValue(os.Getenv(common.ClumioEventPubArn)),
		EventRules: &clumio_aws_manual_connection.EventRules{
			CloudtrailRuleArn: basetypes.NewStringValue(os.Getenv(common.CloudtrailRuleArn)),
			CloudwatchRuleArn: basetypes.NewStringValue(os.Getenv(common.CloudwatchRuleArn)),
		},
		ServiceRoles: &clumio_aws_manual_connection.ServiceRoles{
			S3: &clumio_aws_manual_connection.S3ServiceRoles{
				ContinuousBackupsRoleArn: basetypes.NewStringValue(
					os.Getenv(common.ContinuousBackupsRoleArn)),
			},
			Mssql: &clumio_aws_manual_connection.MssqlServiceRoles{
				Ec2SsmInstanceProfileArn: basetypes.NewStringValue(
					os.Getenv(common.Ec2SsmInstanceProfileArn)),
				SsmNotificationRoleArn:   basetypes.NewStringValue(
					os.Getenv(common.SsmNotificationRoleArn)),
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			clumioPf.UtilTestAccPreCheckClumio(t)
		},
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestClumioAwsManualConnection(
					baseUrl, testAccountId, testAwsRegion, testAssetTypes, testResources),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "connection_status",
						regexp.MustCompile("connected")),
				),
			},
		},
	})
}

// getTestClumioAwsManualConnection returns the Terraform configuration for a basic
// clumio_aws_manual_connection resource.
func getTestClumioAwsManualConnection(
	baseUrl string, accountId string, awsRegion string,
	testAssetTypes map[string]bool,
	testResources *clumio_aws_manual_connection.ResourcesModel) string {
		return fmt.Sprintf(testResourceClumioAwsManualConnection, 
			baseUrl,
			accountId,
			awsRegion,
			testAssetTypes["EBS"],
			testAssetTypes["RDS"],
			testAssetTypes["DynamoDB"],
			testAssetTypes["S3"],
			testAssetTypes["EC2MSSQL"],
			testResources.ClumioIAMRoleArn.ValueString(),
			testResources.ClumioEventPubArn.ValueString(),
			testResources.ClumioSupportRoleArn.ValueString(),
			testResources.EventRules.CloudtrailRuleArn.ValueString(),
			testResources.EventRules.CloudwatchRuleArn.ValueString(),
			testResources.ServiceRoles.S3.ContinuousBackupsRoleArn.ValueString(),
			testResources.ServiceRoles.Mssql.SsmNotificationRoleArn.ValueString(),
			testResources.ServiceRoles.Mssql.Ec2SsmInstanceProfileArn.ValueString(),
		)
}

// testResourceClumioAwsManualConnection is the Terraform configuration for a basic setup of a AWS
// Manual Connection using clumio_aws_connection to create the connection and using the
// clumio_aws_manual_connection resource to deploy the manual resources on it.
const testResourceClumioAwsManualConnection = `
provider clumio{
    clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
    account_native_id = "%s"
    aws_region = "%s"
}

resource "clumio_aws_manual_connection" "test_update_resources" {
    account_id = clumio_aws_connection.test_conn.account_native_id
    aws_region = clumio_aws_connection.test_conn.aws_region
    assets_enabled = {
        ebs = %t
        rds = %t
        ddb = %t
        s3 = %t
        mssql = %t
    }
    resources = {
        clumio_iam_role_arn = "%s"
        clumio_event_pub_arn = "%s"
        clumio_support_role_arn = "%s"
        event_rules = {
            cloudtrail_rule_arn = "%s"
            cloudwatch_rule_arn = "%s"
        }

        service_roles = {
            s3 = {
                continuous_backups_role_arn = "%s"
            }
            mssql = {
                ssm_notification_role_arn = "%s"
                ec2_ssm_instance_profile_arn = "%s"
            }
        }
    }
}
`
