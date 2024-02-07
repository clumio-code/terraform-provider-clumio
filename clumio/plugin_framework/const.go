// Copyright 2023. Clumio, Inc.

package clumio_pf

import "fmt"

const (
	// Provider version key and value
	clumioTfProviderVersionKey = "CLUMIO_TERRAFORM_PROVIDER_VERSION"
	// If the version is being changed here, it must also be changed in the GNUmakefile.
	clumioTfProviderVersionValue = "0.5.9"
	// User-Agent Header
	userAgentHeader = "User-Agent"

	errorFmt = "The provider cannot create the Clumio API client as" +
		" there is an unknown configuration value for the Clumio API %s. " +
		"Either target apply the source of the value first, set the value" +
		" statically in the configuration, or use the %s environment variable."
	baseUrl = "Base URL"
	token   = "Token"
)

var userAgentHeaderValue = fmt.Sprintf("Clumio-Terraform-Provider-%s", clumioTfProviderVersionValue)
