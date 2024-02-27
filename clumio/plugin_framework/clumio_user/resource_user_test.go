// Copyright 2023. Clumio, Inc.
//
// This files holds acceptance tests for the clumio_user Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_user_test

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/controllers/users"
	"github.com/clumio-code/clumio-go-sdk/models"
	clumioPf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	assignedRoleBefore = "30000000-0000-0000-0000-000000000000"
	assignedRoleAfter  = "20000000-0000-0000-0000-000000000000"
)

// Basic test of the clumio_user resource. It tests the following scenarios:
//   - Creates a user and verifies that the plan was applied properly.
//   - Updates the user and verifies that the resource will be updated.
//   - Ensures that updates to the email requires that the resource is re-created as opposed to
//     just updated.
func TestAccResourceClumioUser(t *testing.T) {

	// Retrieve the environment variables required for the test.
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	email := "test@clumio.com"
	name := "acceptance-test-user"

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioUser(baseUrl, name, email, false),
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
							"clumio_user.test_user", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_user.test_user", "assigned_role",
						regexp.MustCompile(assignedRoleBefore)),
				),
			},
			{
				Config: getTestAccResourceClumioUser(baseUrl, name, email, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_user.test_user", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_user.test_user", "assigned_role",
						regexp.MustCompile(assignedRoleAfter)),
				),
			},
			{
				Config: getTestAccResourceClumioUser(baseUrl, name, "test1@clumio.com", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"clumio_user.test_user", plancheck.ResourceActionReplace),
					},
				},
			},
			{
				Config: getTestAccResourceClumioUser(baseUrl, "updated-name", email, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"clumio_user.test_user", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// Tests that an external deletion of a clumio_user resource leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioUserRecreate(t *testing.T) {

	// Retrieve the environment variable required for the test.
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioUser(
					baseUrl, "acceptance-test-user", "test@clumio.com", false),
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
							"clumio_user.test_user", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					deleteUser("clumio_user.test_user"),
				),
				// This attribute is used to denote that the test expects that after the plan is
				// applied and a refresh is run, a non-empty plan is expected due to differences
				// from the state. Without this attribute set, the test would fail as it is unaware
				// that the resource was deleted externally.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Test imports a user by ID and ensures that the import is successful.
func TestAccResourceClumioPolicyImport(t *testing.T) {

	// Create the user to import using the Clumio API.
	clumioPf.UtilTestAccPreCheckClumio(t)
	id, err := createUserUsingSDK()
	if err != nil {
		t.Errorf("Error creating Policy using API: %v", err.Error())
	}

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumioPf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumioPf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceClumioUserImport, os.Getenv(common.ClumioApiBaseUrl)),
				ImportState:   true,
				ResourceName:  "clumio_user.test_user",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported user")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported policy has different ID. Expected: %v, Actual: %v",
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

// deleteUser returns a function that deletes an user using the Clumio API with
// information from the Terraform state. It is used to intentionally cause a difference between the
// Terraform state and the actual state of the resource in the backend.
func deleteUser(resourceName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		// Retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}

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
		userAPI := users.NewUsersV2(client.ClumioConfig)
		userId, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		_, apiErr := userAPI.DeleteUser(userId)
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

// createUserUsingSDK creates a Clumio User using the Clumio API
func createUserUsingSDK() (string, error) {

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
	userAPI := users.NewUsersV2(client.ClumioConfig)
	name := "acceptance-test-user"
	email := "test@clumio.com"

	res, apiErr := userAPI.CreateUser(&models.CreateUserV2Request{
		FullName: &name,
		Email:    &email,
	})
	if apiErr != nil {
		return "", apiErr
	}
	return *res.Id, nil
}

// getTestAccResourceClumioUser returns the Terraform configuration for a basic
// clumio_user resource.
func getTestAccResourceClumioUser(baseUrl string, name string, email string, update bool) string {
	orgUnitId := "clumio_organizational_unit.test_ou1.id"
	assignedRole := assignedRoleBefore
	if update {
		orgUnitId = "clumio_organizational_unit.test_ou2.id"
		assignedRole = assignedRoleAfter
	}
	return fmt.Sprintf(testAccResourceClumioUser, baseUrl, name, email, assignedRole, orgUnitId)
}

// testAccResourceClumioUser is the Terraform configuration for a basic
// clumio_user resource.
const testAccResourceClumioUser = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou1" {
  name = "test_ou1"
  description = "test-ou-1"
}

resource "clumio_organizational_unit" "test_ou2" {
  name = "test_ou2"
  description = "test-ou-2"
}

resource "clumio_user" "test_user" {
  full_name = "%s"
  email = "%s"
  access_control_configuration = [
	{
		role_id = "%s"
		organizational_unit_ids = [%s]
	},
  ]
}
`

// testAccResourceClumioUser is the Terraform configuration for a importing a
// clumio_user resource.
const testAccResourceClumioUserImport = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_user" "test_user" {
  full_name = "acceptance-test-user"
  email = "test@clumio.com"
}
`
