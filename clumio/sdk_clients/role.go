// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK RolesV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkRoles "github.com/clumio-code/clumio-go-sdk/controllers/roles"
)

type RoleClient interface {
	sdkRoles.RolesV1Client
}

func NewRoleClient(config config.Config) RoleClient {
	return sdkRoles.NewRolesV1(config)
}
