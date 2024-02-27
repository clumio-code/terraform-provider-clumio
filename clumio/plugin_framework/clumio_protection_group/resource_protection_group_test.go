// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_protection_group Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_protection_group_test

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/clumio-code/clumio-go-sdk/models"
	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_protection_group resource. It tests the following scenarios:
//   - Creates a protection group and verifies that the plan was applied properly.
//   - Updates the protection group and verifies that the resource will be updated.
func TestAccResourceClumioProtectionGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroup(false),
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
				Config: getTestAccResourceClumioProtectionGroup(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"clumio_protection_group.test_pg", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1_updated"))),
			},
		},
	})
}

// Tests creation of a protection group without specifying optional schema attributes in the config
// such as description and bucket_rule. This test ensures that after creating the resource, when we
// refresh the state it generates an empty plan.
func TestAccResourceClumioProtectionGroupNoOptional(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
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
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceClumioProtectionGroupDuplicateName,
					os.Getenv(common.ClumioApiBaseUrl)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1"))),
				ExpectError: regexp.MustCompile("Error creating Protection Group test_pg_1"),
			},
		},
	})
}

// Tests that an external deletion of a clumio_protection_group resource leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioProtectionGroupRecreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroup(false),
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

// Tests creation of a organizational unit and using that organizational unit to create a
// protection group
func TestAccResourceClumioProtectionGroupWithOU(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioProtectionGroupWithOU(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1"))),
			},
			{
				Config: getTestAccResourceClumioProtectionGroupWithOU(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_protection_group.test_pg", "description",
						regexp.MustCompile("test_pg_1_updated"))),
			},
		},
	})
}

// Tests that empty organizational_unit returns error
func TestAccResourceClumioProtectionGroupWithEmptyOU(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceClumioProtectionGroupWithEmptyOU,
					os.Getenv(common.ClumioApiBaseUrl)),
				ExpectError: regexp.MustCompile(
					"Attribute organizational_unit_id string length must be at least 1"),
			},
		},
	})
}

// Tests importing a Protection Group by ID and ensuring that the import is successful.
func TestAccResourceClumioAwsProtectionGroupImport(t *testing.T) {
	clumioPf.UtilTestAccPreCheckClumio(t)
	id, err := createProtectionGroupUsingSDK()
	if err != nil {
		t.Errorf("Error creating Protection Group using API: %v", err.Error())
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        getTestAccResourceClumioProtectionGroup(false),
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
func getTestAccResourceClumioProtectionGroup(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	description := "test_pg_1"
	if update {
		description = "test_pg_1_updated"
	}
	val := fmt.Sprintf(testAccResourceClumioProtectionGroup, baseUrl, description)
	return val
}

// getTestAccResourceClumioProtectionGroupWithOU returns the Terraform configuration for a
// clumio_protection_group resource with the protection group having a organizational unit.
func getTestAccResourceClumioProtectionGroupWithOU(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	description := "test_pg_1"
	if update {
		description = "test_pg_1_updated"
	}
	val := fmt.Sprintf(testAccResourceClumioProtectionGroupWithOU, baseUrl, description)
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
		ClumioConfig: clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		},
	}
	pg := protectionGroups.NewProtectionGroupsV1(client.ClumioConfig)
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
	return *res.Id, nil
}

// testAccResourceClumioProtectionGroupis the Terraform configuration for a basic
// clumio_protection_group resource.
const testAccResourceClumioProtectionGroup = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "%s"
  object_filter {
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

// testAccResourceClumioProtectionGroupWithOU is the Terraform configuration for a
// clumio_protection_group resource with the protection group having an organizational unit.
const testAccResourceClumioProtectionGroupWithOU = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou2" {
  name = "test_ou2"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "%s"
  organizational_unit_id = clumio_organizational_unit.test_ou2.id
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

// testAccResourceClumioProtectionGroupWithEmptyOU is the Terraform configuration for a
// clumio_protection_group resource with the protection group having an empty organizational unit.
const testAccResourceClumioProtectionGroupWithEmptyOU = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_protection_group" "test_pg"{
  bucket_rule = "{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}"
  name = "test_pg_1"
  description = "some description"
  organizational_unit_id = ""
  object_filter {
	storage_classes = ["S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA", "S3 Reduced Redundancy"]
  }
}
`
