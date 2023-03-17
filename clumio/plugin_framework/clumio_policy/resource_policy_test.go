// Copyright 2023. Clumio, Inc.

// Acceptance test for clumio_policy resource.
package clumio_policy_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyFixedStart(false),
			},
			{
				Config: getTestAccResourceClumioPolicyFixedStart(true),
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

func getTestAccResourceClumioPolicyFixedStart(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-policy-1234"
	timezone := "UTC"
	window := `
	backup_window_tz {
		start_time = "01:00"
	}`
	if update {
		name = "acceptance-test-policy-4321"
		timezone = "US/Pacific"
		window = `
		backup_window_tz {
			start_time = "05:00"
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
		type = "mssql_database_backup"
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
			mssql_database_backup {
				alternative_replica = "sync_secondary"
				preferred_replica = "primary"
			}
		}
	}
	operations {
		action_setting = "immediate"
		type = "mssql_log_backup"
		%s
		advanced_settings {
			mssql_log_backup {
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
