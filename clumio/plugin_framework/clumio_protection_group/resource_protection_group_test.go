// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_protection_group Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_protection_group_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_protection_group resource. It tests the following scenarios:
//   - Creates a protection group and verifies that the plan was applied properly.
//   - Updates the protection group and verifies that the resource will be updated.
//   - Updates the protection group and verifies that the attributes of resource will be deleted.
//   - Updates the protection group with Object filter and verifies that the resource will be updated.
func TestAccResourceClumioProtectionGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroup(true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1"))),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_protection_group.test_pg", plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: getTestAccResourceClumioProtectionGroup(false, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"clumio_protection_group.test_pg", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"clumio_protection_group.test_pg", "description")),
			},
			{
				Config: getTestAccResourceClumioProtectionGroup(false, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"clumio_protection_group.test_pg", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchTypeSetElemNestedAttrs(
						"clumio_protection_group.test_pg", "object_filter.0.prefix_filters.*",
						map[string]*regexp.Regexp{"prefix": regexp.MustCompile("prefix")})),
			},
		},
	})
}

// Tests creation of a protection group without specifying optional schema attributes in the config
// such as description and bucket_rule. This test ensures that after creating the resource, when we
// refresh the state it generates an empty plan.
func TestAccResourceClumioProtectionGroupNoOptional(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceClumioProtectionGroupNoOptional,
					os.Getenv(common.ClumioApiBaseUrl)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "name",
						regexp.MustCompile("test_pg_1"))),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_protection_group.test_pg", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// Tests that if two Protection Groups are created with the same name, the plan should fail and
// throw an error instead.
func TestAccResourceClumioProtectionGroupDuplicateNameError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceClumioProtectionGroupDuplicateName,
					os.Getenv(common.ClumioApiBaseUrl)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1"))),
				ExpectError: regexp.MustCompile("Unable to create"),
			},
		},
	})
}

// Tests that an external deletion of a clumio_protection_group resource leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioProtectionGroupRecreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroup(true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1")),
					common.DeleteProtectionGroup("clumio_protection_group.test_pg", true),
				),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_protection_group.test_pg", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

// Tests importing a Protection Group by ID and ensuring that the import is successful.
func TestAccResourceClumioAwsProtectionGroupImport(t *testing.T) {

	// Return if it is not an acceptance test
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			resource.EnvTfAcc))
		return
	}

	clumiopf.UtilTestAccPreCheckClumio(t)
	id, err := createProtectionGroupUsingSDK()
	if err != nil {
		t.Errorf("Error creating Protection Group using API: %v", err.Error())
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        getTestAccResourceClumioProtectionGroup(false, false),
				ImportState:   true,
				ResourceName:  "clumio_protection_group.test_pg",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported protection group")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported protection group has different ID. Expected: %v, Actual: %v",
							id, instStates[0].ID)
						return errors.New(errMsg)
					}
					return nil
				},
				ImportStatePersist: true,
				Destroy:            true,
			},
		},
	})
}

// getTestAccResourceClumioProtectionGroup returns the Terraform configuration for a basic
// clumio_protection_group resource.
func getTestAccResourceClumioProtectionGroup(description bool, prefixFilter bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	desc := ""
	if description {
		desc = "description = \"test_pg_1\""
	}
	pf := ""
	if prefixFilter {
		pf = "prefix_filters { prefix = \"prefix\" }"
	}
	val := fmt.Sprintf(testAccResourceClumioProtectionGroup, baseUrl, desc, pf)
	return val
}

// createProtectionGroupUsingSDK creates a protection group using Clumio SDK for testing purpose
func createProtectionGroupUsingSDK() (string, error) {

	name := "test_pg_1"
	bucket_rule := "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
	description := "test_description"
	clumioApiToken := os.Getenv(common.ClumioApiToken)
	clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
	clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
	client := &common.ApiClient{
		ClumioConfig: sdkconfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		},
	}
	pg := sdkclients.NewProtectionGroupClient(client.ClumioConfig)
	s3standard := "S3 Standard"
	s3standardia := "S3 Standard-IA"
	res, apiErr := pg.CreateProtectionGroup(models.CreateProtectionGroupV1Request{
		BucketRule:  &bucket_rule,
		Description: &description,
		Name:        &name,
		ObjectFilter: &models.ObjectFilter{
			StorageClasses: []*string{
				&s3standard, &s3standardia,
			},
		},
	})
	if apiErr != nil {
		return "", apiErr
	}

	// Poll till the Protection Group is available for reading.
	_, err := common.PollForProtectionGroup(
		context.Background(), *res.Id, pg, 300*time.Second, 5*time.Second)
	if err != nil {
		return "", err
	}

	return *res.Id, nil
}

// testAccResourceClumioProtectionGroup is the Terraform configuration for a basic
// clumio_protection_group resource.
const testAccResourceClumioProtectionGroup = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  %s
  object_filter {
	%s
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}
`

// testAccResourceClumioProtectionGroupDuplicateName is the Terraform configuration for a
// clumio_protection_group resource with the protection group having a duplicate name.
const testAccResourceClumioProtectionGroupDuplicateName = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "test_pg_1"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}

resource "clumio_protection_group" "test_pg2"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "test_pg_1"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}
`

// testAccResourceClumioProtectionGroupNoOptional is the Terraform configuration for a
// clumio_protection_group resource with optional fields such as description and bucket_rule not set.
const testAccResourceClumioProtectionGroupNoOptional = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  name = "test_pg_1"
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}
`
