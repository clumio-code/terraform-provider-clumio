// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

// Contains the wrapper interface for Clumio GO SDK GcpConnectionsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	gcpconnections "github.com/clumio-code/clumio-go-sdk/controllers/gcp_connections"
)

type GcpConnectionClient interface {
	gcpconnections.GcpConnectionsV1Client
}

func NewGcpConnectionClient(config config.Config) GcpConnectionClient {
	return gcpconnections.NewGcpConnectionsV1(config)
}
