// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK OrganizationalUnitsV2Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	organizationalunits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
)

type OrganizationalUnitClient interface {
	organizationalunits.OrganizationalUnitsV2Client
}

func NewOrganizationalUnitClient(config config.Config) OrganizationalUnitClient {
	return organizationalunits.NewOrganizationalUnitsV2(config)
}
