// Copyright 2023. Clumio, Inc.

// Acceptance test for clumio_policy resource.
package clumio_policy_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyFixedStart(0),
			},
			{
				Config: getTestAccResourceClumioPolicyFixedStart(1),
			},
			{
				Config: getTestAccResourceClumioPolicyFixedStart(2),
			},
			{
				Config: getTestAccResourceClumioPolicyWindow(false),
			},
			{
				Config: getTestAccResourceClumioPolicyWindow(true),
			},
			{
				Config: getTestAccResourceClumioPolicySecureVaultLite(false),
			},
			{
				Config: getTestAccResourceClumioPolicySecureVaultLite(true),
			},
			{
				Config: getTestAccResourceClumioPolicyHourlyMinutely(false),
			},
			{
				Config: getTestAccResourceClumioPolicyHourlyMinutely(true),
			},
			{
				Config: getTestAccResourceClumioPolicyWeekly(false),
			},
			{
				Config: getTestAccResourceClumioPolicyWeekly(true),
			},
			{
				Config: getTestAccResourceClumioPolicyBackupRegion(0),
			},
			{
				Config: getTestAccResourceClumioPolicyBackupRegion(1),
			},
			{
				Config:      getTestAccResourceClumioPolicyBackupRegion(2),
				ExpectError: regexp.MustCompile(".*Error running apply.*"),
			},
		},
	})
}

func getTestAccResourceClumioPolicyWindow(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-policy-1234"
	timezone := "UTC"
	window := `
	backup_window_tz {
		start_time = "01:00"
		end_time = "05:00"
	}`
	if update {
		name = "acceptance-test-policy-4321"
		timezone = "US/Pacific"
		window = `
		backup_window_tz {
			start_time = "03:00"
			end_time = "07:00"
		}`
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, name, timezone, window)
}

func getTestAccResourceClumioPolicyFixedStart(update int) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-policy-1234"
	timezone := "UTC"
	window := `
	backup_window_tz {
		start_time = "01:00"
	}`
	if update == 1 {
		name = "acceptance-test-policy-4321"
		timezone = "US/Pacific"
		window = `
		backup_window_tz {
			start_time = "05:00"
		}`
	} else if update == 2 {
		window = `
		backup_window_tz {
			start_time = "05:00"
			end_time = ""
		}`
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, name, timezone, window)
}

func getTestAccResourceClumioPolicySecureVaultLite(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "SecureVaultLite Test"
	sla := ``
	if update {
		sla = `
		slas {
			retention_duration {
				unit = "months"
				value = 3
			}
			rpo_frequency {
				unit = "months"
				value = 1
			}
		}`
	}
	return fmt.Sprintf(testAccResourceClumioPolicyVaultLite, baseUrl, name, sla, sla)
}

func getTestAccResourceClumioPolicyHourlyMinutely(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Hourly & Minutely Policy Create"
	hourlySla := `
	slas {
		retention_duration {
			unit = "days"
			value = 15
		}
		rpo_frequency {
			unit = "hours"
			value = 4
		}
	}
	`
	minutelySla := `
	slas {
		retention_duration {
			unit = "days"
			value = 5
		}
		rpo_frequency {
			unit = "minutes"
			value = 15
		}
	}
	`
	if update {
		name = "Hourly & Minutely Policy Update"
		hourlySla = `
		slas {
			retention_duration {
				unit = "days"
				value = 15
			}
			rpo_frequency {
				unit = "hours"
				value = 12
			}
		}
		`
		minutelySla = `
		slas {
			retention_duration {
				unit = "days"
				value = 5
			}
			rpo_frequency {
				unit = "minutes"
				value = 30
			}
		}
		`
	}
	return fmt.Sprintf(testAccResourceClumioPolicyHourlyMinutely, baseUrl, name, hourlySla, minutelySla)
}

func getTestAccResourceClumioPolicyWeekly(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Weekly Policy Create"
	weeklySla := `
	slas {
		retention_duration {
			unit = "weeks"
			value = 4
		}
		rpo_frequency {
			unit = "weeks"
			value = 1
			offsets = [1]
		}
	}
	`
	if update {
		name = "Weekly Policy Update"
		weeklySla = `
		slas {
			retention_duration {
				unit = "weeks"
				value = 5
			}
			rpo_frequency {
				unit = "weeks"
				value = 1
				offsets = [3]
			}
		}
		`
	}
	return fmt.Sprintf(testAccResourceClumioPolicyWeekly, baseUrl, name, weeklySla)
}

func getTestAccResourceClumioPolicyBackupRegion(scenario int) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Backup Region Policy Create"
	timezone := "UTC"
	region := `
	backup_aws_region = "us-west-2"`
	if scenario == 1 {
		name = "Backup Region Policy Update"
		region = `` // valid as the region is optional
	} else if scenario == 2 {
		name = "Backup Region Policy Update 2"
		region = `
	backup_aws_region = ""` // invalid as empty region is not allowed as request.
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, name, timezone, region)
}

const testAccResourceClumioPolicy = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
	name = "%s"
	timezone = "%s"
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
		%s
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
`

const testAccResourceClumioPolicyVaultLite = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "secure_vault_lite_success" {
  name = "%s"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 30
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
			}
			%s
			advanced_settings {
			aws_ebs_volume_backup {
				backup_tier = "lite"
			}
		}
	}
	operations {
		action_setting = "immediate"
		type = "aws_ec2_instance_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 30
			}
			rpo_frequency {
				unit = "days"
				value = 1
			}
		}
		%s
		advanced_settings {
			aws_ec2_instance_backup {
				backup_tier = "lite"
			}
		}
	}
}
`

const testAccResourceClumioPolicyHourlyMinutely = `
provider clumio{
	clumio_api_base_url = "%s"
}
resource "clumio_policy" "hourly_minutely_policy" {
	name = "%s"
	operations {
		action_setting = "immediate"
		type = "ec2_mssql_database_backup"
		slas {
			retention_duration {
				unit = "days"
				value = 30
			}
			rpo_frequency {
				unit = "days"
				value = 3
			}
		}
		%s
		advanced_settings {
			ec2_mssql_database_backup {
				alternative_replica = "sync_secondary"
				preferred_replica = "primary"
			}
		}
	}
	operations {
		action_setting = "immediate"
		type = "ec2_mssql_log_backup"
		%s
		advanced_settings {
			ec2_mssql_log_backup {
				alternative_replica = "sync_secondary"
				preferred_replica = "primary"
			}
		}
	}
}
`

const testAccResourceClumioPolicyWeekly = `
provider clumio{
	clumio_api_base_url = "%s"
}
resource "clumio_policy" "weekly_policy" {
	name = "%s"
	operations {
		action_setting = "immediate"
		type = "aws_ebs_volume_backup"
		%s
	}
}
`

func TestRdsComplianceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyRDSCompliance(false),
			},
			{
				Config: getTestAccResourceClumioPolicyRDSCompliance(true),
			},
		},
	})
}

func getTestAccResourceClumioPolicyRDSCompliance(update bool) string {

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "Rds Compliance Policy Create"
	slas := `
	slas {
		retention_duration {
			unit  = "days"
			value = 31
		}
		rpo_frequency {
			unit  = "days"
			value = 7
		}
	}
	advanced_settings {
		aws_rds_resource_granular_backup {
			backup_tier = "frozen"
		}
	}
	`
	if update {
		name = "Rds Compliance Policy Update"
		slas = `
		slas {
			retention_duration {
				unit  = "days"
				value = 28
			}
			rpo_frequency {
				unit  = "days"
				value = 7
			}
		}
		advanced_settings {
			aws_rds_resource_granular_backup {
				backup_tier = "frozen"
			}
		}
		`
	}
	return fmt.Sprintf(testAccResourceClumioRdsPolicy, baseUrl, name, slas)
}

func TestRdsPitrClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// PITR only
				Config: getTestClumioPolicyRds(true, false, false),
			},
			{
				// GRR only
				Config: getTestClumioPolicyRds(false, true, false),
			},
			{
				// Airgap only
				Config: getTestClumioPolicyRds(false, false, true),
			},
			{
				// PITR and GRR both
				Config: getTestClumioPolicyRds(true, true, false),
			},
			{
				// RDS PITR immediate
				Config: getTestClumioPolicyRdsPitrAdv(true),
			},
			{
				// RDS PITR maintenance_window
				Config: getTestClumioPolicyRdsPitrAdv(false),
			},
			// Note: this test is disabled as advanced settings is not supported
			// if `apibackend.EnableRDSComplianceTier` is off.
			// Consider this to be enabled after advanced settings are mandated.
			// {
			// 	// RDS Compliance tier
			// 	Config: getTestClumioPolicyRdsComplianceTier(true),
			// },
		},
	})
}

func getTestClumioPolicyRds(pitr bool, logical bool, airgap bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf-rds-policy"
	operations := ""
	// TODO: add advanced settings on it.
	pitrTemplate := `
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_aws_snapshot"
		slas {
			retention_duration {
				unit  = "days"
				value = 7
			}
			rpo_frequency {
				unit  = "days"
				value = 1
			}
		}
	}`
	logicalTemplate := `
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_granular_backup"
		slas {
			retention_duration {
				unit  = "days"
				value = 31
			}
			rpo_frequency {
				unit  = "days"
				value = 7
			}
		}
		advanced_settings {
			aws_rds_resource_granular_backup {
				backup_tier = "standard"
			}
		}
	}`
	airgapTemplate := `
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_rolling_backup"
		slas {
			retention_duration {
				unit  = "days"
				value = 31
			}
			rpo_frequency {
				unit  = "days"
				value = 7
			}
		}
	}`

	if pitr {
		operations += pitrTemplate
		name += "-pitr"
	}
	if logical {
		operations += logicalTemplate
		name += "-logical"
	}
	if airgap {
		operations += airgapTemplate
		name += "-airgap"
	}
	return fmt.Sprintf(testClumioPolicyRdsPolicyTemplate, baseUrl, name, operations)
}

func getTestClumioPolicyRdsPitrAdv(immediate bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf-rds-pitr-adv-policy"
	rdsPitrConfigAdv := "immediate"
	if !immediate {
		rdsPitrConfigAdv = "maintenance_window"
	}
	operations := fmt.Sprintf(`
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_aws_snapshot"
		slas {
			retention_duration {
				unit  = "days"
				value = 7
			}
			rpo_frequency {
				unit  = "days"
				value = 1
			}
		}
		advanced_settings {
			aws_rds_config_sync {
				apply = "%s"
			}
		}
	}`, rdsPitrConfigAdv)
	return fmt.Sprintf(testClumioPolicyRdsPolicyTemplate, baseUrl, name, operations)
}

func getTestClumioPolicyRdsComplianceTier(complianceTier bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "tf-rds-pitr-adv-policy"
	backupTier := "standard"
	if complianceTier {
		backupTier = "frozen"
	}
	operations := fmt.Sprintf(`
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_granular_backup"
		slas {
			retention_duration {
				unit  = "days"
				value = 31
			}
			rpo_frequency {
				unit  = "days"
				value = 7
			}
		}
		advanced_settings {
			aws_rds_resource_granular_backup {
				backup_tier = "%s"
			}
		}
	}`, backupTier)
	return fmt.Sprintf(testClumioPolicyRdsPolicyTemplate, baseUrl, name, operations)
}

const testClumioPolicyRdsPolicyTemplate = `
provider clumio{
	clumio_api_base_url = "%s"
}
resource "clumio_policy" "tf_rds_policy" {
	name = "%s"
	%s
}
`

const testAccResourceClumioRdsPolicy = `
provider clumio{
	clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
	name = "%s"
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_granular_backup"
		%s
	}
	operations {
		action_setting = "immediate"
		type           = "aws_rds_resource_aws_snapshot"
		slas {
			retention_duration {
				unit  = "days"
				value = 7
			}
			rpo_frequency {
				unit  = "days"
				value = 1
			}
		}
	}
}
`
