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
				Config: getTestAccResourceClumioPolicy(false),
			},
			{
				Config: getTestAccResourceClumioPolicy(true),
			},
		},
	})
}

func getTestAccResourceClumioPolicy(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-policy344542"
	if update {
		name = "acceptance-test-policy-2344542"
	}
	return fmt.Sprintf(testAccResourceClumioPolicy, baseUrl, name)
}

const testAccResourceClumioPolicy = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_policy" "test_policy" {
  name = "%s"
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
