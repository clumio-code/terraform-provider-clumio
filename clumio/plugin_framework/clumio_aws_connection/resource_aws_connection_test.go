// Copyright 2023. Clumio, Inc.

// This files holds acceptance tests for the clumio_aws_connection Terraform resource. Please view
// the README.md file for more information on how to run these tests.

//go:build basic

package clumio_aws_connection_test

import (
	"errors"
	"fmt"
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

// Basic test of the clumio_aws_connection resource. It tests the following scenarios:
//   - Creates a connection and verifies that the plan was applied properly.
//   - Updates the connection and verifies that the resource will be updated.
//   - Ensures that updates to the account ID requires that the resource is re-created as opposed to
//     just updated.
func TestAccResourceClumioAwsConnection(t *testing.T) {
	// Retrieve the environment variables required for the test.
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)
	accountNativeId2 := os.Getenv(common.ClumioTestAwsAccountId2)
	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_aws_connection.test_conn", plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_aws_connection.test_conn", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "account_native_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
			{
				Config: getTestAccResourceClumioAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description_updated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_aws_connection.test_conn", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "account_native_id",
						regexp.MustCompile(accountNativeId)),
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "description",
						regexp.MustCompile("test_description_updated")),
				),
			},
			{
				Config: getTestAccResourceClumioAwsConnection(
					baseUrl, accountNativeId2, testAwsRegion, "test_description_updated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_aws_connection.test_conn", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// Tests creation of a AWS connection without setting the description schema attribute in the config.
// This test is ensures that after creating the resource, when we refresh the state it does not
// generate a non-empty plan.
func TestAccResourceClumioAWSConnectionNoDescription(t *testing.T) {

	// Retrieve the environment variables required for the test.
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceClumioAwsConnectionNoDesc,
					baseUrl, accountNativeId, testAwsRegion),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_aws_connection.test_conn", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// Tests that an external deletion of a clumio_aws_connection resource leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioAwsConnectionRecreate(t *testing.T) {
	// Retrieve the environment variables required for the test.
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	testAwsRegion := os.Getenv(common.AwsRegion)

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
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
							"clumio_aws_connection.test_conn", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_aws_connection.test_conn", "account_native_id",
						regexp.MustCompile(accountNativeId)),
					// Delete the resource using the Clumio API after the plan is applied.
					deleteAWSConnection("clumio_aws_connection.test_conn"),
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

// Test imports an AWS connection by ID and ensures that the import is successful.
func TestAccResourceClumioAwsConnectionImport(t *testing.T) {

	// Return if it is not an acceptance test
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			resource.EnvTfAcc))
		return
	}

	// Retrieve the environment variables required for the test.
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	testAwsRegion := os.Getenv(common.AwsRegion)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	// Create the connection to import using the Clumio API.
	clumiopf.UtilTestAccPreCheckClumio(t)
	id, err := createAWSConnectionUsingSDK(
		accountNativeId, testAwsRegion, "test_description")
	if err != nil {
		t.Errorf("Error creating AWS Connection using API: %v", err.Error())
	}

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumiopf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumiopf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioAwsConnection(
					baseUrl, accountNativeId, testAwsRegion, "test_description"),
				ImportState:   true,
				ResourceName:  "clumio_aws_connection.test_conn",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported connection")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported connection has different ID. Expected: %v, Actual: %v",
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

// createAWSConnectionUsingSDK creates an AWS connection using the Clumio API.
func createAWSConnectionUsingSDK(accountID, region, description string) (string, error) {
	clumioApiToken := os.Getenv(common.ClumioApiToken)
	clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
	clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
	config := sdkconfig.Config{
		Token:                     clumioApiToken,
		BaseUrl:                   clumioApiBaseUrl,
		OrganizationalUnitContext: clumioOrganizationalUnitContext,
		CustomHeaders: map[string]string{
			"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
		},
	}
	awsConnection := sdkclients.NewAWSConnectionClient(config)
	res, apiErr := awsConnection.CreateAwsConnection(&models.CreateAwsConnectionV1Request{
		AccountNativeId: &accountID,
		AwsRegion:       &region,
		Description:     &description,
	})
	if apiErr != nil {
		return "", apiErr
	}
	return *res.Id, nil
}

// deleteAWSConnection returns a function that deletes an AWS connection using the Clumio API with
// information from the Terraform state. It is used to intentionally cause a difference between the
// Terraform state and the actual state of the resource in the backend.
func deleteAWSConnection(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		// Create a Clumio API client and delete the AWS connection.
		clumioApiToken := os.Getenv(common.ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
		config := sdkconfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		awsConnection := sdkclients.NewAWSConnectionClient(config)
		_, apiErr := awsConnection.DeleteAwsConnection(rs.Primary.ID)
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

// getTestAccResourceClumioAwsConnection returns the Terraform configuration for a basic
// clumio_aws_connection resource.
func getTestAccResourceClumioAwsConnection(
	baseUrl string, accountId string, awsRegion string, description string) string {
	return fmt.Sprintf(testAccResourceClumioAwsConnection, baseUrl, accountId,
		awsRegion, description)
}

// testAccResourceClumioAwsConnection is the Terraform configuration for a basic
// clumio_aws_connection resource.
const testAccResourceClumioAwsConnection = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
  account_native_id = "%s"
  aws_region = "%s"
  description = "%s"
}
`

// testAccResourceClumioAwsConnectionNoDesc is the Terraform configuration for a
// clumio_aws_connection resource with description attribute not set.
const testAccResourceClumioAwsConnectionNoDesc = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_aws_connection" "test_conn" {
  account_native_id = "%s"
  aws_region = "%s"
}
`
