// Copyright 2021. Clumio, Inc.

// Acceptance tests for provider.
package clumio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)
const (
	awsAccessKeyId = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	awsRegion = "AWS_REGION"
	clumioTestAwsAccountId = "CLUMIO_TEST_AWS_ACCOUNT_ID"
)
// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"clumio": func() (*schema.Provider, error) {
		return New(!isAcceptanceTest())(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New(true)().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}


// testAccPreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call AWS APIs are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func testAccPreCheckAws(t *testing.T) {
	testFailIfAllEmpty(t,
		[]string{awsAccessKeyId, awsProfile, awsSharedCredsFile},
		"One of AWS_ACCESS_KEY_ID, AWS_PROFILE or AWS_SHARED_CREDENTIALS_FILE must be set")
	if os.Getenv(awsAccessKeyId) != ""{
		testFailIfEmpty(t, awsSecretAccessKey, awsSecretAccessKey+" cannot be empty.")
	}
}

// testAccPreCheckClumio validates that the required environment variables are set before
// the acceptance test is executed.
func testAccPreCheckClumio(t *testing.T) {
	testFailIfEmpty(t, clumioTestAwsAccountId, clumioTestAwsAccountId+" cannot be empty")
	testFailIfEmpty(t, clumioApiToken, clumioApiToken+" cannot be empty.")
	testFailIfEmpty(t, clumioApiBaseUrl, clumioApiBaseUrl+" cannot be empty.")
	testFailIfEmpty(t, awsRegion, awsRegion+" cannot be empty")
}

// TestFailIfAllEmpty verifies that at least one environment variable is non-empty or fails the test.
//
// If at lease one environment variable is non-empty, returns the first name and value.
func testFailIfAllEmpty(t *testing.T, names []string, usageMessage string) (string, string) {
	t.Helper()

	name, value, err := RequireOneOf(names, usageMessage)
	if err != nil {
		t.Fatal(err)
		return "", ""
	}

	return name, value
}

// testFailIfEmpty verifies that an environment variable is non-empty or fails the test.
//
// For acceptance tests, this function must be used outside PreCheck functions to set values for configurations.
func testFailIfEmpty(t *testing.T, name string, usageMessage string) string {
	t.Helper()

	value := os.Getenv(name)

	if value == "" {
		t.Fatalf("environment variable %s must be set. Usage: %s", name, usageMessage)
	}

	return value
}
