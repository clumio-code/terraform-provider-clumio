// Copyright 2023. Clumio, Inc.
//
// This files holds acceptance tests for the clumio_wallet Terraform resource. Please view the
// README.md file for more information on how to run these tests.

//go:build basic

package clumio_wallet_test

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	sdkWallets "github.com/clumio-code/clumio-go-sdk/controllers/wallets"
	"github.com/clumio-code/clumio-go-sdk/models"
	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Basic test of the clumio_wallet resource.
//   - Creates a wallet and verifies that the plan was applied properly.
//   - Ensures that updates to the account ID requires that the resource is re-created as opposed
//     to just updated.
func TestAccResourceWallet(t *testing.T) {
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	accountNativeId2 := os.Getenv(common.ClumioTestAwsAccountId2)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceWallet(baseUrl, accountNativeId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_wallet.test_wallet", plancheck.ResourceActionCreate),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_wallet.test_wallet", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_wallet.test_wallet", "account_native_id",
						regexp.MustCompile(accountNativeId)),
				),
			},
			{
				Config: getTestAccResourceWallet(baseUrl, accountNativeId2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"clumio_wallet.test_wallet", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// Tests that an external deletion of a clumio_wallet resource leads to the resource needing
// to be re-created during the next plan. NOTE the Check function below as it is utilized to delete
// the resource using the Clumio API after the plan is applied.
func TestAccResourceClumioWalletRecreate(t *testing.T) {
	// Retrieve the environment variables required for the test.
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceWallet(baseUrl, accountNativeId),
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
							"clumio_wallet.test_wallet", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_wallet.test_wallet", "account_native_id",
						regexp.MustCompile(accountNativeId)),
					// Delete the resource using the Clumio API after the plan is applied.
					deleteWallet("clumio_wallet.test_wallet"),
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

// Test imports a wallet by ID and ensures that the import is successful.
func TestAccResourceClumioWalletImport(t *testing.T) {
	// Retrieve the environment variables required for the test.
	accountNativeId := os.Getenv(common.ClumioTestAwsAccountId)
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)

	// Create the wallet to import using the Clumio API.
	clumio_pf.UtilTestAccPreCheckClumio(t)
	id, err := createWalletUsingSDK(accountNativeId)
	if err != nil {
		t.Errorf("Error creating wallet using API: %v", err.Error())
	}

	// Run the acceptance test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { clumio_pf.UtilTestAccPreCheckClumio(t) },
		ProtoV6ProviderFactories: clumio_pf.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        getTestAccResourceWallet(baseUrl, accountNativeId),
				ImportState:   true,
				ResourceName:  "clumio_wallet.test_wallet",
				ImportStateId: id,
				ImportStateCheck: func(instStates []*terraform.InstanceState) error {
					if len(instStates) != 1 {
						return errors.New("expected 1 InstanceState for the imported wallet")
					}
					if instStates[0].ID != id {
						errMsg := fmt.Sprintf(
							"Imported wallet has different ID. Expected: %v, Actual: %v",
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

// createWalletUsingSDK creates a wallet using the Clumio API.
func createWalletUsingSDK(accountID string) (string, error) {
	clumioApiToken := os.Getenv(common.ClumioApiToken)
	clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
	clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
	config := clumioConfig.Config{
		Token:                     clumioApiToken,
		BaseUrl:                   clumioApiBaseUrl,
		OrganizationalUnitContext: clumioOrganizationalUnitContext,
		CustomHeaders: map[string]string{
			"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
		},
	}
	wallets := sdkWallets.NewWalletsV1(config)
	res, apiErr := wallets.CreateWallet(&models.CreateWalletV1Request{
		AccountNativeId: &accountID,
	})
	if apiErr != nil {
		return "", apiErr
	}
	if res == nil {
		return "", fmt.Errorf(common.NilErrorMessageDetail)
	}
	return *res.Id, nil
}

// deleteWallet returns a function that deletes a wallet using the Clumio API with
// information from the Terraform state. It is used to intentionally cause a difference between the
// Terraform state and the actual state of the resource in the backend.
func deleteWallet(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		// Create a Clumio API client and delete the wallet.
		clumioApiToken := os.Getenv(common.ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(common.ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(common.ClumioOrganizationalUnitContext)
		config := clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		wallets := sdkWallets.NewWalletsV1(config)
		_, apiErr := wallets.DeleteWallet(rs.Primary.ID)
		if apiErr != nil {
			return apiErr
		}
		return nil
	}
}

// testAccResourcePostWallet returns the Terraform configuration for a basic clumio_wallet resource.
func getTestAccResourceWallet(baseUrl string, accountId string) string {
	return fmt.Sprintf(testAccResourcePostWallet, baseUrl, accountId)
}

// testAccResourcePostWallet is the Terraform configuration for a basic clumio_wallet resource.
const testAccResourcePostWallet = `
provider clumio{
  clumio_api_base_url = "%s"
}

resource "clumio_wallet" "test_wallet" {
  account_native_id = "%s"
}
`
