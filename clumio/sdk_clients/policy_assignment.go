// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsConnectionsV1Client.

package sdkclients

import (
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	sdkpolicyassignments "github.com/clumio-code/clumio-go-sdk/controllers/policy_assignments"
)

type PolicyAssignmentClient interface {
	sdkpolicyassignments.PolicyAssignmentsV1Client
}

func NewPolicyAssignmentClient(config sdkconfig.Config) PolicyAssignmentClient {
	return sdkpolicyassignments.NewPolicyAssignmentsV1(config)
}
