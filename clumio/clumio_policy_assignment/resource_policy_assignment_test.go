// Copyright 2022. Clumio, Inc.

// Acceptance test for clumio_policy_assignment resource.

package clumio_policy_assignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const clumioPolicyIdEnv = "CLUMIO_POLICY_ID"

func testAccPreCheckClumio(t *testing.T) {
	clumio.UtilTestAccPreCheckClumio(t)
	clumio.UtilTestFailIfEmpty(t, clumioPolicyIdEnv, clumioPolicyIdEnv+" cannot be empty.")
}

func TestAccResourceClumioPolicyAssignment(t *testing.T) {
	policyId := os.Getenv("CLUMIO_POLICY_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioPolicyAssignment(policyId),
			},
		},
	})
}

func getTestAccResourceClumioPolicyAssignment(policyId string) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	return fmt.Sprintf(testAccResourceClumioPolicyAssignment, baseUrl, policyId)
}

const testAccResourceClumioPolicyAssignment = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  name = "test_pg_1"
  description = "test-description"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}

resource "clumio_policy_assignment" "test_policy_assignment" {
  entity_id = clumio_protection_group.test_pg.id
  entity_type = "protection_group"
  policy_id = %s
}
`
