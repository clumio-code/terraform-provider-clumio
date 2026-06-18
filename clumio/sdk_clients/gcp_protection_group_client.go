// Copyright (c) 2026 Clumio, a Commvault Company All Rights Reserved

// Contains the wrapper interface for Clumio GO SDK GcpProtectionGroupsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	gcpprotectiongroups "github.com/clumio-code/clumio-go-sdk/controllers/gcp_protection_groups"
)

type GcpProtectionGroupClient interface {
	gcpprotectiongroups.GcpProtectionGroupsV1Client
}

func NewGcpProtectionGroupClient(config config.Config) GcpProtectionGroupClient {
	return gcpprotectiongroups.NewGcpProtectionGroupsV1(config)
}
