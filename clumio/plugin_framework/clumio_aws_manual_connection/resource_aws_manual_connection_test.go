package clumio_aws_manual_connection_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/clumio_aws_manual_connection"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestManualResources(t *testing.T) {
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
	testResources := &clumio_aws_manual_connection.ResourcesModel{
		ClumioIAMRoleArn: basetypes.NewStringValue(os.Getenv(common.ClumioIAMRoleArn)),
		ClumioSupportRoleArn: basetypes.NewStringValue(os.Getenv(common.ClumioSupportRoleArn)),
		ClumioEventPubArn: basetypes.NewStringValue(os.Getenv(common.ClumioEventPubArn)),
		EventRules: &clumio_aws_manual_connection.EventRules{
			CloudtrailRuleArn: basetypes.NewStringValue(os.Getenv(common.CloudtrailRuleArn)),
			CloudwatchRuleArn: basetypes.NewStringValue(os.Getenv(common.CloudwatchRuleArn)),
		},
		ServiceRoles: &clumio_aws_manual_connection.ServiceRoles{
			S3: &clumio_aws_manual_connection.S3ServiceRoles{
				ContinuousBackupsRoleArn: basetypes.NewStringValue(os.Getenv(common.ContinuousBackupsRoleArn)),
			},
			Mssql: &clumio_aws_manual_connection.MssqlServiceRoles{
				Ec2SsmInstanceProfileArn: basetypes.NewStringValue(os.Getenv(common.Ec2SsmInstanceProfileArn)),
				SsmNotificationRoleArn: basetypes.NewStringValue(os.Getenv(common.SsmNotificationRoleArn)),
			},
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { 
			clumio_pf.UtilTestAccPreCheckClumio(t) 
			clumio_pf.UtilTestAwsManualConnectionPreCheckClumio(t)
		},
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestClumioAwsManualConnection(
					baseUrl, testAccountId, testAwsRegion, testAssetTypes, testResources),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_manual_connection.test_update_resources", "account_id",
						regexp.MustCompile(testAccountId)),
						resource.TestMatchResourceAttr(
							"clumio_aws_manual_connection.test_update_resources", "aws_region",
							regexp.MustCompile(testAwsRegion)),
							resource.TestMatchResourceAttr(
								"clumio_aws_manual_connection.test_update_resources", "connection_status",
								regexp.MustCompile("connected")),
				),
			},
		},
	})
}

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
			testResources.ClumioIAMRoleArn,
			testResources.ClumioEventPubArn,
			testResources.ClumioSupportRoleArn,
			testResources.EventRules.CloudtrailRuleArn,
			testResources.EventRules.CloudwatchRuleArn,
			testResources.ServiceRoles.S3.ContinuousBackupsRoleArn,
			testResources.ServiceRoles.Mssql.SsmNotificationRoleArn,
			testResources.ServiceRoles.Mssql.Ec2SsmInstanceProfileArn,
		)
}

const testResourceClumioAwsManualConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_aws_manual_connection" "test_update_resources" {
	
    account_id = "%s"
    aws_region = "%s"
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
