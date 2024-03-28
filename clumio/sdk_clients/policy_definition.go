// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK PolicyDefinitionsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkPolicyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
)

type PolicyDefinitionClient interface {
	sdkPolicyDefinitions.PolicyDefinitionsV1Client
}

func NewPolicyDefinitionClient(config config.Config) PolicyDefinitionClient {
	return sdkPolicyDefinitions.NewPolicyDefinitionsV1(config)
}
