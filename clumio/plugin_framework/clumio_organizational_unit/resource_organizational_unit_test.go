// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_organizational_unit Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_organizational_unit_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	clumiopf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	ouNameBefore   = "acceptance-test-ou"
	ouNameAfter    = "acceptance-test-ou-updated"
	descNameBefore = "test-ou-description"
	descNameAfter  = "test-ou-description-updated"
)

// Basic test of the clumio_organizational_unit resource. It tests the following scenarios:
//   - Creates an organizational unit and verifies that the plan was applied properly.
//   - Updates the organizational unit and verifies that the resource will be updated.
func TestAccResourceClumioOrganizationalUnit(t *testing.T) {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioOrganizationalUnit(baseUrl, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "name",
						regexp.MustCompile(ouNameBefore)),
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "description",
						regexp.MustCompile(descNameBefore)),
				),
			},
			{
				Config: getTestAccResourceClumioOrganizationalUnit(baseUrl, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "name",
						regexp.MustCompile(ouNameAfter)),
					resource.TestMatchResourceAttr(
						"clumio_organizational_unit.test_ou", "description",
						regexp.MustCompile(descNameAfter)),
				),
			},
		},
	})
}

// Tests creation of an organizational unit without setting the description schema attribute in the
// config. This test is ensures that after creating the resource, when we refresh the state it does
// not generate a non-empty plan.
func TestAccResourceClumioOrganizationalUnitNoDescription(t *testing.T) {

	// Retrieve the environment variables required for the test.
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(getTestAccResourceClumioOrganizationalUnit(
					baseUrl, false, true)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_organizational_unit.test_ou", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// Test imports an organizational unit by ID and ensures that the import is successful.
func TestAccResourceClumioOrganizationalUnitImport(t *testing.T) {

	// Return if it is not an acceptance test
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			resource.EnvTfAcc))
		return
	}

	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	clumiopf.UtilTestAccPreCheckClumio(t)
	id, err := createOrganizationalUnitUsingSDK()
	if err != nil {
		t.Errorf("Error creating AWS Connection using API: %v", err.Error())
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        getTestAccResourceClumioOrganizationalUnit(baseUrl, false, false),
				ImportState:   true,
				ResourceName:  "clumio_organizational_unit.test_ou",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported OU")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported OU has different ID. Expected: %v, Actual: %v",
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

// Tests to check if creating a OU with empty parent id returns error.
func TestAccResourceClumioOrganizationalUnitEmptyParentID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioOrganizationalUnitEmptyParentId(),
				ExpectError: regexp.MustCompile(
					"Attribute parent_id string length must be at least 1"),
			},
		},
	})
}

// getTestAccResourceClumioOrganizationalUnit returns the Terraform configuration for a basic
// clumio_organizational_unit resource.
func getTestAccResourceClumioOrganizationalUnit(baseUrl string, update bool, nodesc bool) string {
	content :=
		`name = "acceptance-test-ou"
		 description = "test-ou-description-updated"
		`
	if nodesc {
		content =
			`name = "acceptance-test-ou"`
	}
	if update {
		content =
			`name = "acceptance-test-ou-updated"
			 description = "test-ou-description-updated"
			`
	}
	return fmt.Sprintf(testAccResourceClumioOrganizationalUnit, baseUrl, content)
}

// getTestAccResourceClumioOrganizationalUnitEmptyParentId returns the Terraform configuration for a
// clumio_organizational_unit resource with empty string for parent_id.
func getTestAccResourceClumioOrganizationalUnitEmptyParentId() string {
	content :=
		`name = "acceptance-test-ou"
		 description = "test-ou-description-updated"
         parent_id = ""
		`
	return fmt.Sprintf(testAccResourceClumioOrganizationalUnit, os.Getenv(common.ClumioApiBaseUrl),
		content)
}

// createOrganizationalUnitUsingSDK creates an organizational unit using the Clumio SDK
func createOrganizationalUnitUsingSDK() (string, error) {
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
	ou := sdkclients.NewOrganizationalUnitClient(client.ClumioConfig)
	name := "acceptance-test-ou"
	res, apiErr := ou.CreateOrganizationalUnit(nil,
		&models.CreateOrganizationalUnitV2Request{
			Name: &name,
		})
	if apiErr != nil {
		return "", apiErr
	}
	var id string
	if res.StatusCode == http.StatusOK {
		id = *res.Http200.Id
	} else if res.StatusCode == http.StatusAccepted {
		id = *res.Http202.Id
	}
	return id, nil
}

// testAccResourceClumioOrganizationalUnit is the Terraform configuration for a basic
// clumio_organizational_unit resource.
const testAccResourceClumioOrganizationalUnit = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou" {
   %s
}
`
