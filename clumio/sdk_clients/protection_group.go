// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK ProtectionGroupsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkProtectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
)

type ProtectionGroupClient interface {
	sdkProtectionGroups.ProtectionGroupsV1Client
}

func NewProtectionGroupClient(config config.Config) ProtectionGroupClient {
	return sdkProtectionGroups.NewProtectionGroupsV1(config)
}
