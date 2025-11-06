// Copyright 2023. Clumio, Inc.

package clumio_pf

import "fmt"

const (
	// Provider version header key and value used in Clumio client API requests. If the version is
	// being changed here, it must also be changed in the GNUmakefile.
	clumioTfProviderVersionHeader      = "CLUMIO_TERRAFORM_PROVIDER_VERSION"
	clumioTfProviderVersionHeaderValue = "0.16.1"

	// User-Agent header key used in Clumio client API requests. The value to be set is defined
	// below in: userAgentHeaderValue.
	userAgentHeader = "User-Agent"
)

// userAgentHeaderValue is the value to be set for the User-Agent header in Clumio client API
// requests.
var userAgentHeaderValue = fmt.Sprintf("Clumio-Terraform-Provider-%s", clumioTfProviderVersionHeaderValue)
