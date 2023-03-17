// Copyright 2023. Clumio, Inc.

package clumio_pf

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
)

// ProviderFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"clumio": providerserver.NewProtocol6WithError(New()),
}

// UtilTestAccPreCheckClumio validates that the required environment variables are set before
// the acceptance test is executed.
func UtilTestAccPreCheckClumio(t *testing.T) {
	UtilTestFailIfEmpty(t, common.ClumioTestAwsAccountId, common.ClumioTestAwsAccountId+" cannot be empty")
	UtilTestFailIfEmpty(t, common.ClumioApiToken, common.ClumioApiToken+" cannot be empty.")
	UtilTestFailIfEmpty(t, common.ClumioApiBaseUrl, common.ClumioApiBaseUrl+" cannot be empty.")
	UtilTestFailIfEmpty(t, common.AwsRegion, common.AwsRegion+" cannot be empty")
}

// UtilTestFailIfAllEmpty verifies that at least one environment variable is non-empty or fails the test.
// If at lease one environment variable is non-empty, returns the first name and value.
func UtilTestFailIfAllEmpty(t *testing.T, names []string, usageMessage string) (string, string) {
	t.Helper()

	name, value, err := common.RequireOneOf(names, usageMessage)
	if err != nil {
		t.Fatal(err)
		return "", ""
	}

	return name, value
}

// UtilTestFailIfEmpty verifies that an environment variable is non-empty or fails the test.
//
// For acceptance tests, this function must be used outside PreCheck functions to set values for configurations.
func UtilTestFailIfEmpty(t *testing.T, name string, usageMessage string) string {
	t.Helper()

	value := os.Getenv(name)

	if value == "" {
		t.Fatalf("environment variable %s must be set. Usage: %s", name, usageMessage)
	}

	return value
}
