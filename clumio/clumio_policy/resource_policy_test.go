// Copyright 2021. Clumio, Inc.

// Acceptance test for clumio_policy resource.
package clumio_policy_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyWindow(false),
			},
			{
				Config: getTestAccResourceClumioPolicyWindow(true),
			},
			{
				Config: getTestAccResourceClumioPolicyFixedStart(false),
			},
			{
				Config: getTestAccResourceClumioPolicyFixedStart(true),
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
