// Copyright 2021. Clumio, Inc.

// Acceptance test for clumio_policy resource.
package clumio

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

)

func TestAccResourceClumioPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckClumio(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicy(false),
			},
			{
				Config: getTestAccResourceClumioPolicy(true),
			},
		},
	})
}

func getTestAccResourceClumioPolicy(update bool) string {
	baseUrl := os.Getenv(clumioApiBaseUrl)
	act_status := "activated"
	if update{
		act_status = "deactivated"
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, act_status)
}

const testAccResourceClumioPolicy = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
  name = "acceptance-test-policy"
  activation_status = "%s"
  operations {
	action_setting = "window"
	type = "aws_ebs_volume_backup"
	backup_window {
		start_time = "08:00"
		end_time = "20:00"
	}
	slas {
		retention_duration {
			unit = "days"
			value = 1
		}
		rpo_frequency {
			unit = "days"
			value = 1
		}
	}
  }
}
`
